package gateway

import "time"

// Canonical payment status constants (normalized across gateways)
const (
	StatusPending    = "pending"
	StatusConfirmed  = "confirmed"
	StatusReceived   = "received"
	StatusOverdue    = "overdue"
	StatusRefunded   = "refunded"
	StatusFailed     = "failed"
	StatusCancelled  = "cancelled"
	StatusChargeback = "chargeback"
)

// Canonical billing type constants
const (
	BillingPIX        = "pix"
	BillingBoleto     = "boleto"
	BillingCreditCard = "credit_card"
	BillingDebitCard  = "debit_card"
)

// CreateCustomerRequest is the gateway-agnostic customer creation request.
type CreateCustomerRequest struct {
	Name     string
	Email    string
	Document string // CPF or CNPJ
	Phone    string
}

// CustomerResponse is the gateway-agnostic customer response.
type CustomerResponse struct {
	GatewayID string
	Name      string
	Email     string
	Document  string
}

// CreatePaymentRequest is the gateway-agnostic payment creation request (PIX/Boleto).
type CreatePaymentRequest struct {
	CustomerGatewayID string
	Amount            float64
	Description       string
	DueDate           time.Time
	ExternalReference string
}

// CreateCardPaymentRequest extends CreatePaymentRequest with card data.
type CreateCardPaymentRequest struct {
	CreatePaymentRequest
	CardNumber   string
	CardExpMonth string
	CardExpYear  string
	CardCVV      string
	HolderName   string
	HolderEmail  string
	HolderDoc    string
	HolderZip    string
	HolderPhone  string
	Installments int
}

// PaymentResponse is the gateway-agnostic payment response.
type PaymentResponse struct {
	GatewayPaymentID string
	Status           string // Canonical status
	GatewayRawStatus string // Original gateway status
	Amount           float64
	NetAmount        float64
	BillingType      string
	DueDate          string
	PaidAt           *time.Time
	ConfirmedAt      *time.Time
	InvoiceURL       string

	// PIX
	PixQRCodeBase64 string
	PixCopyPaste    string
	PixExpiration   string

	// Boleto
	BoletoURL     string
	BoletoBarCode string

	// Card
	TransactionReceiptURL string
}

// WebhookEvent is the gateway-agnostic webhook event.
type WebhookEvent struct {
	EventType        string // Canonical event type
	GatewayEvent     string // Original gateway event
	GatewayName      string // Which gateway sent the event
	PaymentID        string // gateway_payment_id
	CustomerID       string
	Amount           float64
	NetAmount        float64
	Status           string // Canonical status
	GatewayRawStatus string
	BillingType      string
	ExternalRef      string
	PaidAt           *time.Time
	RawPayload       []byte
}

// GatewayFees holds fee configuration for a gateway.
type GatewayFees struct {
	PixPercent  float64 // e.g. 0.0099 = 0.99%
	BoletoFixed float64 // e.g. 2.99
	CardPercent float64 // e.g. 0.0299 = 2.99%
	CardFixed   float64 // e.g. 0.49
}

// Canonical webhook event types
const (
	EventPaymentCreated   = "payment_created"
	EventPaymentConfirmed = "payment_confirmed"
	EventPaymentReceived  = "payment_received"
	EventPaymentOverdue   = "payment_overdue"
	EventPaymentRefunded  = "payment_refunded"
	EventPaymentDeleted   = "payment_deleted"
	EventPaymentFailed    = "payment_failed"
	EventPaymentChargeback = "payment_chargeback"
)
