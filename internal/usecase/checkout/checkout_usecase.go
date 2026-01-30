package checkout

import (
	"context"
	"errors"
	"time"

	"github.com/condotrack/api/internal/config"
	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/condotrack/api/internal/infrastructure/database"
	"github.com/condotrack/api/internal/infrastructure/external/asaas"
	"github.com/google/uuid"
)

// CheckoutRequest represents the request to create a checkout
type CheckoutRequest struct {
	// Student info
	StudentID    string `json:"student_id" binding:"required"`
	StudentName  string `json:"student_name" binding:"required"`
	StudentEmail string `json:"student_email" binding:"required,email"`
	StudentCPF   string `json:"student_cpf" binding:"required"`
	StudentPhone string `json:"student_phone,omitempty"`

	// Course info
	CourseID       string  `json:"course_id" binding:"required"`
	CourseName     string  `json:"course_name" binding:"required"`
	InstructorID   string  `json:"instructor_id,omitempty"`
	InstructorName string  `json:"instructor_name,omitempty"`

	// Payment info
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	DiscountCode  string  `json:"discount_code,omitempty"`
	PaymentMethod string  `json:"payment_method" binding:"required"` // pix, boleto, card

	// Card info (required if payment_method is card)
	CreditCard *asaas.CreditCard       `json:"credit_card,omitempty"`
	CardHolder *asaas.CreditCardHolder `json:"card_holder,omitempty"`
}

// CheckoutResponse represents the checkout response
type CheckoutResponse struct {
	EnrollmentID string `json:"enrollment_id"`
	PaymentID    string `json:"payment_id"`
	Status       string `json:"status"`

	// PIX specific
	PixQRCode         string `json:"pix_qr_code,omitempty"`
	PixCopyPaste      string `json:"pix_copy_paste,omitempty"`
	PixExpirationDate string `json:"pix_expiration_date,omitempty"`

	// Boleto specific
	BoletoURL      string `json:"boleto_url,omitempty"`
	BoletoBarCode  string `json:"boleto_bar_code,omitempty"`
	BoletoDueDate  string `json:"boleto_due_date,omitempty"`

	// Revenue split info
	GrossAmount      float64 `json:"gross_amount"`
	PaymentFee       float64 `json:"payment_fee"`
	NetAmount        float64 `json:"net_amount"`
	InstructorAmount float64 `json:"instructor_amount"`
	PlatformAmount   float64 `json:"platform_amount"`
}

// UseCase defines the checkout use case interface
type UseCase interface {
	CreateCheckout(ctx context.Context, req *CheckoutRequest) (*CheckoutResponse, error)
	GetCheckoutStatus(ctx context.Context, enrollmentID string) (*CheckoutResponse, error)
}

type checkoutUseCase struct {
	asaasClient       *asaas.Client
	matriculaRepo     repository.MatriculaRepository
	db                *database.MySQL
	instructorPercent float64
	platformPercent   float64
}

// NewUseCase creates a new checkout use case
func NewUseCase(
	asaasClient *asaas.Client,
	matriculaRepo repository.MatriculaRepository,
	db *database.MySQL,
	cfg *config.Config,
) UseCase {
	return &checkoutUseCase{
		asaasClient:       asaasClient,
		matriculaRepo:     matriculaRepo,
		db:                db,
		instructorPercent: cfg.RevenueInstructorPercent,
		platformPercent:   cfg.RevenuePlatformPercent,
	}
}

