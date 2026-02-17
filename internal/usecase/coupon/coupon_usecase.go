package coupon

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/google/uuid"
)

// UseCase defines the coupon use case interface
type UseCase interface {
	FindByID(ctx context.Context, id string) (*entity.Coupon, error)
	FindAll(ctx context.Context, activeOnly bool, page, perPage int) ([]entity.Coupon, int, error)
	Create(ctx context.Context, req *CreateCouponRequest) (*entity.Coupon, error)
	Update(ctx context.Context, id string, req *UpdateCouponRequest) (*entity.Coupon, error)
	Delete(ctx context.Context, id string) error
	ValidateCoupon(ctx context.Context, req *ValidateCouponRequest) (*ValidateCouponResponse, error)
}

// CreateCouponRequest represents the request to create a coupon
type CreateCouponRequest struct {
	Code               string   `json:"code" binding:"required"`
	Description        string   `json:"description,omitempty"`
	DiscountType       string   `json:"discount_type" binding:"required"` // percentage or fixed
	DiscountValue      float64  `json:"discount_value" binding:"required,gt=0"`
	MaxDiscountAmount  *float64 `json:"max_discount_amount,omitempty"`
	MinimumOrderAmount *float64 `json:"minimum_order_amount,omitempty"`
	MaxUses            *int     `json:"max_uses,omitempty"`
	MaxUsesPerUser     *int     `json:"max_uses_per_user,omitempty"`
	AppliesTo          string   `json:"applies_to,omitempty"` // all_courses or specific_courses
	CourseIDs          string   `json:"course_ids,omitempty"` // JSON array
	StartsAt           *string  `json:"starts_at,omitempty"`  // RFC3339
	ExpiresAt          *string  `json:"expires_at,omitempty"` // RFC3339
	IsActive           bool     `json:"is_active"`
	CreatedBy          string   `json:"-"`
}

// UpdateCouponRequest represents the request to update a coupon
type UpdateCouponRequest struct {
	Code               *string  `json:"code,omitempty"`
	Description        *string  `json:"description,omitempty"`
	DiscountType       *string  `json:"discount_type,omitempty"`
	DiscountValue      *float64 `json:"discount_value,omitempty"`
	MaxDiscountAmount  *float64 `json:"max_discount_amount,omitempty"`
	MinimumOrderAmount *float64 `json:"minimum_order_amount,omitempty"`
	MaxUses            *int     `json:"max_uses,omitempty"`
	MaxUsesPerUser     *int     `json:"max_uses_per_user,omitempty"`
	AppliesTo          *string  `json:"applies_to,omitempty"`
	CourseIDs          *string  `json:"course_ids,omitempty"`
	StartsAt           *string  `json:"starts_at,omitempty"`
	ExpiresAt          *string  `json:"expires_at,omitempty"`
	IsActive           *bool    `json:"is_active,omitempty"`
}

