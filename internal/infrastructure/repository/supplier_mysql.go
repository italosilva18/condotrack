package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type supplierMySQLRepository struct {
	db *sqlx.DB
}

// NewSupplierMySQLRepository creates a new MySQL implementation of SupplierRepository
func NewSupplierMySQLRepository(db *sqlx.DB) repository.SupplierRepository {
	return &supplierMySQLRepository{db: db}
}

func (r *supplierMySQLRepository) FindAll(ctx context.Context, category *string, isActive *bool) ([]entity.Supplier, error) {
	var suppliers []entity.Supplier
	query := `SELECT id, name, cnpj, email, phone, address, category, is_active, notes, created_at, updated_at
			  FROM suppliers
			  WHERE 1=1`
	args := []interface{}{}

	if category != nil && *category != "" {
		query += " AND category = ?"
		args = append(args, *category)
	}

	if isActive != nil {
		query += " AND is_active = ?"
		args = append(args, *isActive)
	}

	query += " ORDER BY name"

	err := r.db.SelectContext(ctx, &suppliers, query, args...)
	if err != nil {
		return nil, err
	}
	return suppliers, nil
}

func (r *supplierMySQLRepository) FindByID(ctx context.Context, id string) (*entity.Supplier, error) {
	var supplier entity.Supplier
	query := `SELECT id, name, cnpj, email, phone, address, category, is_active, notes, created_at, updated_at
			  FROM suppliers
			  WHERE id = ?`
	err := r.db.GetContext(ctx, &supplier, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &supplier, nil
}

func (r *supplierMySQLRepository) FindByCategory(ctx context.Context, category string) ([]entity.Supplier, error) {
	var suppliers []entity.Supplier
	query := `SELECT id, name, cnpj, email, phone, address, category, is_active, notes, created_at, updated_at
			  FROM suppliers
			  WHERE category = ? AND is_active = 1
			  ORDER BY name`
	err := r.db.SelectContext(ctx, &suppliers, query, category)
	if err != nil {
		return nil, err
	}
	return suppliers, nil
}

func (r *supplierMySQLRepository) FindActive(ctx context.Context) ([]entity.Supplier, error) {
	var suppliers []entity.Supplier
	query := `SELECT id, name, cnpj, email, phone, address, category, is_active, notes, created_at, updated_at
			  FROM suppliers
			  WHERE is_active = 1
			  ORDER BY name`
	err := r.db.SelectContext(ctx, &suppliers, query)
	if err != nil {
		return nil, err
	}
	return suppliers, nil
}

func (r *supplierMySQLRepository) FindByCNPJ(ctx context.Context, cnpj string) (*entity.Supplier, error) {
	var supplier entity.Supplier
	query := `SELECT id, name, cnpj, email, phone, address, category, is_active, notes, created_at, updated_at
			  FROM suppliers
			  WHERE cnpj = ?`
	err := r.db.GetContext(ctx, &supplier, query, cnpj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &supplier, nil
}

func (r *supplierMySQLRepository) Create(ctx context.Context, supplier *entity.Supplier) error {
	query := `INSERT INTO suppliers (id, name, cnpj, email, phone, address, category, is_active, notes, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query,
		supplier.ID,
		supplier.Name,
		supplier.CNPJ,
		supplier.Email,
		supplier.Phone,
		supplier.Address,
		supplier.Category,
		supplier.IsActive,
		supplier.Notes,
	)
	return err
}

func (r *supplierMySQLRepository) Update(ctx context.Context, supplier *entity.Supplier) error {
	query := `UPDATE suppliers
			  SET name = ?, cnpj = ?, email = ?, phone = ?, address = ?, category = ?, is_active = ?, notes = ?, updated_at = NOW()
			  WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		supplier.Name,
		supplier.CNPJ,
		supplier.Email,
		supplier.Phone,
		supplier.Address,
		supplier.Category,
		supplier.IsActive,
		supplier.Notes,
		supplier.ID,
	)
	return err
}

func (r *supplierMySQLRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE suppliers SET is_active = 0, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
