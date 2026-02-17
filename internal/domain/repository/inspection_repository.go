package repository

import (
	"context"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
)

// InspectionRepository defines the interface for inspection data access
type InspectionRepository interface {
	// FindAll returns all inspections
	FindAll(ctx context.Context) ([]entity.Inspection, error)

	// FindAllWithFilters returns inspections matching the given filters
	FindAllWithFilters(ctx context.Context, filter *entity.InspectionFilter) ([]entity.Inspection, error)

	// FindByID returns an inspection by ID
	FindByID(ctx context.Context, id string) (*entity.Inspection, error)

	// FindByContract returns all inspections for a specific contract
	FindByContract(ctx context.Context, contractID string) ([]entity.Inspection, error)

	// FindByInspector returns all inspections for a specific inspector
	FindByInspector(ctx context.Context, inspectorID string) ([]entity.Inspection, error)

	// FindByDateRange returns all inspections within a date range
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entity.Inspection, error)

	// FindByStatus returns all inspections with a specific status
	FindByStatus(ctx context.Context, status string) ([]entity.Inspection, error)

	// FindScheduled returns all scheduled inspections (future inspections)
	FindScheduled(ctx context.Context) ([]entity.Inspection, error)

	// Create creates a new inspection
	Create(ctx context.Context, inspection *entity.Inspection) error

	// Update updates an existing inspection
	Update(ctx context.Context, inspection *entity.Inspection) error

	// Delete deletes an inspection by ID
	Delete(ctx context.Context, id string) error

	// CountByContractID returns the number of inspections for a contract
	CountByContractID(ctx context.Context, contractID string) (int, error)

	// CountByInspectorID returns the number of inspections for an inspector
	CountByInspectorID(ctx context.Context, inspectorID string) (int, error)
}
