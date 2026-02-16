# Event Bus System Documentation

## Overview

The Event Bus System is a core component of Pryx that enables loose coupling between different modules and components. It implements a publish-subscribe pattern that allows various parts of the system to communicate without direct dependencies.

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
    async fn handle(&self, event: &Event) -> Result<(), EventBusError> {
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

## Delivery Guarantees

The system supports different delivery guarantee levels:

- **AtMostOnce** (default): Fire-and-forget delivery
- **AtLeastOnce**: Delivery with retry mechanism (requires persistent storage)
- **ExactlyOnce**: Delivery with deduplication (requires persistent storage)

## Factory Pattern

The system includes a factory for creating different types of event buses:

```rust
use pryx::event_bus::{EventBusType, EventBusConfig, create_event_bus};

let config = EventBusConfig {
    bus_type: EventBusType::Enhanced,
    delivery_guarantee: DeliveryGuarantee::AtMostOnce,
    ..Default::default()
};

let bus = create_event_bus(&config);
```

## Performance Characteristics

- High throughput: >1000 events/sec
- Low latency: <10ms average
- Scalable: Works with many subscribers
- Memory efficient: Minimal overhead per event

## Testing

The system includes comprehensive tests:
- Integration tests for end-to-end flows
- Performance tests for throughput and latency
- Error handling tests
- Resource cleanup tests

## Benefits

1. **Loose Coupling**: Components can evolve independently
2. **Extensibility**: Easy to add new event consumers
3. **Monitoring**: Centralized event stream for observability
4. **Reliability**: Asynchronous processing prevents blocking
5. **Maintainability**: Clear separation of concerns