package message

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"pryx-core/internal/bus"
)

// MessageType represents the type of message
type MessageType string

const (
	MessageTypeRequest     MessageType = "request"
	MessageTypeResponse    MessageType = "response"
	MessageTypeEvent       MessageType = "event"
	MessageTypeError       MessageType = "error"
	MessageTypeStreamChunk MessageType = "stream_chunk"
)

// MessagePriority represents message priority levels
type MessagePriority string

const (
	PriorityLow      MessagePriority = "low"
	PriorityNormal   MessagePriority = "normal"
	PriorityHigh     MessagePriority = "high"
	PriorityCritical MessagePriority = "critical"
)

// Message represents an agent-to-agent message
type Message struct {
	ID            string                 `json:"id"`
	Type          MessageType            `json:"type"`
	FromAgent     string                 `json:"from_agent"`
	ToAgent       string                 `json:"to_agent"`
	CorrelationID string                 `json:"correlation_id"`
	Priority      MessagePriority        `json:"priority"`
	Payload       map[string]interface{} `json:"payload"`
	Metadata      MessageMetadata        `json:"metadata"`
	CreatedAt     time.Time              `json:"created_at"`
	ExpiresAt     *time.Time             `json:"expires_at,omitempty"`
}

// MessageMetadata contains additional message information
type MessageMetadata struct {
	ContentType     string            `json:"content_type"`
	ContentEncoding string            `json:"content_encoding"`
	Headers         map[string]string `json:"headers"`
	TraceID         string            `json:"trace_id"`
	SpanID          string            `json:"span_id"`
	RetryCount      int               `json:"retry_count"`
	MaxRetries      int               `json:"max_retries"`
}

// Response represents a message response
type Response struct {
	CorrelationID string                 `json:"correlation_id"`
	StatusCode    int                    `json:"status_code"`
	StatusMessage string                 `json:"status_message"`
	Payload       map[string]interface{} `json:"payload"`
	Headers       map[string]string      `json:"headers"`
	CreatedAt     time.Time              `json:"created_at"`
}

// StreamMessage represents a streaming message chunk
type StreamMessage struct {
	MessageID   string    `json:"message_id"`
	Sequence    int       `json:"sequence"`
	Chunk       []byte    `json:"chunk"`
	ContentType string    `json:"content_type"`
	IsLast      bool      `json:"is_last"`
	Timestamp   time.Time `json:"timestamp"`
}

// MessageRequest represents a request to send a message
type MessageRequest struct {
	ToAgent     string                 `json:"to_agent"`
	Type        MessageType            `json:"type"`
	Priority    MessagePriority        `json:"priority"`
	Payload     map[string]interface{} `json:"payload"`
	ContentType string                 `json:"content_type"`
	Timeout     time.Duration          `json:"timeout"`
	RequiresAck bool                   `json:"requires_ack"`
	MaxRetries  int                    `json:"max_retries"`
}

// MessageResponse represents the result of sending a message
type MessageResponse struct {
	MessageID     string     `json:"message_id"`
	Status        string     `json:"status"` // "queued", "delivered", "failed", "timeout"
	Error         string     `json:"error,omitempty"`
	CorrelationID string     `json:"correlation_id,omitempty"`
	DeliveredAt   *time.Time `json:"delivered_at,omitempty"`
}

// Service manages message exchange between agents
type Service struct {
	mu           sync.RWMutex
	bus          *bus.Bus
	localAgentID string
	messages     map[string]*Message
	responses    map[string]*Response
	inFlight     map[string]*Message
	handlers     map[MessageType]MessageHandler
}

// MessageHandler is a function that handles incoming messages
type MessageHandler func(ctx context.Context, msg *Message) (*Response, error)

// NewService creates a new message exchange service
func NewService(b *bus.Bus, localAgentID string) *Service {
	return &Service{
		bus:          b,
		localAgentID: localAgentID,
		messages:     make(map[string]*Message),
		responses:    make(map[string]*Response),
		inFlight:     make(map[string]*Message),
		handlers:     make(map[MessageType]MessageHandler),
	}
}

// RegisterHandler registers a message handler for a specific message type
func (s *Service) RegisterHandler(msgType MessageType, handler MessageHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[msgType] = handler
}

