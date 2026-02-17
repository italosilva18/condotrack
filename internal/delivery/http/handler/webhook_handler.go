package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"

	"github.com/condotrack/api/internal/config"
	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/gateway"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/condotrack/api/internal/infrastructure/database"
	"github.com/condotrack/api/internal/infrastructure/external"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WebhookHandler handles webhook requests from payment gateways.
type WebhookHandler struct {
	cfg              *config.Config
	db               *database.MySQL
	matriculaRepo    repository.MatriculaRepository
	paymentRepo      repository.PaymentRepository
	paymentTxnRepo   repository.PaymentTransactionRepository
	revenueSplitRepo repository.RevenueSplitRepository
	gatewayFactory   *external.GatewayFactory
}

// NewWebhookHandler creates a new webhook handler.
func NewWebhookHandler(
	cfg *config.Config,
	db *database.MySQL,
	matriculaRepo repository.MatriculaRepository,
	paymentRepo repository.PaymentRepository,
	paymentTxnRepo repository.PaymentTransactionRepository,
	revenueSplitRepo repository.RevenueSplitRepository,
	gatewayFactory *external.GatewayFactory,
) *WebhookHandler {
	return &WebhookHandler{
		cfg:              cfg,
		db:               db,
		matriculaRepo:    matriculaRepo,
		paymentRepo:      paymentRepo,
		paymentTxnRepo:   paymentTxnRepo,
		revenueSplitRepo: revenueSplitRepo,
		gatewayFactory:   gatewayFactory,
	}
}

// HandleAsaasWebhook handles POST /api/v1/webhooks/asaas
func (h *WebhookHandler) HandleAsaasWebhook(c *gin.Context) {
	ctx := c.Request.Context()

	// Read raw body for parsing (limited to 1MB to prevent DoS)
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	if err != nil {
		log.Printf("Failed to read webhook body: %v", err)
		response.BadRequest(c, "Failed to read body")
		return
	}

	// Get the Asaas gateway adapter
	gw, err := h.gatewayFactory.Get("asaas")
	if err != nil {
		log.Printf("Asaas gateway not registered: %v", err)
		response.InternalError(c, "Gateway not configured")
		return
	}

	// Collect headers for validation
	headers := map[string]string{
		"asaas-access-token": c.GetHeader("asaas-access-token"),
	}

	// Validate webhook signature
	if !gw.ValidateWebhookSignature(ctx, headers, body) {
		log.Printf("Invalid webhook signature")
		response.Unauthorized(c, "Invalid webhook token")
		return
	}

	// Parse webhook event into canonical format
	event, err := gw.ParseWebhookEvent(ctx, headers, body)
	if err != nil {
		log.Printf("Failed to parse webhook: %v", err)
		response.BadRequest(c, "Invalid payload: "+err.Error())
		return
	}

	log.Printf("Received webhook: gateway=%s event=%s payment_id=%s status=%s",
		event.GatewayName, event.EventType, event.PaymentID, event.Status)

	// Handle by canonical event type
	switch event.EventType {
	case gateway.EventPaymentConfirmed:
		if err := h.handlePaymentConfirmed(ctx, event); err != nil {
			log.Printf("Failed to handle payment confirmation: %v", err)
			response.InternalError(c, "Failed to process webhook")
			return
		}

	case gateway.EventPaymentOverdue:
		if err := h.handlePaymentOverdue(ctx, event); err != nil {
			log.Printf("Failed to handle payment overdue: %v", err)
			response.InternalError(c, "Failed to process webhook")
			return
		}

	case gateway.EventPaymentRefunded:
		if err := h.handlePaymentRefunded(ctx, event); err != nil {
			log.Printf("Failed to handle payment refund: %v", err)
			response.InternalError(c, "Failed to process webhook")
			return
		}

	case gateway.EventPaymentDeleted:
		if err := h.handlePaymentDeleted(ctx, event); err != nil {
			log.Printf("Failed to handle payment deletion: %v", err)
			response.InternalError(c, "Failed to process webhook")
			return
		}

	default:
		log.Printf("Unhandled webhook event: %s", event.EventType)
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Webhook processed",
	})
}

