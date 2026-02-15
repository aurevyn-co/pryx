package bus

import (
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
)

// Handler is a function that handles an event.
// It receives the event as its only parameter.
type Handler func(Event)

// Subscription represents a subscription to the event bus.
// It contains the channel for receiving events and metadata about the subscription.
type Subscription struct {
	id     string
	ch     chan Event
	topics []EventType
	closer func()
}

// Bus is the central event bus for pub/sub communication.
// It manages subscriptions and routes events to interested subscribers.
// The bus is safe for concurrent use.
type Bus struct {
	mu   sync.RWMutex
	subs map[string]*Subscription
	ver  int64
}

// New creates a new event Bus with no subscribers.
// Use Subscribe to add subscribers and Publish to send events.
func New() *Bus {
	return &Bus{
		subs: make(map[string]*Subscription),
	}
}

// Subscribe subscribes to events. If topics is empty, it subscribes to all events.
// Returns a channel that receives events. The bus owns the channel; use the closer to unsubscribe.
func (b *Bus) Subscribe(topics ...EventType) (<-chan Event, func()) {
	b.mu.Lock()
	defer b.mu.Unlock()

	id := uuid.New().String()
	ch := make(chan Event, 100) // Buffer events

	sub := &Subscription{
		id:     id,
		ch:     ch,
		topics: topics,
		closer: func() {
			b.Unsubscribe(id)
		},
	}

	b.subs[id] = sub
	return ch, sub.closer
}

// Publish publishes an event to all matching subscribers.
// The event is assigned a monotonically increasing version number.
// Events are dropped if a subscriber's channel is full (non-blocking).
func (b *Bus) Publish(event Event) {
	event.Version = int(atomic.AddInt64(&b.ver, 1))

	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, sub := range b.subs {
		if b.matches(sub, event.Event) {
			select {
			case sub.ch <- event:
			default:
				// Drop event if subscriber is too slow to prevent blocking
				// In a real system we might want metrics here
			}
		}
	}
}

// Unsubscribe removes a subscription by its ID.
// This closes the subscription's channel and removes it from the bus.
// Safe to call multiple times for the same ID.
func (b *Bus) Unsubscribe(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if sub, ok := b.subs[id]; ok {
		close(sub.ch)
		delete(b.subs, id)
	}
}

// matches checks if a subscription matches a given topic.
// Returns true if the subscription has no topics (subscribes to all) or if the topic is in the subscription's topic list.
func (b *Bus) matches(sub *Subscription, topic EventType) bool {
	if len(sub.topics) == 0 {
		return true // Subscribe to all
	}
	for _, t := range sub.topics {
		if t == topic {
			return true
		}
	}
	return false
}
