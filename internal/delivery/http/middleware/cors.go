package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS returns a middleware that handles Cross-Origin Resource Sharing.
// If allowedOrigins is non-empty, only those origins are allowed (with credentials).
// If empty, all origins are allowed without credentials (no wildcard + credentials).
func CORS(allowedOrigins string) gin.HandlerFunc {
	allowed := make(map[string]bool)
	if allowedOrigins != "" {
		for _, o := range strings.Split(allowedOrigins, ",") {
			origin := strings.TrimSpace(o)
			if origin != "" {
				allowed[origin] = true
			}
		}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if len(allowed) > 0 {
			// Whitelist mode: only allow configured origins
			if allowed[origin] {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			// If origin not in whitelist, don't set CORS headers (browser will block)
		} else {
			// Development mode: allow all origins but without credentials for safety
			if origin != "" {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
			} else {
				c.Header("Access-Control-Allow-Origin", "*")
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, X-Requested-With")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
