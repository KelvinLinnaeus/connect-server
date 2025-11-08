package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	clients map[string]*client
	mu      sync.RWMutex
	rate    int
	window  time.Duration
}

type client struct {
	tokens     int
	lastRefill time.Time
	mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter
// rate: maximum number of requests allowed per window
// window: time window for rate limiting (typically 1 minute)
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*client),
		rate:    rate,
		window:  window,
	}

	// Cleanup goroutine to remove old clients
	go rl.cleanup()

	return rl
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(rate int) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, time.Minute)

	return func(c *gin.Context) {
		// Use client IP as identifier
		clientID := c.ClientIP()

		// Check if request is allowed
		if !limiter.Allow(clientID) {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests,
				util.NewErrorResponse("rate_limit_exceeded", "Too many requests. Please try again later."))
			return
		}

		c.Next()
	}
}

// Allow checks if a request from the given client is allowed
func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.Lock()
	c, exists := rl.clients[clientID]
	if !exists {
		c = &client{
			tokens:     rl.rate,
			lastRefill: time.Now(),
		}
		rl.clients[clientID] = c
	}
	rl.mu.Unlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Refill tokens based on time passed
	now := time.Now()
	elapsed := now.Sub(c.lastRefill)

	if elapsed >= rl.window {
		// Full window passed, refill to max
		c.tokens = rl.rate
		c.lastRefill = now
	} else {
		// Partial refill based on time elapsed
		tokensToAdd := int(float64(rl.rate) * elapsed.Seconds() / rl.window.Seconds())
		c.tokens += tokensToAdd
		if c.tokens > rl.rate {
			c.tokens = rl.rate
		}
		if tokensToAdd > 0 {
			c.lastRefill = now
		}
	}

	// Check if we have tokens available
	if c.tokens > 0 {
		c.tokens--
		return true
	}

	return false
}

// cleanup removes clients that haven't made requests in a while
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for id, c := range rl.clients {
			c.mu.Lock()
			if now.Sub(c.lastRefill) > 10*time.Minute {
				delete(rl.clients, id)
			}
			c.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}
