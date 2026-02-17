package entity

import "time"

// PaymentTransaction represents an immutable audit trail entry for payment status changes.
type PaymentTransaction struct {
	ID              string     `db:"id" json:"id"`
	PaymentID       string     `db:"payment_id" json:"payment_id"`
	PreviousStatus  *string    `db:"previous_status" json:"previous_status,omitempty"`
	NewStatus       string     `db:"new_status" json:"new_status"`
	EventSource     string     `db:"event_source" json:"event_source"`
	EventType       string     `db:"event_type" json:"event_type"`
	GatewayEventID  *string    `db:"gateway_event_id" json:"gateway_event_id,omitempty"`
	Amount          *float64   `db:"amount" json:"amount,omitempty"`
	Description     *string    `db:"description" json:"description,omitempty"`
	RawPayload      *string    `db:"raw_payload" json:"raw_payload,omitempty"`
	IPAddress       *string    `db:"ip_address" json:"ip_address,omitempty"`
	UserAgent       *string    `db:"user_agent" json:"user_agent,omitempty"`
	TriggeredBy     *string    `db:"triggered_by" json:"triggered_by,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
}

// Event source constants
const (
	EventSourceWebhook   = "webhook"
	EventSourceManual    = "manual"
	EventSourceSystem    = "system"
	EventSourceScheduler = "scheduler"
	EventSourceAPI       = "api"
)

// Transaction event type constants
const (
	TxEventCreated        = "created"
	TxEventStatusChanged  = "status_changed"
	TxEventWebhookReceived = "webhook_received"
	TxEventRefundRequested = "refund_requested"
	TxEventError          = "error"
)
