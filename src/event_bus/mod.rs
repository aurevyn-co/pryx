//! Event bus system for inter-component communication in Pryx.
//!
//! This module provides a centralized pub/sub system that allows different
//! components of Pryx to communicate with each other without tight coupling.

pub mod broadcaster;
pub mod enhanced_in_memory;
pub mod example_daemon;
pub mod factory;
pub mod in_memory;
#[cfg(test)]
pub mod integration_tests;
#[cfg(test)]
pub mod performance_tests;
pub mod traits;

pub use broadcaster::EventBroadcaster;
pub use enhanced_in_memory::InMemoryEventBus as EnhancedEventBus;
pub use factory::{create_event_bus, EventBusConfig, EventBusType};
pub use in_memory::InMemoryEventBus;
pub use traits::{Event, EventBus, EventHandler};
