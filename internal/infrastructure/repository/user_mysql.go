package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type userMySQLRepository struct {
	db *sqlx.DB
}

// NewUserMySQLRepository creates a new MySQL implementation of UserRepository
func NewUserMySQLRepository(db *sqlx.DB) repository.UserRepository {
	return &userMySQLRepository{db: db}
}

func (r *userMySQLRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	query := `SELECT id, email, password_hash, name, role, is_active, phone, cpf, avatar_url, last_login, created_at, updated_at
			  FROM users
			  WHERE id = ?`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userMySQLRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	query := `SELECT id, email, password_hash, name, role, is_active, phone, cpf, avatar_url, last_login, created_at, updated_at
			  FROM users
			  WHERE email = ?`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userMySQLRepository) FindAll(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	query := `SELECT id, email, password_hash, name, role, is_active, phone, cpf, avatar_url, last_login, created_at, updated_at
			  FROM users
			  WHERE is_active = 1
			  ORDER BY name`
	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userMySQLRepository) FindAllWithFilters(ctx context.Context, filters repository.UserFilters) ([]entity.User, error) {
	var users []entity.User
	var conditions []string
	var args []interface{}

	baseQuery := `SELECT id, email, password_hash, name, role, is_active, phone, cpf, avatar_url, last_login, created_at, updated_at
			      FROM users`

	// Apply filters
	if filters.Role != nil {
		conditions = append(conditions, "role = ?")
		args = append(args, string(*filters.Role))
	}

	if filters.IsActive != nil {
		conditions = append(conditions, "is_active = ?")
		args = append(args, *filters.IsActive)
	}

	if filters.Search != nil && *filters.Search != "" {
		conditions = append(conditions, "(name LIKE ? OR email LIKE ?)")
		searchTerm := "%" + *filters.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	// Build query with conditions
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add ordering
	baseQuery += " ORDER BY name"

	// Add pagination
	if filters.Limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT %d", filters.Limit)
		if filters.Offset > 0 {
			baseQuery += fmt.Sprintf(" OFFSET %d", filters.Offset)
		}
	}

	err := r.db.SelectContext(ctx, &users, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userMySQLRepository) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (id, email, password_hash, name, role, is_active, phone, cpf, avatar_url, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.PasswordHash, user.Nome, user.Role, user.IsActive, user.Phone, user.CPF, user.AvatarURL)
	return err
}

func (r *userMySQLRepository) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users
			  SET email = ?, name = ?, role = ?, is_active = ?, phone = ?, cpf = ?, avatar_url = ?, updated_at = NOW()
			  WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, user.Email, user.Nome, user.Role, user.IsActive, user.Phone, user.CPF, user.AvatarURL, user.ID)
	return err
}

func (r *userMySQLRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE users SET is_active = 0, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userMySQLRepository) UpdateLastLogin(ctx context.Context, id string) error {
	query := `UPDATE users SET last_login = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userMySQLRepository) UpdatePassword(ctx context.Context, id string, passwordHash string) error {
	query := `UPDATE users SET password_hash = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, passwordHash, id)
	return err
}

func (r *userMySQLRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = ?`
	err := r.db.GetContext(ctx, &count, query, email)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userMySQLRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM users WHERE is_active = 1`
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}
