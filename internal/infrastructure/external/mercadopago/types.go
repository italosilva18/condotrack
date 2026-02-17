package mercadopago

import (
	"fmt"
	"time"
)

// --- Customer types ---

// MPCustomer represents a Mercado Pago customer.
type MPCustomer struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Identification *MPIdentification `json:"identification,omitempty"`
	Phone         *MPPhone `json:"phone,omitempty"`
}

// MPIdentification represents customer identification (CPF/CNPJ).
type MPIdentification struct {
	Type   string `json:"type"`   // CPF or CNPJ
	Number string `json:"number"`
}

// MPPhone represents a phone number.
type MPPhone struct {
	AreaCode string `json:"area_code"`
	Number   string `json:"number"`
}

// MPCustomerSearchResult represents the search response for customers.
type MPCustomerSearchResult struct {
	Paging  MPPaging     `json:"paging"`
	Results []MPCustomer `json:"results"`
}

// MPPaging represents pagination info.
type MPPaging struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// --- Payment types ---

// MPCreatePaymentRequest represents a request to create a payment on MP.
type MPCreatePaymentRequest struct {
	TransactionAmount float64              `json:"transaction_amount"`
	Description       string               `json:"description,omitempty"`
	PaymentMethodID   string               `json:"payment_method_id"`
	Payer             MPPayer              `json:"payer"`
	ExternalReference string               `json:"external_reference,omitempty"`
	DateOfExpiration  string               `json:"date_of_expiration,omitempty"`
	Installments      int                  `json:"installments,omitempty"`
	Token             string               `json:"token,omitempty"`
	Card              *MPCardRequest       `json:"card,omitempty"`
	NotificationURL   string               `json:"notification_url,omitempty"`
}

// MPPayer represents the payer in a payment request.
type MPPayer struct {
	Email          string             `json:"email"`
	FirstName      string             `json:"first_name,omitempty"`
	LastName       string             `json:"last_name,omitempty"`
	Identification *MPIdentification  `json:"identification,omitempty"`
}

// MPCardRequest represents card info for a payment (tokenized).
type MPCardRequest struct {
	Token string `json:"token,omitempty"`
}

// MPPaymentResponse represents the MP payment response.
type MPPaymentResponse struct {
	ID                int64          `json:"id"`
	Status            string         `json:"status"`
	StatusDetail      string         `json:"status_detail"`
	DateCreated       string         `json:"date_created"`
	DateApproved      *string        `json:"date_approved,omitempty"`
	DateLastUpdated   string         `json:"date_last_updated"`
	TransactionAmount float64        `json:"transaction_amount"`
	NetReceivedAmount float64        `json:"net_received_amount"`
	CurrencyID        string         `json:"currency_id"`
	PaymentMethodID   string         `json:"payment_method_id"`
	PaymentTypeID     string         `json:"payment_type_id"`
	Installments      int            `json:"installments"`
	ExternalReference string         `json:"external_reference,omitempty"`
	Payer             *MPPayer       `json:"payer,omitempty"`
	PointOfInteraction *MPPointOfInteraction `json:"point_of_interaction,omitempty"`
	TransactionDetails *MPTransactionDetails `json:"transaction_details,omitempty"`
	FeeDetails        []MPFeeDetail  `json:"fee_details,omitempty"`
}

// MPPointOfInteraction holds PIX-specific data (QR code, copy-paste).
type MPPointOfInteraction struct {
	TransactionData *MPTransactionData `json:"transaction_data,omitempty"`
}

// MPTransactionData holds the QR code and ticket data.
type MPTransactionData struct {
	QRCode       string `json:"qr_code,omitempty"`
	QRCodeBase64 string `json:"qr_code_base64,omitempty"`
	TicketURL    string `json:"ticket_url,omitempty"`
}

// MPTransactionDetails holds transaction metadata.
type MPTransactionDetails struct {
	NetReceivedAmount   float64 `json:"net_received_amount"`
	TotalPaidAmount     float64 `json:"total_paid_amount"`
	InstallmentAmount   float64 `json:"installment_amount"`
	ExternalResourceURL string  `json:"external_resource_url,omitempty"`
}

// MPFeeDetail represents a fee breakdown item.
type MPFeeDetail struct {
	Type     string  `json:"type"`
	Amount   float64 `json:"amount"`
	FeePayer string  `json:"fee_payer"`
}

// --- Webhook types ---

// MPWebhookNotification represents a webhook notification from MP.
type MPWebhookNotification struct {
	ID          int64  `json:"id"`
	LiveMode    bool   `json:"live_mode"`
	Type        string `json:"type"` // "payment"
	DateCreated string `json:"date_created"`
	UserID      int64  `json:"user_id"`
	APIVersion  string `json:"api_version"`
	Action      string `json:"action"` // "payment.created", "payment.updated"
	Data        struct {
		ID string `json:"id"`
	} `json:"data"`
}

// --- Status constants (Mercado Pago) ---

const (
	MPStatusPending      = "pending"
	MPStatusApproved     = "approved"
	MPStatusAuthorized   = "authorized"
	MPStatusInProcess    = "in_process"
	MPStatusInMediation  = "in_mediation"
	MPStatusRejected     = "rejected"
	MPStatusCancelled    = "cancelled"
	MPStatusRefunded     = "refunded"
	MPStatusChargedBack  = "charged_back"
)

// --- Payment method IDs ---

const (
	MPMethodPIX         = "pix"
	MPMethodBoleto      = "bolbradesco"
	MPMethodVisa        = "visa"
	MPMethodMastercard  = "master"
	MPMethodElo         = "elo"
	MPMethodAmex        = "amex"
	MPMethodHipercard   = "hipercard"
)

// ParseMPDateTime parses a Mercado Pago datetime string.
func ParseMPDateTime(s string) (time.Time, error) {
	// MP uses ISO 8601 with timezone offset
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05.000-07:00",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse MP datetime: %s", s)
}
