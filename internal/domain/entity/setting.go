package entity

import "time"

// SettingType represents the type of a setting value
type SettingType string

const (
	SettingTypeString  SettingType = "string"
	SettingTypeNumber  SettingType = "number"
	SettingTypeBoolean SettingType = "boolean"
	SettingTypeJSON    SettingType = "json"
	SettingTypeSecret  SettingType = "secret"
)

// SettingCategory represents a category of settings
type SettingCategory string

const (
	CategoryGeneral SettingCategory = "general"
	CategoryPayment SettingCategory = "payment"
	CategoryAI      SettingCategory = "ai"
	CategoryRevenue SettingCategory = "revenue"
	CategoryStorage SettingCategory = "storage"
	CategoryEmail   SettingCategory = "email"
)

// Setting represents a system configuration setting
type Setting struct {
	ID              string          `db:"id" json:"id"`
	Key             string          `db:"setting_key" json:"key"`
	Value           *string         `db:"setting_value" json:"value,omitempty"`
	Type            SettingType     `db:"setting_type" json:"type"`
	Category        SettingCategory `db:"category" json:"category"`
	Label           string          `db:"label" json:"label"`
	Description     *string         `db:"description" json:"description,omitempty"`
	IsSecret        bool            `db:"is_secret" json:"is_secret"`
	IsRequired      bool            `db:"is_required" json:"is_required"`
	DefaultValue    *string         `db:"default_value" json:"default_value,omitempty"`
	ValidationRegex *string         `db:"validation_regex" json:"validation_regex,omitempty"`
	DisplayOrder    int             `db:"display_order" json:"display_order"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       *time.Time      `db:"updated_at" json:"updated_at,omitempty"`
}

// SettingPublic represents a setting with masked secret values
type SettingPublic struct {
	ID           string          `json:"id"`
	Key          string          `json:"key"`
	Value        string          `json:"value"`
	Type         SettingType     `json:"type"`
	Category     SettingCategory `json:"category"`
	Label        string          `json:"label"`
	Description  string          `json:"description"`
	IsSecret     bool            `json:"is_secret"`
	IsRequired   bool            `json:"is_required"`
	HasValue     bool            `json:"has_value"`
	DisplayOrder int             `json:"display_order"`
}

// ToPublic converts a Setting to SettingPublic (masks secret values)
func (s *Setting) ToPublic() *SettingPublic {
	value := ""
	hasValue := false

	if s.Value != nil && *s.Value != "" {
		hasValue = true
		if s.IsSecret {
			// Mask secret values
			value = "********"
		} else {
			value = *s.Value
		}
	}

	description := ""
	if s.Description != nil {
		description = *s.Description
	}

	return &SettingPublic{
		ID:           s.ID,
		Key:          s.Key,
		Value:        value,
		Type:         s.Type,
		Category:     s.Category,
		Label:        s.Label,
		Description:  description,
		IsSecret:     s.IsSecret,
		IsRequired:   s.IsRequired,
		HasValue:     hasValue,
		DisplayOrder: s.DisplayOrder,
	}
}

// SettingsByCategory represents settings grouped by category
type SettingsByCategory struct {
	Category string           `json:"category"`
	Label    string           `json:"label"`
	Settings []*SettingPublic `json:"settings"`
}

// UpdateSettingRequest represents a request to update a setting
type UpdateSettingRequest struct {
	Value string `json:"value"`
}

// BulkUpdateSettingsRequest represents a request to update multiple settings
type BulkUpdateSettingsRequest struct {
	Settings map[string]string `json:"settings"` // key -> value
}

// CategoryLabels maps category keys to display labels
var CategoryLabels = map[SettingCategory]string{
	CategoryGeneral: "Geral",
	CategoryPayment: "Pagamentos (Asaas)",
	CategoryAI:      "Inteligência Artificial",
	CategoryRevenue: "Divisão de Receita",
	CategoryStorage: "Armazenamento (MinIO)",
	CategoryEmail:   "Email (SMTP)",
}
