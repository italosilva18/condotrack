package handler

import (
	"context"
	"time"

	"github.com/condotrack/api/internal/infrastructure/database"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *database.MySQL
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.MySQL) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheck handles GET /api/v1/health
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Check database connection
	dbStatus := "healthy"
	if err := h.db.Health(ctx); err != nil {
		dbStatus = "unhealthy: " + err.Error()
	}

	c.JSON(200, gin.H{
		"success":   true,
		"message":   "OK",
		"timestamp": time.Now().Format(time.RFC3339),
		"services": gin.H{
			"database": dbStatus,
		},
	})
}

// Ping handles GET /ping
func (h *HealthHandler) Ping(c *gin.Context) {
	response.Success(c, gin.H{
		"message": "pong",
	})
}
