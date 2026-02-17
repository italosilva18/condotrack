package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// visitor holds rate limiting state for a single IP address
type visitor struct {
	tokens    int
	lastSeen  time.Time
	mu        sync.Mutex
}

// MaxBodySize returns a middleware that limits the request body size.
// Requests exceeding maxBytes will get a 413 error.
func MaxBodySize(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

// RateLimiter returns a middleware that limits requests per IP address
// using an in-memory token bucket algorithm.
// limit is the maximum number of requests allowed within the given window.
func RateLimiter(limit int, window time.Duration) gin.HandlerFunc {
	var visitors sync.Map

	// Background goroutine to clean up expired entries
	go func() {
		for {
			time.Sleep(window * 2)
			visitors.Range(func(key, value interface{}) bool {
				v := value.(*visitor)
				v.mu.Lock()
				if time.Since(v.lastSeen) > window*2 {
					v.mu.Unlock()
					visitors.Delete(key)
					return true
				}
				v.mu.Unlock()
				return true
			})
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		val, _ := visitors.LoadOrStore(ip, &visitor{
			tokens:   limit,
			lastSeen: time.Now(),
		})
		v := val.(*visitor)

		v.mu.Lock()
		defer v.mu.Unlock()

		// Replenish tokens based on elapsed time
		elapsed := time.Since(v.lastSeen)
		if elapsed > window {
			// Full window has passed, reset tokens
			v.tokens = limit
		} else {
			// Proportional replenishment
			replenish := int(float64(limit) * (float64(elapsed) / float64(window)))
			v.tokens += replenish
			if v.tokens > limit {
				v.tokens = limit
			}
		}
		v.lastSeen = time.Now()

		// Check if request is allowed
		if v.tokens <= 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "Too many requests",
			})
			return
		}

		// Consume a token
		v.tokens--

		c.Next()
	}
}
