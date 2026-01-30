package repository

import (
	"context"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
)

// AgendaRepository defines the interface for agenda/calendar data access
type AgendaRepository interface {
	// FindAll returns all events
	FindAll(ctx context.Context) ([]entity.AgendaEvent, error)

	// FindByID returns an event by ID
	FindByID(ctx context.Context, id string) (*entity.AgendaEvent, error)

	// FindByDateRange returns events within a date range
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entity.AgendaEvent, error)

	// FindByContract returns all events for a specific contract
	FindByContract(ctx context.Context, contractID string) ([]entity.AgendaEvent, error)

	// FindByUser returns all events for a specific user
	FindByUser(ctx context.Context, userID string) ([]entity.AgendaEvent, error)

	// FindWithFilters returns events matching the given filters
	FindWithFilters(ctx context.Context, filter *entity.AgendaFilter) ([]entity.AgendaEvent, error)

	// Create creates a new event
	Create(ctx context.Context, event *entity.AgendaEvent) error

	// Update updates an existing event
	Update(ctx context.Context, event *entity.AgendaEvent) error

	// Delete deletes an event by ID
	Delete(ctx context.Context, id string) error
}
