package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger returns a middleware that logs HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if query != "" {
			path = path + "?" + query
		}

		log.Printf("[%d] %s %s %s %v",
			status,
			method,
			path,
			clientIP,
			latency,
		)
	}
}
