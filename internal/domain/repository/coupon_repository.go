package repository

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

// CouponRepository defines the interface for coupon data access.
type CouponRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Coupon, error)
	FindByCode(ctx context.Context, code string) (*entity.Coupon, error)
	FindAll(ctx context.Context, activeOnly bool, page, perPage int) ([]entity.Coupon, int, error)
	Create(ctx context.Context, coupon *entity.Coupon) error
	Update(ctx context.Context, coupon *entity.Coupon) error
	Delete(ctx context.Context, id string) error
	IncrementUsage(ctx context.Context, id string) error
	IncrementUsageWithTx(ctx context.Context, tx *sqlx.Tx, id string) error

	// CouponUsage
	CreateUsage(ctx context.Context, usage *entity.CouponUsage) error
	CreateUsageWithTx(ctx context.Context, tx *sqlx.Tx, usage *entity.CouponUsage) error
	CountUsageByUser(ctx context.Context, couponID, userID string) (int, error)
}
