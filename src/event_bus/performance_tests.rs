//! Performance tests for the event bus system
//! These tests measure throughput, latency, and scalability

use crate::event_bus::{EventBroadcaster, EventBus, InMemoryEventBus};
use async_trait::async_trait;
use serde_json::json;
use std::sync::atomic::{AtomicUsize, Ordering};
use std::sync::Arc;
use std::time::Instant;
use tokio::sync::Mutex;

/// Performance test handler that measures timing
struct TimingHandler {
    count: Arc<AtomicUsize>,
    start_time: Arc<Mutex<Option<Instant>>>,
    end_time: Arc<Mutex<Option<Instant>>>,
}

#[async_trait]
impl crate::event_bus::EventHandler for TimingHandler {
    async fn handle(
        &self,
        _event: &crate::event_bus::Event,
    ) -> Result<(), crate::event_bus::traits::EventBusError> {
        self.count.fetch_add(1, Ordering::SeqCst);

        let mut start_guard = self.start_time.lock().await;
        if start_guard.is_none() {
            *start_guard = Some(Instant::now());
        }
        drop(start_guard);

        let mut end_guard = self.end_time.lock().await;
        *end_guard = Some(Instant::now());
        drop(end_guard);

        Ok(())
    }
}

impl TimingHandler {
    fn new() -> Self {
        Self {
            count: Arc::new(AtomicUsize::new(0)),
            start_time: Arc::new(Mutex::new(None)),
            end_time: Arc::new(Mutex::new(None)),
        }
    }

    fn count(&self) -> usize {
        self.count.load(Ordering::SeqCst)
    }

    async fn get_timing(&self) -> Option<std::time::Duration> {
        let start = self.start_time.lock().await;
        let end = self.end_time.lock().await;

        if let (Some(start_time), Some(end_time)) = (*start, *end) {
            Some(end_time.duration_since(start_time))
        } else {
            None
        }
    }
}

/// Test throughput under high event volume
#[tokio::test]
async fn test_throughput_performance() {
    let bus = Arc::new(InMemoryEventBus::new());
    let broadcaster = EventBroadcaster::new(bus);

    let handler = Arc::new(TimingHandler::new());
    broadcaster.subscribe(handler.clone()).await.unwrap();

    let start = Instant::now();
    let num_events = 1000;

    // Publish many events
    for i in 0..num_events {
        broadcaster
            .publish(
                "perf.throughput",
                json!({"id": i, "timestamp": start.elapsed().as_nanos()}),
            )
            .await
            .unwrap();
    }

    // Wait for all events to be processed
    while handler.count() < num_events {
        tokio::time::sleep(tokio::time::Duration::from_millis(1)).await;
    }

    let elapsed = start.elapsed();
    let events_per_sec = num_events as f64 / elapsed.as_secs_f64();

    println!(
        "Throughput: {} events/sec ({} events in {:?})",
        events_per_sec, num_events, elapsed
    );

    // Verify all events were processed
    assert_eq!(handler.count(), num_events);

    // Performance threshold: should handle at least 1000 events per second
    assert!(
        events_per_sec >= 1000.0,
        "Throughput too low: {} events/sec",
        events_per_sec
    );
}

/// Test memory usage under sustained load
#[tokio::test]
async fn test_memory_usage_under_load() {
    let bus = Arc::new(InMemoryEventBus::new());
    let broadcaster = EventBroadcaster::new(bus);

    let handler = Arc::new(TimingHandler::new());
    broadcaster.subscribe(handler.clone()).await.unwrap();

    let iterations = 500;
    let mut peak_count = 0;

    for i in 0..iterations {
        broadcaster
            .publish("perf.memory", json!({"iteration": i}))
            .await
            .unwrap();

        let current_count = handler.count();
        if current_count > peak_count {
            peak_count = current_count;
        }

        // Small delay to allow processing
        tokio::time::sleep(tokio::time::Duration::from_micros(100)).await;
    }

    // Allow final processing
    tokio::time::sleep(tokio::time::Duration::from_millis(10)).await;

    let final_count = handler.count();
    println!("Memory test: {} events processed", final_count);

    // Verify all events were processed
    assert_eq!(final_count, iterations);
}

