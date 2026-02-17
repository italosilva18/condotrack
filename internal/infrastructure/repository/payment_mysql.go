package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type paymentMySQLRepository struct {
	db *sqlx.DB
}

// NewPaymentMySQLRepository creates a new MySQL implementation of PaymentRepository.
func NewPaymentMySQLRepository(db *sqlx.DB) repository.PaymentRepository {
	return &paymentMySQLRepository{db: db}
}

const paymentColumns = `id, enrollment_id, payer_user_id, payer_name, payer_email, payer_cpf,
	gross_amount, discount_amount, net_amount, gateway_fee, refunded_amount,
	payment_method, gateway, gateway_payment_id, gateway_customer_id,
	gateway_invoice_url, gateway_metadata,
	installment_count, installment_of, installment_number,
	status, coupon_id, due_date, paid_at, refunded_at, cancelled_at, expires_at,
	created_at, updated_at`

func (r *paymentMySQLRepository) FindByID(ctx context.Context, id string) (*entity.Payment, error) {
	var p entity.Payment
	query := fmt.Sprintf(`SELECT %s FROM payments WHERE id = ?`, paymentColumns)
	err := r.db.GetContext(ctx, &p, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *paymentMySQLRepository) FindByEnrollmentID(ctx context.Context, enrollmentID string) ([]entity.Payment, error) {
	var payments []entity.Payment
	query := fmt.Sprintf(`SELECT %s FROM payments WHERE enrollment_id = ? ORDER BY created_at DESC`, paymentColumns)
	err := r.db.SelectContext(ctx, &payments, query, enrollmentID)
	return payments, err
}

func (r *paymentMySQLRepository) FindByGatewayPaymentID(ctx context.Context, gw, gatewayPaymentID string) (*entity.Payment, error) {
	var p entity.Payment
	query := fmt.Sprintf(`SELECT %s FROM payments WHERE gateway = ? AND gateway_payment_id = ?`, paymentColumns)
	err := r.db.GetContext(ctx, &p, query, gw, gatewayPaymentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *paymentMySQLRepository) FindAll(ctx context.Context, filters repository.PaymentFilters) ([]entity.Payment, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if filters.EnrollmentID != "" {
		where = append(where, "enrollment_id = ?")
		args = append(args, filters.EnrollmentID)
	}
	if filters.Gateway != "" {
		where = append(where, "gateway = ?")
		args = append(args, filters.Gateway)
	}
	if filters.Status != "" {
		where = append(where, "status = ?")
		args = append(args, filters.Status)
	}
	if filters.PaymentMethod != "" {
		where = append(where, "payment_method = ?")
		args = append(args, filters.PaymentMethod)
	}
	if filters.DateFrom != nil {
		where = append(where, "created_at >= ?")
		args = append(args, *filters.DateFrom)
	}
	if filters.DateTo != nil {
		where = append(where, "created_at <= ?")
		args = append(args, *filters.DateTo)
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM payments WHERE %s`, whereClause)
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	page := filters.Page
	if page < 1 {
		page = 1
	}
	perPage := filters.PerPage
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	query := fmt.Sprintf(`SELECT %s FROM payments WHERE %s ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		paymentColumns, whereClause)
	args = append(args, perPage, offset)

	var payments []entity.Payment
	if err := r.db.SelectContext(ctx, &payments, query, args...); err != nil {
		return nil, 0, err
	}
	return payments, total, nil
}

func (r *paymentMySQLRepository) Create(ctx context.Context, p *entity.Payment) error {
	query := `INSERT INTO payments (
		id, enrollment_id, payer_user_id, payer_name, payer_email, payer_cpf,
		gross_amount, discount_amount, net_amount, gateway_fee, refunded_amount,
		payment_method, gateway, gateway_payment_id, gateway_customer_id,
		gateway_invoice_url, gateway_metadata,
		installment_count, installment_of, installment_number,
		status, coupon_id, due_date, paid_at, refunded_at, cancelled_at, expires_at,
		created_at
	) VALUES (
		?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?,
		?, ?, ?, ?,
		?, ?,
		?, ?, ?,
		?, ?, ?, ?, ?, ?, ?,
		NOW()
	)`
	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.EnrollmentID, p.PayerUserID, p.PayerName, p.PayerEmail, p.PayerCPF,
		p.GrossAmount, p.DiscountAmount, p.NetAmount, p.GatewayFee, p.RefundedAmount,
		p.PaymentMethod, p.Gateway, p.GatewayPaymentID, p.GatewayCustomerID,
		p.GatewayInvoiceURL, p.GatewayMetadata,
		p.InstallmentCount, p.InstallmentOf, p.InstallmentNumber,
		p.Status, p.CouponID, p.DueDate, p.PaidAt, p.RefundedAt, p.CancelledAt, p.ExpiresAt,
	)
	return err
}

func (r *paymentMySQLRepository) CreateWithTx(ctx context.Context, tx *sqlx.Tx, p *entity.Payment) error {
	query := `INSERT INTO payments (
		id, enrollment_id, payer_user_id, payer_name, payer_email, payer_cpf,
		gross_amount, discount_amount, net_amount, gateway_fee, refunded_amount,
		payment_method, gateway, gateway_payment_id, gateway_customer_id,
		gateway_invoice_url, gateway_metadata,
		installment_count, installment_of, installment_number,
		status, coupon_id, due_date, paid_at, refunded_at, cancelled_at, expires_at,
		created_at
	) VALUES (
		?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?,
		?, ?, ?, ?,
		?, ?,
		?, ?, ?,
		?, ?, ?, ?, ?, ?, ?,
		NOW()
	)`
	_, err := tx.ExecContext(ctx, query,
		p.ID, p.EnrollmentID, p.PayerUserID, p.PayerName, p.PayerEmail, p.PayerCPF,
		p.GrossAmount, p.DiscountAmount, p.NetAmount, p.GatewayFee, p.RefundedAmount,
		p.PaymentMethod, p.Gateway, p.GatewayPaymentID, p.GatewayCustomerID,
		p.GatewayInvoiceURL, p.GatewayMetadata,
		p.InstallmentCount, p.InstallmentOf, p.InstallmentNumber,
		p.Status, p.CouponID, p.DueDate, p.PaidAt, p.RefundedAt, p.CancelledAt, p.ExpiresAt,
	)
	return err
}

func (r *paymentMySQLRepository) Update(ctx context.Context, p *entity.Payment) error {
	query := `UPDATE payments SET
		payer_name = ?, payer_email = ?, payer_cpf = ?,
		gross_amount = ?, discount_amount = ?, net_amount = ?, gateway_fee = ?, refunded_amount = ?,
		gateway_payment_id = ?, gateway_customer_id = ?,
		gateway_invoice_url = ?, gateway_metadata = ?,
		status = ?, paid_at = ?, refunded_at = ?, cancelled_at = ?,
		updated_at = NOW()
		WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		p.PayerName, p.PayerEmail, p.PayerCPF,
		p.GrossAmount, p.DiscountAmount, p.NetAmount, p.GatewayFee, p.RefundedAmount,
		p.GatewayPaymentID, p.GatewayCustomerID,
		p.GatewayInvoiceURL, p.GatewayMetadata,
		p.Status, p.PaidAt, p.RefundedAt, p.CancelledAt,
		p.ID,
	)
	return err
}

func (r *paymentMySQLRepository) UpdateWithTx(ctx context.Context, tx *sqlx.Tx, p *entity.Payment) error {
	query := `UPDATE payments SET
		payer_name = ?, payer_email = ?, payer_cpf = ?,
		gross_amount = ?, discount_amount = ?, net_amount = ?, gateway_fee = ?, refunded_amount = ?,
		gateway_payment_id = ?, gateway_customer_id = ?,
		gateway_invoice_url = ?, gateway_metadata = ?,
		status = ?, paid_at = ?, refunded_at = ?, cancelled_at = ?,
		updated_at = NOW()
		WHERE id = ?`
	_, err := tx.ExecContext(ctx, query,
		p.PayerName, p.PayerEmail, p.PayerCPF,
		p.GrossAmount, p.DiscountAmount, p.NetAmount, p.GatewayFee, p.RefundedAmount,
		p.GatewayPaymentID, p.GatewayCustomerID,
		p.GatewayInvoiceURL, p.GatewayMetadata,
		p.Status, p.PaidAt, p.RefundedAt, p.CancelledAt,
		p.ID,
	)
	return err
}

func (r *paymentMySQLRepository) UpdateStatus(ctx context.Context, id, status string) error {
	query := `UPDATE payments SET status = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *paymentMySQLRepository) UpdateStatusWithTx(ctx context.Context, tx *sqlx.Tx, id, status string) error {
	query := `UPDATE payments SET status = ?, updated_at = NOW() WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, status, id)
	return err
}
