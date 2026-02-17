package payment

import (
	"context"
	"errors"
	"time"

	"github.com/condotrack/api/internal/config"
	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/gateway"
	"github.com/condotrack/api/internal/domain/repository"
)

// UseCase defines the payment use case interface
type UseCase interface {
	CreateCustomer(ctx context.Context, req *CreateCustomerRequest) (*gateway.CustomerResponse, error)
	CreatePixPayment(ctx context.Context, req *CreatePaymentRequest) (*gateway.PaymentResponse, error)
	CreateBoletoPayment(ctx context.Context, req *CreatePaymentRequest) (*gateway.PaymentResponse, error)
	CreateCardPayment(ctx context.Context, req *CreateCardPaymentRequest) (*gateway.PaymentResponse, error)
	GetPaymentStatus(ctx context.Context, paymentID string) (*PaymentStatusResponse, error)
	SimulateRevenueSplit(req *entity.CalculateSplitRequest) *entity.CalculateSplitResponse
	ListPayments(ctx context.Context, filters repository.PaymentFilters) ([]entity.Payment, int, error)
	GetPaymentsByEnrollment(ctx context.Context, enrollmentID string) ([]entity.Payment, error)
}

// CreateCustomerRequest is the handler-level customer request (gateway-agnostic)
type CreateCustomerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Document string `json:"cpf_cnpj" binding:"required"`
	Phone    string `json:"phone,omitempty"`
}

// CreatePaymentRequest is the handler-level payment request (gateway-agnostic)
type CreatePaymentRequest struct {
	CustomerGatewayID string  `json:"customer_id" binding:"required"`
	Amount            float64 `json:"value" binding:"required,gt=0"`
	Description       string  `json:"description,omitempty"`
	DueDate           string  `json:"due_date,omitempty"`
	ExternalRef       string  `json:"external_reference,omitempty"`
}

// CreateCardPaymentRequest extends CreatePaymentRequest with card data
type CreateCardPaymentRequest struct {
	CustomerGatewayID string  `json:"customer_id" binding:"required"`
	Amount            float64 `json:"value" binding:"required,gt=0"`
	Description       string  `json:"description,omitempty"`
	ExternalRef       string  `json:"external_reference,omitempty"`
	Installments      int     `json:"installments,omitempty"`

	// Card info
	HolderName  string `json:"holder_name" binding:"required"`
	CardNumber  string `json:"card_number" binding:"required"`
	ExpiryMonth string `json:"expiry_month" binding:"required"`
	ExpiryYear  string `json:"expiry_year" binding:"required"`
	CCV         string `json:"ccv" binding:"required"`

	// Holder info
	HolderEmail string `json:"holder_email" binding:"required,email"`
	HolderCPF   string `json:"holder_cpf" binding:"required"`
	HolderZip   string `json:"holder_postal_code" binding:"required"`
	HolderPhone string `json:"holder_phone,omitempty"`
}

// PaymentStatusResponse combines local DB data with gateway data
type PaymentStatusResponse struct {
	ID                string  `json:"id"`
	Status            string  `json:"status"`
	GatewayStatus     string  `json:"gateway_status,omitempty"`
	Amount            float64 `json:"amount"`
	NetAmount         float64 `json:"net_amount"`
	BillingType       string  `json:"billing_type,omitempty"`
	DueDate           string  `json:"due_date,omitempty"`
	Gateway           string  `json:"gateway"`
	GatewayPaymentID  string  `json:"gateway_payment_id,omitempty"`
}

type paymentUseCase struct {
	gw                gateway.PaymentGateway
	paymentRepo       repository.PaymentRepository
	instructorPercent float64
	platformPercent   float64
}

// NewUseCase creates a new payment use case
func NewUseCase(gw gateway.PaymentGateway, paymentRepo repository.PaymentRepository, cfg *config.Config) UseCase {
	return &paymentUseCase{
		gw:                gw,
		paymentRepo:       paymentRepo,
		instructorPercent: cfg.RevenueInstructorPercent,
		platformPercent:   cfg.RevenuePlatformPercent,
	}
}

// CreateCustomer creates a new customer via the gateway
func (uc *paymentUseCase) CreateCustomer(ctx context.Context, req *CreateCustomerRequest) (*gateway.CustomerResponse, error) {
	if req.Name == "" {
		return nil, errors.New("customer name is required")
	}
	if req.Document == "" {
		return nil, errors.New("customer CPF/CNPJ is required")
	}

	return uc.gw.CreateCustomer(ctx, gateway.CreateCustomerRequest{
		Name:     req.Name,
		Email:    req.Email,
		Document: req.Document,
		Phone:    req.Phone,
	})
}

