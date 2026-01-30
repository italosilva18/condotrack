package asaas

import "time"

// Customer represents an Asaas customer
type Customer struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	Email                string  `json:"email,omitempty"`
	Phone                string  `json:"phone,omitempty"`
	MobilePhone          string  `json:"mobilePhone,omitempty"`
	CPFCnpj              string  `json:"cpfCnpj"`
	PostalCode           string  `json:"postalCode,omitempty"`
	Address              string  `json:"address,omitempty"`
	AddressNumber        string  `json:"addressNumber,omitempty"`
	Complement           string  `json:"complement,omitempty"`
	Province             string  `json:"province,omitempty"`
	ExternalReference    string  `json:"externalReference,omitempty"`
	NotificationDisabled bool    `json:"notificationDisabled,omitempty"`
	PersonType           string  `json:"personType,omitempty"`
}

// CreateCustomerRequest represents the request to create a customer
type CreateCustomerRequest struct {
	Name                 string `json:"name"`
	Email                string `json:"email,omitempty"`
	Phone                string `json:"phone,omitempty"`
	MobilePhone          string `json:"mobilePhone,omitempty"`
	CPFCnpj              string `json:"cpfCnpj"`
	PostalCode           string `json:"postalCode,omitempty"`
	Address              string `json:"address,omitempty"`
	AddressNumber        string `json:"addressNumber,omitempty"`
	Complement           string `json:"complement,omitempty"`
	Province             string `json:"province,omitempty"`
	ExternalReference    string `json:"externalReference,omitempty"`
	NotificationDisabled bool   `json:"notificationDisabled,omitempty"`
}

// CreatePaymentRequest represents the request to create a payment
type CreatePaymentRequest struct {
	Customer    string  `json:"customer"`
	BillingType string  `json:"billingType"` // BOLETO, CREDIT_CARD, PIX
	Value       float64 `json:"value"`
	DueDate     string  `json:"dueDate"` // Format: YYYY-MM-DD
	Description string  `json:"description,omitempty"`
	ExternalRef string  `json:"externalReference,omitempty"`

	// Discount
	Discount *Discount `json:"discount,omitempty"`

	// Fine and Interest
	Fine     *Fine     `json:"fine,omitempty"`
	Interest *Interest `json:"interest,omitempty"`
}

// CreateCardPaymentRequest represents the request to create a card payment
type CreateCardPaymentRequest struct {
	Customer         string              `json:"customer"`
	BillingType      string              `json:"billingType"` // CREDIT_CARD
	Value            float64             `json:"value"`
	DueDate          string              `json:"dueDate"`
	Description      string              `json:"description,omitempty"`
	ExternalRef      string              `json:"externalReference,omitempty"`
	CreditCard       *CreditCard         `json:"creditCard,omitempty"`
	CreditCardHolder *CreditCardHolder   `json:"creditCardHolderInfo,omitempty"`
	CreditCardToken  string              `json:"creditCardToken,omitempty"`
	Installments     int                 `json:"installmentCount,omitempty"`
}

// CreditCard represents credit card information
type CreditCard struct {
	HolderName  string `json:"holderName"`
	Number      string `json:"number"`
	ExpiryMonth string `json:"expiryMonth"`
	ExpiryYear  string `json:"expiryYear"`
	Ccv         string `json:"ccv"`
}

// CreditCardHolder represents credit card holder information
type CreditCardHolder struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	CPFCnpj    string `json:"cpfCnpj"`
	PostalCode string `json:"postalCode"`
	Address    string `json:"addressNumber"`
	Phone      string `json:"phone,omitempty"`
}

// Discount represents payment discount
type Discount struct {
	Value            float64 `json:"value"`
	DueDateLimitDays int     `json:"dueDateLimitDays,omitempty"`
	Type             string  `json:"type,omitempty"` // PERCENTAGE, FIXED
}

// Fine represents payment fine
type Fine struct {
	Value float64 `json:"value"`
	Type  string  `json:"type,omitempty"` // PERCENTAGE, FIXED
}

// Interest represents payment interest
type Interest struct {
	Value float64 `json:"value"`
	Type  string  `json:"type,omitempty"` // PERCENTAGE
}

// PaymentResponse represents a payment response from Asaas
type PaymentResponse struct {
	ID                string      `json:"id"`
	DateCreated       string      `json:"dateCreated"`
	Customer          string      `json:"customer"`
	PaymentLink       string      `json:"paymentLink,omitempty"`
	DueDate           string      `json:"dueDate"`
	Value             float64     `json:"value"`
	NetValue          float64     `json:"netValue"`
	BillingType       string      `json:"billingType"`
	Status            string      `json:"status"`
	Description       string      `json:"description,omitempty"`
	ExternalReference string      `json:"externalReference,omitempty"`
	ConfirmedDate     string      `json:"confirmedDate,omitempty"`
	OriginalValue     float64     `json:"originalValue,omitempty"`
	InterestValue     float64     `json:"interestValue,omitempty"`
	BankSlipURL       string      `json:"bankSlipUrl,omitempty"`
	InvoiceURL        string      `json:"invoiceUrl,omitempty"`
	InvoiceNumber     string      `json:"invoiceNumber,omitempty"`
	PixQRCode         *PixQRCode  `json:"pix,omitempty"`
	TransactionReceiptURL string  `json:"transactionReceiptUrl,omitempty"`
}

