package handler

import (
	"strconv"

	"github.com/condotrack/api/internal/delivery/http/middleware"
	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getUserIDFromContext gets user_id from JWT context (required auth)
func getUserIDFromContext(c *gin.Context) string {
	if userID, ok := middleware.GetUserID(c); ok && userID != "" {
		return userID
	}
	return ""
}

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	repo repository.NotificacaoRepository
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(repo repository.NotificacaoRepository) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

// ListNotifications handles GET /api/v1/notifications
// Extracts user_id from JWT token or query parameter. Supports is_read filter.
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	ctx := c.Request.Context()

	userID := getUserIDFromContext(c)
	if userID == "" {
		response.BadRequest(c, "user_id is required (provide via Authorization header or query parameter)")
		return
	}

	// Check if filtering by read status
	isReadParam := c.Query("is_read")
	if isReadParam != "" {
		isRead, err := strconv.ParseBool(isReadParam)
		if err != nil {
			response.BadRequest(c, "Invalid is_read parameter: must be 'true' or 'false'")
			return
		}

		if !isRead {
			// Return only unread notifications
			notifications, err := h.repo.FindUnreadByUserID(ctx, userID)
			if err != nil {
				response.SafeInternalError(c, "Failed to fetch notifications", err)
				return
			}
			response.Success(c, notifications)
			return
		}
	}

	// Return all notifications for the user
	notifications, err := h.repo.FindByUserID(ctx, userID)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch notifications", err)
		return
	}

	response.Success(c, notifications)
}

// GetUnreadNotifications handles GET /api/v1/notifications/unread
func (h *NotificationHandler) GetUnreadNotifications(c *gin.Context) {
	ctx := c.Request.Context()

	userID := getUserIDFromContext(c)
	if userID == "" {
		response.BadRequest(c, "user_id is required (provide via Authorization header or query parameter)")
		return
	}

	notifications, err := h.repo.FindUnreadByUserID(ctx, userID)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch unread notifications", err)
		return
	}

	response.Success(c, notifications)
}

// GetUnreadCount handles GET /api/v1/notifications/count
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	ctx := c.Request.Context()

	userID := getUserIDFromContext(c)
	if userID == "" {
		response.BadRequest(c, "user_id is required (provide via Authorization header or query parameter)")
		return
	}

	count, err := h.repo.CountUnread(ctx, userID)
	if err != nil {
		response.SafeInternalError(c, "Failed to count unread notifications", err)
		return
	}

	response.Success(c, gin.H{
		"user_id":      userID,
		"unread_count": count,
	})
}

// CreateNotification handles POST /api/v1/notifications
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Validate notification type
	validTypes := map[string]bool{
		entity.NotificationTypePayment:     true,
		entity.NotificationTypeEnrollment:  true,
		entity.NotificationTypeAudit:       true,
		entity.NotificationTypeCertificate: true,
		entity.NotificationTypeSystem:      true,
	}

	if !validTypes[req.Type] {
		response.BadRequest(c, "Invalid notification type. Valid types: payment, enrollment, audit, certificate, system")
		return
	}

	// Create notification entity
	notification := &entity.Notificacao{
		ID:      uuid.New().String(),
		UserID:  req.UserID,
		Type:    req.Type,
		Title:   req.Title,
		Message: req.Message,
		Data:    req.Data,
		Read:    false,
	}

	if err := h.repo.Create(ctx, notification); err != nil {
		response.SafeInternalError(c, "Failed to create notification", err)
		return
	}

	response.Created(c, notification)
}

// MarkAsRead handles PATCH /api/v1/notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		response.BadRequest(c, "Notification ID is required")
		return
	}

	if err := h.repo.MarkAsRead(ctx, id); err != nil {
		response.SafeInternalError(c, "Failed to mark notification as read", err)
		return
	}

	response.Success(c, gin.H{
		"message": "Notification marked as read",
		"id":      id,
	})
}

// MarkAllAsRead handles PATCH /api/v1/notifications/mark-all-read
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	ctx := c.Request.Context()

	userID := getUserIDFromContext(c)
	if userID == "" {
		response.BadRequest(c, "user_id is required (provide via Authorization header or query parameter)")
		return
	}

	if err := h.repo.MarkAllAsRead(ctx, userID); err != nil {
		response.SafeInternalError(c, "Failed to mark all notifications as read", err)
		return
	}

	response.Success(c, gin.H{
		"message": "All notifications marked as read",
		"user_id": userID,
	})
}

// DeleteNotification handles DELETE /api/v1/notifications/:id
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		response.BadRequest(c, "Notification ID is required")
		return
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		response.SafeInternalError(c, "Failed to delete notification", err)
		return
	}

	response.Success(c, gin.H{
		"message": "Notification deleted successfully",
		"id":      id,
	})
}
