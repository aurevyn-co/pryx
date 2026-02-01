package agentmessaging

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"pryx-core/internal/bus"
	"pryx-core/internal/message"
)

// AgentSession represents an active session with another agent
type AgentSession struct {
	SessionID       string                 `json:"session_id"`
	RemoteAgentID   string                 `json:"remote_agent_id"`
	RemoteAgentName string                 `json:"remote_agent_name"`
	RemoteEndpoint  string                 `json:"remote_endpoint"`
	Status          SessionStatus          `json:"status"`
	Messages        []*message.Message     `json:"messages"`
	CreatedAt       time.Time              `json:"created_at"`
	LastActivity    time.Time              `json:"last_activity"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// SessionStatus represents the status of an agent session
type SessionStatus string

const (
	SessionStatusConnecting   SessionStatus = "connecting"
	SessionStatusConnected    SessionStatus = "connected"
	SessionStatusDisconnected SessionStatus = "disconnected"
	SessionStatusError        SessionStatus = "error"
)

// Conversation represents a conversation thread between agents
type Conversation struct {
	ConversationID string                 `json:"conversation_id"`
	Participants   []string               `json:"participants"` // Agent IDs
	Messages       []*message.Message     `json:"messages"`
	Subject        string                 `json:"subject"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	Status         ConversationStatus     `json:"status"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// ConversationStatus represents the status of a conversation
type ConversationStatus string

const (
	ConversationStatusActive   ConversationStatus = "active"
	ConversationStatusArchived ConversationStatus = "archived"
	ConversationStatusClosed   ConversationStatus = "closed"
)

// MessageReceipt confirms message delivery
type MessageReceipt struct {
	ReceiptID   string    `json:"receipt_id"`
	MessageID   string    `json:"message_id"`
	FromAgent   string    `json:"from_agent"`
	ToAgent     string    `json:"to_agent"`
	DeliveredAt time.Time `json:"delivered_at"`
	Status      string    `json:"status"` // "delivered", "read", "failed"
	LatencyMs   int64     `json:"latency_ms"`
}

// MessageHistoryEntry represents an entry in message history
type MessageHistoryEntry struct {
	EntryID        string    `json:"entry_id"`
	MessageID      string    `json:"message_id"`
	ConversationID string    `json:"conversation_id"`
	FromAgent      string    `json:"from_agent"`
	ToAgent        string    `json:"to_agent"`
	MessageType    string    `json:"message_type"`
	Summary        string    `json:"summary"`
	Timestamp      time.Time `json:"timestamp"`
	Tags           []string  `json:"tags"`
}

// Service manages agent-to-agent messaging
type Service struct {
	mu            sync.RWMutex
	bus           *bus.Bus
	msgService    *message.Service
	sessions      map[string]*AgentSession
	conversations map[string]*Conversation
	history       []MessageHistoryEntry
	replyTo       map[string]string // message_id -> conversation_id mapping
}

// NewService creates a new agent messaging service
func NewService(b *bus.Bus, msgSvc *message.Service) *Service {
	return &Service{
		bus:           b,
		msgService:    msgSvc,
		sessions:      make(map[string]*AgentSession),
		conversations: make(map[string]*Conversation),
		history:       make([]MessageHistoryEntry, 0),
		replyTo:       make(map[string]string),
	}
}

// ConnectSession establishes a session with a remote agent
func (s *Service) ConnectSession(ctx context.Context, agentID, agentName, endpoint string) (*AgentSession, error) {
	session := &AgentSession{
		SessionID:       uuid.New().String(),
		RemoteAgentID:   agentID,
		RemoteAgentName: agentName,
		RemoteEndpoint:  endpoint,
		Status:          SessionStatusConnecting,
		Messages:        make([]*message.Message, 0),
		CreatedAt:       time.Now().UTC(),
		LastActivity:    time.Now().UTC(),
		Metadata:        make(map[string]interface{}),
	}

	// Store session
	s.mu.Lock()
	s.sessions[session.SessionID] = session
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("agent.session.connecting", "", map[string]interface{}{
		"session_id":      session.SessionID,
		"remote_agent_id": agentID,
		"endpoint":        endpoint,
	}))

	// Simulate connection (in real implementation, would establish actual connection)
	session.Status = SessionStatusConnected
	session.LastActivity = time.Now().UTC()

	s.bus.Publish(bus.NewEvent("agent.session.connected", "", map[string]interface{}{
		"session_id":      session.SessionID,
		"remote_agent_id": agentID,
	}))

	return session, nil
}

// DisconnectSession closes a session with a remote agent
func (s *Service) DisconnectSession(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	session.Status = SessionStatusDisconnected
	session.LastActivity = time.Now().UTC()

	// Publish event
	s.bus.Publish(bus.NewEvent("agent.session.disconnected", "", map[string]interface{}{
		"session_id":      sessionID,
		"remote_agent_id": session.RemoteAgentID,
	}))

	return nil
}

// SendMessage sends a message to a remote agent through a session
func (s *Service) SendMessage(ctx context.Context, sessionID string, msgType message.MessageType, payload map[string]interface{}) (*MessageReceipt, error) {
	s.mu.RLock()
	session, exists := s.sessions[sessionID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	if session.Status != SessionStatusConnected {
		return nil, fmt.Errorf("session not connected: %s", session.Status)
	}

	// Create message request
	req := &message.MessageRequest{
		ToAgent:     session.RemoteAgentID,
		Type:        msgType,
		Priority:    message.PriorityNormal,
		Payload:     payload,
		ContentType: "application/json",
		Timeout:     30 * time.Second,
		RequiresAck: true,
		MaxRetries:  3,
	}

	// Send message through message service
	resp, err := s.msgService.Send(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Create receipt
	receipt := &MessageReceipt{
		ReceiptID:   uuid.New().String(),
		MessageID:   resp.MessageID,
		FromAgent:   "", // Will be set by runtime
		ToAgent:     session.RemoteAgentID,
		DeliveredAt: time.Now().UTC(),
		Status:      "delivered",
		LatencyMs:   0, // Would be calculated in real implementation
	}

	// Update session
	s.mu.Lock()
	session.Messages = append(session.Messages, &message.Message{
		ID:        resp.MessageID,
		Type:      msgType,
		ToAgent:   session.RemoteAgentID,
		CreatedAt: time.Now().UTC(),
	})
	session.LastActivity = time.Now().UTC()
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("agent.message.sent", "", map[string]interface{}{
		"session_id": sessionID,
		"message_id": resp.MessageID,
		"to_agent":   session.RemoteAgentID,
	}))

	return receipt, nil
}

// SendDirectMessage sends a message directly to an agent without a session
func (s *Service) SendDirectMessage(ctx context.Context, toAgentID string, msgType message.MessageType, payload map[string]interface{}) (*MessageReceipt, error) {
	req := &message.MessageRequest{
		ToAgent:     toAgentID,
		Type:        msgType,
		Priority:    message.PriorityNormal,
		Payload:     payload,
		ContentType: "application/json",
		Timeout:     30 * time.Second,
		RequiresAck: true,
		MaxRetries:  3,
	}

	resp, err := s.msgService.Send(ctx, req)
	if err != nil {
		return nil, err
	}

	receipt := &MessageReceipt{
		ReceiptID:   uuid.New().String(),
		MessageID:   resp.MessageID,
		FromAgent:   "",
		ToAgent:     toAgentID,
		DeliveredAt: time.Now().UTC(),
		Status:      "delivered",
	}

	return receipt, nil
}

// ReplyToMessage sends a reply to a specific message
func (s *Service) ReplyToMessage(ctx context.Context, originalMessageID string, payload map[string]interface{}) (*MessageReceipt, error) {
	s.mu.RLock()
	conversationID, exists := s.replyTo[originalMessageID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("original message not found: %s", originalMessageID)
	}

	s.mu.RLock()
	conv, exists := s.conversations[conversationID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("conversation not found: %s", conversationID)
	}

	// Get the original message to find the sender
	var originalSender string
	for _, msg := range conv.Messages {
		if msg.ID == originalMessageID {
			originalSender = msg.FromAgent
			break
		}
	}

	if originalSender == "" {
		return nil, fmt.Errorf("original sender not found")
	}

	// Send reply
	return s.SendDirectMessage(ctx, originalSender, message.MessageTypeResponse, payload)
}

// StartConversation starts a new conversation with an agent
func (s *Service) StartConversation(ctx context.Context, participantIDs []string, subject string) (*Conversation, error) {
	conv := &Conversation{
		ConversationID: uuid.New().String(),
		Participants:   participantIDs,
		Messages:       make([]*message.Message, 0),
		Subject:        subject,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Status:         ConversationStatusActive,
		Metadata:       make(map[string]interface{}),
	}

	// Store conversation
	s.mu.Lock()
	s.conversations[conv.ConversationID] = conv
	s.mu.Unlock()

	// Publish event
	s.bus.Publish(bus.NewEvent("agent.conversation.started", "", map[string]interface{}{
		"conversation_id": conv.ConversationID,
		"participants":    participantIDs,
		"subject":         subject,
	}))

	return conv, nil
}

// AddMessageToConversation adds a message to a conversation
func (s *Service) AddMessageToConversation(ctx context.Context, conversationID string, fromAgent string, msg *message.Message) error {
	s.mu.Lock()
	conv, exists := s.conversations[conversationID]
	s.mu.Unlock()

	if !exists {
		return fmt.Errorf("conversation not found: %s", conversationID)
	}

	// Add message to conversation
	s.mu.Lock()
	conv.Messages = append(conv.Messages, msg)
	conv.UpdatedAt = time.Now().UTC()
	s.mu.Unlock()

	// Track reply-to mapping
	s.mu.Lock()
	s.replyTo[msg.ID] = conversationID
	s.mu.Unlock()

	// Add to history
	historyEntry := MessageHistoryEntry{
		EntryID:        uuid.New().String(),
		MessageID:      msg.ID,
		ConversationID: conversationID,
		FromAgent:      fromAgent,
		ToAgent:        msg.ToAgent,
		MessageType:    string(msg.Type),
		Summary:        s.summarizeMessage(msg),
		Timestamp:      time.Now().UTC(),
		Tags:           []string{conversationID},
	}

	s.mu.Lock()
	s.history = append(s.history, historyEntry)
	s.mu.Unlock()

	return nil
}

// GetConversation retrieves a conversation by ID
func (s *Service) GetConversation(conversationID string) (*Conversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conv, exists := s.conversations[conversationID]
	if !exists {
		return nil, fmt.Errorf("conversation not found: %s", conversationID)
	}

	return conv, nil
}

// GetConversationsByParticipant retrieves all conversations for an agent
func (s *Service) GetConversationsByParticipant(agentID string) []*Conversation {
	s.mu.RLock()
	defer s.mu.RUnlock()

	convs := make([]*Conversation, 0)
	for _, conv := range s.conversations {
		for _, participant := range conv.Participants {
			if participant == agentID {
				convs = append(convs, conv)
				break
			}
		}
	}

	return convs
}

// ArchiveConversation archives a conversation
func (s *Service) ArchiveConversation(conversationID string) error {
	s.mu.Lock()
	conv, exists := s.conversations[conversationID]
	s.mu.Unlock()

	if !exists {
		return fmt.Errorf("conversation not found: %s", conversationID)
	}

	s.mu.Lock()
	conv.Status = ConversationStatusArchived
	conv.UpdatedAt = time.Now().UTC()
	s.mu.Unlock()

	return nil
}

// GetMessageHistory retrieves message history for an agent
func (s *Service) GetMessageHistory(agentID string, limit int) []MessageHistoryEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries := make([]MessageHistoryEntry, 0)
	for _, entry := range s.history {
		if entry.FromAgent == agentID || entry.ToAgent == agentID {
			entries = append(entries, entry)
			if len(entries) >= limit {
				break
			}
		}
	}

	return entries
}

// GetActiveSessions retrieves all active sessions
func (s *Service) GetActiveSessions() []*AgentSession {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessions := make([]*AgentSession, 0, len(s.sessions))
	for _, session := range s.sessions {
		if session.Status == SessionStatusConnected {
			sessions = append(sessions, session)
		}
	}

	return sessions
}

// GetSession retrieves a session by ID
func (s *Service) GetSession(sessionID string) (*AgentSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	return session, nil
}

// GetSessionsByAgent retrieves all sessions with a specific agent
func (s *Service) GetSessionsByAgent(agentID string) []*AgentSession {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessions := make([]*AgentSession, 0)
	for _, session := range s.sessions {
		if session.RemoteAgentID == agentID {
			sessions = append(sessions, session)
		}
	}

	return sessions
}

// MarkMessageRead marks a message as read
func (s *Service) MarkMessageRead(messageID string) error {
	s.mu.RLock()
	conversationID, exists := s.replyTo[messageID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("message not found: %s", messageID)
	}

	s.mu.RLock()
	conv, exists := s.conversations[conversationID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("conversation not found: %s", conversationID)
	}

	for _, msg := range conv.Messages {
		if msg.ID == messageID {
			// In a real implementation, would update message status
			break
		}
	}

	return nil
}

// summarizeMessage creates a brief summary of a message
func (s *Service) summarizeMessage(msg *message.Message) string {
	// Simple summarization based on message type and payload
	switch msg.Type {
	case message.MessageTypeRequest:
		return fmt.Sprintf("Request to %s", msg.ToAgent)
	case message.MessageTypeResponse:
		return "Response message"
	case message.MessageTypeEvent:
		return "Event notification"
	case message.MessageTypeError:
		return "Error message"
	default:
		return "Message"
	}
}

// GetActiveSessionCount returns the number of active sessions
func (s *Service) GetActiveSessionCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, session := range s.sessions {
		if session.Status == SessionStatusConnected {
			count++
		}
	}

	return count
}

// GetConversationCount returns the total number of conversations
func (s *Service) GetConversationCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.conversations)
}

// PrettyPrint prints session info
func (s *AgentSession) PrettyPrint() string {
	return fmt.Sprintf("Session{SessionID: %s, Agent: %s (%s), Status: %s, Messages: %d}",
		s.SessionID, s.RemoteAgentName, s.RemoteAgentID, s.Status, len(s.Messages))
}

// PrettyPrint prints conversation info
func (c *Conversation) PrettyPrint() string {
	return fmt.Sprintf("Conversation{ID: %s, Subject: %s, Participants: %d, Messages: %d, Status: %s}",
		c.ConversationID, c.Subject, len(c.Participants), len(c.Messages), c.Status)
}

// MarshalJSON for AgentSession
func (s *AgentSession) MarshalJSON() ([]byte, error) {
	type Alias AgentSession
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	})
}

// MarshalJSON for Conversation
func (c *Conversation) MarshalJSON() ([]byte, error) {
	type Alias Conversation
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	})
}
