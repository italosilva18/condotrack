package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type teamMySQLRepository struct {
	db *sqlx.DB
}

// NewTeamMySQLRepository creates a new MySQL implementation of TeamRepository
func NewTeamMySQLRepository(db *sqlx.DB) repository.TeamRepository {
	return &teamMySQLRepository{db: db}
}

func (r *teamMySQLRepository) FindAll(ctx context.Context, filter *entity.TeamMemberFilter) ([]entity.TeamMember, error) {
	var members []entity.TeamMember

	query := `
		SELECT
			tm.id,
			tm.user_id,
			COALESCE(g.nome, '') as user_name,
			COALESCE(g.email, '') as user_email,
			COALESCE(tm.role, '') as user_role,
			tm.contract_id,
			COALESCE(c.nome, '') as contract_name,
			tm.role,
			tm.start_date,
			tm.end_date,
			tm.is_active,
			tm.created_at,
			tm.updated_at
		FROM team_members tm
		LEFT JOIN gestores g ON tm.user_id = g.id
		LEFT JOIN contratos c ON tm.contract_id = c.id
		WHERE 1=1
	`

	var args []interface{}
	var conditions []string

	if filter != nil {
		if filter.ContractID != nil && *filter.ContractID != "" {
			conditions = append(conditions, "tm.contract_id = ?")
			args = append(args, *filter.ContractID)
		}
		if filter.UserID != nil && *filter.UserID != "" {
			conditions = append(conditions, "tm.user_id = ?")
			args = append(args, *filter.UserID)
		}
		if filter.IsActive != nil {
			conditions = append(conditions, "tm.is_active = ?")
			args = append(args, *filter.IsActive)
		}
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY tm.created_at DESC"

	err := r.db.SelectContext(ctx, &members, query, args...)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (r *teamMySQLRepository) FindByID(ctx context.Context, id string) (*entity.TeamMember, error) {
	var member entity.TeamMember

	query := `
		SELECT
			tm.id,
			tm.user_id,
			COALESCE(g.nome, '') as user_name,
			COALESCE(g.email, '') as user_email,
			COALESCE(tm.role, '') as user_role,
			tm.contract_id,
			COALESCE(c.nome, '') as contract_name,
			tm.role,
			tm.start_date,
			tm.end_date,
			tm.is_active,
			tm.created_at,
			tm.updated_at
		FROM team_members tm
		LEFT JOIN gestores g ON tm.user_id = g.id
		LEFT JOIN contratos c ON tm.contract_id = c.id
		WHERE tm.id = ?
	`

	err := r.db.GetContext(ctx, &member, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &member, nil
}

func (r *teamMySQLRepository) FindByContract(ctx context.Context, contractID string) ([]entity.TeamMember, error) {
	var members []entity.TeamMember

	query := `
		SELECT
			tm.id,
			tm.user_id,
			COALESCE(g.nome, '') as user_name,
			COALESCE(g.email, '') as user_email,
			COALESCE(tm.role, '') as user_role,
			tm.contract_id,
			COALESCE(c.nome, '') as contract_name,
			tm.role,
			tm.start_date,
			tm.end_date,
			tm.is_active,
			tm.created_at,
			tm.updated_at
		FROM team_members tm
		LEFT JOIN gestores g ON tm.user_id = g.id
		LEFT JOIN contratos c ON tm.contract_id = c.id
		WHERE tm.contract_id = ?
		ORDER BY tm.role, g.nome
	`

	err := r.db.SelectContext(ctx, &members, query, contractID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (r *teamMySQLRepository) FindByUser(ctx context.Context, userID string) ([]entity.TeamMember, error) {
	var members []entity.TeamMember

	query := `
		SELECT
			tm.id,
			tm.user_id,
			COALESCE(g.nome, '') as user_name,
			COALESCE(g.email, '') as user_email,
			COALESCE(tm.role, '') as user_role,
			tm.contract_id,
			COALESCE(c.nome, '') as contract_name,
			tm.role,
			tm.start_date,
			tm.end_date,
			tm.is_active,
			tm.created_at,
			tm.updated_at
		FROM team_members tm
		LEFT JOIN gestores g ON tm.user_id = g.id
		LEFT JOIN contratos c ON tm.contract_id = c.id
		WHERE tm.user_id = ?
		ORDER BY c.nome, tm.role
	`

	err := r.db.SelectContext(ctx, &members, query, userID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (r *teamMySQLRepository) FindActiveByContract(ctx context.Context, contractID string) ([]entity.TeamMember, error) {
	var members []entity.TeamMember

	query := `
		SELECT
			tm.id,
			tm.user_id,
			COALESCE(g.nome, '') as user_name,
			COALESCE(g.email, '') as user_email,
			COALESCE(tm.role, '') as user_role,
			tm.contract_id,
			COALESCE(c.nome, '') as contract_name,
			tm.role,
			tm.start_date,
			tm.end_date,
			tm.is_active,
			tm.created_at,
			tm.updated_at
		FROM team_members tm
		LEFT JOIN gestores g ON tm.user_id = g.id
		LEFT JOIN contratos c ON tm.contract_id = c.id
		WHERE tm.contract_id = ? AND tm.is_active = 1
		ORDER BY tm.role, g.nome
	`

	err := r.db.SelectContext(ctx, &members, query, contractID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (r *teamMySQLRepository) FindByUserAndContract(ctx context.Context, userID, contractID string) (*entity.TeamMember, error) {
	var member entity.TeamMember

	query := `
		SELECT
			tm.id,
			tm.user_id,
			COALESCE(g.nome, '') as user_name,
			COALESCE(g.email, '') as user_email,
			COALESCE(tm.role, '') as user_role,
			tm.contract_id,
			COALESCE(c.nome, '') as contract_name,
			tm.role,
			tm.start_date,
			tm.end_date,
			tm.is_active,
			tm.created_at,
			tm.updated_at
		FROM team_members tm
		LEFT JOIN gestores g ON tm.user_id = g.id
		LEFT JOIN contratos c ON tm.contract_id = c.id
		WHERE tm.user_id = ? AND tm.contract_id = ?
	`

	err := r.db.GetContext(ctx, &member, query, userID, contractID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &member, nil
}

func (r *teamMySQLRepository) Create(ctx context.Context, member *entity.TeamMemberDB) error {
	query := `
		INSERT INTO team_members (id, user_id, contract_id, role, start_date, end_date, is_active, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
	`

	_, err := r.db.ExecContext(ctx, query,
		member.ID,
		member.UserID,
		member.ContractID,
		member.Role,
		member.StartDate,
		member.EndDate,
		member.IsActive,
	)

	return err
}

func (r *teamMySQLRepository) Update(ctx context.Context, member *entity.TeamMemberDB) error {
	query := `
		UPDATE team_members
		SET role = ?, start_date = ?, end_date = ?, is_active = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		member.Role,
		member.StartDate,
		member.EndDate,
		member.IsActive,
		member.ID,
	)

	return err
}

func (r *teamMySQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM team_members WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