// Send sends a message to an agent
func (s *Service) Send(ctx context.Context, req *MessageRequest) (*MessageResponse, error) {
	msgID := uuid.New().String()
	correlationID := uuid.New().String()

	msg := &Message{
		ID:            msgID,
		Type:          req.Type,
		FromAgent:     s.localAgentID,
		ToAgent:       req.ToAgent,
		CorrelationID: correlationID,
		Priority:      req.Priority,
		Payload:       req.Payload,
		Metadata: MessageMetadata{
			ContentType: req.ContentType,
			RetryCount:  0,
			MaxRetries:  req.MaxRetries,
		},
		CreatedAt: time.Now().UTC(),
	}

	if req.Timeout > 0 {
		expiresAt := time.Now().UTC().Add(req.Timeout)
		msg.ExpiresAt = &expiresAt
	}

	s.mu.Lock()
	s.messages[msgID] = msg
	if req.RequiresAck {
		s.inFlight[correlationID] = msg
	}
	s.mu.Unlock()

	s.bus.Publish(bus.NewEvent("message.sent", "", map[string]interface{}{
		"message_id":     msgID,
		"to_agent":       req.ToAgent,
		"type":           req.Type,
		"correlation_id": correlationID,
	}))

	return &MessageResponse{
		MessageID:     msgID,
		Status:        "queued",
		CorrelationID: correlationID,
	}, nil
}

// HandleMessage handles an incoming message
func (s *Service) HandleMessage(ctx context.Context, msg *Message) (*Response, error) {
	if msg.ExpiresAt != nil && msg.ExpiresAt.Before(time.Now()) {
		return &Response{
			CorrelationID: msg.CorrelationID,
			StatusCode:    408,
			StatusMessage: "Request expired",
			CreatedAt:     time.Now().UTC(),
		}, nil
	}

	s.mu.RLock()
	handler, exists := s.handlers[msg.Type]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no handler registered for message type: %s", msg.Type)
	}

	response, err := handler(ctx, msg)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	response.CorrelationID = msg.CorrelationID
	return response, nil
}

// SendResponse sends a response to a message
func (s *Service) SendResponse(ctx context.Context, response *Response) error {
	s.mu.Lock()
	s.responses[response.CorrelationID] = response
	delete(s.inFlight, response.CorrelationID)
	s.mu.Unlock()

	s.bus.Publish(bus.NewEvent("message.response", "", map[string]interface{}{
		"correlation_id": response.CorrelationID,
		"status_code":    response.StatusCode,
	}))

	return nil
}

// GetMessage retrieves a message by ID
func (s *Service) GetMessage(msgID string) (*Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msg, exists := s.messages[msgID]
	if !exists {
		return nil, fmt.Errorf("message not found: %s", msgID)
	}

	return msg, nil
}

// NewMessage creates a new message
func NewMessage(msgType MessageType, fromAgent, toAgent string, payload map[string]interface{}) *Message {
	return &Message{
		ID:            uuid.New().String(),
		Type:          msgType,
		FromAgent:     fromAgent,
		ToAgent:       toAgent,
		CorrelationID: uuid.New().String(),
		Priority:      PriorityNormal,
		Payload:       payload,
		Metadata: MessageMetadata{
			ContentType: "application/json",
			Headers:     make(map[string]string),
			RetryCount:  0,
			MaxRetries:  3,
		},
		CreatedAt: time.Now().UTC(),
	}
}

// NewRequest creates a new message request
func NewRequest(toAgent string, msgType MessageType, payload map[string]interface{}) *MessageRequest {
	return &MessageRequest{
		ToAgent:     toAgent,
		Type:        msgType,
		Priority:    PriorityNormal,
		Payload:     payload,
		ContentType: "application/json",
		Timeout:     30 * time.Second,
		RequiresAck: true,
		MaxRetries:  3,
	}
}

// NewResponse creates a new response
func NewResponse(correlationID string, statusCode int, statusMessage string, payload map[string]interface{}) *Response {
	return &Response{
		CorrelationID: correlationID,
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Payload:       payload,
		CreatedAt:     time.Now().UTC(),
	}
}

// MarshalJSON for Message
func (m *Message) MarshalJSON() ([]byte, error) {
	type Alias Message
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// MarshalJSON for Response
func (r *Response) MarshalJSON() ([]byte, error) {
	type Alias Response
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}
