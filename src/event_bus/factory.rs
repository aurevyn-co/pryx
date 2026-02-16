use super::{EnhancedEventBus, EventBus, InMemoryEventBus};
use crate::event_bus::enhanced_in_memory::DeliveryGuarantee;
use std::sync::Arc;

/// Types of event buses that can be created
#[derive(Debug, Clone, PartialEq)]
pub enum EventBusType {
    /// Basic in-memory event bus with simple pub/sub
    InMemory,
    /// Enhanced in-memory event bus with delivery guarantees
    Enhanced,
    // Future: Redis, RabbitMQ, etc.
}

/// Configuration for creating an event bus instance
#[derive(Debug, Clone)]
pub struct EventBusConfig {
    /// Type of event bus to create
    pub bus_type: EventBusType,
    /// Size of internal buffers (currently unused but reserved for future use)
    pub buffer_size: usize,
    /// Delivery guarantee level for the event bus
    pub delivery_guarantee: DeliveryGuarantee,
}

impl Default for EventBusConfig {
    fn default() -> Self {
        Self {
            bus_type: EventBusType::InMemory,
            buffer_size: 100,
            delivery_guarantee: DeliveryGuarantee::AtMostOnce,
        }
    }
}

/// Creates an event bus instance based on the provided configuration
pub fn create_event_bus(config: &EventBusConfig) -> Arc<dyn EventBus> {
    match config.bus_type {
        EventBusType::InMemory => Arc::new(InMemoryEventBus::new()),
        EventBusType::Enhanced => Arc::new(
            EnhancedEventBus::new().with_delivery_guarantee(config.delivery_guarantee.clone()),
        ),
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::event_bus::{Event, EventBroadcaster};
    use async_trait::async_trait;
    use serde_json::json;
    use std::sync::atomic::{AtomicUsize, Ordering};
    use std::sync::Arc;

    struct TestCounterHandler {
        count: Arc<AtomicUsize>,
    }

    #[async_trait]
    impl crate::event_bus::EventHandler for TestCounterHandler {
        async fn handle(
            &self,
            _event: &Event,
        ) -> Result<(), crate::event_bus::traits::EventBusError> {
            self.count.fetch_add(1, Ordering::SeqCst);
            Ok(())
        }
    }

    #[tokio::test]
    async fn test_factory_creates_in_memory_bus() {
        let config = EventBusConfig {
            bus_type: EventBusType::InMemory,
            ..Default::default()
        };

        let bus = create_event_bus(&config);
        let broadcaster = EventBroadcaster::new(bus);

        let counter = Arc::new(AtomicUsize::new(0));
        let handler = Arc::new(TestCounterHandler {
            count: counter.clone(),
        });

        broadcaster.subscribe(handler).await.unwrap();
        broadcaster
            .publish("test", json!({"data": "value"}))
            .await
            .unwrap();

        tokio::time::sleep(tokio::time::Duration::from_millis(10)).await;

        assert_eq!(counter.load(Ordering::SeqCst), 1);
    }

    #[tokio::test]
    async fn test_factory_creates_enhanced_bus() {
        let config = EventBusConfig {
            bus_type: EventBusType::Enhanced,
            delivery_guarantee: DeliveryGuarantee::AtMostOnce,
            ..Default::default()
        };

        let bus = create_event_bus(&config);
        let broadcaster = EventBroadcaster::new(bus);

        let counter = Arc::new(AtomicUsize::new(0));
        let handler = Arc::new(TestCounterHandler {
            count: counter.clone(),
        });

        broadcaster.subscribe(handler).await.unwrap();
        broadcaster
            .publish("test", json!({"data": "value"}))
            .await
            .unwrap();

        tokio::time::sleep(tokio::time::Duration::from_millis(10)).await;

        assert_eq!(counter.load(Ordering::SeqCst), 1);
    }

    #[tokio::test]
    async fn test_config_default_values() {
        let config = EventBusConfig::default();
        assert_eq!(config.bus_type, EventBusType::InMemory);
        assert_eq!(config.buffer_size, 100);
        assert_eq!(config.delivery_guarantee, DeliveryGuarantee::AtMostOnce);
    }
}
