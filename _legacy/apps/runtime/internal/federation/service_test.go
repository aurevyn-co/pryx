package federation

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
	assert.NotNil(t, svc.tools)
	assert.NotNil(t, svc.invocations)
	assert.NotNil(t, svc.federations)
	assert.NotNil(t, svc.connections)
	assert.NotNil(t, svc.grants)
}

func TestRegisterTool(t *testing.T) {
	svc := setupTestService(t)

	tool := &ToolDefinition{
		ToolID:       "",
		Name:         "filesystem.read",
		Description:  "Read files from disk",
		Version:      "1.0",
		Visibility:   ToolVisibilityPublic,
		OwnerAgentID: "agent-001",
	}

	err := svc.RegisterTool(tool)
	require.NoError(t, err)

	assert.NotEmpty(t, tool.ToolID)
	assert.NotZero(t, tool.CreatedAt)
}

func TestRegisterTool_GeneratesID(t *testing.T) {
	svc := setupTestService(t)

	tool := &ToolDefinition{
		Name:         "test.tool",
		OwnerAgentID: "agent-001",
	}

	err := svc.RegisterTool(tool)
	require.NoError(t, err)

	assert.NotEmpty(t, tool.ToolID)
}

func TestGetTool(t *testing.T) {
	svc := setupTestService(t)

	tool := &ToolDefinition{
		ToolID:       "tool-001",
		Name:         "test.tool",
		OwnerAgentID: "agent-001",
	}

	err := svc.RegisterTool(tool)
	require.NoError(t, err)

	retrieved, err := svc.GetTool("tool-001")
	require.NoError(t, err)
	assert.Equal(t, "tool-001", retrieved.ToolID)
	assert.Equal(t, "test.tool", retrieved.Name)
}

func TestGetTool_NotFound(t *testing.T) {
	svc := setupTestService(t)

	_, err := svc.GetTool("non-existent-tool")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tool not found")
}

func TestGetToolsByAgent(t *testing.T) {
	svc := setupTestService(t)

	// Register tools for different agents
	tools := []*ToolDefinition{
		{ToolID: "tool-001", Name: "tool1", OwnerAgentID: "agent-001"},
		{ToolID: "tool-002", Name: "tool2", OwnerAgentID: "agent-001"},
		{ToolID: "tool-003", Name: "tool3", OwnerAgentID: "agent-002"},
	}

	for _, tool := range tools {
		err := svc.RegisterTool(tool)
		require.NoError(t, err)
	}

	agentTools := svc.GetToolsByAgent("agent-001")
	assert.Len(t, agentTools, 2)

	agentTools = svc.GetToolsByAgent("agent-002")
	assert.Len(t, agentTools, 1)
}

func TestGetVisibleTools(t *testing.T) {
	svc := setupTestService(t)

	tools := []*ToolDefinition{
		{ToolID: "tool-001", Name: "public", Visibility: ToolVisibilityPublic, OwnerAgentID: "agent-001"},
		{ToolID: "tool-002", Name: "private", Visibility: ToolVisibilityPrivate, OwnerAgentID: "agent-001"},
		{ToolID: "tool-003", Name: "federated", Visibility: ToolVisibilityFederated, OwnerAgentID: "agent-001"},
		{ToolID: "tool-004", Name: "trusted", Visibility: ToolVisibilityTrusted, OwnerAgentID: "agent-001"},
	}

	for _, tool := range tools {
		err := svc.RegisterTool(tool)
		require.NoError(t, err)
	}

	// Get visible tools for agent-001 (owner)
	visible := svc.GetVisibleTools("agent-001", "low")
	assert.Len(t, visible, 3) // Owner sees all their tools (public, federated, private; trusted requires high trust)

	// Get visible tools for external agent with low trust
	visible = svc.GetVisibleTools("agent-002", "low")
	assert.Len(t, visible, 2) // Public + federated

	// Get visible tools for external agent with high trust
	visible = svc.GetVisibleTools("agent-002", "high")
	assert.Len(t, visible, 3) // Public + federated + trusted
}

func TestRequestFederation(t *testing.T) {
	svc := setupTestService(t)

	req := &FederationRequest{
		FromAgentID:    "agent-001",
		FromAgentName:  "Agent1",
		ToAgentID:      "agent-002",
		ToAgentName:    "Agent2",
		RequestedTools: []string{"tool-001", "tool-002"},
		IntendedUse:    "Data processing",
	}

	response, err := svc.RequestFederation(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, response)

	assert.NotEmpty(t, response.RequestID)
	assert.Equal(t, FederationStatusPending, response.Status)
}

func TestRequestFederation_DefaultDuration(t *testing.T) {
	svc := setupTestService(t)

	req := &FederationRequest{
		FromAgentID: "agent-001",
		ToAgentID:   "agent-002",
	}

	response, err := svc.RequestFederation(context.Background(), req)
	require.NoError(t, err)

	assert.NotNil(t, response)
}

