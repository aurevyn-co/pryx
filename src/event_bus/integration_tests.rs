//! Integration tests for the event bus system
//! These tests verify complete event flows across different system components

use crate::event_bus::{
    create_event_bus, EventBroadcaster, EventBus, EventBusConfig, EventBusType, InMemoryEventBus,
};
use async_trait::async_trait;
use serde_json::json;
use std::sync::atomic::{AtomicUsize, Ordering};
use std::sync::Arc;
use tokio::sync::Mutex;

/// Test handler that collects events for verification
struct CollectingHandler {
    events: Arc<Mutex<Vec<(String, serde_json::Value)>>>,
    call_count: Arc<AtomicUsize>,
}

#[async_trait]
impl crate::event_bus::EventHandler for CollectingHandler {
    async fn handle(
        &self,
        event: &crate::event_bus::Event,
    ) -> Result<(), crate::event_bus::traits::EventBusError> {
        let mut events = self.events.lock().await;
        events.push((event.topic.clone(), event.payload.clone()));
        self.call_count.fetch_add(1, Ordering::SeqCst);
        Ok(())
    }
}

impl CollectingHandler {
    fn new() -> Self {
        Self {
            events: Arc::new(Mutex::new(Vec::new())),
            call_count: Arc::new(AtomicUsize::new(0)),
        }
    }

    fn call_count(&self) -> usize {
        self.call_count.load(Ordering::SeqCst)
    }

    async fn events(&self) -> Vec<(String, serde_json::Value)> {
        self.events.lock().await.clone()
    }
}

/// Test for basic end-to-end event flow
#[tokio::test]
async fn test_end_to_end_event_flow() {
    let bus = Arc::new(InMemoryEventBus::new());
    let broadcaster = EventBroadcaster::new(bus);

    let handler = Arc::new(CollectingHandler::new());
    broadcaster.subscribe(handler.clone()).await.unwrap();

    // Publish an event
    broadcaster
        .publish("test.topic", json!({"message": "hello world"}))
        .await
        .unwrap();

    // Allow time for event processing
    tokio::time::sleep(tokio::time::Duration::from_millis(10)).await;

    // Verify the event was received
    assert_eq!(handler.call_count(), 1);
    let events = handler.events().await;
    assert_eq!(events.len(), 1);
    assert_eq!(events[0].0, "test.topic");
    assert_eq!(events[0].1, json!({"message": "hello world"}));
}

/// Test for cross-module communication simulation
#[tokio::test]
async fn test_cross_module_communication() {
    let bus = Arc::new(InMemoryEventBus::new());
    let broadcaster = EventBroadcaster::new(bus);

    // Create handlers for different "modules"
    let agent_handler = Arc::new(CollectingHandler::new());
    let channel_handler = Arc::new(CollectingHandler::new());
    let memory_handler = Arc::new(CollectingHandler::new());

    // Subscribe handlers to their respective topics
    broadcaster
        .subscribe_to_topics(vec!["agent.start".to_string()], agent_handler.clone())
        .await
        .unwrap();

    broadcaster
        .subscribe_to_topics(vec!["channel.message".to_string()], channel_handler.clone())
        .await
        .unwrap();

    broadcaster
        .subscribe_to_topics(vec!["memory.store".to_string()], memory_handler.clone())
        .await
        .unwrap();

    // Simulate events from different modules
    broadcaster
        .publish("agent.start", json!({"session_id": "123"}))
        .await
        .unwrap();

    broadcaster
        .publish(
            "channel.message",
            json!({"sender": "user", "content": "hello"}),
        )
        .await
        .unwrap();

    broadcaster
        .publish("memory.store", json!({"key": "test", "content": "data"}))
        .await
        .unwrap();

    // Allow time for event processing
    tokio::time::sleep(tokio::time::Duration::from_millis(15)).await;

    // Verify each handler received appropriate events
    assert_eq!(agent_handler.call_count(), 1);
    assert_eq!(channel_handler.call_count(), 1);
    assert_eq!(memory_handler.call_count(), 1);

    let agent_events = agent_handler.events().await;
    assert_eq!(agent_events[0].0, "agent.start");

    let channel_events = channel_handler.events().await;
    assert_eq!(channel_events[0].0, "channel.message");

    let memory_events = memory_handler.events().await;
    assert_eq!(memory_events[0].0, "memory.store");
}

/// Test for concurrent handler processing
#[tokio::test]
async fn test_concurrent_handler_processing() {
    let bus = Arc::new(InMemoryEventBus::new());
    let broadcaster = EventBroadcaster::new(bus);

    // Create multiple handlers
    let handlers: Vec<Arc<CollectingHandler>> =
        (0..5).map(|_| Arc::new(CollectingHandler::new())).collect();

    // Subscribe all handlers
    for handler in &handlers {
        broadcaster.subscribe(handler.clone()).await.unwrap();
    }

    // Publish multiple events
    for i in 0..3 {
        broadcaster
            .publish("concurrent.test", json!({"id": i}))
            .await
            .unwrap();
    }

    // Allow time for event processing
    tokio::time::sleep(tokio::time::Duration::from_millis(20)).await;

    // Verify all handlers received all events
    for handler in &handlers {
        assert_eq!(handler.call_count(), 3, "Handler should receive 3 events");
    }
}

