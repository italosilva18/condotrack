package asaas

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/condotrack/api/internal/domain/gateway"
)

// AsaasAdapter implements gateway.PaymentGateway using the existing Asaas client.
type AsaasAdapter struct {
	client       *Client
	fees         gateway.GatewayFees
	webhookToken string
}

// NewAsaasAdapter creates a new Asaas gateway adapter.
func NewAsaasAdapter(client *Client, fees gateway.GatewayFees, webhookToken string) *AsaasAdapter {
	return &AsaasAdapter{
		client:       client,
		fees:         fees,
		webhookToken: webhookToken,
	}
}

func (a *AsaasAdapter) Name() string { return "asaas" }

// CreateCustomer creates a customer on Asaas.
func (a *AsaasAdapter) CreateCustomer(ctx context.Context, req gateway.CreateCustomerRequest) (*gateway.CustomerResponse, error) {
	customer, err := a.client.FindOrCreateCustomer(ctx, &CreateCustomerRequest{
		Name:    req.Name,
		Email:   req.Email,
		CPFCnpj: req.Document,
		Phone:   req.Phone,
	})
	if err != nil {
		return nil, err
	}
	return &gateway.CustomerResponse{
		GatewayID: customer.ID,
		Name:      customer.Name,
		Email:     customer.Email,
		Document:  customer.CPFCnpj,
	}, nil
}

// FindCustomerByDocument finds a customer by CPF/CNPJ.
func (a *AsaasAdapter) FindCustomerByDocument(ctx context.Context, document string) (*gateway.CustomerResponse, error) {
	customer, err := a.client.FindCustomerByCPF(ctx, document)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, nil
	}
	return &gateway.CustomerResponse{
		GatewayID: customer.ID,
		Name:      customer.Name,
		Email:     customer.Email,
		Document:  customer.CPFCnpj,
	}, nil
}

// CreatePixPayment creates a PIX payment on Asaas.
func (a *AsaasAdapter) CreatePixPayment(ctx context.Context, req gateway.CreatePaymentRequest) (*gateway.PaymentResponse, error) {
	asaasReq := &CreatePaymentRequest{
		Customer:    req.CustomerGatewayID,
		BillingType: "PIX",
		Value:       req.Amount,
		DueDate:     req.DueDate.Format("2006-01-02"),
		Description: req.Description,
		ExternalRef: req.ExternalReference,
	}
	resp, err := a.client.CreatePayment(ctx, asaasReq)
	if err != nil {
		return nil, err
	}
	return a.toCanonicalPayment(resp), nil
}

// CreateBoletoPayment creates a Boleto payment on Asaas.
func (a *AsaasAdapter) CreateBoletoPayment(ctx context.Context, req gateway.CreatePaymentRequest) (*gateway.PaymentResponse, error) {
	asaasReq := &CreatePaymentRequest{
		Customer:    req.CustomerGatewayID,
		BillingType: "BOLETO",
		Value:       req.Amount,
		DueDate:     req.DueDate.Format("2006-01-02"),
		Description: req.Description,
		ExternalRef: req.ExternalReference,
	}
	resp, err := a.client.CreatePayment(ctx, asaasReq)
	if err != nil {
		return nil, err
	}
	result := a.toCanonicalPayment(resp)

	// Try to get boleto barcode
	barcode, err := a.client.GetBoletoIdentificationField(ctx, resp.ID)
	if err == nil {
		result.BoletoBarCode = barcode
	}

	return result, nil
}

// CreateCardPayment creates a credit card payment on Asaas.
func (a *AsaasAdapter) CreateCardPayment(ctx context.Context, req gateway.CreateCardPaymentRequest) (*gateway.PaymentResponse, error) {
	asaasReq := &CreateCardPaymentRequest{
		Customer:    req.CustomerGatewayID,
		BillingType: "CREDIT_CARD",
		Value:       req.Amount,
		DueDate:     req.DueDate.Format("2006-01-02"),
		Description: req.Description,
		ExternalRef: req.ExternalReference,
		CreditCard: &CreditCard{
			HolderName:  req.HolderName,
			Number:      req.CardNumber,
			ExpiryMonth: req.CardExpMonth,
			ExpiryYear:  req.CardExpYear,
			Ccv:         req.CardCVV,
		},
		CreditCardHolder: &CreditCardHolder{
			Name:       req.HolderName,
			Email:      req.HolderEmail,
			CPFCnpj:    req.HolderDoc,
			PostalCode: req.HolderZip,
			Phone:      req.HolderPhone,
		},
		Installments: req.Installments,
	}
	resp, err := a.client.CreateCardPayment(ctx, asaasReq)
	if err != nil {
		return nil, err
	}
	return a.toCanonicalPayment(resp), nil
}

// GetPayment retrieves a payment from Asaas.
func (a *AsaasAdapter) GetPayment(ctx context.Context, gatewayPaymentID string) (*gateway.PaymentResponse, error) {
	resp, err := a.client.GetPayment(ctx, gatewayPaymentID)
	if err != nil {
		return nil, err
	}
	return a.toCanonicalPayment(resp), nil
}

// RefundPayment refunds a payment on Asaas.
func (a *AsaasAdapter) RefundPayment(ctx context.Context, gatewayPaymentID string, amount float64) (*gateway.PaymentResponse, error) {
	resp, err := a.client.RefundPayment(ctx, gatewayPaymentID, amount)
	if err != nil {
		return nil, err
	}
	return a.toCanonicalPayment(resp), nil
}

