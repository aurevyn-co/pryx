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

/// Specific error types for the event bus system
#[derive(Debug, thiserror::Error)]
pub enum EventBusError {
    /// Error occurred while publishing an event
    #[error("Publish error: {0}")]
    PublishError(String),

    /// Error occurred while subscribing to events
    #[error("Subscribe error: {0}")]
    SubscribeError(String),

    /// Error occurred in an event handler
    #[error("Handler error: {0}")]
    HandlerError(String),

    /// Invalid topic provided
    #[error("Invalid topic: {0}")]
    InvalidTopic(String),
}

/// Handler for processing events
#[async_trait]
pub trait EventHandler: Send + Sync {
    /// Handle an incoming event
    ///
    /// Errors returned from this method will be logged by the event bus implementation
    /// but will not stop other handlers from receiving the event.
    async fn handle(&self, event: &Event) -> Result<(), EventBusError>;

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
    ///
    /// This method should return an error only if the event could not be queued for delivery.
    /// Actual delivery failures should be handled internally by the implementation.
    async fn publish(&self, event: Event) -> Result<(), EventBusError>;

    /// Subscribe to events with a handler
    ///
    /// The handler will receive events from all topics.
    async fn subscribe(&self, handler: Arc<dyn EventHandler>) -> Result<(), EventBusError>;

    /// Subscribe to specific topics with a handler
    ///
    /// The handler will only receive events from the specified topics.
    async fn subscribe_to_topics(
        &self,
        topics: Vec<String>,
        handler: Arc<dyn EventHandler>,
    ) -> Result<(), EventBusError>;

    /// Get statistics about the event bus
    ///
    /// Statistics are updated asynchronously and may not reflect the most recent state.
    async fn stats(&self) -> EventBusStats;
}

/// Statistics about the event bus
#[derive(Debug, Clone)]
pub struct EventBusStats {
    /// Total number of events published to the bus
    pub total_events_published: u64,

    /// Total number of events delivered to handlers
    pub total_events_delivered: u64,

    /// Number of currently active handlers (may be approximate)
    pub active_handlers: usize,

    /// Number of events currently queued for delivery
    pub queued_events: usize,
}
