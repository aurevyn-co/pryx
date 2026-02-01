package handoff

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"pryx-core/internal/bus"
)

// HandoffStatus represents the status of a session handoff
type HandoffStatus string

const (
	HandoffStatusPending    HandoffStatus = "pending"
	HandoffStatusAccepted   HandoffStatus = "accepted"
	HandoffStatusInProgress HandoffStatus = "in_progress"
	HandoffStatusCompleted  HandoffStatus = "completed"
	HandoffStatusFailed     HandoffStatus = "failed"
	HandoffStatusCancelled  HandoffStatus = "cancelled"
)

// HandoffPhase represents the phase of handoff
type HandoffPhase string

const (
	HandoffPhaseInitiation HandoffPhase = "initiation"
	HandoffPhaseTransfer   HandoffPhase = "transfer"
	HandoffPhaseValidation HandoffPhase = "validation"
	HandoffPhaseCompletion HandoffPhase = "completion"
)

// ContextType represents the type of context being transferred
type ContextType string

const (
	ContextTypeSession    ContextType = "session"
	ContextTypeMemory     ContextType = "memory"
	ContextTypePolicy     ContextType = "policy"
	ContextTypeCapability ContextType = "capability"
	ContextTypeState      ContextType = "state"
)

// HandoffRequest represents a request to hand off a session
type HandoffRequest struct {
	RequestID     string                 `json:"request_id"`
	FromAgentID   string                 `json:"from_agent_id"`
	FromAgentName string                 `json:"from_agent_name"`
	ToAgentID     string                 `json:"to_agent_id"`
	ToAgentName   string                 `json:"to_agent_name"`
	SessionID     string                 `json:"session_id"`
	ContextTypes  []ContextType          `json:"context_types"`
	Priority      int                    `json:"priority"`
	Timeout       time.Duration          `json:"timeout"`
	RequiresAck   bool                   `json:"requires_ack"`
	Metadata      map[string]interface{} `json:"metadata"`
	CreatedAt     time.Time              `json:"created_at"`
}

// HandoffContext represents context data being transferred
type HandoffContext struct {
	ContextID  string                 `json:"context_id"`
	Type       ContextType            `json:"type"`
	SessionID  string                 `json:"session_id"`
	Data       map[string]interface{} `json:"data"`
	Schema     json.RawMessage        `json:"schema"`
	SizeBytes  int64                  `json:"size_bytes"`
	Compressed bool                   `json:"compressed"`
	Checksum   string                 `json:"checksum"`
	CreatedAt  time.Time              `json:"created_at"`
}

// HandoffTransfer represents a handoff transfer event
type HandoffTransfer struct {
	TransferID       string       `json:"transfer_id"`
	RequestID        string       `json:"request_id"`
	Phase            HandoffPhase `json:"phase"`
	Progress         float64      `json:"progress"`
	BytesTransferred int64        `json:"bytes_transferred"`
	TotalBytes       int64        `json:"total_bytes"`
	Status           string       `json:"status"`
	Error            string       `json:"error,omitempty"`
	Timestamp        time.Time    `json:"timestamp"`
}

