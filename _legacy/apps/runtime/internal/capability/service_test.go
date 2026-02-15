package capability

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pryx-core/internal/bus"
)

func setupTestService(t *testing.T) *Service {
	t.Helper()
	b := bus.New()
	return NewService(b)
}

func TestNewService(t *testing.T) {
	svc := setupTestService(t)

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.advertisements)
	assert.NotNil(t, svc.negotiations)
	assert.NotNil(t, svc.responses)
	assert.NotNil(t, svc.grants)
}

func TestAdvertiseCapabilities(t *testing.T) {
	svc := setupTestService(t)

	advertisement := &CapabilityAdvertisement{
		AgentID:         "agent-001",
		AgentName:       "TestAgent",
		AgentVersion:    "1.0.0",
		ProtocolVersion: "1.0",
		Capabilities: []Capability{
			{
				ID:          "cap-001",
				Name:        "filesystem.read",
				Type:        CapabilityTypeTool,
				Description: "Read files from disk",
				Version:     "1.0",
				Parameters: CapabilityParameters{
					Timeout:   30 * time.Second,
					Streaming: false,
				},
			},
		},
	}

	err := svc.AdvertiseCapabilities(context.Background(), advertisement)
	require.NoError(t, err)

	// Verify advertisement was stored
	retrieved, err := svc.GetAdvertisement("agent-001")
	require.NoError(t, err)
	assert.Equal(t, "agent-001", retrieved.AgentID)
	assert.Equal(t, "TestAgent", retrieved.AgentName)
	assert.Len(t, retrieved.Capabilities, 1)
}

func TestAdvertiseCapabilities_GeneratesID(t *testing.T) {
	svc := setupTestService(t)

	advertisement := &CapabilityAdvertisement{
		AgentID:   "agent-002",
		AgentName: "TestAgent2",
		Capabilities: []Capability{
			{
				Name: "test-capability",
				Type: CapabilityTypeSkill,
			},
		},
	}

	err := svc.AdvertiseCapabilities(context.Background(), advertisement)
	require.NoError(t, err)

	// Verify capability ID was generated
	retrieved, err := svc.GetAdvertisement("agent-002")
	require.NoError(t, err)
	assert.NotEmpty(t, retrieved.Capabilities[0].ID)
}

func TestGetAdvertisement_NotFound(t *testing.T) {
	svc := setupTestService(t)

	_, err := svc.GetAdvertisement("non-existent-agent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no capability advertisement found")
}

func TestGetAdvertisement_Expired(t *testing.T) {
	svc := setupTestService(t)

	advertisement := &CapabilityAdvertisement{
		AgentID:   "agent-001",
		AgentName: "TestAgent",
		ExpiresAt: func() *time.Time {
			t := time.Now().Add(-1 * time.Hour)
			return &t
		}(),
	}

	err := svc.AdvertiseCapabilities(context.Background(), advertisement)
	require.NoError(t, err)

	_, err = svc.GetAdvertisement("agent-001")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}

func TestGetCapabilitiesByType(t *testing.T) {
	svc := setupTestService(t)

	advertisement := &CapabilityAdvertisement{
		AgentID:   "agent-001",
		AgentName: "TestAgent",
		Capabilities: []Capability{
			{Name: "tool-1", Type: CapabilityTypeTool},
			{Name: "tool-2", Type: CapabilityTypeTool},
			{Name: "skill-1", Type: CapabilityTypeSkill},
			{Name: "model-1", Type: CapabilityTypeModel},
		},
	}

	err := svc.AdvertiseCapabilities(context.Background(), advertisement)
	require.NoError(t, err)

	// Get tools
	tools, err := svc.GetCapabilitiesByType("agent-001", CapabilityTypeTool)
	require.NoError(t, err)
	assert.Len(t, tools, 2)

	// Get skills
	skills, err := svc.GetCapabilitiesByType("agent-001", CapabilityTypeSkill)
	require.NoError(t, err)
	assert.Len(t, skills, 1)

	// Get models
	models, err := svc.GetCapabilitiesByType("agent-001", CapabilityTypeModel)
	require.NoError(t, err)
	assert.Len(t, models, 1)
}

