package handler

import (
	"strconv"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
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

// CreateCustomer handles POST /api/v1/payments/customer
func (h *PaymentHandler) CreateCustomer(c *gin.Context) {
	ctx := c.Request.Context()

	var req payment.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	customer, err := h.usecase.CreateCustomer(ctx, &req)
	if err != nil {
		response.SafeInternalError(c, "Failed to create customer", err)
		return
	}

	response.Created(c, customer)
}

// CreatePixPayment handles POST /api/v1/payments/pix
func (h *PaymentHandler) CreatePixPayment(c *gin.Context) {
	ctx := c.Request.Context()

	var req payment.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	gwPayment, err := h.usecase.CreatePixPayment(ctx, &req)
	if err != nil {
		response.SafeInternalError(c, "Failed to create PIX payment", err)
		return
	}

	resp := gin.H{
		"id":      gwPayment.GatewayPaymentID,
		"status":  gwPayment.Status,
		"value":   gwPayment.Amount,
		"due_date": gwPayment.DueDate,
	}

	if gwPayment.PixQRCodeBase64 != "" {
		resp["pix_qr_code"] = gwPayment.PixQRCodeBase64
		resp["pix_copy_paste"] = gwPayment.PixCopyPaste
		resp["pix_expiration"] = gwPayment.PixExpiration
	}

	response.Created(c, resp)
}

// CreateBoletoPayment handles POST /api/v1/payments/boleto
func (h *PaymentHandler) CreateBoletoPayment(c *gin.Context) {
	ctx := c.Request.Context()

	var req payment.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	gwPayment, err := h.usecase.CreateBoletoPayment(ctx, &req)
	if err != nil {
		response.SafeInternalError(c, "Failed to create Boleto payment", err)
		return
	}

	response.Created(c, gin.H{
		"id":          gwPayment.GatewayPaymentID,
		"status":      gwPayment.Status,
		"value":       gwPayment.Amount,
		"due_date":    gwPayment.DueDate,
		"boleto_url":  gwPayment.BoletoURL,
		"invoice_url": gwPayment.InvoiceURL,
		"bar_code":    gwPayment.BoletoBarCode,
	})
}

// CreateCardPayment handles POST /api/v1/payments/card
func (h *PaymentHandler) CreateCardPayment(c *gin.Context) {
	ctx := c.Request.Context()

	var req payment.CreateCardPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	gwPayment, err := h.usecase.CreateCardPayment(ctx, &req)
	if err != nil {
		response.SafeInternalError(c, "Failed to create card payment", err)
		return
	}

	response.Created(c, gin.H{
		"id":                      gwPayment.GatewayPaymentID,
		"status":                  gwPayment.Status,
		"value":                   gwPayment.Amount,
		"transaction_receipt_url": gwPayment.TransactionReceiptURL,
	})
}

// GetPaymentStatus handles GET /api/v1/payments/:id/status
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	ctx := c.Request.Context()
	paymentID := c.Param("id")

	result, err := h.usecase.GetPaymentStatus(ctx, paymentID)
	if err != nil {
		response.SafeInternalError(c, "Failed to get payment status", err)
		return
	}

	response.Success(c, result)
}

// SimulateRevenueSplit handles GET /api/v1/payments/simulate-split
func (h *PaymentHandler) SimulateRevenueSplit(c *gin.Context) {
	var req entity.CalculateSplitRequest

	if value := c.Query("value"); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			req.GrossAmount = parsed
		}
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

// ListPayments handles GET /api/v1/payments
func (h *PaymentHandler) ListPayments(c *gin.Context) {
	ctx := c.Request.Context()

	filters := repository.PaymentFilters{
		EnrollmentID:  c.Query("enrollment_id"),
		Gateway:       c.Query("gateway"),
		Status:        c.Query("status"),
		PaymentMethod: c.Query("payment_method"),
		Page:          1,
		PerPage:       20,
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

	payments, total, err := h.usecase.ListPayments(ctx, filters)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch payments", err)
		return
	}

	response.Success(c, gin.H{
		"payments": payments,
		"total":    total,
		"page":     filters.Page,
		"per_page": filters.PerPage,
	})
}

// GetPaymentsByEnrollment handles GET /api/v1/payments/enrollment/:id
func (h *PaymentHandler) GetPaymentsByEnrollment(c *gin.Context) {
	ctx := c.Request.Context()
	enrollmentID := c.Param("id")

	if enrollmentID == "" {
		response.BadRequest(c, "Enrollment ID is required")
		return
	}

	payments, err := h.usecase.GetPaymentsByEnrollment(ctx, enrollmentID)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch payments", err)
		return
	}

	if len(payments) == 0 {
		response.NotFound(c, "No payments found for this enrollment")
		return
	}

	response.Success(c, payments)
}
