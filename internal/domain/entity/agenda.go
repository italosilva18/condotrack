package entity

import "time"

// EventType represents the type of calendar event
type EventType string

const (
	EventTypeAudit      EventType = "audit"
	EventTypeInspection EventType = "inspection"
	EventTypeMeeting    EventType = "meeting"
	EventTypeTask       EventType = "task"
	EventTypeOther      EventType = "other"
)

// ValidEventTypes returns all valid event types
func ValidEventTypes() []EventType {
	return []EventType{
		EventTypeAudit,
		EventTypeInspection,
		EventTypeMeeting,
		EventTypeTask,
		EventTypeOther,
	}
}

// IsValidEventType checks if the given event type is valid
func IsValidEventType(t EventType) bool {
	for _, valid := range ValidEventTypes() {
		if t == valid {
			return true
		}
	}
	return false
}

// AgendaEvent represents a calendar event entity
type AgendaEvent struct {
	ID             string     `db:"id" json:"id"`
	Title          string     `db:"title" json:"title"`
	Description    *string    `db:"description" json:"description,omitempty"`
	EventType      EventType  `db:"event_type" json:"event_type"`
	StartDatetime  time.Time  `db:"start_datetime" json:"start_datetime"`
	EndDatetime    time.Time  `db:"end_datetime" json:"end_datetime"`
	AllDay         bool       `db:"all_day" json:"all_day"`
	Location       *string    `db:"location" json:"location,omitempty"`
	ContractID     *string    `db:"contract_id" json:"contract_id,omitempty"`
	ContractName   *string    `db:"contract_name" json:"contract_name,omitempty"`
	UserID         *string    `db:"user_id" json:"user_id,omitempty"`
	UserName       *string    `db:"user_name" json:"user_name,omitempty"`
	RecurrenceRule *string    `db:"recurrence_rule" json:"recurrence_rule,omitempty"`
	Color          *string    `db:"color" json:"color,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// CreateEventRequest represents the request to create a calendar event
type CreateEventRequest struct {
	Title          string    `json:"title" binding:"required"`
	Description    *string   `json:"description,omitempty"`
	EventType      EventType `json:"event_type" binding:"required"`
	StartDatetime  time.Time `json:"start_datetime" binding:"required"`
	EndDatetime    time.Time `json:"end_datetime" binding:"required"`
	AllDay         bool      `json:"all_day"`
	Location       *string   `json:"location,omitempty"`
	ContractID     *string   `json:"contract_id,omitempty"`
	UserID         *string   `json:"user_id,omitempty"`
	RecurrenceRule *string   `json:"recurrence_rule,omitempty"`
	Color          *string   `json:"color,omitempty"`
}

// UpdateEventRequest represents the request to update a calendar event
type UpdateEventRequest struct {
	Title          *string    `json:"title,omitempty"`
	Description    *string    `json:"description,omitempty"`
	EventType      *EventType `json:"event_type,omitempty"`
	StartDatetime  *time.Time `json:"start_datetime,omitempty"`
	EndDatetime    *time.Time `json:"end_datetime,omitempty"`
	AllDay         *bool      `json:"all_day,omitempty"`
	Location       *string    `json:"location,omitempty"`
	ContractID     *string    `json:"contract_id,omitempty"`
	UserID         *string    `json:"user_id,omitempty"`
	RecurrenceRule *string    `json:"recurrence_rule,omitempty"`
	Color          *string    `json:"color,omitempty"`
}

// AgendaFilter represents filters for querying events
type AgendaFilter struct {
	StartDate  *time.Time
	EndDate    *time.Time
	ContractID *string
	UserID     *string
	EventType  *EventType
}
