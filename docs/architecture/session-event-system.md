# Session & Event System Architecture

This document describes the session management and event bus architecture in Pryx.

## Overview

Pryx uses a publish/subscribe event bus pattern for loose coupling between components, combined with persistent session storage.

## Event Bus

Located in `apps/runtime/internal/bus/bus.go`

### Core Concepts

- **Publisher-Subscriber Pattern**: Components publish events without knowing who subscribes
- **Topic-Based**: Subscribers can filter by event types (topics)
- **Concurrent-Safe**: Uses RWMutex for thread-safe operations
- **Versioned Events**: Each event has a monotonic version number for ordering
- **Backpressure Handling**: Slow subscribers are dropped to prevent blocking

### Event Types

```go
// Event represents any application event
type Event struct {
    Type    EventType  // Event type identifier
    Payload interface{} // Event data
    Version int64     // Monotonic version for ordering
}

// Common event types
const (
    EventChannelMessage   = "channel:message"
    EventChannelOutbound = "channel:outbound_message"
    EventAgentSpawned   = "agent:spawned"
    EventAgentCompleted = "agent:completed"
    EventErrorOccurred  = "error:occurred"
    EventTraceEvent     = "trace:event"
    // ... add more as needed
)
```

### API

```go
// Create new bus
bus := bus.New()

// Subscribe to topics (returns channel and closer)
ch, closer := bus.Subscribe(EventChannelMessage, EventAgentSpawned)

// Receive events
for event := range ch {
    handleEvent(event)
}

// Unsubscribe when done
closer()  // Or bus.Unsubscribe(id)

// Publish events
bus.Publish(Event{
    Type: EventChannelMessage,
    Payload: Message{...},
})

// Shutdown
bus.Shutdown()
```

## Session Store

Located in `apps/runtime/internal/store/session.go`

### Schema

```sql
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### API

```go
// Create new session
session, err := store.CreateSession("Chat with User")

// Get session by ID
session, err := store.GetSession(id)

// List all sessions (ordered by updated_at, capped at 100)
sessions, err := store.ListSessions()

// Ensure session exists (create if not exists)
err := store.EnsureSession(id, "Chat with User")
```

## Integration

Components interact via the event bus:

```
┌─────────────┐
│   Runtime    │
│             │ Publishes: agent:spawned, agent:completed
│   Event Bus  │─────────┤
│             │ Receives: channel:message, error:occurred
│             │
└─────────────┘
     ↓
┌─────────────┐
│   Session    │
│   Store     │ Reads/writes session data
└─────────────┘
```

## Event Flow Example

1. User sends message via TUI
2. Runtime creates session in store
3. Runtime publishes `session:created` event
4. Channel manager receives event and updates status
5. Agent processes message
6. Runtime publishes `agent:completed` event
7. Session updated in store
8. Runtime publishes `session:updated` event

## Best Practices

1. **Event Naming**: Use descriptive event types with consistent prefixes (e.g., `session:created`, `agent:spawned`)
2. **Subscriber Cleanup**: Always call closer() when done with subscription
3. **Error Handling**: Publish `EventErrorOccurred` for unrecoverable errors
4. **Session Lifecycle**: Create session when conversation starts, ensure session exists before use
5. **Backpressure**: Event bus drops subscribers if channel buffer is full
6. **Concurrent Access**: Use event bus methods safely in goroutines

## Future Enhancements

- [ ] Add session resume functionality
- [ ] Add session export/import
- [ ] Add session analytics (duration, message counts)
- [ ] Add event persistence (event log for debugging)
- [ ] Add wild-card subscriptions
- [ ] Add event filtering and transformation
- [ ] Add metrics and monitoring for event bus