// CancelPayment cancels/deletes a payment on Asaas.
func (a *AsaasAdapter) CancelPayment(ctx context.Context, gatewayPaymentID string) error {
	return a.client.DeletePayment(ctx, gatewayPaymentID)
}

// ParseWebhookEvent parses an Asaas webhook event into the canonical format.
func (a *AsaasAdapter) ParseWebhookEvent(ctx context.Context, headers map[string]string, body []byte) (*gateway.WebhookEvent, error) {
	var event WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("failed to parse Asaas webhook: %w", err)
	}

	if event.Payment == nil {
		return nil, fmt.Errorf("webhook event has no payment data")
	}

	canonical := &gateway.WebhookEvent{
		GatewayEvent:     event.Event,
		GatewayName:      "asaas",
		PaymentID:        event.Payment.ID,
		CustomerID:       event.Payment.Customer,
		Amount:           event.Payment.Value,
		NetAmount:        event.Payment.NetValue,
		GatewayRawStatus: event.Payment.Status,
		Status:           a.NormalizeStatus(event.Payment.Status),
		BillingType:      a.normalizeBillingType(event.Payment.BillingType),
		ExternalRef:      event.Payment.ExternalReference,
		RawPayload:       body,
	}

	// Map event type
	canonical.EventType = a.normalizeEventType(event.Event)

	// Parse paid date
	if event.Payment.PaymentDate != "" {
		if t, err := time.Parse("2006-01-02", event.Payment.PaymentDate); err == nil {
			canonical.PaidAt = &t
		}
	}

	return canonical, nil
}

// ValidateWebhookSignature validates the Asaas webhook token.
func (a *AsaasAdapter) ValidateWebhookSignature(ctx context.Context, headers map[string]string, body []byte) bool {
	if a.webhookToken == "" {
		return true // No token configured, skip validation
	}
	token, ok := headers["asaas-access-token"]
	if !ok {
		return false
	}
	return token == a.webhookToken
}

// GetFees returns the gateway fee configuration.
func (a *AsaasAdapter) GetFees() gateway.GatewayFees {
	return a.fees
}

// NormalizeStatus translates Asaas status to canonical status.
func (a *AsaasAdapter) NormalizeStatus(asaasStatus string) string {
	switch asaasStatus {
	case "PENDING", "AWAITING_RISK_ANALYSIS":
		return gateway.StatusPending
	case "RECEIVED", "CONFIRMED", "RECEIVED_IN_CASH":
		return gateway.StatusConfirmed
	case "OVERDUE":
		return gateway.StatusOverdue
	case "REFUNDED", "REFUND_REQUESTED":
		return gateway.StatusRefunded
	case "CHARGEBACK_REQUESTED", "CHARGEBACK_DISPUTE", "AWAITING_CHARGEBACK_REVERSAL":
		return gateway.StatusChargeback
	default:
		return gateway.StatusFailed
	}
}

// normalizeBillingType translates Asaas billing type to canonical.
func (a *AsaasAdapter) normalizeBillingType(asaasBilling string) string {
	switch asaasBilling {
	case "PIX":
		return gateway.BillingPIX
	case "BOLETO":
		return gateway.BillingBoleto
	case "CREDIT_CARD":
		return gateway.BillingCreditCard
	case "DEBIT_CARD":
		return gateway.BillingDebitCard
	default:
		return asaasBilling
	}
}

// normalizeEventType translates Asaas webhook event to canonical event type.
func (a *AsaasAdapter) normalizeEventType(asaasEvent string) string {
	switch asaasEvent {
	case WebhookEventPaymentCreated:
		return gateway.EventPaymentCreated
	case WebhookEventPaymentConfirmed, WebhookEventPaymentReceived, WebhookEventPaymentReceivedInCash:
		return gateway.EventPaymentConfirmed
	case WebhookEventPaymentOverdue:
		return gateway.EventPaymentOverdue
	case WebhookEventPaymentRefunded:
		return gateway.EventPaymentRefunded
	case WebhookEventPaymentDeleted:
		return gateway.EventPaymentDeleted
	case WebhookEventPaymentChargebackRequested, WebhookEventPaymentChargebackDispute:
		return gateway.EventPaymentChargeback
	default:
		return asaasEvent
	}
}

// toCanonicalPayment converts Asaas PaymentResponse to canonical format.
func (a *AsaasAdapter) toCanonicalPayment(resp *PaymentResponse) *gateway.PaymentResponse {
	result := &gateway.PaymentResponse{
		GatewayPaymentID:      resp.ID,
		Status:                a.NormalizeStatus(resp.Status),
		GatewayRawStatus:      resp.Status,
		Amount:                resp.Value,
		NetAmount:             resp.NetValue,
		BillingType:           a.normalizeBillingType(resp.BillingType),
		DueDate:               resp.DueDate,
		InvoiceURL:            resp.InvoiceURL,
		BoletoURL:             resp.BankSlipURL,
		TransactionReceiptURL: resp.TransactionReceiptURL,
	}

	// Parse confirmed date
	if resp.ConfirmedDate != "" {
		if t, err := ParseAsaasDate(resp.ConfirmedDate); err == nil {
			result.ConfirmedAt = &t
		}
	}

	// PIX QR Code
	if resp.PixQRCode != nil {
		result.PixQRCodeBase64 = resp.PixQRCode.EncodedImage
		result.PixCopyPaste = resp.PixQRCode.Payload
		result.PixExpiration = resp.PixQRCode.ExpirationDate
	}

	return result
}
