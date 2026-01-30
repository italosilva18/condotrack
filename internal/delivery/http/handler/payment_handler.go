package handler

import (
	"strconv"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/condotrack/api/internal/infrastructure/external/asaas"
	"github.com/condotrack/api/internal/usecase/payment"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct {
	usecase       payment.UseCase
	matriculaRepo repository.MatriculaRepository
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(uc payment.UseCase, matriculaRepo repository.MatriculaRepository) *PaymentHandler {
	return &PaymentHandler{
		usecase:       uc,
		matriculaRepo: matriculaRepo,
	}
}

// CreateCustomerRequest represents the request to create a customer
type CreateCustomerRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	CPFCnpj   string `json:"cpf_cnpj" binding:"required"`
	Phone     string `json:"phone,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Address   string `json:"address,omitempty"`
}

// CreatePaymentRequest represents the request to create a payment
type CreatePaymentRequest struct {
	CustomerID  string  `json:"customer_id" binding:"required"`
	Value       float64 `json:"value" binding:"required,gt=0"`
	Description string  `json:"description,omitempty"`
	DueDate     string  `json:"due_date,omitempty"` // Format: YYYY-MM-DD
	ExternalRef string  `json:"external_reference,omitempty"`
}

// CreateCardPaymentRequest represents the request to create a card payment
type CreateCardPaymentRequest struct {
	CustomerID   string  `json:"customer_id" binding:"required"`
	Value        float64 `json:"value" binding:"required,gt=0"`
	Description  string  `json:"description,omitempty"`
	Installments int     `json:"installments,omitempty"`
	ExternalRef  string  `json:"external_reference,omitempty"`

	// Card info
	HolderName  string `json:"holder_name" binding:"required"`
	CardNumber  string `json:"card_number" binding:"required"`
	ExpiryMonth string `json:"expiry_month" binding:"required"`
	ExpiryYear  string `json:"expiry_year" binding:"required"`
	CCV         string `json:"ccv" binding:"required"`

	// Holder info
	HolderEmail      string `json:"holder_email" binding:"required,email"`
	HolderCPF        string `json:"holder_cpf" binding:"required"`
	HolderPostalCode string `json:"holder_postal_code" binding:"required"`
	HolderPhone      string `json:"holder_phone,omitempty"`
}

// CreateCustomer handles POST /api/v1/payments/customer
func (h *PaymentHandler) CreateCustomer(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	asaasReq := &asaas.CreateCustomerRequest{
		Name:       req.Name,
		Email:      req.Email,
		CPFCnpj:    req.CPFCnpj,
		Phone:      req.Phone,
		PostalCode: req.PostalCode,
		Address:    req.Address,
	}

	customer, err := h.usecase.CreateCustomer(ctx, asaasReq)
	if err != nil {
		response.InternalError(c, "Failed to create customer: "+err.Error())
		return
	}

	response.Created(c, customer)
}

// CreatePixPayment handles POST /api/v1/payments/pix
func (h *PaymentHandler) CreatePixPayment(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	asaasReq := &asaas.CreatePaymentRequest{
		Customer:    req.CustomerID,
		Value:       req.Value,
		Description: req.Description,
		DueDate:     req.DueDate,
		ExternalRef: req.ExternalRef,
	}

	payment, err := h.usecase.CreatePixPayment(ctx, asaasReq)
	if err != nil {
		response.InternalError(c, "Failed to create PIX payment: "+err.Error())
		return
	}

	// Build response with PIX details
	resp := gin.H{
		"id":          payment.ID,
		"status":      payment.Status,
		"value":       payment.Value,
		"due_date":    payment.DueDate,
		"description": payment.Description,
	}

	if payment.PixQRCode != nil {
		resp["pix_qr_code"] = payment.PixQRCode.EncodedImage
		resp["pix_copy_paste"] = payment.PixQRCode.Payload
		resp["pix_expiration"] = payment.PixQRCode.ExpirationDate
	}

	response.Created(c, resp)
}

// CreateBoletoPayment handles POST /api/v1/payments/boleto
func (h *PaymentHandler) CreateBoletoPayment(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	asaasReq := &asaas.CreatePaymentRequest{
		Customer:    req.CustomerID,
		Value:       req.Value,
		Description: req.Description,
		DueDate:     req.DueDate,
		ExternalRef: req.ExternalRef,
	}

	payment, err := h.usecase.CreateBoletoPayment(ctx, asaasReq)
	if err != nil {
		response.InternalError(c, "Failed to create Boleto payment: "+err.Error())
		return
	}

	response.Created(c, gin.H{
		"id":           payment.ID,
		"status":       payment.Status,
		"value":        payment.Value,
		"due_date":     payment.DueDate,
		"description":  payment.Description,
		"boleto_url":   payment.BankSlipURL,
		"invoice_url":  payment.InvoiceURL,
	})
}

// CreateCardPayment handles POST /api/v1/payments/card
func (h *PaymentHandler) CreateCardPayment(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateCardPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	asaasReq := &asaas.CreateCardPaymentRequest{
		Customer:     req.CustomerID,
		Value:        req.Value,
		Description:  req.Description,
		ExternalRef:  req.ExternalRef,
		Installments: req.Installments,
		CreditCard: &asaas.CreditCard{
			HolderName:  req.HolderName,
			Number:      req.CardNumber,
			ExpiryMonth: req.ExpiryMonth,
			ExpiryYear:  req.ExpiryYear,
			Ccv:         req.CCV,
		},
		CreditCardHolder: &asaas.CreditCardHolder{
			Name:       req.HolderName,
			Email:      req.HolderEmail,
			CPFCnpj:    req.HolderCPF,
			PostalCode: req.HolderPostalCode,
			Phone:      req.HolderPhone,
		},
	}

	payment, err := h.usecase.CreateCardPayment(ctx, asaasReq)
	if err != nil {
		response.InternalError(c, "Failed to create card payment: "+err.Error())
		return
	}

	response.Created(c, gin.H{
		"id":                      payment.ID,
		"status":                  payment.Status,
		"value":                   payment.Value,
		"description":             payment.Description,
		"transaction_receipt_url": payment.TransactionReceiptURL,
	})
}

// GetPaymentStatus handles GET /api/v1/payments/:id/status
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	ctx := c.Request.Context()
	paymentID := c.Param("id")

	payment, err := h.usecase.GetPaymentStatus(ctx, paymentID)
	if err != nil {
		response.InternalError(c, "Failed to get payment status: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":           payment.ID,
		"status":       payment.Status,
		"value":        payment.Value,
		"net_value":    payment.NetValue,
		"billing_type": payment.BillingType,
		"due_date":     payment.DueDate,
		"confirmed_date": payment.ConfirmedDate,
	})
}

// SimulateRevenueSplit handles GET /api/v1/payments/simulate-split
func (h *PaymentHandler) SimulateRevenueSplit(c *gin.Context) {
	var req entity.CalculateSplitRequest

	// Get from query params
	if value := c.Query("value"); value != "" {
		var v float64
		if _, err := parseFloat(value, &v); err != nil {
			response.BadRequest(c, "Invalid value parameter")
			return
		}
		req.GrossAmount = v
	}

	if method := c.Query("method"); method != "" {
		req.PaymentMethod = method
	}

	if req.GrossAmount <= 0 {
		response.BadRequest(c, "value parameter is required and must be greater than 0")
		return
	}

	if req.PaymentMethod == "" {
		req.PaymentMethod = "pix"
	}

	result := h.usecase.SimulateRevenueSplit(&req)
	response.Success(c, result)
}

// parseFloat parses a string to float64
func parseFloat(s string, v *float64) (float64, error) {
	parsed, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	*v = parsed
	return parsed, nil
}

// PaymentListFilters represents filters for listing payments
type PaymentListFilters struct {
	CustomerID string
	Status     string
	DateFrom   *time.Time
	DateTo     *time.Time
	Page       int
	PerPage    int
}

// PaymentListResponse represents the response for listing payments
type PaymentListResponse struct {
	Payments []entity.PaymentRecord `json:"payments"`
	Total    int                    `json:"total"`
	Page     int                    `json:"page"`
	PerPage  int                    `json:"per_page"`
}

// ListPayments handles GET /api/v1/payments
func (h *PaymentHandler) ListPayments(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse filters from query params
	filters := PaymentListFilters{
		CustomerID: c.Query("customer_id"),
		Status:     c.Query("status"),
		Page:       1,
		PerPage:    20,
	}

	// Parse date range
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filters.DateFrom = &t
		}
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		if t, err := time.Parse("2006-01-02", dateTo); err == nil {
			filters.DateTo = &t
		}
	}

	// Parse pagination
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			filters.Page = parsed
		}
	}
	if pp := c.Query("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 && parsed <= 100 {
			filters.PerPage = parsed
		}
	}

	// Get payments from enrollments (payments are tracked via enrollments)
	enrollments, total, err := h.matriculaRepo.FindAll(ctx, filters.Page, filters.PerPage)
	if err != nil {
		response.InternalError(c, "Failed to fetch payments: "+err.Error())
		return
	}

	// Convert enrollments to payment records with filtering
	var payments []entity.PaymentRecord
	for _, e := range enrollments {
		// Apply filters
		if filters.CustomerID != "" && e.AsaasCustomerID != nil && *e.AsaasCustomerID != filters.CustomerID {
			continue
		}
		if filters.Status != "" && e.PaymentStatus != filters.Status {
			continue
		}
		if filters.DateFrom != nil && e.EnrollmentDate.Before(*filters.DateFrom) {
			continue
		}
		if filters.DateTo != nil && e.EnrollmentDate.After(*filters.DateTo) {
			continue
		}

		payment := entity.PaymentRecord{
			ID:              e.ID,
			EnrollmentID:    e.ID,
			CustomerID:      e.AsaasCustomerID,
			AsaasPaymentID:  e.AsaasPaymentID,
			Amount:          e.FinalAmount,
			Status:          e.PaymentStatus,
			PaymentMethod:   e.PaymentMethod,
			StudentName:     e.StudentName,
			StudentEmail:    e.StudentEmail,
			CourseName:      e.CourseName,
			PaymentDate:     &e.EnrollmentDate,
			CreatedAt:       e.CreatedAt,
		}
		payments = append(payments, payment)
	}

	resp := PaymentListResponse{
		Payments: payments,
		Total:    total,
		Page:     filters.Page,
		PerPage:  filters.PerPage,
	}

	response.Success(c, resp)
}

// GetPaymentsByEnrollment handles GET /api/v1/payments/enrollment/:id
func (h *PaymentHandler) GetPaymentsByEnrollment(c *gin.Context) {
	ctx := c.Request.Context()
	enrollmentID := c.Param("id")

	if enrollmentID == "" {
		response.BadRequest(c, "Enrollment ID is required")
		return
	}

	enrollment, err := h.matriculaRepo.FindByID(ctx, enrollmentID)
	if err != nil {
		response.InternalError(c, "Failed to fetch enrollment: "+err.Error())
		return
	}

	if enrollment == nil {
		response.NotFound(c, "Enrollment not found")
		return
	}

	// Build payment record from enrollment
	payment := entity.PaymentRecord{
		ID:              enrollment.ID,
		EnrollmentID:    enrollment.ID,
		CustomerID:      enrollment.AsaasCustomerID,
		AsaasPaymentID:  enrollment.AsaasPaymentID,
		Amount:          enrollment.FinalAmount,
		OriginalAmount:  enrollment.Amount,
		DiscountAmount:  enrollment.DiscountAmount,
		Status:          enrollment.PaymentStatus,
		PaymentMethod:   enrollment.PaymentMethod,
		StudentName:     enrollment.StudentName,
		StudentEmail:    enrollment.StudentEmail,
		CourseName:      enrollment.CourseName,
		PaymentDate:     &enrollment.EnrollmentDate,
		CreatedAt:       enrollment.CreatedAt,
	}

	// If we have an Asaas payment ID, try to get additional details
	if enrollment.AsaasPaymentID != nil && *enrollment.AsaasPaymentID != "" {
		asaasPayment, err := h.usecase.GetPaymentStatus(ctx, *enrollment.AsaasPaymentID)
		if err == nil && asaasPayment != nil {
			payment.AsaasStatus = &asaasPayment.Status
			payment.NetValue = &asaasPayment.NetValue
			if asaasPayment.ConfirmedDate != "" {
				payment.ConfirmedDate = &asaasPayment.ConfirmedDate
			}
		}
	}

	response.Success(c, payment)
}
