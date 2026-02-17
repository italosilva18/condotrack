package mercadopago

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/condotrack/api/internal/domain/gateway"
)

// MercadoPagoAdapter implements gateway.PaymentGateway using the Mercado Pago client.
type MercadoPagoAdapter struct {
	client        *Client
	fees          gateway.GatewayFees
	webhookSecret string
}

// NewMercadoPagoAdapter creates a new Mercado Pago gateway adapter.
func NewMercadoPagoAdapter(client *Client, fees gateway.GatewayFees, webhookSecret string) *MercadoPagoAdapter {
	return &MercadoPagoAdapter{
		client:        client,
		fees:          fees,
		webhookSecret: webhookSecret,
	}
}

func (a *MercadoPagoAdapter) Name() string { return "mercadopago" }

// CreateCustomer creates a customer on Mercado Pago.
func (a *MercadoPagoAdapter) CreateCustomer(ctx context.Context, req gateway.CreateCustomerRequest) (*gateway.CustomerResponse, error) {
	// Parse name into first/last
	firstName, lastName := splitName(req.Name)

	// Determine document type
	docType := "CPF"
	if len(req.Document) > 11 {
		docType = "CNPJ"
	}

	mpCustomer := &MPCustomer{
		Email:     req.Email,
		FirstName: firstName,
		LastName:  lastName,
		Identification: &MPIdentification{
			Type:   docType,
			Number: req.Document,
		},
	}

	if req.Phone != "" {
		areaCode, number := splitPhone(req.Phone)
		mpCustomer.Phone = &MPPhone{
			AreaCode: areaCode,
			Number:   number,
		}
	}

	result, err := a.client.FindOrCreateCustomer(ctx, mpCustomer)
	if err != nil {
		return nil, err
	}

	return &gateway.CustomerResponse{
		GatewayID: result.ID,
		Name:      result.FirstName + " " + result.LastName,
		Email:     result.Email,
		Document:  getDocumentNumber(result),
	}, nil
}

// FindCustomerByDocument searches for a customer by CPF/CNPJ via email lookup.
func (a *MercadoPagoAdapter) FindCustomerByDocument(ctx context.Context, document string) (*gateway.CustomerResponse, error) {
	// MP doesn't have a direct search by document; return nil
	return nil, nil
}

// CreatePixPayment creates a PIX payment on Mercado Pago.
func (a *MercadoPagoAdapter) CreatePixPayment(ctx context.Context, req gateway.CreatePaymentRequest) (*gateway.PaymentResponse, error) {
	expiration := req.DueDate.Add(24 * time.Hour).Format(time.RFC3339)

	mpReq := &MPCreatePaymentRequest{
		TransactionAmount: req.Amount,
		Description:       req.Description,
		PaymentMethodID:   MPMethodPIX,
		ExternalReference: req.ExternalReference,
		DateOfExpiration:  expiration,
		Payer: MPPayer{
			Email: req.CustomerGatewayID, // For PIX, we pass email as payer
		},
	}

	resp, err := a.client.CreatePayment(ctx, mpReq)
	if err != nil {
		return nil, err
	}

	return a.toCanonicalPayment(resp), nil
}

// CreateBoletoPayment creates a Boleto payment on Mercado Pago.
func (a *MercadoPagoAdapter) CreateBoletoPayment(ctx context.Context, req gateway.CreatePaymentRequest) (*gateway.PaymentResponse, error) {
	expiration := req.DueDate.Add(3 * 24 * time.Hour).Format(time.RFC3339)

	mpReq := &MPCreatePaymentRequest{
		TransactionAmount: req.Amount,
		Description:       req.Description,
		PaymentMethodID:   MPMethodBoleto,
		ExternalReference: req.ExternalReference,
		DateOfExpiration:  expiration,
		Payer: MPPayer{
			Email: req.CustomerGatewayID,
		},
	}

	resp, err := a.client.CreatePayment(ctx, mpReq)
	if err != nil {
		return nil, err
	}

	return a.toCanonicalPayment(resp), nil
}

// CreateCardPayment creates a credit card payment on Mercado Pago.
func (a *MercadoPagoAdapter) CreateCardPayment(ctx context.Context, req gateway.CreateCardPaymentRequest) (*gateway.PaymentResponse, error) {
	docType := "CPF"
	if len(req.HolderDoc) > 11 {
		docType = "CNPJ"
	}

	mpReq := &MPCreatePaymentRequest{
		TransactionAmount: req.Amount,
		Description:       req.Description,
		PaymentMethodID:   MPMethodVisa, // Will be resolved by MP based on token
		ExternalReference: req.ExternalReference,
		Installments:      req.Installments,
		Token:             req.CardNumber, // In MP, this should be a card token from frontend
		Payer: MPPayer{
			Email: req.HolderEmail,
			Identification: &MPIdentification{
				Type:   docType,
				Number: req.HolderDoc,
			},
		},
	}

	if mpReq.Installments < 1 {
		mpReq.Installments = 1
	}

	resp, err := a.client.CreatePayment(ctx, mpReq)
	if err != nil {
		return nil, err
	}

	return a.toCanonicalPayment(resp), nil
}

// GetPayment retrieves a payment from Mercado Pago.
func (a *MercadoPagoAdapter) GetPayment(ctx context.Context, gatewayPaymentID string) (*gateway.PaymentResponse, error) {
	resp, err := a.client.GetPayment(ctx, gatewayPaymentID)
	if err != nil {
		return nil, err
	}
	return a.toCanonicalPayment(resp), nil
}

