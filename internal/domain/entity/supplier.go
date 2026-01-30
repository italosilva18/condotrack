package entity

import "time"

// Supplier represents a supplier/vendor entity
type Supplier struct {
	ID        string     `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	CNPJ      *string    `db:"cnpj" json:"cnpj,omitempty"`
	Email     *string    `db:"email" json:"email,omitempty"`
	Phone     *string    `db:"phone" json:"phone,omitempty"`
	Address   *string    `db:"address" json:"address,omitempty"`
	Category  *string    `db:"category" json:"category,omitempty"`
	IsActive  bool       `db:"is_active" json:"is_active"`
	Notes     *string    `db:"notes" json:"notes,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// CreateSupplierRequest represents the request to create a supplier
type CreateSupplierRequest struct {
	Name     string  `json:"name" binding:"required"`
	CNPJ     *string `json:"cnpj"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Address  *string `json:"address"`
	Category *string `json:"category"`
	Notes    *string `json:"notes"`
}

// UpdateSupplierRequest represents the request to update a supplier
type UpdateSupplierRequest struct {
	Name     *string `json:"name"`
	CNPJ     *string `json:"cnpj"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Address  *string `json:"address"`
	Category *string `json:"category"`
	IsActive *bool   `json:"is_active"`
	Notes    *string `json:"notes"`
}
