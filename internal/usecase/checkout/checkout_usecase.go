package checkout

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/condotrack/api/internal/config"
	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/gateway"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/condotrack/api/internal/infrastructure/database"
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
	CourseID       string `json:"course_id" binding:"required"`
	CourseName     string `json:"course_name" binding:"required"`
	InstructorID   string `json:"instructor_id,omitempty"`
	InstructorName string `json:"instructor_name,omitempty"`

	// Payment info
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	DiscountCode  string  `json:"discount_code,omitempty"`
	PaymentMethod string  `json:"payment_method" binding:"required"` // pix, boleto, card

	// Card info (required if payment_method is card)
	CardNumber   string `json:"card_number,omitempty"`
	CardExpMonth string `json:"card_exp_month,omitempty"`
	CardExpYear  string `json:"card_exp_year,omitempty"`
	CardCVV      string `json:"card_cvv,omitempty"`
	HolderName   string `json:"holder_name,omitempty"`
	HolderEmail  string `json:"holder_email,omitempty"`
	HolderDoc    string `json:"holder_doc,omitempty"`
	HolderZip    string `json:"holder_zip,omitempty"`
	HolderPhone  string `json:"holder_phone,omitempty"`
	Installments int    `json:"installments,omitempty"`
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
	BoletoURL     string `json:"boleto_url,omitempty"`
	BoletoBarCode string `json:"boleto_bar_code,omitempty"`
	BoletoDueDate string `json:"boleto_due_date,omitempty"`

	// Revenue split info
	GrossAmount      float64 `json:"gross_amount"`
	DiscountAmount   float64 `json:"discount_amount"`
	PaymentFee       float64 `json:"payment_fee"`
	NetAmount        float64 `json:"net_amount"`
	InstructorAmount float64 `json:"instructor_amount"`
	PlatformAmount   float64 `json:"platform_amount"`

	// Coupon info
	CouponCode string `json:"coupon_code,omitempty"`
}

// UseCase defines the checkout use case interface
type UseCase interface {
	CreateCheckout(ctx context.Context, req *CheckoutRequest) (*CheckoutResponse, error)
	GetCheckoutStatus(ctx context.Context, enrollmentID string) (*CheckoutResponse, error)
}

type checkoutUseCase struct {
	gw                gateway.PaymentGateway
	matriculaRepo     repository.MatriculaRepository
	paymentRepo       repository.PaymentRepository
	couponRepo        repository.CouponRepository
	paymentTxnRepo    repository.PaymentTransactionRepository
	db                *database.MySQL
	instructorPercent float64
	platformPercent   float64
}

// NewUseCase creates a new checkout use case
func NewUseCase(
	gw gateway.PaymentGateway,
	matriculaRepo repository.MatriculaRepository,
	paymentRepo repository.PaymentRepository,
	couponRepo repository.CouponRepository,
	paymentTxnRepo repository.PaymentTransactionRepository,
	db *database.MySQL,
	cfg *config.Config,
) UseCase {
	return &checkoutUseCase{
		gw:                gw,
		matriculaRepo:     matriculaRepo,
		paymentRepo:       paymentRepo,
		couponRepo:        couponRepo,
		paymentTxnRepo:    paymentTxnRepo,
		db:                db,
		instructorPercent: cfg.RevenueInstructorPercent,
		platformPercent:   cfg.RevenuePlatformPercent,
	}
}

