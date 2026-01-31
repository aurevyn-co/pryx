// Package bus provides an in-process event bus for pub/sub communication.
// It supports typed events, topic filtering, and concurrent subscriptions.
package bus

import "time"

// EventType represents the type of event in the system.
// Event types are used for routing and filtering messages.
type EventType string

// Predefined event types for the Pryx event bus.
const (
	// EventEnvelope is the generic envelope type for all events.
	EventEnvelope EventType = "event"
	// EventSessionMessage is emitted when a new message is added to a session.
	EventSessionMessage EventType = "session.message"
	// EventSessionTyping is emitted when typing indicators change.
	EventSessionTyping EventType = "session.typing"
	// EventToolRequest is emitted when a tool execution is requested.
	EventToolRequest EventType = "tool.request"
	// EventToolExecuting is emitted when a tool starts executing.
	EventToolExecuting EventType = "tool.executing"
	// EventToolComplete is emitted when a tool finishes executing.
	EventToolComplete EventType = "tool.complete"
	// EventApprovalNeeded is emitted when user approval is required.
	EventApprovalNeeded EventType = "approval.needed"
	// EventApprovalResolved is emitted when an approval is resolved.
	EventApprovalResolved EventType = "approval.resolved"
	// EventTraceEvent is emitted for trace/debug events.
	EventTraceEvent EventType = "trace.event"
	// EventErrorOccurred is emitted when an error occurs.
	EventErrorOccurred EventType = "error.occurred"
	// EventChannelStatus is emitted when channel status changes.
	EventChannelStatus EventType = "channel.status"
	// EventChannelMessage is emitted when a message is received from a channel.
	EventChannelMessage EventType = "channel.message"
	// EventChannelOutboundMessage is emitted when a message is sent to a channel.
	EventChannelOutboundMessage EventType = "channel.outbound_message"
	// EventChatRequest is emitted when a chat request is made.
	EventChatRequest EventType = "chat.request"
)

// Event represents a single event in the system.
// Events are the primary mechanism for communication between components.
type Event struct {
	// Type is the envelope type (usually "event").
	Type EventType `json:"type"`
	// Event is the specific event type (e.g., "session.message").
	Event EventType `json:"event"`
	// SessionID identifies the session this event belongs to, if any.
	SessionID string `json:"session_id,omitempty"`
	// Surface identifies the UI surface that originated or should receive this event.
	Surface string `json:"surface,omitempty"`
	// Payload contains the event-specific data.
	Payload interface{} `json:"payload"`
	// Timestamp is when the event was created.
	Timestamp time.Time `json:"timestamp"`
	// Version is a monotonically increasing event version.
	Version int `json:"version"`
}

// NewEvent creates a new event with the current timestamp.
// The event type specifies what kind of event this is, sessionID identifies
// the associated session (can be empty), and payload contains the event data.
func NewEvent(eventType EventType, sessionID string, payload interface{}) Event {
	return Event{
		Type:      EventEnvelope,
		Event:     eventType,
		SessionID: sessionID,
		Payload:   payload,
		Timestamp: time.Now().UTC(),
		Version:   0,
	}
}
