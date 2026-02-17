package entity

import (
	"encoding/json"
	"time"
)

// Audit represents an audit entity
type Audit struct {
	ID            string          `db:"id" json:"id"`
	ContractID    string          `db:"contract_id" json:"contract_id"`
	AuditorName   string          `db:"auditor_name" json:"auditor_name"`
	AuditDate     time.Time       `db:"audit_date" json:"audit_date"`
	Score         float64         `db:"score" json:"score"`
	TargetScore   float64         `db:"target_score" json:"target_score"`
	PreviousScore *float64        `db:"previous_score" json:"previous_score,omitempty"`
	Status        string          `db:"status" json:"status"`
	Observations  *string         `db:"observations" json:"observations,omitempty"`
	DataJSON      json.RawMessage `db:"data_json" json:"data_json,omitempty"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time      `db:"updated_at" json:"updated_at,omitempty"`
}

// AuditStatus constants
const (
	AuditStatusApproved = "approved"
	AuditStatusRejected = "rejected"
	AuditStatusPending  = "pending"
)

// AuditWithContract represents an audit with its contract information
type AuditWithContract struct {
	Audit
	ContractName string `db:"contract_name" json:"contract_name"`
	GestorName   string `db:"gestor_name" json:"gestor_name"`
}

// CreateAuditRequest represents the request to create an audit
type CreateAuditRequest struct {
	ContractID   string          `json:"contract_id" binding:"required"`
	AuditorName  string          `json:"auditor_name" binding:"required"`
	AuditDate    *time.Time      `json:"audit_date,omitempty"`
	Score        float64         `json:"score" binding:"required"`
	TargetScore  float64         `json:"target_score,omitempty"`
	Observations *string         `json:"observations,omitempty"`
	DataJSON     json.RawMessage `json:"data_json,omitempty"`
	Items        []AuditItemInput `json:"items,omitempty"`
}

// AuditItemInput represents an audit item in the create request
type AuditItemInput struct {
	CategoryID  string  `json:"category_id"`
	ItemName    string  `json:"item_name"`
	Score       float64 `json:"score"`
	MaxScore    float64 `json:"max_score"`
	Observation *string `json:"observation,omitempty"`
}

// AuditMeta represents metadata about audits for a contract
type AuditMeta struct {
	TotalAudits     int      `json:"total_audits"`
	AverageScore    float64  `json:"average_score"`
	LastScore       *float64 `json:"last_score,omitempty"`
	TargetScore     float64  `json:"target_score"`
	ApprovedCount   int      `json:"approved_count"`
	RejectedCount   int      `json:"rejected_count"`
	ImprovementRate *float64 `json:"improvement_rate,omitempty"`
}

// CalculateStatus determines the audit status based on score and target
func CalculateStatus(score, targetScore float64, tolerance float64) string {
	if score >= (targetScore - tolerance) {
		return AuditStatusApproved
	}
	return AuditStatusRejected
}

// UpdateAuditRequest represents the request to update an audit
type UpdateAuditRequest struct {
	AuditorName  *string          `json:"auditor_name,omitempty"`
	AuditDate    *time.Time       `json:"audit_date,omitempty"`
	Score        *float64         `json:"score,omitempty"`
	TargetScore  *float64         `json:"target_score,omitempty"`
	Observations *string          `json:"observations,omitempty"`
	DataJSON     json.RawMessage  `json:"data_json,omitempty"`
}
