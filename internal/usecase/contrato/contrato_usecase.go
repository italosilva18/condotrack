package contrato

import (
	"context"
	"errors"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/google/uuid"
)

// UseCase defines the contrato use case interface
type UseCase interface {
	ListContratos(ctx context.Context) ([]entity.Contrato, error)
	ListContratosWithGestor(ctx context.Context) ([]entity.ContratoWithGestor, error)
	ListContratosByGestor(ctx context.Context, gestorID string) ([]entity.Contrato, error)
	GetContratoByID(ctx context.Context, id string) (*entity.Contrato, error)
	CreateContrato(ctx context.Context, req *entity.CreateContratoRequest) (*entity.Contrato, error)
	UpdateContrato(ctx context.Context, id string, req *entity.UpdateContratoRequest) (*entity.Contrato, error)
	DeleteContrato(ctx context.Context, id string) error
}

type contratoUseCase struct {
	repo       repository.ContratoRepository
	gestorRepo repository.GestorRepository
}

// NewUseCase creates a new contrato use case
func NewUseCase(repo repository.ContratoRepository, gestorRepo repository.GestorRepository) UseCase {
	return &contratoUseCase{
		repo:       repo,
		gestorRepo: gestorRepo,
	}
}

// ListContratos returns all active contratos
func (uc *contratoUseCase) ListContratos(ctx context.Context) ([]entity.Contrato, error) {
	return uc.repo.FindAll(ctx)
}

// ListContratosWithGestor returns all contratos with gestor information
func (uc *contratoUseCase) ListContratosWithGestor(ctx context.Context) ([]entity.ContratoWithGestor, error) {
	return uc.repo.FindAllWithGestor(ctx)
}

// ListContratosByGestor returns all contratos for a specific gestor
func (uc *contratoUseCase) ListContratosByGestor(ctx context.Context, gestorID string) ([]entity.Contrato, error) {
	// Verify gestor exists
	gestor, err := uc.gestorRepo.FindByID(ctx, gestorID)
	if err != nil {
		return nil, err
	}
	if gestor == nil {
		return nil, errors.New("gestor not found")
	}

	return uc.repo.FindByGestorID(ctx, gestorID)
}

// GetContratoByID returns a specific contrato by ID
func (uc *contratoUseCase) GetContratoByID(ctx context.Context, id string) (*entity.Contrato, error) {
	return uc.repo.FindByID(ctx, id)
}

// CreateContrato creates a new contrato
func (uc *contratoUseCase) CreateContrato(ctx context.Context, req *entity.CreateContratoRequest) (*entity.Contrato, error) {
	// Verify gestor exists
	gestor, err := uc.gestorRepo.FindByID(ctx, req.GestorID)
	if err != nil {
		return nil, err
	}
	if gestor == nil {
		return nil, errors.New("gestor not found")
	}

	// Create contrato entity
	contrato := &entity.Contrato{
		ID:            uuid.New().String(),
		GestorID:      req.GestorID,
		Nome:          req.Nome,
		Descricao:     req.Descricao,
		Endereco:      req.Endereco,
		Cidade:        req.Cidade,
		Estado:        req.Estado,
		CEP:           req.CEP,
		TotalUnidades: req.TotalUnidades,
		MetaScore:     req.MetaScore,
		Ativo:         true,
		CreatedAt:     time.Now(),
	}

	// Set default meta score if not provided
	if contrato.MetaScore == 0 {
		contrato.MetaScore = 80.0
	}

	if err := uc.repo.Create(ctx, contrato); err != nil {
		return nil, err
	}

	return contrato, nil
}

// UpdateContrato updates an existing contrato
func (uc *contratoUseCase) UpdateContrato(ctx context.Context, id string, req *entity.UpdateContratoRequest) (*entity.Contrato, error) {
	// Verify contrato exists
	contrato, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if contrato == nil {
		return nil, errors.New("contrato not found")
	}

	// If gestor_id is being updated, verify new gestor exists
	if req.GestorID != nil {
		gestor, err := uc.gestorRepo.FindByID(ctx, *req.GestorID)
		if err != nil {
			return nil, err
		}
		if gestor == nil {
			return nil, errors.New("gestor not found")
		}
		contrato.GestorID = *req.GestorID
	}

	// Update fields if provided
	if req.Nome != nil {
		contrato.Nome = *req.Nome
	}
	if req.Descricao != nil {
		contrato.Descricao = req.Descricao
	}
	if req.Endereco != nil {
		contrato.Endereco = req.Endereco
	}
	if req.Cidade != nil {
		contrato.Cidade = req.Cidade
	}
	if req.Estado != nil {
		contrato.Estado = req.Estado
	}
	if req.CEP != nil {
		contrato.CEP = req.CEP
	}
	if req.TotalUnidades != nil {
		contrato.TotalUnidades = *req.TotalUnidades
	}
	if req.MetaScore != nil {
		contrato.MetaScore = *req.MetaScore
	}
	if req.Ativo != nil {
		contrato.Ativo = *req.Ativo
	}

	// Set updated timestamp
	now := time.Now()
	contrato.UpdatedAt = &now

	if err := uc.repo.Update(ctx, contrato); err != nil {
		return nil, err
	}

	return contrato, nil
}

// DeleteContrato deletes a contrato by ID
func (uc *contratoUseCase) DeleteContrato(ctx context.Context, id string) error {
	// Verify contrato exists
	contrato, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if contrato == nil {
		return errors.New("contrato not found")
	}

	return uc.repo.Delete(ctx, id)
}