/// Test for factory pattern with different configurations
#[tokio::test]
async fn test_event_bus_factory_patterns() {
    // Test InMemory configuration
    let config = EventBusConfig {
        bus_type: EventBusType::InMemory,
        ..Default::default()
    };
    let bus = create_event_bus(&config);
    let broadcaster = EventBroadcaster::new(bus);

    let handler = Arc::new(CollectingHandler::new());
    broadcaster.subscribe(handler.clone()).await.unwrap();

    broadcaster
        .publish("factory.test", json!({"source": "in_memory"}))
        .await
        .unwrap();

    tokio::time::sleep(tokio::time::Duration::from_millis(10)).await;
    assert_eq!(handler.call_count(), 1);

    // Test Enhanced configuration
    let config = EventBusConfig {
        bus_type: EventBusType::Enhanced,
        ..Default::default()
    };
    let bus = create_event_bus(&config);
    let broadcaster = EventBroadcaster::new(bus);

    let handler = Arc::new(CollectingHandler::new());
    broadcaster.subscribe(handler.clone()).await.unwrap();

    broadcaster
        .publish("factory.test", json!({"source": "enhanced"}))
        .await
        .unwrap();

    tokio::time::sleep(tokio::time::Duration::from_millis(10)).await;
    assert_eq!(handler.call_count(), 1);
}

/// Test for error propagation when handlers fail
#[tokio::test]
async fn test_error_propagation() {
    use crate::event_bus::traits::EventBusError;

    struct FailingHandler;

    #[async_trait]
    impl crate::event_bus::EventHandler for FailingHandler {
        async fn handle(&self, _event: &crate::event_bus::Event) -> Result<(), EventBusError> {
            Err(EventBusError::HandlerError(
                "Simulated handler failure".to_string(),
            ))
        }
    }

    let bus = Arc::new(InMemoryEventBus::new());
    let broadcaster = EventBroadcaster::new(bus);

    // Add a failing handler and a normal handler
    let failing_handler = Arc::new(FailingHandler);
    let normal_handler = Arc::new(CollectingHandler::new());

    broadcaster.subscribe(failing_handler).await.unwrap();
    broadcaster.subscribe(normal_handler.clone()).await.unwrap();

    // Publish an event - should still reach the normal handler despite the failing one
    broadcaster
        .publish("error.propagation", json!({"test": true}))
        .await
        .unwrap();

    tokio::time::sleep(tokio::time::Duration::from_millis(10)).await;

    // Normal handler should still receive the event
    assert_eq!(normal_handler.call_count(), 1);
}

/// Test for resource cleanup
#[tokio::test]
async fn test_resource_cleanup() {
    let bus = Arc::new(InMemoryEventBus::new());
    let broadcaster = EventBroadcaster::new(bus);

    let handler = Arc::new(CollectingHandler::new());
    let _handler_count_before = handler.call_count();

    broadcaster.subscribe(handler.clone()).await.unwrap();

    // Publish an event
    broadcaster
        .publish("cleanup.test", json!({"step": 1}))
        .await
        .unwrap();

    tokio::time::sleep(tokio::time::Duration::from_millis(5)).await;

    let events_after_first = handler.events().await.len();

    // Publish another event
    broadcaster
        .publish("cleanup.test", json!({"step": 2}))
        .await
        .unwrap();

    tokio::time::sleep(tokio::time::Duration::from_millis(5)).await;

    let events_after_second = handler.events().await.len();

    // Verify events were processed correctly
    assert_eq!(events_after_first, 1);
    assert_eq!(events_after_second, 2);
}

/// Test for topic validation
#[tokio::test]
async fn test_topic_validation() {
    use crate::event_bus::traits::EventBusError;

    let bus = Arc::new(InMemoryEventBus::new());
    let broadcaster = EventBroadcaster::new(bus);

    // Attempt to subscribe with an empty topic (should fail)
    let result = broadcaster
        .subscribe_to_topics(vec!["".to_string()], Arc::new(CollectingHandler::new()))
        .await;

    assert!(matches!(result, Err(EventBusError::InvalidTopic(_))));
}

/// Test for large payload handling
#[tokio::test]
async fn test_large_payload_handling() {
    let bus = Arc::new(InMemoryEventBus::new());
    let broadcaster = EventBroadcaster::new(bus);

    let handler = Arc::new(CollectingHandler::new());
    broadcaster.subscribe(handler.clone()).await.unwrap();

    // Create a large payload
    let large_data = "x".repeat(10000); // 10KB string
    broadcaster
        .publish("large.payload", json!({"data": large_data}))
        .await
        .unwrap();

    tokio::time::sleep(tokio::time::Duration::from_millis(10)).await;

    assert_eq!(handler.call_count(), 1);
    let events = handler.events().await;
    assert_eq!(events[0].0, "large.payload");
}
