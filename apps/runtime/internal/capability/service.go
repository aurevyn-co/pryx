package capability

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"pryx-core/internal/bus"
)

// CapabilityType represents the type of capability
type CapabilityType string

const (
	CapabilityTypeTool   CapabilityType = "tool"
	CapabilityTypeSkill  CapabilityType = "skill"
	CapabilityTypeModel  CapabilityType = "model"
	CapabilityTypeMemory CapabilityType = "memory"
	CapabilityTypePolicy CapabilityType = "policy"
)

// PermissionLevel represents the level of permission granted
type PermissionLevel string

const (
	PermissionLevelNone    PermissionLevel = "none"
	PermissionLevelRead    PermissionLevel = "read"
	PermissionLevelWrite   PermissionLevel = "write"
	PermissionLevelExecute PermissionLevel = "execute"
	PermissionLevelAdmin   PermissionLevel = "admin"
)

// Capability represents an agent's capability
type Capability struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        CapabilityType         `json:"type"`
	Description string                 `json:"description"`
	Version     string                 `json:"version"`
	Parameters  CapabilityParameters   `json:"parameters"`
	Permissions PermissionRequirements `json:"permissions"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// CapabilityParameters defines input/output parameters for a capability
type CapabilityParameters struct {
	InputSchema  json.RawMessage `json:"input_schema"`
	OutputSchema json.RawMessage `json:"output_schema"`
	Timeout      time.Duration   `json:"timeout"`
	Streaming    bool            `json:"streaming"`
}

// PermissionRequirements defines what permissions a capability needs
type PermissionRequirements struct {
	Required   []PermissionRequest `json:"required"`
	Optional   []PermissionRequest `json:"optional"`
	TrustLevel string              `json:"trust_level"`
	RateLimit  *RateLimit          `json:"rate_limit,omitempty"`
}

// PermissionRequest represents a permission request
type PermissionRequest struct {
	Resource    string          `json:"resource"`
	Action      string          `json:"action"`
	Level       PermissionLevel `json:"level"`
	Description string          `json:"description"`
}

// RateLimit defines rate limiting for a capability
type RateLimit struct {
	RequestsPerMinute int `json:"requests_per_minute"`
	RequestsPerHour   int `json:"requests_per_hour"`
	RequestsPerDay    int `json:"requests_per_day"`
	BurstSize         int `json:"burst_size"`
}

// CapabilityAdvertisement represents an agent's capability advertisement
type CapabilityAdvertisement struct {
	AgentID         string                 `json:"agent_id"`
	AgentName       string                 `json:"agent_name"`
	AgentVersion    string                 `json:"agent_version"`
	Capabilities    []Capability           `json:"capabilities"`
	ProtocolVersion string                 `json:"protocol_version"`
	ExpiresAt       *time.Time             `json:"expires_at,omitempty"`
	Metadata        map[string]interface{} `json:"metadata"`
	CreatedAt       time.Time              `json:"created_at"`
}

// NegotiationRequest represents a capability negotiation request
type NegotiationRequest struct {
	RequestID     string              `json:"request_id"`
	FromAgentID   string              `json:"from_agent_id"`
	FromAgentName string              `json:"from_agent_name"`
	RequestedCaps []CapabilityRequest `json:"requested_caps"`
	IntendedUse   string              `json:"intended_use"`
	Duration      time.Duration       `json:"duration"`
	CreatedAt     time.Time           `json:"created_at"`
}

// CapabilityRequest represents a request for a specific capability
type CapabilityRequest struct {
	CapabilityID   string                 `json:"capability_id"`
	CapabilityName string                 `json:"capability_name"`
	Parameters     map[string]interface{} `json:"parameters"`
	Timeout        time.Duration          `json:"timeout"`
}

// NegotiationResponse represents a negotiation response
type NegotiationResponse struct {
	RequestID   string              `json:"request_id"`
	Status      NegotiationStatus   `json:"status"`
	GrantedCaps []GrantedCapability `json:"granted_caps,omitempty"`
	DeniedCaps  []DeniedCapability  `json:"denied_caps,omitempty"`
	Conditions  []string            `json:"conditions,omitempty"`
	Message     string              `json:"message"`
	ExpiresAt   *time.Time          `json:"expires_at,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
}

