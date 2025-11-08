package eventbus

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)



type MemoryBroker struct {
	subscribers map[string][]chan *Event
	mu          sync.RWMutex
	closed      bool
}


func NewMemoryBroker() *MemoryBroker {
	log.Warn().Msg("Using in-memory event bus - not suitable for multi-instance deployments")
	return &MemoryBroker{
		subscribers: make(map[string][]chan *Event),
	}
}


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


func (b *MemoryBroker) Subscribe(ctx context.Context, channel string) (<-chan *Event, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil, fmt.Errorf("event bus is closed")
	}

	
	eventChan := make(chan *Event, 100)

	
	b.subscribers[channel] = append(b.subscribers[channel], eventChan)

	log.Info().
		Str("channel", channel).
		Int("total_subscribers", len(b.subscribers[channel])).
		Msg("Subscribed to channel")

	return eventChan, nil
}



func (b *MemoryBroker) SubscribePattern(ctx context.Context, pattern string) (<-chan *Event, error) {
	log.Warn().
		Str("pattern", pattern).
		Msg("Memory broker doesn't support pattern matching, treating as exact channel")
	return b.Subscribe(ctx, pattern)
}


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

	
	for _, ch := range subscribers {
		close(ch)
	}

	delete(b.subscribers, channel)

	log.Info().Str("channel", channel).Msg("Unsubscribed from channel")

	return nil
}


func (b *MemoryBroker) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil
	}

	b.closed = true

	
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


func (b *MemoryBroker) HealthCheck(ctx context.Context) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return fmt.Errorf("event bus is closed")
	}

	return nil
}
