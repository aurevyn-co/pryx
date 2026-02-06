package federation

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"pryx-core/internal/bus"
)

// ToolVisibility defines who can see a tool
type ToolVisibility string

const (
	ToolVisibilityPrivate   ToolVisibility = "private"
	ToolVisibilityTrusted   ToolVisibility = "trusted"
	ToolVisibilityFederated ToolVisibility = "federated"
	ToolVisibilityPublic    ToolVisibility = "public"
)

// FederationStatus represents the status of a federation request
type FederationStatus string

const (
	FederationStatusPending   FederationStatus = "pending"
	FederationStatusApproved  FederationStatus = "approved"
	FederationStatusActive    FederationStatus = "active"
	FederationStatusSuspended FederationStatus = "suspended"
	FederationStatusRevoked   FederationStatus = "revoked"
)

// ToolDefinition represents a tool that can be federated
type ToolDefinition struct {
	ToolID       string                 `json:"tool_id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Version      string                 `json:"version"`
	Parameters   ToolParameters         `json:"parameters"`
	Visibility   ToolVisibility         `json:"visibility"`
	RequiredCaps []string               `json:"required_caps"`
	RateLimit    *RateLimit             `json:"rate_limit,omitempty"`
	TrustLevel   string                 `json:"trust_level"`
	CostPerUse   float64                `json:"cost_per_use"`
	Metadata     map[string]interface{} `json:"metadata"`
	OwnerAgentID string                 `json:"owner_agent_id"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// ToolParameters defines input/output for a tool
type ToolParameters struct {
	InputSchema  json.RawMessage `json:"input_schema"`
	OutputSchema json.RawMessage `json:"output_schema"`
	Timeout      time.Duration   `json:"timeout"`
	Streaming    bool            `json:"streaming"`
}

// ToolInvocation represents a tool call to a federated tool
type ToolInvocation struct {
	InvocationID  string                 `json:"invocation_id"`
	ToolID        string                 `json:"tool_id"`
	ToolName      string                 `json:"tool_name"`
	CallerAgentID string                 `json:"caller_agent_id"`
	TargetAgentID string                 `json:"target_agent_id"`
	Parameters    map[string]interface{} `json:"parameters"`
	Status        InvocationStatus       `json:"status"`
	Result        map[string]interface{} `json:"result,omitempty"`
	Error         string                 `json:"error,omitempty"`
	StartedAt     time.Time              `json:"started_at"`
	CompletedAt   *time.Time             `json:"completed_at,omitempty"`
	DurationMs    int64                  `json:"duration_ms"`
}

// InvocationStatus represents the status of a tool invocation
type InvocationStatus string

const (
	InvocationStatusPending   InvocationStatus = "pending"
	InvocationStatusRunning   InvocationStatus = "running"
	InvocationStatusCompleted InvocationStatus = "completed"
	InvocationStatusFailed    InvocationStatus = "failed"
	InvocationStatusCancelled InvocationStatus = "cancelled"
)

// FederationRequest represents a request to federate with another agent
type FederationRequest struct {
	RequestID      string        `json:"request_id"`
	FromAgentID    string        `json:"from_agent_id"`
	FromAgentName  string        `json:"from_agent_name"`
	ToAgentID      string        `json:"to_agent_id"`
	ToAgentName    string        `json:"to_agent_name"`
	RequestedTools []string      `json:"requested_tools"`
	IntendedUse    string        `json:"intended_use"`
	Duration       time.Duration `json:"duration"`
	CreatedAt      time.Time     `json:"created_at"`
}

// FederationResponse represents a response to a federation request
type FederationResponse struct {
	RequestID    string           `json:"request_id"`
	Status       FederationStatus `json:"status"`
	GrantedTools []GrantedTool    `json:"granted_tools,omitempty"`
	DeniedTools  []DeniedTool     `json:"denied_tools,omitempty"`
	Message      string           `json:"message"`
	ExpiresAt    *time.Time       `json:"expires_at,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
}

// GrantedTool represents a tool that was granted
type GrantedTool struct {
	ToolID     string     `json:"tool_id"`
	ToolName   string     `json:"tool_name"`
	RateLimit  *RateLimit `json:"rate_limit,omitempty"`
	MaxUses    int        `json:"max_uses,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	Conditions []string   `json:"conditions,omitempty"`
}

// DeniedTool represents a tool that was denied
type DeniedTool struct {
	ToolID    string `json:"tool_id"`
	ToolName  string `json:"tool_name"`
	Reason    string `json:"reason"`
	CanAppeal bool   `json:"can_appeal"`
}

// FederationConnection represents an active federation connection
type FederationConnection struct {
	ConnectionID    string            `json:"connection_id"`
	LocalAgentID    string            `json:"local_agent_id"`
	RemoteAgentID   string            `json:"remote_agent_id"`
	RemoteAgentName string            `json:"remote_agent_name"`
	Status          FederationStatus  `json:"status"`
	GrantedTools    []string          `json:"granted_tools"`
	AvailableTools  []*ToolDefinition `json:"available_tools"`
	LastSyncAt      time.Time         `json:"last_sync_at"`
	CreatedAt       time.Time         `json:"created_at"`
}

// RateLimit defines rate limiting for tool usage
type RateLimit struct {
	RequestsPerMinute int `json:"requests_per_minute"`
	RequestsPerHour   int `json:"requests_per_hour"`
	RequestsPerDay    int `json:"requests_per_day"`
	BurstSize         int `json:"burst_size"`
}

// Service manages tool and skill federation between agents
type Service struct {
	mu          sync.RWMutex
	bus         *bus.Bus
	tools       map[string]*ToolDefinition
	invocations map[string]*ToolInvocation
	federations map[string]*FederationRequest
	connections map[string]*FederationConnection
	grants      map[string][]GrantedTool
}

// NewService creates a new federation service
func NewService(b *bus.Bus) *Service {
	return &Service{
		bus:         b,
		tools:       make(map[string]*ToolDefinition),
		invocations: make(map[string]*ToolInvocation),
		federations: make(map[string]*FederationRequest),
		connections: make(map[string]*FederationConnection),
		grants:      make(map[string][]GrantedTool),
	}
}

// RegisterTool registers a tool for federation
func (s *Service) RegisterTool(tool *ToolDefinition) error {
	if tool.ToolID == "" {
		tool.ToolID = uuid.New().String()
	}
	tool.CreatedAt = time.Now().UTC()
	tool.UpdatedAt = time.Now().UTC()

	s.mu.Lock()
	s.tools[tool.ToolID] = tool
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("federation.tool.registered", "", map[string]interface{}{
		"tool_id":     tool.ToolID,
		"tool_name":   tool.Name,
		"owner_agent": tool.OwnerAgentID,
		"visibility":  tool.Visibility,
	}))

	return nil
}

