package team

import (
	"context"
	"errors"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/google/uuid"
)

// UseCase defines the team use case interface
type UseCase interface {
	// ListTeamMembers returns all team members with optional filters
	ListTeamMembers(ctx context.Context, filter *entity.TeamMemberFilter) ([]entity.TeamMember, error)

	// GetTeamMemberByID returns a specific team member by ID
	GetTeamMemberByID(ctx context.Context, id string) (*entity.TeamMember, error)

	// CreateTeamMember creates a new team member assignment
	CreateTeamMember(ctx context.Context, req *entity.CreateTeamMemberRequest) (*entity.TeamMember, error)

	// UpdateTeamMember updates an existing team member assignment
	UpdateTeamMember(ctx context.Context, id string, req *entity.UpdateTeamMemberRequest) (*entity.TeamMember, error)

	// DeleteTeamMember removes a team member assignment
	DeleteTeamMember(ctx context.Context, id string) error

	// GetTeamByContract returns all team members for a specific contract
	GetTeamByContract(ctx context.Context, contractID string) ([]entity.TeamMember, error)

	// GetContractsByUser returns all contracts a user is assigned to
	GetContractsByUser(ctx context.Context, userID string) ([]entity.TeamMember, error)
}

type teamUseCase struct {
	teamRepo     repository.TeamRepository
	gestorRepo   repository.GestorRepository
	contratoRepo repository.ContratoRepository
}

// NewUseCase creates a new team use case
func NewUseCase(
	teamRepo repository.TeamRepository,
	gestorRepo repository.GestorRepository,
	contratoRepo repository.ContratoRepository,
) UseCase {
	return &teamUseCase{
		teamRepo:     teamRepo,
		gestorRepo:   gestorRepo,
		contratoRepo: contratoRepo,
	}
}

// ListTeamMembers returns all team members with optional filters
func (uc *teamUseCase) ListTeamMembers(ctx context.Context, filter *entity.TeamMemberFilter) ([]entity.TeamMember, error) {
	return uc.teamRepo.FindAll(ctx, filter)
}

// GetTeamMemberByID returns a specific team member by ID
func (uc *teamUseCase) GetTeamMemberByID(ctx context.Context, id string) (*entity.TeamMember, error) {
	return uc.teamRepo.FindByID(ctx, id)
}

// CreateTeamMember creates a new team member assignment
func (uc *teamUseCase) CreateTeamMember(ctx context.Context, req *entity.CreateTeamMemberRequest) (*entity.TeamMember, error) {
	// Validate user exists
	user, err := uc.gestorRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Validate contract exists
	contract, err := uc.contratoRepo.FindByID(ctx, req.ContractID)
	if err != nil {
		return nil, err
	}
	if contract == nil {
		return nil, errors.New("contract not found")
	}

	// Check if user is already assigned to this contract
	existing, err := uc.teamRepo.FindByUserAndContract(ctx, req.UserID, req.ContractID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user is already assigned to this contract")
	}

	// Create team member DB entity
	memberDB := &entity.TeamMemberDB{
		ID:         uuid.New().String(),
		UserID:     req.UserID,
		ContractID: req.ContractID,
		Role:       req.Role,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		IsActive:   true,
		CreatedAt:  time.Now(),
	}

	if err := uc.teamRepo.Create(ctx, memberDB); err != nil {
		return nil, err
	}

	// Fetch and return the created team member with joined data
	return uc.teamRepo.FindByID(ctx, memberDB.ID)
}

// UpdateTeamMember updates an existing team member assignment
func (uc *teamUseCase) UpdateTeamMember(ctx context.Context, id string, req *entity.UpdateTeamMemberRequest) (*entity.TeamMember, error) {
	// Find existing team member
	member, err := uc.teamRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errors.New("team member not found")
	}

	// Build update entity
	memberDB := &entity.TeamMemberDB{
		ID:         member.ID,
		UserID:     member.UserID,
		ContractID: member.ContractID,
		Role:       member.Role,
		StartDate:  member.StartDate,
		EndDate:    member.EndDate,
		IsActive:   member.IsActive,
	}

	// Apply updates
	if req.Role != nil {
		memberDB.Role = *req.Role
	}
	if req.StartDate != nil {
		memberDB.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		memberDB.EndDate = req.EndDate
	}
	if req.IsActive != nil {
		memberDB.IsActive = *req.IsActive
	}

	if err := uc.teamRepo.Update(ctx, memberDB); err != nil {
		return nil, err
	}

	// Fetch and return the updated team member with joined data
	return uc.teamRepo.FindByID(ctx, id)
}

// DeleteTeamMember removes a team member assignment
func (uc *teamUseCase) DeleteTeamMember(ctx context.Context, id string) error {
	// Find existing team member
	member, err := uc.teamRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.New("team member not found")
	}

	return uc.teamRepo.Delete(ctx, id)
}

// GetTeamByContract returns all team members for a specific contract
func (uc *teamUseCase) GetTeamByContract(ctx context.Context, contractID string) ([]entity.TeamMember, error) {
	// Validate contract exists
	contract, err := uc.contratoRepo.FindByID(ctx, contractID)
	if err != nil {
		return nil, err
	}
	if contract == nil {
		return nil, errors.New("contract not found")
	}

	return uc.teamRepo.FindByContract(ctx, contractID)
}

// GetContractsByUser returns all contracts a user is assigned to
func (uc *teamUseCase) GetContractsByUser(ctx context.Context, userID string) ([]entity.TeamMember, error) {
	// Validate user exists
	user, err := uc.gestorRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return uc.teamRepo.FindByUser(ctx, userID)
}
