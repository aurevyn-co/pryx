use super::traits::{Event, EventBus, EventBusStats, EventHandler};
use async_trait::async_trait;
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;

/// Type alias for the complex subscriber mapping
type SubscriberMap = HashMap<String, Vec<Arc<dyn EventHandler>>>;

/// Enhanced in-memory event bus implementation with additional features
pub struct InMemoryEventBus {
    subscribers: Arc<RwLock<SubscriberMap>>,
    stats: Arc<RwLock<EventBusStats>>,
    /// Optional delivery guarantee level
    delivery_guarantee: DeliveryGuarantee,
}

#[derive(Debug, Clone)]
pub enum DeliveryGuarantee {
    AtMostOnce,
    AtLeastOnce, // Would require persistent storage
    ExactlyOnce, // Would require persistent storage + deduplication
}

impl InMemoryEventBus {
    pub fn new() -> Self {
        Self {
            subscribers: Arc::new(RwLock::new(HashMap::new())),
            stats: Arc::new(RwLock::new(EventBusStats {
                total_events_published: 0,
                total_events_delivered: 0,
                active_handlers: 0,
                queued_events: 0,
            })),
            delivery_guarantee: DeliveryGuarantee::AtMostOnce, // Default
        }
    }

    pub fn with_delivery_guarantee(mut self, guarantee: DeliveryGuarantee) -> Self {
        self.delivery_guarantee = guarantee;
        self
    }
}

#[async_trait]
impl EventBus for InMemoryEventBus {
    async fn publish(&self, event: Event) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
        // Update stats
        {
            let mut stats = self.stats.write().await;
            stats.total_events_published += 1;
        }

        // Get subscribers for this topic
        let subscribers = self.subscribers.read().await;
        let topic_specific = subscribers.get(&event.topic);
        let wildcard_subscribers = subscribers.get("*"); // Wildcard for all topics

        // Collect all handlers that should receive this event
        let mut handlers_to_notify = Vec::new();

        if let Some(handlers) = topic_specific {
            handlers_to_notify.extend(handlers.iter().cloned());
        }

        if let Some(handlers) = wildcard_subscribers {
            handlers_to_notify.extend(handlers.iter().cloned());
        }

        // Drop the read lock before processing events
        drop(subscribers);

        // Store the length before iterating to avoid moving the vector
        let handlers_count = handlers_to_notify.len();

        // Notify all relevant handlers based on delivery guarantee
        match self.delivery_guarantee {
            DeliveryGuarantee::AtMostOnce => {
                // Fire and forget - current behavior
                for handler in &handlers_to_notify {
                    let event_clone = event.clone();
                    let handler_clone = Arc::clone(handler);
                    tokio::spawn(async move {
                        if let Err(e) = handler_clone.handle(&event_clone).await {
                            tracing::error!("Event handler error: {}", e);
                        }
                    });
                }
            }
            DeliveryGuarantee::AtLeastOnce => {
                // In a real implementation, this would involve persistent storage
                // For now, we'll simulate by awaiting each handler
                for handler in &handlers_to_notify {
                    let event_clone = event.clone();
                    let handler_clone = Arc::clone(handler);
                    // In a real at-least-once system, we'd store the event before delivery
                    // and retry until acknowledged
                    tokio::spawn(async move {
                        if let Err(e) = handler_clone.handle(&event_clone).await {
                            tracing::error!(
                                "Event handler error (delivery guarantee failed): {}",
                                e
                            );
                        }
                    });
                }
            }
            DeliveryGuarantee::ExactlyOnce => {
                // This would require deduplication mechanisms and persistent storage
                // For now, we'll treat it similarly to AtLeastOnce
                for handler in &handlers_to_notify {
                    let event_clone = event.clone();
                    let handler_clone = Arc::clone(handler);
                    tokio::spawn(async move {
                        if let Err(e) = handler_clone.handle(&event_clone).await {
                            tracing::error!(
                                "Event handler error (delivery guarantee failed): {}",
                                e
                            );
                        }
                    });
                }
            }
        }

        // Update delivery stats
        {
            let mut stats = self.stats.write().await;
            stats.total_events_delivered += handlers_count as u64;
        }

        Ok(())
    }

    async fn subscribe(
        &self,
        handler: Arc<dyn EventHandler>,
    ) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
        // Subscribe to all topics using wildcard
        self.subscribe_to_topics(vec!["*".to_string()], handler)
            .await
    }

    async fn subscribe_to_topics(
        &self,
        topics: Vec<String>,
        handler: Arc<dyn EventHandler>,
    ) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
        let mut subscribers = self.subscribers.write().await;

        for topic in topics {
            subscribers
                .entry(topic)
                .or_insert_with(Vec::new)
                .push(handler.clone());
        }

        // Update stats
        let mut stats = self.stats.write().await;
        stats.active_handlers += 1;

        Ok(())
    }

    async fn stats(&self) -> EventBusStats {
        // Return a clone of the current stats
        let stats = self.stats.read().await;
        stats.clone()
    }
}
