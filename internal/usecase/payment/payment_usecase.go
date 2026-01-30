package payment

import (
	"context"
	"errors"

	"github.com/condotrack/api/internal/config"
	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/infrastructure/external/asaas"
)

// UseCase defines the payment use case interface
type UseCase interface {
	CreateCustomer(ctx context.Context, req *asaas.CreateCustomerRequest) (*asaas.Customer, error)
	CreatePixPayment(ctx context.Context, req *asaas.CreatePaymentRequest) (*asaas.PaymentResponse, error)
	CreateBoletoPayment(ctx context.Context, req *asaas.CreatePaymentRequest) (*asaas.PaymentResponse, error)
	CreateCardPayment(ctx context.Context, req *asaas.CreateCardPaymentRequest) (*asaas.PaymentResponse, error)
	GetPaymentStatus(ctx context.Context, paymentID string) (*asaas.PaymentResponse, error)
	SimulateRevenueSplit(req *entity.CalculateSplitRequest) *entity.CalculateSplitResponse
}

type paymentUseCase struct {
	asaasClient       *asaas.Client
	instructorPercent float64
	platformPercent   float64
}

// NewUseCase creates a new payment use case
func NewUseCase(asaasClient *asaas.Client, cfg *config.Config) UseCase {
	return &paymentUseCase{
		asaasClient:       asaasClient,
		instructorPercent: cfg.RevenueInstructorPercent,
		platformPercent:   cfg.RevenuePlatformPercent,
	}
}

// CreateCustomer creates a new customer in Asaas
func (uc *paymentUseCase) CreateCustomer(ctx context.Context, req *asaas.CreateCustomerRequest) (*asaas.Customer, error) {
	if req.Name == "" {
		return nil, errors.New("customer name is required")
	}
	if req.CPFCnpj == "" {
		return nil, errors.New("customer CPF/CNPJ is required")
	}

	return uc.asaasClient.CreateCustomer(ctx, req)
}

// CreatePixPayment creates a PIX payment in Asaas
func (uc *paymentUseCase) CreatePixPayment(ctx context.Context, req *asaas.CreatePaymentRequest) (*asaas.PaymentResponse, error) {
	if req.Customer == "" {
		return nil, errors.New("customer ID is required")
	}
	if req.Value <= 0 {
		return nil, errors.New("payment value must be greater than 0")
	}

	req.BillingType = "PIX"
	return uc.asaasClient.CreatePayment(ctx, req)
}

// CreateBoletoPayment creates a Boleto payment in Asaas
func (uc *paymentUseCase) CreateBoletoPayment(ctx context.Context, req *asaas.CreatePaymentRequest) (*asaas.PaymentResponse, error) {
	if req.Customer == "" {
		return nil, errors.New("customer ID is required")
	}
	if req.Value <= 0 {
		return nil, errors.New("payment value must be greater than 0")
	}

	req.BillingType = "BOLETO"
	return uc.asaasClient.CreatePayment(ctx, req)
}

// CreateCardPayment creates a credit card payment in Asaas
func (uc *paymentUseCase) CreateCardPayment(ctx context.Context, req *asaas.CreateCardPaymentRequest) (*asaas.PaymentResponse, error) {
	if req.Customer == "" {
		return nil, errors.New("customer ID is required")
	}
	if req.Value <= 0 {
		return nil, errors.New("payment value must be greater than 0")
	}
	if req.CreditCard == nil {
		return nil, errors.New("credit card information is required")
	}

	return uc.asaasClient.CreateCardPayment(ctx, req)
}

// GetPaymentStatus retrieves the status of a payment
func (uc *paymentUseCase) GetPaymentStatus(ctx context.Context, paymentID string) (*asaas.PaymentResponse, error) {
	if paymentID == "" {
		return nil, errors.New("payment ID is required")
	}

	return uc.asaasClient.GetPayment(ctx, paymentID)
}

// SimulateRevenueSplit simulates how revenue would be split
func (uc *paymentUseCase) SimulateRevenueSplit(req *entity.CalculateSplitRequest) *entity.CalculateSplitResponse {
	instructorPercent := req.InstructorPercent
	if instructorPercent == 0 {
		instructorPercent = uc.instructorPercent
	}

	platformPercent := req.PlatformPercent
	if platformPercent == 0 {
		platformPercent = uc.platformPercent
	}

	result := entity.CalculateRevenueSplit(
		req.GrossAmount,
		req.PaymentMethod,
		instructorPercent,
		platformPercent,
	)

	return &result
}
