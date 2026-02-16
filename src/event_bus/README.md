# Event Bus System

The Event Bus System is a core component of Pryx that enables loose coupling between different modules and components. It implements a publish-subscribe pattern that allows various parts of the system to communicate without direct dependencies.

## Overview

The event bus provides:
- **Decoupling**: Components don't need to know about each other directly
- **Scalability**: Easy to add new event producers and consumers
- **Flexibility**: Support for topic-based routing
- **Asynchronous Processing**: Non-blocking event handling

## Architecture

```
[Producer] ----> [Event Bus] ----> [Consumer]
[Producer] ---->    |         ----> [Consumer]
                   |
[Producer] --------> |         ----> [Consumer]
```

## Components

### Event
Represents a message that can be published to the bus:
- `topic`: String identifier for routing
- `payload`: JSON data containing the message content
- `timestamp`: When the event was created

### EventBus Trait
Core interface defining:
- `publish()`: Send an event to the bus
- `subscribe()`: Register a handler for events
- `subscribe_to_topics()`: Register a handler for specific topics
- `stats()`: Get statistics about the bus

### EventHandler Trait
Interface for processing events:
- `handle()`: Process an incoming event
- `topics()`: Optional filter for specific topics

## Usage Examples

### Basic Publishing
```rust
use pryx::event_bus::{InMemoryEventBus, EventBroadcaster};
use serde_json::json;
use std::sync::Arc;

let event_bus = Arc::new(InMemoryEventBus::new());
let broadcaster = EventBroadcaster::new(event_bus);

// Publish an event
broadcaster
    .publish("user.login", json!({"user_id": 123}))
    .await?;
```

### Creating an Event Handler
```rust
use pryx::event_bus::EventHandler;
use async_trait::async_trait;

struct LoginHandler;

#[async_trait]
impl EventHandler for LoginHandler {
    async fn handle(&self, event: &Event) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
        if event.topic == "user.login" {
            println!("User logged in: {}", event.payload);
        }
        Ok(())
    }
}
```

### Subscribing to Events
```rust
use std::sync::Arc;

let handler = Arc::new(LoginHandler);
broadcaster.subscribe(handler).await?;
```

### Subscribing to Specific Topics
```rust
// Subscribe only to specific topics
broadcaster
    .subscribe_to_topics(
        vec!["user.login".to_string(), "user.logout".to_string()],
        handler
    )
    .await?;
```

## Integration Points

The event bus can be integrated into various parts of the system:

- **Agent Module**: Publish events for agent start/end, responses
- **Channel Module**: Emit events for incoming/outgoing messages
- **Memory Module**: Notify when memories are stored/recalled
- **Health Module**: Report component health status
- **Security Module**: Alert on security-related events

## Benefits

1. **Loose Coupling**: Components can evolve independently
2. **Extensibility**: Easy to add new event consumers
3. **Monitoring**: Centralized event stream for observability
4. **Reliability**: Asynchronous processing prevents blocking
5. **Maintainability**: Clear separation of concerns

## Performance Considerations

- Events are processed asynchronously using Tokio tasks
- Topic-based filtering reduces unnecessary processing
- Memory usage scales with active subscriptions
- Backpressure is handled through bounded channels

## Future Enhancements

- Persistent event storage
- Cross-process event distribution
- Event replay capabilities
- Advanced filtering and transformation