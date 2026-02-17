package middleware

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestID returns a middleware that assigns a unique request ID to each request.
// It checks for an existing X-Request-ID header first; if not present, it generates one.
// The request ID is set in the response header and stored in the Gin context.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for existing request ID header
		id := c.GetHeader("X-Request-ID")

		// Generate a new ID if none was provided
		if id == "" {
			id = fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int63())
		}

		// Set the request ID in the response header
		c.Writer.Header().Set("X-Request-ID", id)

		// Store in context for downstream handlers
		c.Set("request_id", id)

		c.Next()
	}
}
