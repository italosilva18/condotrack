package repository

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

// PaymentTransactionRepository defines the interface for payment transaction log access.
type PaymentTransactionRepository interface {
	FindByPaymentID(ctx context.Context, paymentID string) ([]entity.PaymentTransaction, error)
	Create(ctx context.Context, txLog *entity.PaymentTransaction) error
	CreateWithTx(ctx context.Context, tx *sqlx.Tx, txLog *entity.PaymentTransaction) error
}