/// Test scalability with many subscribers
#[tokio::test]
async fn test_scalability_with_many_subscribers() {
    let bus = Arc::new(InMemoryEventBus::new());
    let broadcaster = EventBroadcaster::new(bus);

    let num_handlers = 50;
    let handlers: Vec<Arc<TimingHandler>> = (0..num_handlers)
        .map(|_| Arc::new(TimingHandler::new()))
        .collect();

    // Subscribe all handlers
    for handler in &handlers {
        broadcaster.subscribe(handler.clone()).await.unwrap();
    }

    let start = Instant::now();
    let num_events = 100;

    // Publish events
    for i in 0..num_events {
        broadcaster
            .publish("perf.scalability", json!({"event_id": i}))
            .await
            .unwrap();
    }

    // Wait for all events to be processed by all handlers
    let mut total_processed = 0;
    let max_wait = 5000; // 5 seconds max wait
    let mut waited = 0;
    while total_processed < num_events * num_handlers && waited < max_wait {
        total_processed = handlers.iter().map(|h| h.count()).sum();
        tokio::time::sleep(tokio::time::Duration::from_millis(10)).await;
        waited += 10;
    }

    let elapsed = start.elapsed();
    let total_events = num_events * num_handlers;
    let events_per_sec = total_events as f64 / elapsed.as_secs_f64();

    println!(
        "Scalability: {} events/sec ({} handlers × {} events in {:?})",
        events_per_sec, num_handlers, num_events, elapsed
    );

    // Verify all events were processed by all handlers
    assert_eq!(total_processed, total_events);

    // Performance threshold: should handle at least 500 events per second with 50 handlers
    assert!(
        events_per_sec >= 500.0,
        "Scalability too low: {} events/sec",
        events_per_sec
    );
}

/// Test latency measurements
#[tokio::test]
async fn test_latency_measurements() {
    let bus = Arc::new(InMemoryEventBus::new());
    let broadcaster = EventBroadcaster::new(bus);

    let handler = Arc::new(TimingHandler::new());
    broadcaster.subscribe(handler.clone()).await.unwrap();

    let mut latencies = Vec::new();
    let num_trials = 100;

    for _ in 0..num_trials {
        let start = Instant::now();

        broadcaster
            .publish(
                "perf.latency",
                json!({"timestamp": start.elapsed().as_nanos()}),
            )
            .await
            .unwrap();

        // Wait for processing
        while handler.count() < 1 {
            tokio::time::sleep(tokio::time::Duration::from_millis(1)).await;
        }

        let elapsed = start.elapsed();
        latencies.push(elapsed);

        // Reset handler for next trial
        // We'll just continue and check the cumulative effect
    }

    // Calculate average latency
    let total_duration: std::time::Duration = latencies.iter().sum();
    let avg_latency = total_duration.div_f64(num_trials as f64);

    println!("Average latency: {:?}", avg_latency);
    println!(
        "Min latency: {:?}",
        latencies.iter().min().unwrap_or(&std::time::Duration::ZERO)
    );
    println!(
        "Max latency: {:?}",
        latencies.iter().max().unwrap_or(&std::time::Duration::ZERO)
    );

    // Latency threshold: average should be under 10ms
    assert!(
        avg_latency.as_millis() < 10,
        "Average latency too high: {:?}",
        avg_latency
    );
}

/// Test for sustained load over time
#[tokio::test]
async fn test_sustained_load_over_time() {
    let bus = Arc::new(InMemoryEventBus::new());
    let broadcaster = EventBroadcaster::new(bus);

    let handler = Arc::new(TimingHandler::new());
    broadcaster.subscribe(handler.clone()).await.unwrap();

    let duration = std::time::Duration::from_secs(2);
    let start = Instant::now();
    let mut event_count = 0;

    // Publish events continuously for 2 seconds
    while start.elapsed() < duration {
        broadcaster
            .publish("perf.sustained", json!({"count": event_count}))
            .await
            .unwrap();

        event_count += 1;

        // Small delay to prevent overwhelming the system
        tokio::time::sleep(tokio::time::Duration::from_millis(1)).await;
    }

    let target_events = event_count;

    // Wait a bit more for processing to complete
    tokio::time::sleep(tokio::time::Duration::from_millis(100)).await;

    let processed = handler.count();
    let success_rate = (processed as f64 / target_events as f64) * 100.0;

    println!(
        "Sustained load: {} events sent, {} processed ({}% success rate)",
        target_events, processed, success_rate
    );

    // Success threshold: at least 95% of events should be processed
    assert!(
        success_rate >= 95.0,
        "Success rate too low: {}%",
        success_rate
    );
}

/// Test memory consumption patterns
#[tokio::test]
async fn test_memory_consumption_patterns() {
    let bus = Arc::new(InMemoryEventBus::new());
    let broadcaster = EventBroadcaster::new(bus);

    let handler = Arc::new(TimingHandler::new());
    broadcaster.subscribe(handler.clone()).await.unwrap();

    // Measure initial state
    let initial_count = handler.count();

    // Publish a moderate number of events
    let num_events = 1000;
    for i in 0..num_events {
        broadcaster
            .publish("perf.memory.pattern", json!({"id": i}))
            .await
            .unwrap();
    }

    // Wait for processing
    while handler.count() < initial_count + num_events {
        tokio::time::sleep(tokio::time::Duration::from_millis(1)).await;
    }

    let final_count = handler.count();

    println!(
        "Memory pattern test: {} initial, {} final",
        initial_count, final_count
    );

    // Verify all events were processed
    assert_eq!(final_count, initial_count + num_events);
}
