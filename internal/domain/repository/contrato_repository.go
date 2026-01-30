package repository

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
)

// ContratoRepository defines the interface for contrato data access
type ContratoRepository interface {
	// FindAll returns all contratos
	FindAll(ctx context.Context) ([]entity.Contrato, error)

	// FindByID returns a contrato by ID
	FindByID(ctx context.Context, id string) (*entity.Contrato, error)

	// FindByGestorID returns all contratos for a specific gestor
	FindByGestorID(ctx context.Context, gestorID string) ([]entity.Contrato, error)

	// FindAllWithGestor returns all contratos with gestor information
	FindAllWithGestor(ctx context.Context) ([]entity.ContratoWithGestor, error)

	// Create creates a new contrato
	Create(ctx context.Context, contrato *entity.Contrato) error

	// Update updates an existing contrato
	Update(ctx context.Context, contrato *entity.Contrato) error

	// Delete deletes a contrato by ID
	Delete(ctx context.Context, id string) error

	// CountByGestorID returns the number of contracts for a gestor
	CountByGestorID(ctx context.Context, gestorID string) (int, error)
}
