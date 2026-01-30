package repository

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
)

// SupplierRepository defines the interface for supplier data access
type SupplierRepository interface {
	// FindAll returns all suppliers with optional filters
	FindAll(ctx context.Context, category *string, isActive *bool) ([]entity.Supplier, error)

	// FindByID returns a supplier by ID
	FindByID(ctx context.Context, id string) (*entity.Supplier, error)

	// FindByCategory returns all suppliers in a specific category
	FindByCategory(ctx context.Context, category string) ([]entity.Supplier, error)

	// FindActive returns all active suppliers
	FindActive(ctx context.Context) ([]entity.Supplier, error)

	// FindByCNPJ returns a supplier by CNPJ
	FindByCNPJ(ctx context.Context, cnpj string) (*entity.Supplier, error)

	// Create creates a new supplier
	Create(ctx context.Context, supplier *entity.Supplier) error

	// Update updates an existing supplier
	Update(ctx context.Context, supplier *entity.Supplier) error

	// Delete soft deletes a supplier by ID
	Delete(ctx context.Context, id string) error
}
