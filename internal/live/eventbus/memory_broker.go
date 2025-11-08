package eventbus

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

// MemoryBroker implements EventBus using in-memory channels
// This is useful for development, testing, and single-instance deployments
type MemoryBroker struct {
	subscribers map[string][]chan *Event
	mu          sync.RWMutex
	closed      bool
}

// NewMemoryBroker creates a new in-memory event bus
func NewMemoryBroker() *MemoryBroker {
	log.Warn().Msg("Using in-memory event bus - not suitable for multi-instance deployments")
	return &MemoryBroker{
		subscribers: make(map[string][]chan *Event),
	}
}

// Publish publishes an event to all subscribers of the channel
func (b *MemoryBroker) Publish(ctx context.Context, event *Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return fmt.Errorf("event bus is closed")
	}

	subscribers, exists := b.subscribers[event.Channel]
	if !exists || len(subscribers) == 0 {
		log.Debug().
			Str("event_id", event.ID).
			Str("channel", event.Channel).
			Msg("No subscribers for channel")
		return nil
	}

	// Send to all subscribers (non-blocking)
	for _, ch := range subscribers {
		select {
		case ch <- event:
			log.Debug().
				Str("event_id", event.ID).
				Str("event_type", event.Type).
				Str("channel", event.Channel).
				Msg("Event published to subscriber")
		default:
			log.Warn().
				Str("event_id", event.ID).
				Str("channel", event.Channel).
				Msg("Subscriber channel full, dropping event")
		}
	}

	return nil
}

// Subscribe subscribes to a channel
func (b *MemoryBroker) Subscribe(ctx context.Context, channel string) (<-chan *Event, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil, fmt.Errorf("event bus is closed")
	}

	// Create buffered channel
	eventChan := make(chan *Event, 100)

	// Add to subscribers
	b.subscribers[channel] = append(b.subscribers[channel], eventChan)

	log.Info().
		Str("channel", channel).
		Int("total_subscribers", len(b.subscribers[channel])).
		Msg("Subscribed to channel")

	return eventChan, nil
}

// SubscribePattern subscribes to channels matching a pattern
// Note: In-memory broker doesn't support pattern matching, so this subscribes to exact channel
func (b *MemoryBroker) SubscribePattern(ctx context.Context, pattern string) (<-chan *Event, error) {
	log.Warn().
		Str("pattern", pattern).
		Msg("Memory broker doesn't support pattern matching, treating as exact channel")
	return b.Subscribe(ctx, pattern)
}

// Unsubscribe unsubscribes from a channel
func (b *MemoryBroker) Unsubscribe(ctx context.Context, channel string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return fmt.Errorf("event bus is closed")
	}

	subscribers, exists := b.subscribers[channel]
	if !exists {
		return fmt.Errorf("not subscribed to channel: %s", channel)
	}

	// Close all subscriber channels
	for _, ch := range subscribers {
		close(ch)
	}

	delete(b.subscribers, channel)

	log.Info().Str("channel", channel).Msg("Unsubscribed from channel")

	return nil
}

// Close closes the event bus
func (b *MemoryBroker) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil
	}

	b.closed = true

	// Close all subscriber channels
	for channel, subscribers := range b.subscribers {
		for _, ch := range subscribers {
			close(ch)
		}
		log.Debug().Str("channel", channel).Msg("Closed subscriber channels")
	}

	b.subscribers = make(map[string][]chan *Event)

	log.Info().Msg("Memory event bus closed")

	return nil
}

// HealthCheck always returns healthy for memory broker
func (b *MemoryBroker) HealthCheck(ctx context.Context) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return fmt.Errorf("event bus is closed")
	}

	return nil
}