// handlePaymentConfirmed processes payment confirmation and creates revenue split.
func (h *WebhookHandler) handlePaymentConfirmed(ctx context.Context, event *gateway.WebhookEvent) error {
	// 1. Try to find payment in payments table by gateway payment ID
	payment, err := h.paymentRepo.FindByGatewayPaymentID(ctx, event.GatewayName, event.PaymentID)
	if err != nil {
		return err
	}

	// 2. Also find enrollment (for backwards compatibility and revenue split)
	enrollment, err := h.matriculaRepo.FindByAsaasPaymentID(ctx, event.PaymentID)
	if err != nil {
		return err
	}

	if enrollment == nil && payment == nil {
		log.Printf("No enrollment or payment found for gateway payment ID: %s", event.PaymentID)
		return nil
	}

	// 3. Log webhook receipt
	rawPayload := string(event.RawPayload)
	h.logPaymentTransaction(ctx, payment, entity.TxEventWebhookReceived, event.GatewayEvent,
		nil, &event.Status, &event.Amount, &rawPayload)

	// 4. Start transaction
	tx, err := h.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 5. Update payment record if it exists
	if payment != nil {
		prevStatus := payment.Status
		payment.Status = entity.FinPaymentStatusConfirmed
		payment.PaidAt = event.PaidAt
		netAmount := event.NetAmount
		if netAmount > 0 {
			payment.NetAmount = netAmount
		}
		if err := h.paymentRepo.UpdateWithTx(ctx, tx, payment); err != nil {
			return err
		}

		// Log status change
		h.logPaymentTransaction(ctx, payment, entity.TxEventStatusChanged, event.GatewayEvent,
			&prevStatus, &payment.Status, &event.Amount, nil)
	}

	// 6. Update enrollment status
	if enrollment != nil {
		if err := h.matriculaRepo.UpdatePaymentStatusWithTx(ctx, tx, enrollment.ID, entity.PaymentStatusConfirmed); err != nil {
			return err
		}

		enrollment.Status = entity.EnrollmentStatusActive
		enrollment.PaymentStatus = entity.PaymentStatusConfirmed
		if err := h.matriculaRepo.UpdateWithTx(ctx, tx, enrollment); err != nil {
			return err
		}

		// 7. CREATE REVENUE SPLIT (CRITICAL FIX)
		gw := h.gatewayFactory.GetActive()
		fees := gw.GetFees()
		billingType := event.BillingType
		if payment != nil {
			billingType = payment.PaymentMethod
		}

		grossAmount := event.Amount
		if payment != nil {
			grossAmount = payment.GrossAmount
		}

		gatewayFee := roundCents(calculateGatewayFee(grossAmount, billingType, fees))
		netAmount := roundCents(grossAmount - gatewayFee)
		instructorAmount := roundCents(netAmount * (h.cfg.RevenueInstructorPercent / 100))
		platformAmount := roundCents(netAmount * (h.cfg.RevenuePlatformPercent / 100))

		paymentID := event.PaymentID
		if payment != nil {
			paymentID = payment.ID
		}

		split := &entity.RevenueSplit{
			ID:               uuid.New().String(),
			EnrollmentID:     enrollment.ID,
			PaymentID:        paymentID,
			GrossAmount:      grossAmount,
			NetAmount:        netAmount,
			PlatformFee:      platformAmount,
			PaymentFee:       gatewayFee,
			InstructorAmount: instructorAmount,
			PlatformAmount:   platformAmount,
			InstructorID:     enrollment.InstructorID,
			PaymentMethod:    billingType,
			Status:           entity.RevenueSplitStatusPending,
		}

		if err := h.revenueSplitRepo.CreateWithTx(ctx, tx, split); err != nil {
			log.Printf("Failed to create revenue split: %v", err)
			return fmt.Errorf("failed to create revenue split: %w", err)
		}
	}

	// 8. Commit
	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("Payment confirmed: gateway_id=%s enrollment=%s", event.PaymentID, getEnrollmentID(enrollment))
	return nil
}

// handlePaymentOverdue processes payment overdue events.
func (h *WebhookHandler) handlePaymentOverdue(ctx context.Context, event *gateway.WebhookEvent) error {
	// Update payment record
	payment, err := h.paymentRepo.FindByGatewayPaymentID(ctx, event.GatewayName, event.PaymentID)
	if err != nil {
		log.Printf("Failed to find payment for overdue event: %v", err)
	}
	if payment != nil {
		if err := h.paymentRepo.UpdateStatus(ctx, payment.ID, entity.FinPaymentStatusOverdue); err != nil {
			log.Printf("Failed to update payment status to overdue: %v", err)
		}
	}

	// Update enrollment
	enrollment, err := h.matriculaRepo.FindByAsaasPaymentID(ctx, event.PaymentID)
	if err != nil {
		return err
	}
	if enrollment == nil {
		log.Printf("No enrollment found for payment ID: %s", event.PaymentID)
		return nil
	}

	if err := h.matriculaRepo.UpdatePaymentStatus(ctx, enrollment.ID, entity.PaymentStatusOverdue); err != nil {
		return err
	}

	log.Printf("Payment overdue: enrollment=%s", enrollment.ID)
	return nil
}

