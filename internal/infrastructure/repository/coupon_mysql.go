package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type couponMySQLRepository struct {
	db *sqlx.DB
}

// NewCouponMySQLRepository creates a new MySQL implementation of CouponRepository.
func NewCouponMySQLRepository(db *sqlx.DB) repository.CouponRepository {
	return &couponMySQLRepository{db: db}
}

const couponColumns = `id, code, description, discount_type, discount_value,
	max_discount_amount, minimum_order_amount, max_uses, max_uses_per_user,
	current_uses, applies_to, course_ids, starts_at, expires_at,
	is_active, created_by, created_at, updated_at`

func (r *couponMySQLRepository) FindByID(ctx context.Context, id string) (*entity.Coupon, error) {
	var c entity.Coupon
	query := `SELECT ` + couponColumns + ` FROM coupons WHERE id = ?`
	err := r.db.GetContext(ctx, &c, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *couponMySQLRepository) FindByCode(ctx context.Context, code string) (*entity.Coupon, error) {
	var c entity.Coupon
	query := `SELECT ` + couponColumns + ` FROM coupons WHERE code = ?`
	err := r.db.GetContext(ctx, &c, query, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *couponMySQLRepository) FindAll(ctx context.Context, activeOnly bool, page, perPage int) ([]entity.Coupon, int, error) {
	where := "1=1"
	args := []interface{}{}
	if activeOnly {
		where = "is_active = true"
	}

	var total int
	countQuery := `SELECT COUNT(*) FROM coupons WHERE ` + where
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	query := `SELECT ` + couponColumns + ` FROM coupons WHERE ` + where + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, perPage, offset)

	var coupons []entity.Coupon
	if err := r.db.SelectContext(ctx, &coupons, query, args...); err != nil {
		return nil, 0, err
	}
	return coupons, total, nil
}

func (r *couponMySQLRepository) Create(ctx context.Context, c *entity.Coupon) error {
	query := `INSERT INTO coupons (
		id, code, description, discount_type, discount_value,
		max_discount_amount, minimum_order_amount, max_uses, max_uses_per_user,
		current_uses, applies_to, course_ids, starts_at, expires_at,
		is_active, created_by, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.Code, c.Description, c.DiscountType, c.DiscountValue,
		c.MaxDiscountAmount, c.MinimumOrderAmount, c.MaxUses, c.MaxUsesPerUser,
		c.CurrentUses, c.AppliesTo, c.CourseIDs, c.StartsAt, c.ExpiresAt,
		c.IsActive, c.CreatedBy,
	)
	return err
}

func (r *couponMySQLRepository) Update(ctx context.Context, c *entity.Coupon) error {
	query := `UPDATE coupons SET
		code = ?, description = ?, discount_type = ?, discount_value = ?,
		max_discount_amount = ?, minimum_order_amount = ?, max_uses = ?, max_uses_per_user = ?,
		applies_to = ?, course_ids = ?, starts_at = ?, expires_at = ?,
		is_active = ?, updated_at = NOW()
		WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		c.Code, c.Description, c.DiscountType, c.DiscountValue,
		c.MaxDiscountAmount, c.MinimumOrderAmount, c.MaxUses, c.MaxUsesPerUser,
		c.AppliesTo, c.CourseIDs, c.StartsAt, c.ExpiresAt,
		c.IsActive, c.ID,
	)
	return err
}

func (r *couponMySQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM coupons WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *couponMySQLRepository) IncrementUsage(ctx context.Context, id string) error {
	query := `UPDATE coupons SET current_uses = current_uses + 1 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *couponMySQLRepository) IncrementUsageWithTx(ctx context.Context, tx *sqlx.Tx, id string) error {
	query := `UPDATE coupons SET current_uses = current_uses + 1 WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, id)
	return err
}

func (r *couponMySQLRepository) CreateUsage(ctx context.Context, u *entity.CouponUsage) error {
	query := `INSERT INTO coupon_usage (
		id, coupon_id, user_id, payment_id, enrollment_id, course_id,
		discount_type, discount_value, discount_applied,
		original_amount, final_amount, used_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.CouponID, u.UserID, u.PaymentID, u.EnrollmentID, u.CourseID,
		u.DiscountType, u.DiscountValue, u.DiscountApplied,
		u.OriginalAmount, u.FinalAmount,
	)
	return err
}

func (r *couponMySQLRepository) CreateUsageWithTx(ctx context.Context, tx *sqlx.Tx, u *entity.CouponUsage) error {
	query := `INSERT INTO coupon_usage (
		id, coupon_id, user_id, payment_id, enrollment_id, course_id,
		discount_type, discount_value, discount_applied,
		original_amount, final_amount, used_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := tx.ExecContext(ctx, query,
		u.ID, u.CouponID, u.UserID, u.PaymentID, u.EnrollmentID, u.CourseID,
		u.DiscountType, u.DiscountValue, u.DiscountApplied,
		u.OriginalAmount, u.FinalAmount,
	)
	return err
}

func (r *couponMySQLRepository) CountUsageByUser(ctx context.Context, couponID, userID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM coupon_usage WHERE coupon_id = ? AND user_id = ?`
	err := r.db.GetContext(ctx, &count, query, couponID, userID)
	return count, err
}
