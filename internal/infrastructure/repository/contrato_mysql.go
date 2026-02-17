package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type contratoMySQLRepository struct {
	db *sqlx.DB
}

// NewContratoMySQLRepository creates a new MySQL implementation of ContratoRepository
func NewContratoMySQLRepository(db *sqlx.DB) repository.ContratoRepository {
	return &contratoMySQLRepository{db: db}
}

func (r *contratoMySQLRepository) FindAll(ctx context.Context) ([]entity.Contrato, error) {
	var contratos []entity.Contrato
	query := `SELECT id, gestor_id, nome, descricao, endereco, cidade, estado, cep,
			  total_unidades, meta_score, data_inicio, data_fim, ativo, created_at, updated_at
			  FROM contratos
			  WHERE ativo = 1
			  ORDER BY nome`
	err := r.db.SelectContext(ctx, &contratos, query)
	if err != nil {
		return nil, err
	}
	return contratos, nil
}

func (r *contratoMySQLRepository) FindByID(ctx context.Context, id string) (*entity.Contrato, error) {
	var contrato entity.Contrato
	query := `SELECT id, gestor_id, nome, descricao, endereco, cidade, estado, cep,
			  total_unidades, meta_score, data_inicio, data_fim, ativo, created_at, updated_at
			  FROM contratos
			  WHERE id = ?`
	err := r.db.GetContext(ctx, &contrato, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &contrato, nil
}

func (r *contratoMySQLRepository) FindByGestorID(ctx context.Context, gestorID string) ([]entity.Contrato, error) {
	var contratos []entity.Contrato
	query := `SELECT id, gestor_id, nome, descricao, endereco, cidade, estado, cep,
			  total_unidades, meta_score, data_inicio, data_fim, ativo, created_at, updated_at
			  FROM contratos
			  WHERE gestor_id = ? AND ativo = 1
			  ORDER BY nome`
	err := r.db.SelectContext(ctx, &contratos, query, gestorID)
	if err != nil {
		return nil, err
	}
	return contratos, nil
}

func (r *contratoMySQLRepository) FindAllWithGestor(ctx context.Context) ([]entity.ContratoWithGestor, error) {
	var contratos []entity.ContratoWithGestor
	query := `SELECT c.id, c.gestor_id, c.nome, c.descricao, c.endereco, c.cidade, c.estado, c.cep,
			  c.total_unidades, c.meta_score, c.data_inicio, c.data_fim, c.ativo, c.created_at, c.updated_at,
			  g.nome as gestor_nome, g.email as gestor_email
			  FROM contratos c
			  INNER JOIN gestores g ON g.id = c.gestor_id
			  WHERE c.ativo = 1
			  ORDER BY c.nome`
	err := r.db.SelectContext(ctx, &contratos, query)
	if err != nil {
		return nil, err
	}
	return contratos, nil
}

func (r *contratoMySQLRepository) Create(ctx context.Context, contrato *entity.Contrato) error {
	query := `INSERT INTO contratos (id, gestor_id, nome, descricao, endereco, cidade, estado, cep,
			  total_unidades, meta_score, data_inicio, data_fim, ativo, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query,
		contrato.ID, contrato.GestorID, contrato.Nome, contrato.Descricao,
		contrato.Endereco, contrato.Cidade, contrato.Estado, contrato.CEP,
		contrato.TotalUnidades, contrato.MetaScore, contrato.DataInicio, contrato.DataFim, contrato.Ativo)
	return err
}

func (r *contratoMySQLRepository) Update(ctx context.Context, contrato *entity.Contrato) error {
	query := `UPDATE contratos
			  SET gestor_id = ?, nome = ?, descricao = ?, endereco = ?, cidade = ?, estado = ?, cep = ?,
			  total_unidades = ?, meta_score = ?, data_inicio = ?, data_fim = ?, ativo = ?, updated_at = NOW()
			  WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		contrato.GestorID, contrato.Nome, contrato.Descricao,
		contrato.Endereco, contrato.Cidade, contrato.Estado, contrato.CEP,
		contrato.TotalUnidades, contrato.MetaScore, contrato.DataInicio, contrato.DataFim, contrato.Ativo, contrato.ID)
	return err
}

func (r *contratoMySQLRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE contratos SET ativo = 0, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *contratoMySQLRepository) CountByGestorID(ctx context.Context, gestorID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM contratos WHERE gestor_id = ? AND ativo = 1`
	err := r.db.GetContext(ctx, &count, query, gestorID)
	return count, err
}
