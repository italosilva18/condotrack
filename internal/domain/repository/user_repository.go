package repository

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// FindByID finds a user by their ID
	FindByID(ctx context.Context, id string) (*entity.User, error)

	// FindByEmail finds a user by their email
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// FindAll returns all active users
	FindAll(ctx context.Context) ([]entity.User, error)

	// FindAllWithFilters returns users with optional filters
	FindAllWithFilters(ctx context.Context, filters UserFilters) ([]entity.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete soft deletes a user (sets ativo = false)
	Delete(ctx context.Context, id string) error

	// UpdateLastLogin updates the user's last login timestamp
	UpdateLastLogin(ctx context.Context, id string) error

	// UpdatePassword updates the user's password hash
	UpdatePassword(ctx context.Context, id string, passwordHash string) error

	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// Count returns the total count of active users
	Count(ctx context.Context) (int64, error)
}

// UserFilters represents filters for querying users
type UserFilters struct {
	Role     *entity.UserRole
	IsActive *bool
	Search   *string
	Limit    int
	Offset   int
}