// CreateCheckout creates a complete checkout with enrollment and payment
func (uc *checkoutUseCase) CreateCheckout(ctx context.Context, req *CheckoutRequest) (*CheckoutResponse, error) {
	// Validate payment method
	if req.PaymentMethod != "pix" && req.PaymentMethod != "boleto" && req.PaymentMethod != "card" {
		return nil, errors.New("invalid payment method. Use: pix, boleto, or card")
	}

	// Validate card info if payment method is card
	if req.PaymentMethod == "card" && (req.CreditCard == nil || req.CardHolder == nil) {
		return nil, errors.New("credit card information is required for card payment")
	}

	// Start transaction
	tx, err := uc.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create or get Asaas customer
	customer, err := uc.asaasClient.FindOrCreateCustomer(ctx, &asaas.CreateCustomerRequest{
		Name:    req.StudentName,
		Email:   req.StudentEmail,
		CPFCnpj: req.StudentCPF,
		Phone:   req.StudentPhone,
	})
	if err != nil {
		return nil, err
	}

	// Create enrollment
	enrollmentID := uuid.New().String()
	enrollment := &entity.Matricula{
		ID:              enrollmentID,
		StudentID:       req.StudentID,
		StudentName:     req.StudentName,
		StudentEmail:    req.StudentEmail,
		StudentCPF:      &req.StudentCPF,
		StudentPhone:    &req.StudentPhone,
		CourseID:        req.CourseID,
		CourseName:      req.CourseName,
		InstructorID:    &req.InstructorID,
		InstructorName:  &req.InstructorName,
		PaymentStatus:   entity.PaymentStatusPending,
		Amount:          req.Amount,
		DiscountAmount:  0,
		FinalAmount:     req.Amount,
		PaymentMethod:   &req.PaymentMethod,
		EnrollmentDate:  time.Now(),
		Status:          entity.EnrollmentStatusPending,
		Progress:        0,
		AsaasCustomerID: &customer.ID,
		CreatedAt:       time.Now(),
	}

	if err := uc.matriculaRepo.CreateWithTx(ctx, tx, enrollment); err != nil {
		return nil, err
	}

	// Create payment based on method
	var paymentResponse *asaas.PaymentResponse
	dueDate := time.Now().AddDate(0, 0, 3).Format("2006-01-02") // 3 days from now

	switch req.PaymentMethod {
	case "pix":
		paymentResponse, err = uc.asaasClient.CreatePayment(ctx, &asaas.CreatePaymentRequest{
			Customer:     customer.ID,
			BillingType:  "PIX",
			Value:        req.Amount,
			DueDate:      dueDate,
			Description:  "Matrícula: " + req.CourseName,
			ExternalRef:  enrollmentID,
		})
	case "boleto":
		paymentResponse, err = uc.asaasClient.CreatePayment(ctx, &asaas.CreatePaymentRequest{
			Customer:     customer.ID,
			BillingType:  "BOLETO",
			Value:        req.Amount,
			DueDate:      dueDate,
			Description:  "Matrícula: " + req.CourseName,
			ExternalRef:  enrollmentID,
		})
	case "card":
		paymentResponse, err = uc.asaasClient.CreateCardPayment(ctx, &asaas.CreateCardPaymentRequest{
			Customer:         customer.ID,
			BillingType:      "CREDIT_CARD",
			Value:            req.Amount,
			DueDate:          dueDate,
			Description:      "Matrícula: " + req.CourseName,
			ExternalRef:      enrollmentID,
			CreditCard:       req.CreditCard,
			CreditCardHolder: req.CardHolder,
		})
	}

	if err != nil {
		return nil, err
	}

	// Update enrollment with payment info
	enrollment.AsaasPaymentID = &paymentResponse.ID
	if err := uc.matriculaRepo.UpdateWithTx(ctx, tx, enrollment); err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Calculate revenue split
	split := entity.CalculateRevenueSplit(
		req.Amount,
		req.PaymentMethod,
		uc.instructorPercent,
		uc.platformPercent,
	)

	// Build response
	response := &CheckoutResponse{
		EnrollmentID:     enrollmentID,
		PaymentID:        paymentResponse.ID,
		Status:           paymentResponse.Status,
		GrossAmount:      split.GrossAmount,
		PaymentFee:       split.PaymentFee,
		NetAmount:        split.NetAmount,
		InstructorAmount: split.InstructorAmount,
		PlatformAmount:   split.PlatformAmount,
	}

	// Add payment-specific info
	if req.PaymentMethod == "pix" && paymentResponse.PixQRCode != nil {
		response.PixQRCode = paymentResponse.PixQRCode.EncodedImage
		response.PixCopyPaste = paymentResponse.PixQRCode.Payload
		response.PixExpirationDate = paymentResponse.PixQRCode.ExpirationDate
	}

	if req.PaymentMethod == "boleto" {
		response.BoletoURL = paymentResponse.BankSlipURL
		response.BoletoDueDate = paymentResponse.DueDate
	}

	return response, nil
}

// GetCheckoutStatus retrieves the status of a checkout
func (uc *checkoutUseCase) GetCheckoutStatus(ctx context.Context, enrollmentID string) (*CheckoutResponse, error) {
	enrollment, err := uc.matriculaRepo.FindByID(ctx, enrollmentID)
	if err != nil {
		return nil, err
	}
	if enrollment == nil {
		return nil, errors.New("enrollment not found")
	}

	response := &CheckoutResponse{
		EnrollmentID: enrollmentID,
		Status:       enrollment.PaymentStatus,
	}

	// Get payment info from Asaas if available
	if enrollment.AsaasPaymentID != nil && *enrollment.AsaasPaymentID != "" {
		payment, err := uc.asaasClient.GetPayment(ctx, *enrollment.AsaasPaymentID)
		if err == nil {
			response.PaymentID = payment.ID
			response.Status = payment.Status

			if payment.PixQRCode != nil {
				response.PixQRCode = payment.PixQRCode.EncodedImage
				response.PixCopyPaste = payment.PixQRCode.Payload
			}
			if payment.BankSlipURL != "" {
				response.BoletoURL = payment.BankSlipURL
			}
		}
	}

	// Calculate revenue split
	paymentMethod := "pix"
	if enrollment.PaymentMethod != nil {
		paymentMethod = *enrollment.PaymentMethod
	}

	split := entity.CalculateRevenueSplit(
		enrollment.FinalAmount,
		paymentMethod,
		uc.instructorPercent,
		uc.platformPercent,
	)

	response.GrossAmount = split.GrossAmount
	response.PaymentFee = split.PaymentFee
	response.NetAmount = split.NetAmount
	response.InstructorAmount = split.InstructorAmount
	response.PlatformAmount = split.PlatformAmount

	return response, nil
}