func TestCheckCompatibility(t *testing.T) {
	svc := setupTestService(t)

	// Create two advertisements
	adv1 := &CapabilityAdvertisement{
		AgentID:         "agent-001",
		AgentName:       "Agent1",
		ProtocolVersion: "1.0",
		Capabilities: []Capability{
			{
				Name: "tool-1",
				Type: CapabilityTypeTool,
			},
		},
	}

	adv2 := &CapabilityAdvertisement{
		AgentID:         "agent-002",
		AgentName:       "Agent2",
		ProtocolVersion: "1.0",
		Capabilities: []Capability{
			{
				Name: "skill-1",
				Type: CapabilityTypeSkill,
			},
		},
	}

	err := svc.AdvertiseCapabilities(context.Background(), adv1)
	require.NoError(t, err)

	err = svc.AdvertiseCapabilities(context.Background(), adv2)
	require.NoError(t, err)

	// Check compatibility
	result, err := svc.CheckCompatibility("agent-001", "agent-002")
	require.NoError(t, err)

	assert.True(t, result.Compatible)
	assert.Greater(t, result.Score, 0.0)
}

func TestCheckCompatibility_ProtocolMismatch(t *testing.T) {
	svc := setupTestService(t)

	adv1 := &CapabilityAdvertisement{
		AgentID:         "agent-001",
		AgentName:       "Agent1",
		ProtocolVersion: "1.0",
	}

	adv2 := &CapabilityAdvertisement{
		AgentID:         "agent-002",
		AgentName:       "Agent2",
		ProtocolVersion: "2.0",
	}

	err := svc.AdvertiseCapabilities(context.Background(), adv1)
	require.NoError(t, err)

	err = svc.AdvertiseCapabilities(context.Background(), adv2)
	require.NoError(t, err)

	result, err := svc.CheckCompatibility("agent-001", "agent-002")
	require.NoError(t, err)

	// Should still be compatible but with warnings
	assert.True(t, result.Compatible)
	assert.Contains(t, result.Issues[0].Code, "PROTOCOL_VERSION_MISMATCH")
}

func TestCheckCompatibility_AgentNotFound(t *testing.T) {
	svc := setupTestService(t)

	_, err := svc.CheckCompatibility("non-existent-1", "non-existent-2")
	assert.Error(t, err)
}

func TestRequestNegotiation(t *testing.T) {
	svc := setupTestService(t)

	req := &NegotiationRequest{
		FromAgentID:   "agent-001",
		FromAgentName: "TestAgent",
		RequestedCaps: []CapabilityRequest{
			{
				CapabilityID:   "cap-001",
				CapabilityName: "filesystem.read",
			},
		},
		IntendedUse: "Reading configuration files",
	}

	response, err := svc.RequestNegotiation(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, response)

	assert.NotEmpty(t, response.RequestID)
	assert.Equal(t, NegotiationStatusPending, response.Status)
}

func TestRequestNegotiation_DefaultDuration(t *testing.T) {
	svc := setupTestService(t)

	req := &NegotiationRequest{
		FromAgentID: "agent-001",
	}

	response, err := svc.RequestNegotiation(context.Background(), req)
	require.NoError(t, err)

	// Default duration should be 24 hours
	assert.NotNil(t, response)
}

