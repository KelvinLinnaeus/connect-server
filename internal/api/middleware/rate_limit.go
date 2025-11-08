package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/gin-gonic/gin"
)


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




func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*client),
		rate:    rate,
		window:  window,
	}

	
	go rl.cleanup()

	return rl
}


func RateLimitMiddleware(rate int) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, time.Minute)

	return func(c *gin.Context) {
		
		clientID := c.ClientIP()

		
		if !limiter.Allow(clientID) {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests,
				util.NewErrorResponse(http.StatusTooManyRequests, "Too many requests. Please try again later."))
			return
		}

		c.Next()
	}
}


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

	
	now := time.Now()
	elapsed := now.Sub(c.lastRefill)

	if elapsed >= rl.window {
		
		c.tokens = rl.rate
		c.lastRefill = now
	} else {
		
		tokensToAdd := int(float64(rl.rate) * elapsed.Seconds() / rl.window.Seconds())
		c.tokens += tokensToAdd
		if c.tokens > rl.rate {
			c.tokens = rl.rate
		}
		if tokensToAdd > 0 {
			c.lastRefill = now
		}
	}

	
	if c.tokens > 0 {
		c.tokens--
		return true
	}

	return false
}


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
