package mercadopago

import (
	"testing"

	"github.com/condotrack/api/internal/domain/gateway"
)

func newTestAdapter() *MercadoPagoAdapter {
	return NewMercadoPagoAdapter(nil, gateway.GatewayFees{
		PixPercent:  0.0099,
		BoletoFixed: 3.49,
		CardPercent: 0.0499,
		CardFixed:   0.39,
	}, "test-secret")
}

func TestName(t *testing.T) {
	a := newTestAdapter()
	if a.Name() != "mercadopago" {
		t.Errorf("expected name 'mercadopago', got %s", a.Name())
	}
}

func TestNormalizeStatus_Pending(t *testing.T) {
	a := newTestAdapter()
	if a.NormalizeStatus(MPStatusPending) != gateway.StatusPending {
		t.Error("pending should map to pending")
	}
}

func TestNormalizeStatus_Approved(t *testing.T) {
	a := newTestAdapter()
	if a.NormalizeStatus(MPStatusApproved) != gateway.StatusConfirmed {
		t.Error("approved should map to confirmed")
	}
}

func TestNormalizeStatus_Authorized(t *testing.T) {
	a := newTestAdapter()
	if a.NormalizeStatus(MPStatusAuthorized) != gateway.StatusPending {
		t.Error("authorized should map to pending")
	}
}

func TestNormalizeStatus_InProcess(t *testing.T) {
	a := newTestAdapter()
	if a.NormalizeStatus(MPStatusInProcess) != gateway.StatusPending {
		t.Error("in_process should map to pending")
	}
}

func TestNormalizeStatus_InMediation(t *testing.T) {
	a := newTestAdapter()
	if a.NormalizeStatus(MPStatusInMediation) != gateway.StatusChargeback {
		t.Error("in_mediation should map to chargeback")
	}
}

func TestNormalizeStatus_Rejected(t *testing.T) {
	a := newTestAdapter()
	if a.NormalizeStatus(MPStatusRejected) != gateway.StatusFailed {
		t.Error("rejected should map to failed")
	}
}

func TestNormalizeStatus_Cancelled(t *testing.T) {
	a := newTestAdapter()
	if a.NormalizeStatus(MPStatusCancelled) != gateway.StatusCancelled {
		t.Error("cancelled should map to cancelled")
	}
}

func TestNormalizeStatus_Refunded(t *testing.T) {
	a := newTestAdapter()
	if a.NormalizeStatus(MPStatusRefunded) != gateway.StatusRefunded {
		t.Error("refunded should map to refunded")
	}
}

func TestNormalizeStatus_ChargedBack(t *testing.T) {
	a := newTestAdapter()
	if a.NormalizeStatus(MPStatusChargedBack) != gateway.StatusChargeback {
		t.Error("charged_back should map to chargeback")
	}
}

func TestNormalizeStatus_Unknown(t *testing.T) {
	a := newTestAdapter()
	if a.NormalizeStatus("unknown_status") != gateway.StatusFailed {
		t.Error("unknown status should map to failed")
	}
}

func TestNormalizeBillingType_BankTransfer(t *testing.T) {
	a := newTestAdapter()
	if a.normalizeBillingType("bank_transfer") != gateway.BillingPIX {
		t.Error("bank_transfer should map to pix")
	}
}

func TestNormalizeBillingType_Ticket(t *testing.T) {
	a := newTestAdapter()
	if a.normalizeBillingType("ticket") != gateway.BillingBoleto {
		t.Error("ticket should map to boleto")
	}
}

func TestNormalizeBillingType_CreditCard(t *testing.T) {
	a := newTestAdapter()
	if a.normalizeBillingType("credit_card") != gateway.BillingCreditCard {
		t.Error("credit_card should map to credit_card")
	}
}

func TestNormalizeBillingType_DebitCard(t *testing.T) {
	a := newTestAdapter()
	if a.normalizeBillingType("debit_card") != gateway.BillingDebitCard {
		t.Error("debit_card should map to debit_card")
	}
}

func TestNormalizeEventType_Created(t *testing.T) {
	a := newTestAdapter()
	if a.normalizeEventType("payment.created", "") != gateway.EventPaymentCreated {
		t.Error("payment.created should map to payment_created")
	}
}

func TestNormalizeEventType_Updated_Approved(t *testing.T) {
	a := newTestAdapter()
	if a.normalizeEventType("payment.updated", MPStatusApproved) != gateway.EventPaymentConfirmed {
		t.Error("payment.updated + approved should map to payment_confirmed")
	}
}

func TestNormalizeEventType_Updated_Refunded(t *testing.T) {
	a := newTestAdapter()
	if a.normalizeEventType("payment.updated", MPStatusRefunded) != gateway.EventPaymentRefunded {
		t.Error("payment.updated + refunded should map to payment_refunded")
	}
}

func TestNormalizeEventType_Updated_Cancelled(t *testing.T) {
	a := newTestAdapter()
	if a.normalizeEventType("payment.updated", MPStatusCancelled) != gateway.EventPaymentDeleted {
		t.Error("payment.updated + cancelled should map to payment_deleted")
	}
}

