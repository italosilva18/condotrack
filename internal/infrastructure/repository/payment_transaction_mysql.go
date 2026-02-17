package repository

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type paymentTransactionMySQLRepository struct {
	db *sqlx.DB
}

// NewPaymentTransactionMySQLRepository creates a new MySQL implementation of PaymentTransactionRepository.
func NewPaymentTransactionMySQLRepository(db *sqlx.DB) repository.PaymentTransactionRepository {
	return &paymentTransactionMySQLRepository{db: db}
}

func (r *paymentTransactionMySQLRepository) FindByPaymentID(ctx context.Context, paymentID string) ([]entity.PaymentTransaction, error) {
	var txns []entity.PaymentTransaction
	query := `SELECT id, payment_id, previous_status, new_status,
		event_source, event_type, gateway_event_id,
		amount, description, raw_payload,
		ip_address, user_agent, triggered_by, created_at
		FROM payment_transactions
		WHERE payment_id = ?
		ORDER BY created_at ASC`
	err := r.db.SelectContext(ctx, &txns, query, paymentID)
	return txns, err
}

func (r *paymentTransactionMySQLRepository) Create(ctx context.Context, txLog *entity.PaymentTransaction) error {
	query := `INSERT INTO payment_transactions (
		id, payment_id, previous_status, new_status,
		event_source, event_type, gateway_event_id,
		amount, description, raw_payload,
		ip_address, user_agent, triggered_by, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query,
		txLog.ID, txLog.PaymentID, txLog.PreviousStatus, txLog.NewStatus,
		txLog.EventSource, txLog.EventType, txLog.GatewayEventID,
		txLog.Amount, txLog.Description, txLog.RawPayload,
		txLog.IPAddress, txLog.UserAgent, txLog.TriggeredBy,
	)
	return err
}

func (r *paymentTransactionMySQLRepository) CreateWithTx(ctx context.Context, tx *sqlx.Tx, txLog *entity.PaymentTransaction) error {
	query := `INSERT INTO payment_transactions (
		id, payment_id, previous_status, new_status,
		event_source, event_type, gateway_event_id,
		amount, description, raw_payload,
		ip_address, user_agent, triggered_by, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := tx.ExecContext(ctx, query,
		txLog.ID, txLog.PaymentID, txLog.PreviousStatus, txLog.NewStatus,
		txLog.EventSource, txLog.EventType, txLog.GatewayEventID,
		txLog.Amount, txLog.Description, txLog.RawPayload,
		txLog.IPAddress, txLog.UserAgent, txLog.TriggeredBy,
	)
	return err
}
