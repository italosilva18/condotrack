package entity

import (
	"encoding/json"
	"time"
)

// InspectionType constants
const (
	InspectionTypeRoutine    = "routine"
	InspectionTypePreventive = "preventive"
	InspectionTypeCorrective = "corrective"
	InspectionTypeEmergency  = "emergency"
)

// InspectionStatus constants
const (
	InspectionStatusScheduled  = "scheduled"
	InspectionStatusInProgress = "in_progress"
	InspectionStatusCompleted  = "completed"
	InspectionStatusCancelled  = "cancelled"
)

// Inspection represents an inspection entity
type Inspection struct {
	ID              string          `db:"id" json:"id"`
	ContractID      string          `db:"contract_id" json:"contract_id"`
	ContractName    string          `db:"contract_name" json:"contract_name"`
	InspectorID     string          `db:"inspector_id" json:"inspector_id"`
	InspectorName   string          `db:"inspector_name" json:"inspector_name"`
	InspectionDate  time.Time       `db:"inspection_date" json:"inspection_date"`
	InspectionType  string          `db:"inspection_type" json:"inspection_type"`
	Status          string          `db:"status" json:"status"`
	Findings        *string         `db:"findings" json:"findings,omitempty"`
	Recommendations *string         `db:"recommendations" json:"recommendations,omitempty"`
	Photos          json.RawMessage `db:"photos" json:"photos,omitempty"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       *time.Time      `db:"updated_at" json:"updated_at,omitempty"`
}

// InspectionWithDetails represents an inspection with full contract and inspector information
type InspectionWithDetails struct {
	Inspection
	GestorName string `db:"gestor_name" json:"gestor_name,omitempty"`
}

// CreateInspectionRequest represents the request to create an inspection
type CreateInspectionRequest struct {
	ContractID      string          `json:"contract_id" binding:"required"`
	InspectorID     string          `json:"inspector_id" binding:"required"`
	InspectionDate  *time.Time      `json:"inspection_date,omitempty"`
	InspectionType  string          `json:"inspection_type" binding:"required"`
	Status          string          `json:"status,omitempty"`
	Findings        *string         `json:"findings,omitempty"`
	Recommendations *string         `json:"recommendations,omitempty"`
	Photos          json.RawMessage `json:"photos,omitempty"`
}

// UpdateInspectionRequest represents the request to update an inspection
type UpdateInspectionRequest struct {
	ContractID      *string         `json:"contract_id,omitempty"`
	InspectorID     *string         `json:"inspector_id,omitempty"`
	InspectionDate  *time.Time      `json:"inspection_date,omitempty"`
	InspectionType  *string         `json:"inspection_type,omitempty"`
	Status          *string         `json:"status,omitempty"`
	Findings        *string         `json:"findings,omitempty"`
	Recommendations *string         `json:"recommendations,omitempty"`
	Photos          json.RawMessage `json:"photos,omitempty"`
}

// InspectionFilter represents filters for listing inspections
type InspectionFilter struct {
	ContractID  string
	InspectorID string
	Status      string
	StartDate   *time.Time
	EndDate     *time.Time
}

// ValidInspectionTypes returns all valid inspection types
func ValidInspectionTypes() []string {
	return []string{
		InspectionTypeRoutine,
		InspectionTypePreventive,
		InspectionTypeCorrective,
		InspectionTypeEmergency,
	}
}

// ValidInspectionStatuses returns all valid inspection statuses
func ValidInspectionStatuses() []string {
	return []string{
		InspectionStatusScheduled,
		InspectionStatusInProgress,
		InspectionStatusCompleted,
		InspectionStatusCancelled,
	}
}

// IsValidInspectionType checks if the given type is valid
func IsValidInspectionType(t string) bool {
	for _, valid := range ValidInspectionTypes() {
		if t == valid {
			return true
		}
	}
	return false
}

// IsValidInspectionStatus checks if the given status is valid
func IsValidInspectionStatus(s string) bool {
	for _, valid := range ValidInspectionStatuses() {
		if s == valid {
			return true
		}
	}
	return false
}
