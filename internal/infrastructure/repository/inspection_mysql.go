package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type inspectionMySQLRepository struct {
	db *sqlx.DB
}

// NewInspectionMySQLRepository creates a new MySQL implementation of InspectionRepository
func NewInspectionMySQLRepository(db *sqlx.DB) repository.InspectionRepository {
	return &inspectionMySQLRepository{db: db}
}

func (r *inspectionMySQLRepository) FindAll(ctx context.Context) ([]entity.Inspection, error) {
	var inspections []entity.Inspection
	query := `SELECT i.id, i.contract_id, c.nome as contract_name,
			  i.inspector_id, COALESCE(u.name, '') as inspector_name,
			  i.inspection_date, i.inspection_type, i.status,
			  i.findings, i.recommendations, COALESCE(i.photos, '[]') as photos,
			  i.created_at, i.updated_at
			  FROM inspections i
			  INNER JOIN contratos c ON c.id = i.contract_id
			  LEFT JOIN users u ON u.id = i.inspector_id
			  ORDER BY i.inspection_date DESC`
	err := r.db.SelectContext(ctx, &inspections, query)
	if err != nil {
		return nil, err
	}
	return inspections, nil
}

func (r *inspectionMySQLRepository) FindAllWithFilters(ctx context.Context, filter *entity.InspectionFilter) ([]entity.Inspection, error) {
	var inspections []entity.Inspection
	query := `SELECT i.id, i.contract_id, c.nome as contract_name,
			  i.inspector_id, COALESCE(u.name, '') as inspector_name,
			  i.inspection_date, i.inspection_type, i.status,
			  i.findings, i.recommendations, COALESCE(i.photos, '[]') as photos,
			  i.created_at, i.updated_at
			  FROM inspections i
			  INNER JOIN contratos c ON c.id = i.contract_id
			  LEFT JOIN users u ON u.id = i.inspector_id
			  WHERE 1=1`

	args := []interface{}{}

	if filter != nil {
		if filter.ContractID != "" {
			query += " AND i.contract_id = ?"
			args = append(args, filter.ContractID)
		}
		if filter.InspectorID != "" {
			query += " AND i.inspector_id = ?"
			args = append(args, filter.InspectorID)
		}
		if filter.Status != "" {
			query += " AND i.status = ?"
			args = append(args, filter.Status)
		}
		if filter.StartDate != nil {
			query += " AND i.inspection_date >= ?"
			args = append(args, filter.StartDate)
		}
		if filter.EndDate != nil {
			query += " AND i.inspection_date <= ?"
			args = append(args, filter.EndDate)
		}
	}

	query += " ORDER BY i.inspection_date DESC"

	err := r.db.SelectContext(ctx, &inspections, query, args...)
	if err != nil {
		return nil, err
	}
	return inspections, nil
}