func TestApproveNegotiation(t *testing.T) {
	svc := setupTestService(t)

	// Request negotiation
	req := &NegotiationRequest{
		FromAgentID:   "agent-001",
		FromAgentName: "TestAgent",
		RequestedCaps: []CapabilityRequest{
			{CapabilityID: "cap-001"},
		},
	}

	requestResponse, err := svc.RequestNegotiation(context.Background(), req)
	require.NoError(t, err)

	// Approve negotiation
	grantedCaps := []GrantedCapability{
		{
			CapabilityID:    "cap-001",
			CapabilityName:  "filesystem.read",
			PermissionLevel: PermissionLevelRead,
		},
	}

	conditions := []string{"Only access config files"}

	response, err := svc.ApproveNegotiation(context.Background(), requestResponse.RequestID, grantedCaps, conditions)
	require.NoError(t, err)
	require.NotNil(t, response)

	assert.Equal(t, NegotiationStatusApproved, response.Status)
	assert.Len(t, response.GrantedCaps, 1)
	assert.Equal(t, PermissionLevelRead, response.GrantedCaps[0].PermissionLevel)
	assert.Len(t, response.Conditions, 1)

	// Verify grant was stored
	grants := svc.GetGrantedCapabilities("agent-001")
	assert.Len(t, grants, 1)
}

func TestDenyNegotiation(t *testing.T) {
	svc := setupTestService(t)

	// Request negotiation
	req := &NegotiationRequest{
		FromAgentID:   "agent-001",
		FromAgentName: "TestAgent",
	}

	requestResponse, err := svc.RequestNegotiation(context.Background(), req)
	require.NoError(t, err)

	// Deny negotiation
	deniedCaps := []DeniedCapability{
		{
			CapabilityID:   "cap-001",
			CapabilityName: "dangerous.tool",
			Reason:         "Too risky",
			CanAppeal:      true,
		},
	}

	response, err := svc.DenyNegotiation(context.Background(), requestResponse.RequestID, deniedCaps, "Request denied for security reasons")
	require.NoError(t, err)
	require.NotNil(t, response)

	assert.Equal(t, NegotiationStatusDenied, response.Status)
	assert.Len(t, response.DeniedCaps, 1)
}

func TestGetNegotiationStatus_NotFound(t *testing.T) {
	svc := setupTestService(t)

	_, err := svc.GetNegotiationStatus("non-existent-request")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "negotiation response not found")
}

func TestRevokeCapabilities(t *testing.T) {
	svc := setupTestService(t)

	// Approve a negotiation first
	req := &NegotiationRequest{
		FromAgentID: "agent-001",
	}

	requestResponse, err := svc.RequestNegotiation(context.Background(), req)
	require.NoError(t, err)

	grantedCaps := []GrantedCapability{
		{CapabilityID: "cap-001"},
		{CapabilityID: "cap-002"},
	}

	_, err = svc.ApproveNegotiation(context.Background(), requestResponse.RequestID, grantedCaps, nil)
	require.NoError(t, err)

	// Verify grants exist
	grants := svc.GetGrantedCapabilities("agent-001")
	assert.Len(t, grants, 2)

	// Revoke one capability
	err = svc.RevokeCapabilities("agent-001", []string{"cap-001"})
	require.NoError(t, err)

	// Verify one grant remains
	grants = svc.GetGrantedCapabilities("agent-001")
	assert.Len(t, grants, 1)
	assert.Equal(t, "cap-002", grants[0].CapabilityID)
}

func TestDiscoverCapabilities(t *testing.T) {
	svc := setupTestService(t)

	// Create advertisements with different capabilities
	advertisements := []*CapabilityAdvertisement{
		{
			AgentID:   "agent-001",
			AgentName: "Agent1",
			Capabilities: []Capability{
				{Name: "filesystem.read", Type: CapabilityTypeTool},
			},
		},
		{
			AgentID:   "agent-002",
			AgentName: "Agent2",
			Capabilities: []Capability{
				{Name: "http.request", Type: CapabilityTypeTool},
				{Name: "database.query", Type: CapabilityTypeTool},
			},
		},
		{
			AgentID:   "agent-003",
			AgentName: "Agent3",
			Capabilities: []Capability{
				{Name: "analysis.skill", Type: CapabilityTypeSkill},
			},
		},
	}

	for _, adv := range advertisements {
		err := svc.AdvertiseCapabilities(context.Background(), adv)
		require.NoError(t, err)
	}

	// Discover all tools
	tools := svc.DiscoverCapabilities(map[string]interface{}{
		"type": string(CapabilityTypeTool),
	})
	assert.Len(t, tools, 3)

	// Discover by name
	specific := svc.DiscoverCapabilities(map[string]interface{}{
		"name": "filesystem.read",
	})
	assert.Len(t, specific, 1)
}

