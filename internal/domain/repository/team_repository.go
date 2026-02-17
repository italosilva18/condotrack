package repository

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
)

// TeamRepository defines the interface for team member data access
type TeamRepository interface {
	// FindAll returns all team members with optional filters
	FindAll(ctx context.Context, filter *entity.TeamMemberFilter) ([]entity.TeamMember, error)

	// FindByID returns a team member by ID
	FindByID(ctx context.Context, id string) (*entity.TeamMember, error)

	// FindByContract returns all team members for a specific contract
	FindByContract(ctx context.Context, contractID string) ([]entity.TeamMember, error)

	// FindByUser returns all team member assignments for a specific user
	FindByUser(ctx context.Context, userID string) ([]entity.TeamMember, error)

	// FindActiveByContract returns all active team members for a specific contract
	FindActiveByContract(ctx context.Context, contractID string) ([]entity.TeamMember, error)

	// FindByUserAndContract returns a team member by user ID and contract ID
	FindByUserAndContract(ctx context.Context, userID, contractID string) (*entity.TeamMember, error)

	// Create creates a new team member assignment
	Create(ctx context.Context, member *entity.TeamMemberDB) error

	// Update updates an existing team member assignment
	Update(ctx context.Context, member *entity.TeamMemberDB) error

	// Delete removes a team member assignment by ID
	Delete(ctx context.Context, id string) error
}