// NegotiationStatus represents the status of a negotiation
type NegotiationStatus string

const (
	NegotiationStatusPending  NegotiationStatus = "pending"
	NegotiationStatusApproved NegotiationStatus = "approved"
	NegotiationStatusPartial  NegotiationStatus = "partial"
	NegotiationStatusDenied   NegotiationStatus = "denied"
	NegotiationStatusExpired  NegotiationStatus = "expired"
)

// GrantedCapability represents a granted capability
type GrantedCapability struct {
	CapabilityID    string          `json:"capability_id"`
	CapabilityName  string          `json:"capability_name"`
	PermissionLevel PermissionLevel `json:"permission_level"`
	ExpiresAt       *time.Time      `json:"expires_at,omitempty"`
	Conditions      []string        `json:"conditions,omitempty"`
}

// DeniedCapability represents a denied capability
type DeniedCapability struct {
	CapabilityID   string `json:"capability_id"`
	CapabilityName string `json:"capability_name"`
	Reason         string `json:"reason"`
	CanAppeal      bool   `json:"can_appeal"`
}

// CompatibilityResult represents the result of a compatibility check
type CompatibilityResult struct {
	Compatible      bool                 `json:"compatible"`
	Score           float64              `json:"score"`
	Issues          []CompatibilityIssue `json:"issues"`
	Warnings        []string             `json:"warnings"`
	Recommendations []string             `json:"recommendations"`
}

// CompatibilityIssue represents a compatibility issue
type CompatibilityIssue struct {
	Severity string `json:"severity"` // error, warning, info
	Code     string `json:"code"`
	Message  string `json:"message"`
	Field    string `json:"field,omitempty"`
}

// Service manages capability advertisement and negotiation
type Service struct {
	mu             sync.RWMutex
	bus            *bus.Bus
	advertisements map[string]*CapabilityAdvertisement
	negotiations   map[string]*NegotiationRequest
	responses      map[string]*NegotiationResponse
	grants         map[string][]GrantedCapability
}

// NewService creates a new capability service
func NewService(b *bus.Bus) *Service {
	return &Service{
		bus:            b,
		advertisements: make(map[string]*CapabilityAdvertisement),
		negotiations:   make(map[string]*NegotiationRequest),
		responses:      make(map[string]*NegotiationResponse),
		grants:         make(map[string][]GrantedCapability),
	}
}

// AdvertiseCapabilities registers an agent's capabilities
func (s *Service) AdvertiseCapabilities(ctx context.Context, advertisement *CapabilityAdvertisement) error {
	advertisement.CreatedAt = time.Now().UTC()
	advertisement.CreatedAt = time.Now().UTC()

	if advertisement.ProtocolVersion == "" {
		advertisement.ProtocolVersion = "1.0"
	}

	// Validate capabilities
	for i := range advertisement.Capabilities {
		cap := &advertisement.Capabilities[i]
		if cap.ID == "" {
			cap.ID = uuid.New().String()
		}
		cap.CreatedAt = time.Now().UTC()
		cap.UpdatedAt = time.Now().UTC()
	}

	s.mu.Lock()
	s.advertisements[advertisement.AgentID] = advertisement
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("capability.advertised", "", map[string]interface{}{
		"agent_id":     advertisement.AgentID,
		"agent_name":   advertisement.AgentName,
		"capabilities": len(advertisement.Capabilities),
	}))

	return nil
}

// GetAdvertisement retrieves an agent's capability advertisement
func (s *Service) GetAdvertisement(agentID string) (*CapabilityAdvertisement, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	adv, exists := s.advertisements[agentID]
	if !exists {
		return nil, fmt.Errorf("no capability advertisement found for agent: %s", agentID)
	}

	// Check expiration
	if adv.ExpiresAt != nil && adv.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("capability advertisement expired for agent: %s", agentID)
	}

	return adv, nil
}

// GetCapabilitiesByType retrieves all capabilities of a specific type from an agent
func (s *Service) GetCapabilitiesByType(agentID string, capType CapabilityType) ([]Capability, error) {
	adv, err := s.GetAdvertisement(agentID)
	if err != nil {
		return nil, err
	}

	capabilities := make([]Capability, 0)
	for _, cap := range adv.Capabilities {
		if cap.Type == capType {
			capabilities = append(capabilities, cap)
		}
	}

	return capabilities, nil
}

