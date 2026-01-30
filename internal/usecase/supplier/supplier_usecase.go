package supplier

import (
	"context"
	"errors"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/google/uuid"
)

// UseCase defines the supplier use case interface
type UseCase interface {
	ListSuppliers(ctx context.Context, category *string, isActive *bool) ([]entity.Supplier, error)
	ListActiveSuppliers(ctx context.Context) ([]entity.Supplier, error)
	GetSupplierByID(ctx context.Context, id string) (*entity.Supplier, error)
	CreateSupplier(ctx context.Context, req *entity.CreateSupplierRequest) (*entity.Supplier, error)
	UpdateSupplier(ctx context.Context, id string, req *entity.UpdateSupplierRequest) (*entity.Supplier, error)
	DeleteSupplier(ctx context.Context, id string) error
}

type supplierUseCase struct {
	repo repository.SupplierRepository
}

// NewUseCase creates a new supplier use case
func NewUseCase(repo repository.SupplierRepository) UseCase {
	return &supplierUseCase{repo: repo}
}

// ListSuppliers returns all suppliers with optional category and isActive filters
func (uc *supplierUseCase) ListSuppliers(ctx context.Context, category *string, isActive *bool) ([]entity.Supplier, error) {
	return uc.repo.FindAll(ctx, category, isActive)
}

// ListActiveSuppliers returns all active suppliers
func (uc *supplierUseCase) ListActiveSuppliers(ctx context.Context) ([]entity.Supplier, error) {
	return uc.repo.FindActive(ctx)
}

// GetSupplierByID returns a specific supplier by ID
func (uc *supplierUseCase) GetSupplierByID(ctx context.Context, id string) (*entity.Supplier, error) {
	return uc.repo.FindByID(ctx, id)
}

// CreateSupplier creates a new supplier
func (uc *supplierUseCase) CreateSupplier(ctx context.Context, req *entity.CreateSupplierRequest) (*entity.Supplier, error) {
	// Check if CNPJ already exists (if provided)
	if req.CNPJ != nil && *req.CNPJ != "" {
		existing, err := uc.repo.FindByCNPJ(ctx, *req.CNPJ)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("supplier with this CNPJ already exists")
		}
	}

	// Create supplier entity
	supplier := &entity.Supplier{
		ID:        uuid.New().String(),
		Name:      req.Name,
		CNPJ:      req.CNPJ,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		Category:  req.Category,
		IsActive:  true,
		Notes:     req.Notes,
		CreatedAt: time.Now(),
	}

	if err := uc.repo.Create(ctx, supplier); err != nil {
		return nil, err
	}

	return supplier, nil
}

// UpdateSupplier updates an existing supplier
func (uc *supplierUseCase) UpdateSupplier(ctx context.Context, id string, req *entity.UpdateSupplierRequest) (*entity.Supplier, error) {
	// Find existing supplier
	supplier, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if supplier == nil {
		return nil, errors.New("supplier not found")
	}

	// Check if CNPJ is being changed and if new CNPJ already exists
	if req.CNPJ != nil && *req.CNPJ != "" {
		if supplier.CNPJ == nil || *req.CNPJ != *supplier.CNPJ {
			existing, err := uc.repo.FindByCNPJ(ctx, *req.CNPJ)
			if err != nil {
				return nil, err
			}
			if existing != nil && existing.ID != id {
				return nil, errors.New("supplier with this CNPJ already exists")
			}
		}
		supplier.CNPJ = req.CNPJ
	}

	// Update fields if provided
	if req.Name != nil {
		supplier.Name = *req.Name
	}
	if req.Email != nil {
		supplier.Email = req.Email
	}
	if req.Phone != nil {
		supplier.Phone = req.Phone
	}
	if req.Address != nil {
		supplier.Address = req.Address
	}
	if req.Category != nil {
		supplier.Category = req.Category
	}
	if req.IsActive != nil {
		supplier.IsActive = *req.IsActive
	}
	if req.Notes != nil {
		supplier.Notes = req.Notes
	}

	// Set updated timestamp
	now := time.Now()
	supplier.UpdatedAt = &now

	if err := uc.repo.Update(ctx, supplier); err != nil {
		return nil, err
	}

	return supplier, nil
}

// DeleteSupplier soft deletes a supplier by setting is_active to false
func (uc *supplierUseCase) DeleteSupplier(ctx context.Context, id string) error {
	// Find existing supplier
	supplier, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if supplier == nil {
		return errors.New("supplier not found")
	}

	// Soft delete by calling repository Delete method
	return uc.repo.Delete(ctx, id)
}
