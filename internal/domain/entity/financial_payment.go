package entity

import "time"

// Payment represents a dedicated payment record in the payments table.
// This is the new entity decoupled from enrollments.
type Payment struct {
	ID                string     `db:"id" json:"id"`
	EnrollmentID      string     `db:"enrollment_id" json:"enrollment_id"`
	PayerUserID       *string    `db:"payer_user_id" json:"payer_user_id,omitempty"`
	PayerName         string     `db:"payer_name" json:"payer_name"`
	PayerEmail        string     `db:"payer_email" json:"payer_email"`
	PayerCPF          *string    `db:"payer_cpf" json:"payer_cpf,omitempty"`
	GrossAmount       float64    `db:"gross_amount" json:"gross_amount"`
	DiscountAmount    float64    `db:"discount_amount" json:"discount_amount"`
	NetAmount         float64    `db:"net_amount" json:"net_amount"`
	GatewayFee        float64    `db:"gateway_fee" json:"gateway_fee"`
	RefundedAmount    float64    `db:"refunded_amount" json:"refunded_amount"`
	PaymentMethod     string     `db:"payment_method" json:"payment_method"`
	Gateway           string     `db:"gateway" json:"gateway"`
	GatewayPaymentID  *string    `db:"gateway_payment_id" json:"gateway_payment_id,omitempty"`
	GatewayCustomerID *string    `db:"gateway_customer_id" json:"gateway_customer_id,omitempty"`
	GatewayInvoiceURL *string    `db:"gateway_invoice_url" json:"gateway_invoice_url,omitempty"`
	GatewayMetadata   *string    `db:"gateway_metadata" json:"gateway_metadata,omitempty"`
	InstallmentCount  int        `db:"installment_count" json:"installment_count"`
	InstallmentOf     *string    `db:"installment_of" json:"installment_of,omitempty"`
	InstallmentNumber *int       `db:"installment_number" json:"installment_number,omitempty"`
	Status            string     `db:"status" json:"status"`
	CouponID          *string    `db:"coupon_id" json:"coupon_id,omitempty"`
	DueDate           *time.Time `db:"due_date" json:"due_date,omitempty"`
	PaidAt            *time.Time `db:"paid_at" json:"paid_at,omitempty"`
	RefundedAt        *time.Time `db:"refunded_at" json:"refunded_at,omitempty"`
	CancelledAt       *time.Time `db:"cancelled_at" json:"cancelled_at,omitempty"`
	ExpiresAt         *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt         *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// Payment status constants matching the ENUM in 006_financial_module.sql
const (
	FinPaymentStatusPending          = "pending"
	FinPaymentStatusAwaitingPayment  = "awaiting_payment"
	FinPaymentStatusConfirmed        = "confirmed"
	FinPaymentStatusReceived         = "received"
	FinPaymentStatusOverdue          = "overdue"
	FinPaymentStatusRefundRequested  = "refund_requested"
	FinPaymentStatusRefunded         = "refunded"
	FinPaymentStatusPartiallyRefunded = "partially_refunded"
	FinPaymentStatusChargeback       = "chargeback"
	FinPaymentStatusFailed           = "failed"
	FinPaymentStatusCancelled        = "cancelled"
)

// Payment method constants
const (
	MethodCreditCard   = "credit_card"
	MethodDebitCard    = "debit_card"
	MethodBoleto       = "boleto"
	MethodPIX          = "pix"
	MethodBankTransfer = "bank_transfer"
	MethodFree         = "free"
)

// Gateway constants
const (
	GatewayAsaas      = "asaas"
	GatewayMercadoPago = "mercadopago"
	GatewayManual     = "manual"
	GatewayFree       = "free"
)