// CheckCompatibility checks if two agents are compatible
func (s *Service) CheckCompatibility(agentID1, agentID2 string) (*CompatibilityResult, error) {
	result := &CompatibilityResult{
		Compatible:      true,
		Score:           1.0,
		Issues:          make([]CompatibilityIssue, 0),
		Warnings:        make([]string, 0),
		Recommendations: make([]string, 0),
	}

	adv1, err := s.GetAdvertisement(agentID1)
	if err != nil {
		return nil, fmt.Errorf("failed to get advertisement for agent %s: %w", agentID1, err)
	}

	adv2, err := s.GetAdvertisement(agentID2)
	if err != nil {
		return nil, fmt.Errorf("failed to get advertisement for agent %s: %w", agentID2, err)
	}

	// Check protocol version compatibility
	if adv1.ProtocolVersion != adv2.ProtocolVersion {
		result.Issues = append(result.Issues, CompatibilityIssue{
			Severity: "warning",
			Code:     "PROTOCOL_VERSION_MISMATCH",
			Message:  fmt.Sprintf("Protocol version mismatch: %s vs %s", adv1.ProtocolVersion, adv2.ProtocolVersion),
			Field:    "protocol_version",
		})
		result.Score -= 0.1
		result.Recommendations = append(result.Recommendations, "Consider upgrading to the same protocol version")
	}

	// Check for overlapping capabilities
	cap1Types := make(map[CapabilityType]int)
	cap2Types := make(map[CapabilityType]int)

	for _, cap := range adv1.Capabilities {
		cap1Types[cap.Type]++
	}

	for _, cap := range adv2.Capabilities {
		cap2Types[cap.Type]++
	}

	// Look for complementary capabilities
	for capType, count1 := range cap1Types {
		if count2, exists := cap2Types[capType]; exists {
			if count1 > 0 && count2 > 0 {
				result.Warnings = append(result.Warnings, fmt.Sprintf("Both agents have %s capabilities - potential overlap", capType))
			}
		}
	}

	// Check capability requirements
	for _, cap := range adv1.Capabilities {
		if len(cap.Permissions.Required) > 0 {
			trustLevel := cap.Permissions.TrustLevel
			if trustLevel != "" && trustLevel != "any" {
				result.Warnings = append(result.Warnings, fmt.Sprintf("Capability %s requires trust level %s", cap.Name, trustLevel))
			}
		}
	}

	return result, nil
}

// RequestNegotiation initiates a capability negotiation
func (s *Service) RequestNegotiation(ctx context.Context, req *NegotiationRequest) (*NegotiationResponse, error) {
	req.RequestID = uuid.New().String()
	req.CreatedAt = time.Now().UTC()

	if req.Duration == 0 {
		req.Duration = 24 * time.Hour
	}

	s.mu.Lock()
	s.negotiations[req.RequestID] = req
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("capability.negotiation.requested", "", map[string]interface{}{
		"request_id":     req.RequestID,
		"from_agent_id":  req.FromAgentID,
		"requested_caps": len(req.RequestedCaps),
	}))

	// Create pending response
	response := &NegotiationResponse{
		RequestID: req.RequestID,
		Status:    NegotiationStatusPending,
		CreatedAt: time.Now().UTC(),
	}

	return response, nil
}

// ApproveNegotiation approves a negotiation request
func (s *Service) ApproveNegotiation(ctx context.Context, requestID string, grantedCaps []GrantedCapability, conditions []string) (*NegotiationResponse, error) {
	s.mu.RLock()
	req, exists := s.negotiations[requestID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("negotiation request not found: %s", requestID)
	}

	expiresAt := time.Now().UTC().Add(req.Duration)

	response := &NegotiationResponse{
		RequestID:   requestID,
		Status:      NegotiationStatusApproved,
		GrantedCaps: grantedCaps,
		Conditions:  conditions,
		ExpiresAt:   &expiresAt,
		CreatedAt:   time.Now().UTC(),
	}

	s.mu.Lock()
	s.responses[requestID] = response
	s.grants[req.FromAgentID] = append(s.grants[req.FromAgentID], grantedCaps...)
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("capability.negotiation.approved", "", map[string]interface{}{
		"request_id":   requestID,
		"from_agent":   req.FromAgentID,
		"granted_caps": len(grantedCaps),
	}))

	return response, nil
}

