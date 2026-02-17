package middleware

import (
	"strings"

	"github.com/condotrack/api/internal/infrastructure/auth"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

const (
	// AuthorizationHeader is the header key for authorization
	AuthorizationHeader = "Authorization"
	// BearerPrefix is the prefix for bearer tokens
	BearerPrefix = "Bearer "
	// UserIDKey is the context key for user ID
	UserIDKey = "user_id"
	// UserEmailKey is the context key for user email
	UserEmailKey = "user_email"
	// UserRoleKey is the context key for user role
	UserRoleKey = "user_role"
	// ClaimsKey is the context key for JWT claims
	ClaimsKey = "claims"
)

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			response.Unauthorized(c, "Invalid authorization header format. Use: Bearer <token>")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == "" {
			response.Unauthorized(c, "Token is required")
			c.Abort()
			return
		}

		// Check if token has been blacklisted (logout)
		if jwtManager.IsBlacklisted(tokenString) {
			response.Unauthorized(c, "Token has been revoked")
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			switch err {
			case auth.ErrExpiredToken:
				response.Unauthorized(c, "Token has expired")
			case auth.ErrInvalidToken:
				response.Unauthorized(c, "Invalid token")
			case auth.ErrInvalidClaims:
				response.Unauthorized(c, "Invalid token claims")
			default:
				response.Unauthorized(c, "Authentication failed")
			}
			c.Abort()
			return
		}

		// Set user information in context
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)
		c.Set(UserRoleKey, claims.Role)
		c.Set(ClaimsKey, claims)

		c.Next()
	}
}

// RequireRole creates a middleware that checks if the user has one of the required roles
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(UserRoleKey)
		if !exists {
			response.Unauthorized(c, "User role not found in context")
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			response.InternalError(c, "Invalid user role type")
			c.Abort()
			return
		}

		// Check if user role is in the allowed roles
		authorized := false
		for _, allowedRole := range roles {
			if role == allowedRole {
				authorized = true
				break
			}
		}

		if !authorized {
			response.Forbidden(c, "You don't have permission to access this resource")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin creates a middleware that allows only admin users
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

// RequireAdminOrManager creates a middleware that allows admin or manager users
func RequireAdminOrManager() gin.HandlerFunc {
	return RequireRole("admin", "manager")
}

// OptionalAuth creates a middleware that validates JWT tokens if present, but doesn't require them
func OptionalAuth(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == "" {
			c.Next()
			return
		}

		if jwtManager.IsBlacklisted(tokenString) {
			c.Next()
			return
		}

		claims, err := jwtManager.ValidateToken(tokenString)
		if err == nil {
			// Set user information in context if token is valid
			c.Set(UserIDKey, claims.UserID)
			c.Set(UserEmailKey, claims.Email)
			c.Set(UserRoleKey, claims.Role)
			c.Set(ClaimsKey, claims)
		}

		c.Next()
	}
}

// GetUserID retrieves the user ID from context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}

// GetUserEmail retrieves the user email from context
func GetUserEmail(c *gin.Context) (string, bool) {
	userEmail, exists := c.Get(UserEmailKey)
	if !exists {
		return "", false
	}
	email, ok := userEmail.(string)
	return email, ok
}

// GetUserRole retrieves the user role from context
func GetUserRole(c *gin.Context) (string, bool) {
	userRole, exists := c.Get(UserRoleKey)
	if !exists {
		return "", false
	}
	role, ok := userRole.(string)
	return role, ok
}

// GetClaims retrieves the JWT claims from context
func GetClaims(c *gin.Context) (*auth.Claims, bool) {
	claims, exists := c.Get(ClaimsKey)
	if !exists {
		return nil, false
	}
	c2, ok := claims.(*auth.Claims)
	return c2, ok
}