func TestGetAdvertisementCount(t *testing.T) {
	svc := setupTestService(t)

	assert.Equal(t, 0, svc.GetAdvertisementCount())

	// Add advertisements
	adv1 := &CapabilityAdvertisement{AgentID: "agent-001", AgentName: "Agent1"}
	adv2 := &CapabilityAdvertisement{AgentID: "agent-002", AgentName: "Agent2"}

	err := svc.AdvertiseCapabilities(context.Background(), adv1)
	require.NoError(t, err)
	assert.Equal(t, 1, svc.GetAdvertisementCount())

	err = svc.AdvertiseCapabilities(context.Background(), adv2)
	require.NoError(t, err)
	assert.Equal(t, 2, svc.GetAdvertisementCount())
}

func TestGetNegotiationCount(t *testing.T) {
	svc := setupTestService(t)

	assert.Equal(t, 0, svc.GetNegotiationCount())

	// Request negotiations
	req1 := &NegotiationRequest{FromAgentID: "agent-001"}
	req2 := &NegotiationRequest{FromAgentID: "agent-002"}

	_, err := svc.RequestNegotiation(context.Background(), req1)
	require.NoError(t, err)
	assert.Equal(t, 1, svc.GetNegotiationCount())

	_, err = svc.RequestNegotiation(context.Background(), req2)
	require.NoError(t, err)
	assert.Equal(t, 2, svc.GetNegotiationCount())
}

func TestPrettyPrint(t *testing.T) {
	_ = setupTestService(t)

	cap := &Capability{
		Name:    "test.capability",
		Type:    CapabilityTypeTool,
		Version: "1.0.0",
	}

	output := cap.PrettyPrint()
	assert.Contains(t, output, "test.capability")
	assert.Contains(t, output, "tool")
	assert.Contains(t, output, "1.0.0")

	adv := &CapabilityAdvertisement{
		AgentID:         "agent-001",
		AgentName:       "TestAgent",
		ProtocolVersion: "1.0",
		Capabilities:    []Capability{{Name: "cap1"}},
	}

	output = adv.PrettyPrint()
	assert.Contains(t, output, "TestAgent")
	assert.Contains(t, output, "agent-001")
	assert.Contains(t, output, "1")
}

func TestNegotiationExpiry(t *testing.T) {
	svc := setupTestService(t)

	req := &NegotiationRequest{
		FromAgentID: "agent-001",
		Duration:    1 * time.Hour,
	}

	response, err := svc.RequestNegotiation(context.Background(), req)
	require.NoError(t, err)

	// Approve with 1 hour duration
	grantedCaps := []GrantedCapability{
		{CapabilityID: "cap-001"},
	}

	approveResp, err := svc.ApproveNegotiation(context.Background(), response.RequestID, grantedCaps, nil)
	require.NoError(t, err)

	// Verify expiry is set
	assert.NotNil(t, approveResp.ExpiresAt)
	assert.True(t, approveResp.ExpiresAt.After(time.Now()))
}

func TestConcurrentNegotiations(t *testing.T) {
	svc := setupTestService(t)

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(agentID string) {
			req := &NegotiationRequest{
				FromAgentID: agentID,
			}
			_, _ = svc.RequestNegotiation(context.Background(), req)
			done <- true
		}("agent-" + string(rune('A'+i)))
	}

	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for concurrent negotiations")
		}
	}

	assert.Equal(t, 10, svc.GetNegotiationCount())
}
