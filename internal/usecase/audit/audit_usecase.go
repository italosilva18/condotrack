package audit

import (
	"context"
	"errors"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/condotrack/api/internal/infrastructure/database"
	"github.com/google/uuid"
)

const (
	// DefaultTolerance is the default score tolerance for approval (5 points)
	DefaultTolerance = 5.0
)

// UseCase defines the audit use case interface
type UseCase interface {
	ListAudits(ctx context.Context) ([]entity.Audit, error)
	ListAuditsWithContract(ctx context.Context) ([]entity.AuditWithContract, error)
	ListAuditsByContract(ctx context.Context, contractID string) ([]entity.Audit, error)
	GetAuditByID(ctx context.Context, id string) (*entity.Audit, error)
	GetAuditMeta(ctx context.Context, contractID string) (*entity.AuditMeta, error)
	CreateAudit(ctx context.Context, req *entity.CreateAuditRequest) (*entity.Audit, error)
	UpdateAudit(ctx context.Context, id string, req *entity.UpdateAuditRequest) (*entity.Audit, error)
	DeleteAudit(ctx context.Context, id string) error
}

type auditUseCase struct {
	repo         repository.AuditRepository
	itemRepo     repository.AuditItemRepository
	contratoRepo repository.ContratoRepository
	db           *database.MySQL
}

// NewUseCase creates a new audit use case
func NewUseCase(
	repo repository.AuditRepository,
	itemRepo repository.AuditItemRepository,
	contratoRepo repository.ContratoRepository,
	db *database.MySQL,
) UseCase {
	return &auditUseCase{
		repo:         repo,
		itemRepo:     itemRepo,
		contratoRepo: contratoRepo,
		db:           db,
	}
}

// ListAudits returns all audits
func (uc *auditUseCase) ListAudits(ctx context.Context) ([]entity.Audit, error) {
	return uc.repo.FindAll(ctx)
}

// ListAuditsWithContract returns all audits with contract information
func (uc *auditUseCase) ListAuditsWithContract(ctx context.Context) ([]entity.AuditWithContract, error) {
	return uc.repo.FindAllWithContract(ctx)
}

// ListAuditsByContract returns all audits for a specific contract
func (uc *auditUseCase) ListAuditsByContract(ctx context.Context, contractID string) ([]entity.Audit, error) {
	return uc.repo.FindByContractID(ctx, contractID)
}

// GetAuditByID returns a specific audit by ID
func (uc *auditUseCase) GetAuditByID(ctx context.Context, id string) (*entity.Audit, error) {
	return uc.repo.FindByID(ctx, id)
}

// GetAuditMeta returns metadata about audits for a contract
func (uc *auditUseCase) GetAuditMeta(ctx context.Context, contractID string) (*entity.AuditMeta, error) {
	// Verify contract exists
	contrato, err := uc.contratoRepo.FindByID(ctx, contractID)
	if err != nil {
		return nil, err
	}
	if contrato == nil {
		return nil, errors.New("contract not found")
	}

	meta, err := uc.repo.GetMeta(ctx, contractID)
	if err != nil {
		return nil, err
	}

	// Set target score from contract if not available from audits
	if meta.TargetScore == 0 {
		meta.TargetScore = contrato.MetaScore
	}

	return meta, nil
}

// CreateAudit creates a new audit
func (uc *auditUseCase) CreateAudit(ctx context.Context, req *entity.CreateAuditRequest) (*entity.Audit, error) {
	// Verify contract exists
	contrato, err := uc.contratoRepo.FindByID(ctx, req.ContractID)
	if err != nil {
		return nil, err
	}
	if contrato == nil {
		return nil, errors.New("contract not found")
	}

	// Get previous score if exists
	var previousScore *float64
	lastAudit, err := uc.repo.FindLastByContractID(ctx, req.ContractID)
	if err != nil {
		return nil, err
	}
	if lastAudit != nil {
		previousScore = &lastAudit.Score
	}

	// Determine target score
	targetScore := req.TargetScore
	if targetScore == 0 {
		targetScore = contrato.MetaScore
	}

	// Calculate status
	status := entity.CalculateStatus(req.Score, targetScore, DefaultTolerance)

	// Set audit date
	auditDate := time.Now()
	if req.AuditDate != nil {
		auditDate = *req.AuditDate
	}

	// Create audit entity
	audit := &entity.Audit{
		ID:            uuid.New().String(),
		ContractID:    req.ContractID,
		AuditorName:   req.AuditorName,
		AuditDate:     auditDate,
		Score:         req.Score,
		TargetScore:   targetScore,
		PreviousScore: previousScore,
		Status:        status,
		Observations:  req.Observations,
		DataJSON:      req.DataJSON,
		CreatedAt:     time.Now(),
	}

	// Start transaction if there are items
	if len(req.Items) > 0 {
		tx, err := uc.db.BeginTx(ctx)
		if err != nil {
			return nil, err
		}
		defer tx.Rollback()

		// Create audit
		if err := uc.repo.CreateWithTx(ctx, tx, audit); err != nil {
			return nil, err
		}

		// Create audit items
		var items []entity.AuditItem
		for _, itemReq := range req.Items {
			item := entity.AuditItem{
				ID:          uuid.New().String(),
				AuditID:     audit.ID,
				CategoryID:  itemReq.CategoryID,
				ItemName:    itemReq.ItemName,
				Score:       itemReq.Score,
				MaxScore:    itemReq.MaxScore,
				Observation: itemReq.Observation,
				CreatedAt:   time.Now(),
			}
			items = append(items, item)
		}

		if err := uc.itemRepo.CreateBatchWithTx(ctx, tx, items); err != nil {
			return nil, err
		}

		if err := tx.Commit(); err != nil {
			return nil, err
		}
	} else {
		// Create audit without transaction
		if err := uc.repo.Create(ctx, audit); err != nil {
			return nil, err
		}
	}

	return audit, nil
}

// UpdateAudit updates an existing audit
func (uc *auditUseCase) UpdateAudit(ctx context.Context, id string, req *entity.UpdateAuditRequest) (*entity.Audit, error) {
	// Verify audit exists
	audit, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if audit == nil {
		return nil, errors.New("audit not found")
	}

	// Update fields if provided
	if req.AuditorName != nil {
		audit.AuditorName = *req.AuditorName
	}
	if req.AuditDate != nil {
		audit.AuditDate = *req.AuditDate
	}
	if req.Score != nil {
		audit.Score = *req.Score
	}
	if req.TargetScore != nil {
		audit.TargetScore = *req.TargetScore
	}
	if req.Observations != nil {
		audit.Observations = req.Observations
	}
	if req.DataJSON != nil {
		audit.DataJSON = req.DataJSON
	}

	// Recalculate status if score or target changed
	if req.Score != nil || req.TargetScore != nil {
		audit.Status = entity.CalculateStatus(audit.Score, audit.TargetScore, DefaultTolerance)
	}

	// Set updated timestamp
	now := time.Now()
	audit.UpdatedAt = &now

	if err := uc.repo.Update(ctx, audit); err != nil {
		return nil, err
	}

	return audit, nil
}

// DeleteAudit deletes an audit by ID
func (uc *auditUseCase) DeleteAudit(ctx context.Context, id string) error {
	// Verify audit exists
	audit, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if audit == nil {
		return errors.New("audit not found")
	}

	// Delete associated audit items first
	if err := uc.itemRepo.DeleteByAuditID(ctx, id); err != nil {
		return err
	}

	return uc.repo.Delete(ctx, id)
}