// CreatePixPayment creates a PIX payment via the gateway
func (uc *paymentUseCase) CreatePixPayment(ctx context.Context, req *CreatePaymentRequest) (*gateway.PaymentResponse, error) {
	if req.CustomerGatewayID == "" {
		return nil, errors.New("customer ID is required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("payment value must be greater than 0")
	}

	gwReq, err := req.toGatewayRequest()
	if err != nil {
		return nil, err
	}

	return uc.gw.CreatePixPayment(ctx, *gwReq)
}

// CreateBoletoPayment creates a Boleto payment via the gateway
func (uc *paymentUseCase) CreateBoletoPayment(ctx context.Context, req *CreatePaymentRequest) (*gateway.PaymentResponse, error) {
	if req.CustomerGatewayID == "" {
		return nil, errors.New("customer ID is required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("payment value must be greater than 0")
	}

	gwReq, err := req.toGatewayRequest()
	if err != nil {
		return nil, err
	}

	return uc.gw.CreateBoletoPayment(ctx, *gwReq)
}

// CreateCardPayment creates a credit card payment via the gateway
func (uc *paymentUseCase) CreateCardPayment(ctx context.Context, req *CreateCardPaymentRequest) (*gateway.PaymentResponse, error) {
	if req.CustomerGatewayID == "" {
		return nil, errors.New("customer ID is required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("payment value must be greater than 0")
	}
	if req.CardNumber == "" {
		return nil, errors.New("credit card information is required")
	}

	gwReq := req.toGatewayRequest()
	return uc.gw.CreateCardPayment(ctx, gwReq)
}

// GetPaymentStatus retrieves the status of a payment (local DB + gateway)
func (uc *paymentUseCase) GetPaymentStatus(ctx context.Context, paymentID string) (*PaymentStatusResponse, error) {
	if paymentID == "" {
		return nil, errors.New("payment ID is required")
	}

	// Try local DB first
	payment, err := uc.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	if payment != nil {
		resp := &PaymentStatusResponse{
			ID:        payment.ID,
			Status:    payment.Status,
			Amount:    payment.GrossAmount,
			NetAmount: payment.NetAmount,
			Gateway:   payment.Gateway,
		}
		if payment.GatewayPaymentID != nil {
			resp.GatewayPaymentID = *payment.GatewayPaymentID

			// Get live status from gateway
			gwPayment, err := uc.gw.GetPayment(ctx, *payment.GatewayPaymentID)
			if err == nil {
				resp.GatewayStatus = gwPayment.Status
				resp.BillingType = gwPayment.BillingType
				resp.DueDate = gwPayment.DueDate
			}
		}
		return resp, nil
	}

	// Fallback: try as gateway payment ID
	gwPayment, err := uc.gw.GetPayment(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	return &PaymentStatusResponse{
		ID:               gwPayment.GatewayPaymentID,
		Status:           gwPayment.Status,
		GatewayStatus:    gwPayment.GatewayRawStatus,
		Amount:           gwPayment.Amount,
		NetAmount:        gwPayment.NetAmount,
		BillingType:      gwPayment.BillingType,
		DueDate:          gwPayment.DueDate,
		Gateway:          uc.gw.Name(),
		GatewayPaymentID: gwPayment.GatewayPaymentID,
	}, nil
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

// ListPayments retrieves payments from the payments table
func (uc *paymentUseCase) ListPayments(ctx context.Context, filters repository.PaymentFilters) ([]entity.Payment, int, error) {
	return uc.paymentRepo.FindAll(ctx, filters)
}

// GetPaymentsByEnrollment retrieves payments for an enrollment
func (uc *paymentUseCase) GetPaymentsByEnrollment(ctx context.Context, enrollmentID string) ([]entity.Payment, error) {
	return uc.paymentRepo.FindByEnrollmentID(ctx, enrollmentID)
}

// toGatewayRequest converts handler-level request to gateway request
func (req *CreatePaymentRequest) toGatewayRequest() (*gateway.CreatePaymentRequest, error) {
	gwReq := &gateway.CreatePaymentRequest{
		CustomerGatewayID: req.CustomerGatewayID,
		Amount:            req.Amount,
		Description:       req.Description,
		ExternalReference: req.ExternalRef,
	}

	if req.DueDate != "" {
		t, err := parseDateString(req.DueDate)
		if err != nil {
			return nil, errors.New("invalid due_date format (use YYYY-MM-DD)")
		}
		gwReq.DueDate = t
	}

	return gwReq, nil
}

func (req *CreateCardPaymentRequest) toGatewayRequest() gateway.CreateCardPaymentRequest {
	return gateway.CreateCardPaymentRequest{
		CreatePaymentRequest: gateway.CreatePaymentRequest{
			CustomerGatewayID: req.CustomerGatewayID,
			Amount:            req.Amount,
			Description:       req.Description,
			ExternalReference: req.ExternalRef,
		},
		CardNumber:   req.CardNumber,
		CardExpMonth: req.ExpiryMonth,
		CardExpYear:  req.ExpiryYear,
		CardCVV:      req.CCV,
		HolderName:   req.HolderName,
		HolderEmail:  req.HolderEmail,
		HolderDoc:    req.HolderCPF,
		HolderZip:    req.HolderZip,
		HolderPhone:  req.HolderPhone,
		Installments: req.Installments,
	}
}

func parseDateString(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}