func TestNormalizeEventType_Updated_ChargedBack(t *testing.T) {
	a := newTestAdapter()
	if a.normalizeEventType("payment.updated", MPStatusChargedBack) != gateway.EventPaymentChargeback {
		t.Error("payment.updated + charged_back should map to payment_chargeback")
	}
}

func TestGetFees(t *testing.T) {
	a := newTestAdapter()
	fees := a.GetFees()
	if fees.PixPercent != 0.0099 {
		t.Errorf("expected PIX percent 0.0099, got %f", fees.PixPercent)
	}
	if fees.BoletoFixed != 3.49 {
		t.Errorf("expected boleto fixed 3.49, got %f", fees.BoletoFixed)
	}
	if fees.CardPercent != 0.0499 {
		t.Errorf("expected card percent 0.0499, got %f", fees.CardPercent)
	}
	if fees.CardFixed != 0.39 {
		t.Errorf("expected card fixed 0.39, got %f", fees.CardFixed)
	}
}

func TestToCanonicalPayment(t *testing.T) {
	a := newTestAdapter()
	approvedDate := "2026-02-10T15:30:00.000-03:00"
	resp := &MPPaymentResponse{
		ID:                12345,
		Status:            MPStatusApproved,
		StatusDetail:      "accredited",
		TransactionAmount: 100.50,
		NetReceivedAmount: 95.50,
		PaymentTypeID:     "bank_transfer",
		DateApproved:      &approvedDate,
		PointOfInteraction: &MPPointOfInteraction{
			TransactionData: &MPTransactionData{
				QRCode:       "pix-copy-paste-code",
				QRCodeBase64: "base64-qr-image",
			},
		},
	}

	result := a.toCanonicalPayment(resp)

	if result.GatewayPaymentID != "12345" {
		t.Errorf("expected payment ID '12345', got %s", result.GatewayPaymentID)
	}
	if result.Status != gateway.StatusConfirmed {
		t.Errorf("expected status confirmed, got %s", result.Status)
	}
	if result.GatewayRawStatus != MPStatusApproved {
		t.Errorf("expected raw status approved, got %s", result.GatewayRawStatus)
	}
	if result.Amount != 100.50 {
		t.Errorf("expected amount 100.50, got %f", result.Amount)
	}
	if result.NetAmount != 95.50 {
		t.Errorf("expected net amount 95.50, got %f", result.NetAmount)
	}
	if result.BillingType != gateway.BillingPIX {
		t.Errorf("expected billing type pix, got %s", result.BillingType)
	}
	if result.PixQRCodeBase64 != "base64-qr-image" {
		t.Error("expected PIX QR code base64")
	}
	if result.PixCopyPaste != "pix-copy-paste-code" {
		t.Error("expected PIX copy paste code")
	}
	if result.PaidAt == nil {
		t.Error("expected PaidAt to be set")
	}
	if result.ConfirmedAt == nil {
		t.Error("expected ConfirmedAt to be set")
	}
}

func TestSplitName(t *testing.T) {
	tests := []struct {
		input     string
		firstName string
		lastName  string
	}{
		{"John Doe", "John", "Doe"},
		{"John", "John", ""},
		{"Maria da Silva", "Maria", "da Silva"},
		{"  Spaces  ", "Spaces", ""},
	}

	for _, tt := range tests {
		f, l := splitName(tt.input)
		if f != tt.firstName || l != tt.lastName {
			t.Errorf("splitName(%q) = (%q, %q), want (%q, %q)", tt.input, f, l, tt.firstName, tt.lastName)
		}
	}
}

func TestSplitPhone(t *testing.T) {
	area, num := splitPhone("11987654321")
	if area != "11" || num != "987654321" {
		t.Errorf("splitPhone: got (%s, %s), want (11, 987654321)", area, num)
	}

	area, num = splitPhone("(11) 98765-4321")
	if area != "11" || num != "987654321" {
		t.Errorf("splitPhone with formatting: got (%s, %s)", area, num)
	}
}

func TestValidateWebhookSignature_EmptySignature(t *testing.T) {
	result := ValidateWebhookSignature("secret", map[string]string{}, []byte("{}"))
	if result {
		t.Error("expected false for empty signature")
	}
}

func TestValidateWebhookSignature_MissingParts(t *testing.T) {
	headers := map[string]string{
		"x-signature":  "invalid",
		"x-request-id": "req-123",
	}
	result := ValidateWebhookSignature("secret", headers, []byte("{}"))
	if result {
		t.Error("expected false for missing ts/v1")
	}
}

func TestParseMPDateTime(t *testing.T) {
	tests := []struct {
		input string
		ok    bool
	}{
		{"2026-02-10T15:30:00.000-03:00", true},
		{"2026-02-10T15:30:00Z", true},
		{"2026-02-10T15:30:00.000Z", true},
		{"invalid", false},
	}

	for _, tt := range tests {
		_, err := ParseMPDateTime(tt.input)
		if (err == nil) != tt.ok {
			t.Errorf("ParseMPDateTime(%q): got err=%v, want ok=%v", tt.input, err, tt.ok)
		}
	}
}
