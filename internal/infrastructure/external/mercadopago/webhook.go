package mercadopago

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/condotrack/api/internal/domain/gateway"
)

// ValidateWebhookSignature validates a Mercado Pago webhook signature using HMAC-SHA256.
// MP sends the signature in the "x-signature" header with format:
// "ts=<timestamp>,v1=<hmac_hash>"
// The data to sign is: "id:<data.id>;request-id:<x-request-id>;ts:<timestamp>;"
func ValidateWebhookSignature(secret string, headers map[string]string, body []byte) bool {
	signature := headers["x-signature"]
	requestID := headers["x-request-id"]
	if signature == "" {
		return false
	}

	// Parse x-signature header
	parts := strings.Split(signature, ",")
	var ts, v1 string
	for _, part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "ts":
			ts = kv[1]
		case "v1":
			v1 = kv[1]
		}
	}

	if ts == "" || v1 == "" {
		return false
	}

	// Extract data.id from body
	var notification MPWebhookNotification
	if err := json.Unmarshal(body, &notification); err != nil {
		return false
	}

	// Build the manifest string to validate
	// Format: "id:<data.id>;request-id:<x-request-id>;ts:<ts>;"
	manifest := fmt.Sprintf("id:%s;request-id:%s;ts:%s;", notification.Data.ID, requestID, ts)

	// Compute HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(manifest))
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(v1))
}

// ParseWebhookEvent parses a Mercado Pago webhook notification and fetches
// the full payment data from the API to build a canonical WebhookEvent.
func ParseWebhookEvent(ctx context.Context, adapter *MercadoPagoAdapter, headers map[string]string, body []byte) (*gateway.WebhookEvent, error) {
	var notification MPWebhookNotification
	if err := json.Unmarshal(body, &notification); err != nil {
		return nil, fmt.Errorf("failed to parse MP webhook notification: %w", err)
	}

	// Only handle payment notifications
	if notification.Type != "payment" {
		return nil, fmt.Errorf("unsupported webhook type: %s", notification.Type)
	}

	paymentID := notification.Data.ID
	if paymentID == "" {
		return nil, fmt.Errorf("webhook notification has no payment ID")
	}

	// Fetch the full payment details from MP API
	payment, err := adapter.GetClient().GetPayment(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch MP payment %s: %w", paymentID, err)
	}

	canonical := &gateway.WebhookEvent{
		GatewayEvent:     notification.Action,
		GatewayName:      "mercadopago",
		PaymentID:        paymentID,
		Amount:           payment.TransactionAmount,
		NetAmount:        payment.NetReceivedAmount,
		GatewayRawStatus: payment.Status,
		Status:           adapter.NormalizeStatus(payment.Status),
		BillingType:      adapter.normalizeBillingType(payment.PaymentTypeID),
		ExternalRef:      payment.ExternalReference,
		RawPayload:       body,
	}

	// Set event type based on action + payment status
	canonical.EventType = adapter.normalizeEventType(notification.Action, payment.Status)

	// Parse approved date
	if payment.DateApproved != nil && *payment.DateApproved != "" {
		if t, err := ParseMPDateTime(*payment.DateApproved); err == nil {
			canonical.PaidAt = &t
		}
	}

	// Extract payer/customer info
	if payment.Payer != nil {
		canonical.CustomerID = payment.Payer.Email
	}

	return canonical, nil
}
