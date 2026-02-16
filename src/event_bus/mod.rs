//! Event bus system for inter-component communication in Pryx.
//!
//! This module provides a centralized pub/sub system that allows different
//! components of Pryx to communicate with each other without tight coupling.

pub mod traits;
pub mod in_memory;
pub mod enhanced_in_memory;
pub mod broadcaster;
pub mod example_daemon;
pub mod factory;

pub use traits::{EventBus, Event, EventHandler};
pub use in_memory::InMemoryEventBus;
pub use enhanced_in_memory::InMemoryEventBus as EnhancedEventBus;
pub use broadcaster::EventBroadcaster;
pub use factory::{create_event_bus, EventBusConfig, EventBusType};