func (r *inspectionMySQLRepository) FindByID(ctx context.Context, id string) (*entity.Inspection, error) {
	var inspection entity.Inspection
	query := `SELECT i.id, i.contract_id, c.nome as contract_name,
			  i.inspector_id, COALESCE(u.name, '') as inspector_name,
			  i.inspection_date, i.inspection_type, i.status,
			  i.findings, i.recommendations, COALESCE(i.photos, '[]') as photos,
			  i.created_at, i.updated_at
			  FROM inspections i
			  INNER JOIN contratos c ON c.id = i.contract_id
			  LEFT JOIN users u ON u.id = i.inspector_id
			  WHERE i.id = ?`
	err := r.db.GetContext(ctx, &inspection, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &inspection, nil
}

func (r *inspectionMySQLRepository) FindByContract(ctx context.Context, contractID string) ([]entity.Inspection, error) {
	var inspections []entity.Inspection
	query := `SELECT i.id, i.contract_id, c.nome as contract_name,
			  i.inspector_id, COALESCE(u.name, '') as inspector_name,
			  i.inspection_date, i.inspection_type, i.status,
			  i.findings, i.recommendations, COALESCE(i.photos, '[]') as photos,
			  i.created_at, i.updated_at
			  FROM inspections i
			  INNER JOIN contratos c ON c.id = i.contract_id
			  LEFT JOIN users u ON u.id = i.inspector_id
			  WHERE i.contract_id = ?
			  ORDER BY i.inspection_date DESC`
	err := r.db.SelectContext(ctx, &inspections, query, contractID)
	if err != nil {
		return nil, err
	}
	return inspections, nil
}

func (r *inspectionMySQLRepository) FindByInspector(ctx context.Context, inspectorID string) ([]entity.Inspection, error) {
	var inspections []entity.Inspection
	query := `SELECT i.id, i.contract_id, c.nome as contract_name,
			  i.inspector_id, COALESCE(u.name, '') as inspector_name,
			  i.inspection_date, i.inspection_type, i.status,
			  i.findings, i.recommendations, COALESCE(i.photos, '[]') as photos,
			  i.created_at, i.updated_at
			  FROM inspections i
			  INNER JOIN contratos c ON c.id = i.contract_id
			  LEFT JOIN users u ON u.id = i.inspector_id
			  WHERE i.inspector_id = ?
			  ORDER BY i.inspection_date DESC`
	err := r.db.SelectContext(ctx, &inspections, query, inspectorID)
	if err != nil {
		return nil, err
	}
	return inspections, nil
}

func (r *inspectionMySQLRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entity.Inspection, error) {
	var inspections []entity.Inspection
	query := `SELECT i.id, i.contract_id, c.nome as contract_name,
			  i.inspector_id, COALESCE(u.name, '') as inspector_name,
			  i.inspection_date, i.inspection_type, i.status,
			  i.findings, i.recommendations, COALESCE(i.photos, '[]') as photos,
			  i.created_at, i.updated_at
			  FROM inspections i
			  INNER JOIN contratos c ON c.id = i.contract_id
			  LEFT JOIN users u ON u.id = i.inspector_id
			  WHERE i.inspection_date >= ? AND i.inspection_date <= ?
			  ORDER BY i.inspection_date DESC`
	err := r.db.SelectContext(ctx, &inspections, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return inspections, nil
}

func (r *inspectionMySQLRepository) FindByStatus(ctx context.Context, status string) ([]entity.Inspection, error) {
	var inspections []entity.Inspection
	query := `SELECT i.id, i.contract_id, c.nome as contract_name,
			  i.inspector_id, COALESCE(u.name, '') as inspector_name,
			  i.inspection_date, i.inspection_type, i.status,
			  i.findings, i.recommendations, COALESCE(i.photos, '[]') as photos,
			  i.created_at, i.updated_at
			  FROM inspections i
			  INNER JOIN contratos c ON c.id = i.contract_id
			  LEFT JOIN users u ON u.id = i.inspector_id
			  WHERE i.status = ?
			  ORDER BY i.inspection_date DESC`
	err := r.db.SelectContext(ctx, &inspections, query, status)
	if err != nil {
		return nil, err
	}
	return inspections, nil
}

func (r *inspectionMySQLRepository) FindScheduled(ctx context.Context) ([]entity.Inspection, error) {
	var inspections []entity.Inspection
	query := `SELECT i.id, i.contract_id, c.nome as contract_name,
			  i.inspector_id, COALESCE(u.name, '') as inspector_name,
			  i.inspection_date, i.inspection_type, i.status,
			  i.findings, i.recommendations, COALESCE(i.photos, '[]') as photos,
			  i.created_at, i.updated_at
			  FROM inspections i
			  INNER JOIN contratos c ON c.id = i.contract_id
			  LEFT JOIN users u ON u.id = i.inspector_id
			  WHERE i.status = ? AND i.inspection_date >= NOW()
			  ORDER BY i.inspection_date ASC`
	err := r.db.SelectContext(ctx, &inspections, query, entity.InspectionStatusScheduled)
	if err != nil {
		return nil, err
	}
	return inspections, nil
}

func (r *inspectionMySQLRepository) Create(ctx context.Context, inspection *entity.Inspection) error {
	query := `INSERT INTO inspections (id, contract_id, inspector_id, inspection_date,
			  inspection_type, status, findings, recommendations, photos, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query,
		inspection.ID, inspection.ContractID, inspection.InspectorID, inspection.InspectionDate,
		inspection.InspectionType, inspection.Status, inspection.Findings, inspection.Recommendations,
		inspection.Photos)
	return err
}

func (r *inspectionMySQLRepository) Update(ctx context.Context, inspection *entity.Inspection) error {
	query := `UPDATE inspections
			  SET contract_id = ?, inspector_id = ?, inspection_date = ?,
			  inspection_type = ?, status = ?, findings = ?, recommendations = ?,
			  photos = ?, updated_at = NOW()
			  WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		inspection.ContractID, inspection.InspectorID, inspection.InspectionDate,
		inspection.InspectionType, inspection.Status, inspection.Findings, inspection.Recommendations,
		inspection.Photos, inspection.ID)
	return err
}

func (r *inspectionMySQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM inspections WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *inspectionMySQLRepository) CountByContractID(ctx context.Context, contractID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM inspections WHERE contract_id = ?`
	err := r.db.GetContext(ctx, &count, query, contractID)
	return count, err
}

func (r *inspectionMySQLRepository) CountByInspectorID(ctx context.Context, inspectorID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM inspections WHERE inspector_id = ?`
	err := r.db.GetContext(ctx, &count, query, inspectorID)
	return count, err
}