// handlePaymentRefunded processes payment refund events.
func (h *WebhookHandler) handlePaymentRefunded(ctx context.Context, event *gateway.WebhookEvent) error {
	// Update payment record
	payment, err := h.paymentRepo.FindByGatewayPaymentID(ctx, event.GatewayName, event.PaymentID)
	if err != nil {
		log.Printf("Failed to find payment for refund event: %v", err)
	}
	if payment != nil {
		payment.Status = entity.FinPaymentStatusRefunded
		payment.RefundedAmount = event.Amount
		if err := h.paymentRepo.Update(ctx, payment); err != nil {
			log.Printf("Failed to update payment for refund: %v", err)
		}
	}

	enrollment, err := h.matriculaRepo.FindByAsaasPaymentID(ctx, event.PaymentID)
	if err != nil {
		return err
	}
	if enrollment == nil {
		log.Printf("No enrollment found for payment ID: %s", event.PaymentID)
		return nil
	}

	tx, err := h.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := h.matriculaRepo.UpdatePaymentStatusWithTx(ctx, tx, enrollment.ID, entity.PaymentStatusRefunded); err != nil {
		return err
	}

	enrollment.Status = entity.EnrollmentStatusCancelled
	enrollment.PaymentStatus = entity.PaymentStatusRefunded
	if err := h.matriculaRepo.UpdateWithTx(ctx, tx, enrollment); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("Payment refunded: enrollment=%s", enrollment.ID)
	return nil
}

// handlePaymentDeleted processes payment deletion events.
func (h *WebhookHandler) handlePaymentDeleted(ctx context.Context, event *gateway.WebhookEvent) error {
	// Update payment record
	payment, err := h.paymentRepo.FindByGatewayPaymentID(ctx, event.GatewayName, event.PaymentID)
	if err != nil {
		log.Printf("Failed to find payment for deleted event: %v", err)
	}
	if payment != nil {
		if err := h.paymentRepo.UpdateStatus(ctx, payment.ID, entity.FinPaymentStatusCancelled); err != nil {
			log.Printf("Failed to update payment status to cancelled: %v", err)
		}
	}

	enrollment, err := h.matriculaRepo.FindByAsaasPaymentID(ctx, event.PaymentID)
	if err != nil {
		return err
	}
	if enrollment == nil {
		log.Printf("No enrollment found for payment ID: %s", event.PaymentID)
		return nil
	}

	if err := h.matriculaRepo.UpdateStatus(ctx, enrollment.ID, entity.EnrollmentStatusCancelled); err != nil {
		return err
	}

	log.Printf("Payment deleted: enrollment=%s", enrollment.ID)
	return nil
}

// logPaymentTransaction creates a payment transaction log entry.
func (h *WebhookHandler) logPaymentTransaction(
	ctx context.Context, payment *entity.Payment,
	eventType, gatewayEvent string,
	prevStatus, newStatus *string, amount *float64, rawPayload *string,
) {
	if payment == nil {
		return
	}

	txLog := &entity.PaymentTransaction{
		ID:          uuid.New().String(),
		PaymentID:   payment.ID,
		EventSource: entity.EventSourceWebhook,
		EventType:   eventType,
		NewStatus:   derefStr(newStatus),
	}
	if prevStatus != nil {
		txLog.PreviousStatus = prevStatus
	}
	if amount != nil {
		txLog.Amount = amount
	}
	if rawPayload != nil {
		txLog.RawPayload = rawPayload
	}
	ge := gatewayEvent
	txLog.GatewayEventID = &ge

	if err := h.paymentTxnRepo.Create(ctx, txLog); err != nil {
		log.Printf("Failed to log payment transaction: %v", err)
	}
}

// calculateGatewayFee calculates the fee for a given payment method.
func calculateGatewayFee(amount float64, billingType string, fees gateway.GatewayFees) float64 {
	switch billingType {
	case "pix", "PIX":
		return amount * fees.PixPercent
	case "boleto", "BOLETO":
		return fees.BoletoFixed
	case "credit_card", "CREDIT_CARD", "card":
		return (amount * fees.CardPercent) + fees.CardFixed
	default:
		return 0
	}
}

