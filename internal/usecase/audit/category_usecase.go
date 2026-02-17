package audit

import (
	"context"
	"errors"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/google/uuid"
)

// CategoryUseCase defines the audit category use case interface
type CategoryUseCase interface {
	ListCategories(ctx context.Context) ([]entity.AuditCategory, error)
	GetCategoryByID(ctx context.Context, id string) (*entity.AuditCategory, error)
	CreateCategory(ctx context.Context, req *entity.CreateAuditCategoryRequest) (*entity.AuditCategory, error)
	UpdateCategory(ctx context.Context, id string, req *entity.UpdateAuditCategoryRequest) (*entity.AuditCategory, error)
	DeleteCategory(ctx context.Context, id string) error
}

type categoryUseCase struct {
	repo repository.AuditCategoryRepository
}

// NewCategoryUseCase creates a new audit category use case
func NewCategoryUseCase(repo repository.AuditCategoryRepository) CategoryUseCase {
	return &categoryUseCase{repo: repo}
}

// ListCategories returns all audit categories
func (uc *categoryUseCase) ListCategories(ctx context.Context) ([]entity.AuditCategory, error) {
	return uc.repo.FindAll(ctx)
}

// GetCategoryByID returns a specific category by ID
func (uc *categoryUseCase) GetCategoryByID(ctx context.Context, id string) (*entity.AuditCategory, error) {
	return uc.repo.FindByID(ctx, id)
}

// CreateCategory creates a new audit category
func (uc *categoryUseCase) CreateCategory(ctx context.Context, req *entity.CreateAuditCategoryRequest) (*entity.AuditCategory, error) {
	// Check if category with same name already exists
	existing, err := uc.repo.FindByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("category with this name already exists")
	}

	// Create category entity
	category := &entity.AuditCategory{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Weight:      req.Weight,
		Order:       req.Order,
	}

	// Set default weight if not provided
	if category.Weight == 0 {
		category.Weight = 1.0
	}

	if err := uc.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory updates an existing audit category
func (uc *categoryUseCase) UpdateCategory(ctx context.Context, id string, req *entity.UpdateAuditCategoryRequest) (*entity.AuditCategory, error) {
	// Get existing category
	category, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	// Check if new name conflicts with another category
	if req.Name != nil && *req.Name != category.Name {
		existing, err := uc.repo.FindByName(ctx, *req.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, errors.New("category with this name already exists")
		}
		category.Name = *req.Name
	}

	// Update fields if provided
	if req.Description != nil {
		category.Description = req.Description
	}
	if req.Weight != nil {
		category.Weight = *req.Weight
	}
	if req.Order != nil {
		category.Order = *req.Order
	}

	if err := uc.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory deletes an audit category
func (uc *categoryUseCase) DeleteCategory(ctx context.Context, id string) error {
	// Check if category exists
	category, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("category not found")
	}

	// Check if category is being used by audit items
	count, err := uc.repo.CountItemsByCategory(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("cannot delete category: it is being used by audit items")
	}

	return uc.repo.Delete(ctx, id)
}
