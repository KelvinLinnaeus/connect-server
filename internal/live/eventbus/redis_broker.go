package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)


type RedisBroker struct {
	client        *redis.Client
	subscriptions map[string]*redis.PubSub
	mu            sync.RWMutex
	closed        bool
	redisURL      string
	ctx           context.Context
	cancel        context.CancelFunc

	
	reconnecting  atomic.Bool
	reconnectChan chan struct{}

	
	metrics       *BrokerMetrics
}


type BrokerMetrics struct {
	EventsPublished   atomic.Int64
	EventsReceived    atomic.Int64
	PublishErrors     atomic.Int64
	ReconnectCount    atomic.Int64
	LastReconnectTime time.Time
	mu                sync.RWMutex
}


func (m *BrokerMetrics) GetEventsPublished() int64 {
	return m.EventsPublished.Load()
}


func (m *BrokerMetrics) GetEventsReceived() int64 {
	return m.EventsReceived.Load()
}


func (m *BrokerMetrics) GetPublishErrors() int64 {
	return m.PublishErrors.Load()
}


func (m *BrokerMetrics) GetReconnectCount() int64 {
	return m.ReconnectCount.Load()
}


func NewRedisBroker(redisURL string) (*RedisBroker, error) {
	ctx, cancel := context.WithCancel(context.Background())

	broker := &RedisBroker{
		subscriptions: make(map[string]*redis.PubSub),
		redisURL:      redisURL,
		ctx:           ctx,
		cancel:        cancel,
		reconnectChan: make(chan struct{}, 1),
		metrics:       &BrokerMetrics{},
	}

	if err := broker.connect(); err != nil {
		cancel()
		return nil, err
	}

	
	go broker.monitorConnection()

	log.Info().Str("redis_url", redisURL).Msg("Redis live broker active (multi-instance ready)")

	return broker, nil
}


func (b *RedisBroker) connect() error {
	opts, err := redis.ParseURL(b.redisURL)
	if err != nil {
		return fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	
	opts.PoolSize = 10
	opts.MinIdleConns = 5
	opts.MaxRetries = 3
	opts.DialTimeout = 5 * time.Second
	opts.ReadTimeout = 3 * time.Second
	opts.WriteTimeout = 3 * time.Second

	client := redis.NewClient(opts)

	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	b.client = client
	return nil
}


func (b *RedisBroker) monitorConnection() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-b.ctx.Done():
			return
		case <-ticker.C:
			if err := b.HealthCheck(b.ctx); err != nil {
				log.Warn().Err(err).Msg("Redis health check failed, initiating reconnect")
				select {
				case b.reconnectChan <- struct{}{}:
				default:
				}
			}
		case <-b.reconnectChan:
			b.reconnect()
		}
	}
}


func (b *RedisBroker) reconnect() {
	if !b.reconnecting.CompareAndSwap(false, true) {
		
		return
	}
	defer b.reconnecting.Store(false)

	log.Warn().Msg("Redis connection lost, attempting to reconnect")

	
	backoff := time.Second
	maxBackoff := 30 * time.Second
	attempt := 1

	for {
		select {
		case <-b.ctx.Done():
			return
		default:
		}

		log.Info().
			Int("attempt", attempt).
			Dur("backoff", backoff).
			Msg("Attempting Redis reconnect")

		if err := b.connect(); err != nil {
			log.Error().Err(err).
				Int("attempt", attempt).
				Msg("Redis reconnect failed")

			
			time.Sleep(backoff)
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			attempt++
			continue
		}

		
		b.metrics.ReconnectCount.Add(1)
		b.metrics.mu.Lock()
		b.metrics.LastReconnectTime = time.Now()
		b.metrics.mu.Unlock()

		log.Info().
			Int("attempt", attempt).
			Int64("total_reconnects", b.metrics.ReconnectCount.Load()).
			Msg("Redis reconnected successfully")

		
		b.resubscribeAll()
		return
	}
}


func (b *RedisBroker) resubscribeAll() {
	b.mu.RLock()
	channels := make([]string, 0, len(b.subscriptions))
	for channel := range b.subscriptions {
		channels = append(channels, channel)
	}
	b.mu.RUnlock()

	for _, channel := range channels {
		log.Info().Str("channel", channel).Msg("Resubscribing to channel")
		
	}
}