// GetTool retrieves a tool definition
func (s *Service) GetTool(toolID string) (*ToolDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tool, exists := s.tools[toolID]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", toolID)
	}

	return tool, nil
}

// GetToolsByAgent retrieves all tools owned by an agent
func (s *Service) GetToolsByAgent(agentID string) []*ToolDefinition {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tools := make([]*ToolDefinition, 0)
	for _, tool := range s.tools {
		if tool.OwnerAgentID == agentID {
			tools = append(tools, tool)
		}
	}

	return tools
}

// GetVisibleTools retrieves tools visible to a specific agent
func (s *Service) GetVisibleTools(agentID string, trustLevel string) []*ToolDefinition {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tools := make([]*ToolDefinition, 0)
	for _, tool := range s.tools {
		if s.isToolVisible(tool, agentID, trustLevel) {
			tools = append(tools, tool)
		}
	}

	return tools
}

// isToolVisible checks if a tool is visible to an agent
func (s *Service) isToolVisible(tool *ToolDefinition, agentID string, trustLevel string) bool {
	switch tool.Visibility {
	case ToolVisibilityPublic:
		return true
	case ToolVisibilityFederated:
		return agentID != ""
	case ToolVisibilityTrusted:
		return trustLevel == "high" || trustLevel == "trusted"
	case ToolVisibilityPrivate:
		return tool.OwnerAgentID == agentID
	default:
		return false
	}
}

// RequestFederation requests federation with another agent
func (s *Service) RequestFederation(ctx context.Context, req *FederationRequest) (*FederationResponse, error) {
	req.RequestID = uuid.New().String()
	req.CreatedAt = time.Now().UTC()

	if req.Duration == 0 {
		req.Duration = 24 * time.Hour
	}

	s.mu.Lock()
	s.federations[req.RequestID] = req
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("federation.requested", "", map[string]interface{}{
		"request_id":      req.RequestID,
		"from_agent":      req.FromAgentID,
		"to_agent":        req.ToAgentID,
		"requested_tools": len(req.RequestedTools),
	}))

	return &FederationResponse{
		RequestID: req.RequestID,
		Status:    FederationStatusPending,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// ApproveFederation approves a federation request
func (s *Service) ApproveFederation(ctx context.Context, requestID string, grantedTools []GrantedTool, message string) (*FederationResponse, error) {
	s.mu.RLock()
	req, exists := s.federations[requestID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("federation request not found: %s", requestID)
	}

	expiresAt := time.Now().UTC().Add(req.Duration)

	response := &FederationResponse{
		RequestID:    requestID,
		Status:       FederationStatusApproved,
		GrantedTools: grantedTools,
		Message:      message,
		ExpiresAt:    &expiresAt,
		CreatedAt:    time.Now().UTC(),
	}

	// Create connection
	connection := &FederationConnection{
		ConnectionID:    uuid.New().String(),
		LocalAgentID:    req.ToAgentID,
		RemoteAgentID:   req.FromAgentID,
		RemoteAgentName: req.FromAgentName,
		Status:          FederationStatusActive,
		GrantedTools:    make([]string, 0),
		CreatedAt:       time.Now().UTC(),
		LastSyncAt:      time.Now().UTC(),
	}

	for _, gt := range grantedTools {
		connection.GrantedTools = append(connection.GrantedTools, gt.ToolID)
	}

	s.mu.Lock()
	s.connections[connection.ConnectionID] = connection
	s.grants[req.FromAgentID] = append(s.grants[req.FromAgentID], grantedTools...)
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("federation.approved", "", map[string]interface{}{
		"request_id":    requestID,
		"from_agent":    req.FromAgentID,
		"to_agent":      req.ToAgentID,
		"granted_tools": len(grantedTools),
	}))

	return response, nil
}