// DenyNegotiation denies a negotiation request
func (s *Service) DenyNegotiation(ctx context.Context, requestID string, deniedCaps []DeniedCapability, message string) (*NegotiationResponse, error) {
	s.mu.RLock()
	req, exists := s.negotiations[requestID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("negotiation request not found: %s", requestID)
	}

	response := &NegotiationResponse{
		RequestID:  requestID,
		Status:     NegotiationStatusDenied,
		DeniedCaps: deniedCaps,
		Message:    message,
		CreatedAt:  time.Now().UTC(),
	}

	s.mu.Lock()
	s.responses[requestID] = response
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("capability.negotiation.denied", "", map[string]interface{}{
		"request_id":  requestID,
		"from_agent":  req.FromAgentID,
		"denied_caps": len(deniedCaps),
	}))

	return response, nil
}

// GetNegotiationStatus gets the status of a negotiation
func (s *Service) GetNegotiationStatus(requestID string) (*NegotiationResponse, error) {
	s.mu.RLock()
	response, exists := s.responses[requestID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("negotiation response not found: %s", requestID)
	}

	return response, nil
}

// GetGrantedCapabilities gets all granted capabilities for an agent
func (s *Service) GetGrantedCapabilities(agentID string) []GrantedCapability {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.grants[agentID]
}

// RevokeCapabilities revokes granted capabilities
func (s *Service) RevokeCapabilities(agentID string, capabilityIDs []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	grants := s.grants[agentID]
	remaining := make([]GrantedCapability, 0)

	for _, grant := range grants {
		found := false
		for _, id := range capabilityIDs {
			if grant.CapabilityID == id {
				found = true
				break
			}
		}
		if !found {
			remaining = append(remaining, grant)
		}
	}

	s.grants[agentID] = remaining

	// Publish event
	s.bus.Publish(bus.NewEvent("capability.revoked", "", map[string]interface{}{
		"agent_id":      agentID,
		"revoked_count": len(capabilityIDs),
	}))

	return nil
}

// DiscoverCapabilities discovers capabilities matching criteria
func (s *Service) DiscoverCapabilities(criteria map[string]interface{}) []Capability {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var matchingCaps []Capability

	for _, adv := range s.advertisements {
		for _, cap := range adv.Capabilities {
			if s.matchesCriteria(&cap, criteria) {
				matchingCaps = append(matchingCaps, cap)
			}
		}
	}

	return matchingCaps
}

// matchesCriteria checks if a capability matches the given criteria
func (s *Service) matchesCriteria(cap *Capability, criteria map[string]interface{}) bool {
	for key, value := range criteria {
		switch key {
		case "type":
			if string(cap.Type) != value.(string) {
				return false
			}
		case "name":
			if cap.Name != value.(string) {
				return false
			}
		case "version":
			if cap.Version != value.(string) {
				return false
			}
		}
	}
	return true
}

// GetAdvertisementCount returns the number of registered advertisements
func (s *Service) GetAdvertisementCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.advertisements)
}

// GetNegotiationCount returns the number of pending negotiations
func (s *Service) GetNegotiationCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.negotiations)
}

// PrettyPrint prints capability info
func (c *Capability) PrettyPrint() string {
	return fmt.Sprintf("Capability{Name: %s, Type: %s, Version: %s}", c.Name, c.Type, c.Version)
}

// PrettyPrint prints advertisement info
func (a *CapabilityAdvertisement) PrettyPrint() string {
	return fmt.Sprintf("Advertisement{Agent: %s (%s), Capabilities: %d, Protocol: %s}",
		a.AgentName, a.AgentID, len(a.Capabilities), a.ProtocolVersion)
}

// MarshalJSON for Capability
func (c *Capability) MarshalJSON() ([]byte, error) {
	type Alias Capability
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}

// MarshalJSON for CapabilityAdvertisement
func (a *CapabilityAdvertisement) MarshalJSON() ([]byte, error) {
	type Alias CapabilityAdvertisement
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	})
}