// CreateCheckout creates a complete checkout with enrollment, payment record, and gateway charge
func (uc *checkoutUseCase) CreateCheckout(ctx context.Context, req *CheckoutRequest) (*CheckoutResponse, error) {
	// Validate payment method
	if req.PaymentMethod != "pix" && req.PaymentMethod != "boleto" && req.PaymentMethod != "card" {
		return nil, errors.New("invalid payment method. Use: pix, boleto, or card")
	}

	// Validate card info if payment method is card
	if req.PaymentMethod == "card" && (req.CardNumber == "" || req.CardCVV == "") {
		return nil, errors.New("credit card information is required for card payment")
	}

	// --- Coupon validation ---
	var coupon *entity.Coupon
	var discountAmount float64
	finalAmount := req.Amount

	if req.DiscountCode != "" {
		var err error
		coupon, err = uc.couponRepo.FindByCode(ctx, req.DiscountCode)
		if err != nil {
			return nil, err
		}
		if coupon == nil {
			return nil, errors.New("invalid coupon code")
		}

		// Check per-user usage limit
		if coupon.MaxUsesPerUser != nil && req.StudentID != "" {
			usageCount, err := uc.couponRepo.CountUsageByUser(ctx, coupon.ID, req.StudentID)
			if err != nil {
				return nil, err
			}
			if usageCount >= *coupon.MaxUsesPerUser {
				return nil, errors.New("coupon usage limit exceeded for this user")
			}
		}

		discountAmount = coupon.CalculateDiscount(req.Amount)
		if discountAmount <= 0 {
			return nil, errors.New("coupon is not applicable to this order")
		}
		finalAmount = req.Amount - discountAmount
	}

	// Start transaction
	tx, err := uc.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create or get gateway customer
	customer, err := uc.gw.CreateCustomer(ctx, gateway.CreateCustomerRequest{
		Name:     req.StudentName,
		Email:    req.StudentEmail,
		Document: req.StudentCPF,
		Phone:    req.StudentPhone,
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
		DiscountAmount:  discountAmount,
		FinalAmount:     finalAmount,
		PaymentMethod:   &req.PaymentMethod,
		EnrollmentDate:  time.Now(),
		Status:          entity.EnrollmentStatusPending,
		Progress:        0,
		AsaasCustomerID: &customer.GatewayID,
		CreatedAt:       time.Now(),
	}

	if err := uc.matriculaRepo.CreateWithTx(ctx, tx, enrollment); err != nil {
		return nil, err
	}

	// Create payment on gateway
	var gatewayResp *gateway.PaymentResponse
	dueDate := time.Now().AddDate(0, 0, 3) // 3 days from now
	description := "Matrícula: " + req.CourseName

	switch req.PaymentMethod {
	case "pix":
		gatewayResp, err = uc.gw.CreatePixPayment(ctx, gateway.CreatePaymentRequest{
			CustomerGatewayID: customer.GatewayID,
			Amount:            finalAmount,
			Description:       description,
			DueDate:           dueDate,
			ExternalReference: enrollmentID,
		})
	case "boleto":
		gatewayResp, err = uc.gw.CreateBoletoPayment(ctx, gateway.CreatePaymentRequest{
			CustomerGatewayID: customer.GatewayID,
			Amount:            finalAmount,
			Description:       description,
			DueDate:           dueDate,
			ExternalReference: enrollmentID,
		})
	case "card":
		gatewayResp, err = uc.gw.CreateCardPayment(ctx, gateway.CreateCardPaymentRequest{
			CreatePaymentRequest: gateway.CreatePaymentRequest{
				CustomerGatewayID: customer.GatewayID,
				Amount:            finalAmount,
				Description:       description,
				DueDate:           dueDate,
				ExternalReference: enrollmentID,
			},
			CardNumber:   req.CardNumber,
			CardExpMonth: req.CardExpMonth,
			CardExpYear:  req.CardExpYear,
			CardCVV:      req.CardCVV,
			HolderName:   req.HolderName,
			HolderEmail:  req.HolderEmail,
			HolderDoc:    req.HolderDoc,
			HolderZip:    req.HolderZip,
			HolderPhone:  req.HolderPhone,
			Installments: req.Installments,
		})
	}

	if err != nil {
		return nil, err
	}

	// Update enrollment with gateway payment info
	enrollment.AsaasPaymentID = &gatewayResp.GatewayPaymentID
	if err := uc.matriculaRepo.UpdateWithTx(ctx, tx, enrollment); err != nil {
		return nil, err
	}

	// Calculate fees using the gateway's fee config
	fees := uc.gw.GetFees()
	gatewayFee := calculateGatewayFee(finalAmount, req.PaymentMethod, fees)

	// Create payment record in payments table
	paymentID := uuid.New().String()
	var couponID *string
	if coupon != nil {
		couponID = &coupon.ID
	}
	gwPaymentID := gatewayResp.GatewayPaymentID
	invoiceURL := gatewayResp.InvoiceURL
	dueDatePtr := dueDate

	paymentRecord := &entity.Payment{
		ID:                paymentID,
		EnrollmentID:      enrollmentID,
		PayerUserID:       &req.StudentID,
		PayerName:         req.StudentName,
		PayerEmail:        req.StudentEmail,
		PayerCPF:          &req.StudentCPF,
		GrossAmount:       req.Amount,
		DiscountAmount:    discountAmount,
		NetAmount:         finalAmount,
		GatewayFee:        gatewayFee,
		PaymentMethod:     req.PaymentMethod,
		Gateway:           uc.gw.Name(),
		GatewayPaymentID:  &gwPaymentID,
		GatewayCustomerID: &customer.GatewayID,
		GatewayInvoiceURL: nilIfEmpty(invoiceURL),
		InstallmentCount:  maxInt(req.Installments, 1),
		Status:            entity.FinPaymentStatusPending,
		CouponID:          couponID,
		DueDate:           &dueDatePtr,
		CreatedAt:         time.Now(),
	}

	if err := uc.paymentRepo.CreateWithTx(ctx, tx, paymentRecord); err != nil {
		return nil, err
	}

	// Record coupon usage
	if coupon != nil {
		usage := &entity.CouponUsage{
			ID:              uuid.New().String(),
			CouponID:        coupon.ID,
			UserID:          &req.StudentID,
			PaymentID:       &paymentID,
			EnrollmentID:    &enrollmentID,
			CourseID:        &req.CourseID,
			DiscountType:    coupon.DiscountType,
			DiscountValue:   coupon.DiscountValue,
			DiscountApplied: discountAmount,
			OriginalAmount:  req.Amount,
			FinalAmount:     finalAmount,
		}
		if err := uc.couponRepo.CreateUsageWithTx(ctx, tx, usage); err != nil {
			log.Printf("Failed to create coupon usage: %v", err)
		}
		if err := uc.couponRepo.IncrementUsageWithTx(ctx, tx, coupon.ID); err != nil {
			log.Printf("Failed to increment coupon usage: %v", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Log payment transaction (outside tx, non-critical)
	uc.logPaymentCreated(ctx, paymentRecord, gatewayResp)

	// Calculate revenue split for response
	netAfterFee := finalAmount - gatewayFee
	instructorAmount := netAfterFee * (uc.instructorPercent / 100)
	platformAmount := netAfterFee * (uc.platformPercent / 100)

	// Build response
	response := &CheckoutResponse{
		EnrollmentID:     enrollmentID,
		PaymentID:        paymentID,
		Status:           gatewayResp.Status,
		GrossAmount:      req.Amount,
		DiscountAmount:   discountAmount,
		PaymentFee:       gatewayFee,
		NetAmount:        netAfterFee,
		InstructorAmount: instructorAmount,
		PlatformAmount:   platformAmount,
	}

	if coupon != nil {
		response.CouponCode = coupon.Code
	}

	// Add payment-specific info
	if req.PaymentMethod == "pix" {
		response.PixQRCode = gatewayResp.PixQRCodeBase64
		response.PixCopyPaste = gatewayResp.PixCopyPaste
		response.PixExpirationDate = gatewayResp.PixExpiration
	}

	if req.PaymentMethod == "boleto" {
		response.BoletoURL = gatewayResp.BoletoURL
		response.BoletoBarCode = gatewayResp.BoletoBarCode
		response.BoletoDueDate = gatewayResp.DueDate
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

	// Try to get payment info from payments table
	payments, err := uc.paymentRepo.FindByEnrollmentID(ctx, enrollmentID)
	if err == nil && len(payments) > 0 {
		p := payments[0] // Most recent payment
		response.PaymentID = p.ID
		response.Status = p.Status
		response.GrossAmount = p.GrossAmount
		response.DiscountAmount = p.DiscountAmount

		// Get live status from gateway if we have a gateway payment ID
		if p.GatewayPaymentID != nil && *p.GatewayPaymentID != "" {
			gwPayment, err := uc.gw.GetPayment(ctx, *p.GatewayPaymentID)
			if err == nil {
				response.Status = gwPayment.Status
				if gwPayment.PixQRCodeBase64 != "" {
					response.PixQRCode = gwPayment.PixQRCodeBase64
					response.PixCopyPaste = gwPayment.PixCopyPaste
				}
				if gwPayment.BoletoURL != "" {
					response.BoletoURL = gwPayment.BoletoURL
				}
			}
		}
	} else {
		// Fallback: get info from gateway via enrollment's asaas payment ID
		if enrollment.AsaasPaymentID != nil && *enrollment.AsaasPaymentID != "" {
			gwPayment, err := uc.gw.GetPayment(ctx, *enrollment.AsaasPaymentID)
			if err == nil {
				response.PaymentID = gwPayment.GatewayPaymentID
				response.Status = gwPayment.Status
				if gwPayment.PixQRCodeBase64 != "" {
					response.PixQRCode = gwPayment.PixQRCodeBase64
					response.PixCopyPaste = gwPayment.PixCopyPaste
				}
				if gwPayment.BoletoURL != "" {
					response.BoletoURL = gwPayment.BoletoURL
				}
			}
		}
	}

	// Calculate revenue split
	paymentMethod := "pix"
	if enrollment.PaymentMethod != nil {
		paymentMethod = *enrollment.PaymentMethod
	}

	fees := uc.gw.GetFees()
	gatewayFee := calculateGatewayFee(enrollment.FinalAmount, paymentMethod, fees)
	netAfterFee := enrollment.FinalAmount - gatewayFee

	response.PaymentFee = gatewayFee
	response.NetAmount = netAfterFee
	response.InstructorAmount = netAfterFee * (uc.instructorPercent / 100)
	response.PlatformAmount = netAfterFee * (uc.platformPercent / 100)

	return response, nil
}

// logPaymentCreated creates a transaction log for payment creation (non-critical)
func (uc *checkoutUseCase) logPaymentCreated(ctx context.Context, payment *entity.Payment, gwResp *gateway.PaymentResponse) {
	txLog := &entity.PaymentTransaction{
		ID:          uuid.New().String(),
		PaymentID:   payment.ID,
		NewStatus:   payment.Status,
		EventSource: entity.EventSourceAPI,
		EventType:   entity.TxEventCreated,
	}
	amount := payment.GrossAmount
	txLog.Amount = &amount

	gwEvent := "checkout_created"
	txLog.GatewayEventID = &gwEvent

	if err := uc.paymentTxnRepo.Create(ctx, txLog); err != nil {
		log.Printf("Failed to log payment transaction: %v", err)
	}
}

// calculateGatewayFee calculates fee based on payment method and gateway fees
func calculateGatewayFee(amount float64, method string, fees gateway.GatewayFees) float64 {
	switch method {
	case "pix", "PIX":
		return amount * fees.PixPercent
	case "boleto", "BOLETO":
		return fees.BoletoFixed
	case "card", "credit_card", "CREDIT_CARD":
		return (amount * fees.CardPercent) + fees.CardFixed
	default:
		return 0
	}
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