// RefundPayment refunds a payment on Mercado Pago.
func (a *MercadoPagoAdapter) RefundPayment(ctx context.Context, gatewayPaymentID string, amount float64) (*gateway.PaymentResponse, error) {
	resp, err := a.client.RefundPayment(ctx, gatewayPaymentID, amount)
	if err != nil {
		return nil, err
	}
	return a.toCanonicalPayment(resp), nil
}

// CancelPayment cancels a payment on Mercado Pago.
func (a *MercadoPagoAdapter) CancelPayment(ctx context.Context, gatewayPaymentID string) error {
	return a.client.CancelPayment(ctx, gatewayPaymentID)
}

// ParseWebhookEvent parses a Mercado Pago webhook event into the canonical format.
func (a *MercadoPagoAdapter) ParseWebhookEvent(ctx context.Context, headers map[string]string, body []byte) (*gateway.WebhookEvent, error) {
	return ParseWebhookEvent(ctx, a, headers, body)
}

// ValidateWebhookSignature validates the Mercado Pago webhook signature.
func (a *MercadoPagoAdapter) ValidateWebhookSignature(ctx context.Context, headers map[string]string, body []byte) bool {
	if a.webhookSecret == "" {
		return true // No secret configured, skip validation
	}
	return ValidateWebhookSignature(a.webhookSecret, headers, body)
}

// GetFees returns the gateway fee configuration.
func (a *MercadoPagoAdapter) GetFees() gateway.GatewayFees {
	return a.fees
}

// NormalizeStatus translates Mercado Pago status to canonical status.
func (a *MercadoPagoAdapter) NormalizeStatus(mpStatus string) string {
	switch mpStatus {
	case MPStatusPending:
		return gateway.StatusPending
	case MPStatusApproved:
		return gateway.StatusConfirmed
	case MPStatusAuthorized:
		return gateway.StatusPending
	case MPStatusInProcess:
		return gateway.StatusPending
	case MPStatusInMediation:
		return gateway.StatusChargeback
	case MPStatusRejected:
		return gateway.StatusFailed
	case MPStatusCancelled:
		return gateway.StatusCancelled
	case MPStatusRefunded:
		return gateway.StatusRefunded
	case MPStatusChargedBack:
		return gateway.StatusChargeback
	default:
		return gateway.StatusFailed
	}
}

// normalizeBillingType maps MP payment_type_id to canonical billing type.
func (a *MercadoPagoAdapter) normalizeBillingType(paymentTypeID string) string {
	switch paymentTypeID {
	case "bank_transfer":
		return gateway.BillingPIX
	case "ticket":
		return gateway.BillingBoleto
	case "credit_card":
		return gateway.BillingCreditCard
	case "debit_card":
		return gateway.BillingDebitCard
	default:
		return paymentTypeID
	}
}

// normalizeEventType maps MP webhook action to canonical event type.
func (a *MercadoPagoAdapter) normalizeEventType(action, paymentStatus string) string {
	if action == "payment.created" {
		return gateway.EventPaymentCreated
	}
	// For payment.updated, derive from payment status
	switch paymentStatus {
	case MPStatusApproved:
		return gateway.EventPaymentConfirmed
	case MPStatusRefunded:
		return gateway.EventPaymentRefunded
	case MPStatusCancelled:
		return gateway.EventPaymentDeleted
	case MPStatusChargedBack, MPStatusInMediation:
		return gateway.EventPaymentChargeback
	case MPStatusRejected:
		return gateway.EventPaymentFailed
	default:
		return gateway.EventPaymentCreated
	}
}

// toCanonicalPayment converts MP PaymentResponse to canonical format.
func (a *MercadoPagoAdapter) toCanonicalPayment(resp *MPPaymentResponse) *gateway.PaymentResponse {
	paymentID := strconv.FormatInt(resp.ID, 10)

	result := &gateway.PaymentResponse{
		GatewayPaymentID: paymentID,
		Status:           a.NormalizeStatus(resp.Status),
		GatewayRawStatus: resp.Status,
		Amount:           resp.TransactionAmount,
		NetAmount:        resp.NetReceivedAmount,
		BillingType:      a.normalizeBillingType(resp.PaymentTypeID),
	}

	// Parse dates
	if resp.DateApproved != nil && *resp.DateApproved != "" {
		if t, err := ParseMPDateTime(*resp.DateApproved); err == nil {
			result.ConfirmedAt = &t
			result.PaidAt = &t
		}
	}

	// PIX QR Code
	if resp.PointOfInteraction != nil && resp.PointOfInteraction.TransactionData != nil {
		td := resp.PointOfInteraction.TransactionData
		result.PixQRCodeBase64 = td.QRCodeBase64
		result.PixCopyPaste = td.QRCode
	}

	// Boleto / external URL
	if resp.TransactionDetails != nil && resp.TransactionDetails.ExternalResourceURL != "" {
		result.BoletoURL = resp.TransactionDetails.ExternalResourceURL
		result.InvoiceURL = resp.TransactionDetails.ExternalResourceURL
	}

	return result
}

// splitName splits a full name into first and last name.
func splitName(fullName string) (string, string) {
	parts := strings.SplitN(strings.TrimSpace(fullName), " ", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

// splitPhone splits a phone number into area code and number.
func splitPhone(phone string) (string, string) {
	cleaned := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)
	if len(cleaned) >= 10 {
		return cleaned[:2], cleaned[2:]
	}
	return "", cleaned
}

// getDocumentNumber extracts document number from customer.
func getDocumentNumber(customer *MPCustomer) string {
	if customer.Identification != nil {
		return customer.Identification.Number
	}
	return ""
}

// GetClient returns the underlying MP client (for webhook processing).
func (a *MercadoPagoAdapter) GetClient() *Client {
	return a.client
}