// HandoffResponse represents a response to a handoff request
type HandoffResponse struct {
	RequestID    string        `json:"request_id"`
	Status       HandoffStatus `json:"status"`
	AcceptedCaps []ContextType `json:"accepted_caps"`
	RejectedCaps []ContextType `json:"rejected_caps"`
	Conditions   []string      `json:"conditions,omitempty"`
	Message      string        `json:"message"`
	TransferURL  string        `json:"transfer_url,omitempty"`
	ExpiresAt    *time.Time    `json:"expires_at,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
}

// SessionState represents the state of a session being handed off
type SessionState struct {
	SessionID    string                 `json:"session_id"`
	AgentID      string                 `json:"agent_id"`
	Messages     []SessionMessage       `json:"messages"`
	Metadata     map[string]interface{} `json:"metadata"`
	Capabilities []string               `json:"capabilities"`
	Permissions  map[string]string      `json:"permissions"`
	LastActiveAt time.Time              `json:"last_active_at"`
	CreatedAt    time.Time              `json:"created_at"`
}

// SessionMessage represents a message in a session
type SessionMessage struct {
	MessageID string                 `json:"message_id"`
	Role      string                 `json:"role"`
	Content   string                 `json:"content"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// Service manages session handoffs between agents
type Service struct {
	mu             sync.RWMutex
	bus            *bus.Bus
	requests       map[string]*HandoffRequest
	responses      map[string]*HandoffResponse
	transfers      map[string]*HandoffTransfer
	history        []HandoffHistoryEntry
	activeHandoffs map[string]*ActiveHandoff
}

// HandoffHistoryEntry represents an entry in handoff history
type HandoffHistoryEntry struct {
	EntryID     string                 `json:"entry_id"`
	RequestID   string                 `json:"request_id"`
	FromAgentID string                 `json:"from_agent_id"`
	ToAgentID   string                 `json:"to_agent_id"`
	SessionID   string                 `json:"session_id"`
	Status      HandoffStatus          `json:"status"`
	Duration    time.Duration          `json:"duration"`
	Timestamp   time.Time              `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ActiveHandoff tracks an active handoff operation
type ActiveHandoff struct {
	Request    *HandoffRequest
	Response   *HandoffResponse
	Transfers  []*HandoffTransfer
	Contexts   []*HandoffContext
	StartTime  time.Time
	LastUpdate time.Time
}

// NewService creates a new handoff service
func NewService(b *bus.Bus) *Service {
	return &Service{
		bus:            b,
		requests:       make(map[string]*HandoffRequest),
		responses:      make(map[string]*HandoffResponse),
		transfers:      make(map[string]*HandoffTransfer),
		history:        make([]HandoffHistoryEntry, 0),
		activeHandoffs: make(map[string]*ActiveHandoff),
	}
}

// RequestHandoff initiates a session handoff request
func (s *Service) RequestHandoff(ctx context.Context, req *HandoffRequest) (*HandoffResponse, error) {
	req.RequestID = uuid.New().String()
	req.CreatedAt = time.Now().UTC()

	if req.Timeout == 0 {
		req.Timeout = 5 * time.Minute
	}

	s.mu.Lock()
	s.requests[req.RequestID] = req
	s.activeHandoffs[req.RequestID] = &ActiveHandoff{
		Request:    req,
		StartTime:  time.Now().UTC(),
		LastUpdate: time.Now().UTC(),
	}
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("handoff.requested", "", map[string]interface{}{
		"request_id": req.RequestID,
		"from_agent": req.FromAgentID,
		"to_agent":   req.ToAgentID,
		"session_id": req.SessionID,
	}))

	// Create pending response
	response := &HandoffResponse{
		RequestID: req.RequestID,
		Status:    HandoffStatusPending,
		CreatedAt: time.Now().UTC(),
	}

	return response, nil
}

// AcceptHandoff accepts a handoff request
func (s *Service) AcceptHandoff(ctx context.Context, requestID string, acceptedCaps []ContextType, conditions []string) (*HandoffResponse, error) {
	s.mu.RLock()
	req, exists := s.requests[requestID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("handoff request not found: %s", requestID)
	}

	expiresAt := time.Now().UTC().Add(req.Timeout)

	response := &HandoffResponse{
		RequestID:    requestID,
		Status:       HandoffStatusAccepted,
		AcceptedCaps: acceptedCaps,
		Conditions:   conditions,
		TransferURL:  fmt.Sprintf("/handoff/transfer/%s", requestID),
		ExpiresAt:    &expiresAt,
		CreatedAt:    time.Now().UTC(),
	}

	s.mu.Lock()
	s.responses[requestID] = response
	if active, ok := s.activeHandoffs[requestID]; ok {
		active.Response = response
		active.LastUpdate = time.Now().UTC()
	}
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("handoff.accepted", "", map[string]interface{}{
		"request_id":    requestID,
		"from_agent":    req.FromAgentID,
		"to_agent":      req.ToAgentID,
		"accepted_caps": len(acceptedCaps),
	}))

	return response, nil
}

// RejectHandoff rejects a handoff request
func (s *Service) RejectHandoff(ctx context.Context, requestID string, rejectedCaps []ContextType, message string) (*HandoffResponse, error) {
	s.mu.RLock()
	req, exists := s.requests[requestID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("handoff request not found: %s", requestID)
	}

	response := &HandoffResponse{
		RequestID:    requestID,
		Status:       HandoffStatusFailed,
		RejectedCaps: rejectedCaps,
		Message:      message,
		CreatedAt:    time.Now().UTC(),
	}

	s.mu.Lock()
	s.responses[requestID] = response
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("handoff.rejected", "", map[string]interface{}{
		"request_id": requestID,
		"from_agent": req.FromAgentID,
		"to_agent":   req.ToAgentID,
		"reason":     message,
	}))

	return response, nil
}

// StartTransfer starts the actual transfer of context data
func (s *Service) StartTransfer(ctx context.Context, requestID string, contexts []*HandoffContext) ([]*HandoffTransfer, error) {
	s.mu.RLock()
	req, exists := s.requests[requestID]
	active, activeExists := s.activeHandoffs[requestID]
	s.mu.RUnlock()

	if !exists || !activeExists {
		return nil, fmt.Errorf("handoff request not found: %s", requestID)
	}

	transfers := make([]*HandoffTransfer, 0)

	for _, ctx := range contexts {
		transfer := &HandoffTransfer{
			TransferID:       uuid.New().String(),
			RequestID:        requestID,
			Phase:            HandoffPhaseTransfer,
			Progress:         0,
			BytesTransferred: 0,
			TotalBytes:       ctx.SizeBytes,
			Status:           "in_progress",
			Timestamp:        time.Now().UTC(),
		}

		// Simulate transfer progress
		transfer.BytesTransferred = ctx.SizeBytes
		transfer.Progress = 1.0
		transfer.Status = "completed"

		transfers = append(transfers, transfer)

		// Update active handoff
		s.mu.Lock()
		active.Contexts = append(active.Contexts, ctx)
		active.Transfers = append(active.Transfers, transfer)
		active.LastUpdate = time.Now().UTC()
		s.mu.Unlock()
	}

	// Update handoff status
	s.mu.Lock()
	if active, ok := s.activeHandoffs[requestID]; ok {
		active.Request.ContextTypes = req.ContextTypes
	}
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("handoff.transfer.started", "", map[string]interface{}{
		"request_id":  requestID,
		"contexts":    len(contexts),
		"total_bytes": contexts[0].SizeBytes,
	}))

	return transfers, nil
}

// ValidateTransfer validates the transferred context
func (s *Service) ValidateTransfer(ctx context.Context, requestID string) (*HandoffTransfer, error) {
	s.mu.RLock()
	active, exists := s.activeHandoffs[requestID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("handoff not found: %s", requestID)
	}

	transfer := &HandoffTransfer{
		TransferID: uuid.New().String(),
		RequestID:  requestID,
		Phase:      HandoffPhaseValidation,
		Progress:   1.0,
		Status:     "completed",
		Timestamp:  time.Now().UTC(),
	}

	// Update active handoff
	s.mu.Lock()
	active.Transfers = append(active.Transfers, transfer)
	active.LastUpdate = time.Now().UTC()
	s.mu.Unlock()

	return transfer, nil
}

// CompleteHandoff marks a handoff as completed
func (s *Service) CompleteHandoff(ctx context.Context, requestID string) error {
	s.mu.RLock()
	req, reqExists := s.requests[requestID]
	active, activeExists := s.activeHandoffs[requestID]
	s.mu.RUnlock()

	if !reqExists || !activeExists {
		return fmt.Errorf("handoff not found: %s", requestID)
	}

	duration := time.Since(active.StartTime)

	// Add to history
	entry := HandoffHistoryEntry{
		EntryID:     uuid.New().String(),
		RequestID:   requestID,
		FromAgentID: req.FromAgentID,
		ToAgentID:   req.ToAgentID,
		SessionID:   req.SessionID,
		Status:      HandoffStatusCompleted,
		Duration:    duration,
		Timestamp:   time.Now().UTC(),
	}

	s.mu.Lock()
	s.history = append(s.history, entry)
	delete(s.activeHandoffs, requestID)
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("handoff.completed", "", map[string]interface{}{
		"request_id":  requestID,
		"from_agent":  req.FromAgentID,
		"to_agent":    req.ToAgentID,
		"session_id":  req.SessionID,
		"duration_ms": duration.Milliseconds(),
	}))

	return nil
}

// CancelHandoff cancels a pending handoff
func (s *Service) CancelHandoff(ctx context.Context, requestID string, reason string) error {
	s.mu.RLock()
	req, exists := s.requests[requestID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("handoff request not found: %s", requestID)
	}

	s.mu.Lock()
	delete(s.activeHandoffs, requestID)
	s.mu.Unlock()

	// Add failed entry to history
	entry := HandoffHistoryEntry{
		EntryID:     uuid.New().String(),
		RequestID:   requestID,
		FromAgentID: req.FromAgentID,
		ToAgentID:   req.ToAgentID,
		SessionID:   req.SessionID,
		Status:      HandoffStatusCancelled,
		Timestamp:   time.Now().UTC(),
		Metadata: map[string]interface{}{
			"reason": reason,
		},
	}

	s.mu.Lock()
	s.history = append(s.history, entry)
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("handoff.cancelled", "", map[string]interface{}{
		"request_id": requestID,
		"reason":     reason,
	}))

	return nil
}

// GetHandoffStatus gets the current status of a handoff
func (s *Service) GetHandoffStatus(requestID string) (*HandoffResponse, error) {
	s.mu.RLock()
	response, exists := s.responses[requestID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("handoff response not found: %s", requestID)
	}

	return response, nil
}

// GetActiveHandoffs gets all active handoffs
func (s *Service) GetActiveHandoffs() []*ActiveHandoff {
	s.mu.RLock()
	defer s.mu.RUnlock()

	handoffs := make([]*ActiveHandoff, 0, len(s.activeHandoffs))
	for _, h := range s.activeHandoffs {
		handoffs = append(handoffs, h)
	}

	return handoffs
}

// GetHandoffHistory gets the handoff history for an agent
func (s *Service) GetHandoffHistory(agentID string, limit int) []HandoffHistoryEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries := make([]HandoffHistoryEntry, 0)
	for _, entry := range s.history {
		if entry.FromAgentID == agentID || entry.ToAgentID == agentID {
			entries = append(entries, entry)
			if len(entries) >= limit {
				break
			}
		}
	}

	return entries
}

// CreateSessionState creates a session state for handoff
func (s *Service) CreateSessionState(sessionID string, agentID string, messages []SessionMessage, metadata map[string]interface{}) *SessionState {
	return &SessionState{
		SessionID:    sessionID,
		AgentID:      agentID,
		Messages:     messages,
		Metadata:     metadata,
		LastActiveAt: time.Now().UTC(),
		CreatedAt:    time.Now().UTC(),
	}
}

// CreateContext creates a handoff context
func (s *Service) CreateContext(ctxType ContextType, sessionID string, data map[string]interface{}) *HandoffContext {
	return &HandoffContext{
		ContextID: uuid.New().String(),
		Type:      ctxType,
		SessionID: sessionID,
		Data:      data,
		CreatedAt: time.Now().UTC(),
	}
}

// GetRequestCount returns the number of pending requests
func (s *Service) GetRequestCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.requests)
}

// GetActiveCount returns the number of active handoffs
func (s *Service) GetActiveCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.activeHandoffs)
}

// GetHistoryCount returns the total number of handoffs in history
func (s *Service) GetHistoryCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.history)
}

// PrettyPrint prints request info
func (r *HandoffRequest) PrettyPrint() string {
	return fmt.Sprintf("Handoff{RequestID: %s, From: %s, To: %s, Session: %s}",
		r.RequestID, r.FromAgentID, r.ToAgentID, r.SessionID)
}

// MarshalJSON for HandoffRequest
func (r *HandoffRequest) MarshalJSON() ([]byte, error) {
	type Alias HandoffRequest
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}

// MarshalJSON for HandoffResponse
func (r *HandoffResponse) MarshalJSON() ([]byte, error) {
	type Alias HandoffResponse
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}
