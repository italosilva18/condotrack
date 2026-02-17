package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

// SettingMySQLRepository implements SettingRepository for MySQL
type SettingMySQLRepository struct {
	db *sqlx.DB
}

// NewSettingMySQLRepository creates a new SettingMySQLRepository
func NewSettingMySQLRepository(db *sqlx.DB) repository.SettingRepository {
	return &SettingMySQLRepository{db: db}
}

// GetAll returns all settings ordered by category and display_order
func (r *SettingMySQLRepository) GetAll(ctx context.Context) ([]*entity.Setting, error) {
	query := `
		SELECT id, setting_key, setting_value, setting_type, category, label, description,
		       is_secret, is_required, default_value, validation_regex, display_order,
		       created_at, updated_at
		FROM settings
		ORDER BY category, display_order, setting_key
	`

	var settings []*entity.Setting
	err := r.db.SelectContext(ctx, &settings, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all settings: %w", err)
	}

	return settings, nil
}

// GetByCategory returns all settings for a specific category
func (r *SettingMySQLRepository) GetByCategory(ctx context.Context, category string) ([]*entity.Setting, error) {
	query := `
		SELECT id, setting_key, setting_value, setting_type, category, label, description,
		       is_secret, is_required, default_value, validation_regex, display_order,
		       created_at, updated_at
		FROM settings
		WHERE category = ?
		ORDER BY display_order, setting_key
	`

	var settings []*entity.Setting
	err := r.db.SelectContext(ctx, &settings, query, category)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings by category: %w", err)
	}

	return settings, nil
}

// GetByKey returns a setting by its key
func (r *SettingMySQLRepository) GetByKey(ctx context.Context, key string) (*entity.Setting, error) {
	query := `
		SELECT id, setting_key, setting_value, setting_type, category, label, description,
		       is_secret, is_required, default_value, validation_regex, display_order,
		       created_at, updated_at
		FROM settings
		WHERE setting_key = ?
	`

	var setting entity.Setting
	err := r.db.GetContext(ctx, &setting, query, key)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get setting by key: %w", err)
	}

	return &setting, nil
}

// GetByID returns a setting by its ID
func (r *SettingMySQLRepository) GetByID(ctx context.Context, id string) (*entity.Setting, error) {
	query := `
		SELECT id, setting_key, setting_value, setting_type, category, label, description,
		       is_secret, is_required, default_value, validation_regex, display_order,
		       created_at, updated_at
		FROM settings
		WHERE id = ?
	`

	var setting entity.Setting
	err := r.db.GetContext(ctx, &setting, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get setting by id: %w", err)
	}

	return &setting, nil
}

// Update updates a setting value by key
func (r *SettingMySQLRepository) Update(ctx context.Context, key string, value string) error {
	query := `UPDATE settings SET setting_value = ?, updated_at = NOW() WHERE setting_key = ?`

	result, err := r.db.ExecContext(ctx, query, value, key)
	if err != nil {
		return fmt.Errorf("failed to update setting: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("setting not found: %s", key)
	}

	return nil
}

// UpdateByID updates a setting value by ID
func (r *SettingMySQLRepository) UpdateByID(ctx context.Context, id string, value string) error {
	query := `UPDATE settings SET setting_value = ?, updated_at = NOW() WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, value, id)
	if err != nil {
		return fmt.Errorf("failed to update setting: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("setting not found: %s", id)
	}

	return nil
}

// BulkUpdate updates multiple settings at once
func (r *SettingMySQLRepository) BulkUpdate(ctx context.Context, settings map[string]string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE settings SET setting_value = ?, updated_at = NOW() WHERE setting_key = ?`

	for key, value := range settings {
		_, err := tx.ExecContext(ctx, query, value, key)
		if err != nil {
			return fmt.Errorf("failed to update setting %s: %w", key, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetValue returns the raw value of a setting by key
func (r *SettingMySQLRepository) GetValue(ctx context.Context, key string) (string, error) {
	query := `SELECT COALESCE(setting_value, '') FROM settings WHERE setting_key = ?`

	var value string
	err := r.db.GetContext(ctx, &value, query, key)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get setting value: %w", err)
	}

	return value, nil
}

// GetCategories returns all distinct categories
func (r *SettingMySQLRepository) GetCategories(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT category FROM settings ORDER BY category`

	var categories []string
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}