// PixQRCode represents PIX QR code information
type PixQRCode struct {
	EncodedImage   string `json:"encodedImage"`
	Payload        string `json:"payload"`
	ExpirationDate string `json:"expirationDate"`
}

// PaymentStatus constants
const (
	PaymentStatusPending         = "PENDING"
	PaymentStatusReceived        = "RECEIVED"
	PaymentStatusConfirmed       = "CONFIRMED"
	PaymentStatusOverdue         = "OVERDUE"
	PaymentStatusRefunded        = "REFUNDED"
	PaymentStatusReceivedInCash  = "RECEIVED_IN_CASH"
	PaymentStatusRefundRequested = "REFUND_REQUESTED"
	PaymentStatusChargeback      = "CHARGEBACK_REQUESTED"
	PaymentStatusChargebackDispute = "CHARGEBACK_DISPUTE"
	PaymentStatusAwaitingChargeback = "AWAITING_CHARGEBACK_REVERSAL"
	PaymentStatusDunningRequested   = "DUNNING_REQUESTED"
	PaymentStatusDunningReceived    = "DUNNING_RECEIVED"
	PaymentStatusAwaitingRisk       = "AWAITING_RISK_ANALYSIS"
)

// WebhookEvent represents an Asaas webhook event
type WebhookEvent struct {
	Event   string          `json:"event"`
	Payment *WebhookPayment `json:"payment,omitempty"`
}

// WebhookPayment represents payment data in webhook
type WebhookPayment struct {
	ID                string  `json:"id"`
	Customer          string  `json:"customer"`
	Value             float64 `json:"value"`
	NetValue          float64 `json:"netValue"`
	Status            string  `json:"status"`
	BillingType       string  `json:"billingType"`
	ExternalReference string  `json:"externalReference,omitempty"`
	PaymentDate       string  `json:"paymentDate,omitempty"`
	ConfirmedDate     string  `json:"confirmedDate,omitempty"`
}

// Webhook event types
const (
	WebhookEventPaymentCreated          = "PAYMENT_CREATED"
	WebhookEventPaymentUpdated          = "PAYMENT_UPDATED"
	WebhookEventPaymentConfirmed        = "PAYMENT_CONFIRMED"
	WebhookEventPaymentReceived         = "PAYMENT_RECEIVED"
	WebhookEventPaymentOverdue          = "PAYMENT_OVERDUE"
	WebhookEventPaymentDeleted          = "PAYMENT_DELETED"
	WebhookEventPaymentRestored         = "PAYMENT_RESTORED"
	WebhookEventPaymentRefunded         = "PAYMENT_REFUNDED"
	WebhookEventPaymentReceivedInCash   = "PAYMENT_RECEIVED_IN_CASH"
	WebhookEventPaymentChargebackRequested = "PAYMENT_CHARGEBACK_REQUESTED"
	WebhookEventPaymentChargebackDispute   = "PAYMENT_CHARGEBACK_DISPUTE"
	WebhookEventPaymentAwaitingChargeback  = "PAYMENT_AWAITING_CHARGEBACK_REVERSAL"
	WebhookEventPaymentDunningReceived     = "PAYMENT_DUNNING_RECEIVED"
	WebhookEventPaymentDunningRequested    = "PAYMENT_DUNNING_REQUESTED"
	WebhookEventPaymentBankSlipViewed      = "PAYMENT_BANK_SLIP_VIEWED"
	WebhookEventPaymentPixAddressKeyCreated = "PAYMENT_CHECKOUT_VIEWED"
)

// CustomerListResponse represents the response when listing customers
type CustomerListResponse struct {
	Object     string     `json:"object"`
	HasMore    bool       `json:"hasMore"`
	TotalCount int        `json:"totalCount"`
	Limit      int        `json:"limit"`
	Offset     int        `json:"offset"`
	Data       []Customer `json:"data"`
}

// APIError represents an Asaas API error
type APIError struct {
	Errors []struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"errors"`
}

// Error returns the error message
func (e *APIError) Error() string {
	if len(e.Errors) > 0 {
		return e.Errors[0].Description
	}
	return "unknown Asaas API error"
}

// IsPaymentConfirmed checks if payment status indicates confirmation
func IsPaymentConfirmed(status string) bool {
	return status == PaymentStatusReceived ||
		   status == PaymentStatusConfirmed ||
		   status == PaymentStatusReceivedInCash
}

// IsPaymentPending checks if payment is still pending
func IsPaymentPending(status string) bool {
	return status == PaymentStatusPending ||
		   status == PaymentStatusAwaitingRisk
}

// IsPaymentFailed checks if payment failed
func IsPaymentFailed(status string) bool {
	return status == PaymentStatusRefunded ||
		   status == PaymentStatusChargeback ||
		   status == PaymentStatusChargebackDispute
}

// ParseAsaasDate parses a date string from Asaas format
func ParseAsaasDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}
