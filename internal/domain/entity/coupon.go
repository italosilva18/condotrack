package entity

import "time"

// Coupon represents a discount coupon.
type Coupon struct {
	ID                 string     `db:"id" json:"id"`
	Code               string     `db:"code" json:"code"`
	Description        *string    `db:"description" json:"description,omitempty"`
	DiscountType       string     `db:"discount_type" json:"discount_type"`
	DiscountValue      float64    `db:"discount_value" json:"discount_value"`
	MaxDiscountAmount  *float64   `db:"max_discount_amount" json:"max_discount_amount,omitempty"`
	MinimumOrderAmount *float64   `db:"minimum_order_amount" json:"minimum_order_amount,omitempty"`
	MaxUses            *int       `db:"max_uses" json:"max_uses,omitempty"`
	MaxUsesPerUser     *int       `db:"max_uses_per_user" json:"max_uses_per_user,omitempty"`
	CurrentUses        int        `db:"current_uses" json:"current_uses"`
	AppliesTo          string     `db:"applies_to" json:"applies_to"`
	CourseIDs          *string    `db:"course_ids" json:"course_ids,omitempty"` // JSON array
	StartsAt           *time.Time `db:"starts_at" json:"starts_at,omitempty"`
	ExpiresAt          *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	IsActive           bool       `db:"is_active" json:"is_active"`
	CreatedBy          *string    `db:"created_by" json:"created_by,omitempty"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt          *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// CouponUsage represents a coupon redemption record.
type CouponUsage struct {
	ID              string     `db:"id" json:"id"`
	CouponID        string     `db:"coupon_id" json:"coupon_id"`
	UserID          *string    `db:"user_id" json:"user_id,omitempty"`
	PaymentID       *string    `db:"payment_id" json:"payment_id,omitempty"`
	EnrollmentID    *string    `db:"enrollment_id" json:"enrollment_id,omitempty"`
	CourseID        *string    `db:"course_id" json:"course_id,omitempty"`
	DiscountType    string     `db:"discount_type" json:"discount_type"`
	DiscountValue   float64    `db:"discount_value" json:"discount_value"`
	DiscountApplied float64    `db:"discount_applied" json:"discount_applied"`
	OriginalAmount  float64    `db:"original_amount" json:"original_amount"`
	FinalAmount     float64    `db:"final_amount" json:"final_amount"`
	UsedAt          time.Time  `db:"used_at" json:"used_at"`
}

// Coupon discount type constants
const (
	DiscountTypePercentage = "percentage"
	DiscountTypeFixed      = "fixed"
)

// Coupon applies_to constants
const (
	CouponAppliesToAll      = "all_courses"
	CouponAppliesToSpecific = "specific_courses"
)

// CalculateDiscount calculates the actual discount amount for a given order amount.
func (c *Coupon) CalculateDiscount(orderAmount float64) float64 {
	if !c.IsActive {
		return 0
	}

	// Check minimum order amount
	if c.MinimumOrderAmount != nil && orderAmount < *c.MinimumOrderAmount {
		return 0
	}

	// Check usage limits
	if c.MaxUses != nil && c.CurrentUses >= *c.MaxUses {
		return 0
	}

	// Check validity period
	now := time.Now()
	if c.StartsAt != nil && now.Before(*c.StartsAt) {
		return 0
	}
	if c.ExpiresAt != nil && now.After(*c.ExpiresAt) {
		return 0
	}

	var discount float64
	switch c.DiscountType {
	case DiscountTypePercentage:
		discount = orderAmount * (c.DiscountValue / 100)
		if c.MaxDiscountAmount != nil && discount > *c.MaxDiscountAmount {
			discount = *c.MaxDiscountAmount
		}
	case DiscountTypeFixed:
		discount = c.DiscountValue
	}

	// Discount cannot exceed order amount
	if discount > orderAmount {
		discount = orderAmount
	}

	return discount
}
