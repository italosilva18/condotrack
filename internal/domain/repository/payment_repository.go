package repository

import (
	"context"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

// PaymentRepository defines the interface for payment data access (payments table).
type PaymentRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Payment, error)
	FindByEnrollmentID(ctx context.Context, enrollmentID string) ([]entity.Payment, error)
	FindByGatewayPaymentID(ctx context.Context, gateway, gatewayPaymentID string) (*entity.Payment, error)
	FindAll(ctx context.Context, filters PaymentFilters) ([]entity.Payment, int, error)
	Create(ctx context.Context, payment *entity.Payment) error
	CreateWithTx(ctx context.Context, tx *sqlx.Tx, payment *entity.Payment) error
	Update(ctx context.Context, payment *entity.Payment) error
	UpdateWithTx(ctx context.Context, tx *sqlx.Tx, payment *entity.Payment) error
	UpdateStatus(ctx context.Context, id, status string) error
	UpdateStatusWithTx(ctx context.Context, tx *sqlx.Tx, id, status string) error
}

// PaymentFilters holds filter parameters for listing payments.
type PaymentFilters struct {
	EnrollmentID string
	Gateway      string
	Status       string
	PaymentMethod string
	DateFrom     *time.Time
	DateTo       *time.Time
	Page         int
	PerPage      int
}
