package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type auditMySQLRepository struct {
	db *sqlx.DB
}

// NewAuditMySQLRepository creates a new MySQL implementation of AuditRepository
func NewAuditMySQLRepository(db *sqlx.DB) repository.AuditRepository {
	return &auditMySQLRepository{db: db}
}

func (r *auditMySQLRepository) FindAll(ctx context.Context) ([]entity.Audit, error) {
	var audits []entity.Audit
	query := `SELECT id, contract_id, auditor_name, audit_date, score, target_score,
			  previous_score, status, observations, COALESCE(data_json, '{}') as data_json, created_at, updated_at
			  FROM audits
			  ORDER BY audit_date DESC`
	err := r.db.SelectContext(ctx, &audits, query)
	if err != nil {
		return nil, err
	}
	return audits, nil
}

func (r *auditMySQLRepository) FindByID(ctx context.Context, id string) (*entity.Audit, error) {
	var audit entity.Audit
	query := `SELECT id, contract_id, auditor_name, audit_date, score, target_score,
			  previous_score, status, observations, COALESCE(data_json, '{}') as data_json, created_at, updated_at
			  FROM audits
			  WHERE id = ?`
	err := r.db.GetContext(ctx, &audit, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &audit, nil
}

func (r *auditMySQLRepository) FindByContractID(ctx context.Context, contractID string) ([]entity.Audit, error) {
	var audits []entity.Audit
	query := `SELECT id, contract_id, auditor_name, audit_date, score, target_score,
			  previous_score, status, observations, COALESCE(data_json, '{}') as data_json, created_at, updated_at
			  FROM audits
			  WHERE contract_id = ?
			  ORDER BY audit_date DESC`
	err := r.db.SelectContext(ctx, &audits, query, contractID)
	if err != nil {
		return nil, err
	}
	return audits, nil
}

func (r *auditMySQLRepository) FindAllWithContract(ctx context.Context) ([]entity.AuditWithContract, error) {
	var audits []entity.AuditWithContract
	query := `SELECT a.id, a.contract_id, a.auditor_name, a.audit_date, a.score, a.target_score,
			  a.previous_score, a.status, a.observations, COALESCE(a.data_json, '{}') as data_json, a.created_at, a.updated_at,
			  c.nome as contract_name, g.nome as gestor_name
			  FROM audits a
			  INNER JOIN contratos c ON c.id = a.contract_id
			  INNER JOIN gestores g ON g.id = c.gestor_id
			  ORDER BY a.audit_date DESC`
	err := r.db.SelectContext(ctx, &audits, query)
	if err != nil {
		return nil, err
	}
	return audits, nil
}

func (r *auditMySQLRepository) FindLastByContractID(ctx context.Context, contractID string) (*entity.Audit, error) {
	var audit entity.Audit
	query := `SELECT id, contract_id, auditor_name, audit_date, score, target_score,
			  previous_score, status, observations, COALESCE(data_json, '{}') as data_json, created_at, updated_at
			  FROM audits
			  WHERE contract_id = ?
			  ORDER BY audit_date DESC
			  LIMIT 1`
	err := r.db.GetContext(ctx, &audit, query, contractID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &audit, nil
}

func (r *auditMySQLRepository) Create(ctx context.Context, audit *entity.Audit) error {
	query := `INSERT INTO audits (id, contract_id, auditor_name, audit_date, score, target_score,
			  previous_score, status, observations, data_json, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query,
		audit.ID, audit.ContractID, audit.AuditorName, audit.AuditDate, audit.Score,
		audit.TargetScore, audit.PreviousScore, audit.Status, audit.Observations, audit.DataJSON)
	return err
}

func (r *auditMySQLRepository) CreateWithTx(ctx context.Context, tx *sqlx.Tx, audit *entity.Audit) error {
	query := `INSERT INTO audits (id, contract_id, auditor_name, audit_date, score, target_score,
			  previous_score, status, observations, data_json, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := tx.ExecContext(ctx, query,
		audit.ID, audit.ContractID, audit.AuditorName, audit.AuditDate, audit.Score,
		audit.TargetScore, audit.PreviousScore, audit.Status, audit.Observations, audit.DataJSON)
	return err
}

func (r *auditMySQLRepository) Update(ctx context.Context, audit *entity.Audit) error {
	query := `UPDATE audits
			  SET contract_id = ?, auditor_name = ?, audit_date = ?, score = ?, target_score = ?,
			  previous_score = ?, status = ?, observations = ?, data_json = ?, updated_at = NOW()
			  WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		audit.ContractID, audit.AuditorName, audit.AuditDate, audit.Score,
		audit.TargetScore, audit.PreviousScore, audit.Status, audit.Observations, audit.DataJSON, audit.ID)
	return err
}

func (r *auditMySQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM audits WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *auditMySQLRepository) GetMeta(ctx context.Context, contractID string) (*entity.AuditMeta, error) {
	meta := &entity.AuditMeta{}

	// Get total and average
	query := `SELECT COUNT(*) as total, COALESCE(AVG(score), 0) as avg_score
			  FROM audits WHERE contract_id = ?`
	var stats struct {
		Total    int     `db:"total"`
		AvgScore float64 `db:"avg_score"`
	}
	err := r.db.GetContext(ctx, &stats, query, contractID)
	if err != nil {
		return nil, err
	}
	meta.TotalAudits = stats.Total
	meta.AverageScore = stats.AvgScore

	// Get last score
	lastAudit, err := r.FindLastByContractID(ctx, contractID)
	if err != nil {
		return nil, err
	}
	if lastAudit != nil {
		meta.LastScore = &lastAudit.Score
		meta.TargetScore = lastAudit.TargetScore
	}

	// Get approved/rejected counts
	countQuery := `SELECT
			  SUM(CASE WHEN status = 'approved' THEN 1 ELSE 0 END) as approved,
			  SUM(CASE WHEN status = 'rejected' THEN 1 ELSE 0 END) as rejected
			  FROM audits WHERE contract_id = ?`
	var counts struct {
		Approved int `db:"approved"`
		Rejected int `db:"rejected"`
	}
	err = r.db.GetContext(ctx, &counts, countQuery, contractID)
	if err != nil {
		return nil, err
	}
	meta.ApprovedCount = counts.Approved
	meta.RejectedCount = counts.Rejected

	// Calculate improvement rate if there are at least 2 audits
	if meta.TotalAudits >= 2 {
		var audits []entity.Audit
		auditsQuery := `SELECT score FROM audits WHERE contract_id = ? ORDER BY audit_date ASC LIMIT 2`
		if err := r.db.SelectContext(ctx, &audits, auditsQuery, contractID); err == nil && len(audits) == 2 {
			if audits[0].Score > 0 {
				rate := ((audits[1].Score - audits[0].Score) / audits[0].Score) * 100
				meta.ImprovementRate = &rate
			}
		}
	}

	return meta, nil
}

func (r *auditMySQLRepository) CountByContractID(ctx context.Context, contractID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM audits WHERE contract_id = ?`
	err := r.db.GetContext(ctx, &count, query, contractID)
	return count, err
}

func (r *auditMySQLRepository) GetAverageScoreByContractID(ctx context.Context, contractID string) (float64, error) {
	var avg float64
	query := `SELECT COALESCE(AVG(score), 0) FROM audits WHERE contract_id = ?`
	err := r.db.GetContext(ctx, &avg, query, contractID)
	return avg, err
}

// AuditItemMySQLRepository implementation
type auditItemMySQLRepository struct {
	db *sqlx.DB
}

// NewAuditItemMySQLRepository creates a new MySQL implementation of AuditItemRepository
func NewAuditItemMySQLRepository(db *sqlx.DB) repository.AuditItemRepository {
	return &auditItemMySQLRepository{db: db}
}

func (r *auditItemMySQLRepository) FindByAuditID(ctx context.Context, auditID string) ([]entity.AuditItem, error) {
	var items []entity.AuditItem
	query := `SELECT id, audit_id, category_id, item_name, score, max_score, observation, created_at
			  FROM audit_items
			  WHERE audit_id = ?`
	err := r.db.SelectContext(ctx, &items, query, auditID)
	return items, err
}

func (r *auditItemMySQLRepository) FindByAuditIDWithCategory(ctx context.Context, auditID string) ([]entity.AuditItemWithCategory, error) {
	var items []entity.AuditItemWithCategory
	query := `SELECT ai.id, ai.audit_id, ai.category_id, ai.item_name, ai.score, ai.max_score,
			  ai.observation, ai.created_at, COALESCE(ac.name, 'Geral') as category_name
			  FROM audit_items ai
			  LEFT JOIN audit_categories ac ON ac.id = ai.category_id
			  WHERE ai.audit_id = ?`
	err := r.db.SelectContext(ctx, &items, query, auditID)
	return items, err
}

func (r *auditItemMySQLRepository) Create(ctx context.Context, item *entity.AuditItem) error {
	query := `INSERT INTO audit_items (id, audit_id, category_id, item_name, score, max_score, observation, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query, item.ID, item.AuditID, item.CategoryID, item.ItemName, item.Score, item.MaxScore, item.Observation)
	return err
}

func (r *auditItemMySQLRepository) CreateWithTx(ctx context.Context, tx *sqlx.Tx, item *entity.AuditItem) error {
	query := `INSERT INTO audit_items (id, audit_id, category_id, item_name, score, max_score, observation, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := tx.ExecContext(ctx, query, item.ID, item.AuditID, item.CategoryID, item.ItemName, item.Score, item.MaxScore, item.Observation)
	return err
}

func (r *auditItemMySQLRepository) CreateBatch(ctx context.Context, items []entity.AuditItem) error {
	if len(items) == 0 {
		return nil
	}
	query := `INSERT INTO audit_items (id, audit_id, category_id, item_name, score, max_score, observation, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, NOW())`
	for _, item := range items {
		if _, err := r.db.ExecContext(ctx, query, item.ID, item.AuditID, item.CategoryID, item.ItemName, item.Score, item.MaxScore, item.Observation); err != nil {
			return err
		}
	}
	return nil
}

func (r *auditItemMySQLRepository) CreateBatchWithTx(ctx context.Context, tx *sqlx.Tx, items []entity.AuditItem) error {
	if len(items) == 0 {
		return nil
	}
	query := `INSERT INTO audit_items (id, audit_id, category_id, item_name, score, max_score, observation, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, NOW())`
	for _, item := range items {
		if _, err := tx.ExecContext(ctx, query, item.ID, item.AuditID, item.CategoryID, item.ItemName, item.Score, item.MaxScore, item.Observation); err != nil {
			return err
		}
	}
	return nil
}

func (r *auditItemMySQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM audit_items WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *auditItemMySQLRepository) DeleteByAuditID(ctx context.Context, auditID string) error {
	query := `DELETE FROM audit_items WHERE audit_id = ?`
	_, err := r.db.ExecContext(ctx, query, auditID)
	return err
}