// InvokeTool invokes a federated tool
func (s *Service) InvokeTool(ctx context.Context, toolID string, callerAgentID string, parameters map[string]interface{}) (*ToolInvocation, error) {
	// Get tool definition
	tool, err := s.GetTool(toolID)
	if err != nil {
		return nil, err
	}

	invocation := &ToolInvocation{
		InvocationID:  uuid.New().String(),
		ToolID:        toolID,
		ToolName:      tool.Name,
		CallerAgentID: callerAgentID,
		TargetAgentID: tool.OwnerAgentID,
		Parameters:    parameters,
		Status:        InvocationStatusPending,
		StartedAt:     time.Now().UTC(),
	}

	s.mu.Lock()
	s.invocations[invocation.InvocationID] = invocation
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("federation.tool.invoked", "", map[string]interface{}{
		"invocation_id": invocation.InvocationID,
		"tool_id":       toolID,
		"tool_name":     tool.Name,
		"caller":        callerAgentID,
	}))

	// Simulate tool execution (in real implementation, would call actual tool)
	invocation.Status = InvocationStatusCompleted
	completedAt := time.Now().UTC()
	invocation.CompletedAt = &completedAt
	invocation.DurationMs = time.Since(invocation.StartedAt).Milliseconds()
	invocation.Result = map[string]interface{}{
		"success": true,
		"output":  "Tool executed successfully",
	}

	s.mu.Lock()
	s.invocations[invocation.InvocationID] = invocation
	s.mu.Unlock()

	return invocation, nil
}

// GetInvocationStatus gets the status of a tool invocation
func (s *Service) GetInvocationStatus(invocationID string) (*ToolInvocation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	invocation, exists := s.invocations[invocationID]
	if !exists {
		return nil, fmt.Errorf("invocation not found: %s", invocationID)
	}

	return invocation, nil
}

// GetFederationConnection gets a federation connection
func (s *Service) GetFederationConnection(connectionID string) (*FederationConnection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conn, exists := s.connections[connectionID]
	if !exists {
		return nil, fmt.Errorf("connection not found: %s", connectionID)
	}

	return conn, nil
}

// GetConnectionsByAgent gets all federation connections for an agent
func (s *Service) GetConnectionsByAgent(agentID string) []*FederationConnection {
	s.mu.RLock()
	defer s.mu.RUnlock()

	connections := make([]*FederationConnection, 0)
	for _, conn := range s.connections {
		if conn.LocalAgentID == agentID || conn.RemoteAgentID == agentID {
			connections = append(connections, conn)
		}
	}

	return connections
}

// RevokeTool revokes a granted tool
func (s *Service) RevokeTool(agentID string, toolID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	grants := s.grants[agentID]
	remaining := make([]GrantedTool, 0)

	for _, grant := range grants {
		if grant.ToolID != toolID {
			remaining = append(remaining, grant)
		}
	}

	s.grants[agentID] = remaining

	// Publish event
	s.bus.Publish(bus.NewEvent("federation.tool.revoked", "", map[string]interface{}{
		"agent_id": agentID,
		"tool_id":  toolID,
	}))

	return nil
}

// GetGrantedTools gets all granted tools for an agent
func (s *Service) GetGrantedTools(agentID string) []GrantedTool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.grants[agentID]
}

// SyncTools syncs tools with a federated agent
func (s *Service) SyncTools(connectionID string) ([]*ToolDefinition, error) {
	s.mu.RLock()
	conn, exists := s.connections[connectionID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("connection not found: %s", connectionID)
	}

	// Get available tools from remote agent
	availableTools := s.GetVisibleTools(conn.RemoteAgentID, "medium")

	// Update connection
	s.mu.Lock()
	conn.AvailableTools = availableTools
	conn.LastSyncAt = time.Now().UTC()
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("federation.synced", "", map[string]interface{}{
		"connection_id": connectionID,
		"tools_count":   len(availableTools),
	}))

	return availableTools, nil
}

// GetToolCount returns the number of registered tools
func (s *Service) GetToolCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.tools)
}

// GetConnectionCount returns the number of active connections
func (s *Service) GetConnectionCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.connections)
}

// GetInvocationCount returns the number of invocations
func (s *Service) GetInvocationCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.invocations)
}

// PrettyPrint prints tool info
func (t *ToolDefinition) PrettyPrint() string {
	return fmt.Sprintf("Tool{Name: %s, Version: %s, Visibility: %s, Owner: %s}",
		t.Name, t.Version, t.Visibility, t.OwnerAgentID)
}

// MarshalJSON for ToolDefinition
func (t *ToolDefinition) MarshalJSON() ([]byte, error) {
	type Alias ToolDefinition
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

// MarshalJSON for ToolInvocation
func (t *ToolInvocation) MarshalJSON() ([]byte, error) {
	type Alias ToolInvocation
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}