func TestApproveFederation(t *testing.T) {
	svc := setupTestService(t)

	// Request federation
	req := &FederationRequest{
		FromAgentID:   "agent-001",
		FromAgentName: "Agent1",
		ToAgentID:     "agent-002",
		ToAgentName:   "Agent2",
	}

	requestResponse, err := svc.RequestFederation(context.Background(), req)
	require.NoError(t, err)

	// Approve federation
	grantedTools := []GrantedTool{
		{ToolID: "tool-001", ToolName: "filesystem.read"},
	}

	response, err := svc.ApproveFederation(context.Background(), requestResponse.RequestID, grantedTools, "Approved for data processing")
	require.NoError(t, err)
	require.NotNil(t, response)

	assert.Equal(t, FederationStatusApproved, response.Status)
	assert.Len(t, response.GrantedTools, 1)

	// Verify connection was created
	connections := svc.GetConnectionsByAgent("agent-001")
	assert.Len(t, connections, 1)
}

func TestApproveFederation_NotFound(t *testing.T) {
	svc := setupTestService(t)

	_, err := svc.ApproveFederation(context.Background(), "non-existent-request", nil, "message")
	assert.Error(t, err)
}

func TestInvokeTool(t *testing.T) {
	svc := setupTestService(t)

	// Register a tool
	tool := &ToolDefinition{
		ToolID:       "tool-001",
		Name:         "test.tool",
		OwnerAgentID: "agent-001",
	}

	err := svc.RegisterTool(tool)
	require.NoError(t, err)

	// Invoke the tool
	invocation, err := svc.InvokeTool(context.Background(), "tool-001", "agent-002", map[string]interface{}{
		"param": "value",
	})
	require.NoError(t, err)
	require.NotNil(t, invocation)

	assert.NotEmpty(t, invocation.InvocationID)
	assert.Equal(t, "tool-001", invocation.ToolID)
	assert.Equal(t, "agent-002", invocation.CallerAgentID)
	assert.Equal(t, InvocationStatusCompleted, invocation.Status)
}

func TestInvokeTool_NotFound(t *testing.T) {
	svc := setupTestService(t)

	_, err := svc.InvokeTool(context.Background(), "non-existent-tool", "agent-001", nil)
	assert.Error(t, err)
}

func TestGetInvocationStatus(t *testing.T) {
	svc := setupTestService(t)

	// Register and invoke a tool
	tool := &ToolDefinition{
		ToolID:       "tool-001",
		Name:         "test.tool",
		OwnerAgentID: "agent-001",
	}

	err := svc.RegisterTool(tool)
	require.NoError(t, err)

	invocation, err := svc.InvokeTool(context.Background(), "tool-001", "agent-002", nil)
	require.NoError(t, err)

	// Get status
	status, err := svc.GetInvocationStatus(invocation.InvocationID)
	require.NoError(t, err)
	assert.Equal(t, InvocationStatusCompleted, status.Status)
}

func TestGetInvocationStatus_NotFound(t *testing.T) {
	svc := setupTestService(t)

	_, err := svc.GetInvocationStatus("non-existent-invocation")
	assert.Error(t, err)
}

func TestGetFederationConnection(t *testing.T) {
	svc := setupTestService(t)

	// Create and approve federation
	req := &FederationRequest{
		FromAgentID: "agent-001",
		ToAgentID:   "agent-002",
	}

	requestResponse, err := svc.RequestFederation(context.Background(), req)
	require.NoError(t, err)

	_, err = svc.ApproveFederation(context.Background(), requestResponse.RequestID, []GrantedTool{{ToolID: "tool-001"}}, "Approved")
	require.NoError(t, err)

	// Get connection
	connections := svc.GetConnectionsByAgent("agent-001")
	require.Len(t, connections, 1)

	retrieved, err := svc.GetFederationConnection(connections[0].ConnectionID)
	require.NoError(t, err)
	assert.Equal(t, connections[0].ConnectionID, retrieved.ConnectionID)
}

func TestGetConnectionsByAgent(t *testing.T) {
	svc := setupTestService(t)

	// Create multiple federations
	for i := 0; i < 3; i++ {
		req := &FederationRequest{
			FromAgentID: "agent-001",
			ToAgentID:   "agent-00" + string(rune('2'+i)),
		}

		resp, err := svc.RequestFederation(context.Background(), req)
		require.NoError(t, err)

		_, err = svc.ApproveFederation(context.Background(), resp.RequestID, nil, "Approved")
		require.NoError(t, err)
	}

	connections := svc.GetConnectionsByAgent("agent-001")
	assert.Len(t, connections, 3)
}

