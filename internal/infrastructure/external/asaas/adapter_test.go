package asaas

import (
	"testing"

	"github.com/condotrack/api/internal/domain/gateway"
)

func newTestAsaasAdapter() *AsaasAdapter {
	return NewAsaasAdapter(nil, gateway.GatewayFees{
		PixPercent:  0.0099,
		BoletoFixed: 2.99,
		CardPercent: 0.0499,
		CardFixed:   0.49,
	}, "test-webhook-token")
}

func TestAsaasAdapter_Name(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.Name() != "asaas" {
		t.Errorf("Name() = %q, want %q", a.Name(), "asaas")
	}
}

// --- NormalizeStatus tests ---

func TestNormalizeStatus_Pending(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.NormalizeStatus("PENDING") != gateway.StatusPending {
		t.Error("PENDING should map to pending")
	}
}

func TestNormalizeStatus_AwaitingRisk(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.NormalizeStatus("AWAITING_RISK_ANALYSIS") != gateway.StatusPending {
		t.Error("AWAITING_RISK_ANALYSIS should map to pending")
	}
}

func TestNormalizeStatus_Received(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.NormalizeStatus("RECEIVED") != gateway.StatusConfirmed {
		t.Error("RECEIVED should map to confirmed")
	}
}

func TestNormalizeStatus_Confirmed(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.NormalizeStatus("CONFIRMED") != gateway.StatusConfirmed {
		t.Error("CONFIRMED should map to confirmed")
	}
}

func TestNormalizeStatus_ReceivedInCash(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.NormalizeStatus("RECEIVED_IN_CASH") != gateway.StatusConfirmed {
		t.Error("RECEIVED_IN_CASH should map to confirmed")
	}
}

func TestNormalizeStatus_Overdue(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.NormalizeStatus("OVERDUE") != gateway.StatusOverdue {
		t.Error("OVERDUE should map to overdue")
	}
}

func TestNormalizeStatus_Refunded(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.NormalizeStatus("REFUNDED") != gateway.StatusRefunded {
		t.Error("REFUNDED should map to refunded")
	}
}

func TestNormalizeStatus_RefundRequested(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.NormalizeStatus("REFUND_REQUESTED") != gateway.StatusRefunded {
		t.Error("REFUND_REQUESTED should map to refunded")
	}
}

func TestNormalizeStatus_ChargebackRequested(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.NormalizeStatus("CHARGEBACK_REQUESTED") != gateway.StatusChargeback {
		t.Error("CHARGEBACK_REQUESTED should map to chargeback")
	}
}

func TestNormalizeStatus_ChargebackDispute(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.NormalizeStatus("CHARGEBACK_DISPUTE") != gateway.StatusChargeback {
		t.Error("CHARGEBACK_DISPUTE should map to chargeback")
	}
}

func TestNormalizeStatus_AwaitingChargeback(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.NormalizeStatus("AWAITING_CHARGEBACK_REVERSAL") != gateway.StatusChargeback {
		t.Error("AWAITING_CHARGEBACK_REVERSAL should map to chargeback")
	}
}

func TestNormalizeStatus_Unknown(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.NormalizeStatus("UNKNOWN_STATUS") != gateway.StatusFailed {
		t.Error("unknown status should map to failed")
	}
}

// --- normalizeBillingType tests ---

func TestNormalizeBillingType_PIX(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.normalizeBillingType("PIX") != gateway.BillingPIX {
		t.Error("PIX should map to pix")
	}
}

func TestNormalizeBillingType_Boleto(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.normalizeBillingType("BOLETO") != gateway.BillingBoleto {
		t.Error("BOLETO should map to boleto")
	}
}

func TestNormalizeBillingType_CreditCard(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.normalizeBillingType("CREDIT_CARD") != gateway.BillingCreditCard {
		t.Error("CREDIT_CARD should map to credit_card")
	}
}

func TestNormalizeBillingType_DebitCard(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.normalizeBillingType("DEBIT_CARD") != gateway.BillingDebitCard {
		t.Error("DEBIT_CARD should map to debit_card")
	}
}

