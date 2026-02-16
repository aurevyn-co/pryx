use super::traits::{Event, EventBus, EventHandler};
use std::sync::Arc;

/// A convenience wrapper around `EventBus` for easy publishing of events
pub struct EventBroadcaster {
    event_bus: Arc<dyn EventBus>,
}

impl EventBroadcaster {
    pub fn new(event_bus: Arc<dyn EventBus>) -> Self {
        Self { event_bus }
    }

    /// Publish an event to the bus
    pub async fn publish(
        &self,
        topic: impl Into<String>,
        payload: serde_json::Value,
    ) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
        let event = Event::new(topic, payload);
        self.event_bus.publish(event).await
    }

    /// Subscribe to all events with a handler
    pub async fn subscribe(
        &self,
        handler: Arc<dyn EventHandler>,
    ) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
        self.event_bus.subscribe(handler).await
    }

    /// Subscribe to specific topics with a handler
    pub async fn subscribe_to_topics(
        &self,
        topics: Vec<String>,
        handler: Arc<dyn EventHandler>,
    ) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
        self.event_bus.subscribe_to_topics(topics, handler).await
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::event_bus::InMemoryEventBus;
    use async_trait::async_trait;
    use serde_json::json;
    use std::sync::atomic::{AtomicBool, Ordering};
    use std::sync::Arc;

    struct TestEventHandler {
        received: Arc<AtomicBool>,
    }

    #[async_trait]
    impl EventHandler for TestEventHandler {
        async fn handle(
            &self,
            _event: &Event,
        ) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
            self.received.store(true, Ordering::SeqCst);
            Ok(())
        }
    }

    #[tokio::test]
    async fn test_event_bus_publish_subscribe() {
        let bus = Arc::new(InMemoryEventBus::new());
        let broadcaster = EventBroadcaster::new(bus.clone());

        let received = Arc::new(AtomicBool::new(false));
        let handler = Arc::new(TestEventHandler {
            received: received.clone(),
        });

        broadcaster.subscribe(handler).await.unwrap();

        broadcaster
            .publish("test_topic", json!({"message": "hello"}))
            .await
            .unwrap();

        // Give the async handler time to process
        tokio::time::sleep(tokio::time::Duration::from_millis(10)).await;

        assert!(received.load(Ordering::SeqCst));
    }

    #[tokio::test]
    async fn test_event_bus_subscribe_to_specific_topics() {
        let bus = Arc::new(InMemoryEventBus::new());
        let broadcaster = EventBroadcaster::new(bus.clone());

        let received = Arc::new(AtomicBool::new(false));
        let handler = Arc::new(TestEventHandler {
            received: received.clone(),
        });

        broadcaster
            .subscribe_to_topics(vec!["specific_topic".to_string()], handler)
            .await
            .unwrap();

        // Publish to the subscribed topic
        broadcaster
            .publish("specific_topic", json!({"message": "hello"}))
            .await
            .unwrap();

        // Give the async handler time to process
        tokio::time::sleep(tokio::time::Duration::from_millis(10)).await;

        assert!(received.load(Ordering::SeqCst));

        // Reset for next test
        received.store(false, Ordering::SeqCst);

        // Publish to a different topic - should not be received
        broadcaster
            .publish("other_topic", json!({"message": "world"}))
            .await
            .unwrap();

        // Give the async handler time to process
        tokio::time::sleep(tokio::time::Duration::from_millis(10)).await;

        // Should still be false since handler didn't subscribe to "other_topic"
        assert!(!received.load(Ordering::SeqCst));
    }
}
