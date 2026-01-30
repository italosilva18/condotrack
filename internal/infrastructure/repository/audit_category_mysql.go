package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type auditCategoryMySQLRepository struct {
	db *sqlx.DB
}

// NewAuditCategoryMySQLRepository creates a new MySQL implementation of AuditCategoryRepository
func NewAuditCategoryMySQLRepository(db *sqlx.DB) repository.AuditCategoryRepository {
	return &auditCategoryMySQLRepository{db: db}
}

func (r *auditCategoryMySQLRepository) FindAll(ctx context.Context) ([]entity.AuditCategory, error) {
	var categories []entity.AuditCategory
	query := `SELECT id, name, description, weight, order_num
			  FROM audit_categories
			  ORDER BY order_num, name`
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *auditCategoryMySQLRepository) FindByID(ctx context.Context, id string) (*entity.AuditCategory, error) {
	var category entity.AuditCategory
	query := `SELECT id, name, description, weight, order_num
			  FROM audit_categories
			  WHERE id = ?`
	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *auditCategoryMySQLRepository) FindByName(ctx context.Context, name string) (*entity.AuditCategory, error) {
	var category entity.AuditCategory
	query := `SELECT id, name, description, weight, order_num
			  FROM audit_categories
			  WHERE name = ?`
	err := r.db.GetContext(ctx, &category, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *auditCategoryMySQLRepository) Create(ctx context.Context, category *entity.AuditCategory) error {
	query := `INSERT INTO audit_categories (id, name, description, weight, order_num, created_at)
			  VALUES (?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query,
		category.ID, category.Name, category.Description, category.Weight, category.Order)
	return err
}

func (r *auditCategoryMySQLRepository) Update(ctx context.Context, category *entity.AuditCategory) error {
	query := `UPDATE audit_categories
			  SET name = ?, description = ?, weight = ?, order_num = ?
			  WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		category.Name, category.Description, category.Weight, category.Order, category.ID)
	return err
}

func (r *auditCategoryMySQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM audit_categories WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *auditCategoryMySQLRepository) CountItemsByCategory(ctx context.Context, categoryID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM audit_items WHERE category_id = ?`
	err := r.db.GetContext(ctx, &count, query, categoryID)
	return count, err
}
