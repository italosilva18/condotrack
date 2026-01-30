package inspection

import (
	"context"
	"errors"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/google/uuid"
)

// UseCase defines the inspection use case interface
type UseCase interface {
	ListInspections(ctx context.Context, filter *entity.InspectionFilter) ([]entity.Inspection, error)
	GetInspectionByID(ctx context.Context, id string) (*entity.Inspection, error)
	CreateInspection(ctx context.Context, req *entity.CreateInspectionRequest) (*entity.Inspection, error)
	UpdateInspection(ctx context.Context, id string, req *entity.UpdateInspectionRequest) (*entity.Inspection, error)
	DeleteInspection(ctx context.Context, id string) error
	GetInspectionsByContract(ctx context.Context, contractID string) ([]entity.Inspection, error)
	GetInspectionsByInspector(ctx context.Context, inspectorID string) ([]entity.Inspection, error)
	GetScheduledInspections(ctx context.Context) ([]entity.Inspection, error)
}

type inspectionUseCase struct {
	repo         repository.InspectionRepository
	contratoRepo repository.ContratoRepository
	gestorRepo   repository.GestorRepository
}

// NewUseCase creates a new inspection use case
func NewUseCase(
	repo repository.InspectionRepository,
	contratoRepo repository.ContratoRepository,
	gestorRepo repository.GestorRepository,
) UseCase {
	return &inspectionUseCase{
		repo:         repo,
		contratoRepo: contratoRepo,
		gestorRepo:   gestorRepo,
	}
}

// ListInspections returns inspections based on filters
func (uc *inspectionUseCase) ListInspections(ctx context.Context, filter *entity.InspectionFilter) ([]entity.Inspection, error) {
	if filter == nil || (filter.ContractID == "" && filter.InspectorID == "" && filter.Status == "" && filter.StartDate == nil && filter.EndDate == nil) {
		return uc.repo.FindAll(ctx)
	}
	return uc.repo.FindAllWithFilters(ctx, filter)
}

// GetInspectionByID returns a specific inspection by ID
func (uc *inspectionUseCase) GetInspectionByID(ctx context.Context, id string) (*entity.Inspection, error) {
	return uc.repo.FindByID(ctx, id)
}

// CreateInspection creates a new inspection
func (uc *inspectionUseCase) CreateInspection(ctx context.Context, req *entity.CreateInspectionRequest) (*entity.Inspection, error) {
	// Validate inspection type
	if !entity.IsValidInspectionType(req.InspectionType) {
		return nil, errors.New("invalid inspection type")
	}

	// Verify contract exists
	contrato, err := uc.contratoRepo.FindByID(ctx, req.ContractID)
	if err != nil {
		return nil, err
	}
	if contrato == nil {
		return nil, errors.New("contract not found")
	}

	// Verify inspector (gestor) exists
	inspector, err := uc.gestorRepo.FindByID(ctx, req.InspectorID)
	if err != nil {
		return nil, err
	}
	if inspector == nil {
		return nil, errors.New("inspector not found")
	}

	// Set inspection date
	inspectionDate := time.Now()
	if req.InspectionDate != nil {
		inspectionDate = *req.InspectionDate
	}

	// Set default status
	status := entity.InspectionStatusScheduled
	if req.Status != "" {
		if !entity.IsValidInspectionStatus(req.Status) {
			return nil, errors.New("invalid inspection status")
		}
		status = req.Status
	}

	// Create inspection entity
	inspection := &entity.Inspection{
		ID:              uuid.New().String(),
		ContractID:      req.ContractID,
		ContractName:    contrato.Nome,
		InspectorID:     req.InspectorID,
		InspectorName:   inspector.Nome,
		InspectionDate:  inspectionDate,
		InspectionType:  req.InspectionType,
		Status:          status,
		Findings:        req.Findings,
		Recommendations: req.Recommendations,
		Photos:          req.Photos,
		CreatedAt:       time.Now(),
	}

	if err := uc.repo.Create(ctx, inspection); err != nil {
		return nil, err
	}

	return inspection, nil
}

// UpdateInspection updates an existing inspection
func (uc *inspectionUseCase) UpdateInspection(ctx context.Context, id string, req *entity.UpdateInspectionRequest) (*entity.Inspection, error) {
	// Verify inspection exists
	inspection, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if inspection == nil {
		return nil, errors.New("inspection not found")
	}

	// Update fields if provided
	if req.ContractID != nil {
		// Verify new contract exists
		contrato, err := uc.contratoRepo.FindByID(ctx, *req.ContractID)
		if err != nil {
			return nil, err
		}
		if contrato == nil {
			return nil, errors.New("contract not found")
		}
		inspection.ContractID = *req.ContractID
		inspection.ContractName = contrato.Nome
	}

	if req.InspectorID != nil {
		// Verify new inspector exists
		inspector, err := uc.gestorRepo.FindByID(ctx, *req.InspectorID)
		if err != nil {
			return nil, err
		}
		if inspector == nil {
			return nil, errors.New("inspector not found")
		}
		inspection.InspectorID = *req.InspectorID
		inspection.InspectorName = inspector.Nome
	}

	if req.InspectionDate != nil {
		inspection.InspectionDate = *req.InspectionDate
	}

	if req.InspectionType != nil {
		if !entity.IsValidInspectionType(*req.InspectionType) {
			return nil, errors.New("invalid inspection type")
		}
		inspection.InspectionType = *req.InspectionType
	}

	if req.Status != nil {
		if !entity.IsValidInspectionStatus(*req.Status) {
			return nil, errors.New("invalid inspection status")
		}
		inspection.Status = *req.Status
	}

	if req.Findings != nil {
		inspection.Findings = req.Findings
	}

	if req.Recommendations != nil {
		inspection.Recommendations = req.Recommendations
	}

	if req.Photos != nil {
		inspection.Photos = req.Photos
	}

	// Set updated timestamp
	now := time.Now()
	inspection.UpdatedAt = &now

	if err := uc.repo.Update(ctx, inspection); err != nil {
		return nil, err
	}

	return inspection, nil
}

// DeleteInspection deletes an inspection by ID
func (uc *inspectionUseCase) DeleteInspection(ctx context.Context, id string) error {
	// Verify inspection exists
	inspection, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if inspection == nil {
		return errors.New("inspection not found")
	}

	return uc.repo.Delete(ctx, id)
}

// GetInspectionsByContract returns all inspections for a specific contract
func (uc *inspectionUseCase) GetInspectionsByContract(ctx context.Context, contractID string) ([]entity.Inspection, error) {
	// Verify contract exists
	contrato, err := uc.contratoRepo.FindByID(ctx, contractID)
	if err != nil {
		return nil, err
	}
	if contrato == nil {
		return nil, errors.New("contract not found")
	}

	return uc.repo.FindByContract(ctx, contractID)
}

// GetInspectionsByInspector returns all inspections for a specific inspector
func (uc *inspectionUseCase) GetInspectionsByInspector(ctx context.Context, inspectorID string) ([]entity.Inspection, error) {
	// Verify inspector exists
	inspector, err := uc.gestorRepo.FindByID(ctx, inspectorID)
	if err != nil {
		return nil, err
	}
	if inspector == nil {
		return nil, errors.New("inspector not found")
	}

	return uc.repo.FindByInspector(ctx, inspectorID)
}

// GetScheduledInspections returns all scheduled (upcoming) inspections
func (uc *inspectionUseCase) GetScheduledInspections(ctx context.Context) ([]entity.Inspection, error) {
	return uc.repo.FindScheduled(ctx)
}
