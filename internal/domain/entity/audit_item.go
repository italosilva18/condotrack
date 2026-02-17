package entity

import "time"

// AuditItem represents an individual item within an audit
type AuditItem struct {
	ID          string     `db:"id" json:"id"`
	AuditID     string     `db:"audit_id" json:"audit_id"`
	CategoryID  string     `db:"category_id" json:"category_id"`
	ItemName    string     `db:"item_name" json:"item_name"`
	Score       float64    `db:"score" json:"score"`
	MaxScore    float64    `db:"max_score" json:"max_score"`
	Observation *string    `db:"observation" json:"observation,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
}

// AuditCategory represents a category of audit items
type AuditCategory struct {
	ID          string  `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Description *string `db:"description" json:"description,omitempty"`
	Weight      float64 `db:"weight" json:"weight"`
	Order       int     `db:"order_num" json:"order"`
}

// AuditItemWithCategory represents an audit item with its category info
type AuditItemWithCategory struct {
	AuditItem
	CategoryName string `db:"category_name" json:"category_name"`
}

// CalculateItemPercentage returns the percentage score for an item
func (ai *AuditItem) CalculateItemPercentage() float64 {
	if ai.MaxScore == 0 {
		return 0
	}
	return (ai.Score / ai.MaxScore) * 100
}

// CreateAuditCategoryRequest represents the request to create an audit category
type CreateAuditCategoryRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	Weight      float64 `json:"weight,omitempty"`
	Order       int     `json:"order,omitempty"`
}

// UpdateAuditCategoryRequest represents the request to update an audit category
type UpdateAuditCategoryRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Weight      *float64 `json:"weight,omitempty"`
	Order       *int     `json:"order,omitempty"`
}