func (b *RedisBroker) Publish(ctx context.Context, event *Event) error {
	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return fmt.Errorf("event bus is closed")
	}
	client := b.client
	b.mu.RUnlock()

	
	data, err := json.Marshal(event)
	if err != nil {
		b.metrics.PublishErrors.Add(1)
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	
	maxRetries := 3
	var publishErr error
	for i := 0; i < maxRetries; i++ {
		publishErr = client.Publish(ctx, event.Channel, data).Err()
		if publishErr == nil {
			break
		}

		
		if i < maxRetries-1 {
			log.Warn().Err(publishErr).
				Int("retry", i+1).
				Msg("Publish failed, retrying")
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
		}
	}

	if publishErr != nil {
		b.metrics.PublishErrors.Add(1)
		
		select {
		case b.reconnectChan <- struct{}{}:
		default:
		}
		return fmt.Errorf("failed to publish event: %w", publishErr)
	}

	b.metrics.EventsPublished.Add(1)

	log.Debug().
		Str("event_id", event.ID).
		Str("event_type", event.Type).
		Str("channel", event.Channel).
		Msg("Event published")

	return nil
}


func (b *RedisBroker) Subscribe(ctx context.Context, channel string) (<-chan *Event, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil, fmt.Errorf("event bus is closed")
	}

	
	if _, exists := b.subscriptions[channel]; exists {
		return nil, fmt.Errorf("already subscribed to channel: %s", channel)
	}

	
	pubsub := b.client.Subscribe(ctx, channel)
	b.subscriptions[channel] = pubsub

	
	eventChan := make(chan *Event, 100) 

	
	go b.receiveMessages(ctx, pubsub, eventChan, channel)

	log.Info().Str("channel", channel).Msg("Subscribed to channel")

	return eventChan, nil
}


func (b *RedisBroker) SubscribePattern(ctx context.Context, pattern string) (<-chan *Event, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil, fmt.Errorf("event bus is closed")
	}

	
	pubsub := b.client.PSubscribe(ctx, pattern)
	b.subscriptions[pattern] = pubsub

	
	eventChan := make(chan *Event, 100)

	
	go b.receiveMessages(ctx, pubsub, eventChan, pattern)

	log.Info().Str("pattern", pattern).Msg("Subscribed to pattern")

	return eventChan, nil
}


func (b *RedisBroker) receiveMessages(ctx context.Context, pubsub *redis.PubSub, eventChan chan<- *Event, identifier string) {
	defer close(eventChan)

	ch := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			log.Debug().Str("identifier", identifier).Msg("Context cancelled, stopping message receiver")
			return
		case <-b.ctx.Done():
			log.Debug().Str("identifier", identifier).Msg("Broker context cancelled, stopping message receiver")
			return
		case msg, ok := <-ch:
			if !ok {
				log.Debug().Str("identifier", identifier).Msg("Redis channel closed")
				
				select {
				case b.reconnectChan <- struct{}{}:
				default:
				}
				return
			}

			
			var event Event
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				log.Error().Err(err).Str("payload", msg.Payload).Msg("Failed to unmarshal event")
				continue
			}

			b.metrics.EventsReceived.Add(1)

			
			select {
			case eventChan <- &event:
				log.Debug().
					Str("event_id", event.ID).
					Str("event_type", event.Type).
					Str("channel", event.Channel).
					Msg("Event received")
			default:
				log.Warn().
					Str("event_id", event.ID).
					Str("channel", event.Channel).
					Msg("Event channel full, dropping event")
			}
		}
	}
}


func (b *RedisBroker) Unsubscribe(ctx context.Context, channel string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return fmt.Errorf("event bus is closed")
	}

	pubsub, exists := b.subscriptions[channel]
	if !exists {
		return fmt.Errorf("not subscribed to channel: %s", channel)
	}

	
	if err := pubsub.Unsubscribe(ctx, channel); err != nil {
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}

	
	if err := pubsub.Close(); err != nil {
		log.Warn().Err(err).Str("channel", channel).Msg("Error closing pubsub")
	}

	delete(b.subscriptions, channel)

	log.Info().Str("channel", channel).Msg("Unsubscribed from channel")

	return nil
}


func (b *RedisBroker) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil
	}

	b.closed = true

	
	b.cancel()

	
	for channel, pubsub := range b.subscriptions {
		if err := pubsub.Close(); err != nil {
			log.Warn().Err(err).Str("channel", channel).Msg("Error closing subscription")
		}
	}

	
	if err := b.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis client: %w", err)
	}

	log.Info().Msg("Event bus closed")

	return nil
}


func (b *RedisBroker) HealthCheck(ctx context.Context) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return fmt.Errorf("event bus is closed")
	}

	if err := b.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}

	return nil
}


func (b *RedisBroker) GetMetrics() BrokerMetrics {
	b.metrics.mu.RLock()
	defer b.metrics.mu.RUnlock()

	metrics := BrokerMetrics{
		LastReconnectTime: b.metrics.LastReconnectTime,
	}
	
	metrics.EventsPublished.Store(b.metrics.EventsPublished.Load())
	metrics.EventsReceived.Store(b.metrics.EventsReceived.Load())
	metrics.PublishErrors.Store(b.metrics.PublishErrors.Load())
	metrics.ReconnectCount.Store(b.metrics.ReconnectCount.Load())

	return metrics
}
