//! Example of how to integrate the event bus into the daemon module
//!
//! This shows how different components can communicate through the event bus

use crate::config::Config;
use crate::event_bus::traits::EventBusError;
use crate::event_bus::{Event, EventBroadcaster, EventBus, InMemoryEventBus};
use anyhow::Result;
use async_trait::async_trait;
use serde_json::json;
use std::sync::Arc;

/// Example event handler that logs system events
pub struct SystemLoggerHandler;

#[async_trait]
impl crate::event_bus::EventHandler for SystemLoggerHandler {
    async fn handle(&self, event: &Event) -> Result<(), EventBusError> {
        match event.topic.as_str() {
            "agent.start" => {
                tracing::info!("Agent started: {}", event.payload);
            }
            "agent.end" => {
                tracing::info!("Agent ended: {}", event.payload);
            }
            "channel.message" => {
                tracing::info!("Channel message received: {}", event.payload);
            }
            "memory.store" => {
                tracing::info!("Memory stored: {}", event.payload);
            }
            "memory.recall" => {
                tracing::info!("Memory recalled: {}", event.payload);
            }
            _ => {
                tracing::debug!(
                    "Received event on topic '{}': {}",
                    event.topic,
                    event.payload
                );
            }
        }
        Ok(())
    }
}

/// Example event handler that monitors for errors
pub struct ErrorHandler;

#[async_trait]
impl crate::event_bus::EventHandler for ErrorHandler {
    async fn handle(&self, event: &Event) -> Result<(), EventBusError> {
        if event.topic == "system.error" {
            tracing::error!("System error occurred: {}", event.payload);
            // Could trigger alerting, recovery mechanisms, etc.
        }
        Ok(())
    }
}

/// Example of how to initialize and use the event bus in a system component
pub async fn setup_event_bus_example(config: &Config) -> anyhow::Result<()> {
    // Create the event bus
    let event_bus: Arc<dyn EventBus> = Arc::new(InMemoryEventBus::new());
    let broadcaster = EventBroadcaster::new(event_bus.clone());

    // Register event handlers
    let logger_handler = Arc::new(SystemLoggerHandler);
    let error_handler = Arc::new(ErrorHandler);

    // Convert the error types to anyhow::Error
    match broadcaster.subscribe(logger_handler).await {
        Ok(()) => (),
        Err(EventBusError::SubscribeError(e)) => return Err(anyhow::anyhow!(e)),
        Err(e) => return Err(anyhow::anyhow!("{e}")),
    }

    match broadcaster.subscribe(error_handler).await {
        Ok(()) => (),
        Err(EventBusError::SubscribeError(e)) => return Err(anyhow::anyhow!(e)),
        Err(e) => return Err(anyhow::anyhow!("{e}")),
    }

    // Example: Publish some events
    match broadcaster
        .publish(
            "system.startup",
            json!({
                "timestamp": chrono::Utc::now().to_rfc3339(),
                "version": env!("CARGO_PKG_VERSION"),
                "config_path": config.config_path.display().to_string()
            }),
        )
        .await
    {
        Ok(()) => (),
        Err(EventBusError::PublishError(e)) => return Err(anyhow::anyhow!(e)),
        Err(e) => return Err(anyhow::anyhow!("{e}")),
    }

    match broadcaster
        .publish(
            "agent.ready",
            json!({
                "ready": true,
                "provider": config.default_provider.as_deref().unwrap_or("unknown"),
                "model": config.default_model.as_deref().unwrap_or("unknown")
            }),
        )
        .await
    {
        Ok(()) => (),
        Err(EventBusError::PublishError(e)) => return Err(anyhow::anyhow!(e)),
        Err(e) => return Err(anyhow::anyhow!("{e}")),
    }

    tracing::info!("Event bus initialized and example events published");

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde_json::json;

    #[tokio::test]
    async fn test_daemon_event_bus_integration() {
        let config = Config::default();
        let result = setup_event_bus_example(&config).await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_event_bus_with_specific_handlers() {
        let event_bus: Arc<dyn EventBus> = Arc::new(InMemoryEventBus::new());
        let broadcaster = EventBroadcaster::new(event_bus.clone());

        // Create a handler that only listens to specific topics
        struct SpecificHandler {
            received_events: Arc<tokio::sync::Mutex<Vec<String>>>,
        }

        #[async_trait]
        impl crate::event_bus::EventHandler for SpecificHandler {
            async fn handle(&self, event: &Event) -> Result<(), EventBusError> {
                let mut events = self.received_events.lock().await;
                events.push(event.topic.clone());
                Ok(())
            }

            fn topics(&self) -> Option<Vec<String>> {
                Some(vec!["test.specific".to_string()])
            }
        }

        let received_events = Arc::new(tokio::sync::Mutex::new(Vec::new()));
        let specific_handler = Arc::new(SpecificHandler {
            received_events: received_events.clone(),
        });

        // Subscribe to specific topics
        broadcaster
            .subscribe_to_topics(vec!["test.specific".to_string()], specific_handler)
            .await
            .unwrap();

        // Publish events - only the specific topic should be handled
        broadcaster
            .publish("test.other", json!({"data": "other"}))
            .await
            .unwrap();

        broadcaster
            .publish("test.specific", json!({"data": "specific"}))
            .await
            .unwrap();

        // Give handlers time to process
        tokio::time::sleep(tokio::time::Duration::from_millis(10)).await;

        let events = received_events.lock().await;
        assert_eq!(events.len(), 1);
        assert_eq!(events[0], "test.specific");
    }
}
