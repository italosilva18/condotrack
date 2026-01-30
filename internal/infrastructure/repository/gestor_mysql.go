package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type gestorMySQLRepository struct {
	db *sqlx.DB
}

// NewGestorMySQLRepository creates a new MySQL implementation of GestorRepository
func NewGestorMySQLRepository(db *sqlx.DB) repository.GestorRepository {
	return &gestorMySQLRepository{db: db}
}

func (r *gestorMySQLRepository) FindAll(ctx context.Context) ([]entity.Gestor, error) {
	var gestores []entity.Gestor
	query := `SELECT id, nome, email, telefone, cpf, ativo, created_at, updated_at
			  FROM gestores
			  WHERE ativo = 1
			  ORDER BY nome`
	err := r.db.SelectContext(ctx, &gestores, query)
	if err != nil {
		return nil, err
	}
	return gestores, nil
}

func (r *gestorMySQLRepository) FindByID(ctx context.Context, id string) (*entity.Gestor, error) {
	var gestor entity.Gestor
	query := `SELECT id, nome, email, telefone, cpf, ativo, created_at, updated_at
			  FROM gestores
			  WHERE id = ?`
	err := r.db.GetContext(ctx, &gestor, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &gestor, nil
}

func (r *gestorMySQLRepository) FindByEmail(ctx context.Context, email string) (*entity.Gestor, error) {
	var gestor entity.Gestor
	query := `SELECT id, nome, email, telefone, cpf, ativo, created_at, updated_at
			  FROM gestores
			  WHERE email = ?`
	err := r.db.GetContext(ctx, &gestor, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &gestor, nil
}

func (r *gestorMySQLRepository) FindAllWithContracts(ctx context.Context) ([]entity.GestorWithContracts, error) {
	var gestores []entity.GestorWithContracts
	query := `SELECT g.id, g.nome, g.email, g.telefone, g.cpf, g.ativo, g.created_at, g.updated_at,
			  COALESCE(COUNT(c.id), 0) as total_contratos
			  FROM gestores g
			  LEFT JOIN contratos c ON c.gestor_id = g.id AND c.ativo = 1
			  WHERE g.ativo = 1
			  GROUP BY g.id
			  ORDER BY g.nome`
	err := r.db.SelectContext(ctx, &gestores, query)
	if err != nil {
		return nil, err
	}
	return gestores, nil
}

func (r *gestorMySQLRepository) Create(ctx context.Context, gestor *entity.Gestor) error {
	query := `INSERT INTO gestores (id, nome, email, telefone, cpf, ativo, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query, gestor.ID, gestor.Nome, gestor.Email, gestor.Telefone, gestor.CPF, gestor.Ativo)
	return err
}

func (r *gestorMySQLRepository) Update(ctx context.Context, gestor *entity.Gestor) error {
	query := `UPDATE gestores
			  SET nome = ?, email = ?, telefone = ?, cpf = ?, ativo = ?, updated_at = NOW()
			  WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, gestor.Nome, gestor.Email, gestor.Telefone, gestor.CPF, gestor.Ativo, gestor.ID)
	return err
}

func (r *gestorMySQLRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE gestores SET ativo = 0, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