// ValidateCouponRequest represents the request to validate a coupon
type ValidateCouponRequest struct {
	Code     string  `json:"code" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	CourseID string  `json:"course_id,omitempty"`
	UserID   string  `json:"user_id,omitempty"`
}

// ValidateCouponResponse represents the coupon validation result
type ValidateCouponResponse struct {
	Valid           bool    `json:"valid"`
	CouponID        string  `json:"coupon_id,omitempty"`
	Code            string  `json:"code"`
	DiscountType    string  `json:"discount_type,omitempty"`
	DiscountValue   float64 `json:"discount_value,omitempty"`
	DiscountAmount  float64 `json:"discount_amount"`
	FinalAmount     float64 `json:"final_amount"`
	Message         string  `json:"message,omitempty"`
}

type couponUseCase struct {
	repo repository.CouponRepository
}

// NewUseCase creates a new coupon use case
func NewUseCase(repo repository.CouponRepository) UseCase {
	return &couponUseCase{repo: repo}
}

func (uc *couponUseCase) FindByID(ctx context.Context, id string) (*entity.Coupon, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *couponUseCase) FindAll(ctx context.Context, activeOnly bool, page, perPage int) ([]entity.Coupon, int, error) {
	return uc.repo.FindAll(ctx, activeOnly, page, perPage)
}

func (uc *couponUseCase) Create(ctx context.Context, req *CreateCouponRequest) (*entity.Coupon, error) {
	// Validate discount type
	if req.DiscountType != entity.DiscountTypePercentage && req.DiscountType != entity.DiscountTypeFixed {
		return nil, errors.New("discount_type must be 'percentage' or 'fixed'")
	}

	if req.DiscountType == entity.DiscountTypePercentage && req.DiscountValue > 100 {
		return nil, errors.New("percentage discount cannot exceed 100")
	}

	// Check code uniqueness
	existing, err := uc.repo.FindByCode(ctx, strings.ToUpper(req.Code))
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("coupon code already exists")
	}

	appliesTo := req.AppliesTo
	if appliesTo == "" {
		appliesTo = entity.CouponAppliesToAll
	}

	coupon := &entity.Coupon{
		ID:                 uuid.New().String(),
		Code:               strings.ToUpper(req.Code),
		DiscountType:       req.DiscountType,
		DiscountValue:      req.DiscountValue,
		MaxDiscountAmount:  req.MaxDiscountAmount,
		MinimumOrderAmount: req.MinimumOrderAmount,
		MaxUses:            req.MaxUses,
		MaxUsesPerUser:     req.MaxUsesPerUser,
		CurrentUses:        0,
		AppliesTo:          appliesTo,
		IsActive:           req.IsActive,
		CreatedAt:          time.Now(),
	}

	if req.Description != "" {
		coupon.Description = &req.Description
	}
	if req.CourseIDs != "" {
		coupon.CourseIDs = &req.CourseIDs
	}
	if req.CreatedBy != "" {
		coupon.CreatedBy = &req.CreatedBy
	}

	// Parse dates
	if req.StartsAt != nil {
		if t, err := time.Parse(time.RFC3339, *req.StartsAt); err == nil {
			coupon.StartsAt = &t
		}
	}
	if req.ExpiresAt != nil {
		if t, err := time.Parse(time.RFC3339, *req.ExpiresAt); err == nil {
			coupon.ExpiresAt = &t
		}
	}

	if err := uc.repo.Create(ctx, coupon); err != nil {
		return nil, err
	}

	return coupon, nil
}

func (uc *couponUseCase) Update(ctx context.Context, id string, req *UpdateCouponRequest) (*entity.Coupon, error) {
	coupon, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if coupon == nil {
		return nil, errors.New("coupon not found")
	}

	if req.Code != nil {
		code := strings.ToUpper(*req.Code)
		// Check uniqueness if code changed
		if code != coupon.Code {
			existing, err := uc.repo.FindByCode(ctx, code)
			if err != nil {
				return nil, err
			}
			if existing != nil {
				return nil, errors.New("coupon code already exists")
			}
		}
		coupon.Code = code
	}
	if req.Description != nil {
		coupon.Description = req.Description
	}
	if req.DiscountType != nil {
		coupon.DiscountType = *req.DiscountType
	}
	if req.DiscountValue != nil {
		coupon.DiscountValue = *req.DiscountValue
	}
	if req.MaxDiscountAmount != nil {
		coupon.MaxDiscountAmount = req.MaxDiscountAmount
	}
	if req.MinimumOrderAmount != nil {
		coupon.MinimumOrderAmount = req.MinimumOrderAmount
	}
	if req.MaxUses != nil {
		coupon.MaxUses = req.MaxUses
	}
	if req.MaxUsesPerUser != nil {
		coupon.MaxUsesPerUser = req.MaxUsesPerUser
	}
	if req.AppliesTo != nil {
		coupon.AppliesTo = *req.AppliesTo
	}
	if req.CourseIDs != nil {
		coupon.CourseIDs = req.CourseIDs
	}
	if req.IsActive != nil {
		coupon.IsActive = *req.IsActive
	}
	if req.StartsAt != nil {
		if t, err := time.Parse(time.RFC3339, *req.StartsAt); err == nil {
			coupon.StartsAt = &t
		}
	}
	if req.ExpiresAt != nil {
		if t, err := time.Parse(time.RFC3339, *req.ExpiresAt); err == nil {
			coupon.ExpiresAt = &t
		}
	}

	if err := uc.repo.Update(ctx, coupon); err != nil {
		return nil, err
	}

	return coupon, nil
}

func (uc *couponUseCase) Delete(ctx context.Context, id string) error {
	coupon, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if coupon == nil {
		return errors.New("coupon not found")
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *couponUseCase) ValidateCoupon(ctx context.Context, req *ValidateCouponRequest) (*ValidateCouponResponse, error) {
	coupon, err := uc.repo.FindByCode(ctx, strings.ToUpper(req.Code))
	if err != nil {
		return nil, err
	}

	if coupon == nil {
		return &ValidateCouponResponse{
			Valid:       false,
			Code:        req.Code,
			FinalAmount: req.Amount,
			Message:     "Cupom não encontrado",
		}, nil
	}

	// Check per-user usage
	if coupon.MaxUsesPerUser != nil && req.UserID != "" {
		usageCount, err := uc.repo.CountUsageByUser(ctx, coupon.ID, req.UserID)
		if err != nil {
			return nil, err
		}
		if usageCount >= *coupon.MaxUsesPerUser {
			return &ValidateCouponResponse{
				Valid:       false,
				Code:        req.Code,
				FinalAmount: req.Amount,
				Message:     "Limite de uso deste cupom excedido",
			}, nil
		}
	}

	discount := coupon.CalculateDiscount(req.Amount)
	if discount <= 0 {
		return &ValidateCouponResponse{
			Valid:       false,
			Code:        req.Code,
			FinalAmount: req.Amount,
			Message:     "Cupom não aplicável a este pedido",
		}, nil
	}

	return &ValidateCouponResponse{
		Valid:          true,
		CouponID:       coupon.ID,
		Code:           coupon.Code,
		DiscountType:   coupon.DiscountType,
		DiscountValue:  coupon.DiscountValue,
		DiscountAmount: discount,
		FinalAmount:    req.Amount - discount,
		Message:        "Cupom válido",
	}, nil
}
