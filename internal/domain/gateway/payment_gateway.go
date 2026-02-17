package gateway

import "context"

// PaymentGateway defines the interface that every payment gateway must implement.
// This allows swapping between Asaas, Mercado Pago, or any other gateway
// without changing business logic.
type PaymentGateway interface {
	// Name returns the gateway identifier (e.g. "asaas", "mercadopago")
	Name() string

	// Customers
	CreateCustomer(ctx context.Context, req CreateCustomerRequest) (*CustomerResponse, error)
	FindCustomerByDocument(ctx context.Context, document string) (*CustomerResponse, error)

	// Payments
	CreatePixPayment(ctx context.Context, req CreatePaymentRequest) (*PaymentResponse, error)
	CreateBoletoPayment(ctx context.Context, req CreatePaymentRequest) (*PaymentResponse, error)
	CreateCardPayment(ctx context.Context, req CreateCardPaymentRequest) (*PaymentResponse, error)
	GetPayment(ctx context.Context, gatewayPaymentID string) (*PaymentResponse, error)

	// Refund / Cancel
	RefundPayment(ctx context.Context, gatewayPaymentID string, amount float64) (*PaymentResponse, error)
	CancelPayment(ctx context.Context, gatewayPaymentID string) error

	// Webhook
	ParseWebhookEvent(ctx context.Context, headers map[string]string, body []byte) (*WebhookEvent, error)
	ValidateWebhookSignature(ctx context.Context, headers map[string]string, body []byte) bool

	// Fees
	GetFees() GatewayFees
}
