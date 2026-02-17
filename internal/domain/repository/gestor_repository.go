package repository

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
)

// GestorRepository defines the interface for gestor data access
type GestorRepository interface {
	// FindAll returns all gestores
	FindAll(ctx context.Context) ([]entity.Gestor, error)

	// FindByID returns a gestor by ID
	FindByID(ctx context.Context, id string) (*entity.Gestor, error)

	// FindByEmail returns a gestor by email
	FindByEmail(ctx context.Context, email string) (*entity.Gestor, error)

	// FindAllWithContracts returns all gestores with their contract counts
	FindAllWithContracts(ctx context.Context) ([]entity.GestorWithContracts, error)

	// Create creates a new gestor
	Create(ctx context.Context, gestor *entity.Gestor) error

	// Update updates an existing gestor
	Update(ctx context.Context, gestor *entity.Gestor) error

	// Delete deletes a gestor by ID
	Delete(ctx context.Context, id string) error
}
