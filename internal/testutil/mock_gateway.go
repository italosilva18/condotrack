package testutil

import (
	"context"

	"github.com/condotrack/api/internal/domain/gateway"
)

// MockGateway is a mock implementation of gateway.PaymentGateway for testing.
type MockGateway struct {
	NameFunc                    func() string
	CreateCustomerFunc          func(ctx context.Context, req gateway.CreateCustomerRequest) (*gateway.CustomerResponse, error)
	FindCustomerByDocumentFunc  func(ctx context.Context, document string) (*gateway.CustomerResponse, error)
	CreatePixPaymentFunc        func(ctx context.Context, req gateway.CreatePaymentRequest) (*gateway.PaymentResponse, error)
	CreateBoletoPaymentFunc     func(ctx context.Context, req gateway.CreatePaymentRequest) (*gateway.PaymentResponse, error)
	CreateCardPaymentFunc       func(ctx context.Context, req gateway.CreateCardPaymentRequest) (*gateway.PaymentResponse, error)
	GetPaymentFunc              func(ctx context.Context, gatewayPaymentID string) (*gateway.PaymentResponse, error)
	RefundPaymentFunc           func(ctx context.Context, gatewayPaymentID string, amount float64) (*gateway.PaymentResponse, error)
	CancelPaymentFunc           func(ctx context.Context, gatewayPaymentID string) error
	ParseWebhookEventFunc       func(ctx context.Context, headers map[string]string, body []byte) (*gateway.WebhookEvent, error)
	ValidateWebhookSignatureFunc func(ctx context.Context, headers map[string]string, body []byte) bool
	GetFeesFunc                 func() gateway.GatewayFees
}

func (m *MockGateway) Name() string {
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return "mock"
}

func (m *MockGateway) CreateCustomer(ctx context.Context, req gateway.CreateCustomerRequest) (*gateway.CustomerResponse, error) {
	if m.CreateCustomerFunc != nil {
		return m.CreateCustomerFunc(ctx, req)
	}
	return &gateway.CustomerResponse{GatewayID: "cust_mock_123", Name: req.Name, Email: req.Email}, nil
}

func (m *MockGateway) FindCustomerByDocument(ctx context.Context, document string) (*gateway.CustomerResponse, error) {
	if m.FindCustomerByDocumentFunc != nil {
		return m.FindCustomerByDocumentFunc(ctx, document)
	}
	return nil, nil
}

func (m *MockGateway) CreatePixPayment(ctx context.Context, req gateway.CreatePaymentRequest) (*gateway.PaymentResponse, error) {
	if m.CreatePixPaymentFunc != nil {
		return m.CreatePixPaymentFunc(ctx, req)
	}
	return &gateway.PaymentResponse{
		GatewayPaymentID: "pay_pix_mock_123",
		Status:           gateway.StatusPending,
		Amount:           req.Amount,
		BillingType:      gateway.BillingPIX,
	}, nil
}

func (m *MockGateway) CreateBoletoPayment(ctx context.Context, req gateway.CreatePaymentRequest) (*gateway.PaymentResponse, error) {
	if m.CreateBoletoPaymentFunc != nil {
		return m.CreateBoletoPaymentFunc(ctx, req)
	}
	return &gateway.PaymentResponse{
		GatewayPaymentID: "pay_boleto_mock_123",
		Status:           gateway.StatusPending,
		Amount:           req.Amount,
		BillingType:      gateway.BillingBoleto,
	}, nil
}

func (m *MockGateway) CreateCardPayment(ctx context.Context, req gateway.CreateCardPaymentRequest) (*gateway.PaymentResponse, error) {
	if m.CreateCardPaymentFunc != nil {
		return m.CreateCardPaymentFunc(ctx, req)
	}
	return &gateway.PaymentResponse{
		GatewayPaymentID: "pay_card_mock_123",
		Status:           gateway.StatusConfirmed,
		Amount:           req.Amount,
		BillingType:      gateway.BillingCreditCard,
	}, nil
}

func (m *MockGateway) GetPayment(ctx context.Context, gatewayPaymentID string) (*gateway.PaymentResponse, error) {
	if m.GetPaymentFunc != nil {
		return m.GetPaymentFunc(ctx, gatewayPaymentID)
	}
	return &gateway.PaymentResponse{
		GatewayPaymentID: gatewayPaymentID,
		Status:           gateway.StatusPending,
		Amount:           100,
	}, nil
}

func (m *MockGateway) RefundPayment(ctx context.Context, gatewayPaymentID string, amount float64) (*gateway.PaymentResponse, error) {
	if m.RefundPaymentFunc != nil {
		return m.RefundPaymentFunc(ctx, gatewayPaymentID, amount)
	}
	return &gateway.PaymentResponse{
		GatewayPaymentID: gatewayPaymentID,
		Status:           gateway.StatusRefunded,
	}, nil
}

func (m *MockGateway) CancelPayment(ctx context.Context, gatewayPaymentID string) error {
	if m.CancelPaymentFunc != nil {
		return m.CancelPaymentFunc(ctx, gatewayPaymentID)
	}
	return nil
}

func (m *MockGateway) ParseWebhookEvent(ctx context.Context, headers map[string]string, body []byte) (*gateway.WebhookEvent, error) {
	if m.ParseWebhookEventFunc != nil {
		return m.ParseWebhookEventFunc(ctx, headers, body)
	}
	return &gateway.WebhookEvent{}, nil
}

func (m *MockGateway) ValidateWebhookSignature(ctx context.Context, headers map[string]string, body []byte) bool {
	if m.ValidateWebhookSignatureFunc != nil {
		return m.ValidateWebhookSignatureFunc(ctx, headers, body)
	}
	return true
}

func (m *MockGateway) GetFees() gateway.GatewayFees {
	if m.GetFeesFunc != nil {
		return m.GetFeesFunc()
	}
	return gateway.GatewayFees{
		PixPercent:  0.0099,
		BoletoFixed: 2.99,
		CardPercent: 0.0299,
		CardFixed:   0.49,
	}
}