// roundCents rounds a float64 to 2 decimal places for financial precision.
func roundCents(v float64) float64 {
	return math.Round(v*100) / 100
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func getEnrollmentID(enrollment *entity.Matricula) string {
	if enrollment == nil {
		return "unknown"
	}
	return enrollment.ID
}

// HandleMercadoPagoWebhook handles POST /api/v1/webhooks/mercadopago
func (h *WebhookHandler) HandleMercadoPagoWebhook(c *gin.Context) {
	ctx := c.Request.Context()

	// Read raw body (limited to 1MB to prevent DoS)
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	if err != nil {
		log.Printf("Failed to read MP webhook body: %v", err)
		response.BadRequest(c, "Failed to read body")
		return
	}

	// Get the Mercado Pago gateway adapter
	gw, err := h.gatewayFactory.Get("mercadopago")
	if err != nil {
		log.Printf("Mercado Pago gateway not registered: %v", err)
		response.InternalError(c, "Gateway not configured")
		return
	}

	// Collect headers for validation
	headers := map[string]string{
		"x-signature":  c.GetHeader("x-signature"),
		"x-request-id": c.GetHeader("x-request-id"),
	}

	// Validate webhook signature
	if !gw.ValidateWebhookSignature(ctx, headers, body) {
		log.Printf("Invalid MP webhook signature")
		response.Unauthorized(c, "Invalid webhook signature")
		return
	}

	// Parse webhook event into canonical format
	event, err := gw.ParseWebhookEvent(ctx, headers, body)
	if err != nil {
		log.Printf("Failed to parse MP webhook: %v", err)
		// Return 200 for unsupported event types (MP expects 200)
		c.JSON(200, gin.H{"success": true, "message": "Event type not handled"})
		return
	}

	log.Printf("Received MP webhook: event=%s payment_id=%s status=%s",
		event.EventType, event.PaymentID, event.Status)

	// Reuse the same canonical event handlers
	switch event.EventType {
	case gateway.EventPaymentConfirmed:
		if err := h.handlePaymentConfirmed(ctx, event); err != nil {
			log.Printf("Failed to handle MP payment confirmation: %v", err)
			response.InternalError(c, "Failed to process webhook")
			return
		}

	case gateway.EventPaymentOverdue:
		if err := h.handlePaymentOverdue(ctx, event); err != nil {
			log.Printf("Failed to handle MP payment overdue: %v", err)
			response.InternalError(c, "Failed to process webhook")
			return
		}

	case gateway.EventPaymentRefunded:
		if err := h.handlePaymentRefunded(ctx, event); err != nil {
			log.Printf("Failed to handle MP payment refund: %v", err)
			response.InternalError(c, "Failed to process webhook")
			return
		}

	case gateway.EventPaymentDeleted:
		if err := h.handlePaymentDeleted(ctx, event); err != nil {
			log.Printf("Failed to handle MP payment deletion: %v", err)
			response.InternalError(c, "Failed to process webhook")
			return
		}

	case gateway.EventPaymentChargeback:
		if err := h.handlePaymentChargeback(ctx, event); err != nil {
			log.Printf("Failed to handle MP chargeback: %v", err)
			response.InternalError(c, "Failed to process webhook")
			return
		}

	default:
		log.Printf("Unhandled MP webhook event: %s", event.EventType)
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Webhook processed",
	})
}

// handlePaymentChargeback processes chargeback events.
func (h *WebhookHandler) handlePaymentChargeback(ctx context.Context, event *gateway.WebhookEvent) error {
	// Update payment record
	payment, err := h.paymentRepo.FindByGatewayPaymentID(ctx, event.GatewayName, event.PaymentID)
	if err != nil {
		log.Printf("Failed to find payment for chargeback event: %v", err)
	}
	if payment != nil {
		payment.Status = entity.FinPaymentStatusChargeback
		if err := h.paymentRepo.Update(ctx, payment); err != nil {
			log.Printf("Failed to update payment for chargeback: %v", err)
		}
	}

	// Update enrollment
	enrollment, err := h.matriculaRepo.FindByAsaasPaymentID(ctx, event.PaymentID)
	if err != nil {
		return err
	}
	if enrollment == nil {
		log.Printf("No enrollment found for payment ID: %s", event.PaymentID)
		return nil
	}

	if err := h.matriculaRepo.UpdatePaymentStatus(ctx, enrollment.ID, entity.PaymentStatusChargeback); err != nil {
		return err
	}

	log.Printf("Payment chargeback: enrollment=%s", enrollment.ID)
	return nil
}
