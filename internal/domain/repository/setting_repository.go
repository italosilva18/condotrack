package repository

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
)

// SettingRepository defines the interface for setting data access
type SettingRepository interface {
	// GetAll returns all settings
	GetAll(ctx context.Context) ([]*entity.Setting, error)

	// GetByCategory returns all settings for a specific category
	GetByCategory(ctx context.Context, category string) ([]*entity.Setting, error)

	// GetByKey returns a setting by its key
	GetByKey(ctx context.Context, key string) (*entity.Setting, error)

	// GetByID returns a setting by its ID
	GetByID(ctx context.Context, id string) (*entity.Setting, error)

	// Update updates a setting value
	Update(ctx context.Context, key string, value string) error

	// UpdateByID updates a setting value by ID
	UpdateByID(ctx context.Context, id string, value string) error

	// BulkUpdate updates multiple settings at once
	BulkUpdate(ctx context.Context, settings map[string]string) error

	// GetValue returns the raw value of a setting by key
	GetValue(ctx context.Context, key string) (string, error)

	// GetCategories returns all distinct categories
	GetCategories(ctx context.Context) ([]string, error)
}
