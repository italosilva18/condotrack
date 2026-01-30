package handler

import (
	"context"
	"log"

	"github.com/condotrack/api/internal/config"
	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/condotrack/api/internal/infrastructure/database"
	"github.com/condotrack/api/internal/infrastructure/external/asaas"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// WebhookHandler handles webhook requests from Asaas
type WebhookHandler struct {
	cfg           *config.Config
	db            *database.MySQL
	matriculaRepo repository.MatriculaRepository
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(cfg *config.Config, db *database.MySQL, matriculaRepo repository.MatriculaRepository) *WebhookHandler {
	return &WebhookHandler{
		cfg:           cfg,
		db:            db,
		matriculaRepo: matriculaRepo,
	}
}

// HandleAsaasWebhook handles POST /api/v1/webhooks/asaas
func (h *WebhookHandler) HandleAsaasWebhook(c *gin.Context) {
	ctx := c.Request.Context()

	// Verify webhook token if configured
	if h.cfg.AsaasWebhookToken != "" {
		token := c.GetHeader("asaas-access-token")
		if token != h.cfg.AsaasWebhookToken {
			log.Printf("Invalid webhook token received: %s", token)
			response.Unauthorized(c, "Invalid webhook token")
			return
		}
	}

	var event asaas.WebhookEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		log.Printf("Failed to parse webhook payload: %v", err)
		response.BadRequest(c, "Invalid payload: "+err.Error())
		return
	}

	if event.Payment == nil {
		log.Printf("Webhook event without payment data: %s", event.Event)
		c.JSON(200, gin.H{"success": true, "message": "No payment data"})
		return
	}

	log.Printf("Received Asaas webhook: event=%s, payment_id=%s", event.Event, event.Payment.ID)

	// Handle different event types
	switch event.Event {
	case asaas.WebhookEventPaymentConfirmed, asaas.WebhookEventPaymentReceived:
		if err := h.handlePaymentConfirmed(ctx, event.Payment); err != nil {
			log.Printf("Failed to handle payment confirmation: %v", err)
			response.InternalError(c, "Failed to process webhook")
			return
		}

	case asaas.WebhookEventPaymentOverdue:
		if err := h.handlePaymentOverdue(ctx, event.Payment); err != nil {
			log.Printf("Failed to handle payment overdue: %v", err)
			response.InternalError(c, "Failed to process webhook")
			return
		}

	case asaas.WebhookEventPaymentRefunded:
		if err := h.handlePaymentRefunded(ctx, event.Payment); err != nil {
			log.Printf("Failed to handle payment refund: %v", err)
			response.InternalError(c, "Failed to process webhook")
			return
		}

	case asaas.WebhookEventPaymentDeleted:
		if err := h.handlePaymentDeleted(ctx, event.Payment); err != nil {
			log.Printf("Failed to handle payment deletion: %v", err)
			response.InternalError(c, "Failed to process webhook")
			return
		}

	default:
		log.Printf("Unhandled webhook event: %s", event.Event)
	}

	// Return success to Asaas
	c.JSON(200, gin.H{
		"success": true,
		"message": "Webhook processed",
	})
}

// handlePaymentConfirmed handles payment confirmation events
func (h *WebhookHandler) handlePaymentConfirmed(ctx context.Context, payment *asaas.WebhookPayment) error {
	// Find enrollment by Asaas payment ID
	enrollment, err := h.matriculaRepo.FindByAsaasPaymentID(ctx, payment.ID)
	if err != nil {
		return err
	}

	if enrollment == nil {
		log.Printf("No enrollment found for payment ID: %s", payment.ID)
		return nil
	}

	// Start transaction
	tx, err := h.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update payment status to confirmed
	if err := h.matriculaRepo.UpdatePaymentStatusWithTx(ctx, tx, enrollment.ID, entity.PaymentStatusConfirmed); err != nil {
		return err
	}

	// Activate enrollment
	enrollment.Status = entity.EnrollmentStatusActive
	enrollment.PaymentStatus = entity.PaymentStatusConfirmed
	if err := h.matriculaRepo.UpdateWithTx(ctx, tx, enrollment); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("Payment confirmed for enrollment: %s", enrollment.ID)
	return nil
}

// handlePaymentOverdue handles payment overdue events
func (h *WebhookHandler) handlePaymentOverdue(ctx context.Context, payment *asaas.WebhookPayment) error {
	enrollment, err := h.matriculaRepo.FindByAsaasPaymentID(ctx, payment.ID)
	if err != nil {
		return err
	}

	if enrollment == nil {
		log.Printf("No enrollment found for payment ID: %s", payment.ID)
		return nil
	}

	// Update enrollment status
	if err := h.matriculaRepo.UpdatePaymentStatus(ctx, enrollment.ID, "overdue"); err != nil {
		return err
	}

	log.Printf("Payment overdue for enrollment: %s", enrollment.ID)
	return nil
}

// handlePaymentRefunded handles payment refund events
func (h *WebhookHandler) handlePaymentRefunded(ctx context.Context, payment *asaas.WebhookPayment) error {
	enrollment, err := h.matriculaRepo.FindByAsaasPaymentID(ctx, payment.ID)
	if err != nil {
		return err
	}

	if enrollment == nil {
		log.Printf("No enrollment found for payment ID: %s", payment.ID)
		return nil
	}

	// Start transaction
	tx, err := h.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update payment status to refunded
	if err := h.matriculaRepo.UpdatePaymentStatusWithTx(ctx, tx, enrollment.ID, entity.PaymentStatusRefunded); err != nil {
		return err
	}

	// Cancel enrollment
	enrollment.Status = entity.EnrollmentStatusCancelled
	enrollment.PaymentStatus = entity.PaymentStatusRefunded
	if err := h.matriculaRepo.UpdateWithTx(ctx, tx, enrollment); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("Payment refunded for enrollment: %s", enrollment.ID)
	return nil
}

// handlePaymentDeleted handles payment deletion events
func (h *WebhookHandler) handlePaymentDeleted(ctx context.Context, payment *asaas.WebhookPayment) error {
	enrollment, err := h.matriculaRepo.FindByAsaasPaymentID(ctx, payment.ID)
	if err != nil {
		return err
	}

	if enrollment == nil {
		log.Printf("No enrollment found for payment ID: %s", payment.ID)
		return nil
	}

	// Update enrollment status to cancelled
	if err := h.matriculaRepo.UpdateStatus(ctx, enrollment.ID, entity.EnrollmentStatusCancelled); err != nil {
		return err
	}

	log.Printf("Payment deleted for enrollment: %s", enrollment.ID)
	return nil
}