func TestNormalizeBillingType_Unknown(t *testing.T) {
	a := newTestAsaasAdapter()
	result := a.normalizeBillingType("CRYPTO")
	if result != "CRYPTO" {
		t.Errorf("unknown billing type should pass through, got %q", result)
	}
}

// --- normalizeEventType tests ---

func TestNormalizeEventType_Created(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.normalizeEventType(WebhookEventPaymentCreated) != gateway.EventPaymentCreated {
		t.Error("PAYMENT_CREATED should map to payment_created")
	}
}

func TestNormalizeEventType_Confirmed(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.normalizeEventType(WebhookEventPaymentConfirmed) != gateway.EventPaymentConfirmed {
		t.Error("PAYMENT_CONFIRMED should map to payment_confirmed")
	}
}

func TestNormalizeEventType_Received(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.normalizeEventType(WebhookEventPaymentReceived) != gateway.EventPaymentConfirmed {
		t.Error("PAYMENT_RECEIVED should map to payment_confirmed")
	}
}

func TestNormalizeEventType_ReceivedInCash(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.normalizeEventType(WebhookEventPaymentReceivedInCash) != gateway.EventPaymentConfirmed {
		t.Error("PAYMENT_RECEIVED_IN_CASH should map to payment_confirmed")
	}
}

func TestNormalizeEventType_Overdue(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.normalizeEventType(WebhookEventPaymentOverdue) != gateway.EventPaymentOverdue {
		t.Error("PAYMENT_OVERDUE should map to payment_overdue")
	}
}

func TestNormalizeEventType_Refunded(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.normalizeEventType(WebhookEventPaymentRefunded) != gateway.EventPaymentRefunded {
		t.Error("PAYMENT_REFUNDED should map to payment_refunded")
	}
}

func TestNormalizeEventType_Deleted(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.normalizeEventType(WebhookEventPaymentDeleted) != gateway.EventPaymentDeleted {
		t.Error("PAYMENT_DELETED should map to payment_deleted")
	}
}

func TestNormalizeEventType_ChargebackRequested(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.normalizeEventType(WebhookEventPaymentChargebackRequested) != gateway.EventPaymentChargeback {
		t.Error("PAYMENT_CHARGEBACK_REQUESTED should map to payment_chargeback")
	}
}

func TestNormalizeEventType_ChargebackDispute(t *testing.T) {
	a := newTestAsaasAdapter()
	if a.normalizeEventType(WebhookEventPaymentChargebackDispute) != gateway.EventPaymentChargeback {
		t.Error("PAYMENT_CHARGEBACK_DISPUTE should map to payment_chargeback")
	}
}

func TestNormalizeEventType_Unknown(t *testing.T) {
	a := newTestAsaasAdapter()
	result := a.normalizeEventType("PAYMENT_UNKNOWN_EVENT")
	if result != "PAYMENT_UNKNOWN_EVENT" {
		t.Errorf("unknown event should pass through, got %q", result)
	}
}

// --- ValidateWebhookSignature tests ---

func TestValidateWebhookSignature_ValidToken(t *testing.T) {
	a := newTestAsaasAdapter()
	headers := map[string]string{"asaas-access-token": "test-webhook-token"}
	if !a.ValidateWebhookSignature(nil, headers, nil) {
		t.Error("valid token should pass validation")
	}
}

func TestValidateWebhookSignature_InvalidToken(t *testing.T) {
	a := newTestAsaasAdapter()
	headers := map[string]string{"asaas-access-token": "wrong-token"}
	if a.ValidateWebhookSignature(nil, headers, nil) {
		t.Error("invalid token should fail validation")
	}
}

func TestValidateWebhookSignature_MissingToken(t *testing.T) {
	a := newTestAsaasAdapter()
	headers := map[string]string{}
	if a.ValidateWebhookSignature(nil, headers, nil) {
		t.Error("missing token should fail validation")
	}
}

func TestValidateWebhookSignature_NoTokenConfigured(t *testing.T) {
	a := NewAsaasAdapter(nil, gateway.GatewayFees{}, "")
	headers := map[string]string{}
	if !a.ValidateWebhookSignature(nil, headers, nil) {
		t.Error("no token configured should skip validation (return true)")
	}
}

// --- GetFees tests ---