func TestRevokeTool(t *testing.T) {
	svc := setupTestService(t)

	// Approve federation with tools
	req := &FederationRequest{
		FromAgentID: "agent-001",
		ToAgentID:   "agent-002",
	}

	resp, err := svc.RequestFederation(context.Background(), req)
	require.NoError(t, err)

	_, err = svc.ApproveFederation(context.Background(), resp.RequestID, []GrantedTool{
		{ToolID: "tool-001"},
		{ToolID: "tool-002"},
	}, "Approved")
	require.NoError(t, err)

	// Verify grants
	grants := svc.GetGrantedTools("agent-001")
	assert.Len(t, grants, 2)

	// Revoke one tool
	err = svc.RevokeTool("agent-001", "tool-001")
	require.NoError(t, err)

	// Verify one grant remains
	grants = svc.GetGrantedTools("agent-001")
	assert.Len(t, grants, 1)
}

func TestGetGrantedTools(t *testing.T) {
	svc := setupTestService(t)

	// Approve federation
	req := &FederationRequest{
		FromAgentID: "agent-001",
		ToAgentID:   "agent-002",
	}

	resp, err := svc.RequestFederation(context.Background(), req)
	require.NoError(t, err)

	_, err = svc.ApproveFederation(context.Background(), resp.RequestID, []GrantedTool{
		{ToolID: "tool-001"},
	}, "Approved")
	require.NoError(t, err)

	grants := svc.GetGrantedTools("agent-001")
	assert.Len(t, grants, 1)
}

func TestSyncTools(t *testing.T) {
	svc := setupTestService(t)

	// Create federation
	req := &FederationRequest{
		FromAgentID: "agent-001",
		ToAgentID:   "agent-002",
	}

	_, err := svc.RequestFederation(context.Background(), req)
	require.NoError(t, err)

	_, err = svc.ApproveFederation(context.Background(), req.RequestID, nil, "Approved")
	require.NoError(t, err)

	// Register tools on agent-002 (the remote agent)
	tool := &ToolDefinition{
		ToolID:       "remote-tool-001",
		Name:         "remote.tool",
		OwnerAgentID: "agent-002",
		Visibility:   ToolVisibilityFederated,
	}

	err = svc.RegisterTool(tool)
	require.NoError(t, err)

	// Get the connection ID
	connections := svc.GetConnectionsByAgent("agent-002")
	require.Len(t, connections, 1)
	connectionID := connections[0].ConnectionID

	// Sync tools
	tools, err := svc.SyncTools(connectionID)
	require.NoError(t, err)
	assert.Len(t, tools, 1)
}

func TestGetToolCount(t *testing.T) {
	svc := setupTestService(t)

	assert.Equal(t, 0, svc.GetToolCount())

	for i := 0; i < 5; i++ {
		tool := &ToolDefinition{
			ToolID:       "tool-" + string(rune('0'+i)),
			OwnerAgentID: "agent-001",
		}
		err := svc.RegisterTool(tool)
		require.NoError(t, err)
	}

	assert.Equal(t, 5, svc.GetToolCount())
}

func TestGetConnectionCount(t *testing.T) {
	svc := setupTestService(t)

	assert.Equal(t, 0, svc.GetConnectionCount())

	for i := 0; i < 3; i++ {
		req := &FederationRequest{
			FromAgentID: "agent-001",
			ToAgentID:   "agent-00" + string(rune('2'+i)),
		}

		resp, err := svc.RequestFederation(context.Background(), req)
		require.NoError(t, err)

		_, err = svc.ApproveFederation(context.Background(), resp.RequestID, nil, "Approved")
		require.NoError(t, err)
	}

	assert.Equal(t, 3, svc.GetConnectionCount())
}

func TestGetInvocationCount(t *testing.T) {
	svc := setupTestService(t)

	tool := &ToolDefinition{
		ToolID:       "tool-001",
		OwnerAgentID: "agent-001",
	}

	err := svc.RegisterTool(tool)
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		_, err = svc.InvokeTool(context.Background(), "tool-001", "agent-00"+string(rune('2'+i)), nil)
		require.NoError(t, err)
	}

	assert.Equal(t, 3, svc.GetInvocationCount())
}

func TestPrettyPrint(t *testing.T) {
	_ = setupTestService(t)

	tool := &ToolDefinition{
		ToolID:       "tool-001",
		Name:         "test.tool",
		Version:      "1.0.0",
		Visibility:   ToolVisibilityPublic,
		OwnerAgentID: "agent-001",
	}

	output := tool.PrettyPrint()
	assert.Contains(t, output, "test.tool")
	assert.Contains(t, output, "1.0.0")
	assert.Contains(t, output, "public")
}

func TestConcurrentFederations(t *testing.T) {
	svc := setupTestService(t)

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(agentID string) {
			req := &FederationRequest{
				FromAgentID: agentID,
				ToAgentID:   "agent-dest",
			}
			_, _ = svc.RequestFederation(context.Background(), req)
			done <- true
		}("agent-" + string(rune('A'+i)))
	}

	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for concurrent federations")
		}
	}

	assert.Equal(t, 10, len(svc.federations))
}
