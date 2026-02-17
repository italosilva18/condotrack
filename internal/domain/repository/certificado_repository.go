package repository

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

// CertificadoRepository defines the interface for certificado data access
type CertificadoRepository interface {
	// FindByID returns a certificate by ID
	FindByID(ctx context.Context, id string) (*entity.Certificado, error)

	// FindByValidationCode returns a certificate by validation code
	FindByValidationCode(ctx context.Context, code string) (*entity.Certificado, error)

	// FindByStudentID returns all certificates for a student
	FindByStudentID(ctx context.Context, studentID string) ([]entity.Certificado, error)

	// FindByEnrollmentID returns a certificate by enrollment ID
	FindByEnrollmentID(ctx context.Context, enrollmentID string) (*entity.Certificado, error)

	// Create creates a new certificate
	Create(ctx context.Context, cert *entity.Certificado) error

	// Update updates an existing certificate
	Update(ctx context.Context, cert *entity.Certificado) error

	// UpdateStatus updates the certificate status
	UpdateStatus(ctx context.Context, id, status string) error

	// Delete deletes a certificate by ID
	Delete(ctx context.Context, id string) error

	// Exists checks if a certificate exists for an enrollment
	Exists(ctx context.Context, enrollmentID string) (bool, error)
}

// NotificacaoRepository defines the interface for notification data access
type NotificacaoRepository interface {
	// FindByUserID returns all notifications for a user
	FindByUserID(ctx context.Context, userID string) ([]entity.Notificacao, error)

	// FindUnreadByUserID returns unread notifications for a user
	FindUnreadByUserID(ctx context.Context, userID string) ([]entity.Notificacao, error)

	// Create creates a new notification
	Create(ctx context.Context, notif *entity.Notificacao) error

	// MarkAsRead marks a notification as read
	MarkAsRead(ctx context.Context, id string) error

	// MarkAllAsRead marks all notifications as read for a user
	MarkAllAsRead(ctx context.Context, userID string) error

	// Delete deletes a notification by ID
	Delete(ctx context.Context, id string) error

	// CountUnread returns the count of unread notifications for a user
	CountUnread(ctx context.Context, userID string) (int, error)
}

// RevenueSplitRepository defines the interface for revenue split data access
type RevenueSplitRepository interface {
	// FindByID returns a revenue split by ID
	FindByID(ctx context.Context, id string) (*entity.RevenueSplit, error)

	// FindAll returns all revenue splits with optional status filter
	FindAll(ctx context.Context, status string) ([]entity.RevenueSplit, error)

	// FindByEnrollmentID returns revenue split by enrollment ID
	FindByEnrollmentID(ctx context.Context, enrollmentID string) (*entity.RevenueSplit, error)

	// FindByPaymentID returns revenue split by payment ID
	FindByPaymentID(ctx context.Context, paymentID string) (*entity.RevenueSplit, error)

	// FindByInstructorID returns all revenue splits for an instructor
	FindByInstructorID(ctx context.Context, instructorID string) ([]entity.RevenueSplit, error)

	// Create creates a new revenue split
	Create(ctx context.Context, split *entity.RevenueSplit) error

	// CreateWithTx creates a new revenue split within a transaction
	CreateWithTx(ctx context.Context, tx *sqlx.Tx, split *entity.RevenueSplit) error

	// UpdateStatus updates the status of a revenue split
	UpdateStatus(ctx context.Context, id, status string) error

	// GetTotalByInstructor returns total earnings for an instructor
	GetTotalByInstructor(ctx context.Context, instructorID string) (float64, error)
}