func TestGetFees(t *testing.T) {
	a := newTestAsaasAdapter()
	fees := a.GetFees()
	if fees.PixPercent != 0.0099 {
		t.Errorf("PixPercent = %f, want 0.0099", fees.PixPercent)
	}
	if fees.BoletoFixed != 2.99 {
		t.Errorf("BoletoFixed = %f, want 2.99", fees.BoletoFixed)
	}
	if fees.CardPercent != 0.0499 {
		t.Errorf("CardPercent = %f, want 0.0499", fees.CardPercent)
	}
	if fees.CardFixed != 0.49 {
		t.Errorf("CardFixed = %f, want 0.49", fees.CardFixed)
	}
}

// --- toCanonicalPayment tests ---

func TestToCanonicalPayment_Basic(t *testing.T) {
	a := newTestAsaasAdapter()
	resp := &PaymentResponse{
		ID:          "pay_123",
		Status:      "CONFIRMED",
		Value:       100.50,
		NetValue:    97.51,
		BillingType: "PIX",
		DueDate:     "2026-03-01",
		InvoiceURL:  "https://invoice.url",
		BankSlipURL: "https://boleto.url",
	}

	result := a.toCanonicalPayment(resp)

	if result.GatewayPaymentID != "pay_123" {
		t.Errorf("GatewayPaymentID = %q, want %q", result.GatewayPaymentID, "pay_123")
	}
	if result.Status != gateway.StatusConfirmed {
		t.Errorf("Status = %q, want %q", result.Status, gateway.StatusConfirmed)
	}
	if result.GatewayRawStatus != "CONFIRMED" {
		t.Errorf("GatewayRawStatus = %q, want %q", result.GatewayRawStatus, "CONFIRMED")
	}
	if result.Amount != 100.50 {
		t.Errorf("Amount = %f, want 100.50", result.Amount)
	}
	if result.NetAmount != 97.51 {
		t.Errorf("NetAmount = %f, want 97.51", result.NetAmount)
	}
	if result.BillingType != gateway.BillingPIX {
		t.Errorf("BillingType = %q, want %q", result.BillingType, gateway.BillingPIX)
	}
	if result.DueDate != "2026-03-01" {
		t.Errorf("DueDate = %q, want %q", result.DueDate, "2026-03-01")
	}
	if result.InvoiceURL != "https://invoice.url" {
		t.Errorf("InvoiceURL = %q", result.InvoiceURL)
	}
	if result.BoletoURL != "https://boleto.url" {
		t.Errorf("BoletoURL = %q", result.BoletoURL)
	}
}

func TestToCanonicalPayment_WithPixQRCode(t *testing.T) {
	a := newTestAsaasAdapter()
	resp := &PaymentResponse{
		ID:          "pay_pix",
		Status:      "PENDING",
		Value:       50.00,
		BillingType: "PIX",
		PixQRCode: &PixQRCode{
			EncodedImage:   "base64-image-data",
			Payload:        "pix-copy-paste-code",
			ExpirationDate: "2026-03-01 23:59:59",
		},
	}

	result := a.toCanonicalPayment(resp)

	if result.PixQRCodeBase64 != "base64-image-data" {
		t.Errorf("PixQRCodeBase64 = %q", result.PixQRCodeBase64)
	}
	if result.PixCopyPaste != "pix-copy-paste-code" {
		t.Errorf("PixCopyPaste = %q", result.PixCopyPaste)
	}
	if result.PixExpiration != "2026-03-01 23:59:59" {
		t.Errorf("PixExpiration = %q", result.PixExpiration)
	}
}

func TestToCanonicalPayment_WithConfirmedDate(t *testing.T) {
	a := newTestAsaasAdapter()
	resp := &PaymentResponse{
		ID:            "pay_conf",
		Status:        "RECEIVED",
		Value:         200.00,
		BillingType:   "BOLETO",
		ConfirmedDate: "2026-02-15",
	}

	result := a.toCanonicalPayment(resp)

	if result.ConfirmedAt == nil {
		t.Fatal("ConfirmedAt should not be nil")
	}
	if result.ConfirmedAt.Year() != 2026 || result.ConfirmedAt.Month() != 2 || result.ConfirmedAt.Day() != 15 {
		t.Errorf("ConfirmedAt = %v, want 2026-02-15", result.ConfirmedAt)
	}
}

