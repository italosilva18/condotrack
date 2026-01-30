package revenue

import (
	"context"
	"errors"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
)

// UseCase defines the revenue split use case interface
type UseCase interface {
	// ListRevenueSplits returns all revenue splits with optional filters
	ListRevenueSplits(ctx context.Context, filters RevenueSplitFilters) ([]entity.RevenueSplit, error)

	// GetRevenueSplitByID returns a revenue split by ID
	GetRevenueSplitByID(ctx context.Context, id string) (*entity.RevenueSplit, error)

	// GetRevenueSplitByEnrollment returns revenue split by enrollment ID
	GetRevenueSplitByEnrollment(ctx context.Context, enrollmentID string) (*entity.RevenueSplit, error)

	// GetInstructorEarnings returns all revenue splits for an instructor
	GetInstructorEarnings(ctx context.Context, instructorID string) ([]entity.RevenueSplit, error)

	// GetInstructorTotalEarnings returns total earnings for an instructor
	GetInstructorTotalEarnings(ctx context.Context, instructorID string) (*InstructorTotalResponse, error)

	// UpdateStatus updates the status of a revenue split
	UpdateStatus(ctx context.Context, id, status string) error
}

// RevenueSplitFilters holds the filter parameters for listing revenue splits
type RevenueSplitFilters struct {
	EnrollmentID string
	InstructorID string
	Status       string
}

// InstructorTotalResponse represents the total earnings response for an instructor
type InstructorTotalResponse struct {
	InstructorID  string  `json:"instructor_id"`
	TotalEarnings float64 `json:"total_earnings"`
	Currency      string  `json:"currency"`
}

type revenueUseCase struct {
	revenueSplitRepo repository.RevenueSplitRepository
}

// NewUseCase creates a new revenue use case
func NewUseCase(revenueSplitRepo repository.RevenueSplitRepository) UseCase {
	return &revenueUseCase{
		revenueSplitRepo: revenueSplitRepo,
	}
}

// ListRevenueSplits returns all revenue splits with optional filters
func (uc *revenueUseCase) ListRevenueSplits(ctx context.Context, filters RevenueSplitFilters) ([]entity.RevenueSplit, error) {
	// If enrollment_id filter is provided, use FindByEnrollmentID
	if filters.EnrollmentID != "" {
		split, err := uc.revenueSplitRepo.FindByEnrollmentID(ctx, filters.EnrollmentID)
		if err != nil {
			return nil, err
		}
		if split == nil {
			return []entity.RevenueSplit{}, nil
		}
		// Apply status filter if provided
		if filters.Status != "" && split.Status != filters.Status {
			return []entity.RevenueSplit{}, nil
		}
		return []entity.RevenueSplit{*split}, nil
	}

	// If instructor_id filter is provided, use FindByInstructorID
	if filters.InstructorID != "" {
		splits, err := uc.revenueSplitRepo.FindByInstructorID(ctx, filters.InstructorID)
		if err != nil {
			return nil, err
		}
		// Apply status filter if provided
		if filters.Status != "" {
			filtered := make([]entity.RevenueSplit, 0)
			for _, s := range splits {
				if s.Status == filters.Status {
					filtered = append(filtered, s)
				}
			}
			return filtered, nil
		}
		return splits, nil
	}

	// If only status filter is provided or no filters, use FindAll
	return uc.revenueSplitRepo.FindAll(ctx, filters.Status)
}

// GetRevenueSplitByID returns a revenue split by ID
func (uc *revenueUseCase) GetRevenueSplitByID(ctx context.Context, id string) (*entity.RevenueSplit, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return uc.revenueSplitRepo.FindByID(ctx, id)
}

// GetRevenueSplitByEnrollment returns revenue split by enrollment ID
func (uc *revenueUseCase) GetRevenueSplitByEnrollment(ctx context.Context, enrollmentID string) (*entity.RevenueSplit, error) {
	if enrollmentID == "" {
		return nil, errors.New("enrollment_id is required")
	}
	return uc.revenueSplitRepo.FindByEnrollmentID(ctx, enrollmentID)
}

// GetInstructorEarnings returns all revenue splits for an instructor
func (uc *revenueUseCase) GetInstructorEarnings(ctx context.Context, instructorID string) ([]entity.RevenueSplit, error) {
	if instructorID == "" {
		return nil, errors.New("instructor_id is required")
	}
	return uc.revenueSplitRepo.FindByInstructorID(ctx, instructorID)
}

// GetInstructorTotalEarnings returns total earnings for an instructor
func (uc *revenueUseCase) GetInstructorTotalEarnings(ctx context.Context, instructorID string) (*InstructorTotalResponse, error) {
	if instructorID == "" {
		return nil, errors.New("instructor_id is required")
	}

	total, err := uc.revenueSplitRepo.GetTotalByInstructor(ctx, instructorID)
	if err != nil {
		return nil, err
	}

	return &InstructorTotalResponse{
		InstructorID:  instructorID,
		TotalEarnings: total,
		Currency:      "BRL",
	}, nil
}

// UpdateStatus updates the status of a revenue split
func (uc *revenueUseCase) UpdateStatus(ctx context.Context, id, status string) error {
	if id == "" {
		return errors.New("id is required")
	}

	// Validate status
	validStatuses := map[string]bool{
		entity.RevenueSplitStatusPending:   true,
		entity.RevenueSplitStatusProcessed: true,
		entity.RevenueSplitStatusFailed:    true,
	}

	if !validStatuses[status] {
		return errors.New("invalid status: must be 'pending', 'processed', or 'failed'")
	}

	// Check if revenue split exists
	split, err := uc.revenueSplitRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if split == nil {
		return errors.New("revenue split not found")
	}

	return uc.revenueSplitRepo.UpdateStatus(ctx, id, status)
}
