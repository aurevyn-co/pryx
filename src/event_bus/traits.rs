use async_trait::async_trait;
use std::sync::Arc;

/// Represents an event that can be published to the event bus
#[derive(Debug, Clone)]
pub struct Event {
    pub topic: String,
    pub payload: serde_json::Value,
    pub timestamp: std::time::SystemTime,
}

impl Event {
    pub fn new(topic: impl Into<String>, payload: serde_json::Value) -> Self {
        Self {
            topic: topic.into(),
            payload,
            timestamp: std::time::SystemTime::now(),
        }
    }
}

/// Handler for processing events
#[async_trait]
pub trait EventHandler: Send + Sync {
    async fn handle(&self, event: &Event) -> Result<(), Box<dyn std::error::Error + Send + Sync>>;

    /// Optional method to specify which topics this handler is interested in
    /// Return None to receive all events, or Some(Vec<String>) to receive only specific topics
    fn topics(&self) -> Option<Vec<String>> {
        None
    }
}

/// Core event bus trait - implement for any event bus backend
#[async_trait]
pub trait EventBus: Send + Sync {
    /// Publish an event to the bus
    async fn publish(&self, event: Event) -> Result<(), Box<dyn std::error::Error + Send + Sync>>;

    /// Subscribe to events with a handler
    async fn subscribe(
        &self,
        handler: Arc<dyn EventHandler>,
    ) -> Result<(), Box<dyn std::error::Error + Send + Sync>>;

    /// Subscribe to specific topics with a handler
    async fn subscribe_to_topics(
        &self,
        topics: Vec<String>,
        handler: Arc<dyn EventHandler>,
    ) -> Result<(), Box<dyn std::error::Error + Send + Sync>>;

    /// Get statistics about the event bus
    async fn stats(&self) -> EventBusStats;
}

/// Statistics about the event bus
#[derive(Debug, Clone)]
pub struct EventBusStats {
    pub total_events_published: u64,
    pub total_events_delivered: u64,
    pub active_handlers: usize,
    pub queued_events: usize,
}