func TestToCanonicalPayment_NoPixQRCode(t *testing.T) {
	a := newTestAsaasAdapter()
	resp := &PaymentResponse{
		ID:          "pay_no_pix",
		Status:      "PENDING",
		Value:       100.00,
		BillingType: "BOLETO",
	}

	result := a.toCanonicalPayment(resp)

	if result.PixQRCodeBase64 != "" {
		t.Error("PixQRCodeBase64 should be empty for non-PIX payment")
	}
	if result.PixCopyPaste != "" {
		t.Error("PixCopyPaste should be empty for non-PIX payment")
	}
}

// --- Helper function tests ---

func TestIsPaymentConfirmed(t *testing.T) {
	confirmedStatuses := []string{PaymentStatusReceived, PaymentStatusConfirmed, PaymentStatusReceivedInCash}
	for _, s := range confirmedStatuses {
		if !IsPaymentConfirmed(s) {
			t.Errorf("IsPaymentConfirmed(%q) = false, want true", s)
		}
	}

	notConfirmed := []string{PaymentStatusPending, PaymentStatusOverdue, PaymentStatusRefunded, "UNKNOWN"}
	for _, s := range notConfirmed {
		if IsPaymentConfirmed(s) {
			t.Errorf("IsPaymentConfirmed(%q) = true, want false", s)
		}
	}
}

func TestIsPaymentPending(t *testing.T) {
	pendingStatuses := []string{PaymentStatusPending, PaymentStatusAwaitingRisk}
	for _, s := range pendingStatuses {
		if !IsPaymentPending(s) {
			t.Errorf("IsPaymentPending(%q) = false, want true", s)
		}
	}

	notPending := []string{PaymentStatusConfirmed, PaymentStatusOverdue, "UNKNOWN"}
	for _, s := range notPending {
		if IsPaymentPending(s) {
			t.Errorf("IsPaymentPending(%q) = true, want false", s)
		}
	}
}

func TestIsPaymentFailed(t *testing.T) {
	failedStatuses := []string{PaymentStatusRefunded, PaymentStatusChargeback, PaymentStatusChargebackDispute}
	for _, s := range failedStatuses {
		if !IsPaymentFailed(s) {
			t.Errorf("IsPaymentFailed(%q) = false, want true", s)
		}
	}

	notFailed := []string{PaymentStatusPending, PaymentStatusConfirmed, "UNKNOWN"}
	for _, s := range notFailed {
		if IsPaymentFailed(s) {
			t.Errorf("IsPaymentFailed(%q) = true, want false", s)
		}
	}
}

func TestParseAsaasDate_Valid(t *testing.T) {
	parsed, err := ParseAsaasDate("2026-02-15")
	if err != nil {
		t.Fatalf("ParseAsaasDate() error = %v", err)
	}
	if parsed.Year() != 2026 || parsed.Month() != 2 || parsed.Day() != 15 {
		t.Errorf("ParseAsaasDate() = %v, want 2026-02-15", parsed)
	}
}

func TestParseAsaasDate_Invalid(t *testing.T) {
	_, err := ParseAsaasDate("invalid-date")
	if err == nil {
		t.Error("ParseAsaasDate() with invalid date should return error")
	}
}

func TestParseAsaasDate_WrongFormat(t *testing.T) {
	_, err := ParseAsaasDate("15/02/2026")
	if err == nil {
		t.Error("ParseAsaasDate() with wrong format should return error")
	}
}

// --- APIError tests ---

func TestAPIError_WithErrors(t *testing.T) {
	apiErr := &APIError{
		Errors: []struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		}{
			{Code: "invalid_customer", Description: "Customer not found"},
		},
	}

	if apiErr.Error() != "Customer not found" {
		t.Errorf("APIError.Error() = %q, want %q", apiErr.Error(), "Customer not found")
	}
}

func TestAPIError_NoErrors(t *testing.T) {
	apiErr := &APIError{}
	if apiErr.Error() != "unknown Asaas API error" {
		t.Errorf("APIError.Error() = %q, want %q", apiErr.Error(), "unknown Asaas API error")
	}
}
