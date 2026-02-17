package entity

import "time"

// Notificacao represents a notification entity
type Notificacao struct {
	ID          string     `db:"id" json:"id"`
	UserID      string     `db:"user_id" json:"user_id"`
	Type        string     `db:"type" json:"type"`
	Title       string     `db:"title" json:"title"`
	Message     string     `db:"message" json:"message"`
	Data        *string    `db:"data" json:"data,omitempty"`
	Read        bool       `db:"is_read" json:"is_read"`
	ReadAt      *time.Time `db:"read_at" json:"read_at,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
}

// Notification type constants
const (
	NotificationTypePayment    = "payment"
	NotificationTypeEnrollment = "enrollment"
	NotificationTypeAudit      = "audit"
	NotificationTypeCertificate = "certificate"
	NotificationTypeSystem     = "system"
)

// CreateNotificationRequest represents the request to create a notification
type CreateNotificationRequest struct {
	UserID  string  `json:"user_id" binding:"required"`
	Type    string  `json:"type" binding:"required"`
	Title   string  `json:"title" binding:"required"`
	Message string  `json:"message" binding:"required"`
	Data    *string `json:"data,omitempty"`
}
