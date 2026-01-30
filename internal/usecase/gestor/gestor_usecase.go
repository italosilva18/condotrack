package gestor

import (
	"context"
	"errors"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/google/uuid"
)

// CreateGestorRequest represents the request to create a gestor
type CreateGestorRequest struct {
	Nome     string  `json:"nome" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Telefone *string `json:"telefone"`
	CPF      *string `json:"cpf"`
}

// UpdateGestorRequest represents the request to update a gestor
type UpdateGestorRequest struct {
	Nome     *string `json:"nome"`
	Email    *string `json:"email"`
	Telefone *string `json:"telefone"`
	CPF      *string `json:"cpf"`
	Ativo    *bool   `json:"ativo"`
}

// UseCase defines the gestor use case interface
type UseCase interface {
	ListGestores(ctx context.Context) ([]entity.Gestor, error)
	ListGestoresWithContracts(ctx context.Context) ([]entity.GestorWithContracts, error)
	GetGestorByID(ctx context.Context, id string) (*entity.Gestor, error)
	CreateGestor(ctx context.Context, req *CreateGestorRequest) (*entity.Gestor, error)
	UpdateGestor(ctx context.Context, id string, req *UpdateGestorRequest) (*entity.Gestor, error)
	DeleteGestor(ctx context.Context, id string) error
}

type gestorUseCase struct {
	repo repository.GestorRepository
}

// NewUseCase creates a new gestor use case
func NewUseCase(repo repository.GestorRepository) UseCase {
	return &gestorUseCase{repo: repo}
}

// ListGestores returns all active gestores
func (uc *gestorUseCase) ListGestores(ctx context.Context) ([]entity.Gestor, error) {
	return uc.repo.FindAll(ctx)
}

// ListGestoresWithContracts returns all gestores with their contract counts
func (uc *gestorUseCase) ListGestoresWithContracts(ctx context.Context) ([]entity.GestorWithContracts, error) {
	return uc.repo.FindAllWithContracts(ctx)
}

// GetGestorByID returns a specific gestor by ID
func (uc *gestorUseCase) GetGestorByID(ctx context.Context, id string) (*entity.Gestor, error) {
	return uc.repo.FindByID(ctx, id)
}

// CreateGestor creates a new gestor
func (uc *gestorUseCase) CreateGestor(ctx context.Context, req *CreateGestorRequest) (*entity.Gestor, error) {
	// Check if email already exists
	existing, err := uc.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("gestor with this email already exists")
	}

	// Create gestor entity
	gestor := &entity.Gestor{
		ID:        uuid.New().String(),
		Nome:      req.Nome,
		Email:     req.Email,
		Telefone:  req.Telefone,
		CPF:       req.CPF,
		Ativo:     true,
		CreatedAt: time.Now(),
	}

	if err := uc.repo.Create(ctx, gestor); err != nil {
		return nil, err
	}

	return gestor, nil
}

// UpdateGestor updates an existing gestor
func (uc *gestorUseCase) UpdateGestor(ctx context.Context, id string, req *UpdateGestorRequest) (*entity.Gestor, error) {
	// Find existing gestor
	gestor, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if gestor == nil {
		return nil, errors.New("gestor not found")
	}

	// Check if email is being changed and if new email already exists
	if req.Email != nil && *req.Email != gestor.Email {
		existing, err := uc.repo.FindByEmail(ctx, *req.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, errors.New("gestor with this email already exists")
		}
		gestor.Email = *req.Email
	}

	// Update fields if provided
	if req.Nome != nil {
		gestor.Nome = *req.Nome
	}
	if req.Telefone != nil {
		gestor.Telefone = req.Telefone
	}
	if req.CPF != nil {
		gestor.CPF = req.CPF
	}
	if req.Ativo != nil {
		gestor.Ativo = *req.Ativo
	}

	// Set updated timestamp
	now := time.Now()
	gestor.UpdatedAt = &now

	if err := uc.repo.Update(ctx, gestor); err != nil {
		return nil, err
	}

	return gestor, nil
}

// DeleteGestor soft deletes a gestor by setting ativo to false
func (uc *gestorUseCase) DeleteGestor(ctx context.Context, id string) error {
	// Find existing gestor
	gestor, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if gestor == nil {
		return errors.New("gestor not found")
	}

	// Soft delete by calling repository Delete method
	return uc.repo.Delete(ctx, id)
}
