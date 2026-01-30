package repository

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
)

// AuditCategoryRepository defines the interface for audit category data access
type AuditCategoryRepository interface {
	// FindAll returns all audit categories ordered by order_num
	FindAll(ctx context.Context) ([]entity.AuditCategory, error)

	// FindByID returns an audit category by ID
	FindByID(ctx context.Context, id string) (*entity.AuditCategory, error)

	// FindByName returns an audit category by name
	FindByName(ctx context.Context, name string) (*entity.AuditCategory, error)

	// Create creates a new audit category
	Create(ctx context.Context, category *entity.AuditCategory) error

	// Update updates an existing audit category
	Update(ctx context.Context, category *entity.AuditCategory) error

	// Delete deletes an audit category by ID
	Delete(ctx context.Context, id string) error

	// CountItemsByCategory returns the number of audit items using this category
	CountItemsByCategory(ctx context.Context, categoryID string) (int, error)
}
