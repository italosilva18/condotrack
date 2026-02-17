package entity

import (
	"testing"
	"time"
)

func TestCoupon_CalculateDiscount_Percentage(t *testing.T) {
	coupon := &Coupon{
		ID:           "c1",
		Code:         "PERC10",
		DiscountType: DiscountTypePercentage,
		DiscountValue: 10,
		IsActive:     true,
	}

	discount := coupon.CalculateDiscount(200)
	if discount != 20 {
		t.Errorf("expected 20, got %f", discount)
	}
}

func TestCoupon_CalculateDiscount_Fixed(t *testing.T) {
	coupon := &Coupon{
		ID:           "c2",
		Code:         "FIX25",
		DiscountType: DiscountTypeFixed,
		DiscountValue: 25,
		IsActive:     true,
	}

	discount := coupon.CalculateDiscount(100)
	if discount != 25 {
		t.Errorf("expected 25, got %f", discount)
	}
}

func TestCoupon_CalculateDiscount_Inactive(t *testing.T) {
	coupon := &Coupon{
		ID:           "c3",
		Code:         "INACTIVE",
		DiscountType: DiscountTypePercentage,
		DiscountValue: 50,
		IsActive:     false,
	}

	discount := coupon.CalculateDiscount(100)
	if discount != 0 {
		t.Errorf("expected 0 for inactive coupon, got %f", discount)
	}
}

func TestCoupon_CalculateDiscount_MinimumOrderAmount(t *testing.T) {
	minAmount := 100.0
	coupon := &Coupon{
		ID:                 "c4",
		Code:               "MIN100",
		DiscountType:       DiscountTypePercentage,
		DiscountValue:      10,
		MinimumOrderAmount: &minAmount,
		IsActive:           true,
	}

	// Below minimum
	discount := coupon.CalculateDiscount(50)
	if discount != 0 {
		t.Errorf("expected 0 for order below minimum, got %f", discount)
	}

	// At minimum
	discount = coupon.CalculateDiscount(100)
	if discount != 10 {
		t.Errorf("expected 10, got %f", discount)
	}

	// Above minimum
	discount = coupon.CalculateDiscount(200)
	if discount != 20 {
		t.Errorf("expected 20, got %f", discount)
	}
}

func TestCoupon_CalculateDiscount_MaxDiscountAmount(t *testing.T) {
	maxDiscount := 15.0
	coupon := &Coupon{
		ID:                "c5",
		Code:              "MAX15",
		DiscountType:      DiscountTypePercentage,
		DiscountValue:     50,
		MaxDiscountAmount: &maxDiscount,
		IsActive:          true,
	}

	// 50% of 100 = 50, but capped at 15
	discount := coupon.CalculateDiscount(100)
	if discount != 15 {
		t.Errorf("expected 15 (capped), got %f", discount)
	}

	// 50% of 20 = 10, within cap
	discount = coupon.CalculateDiscount(20)
	if discount != 10 {
		t.Errorf("expected 10 (within cap), got %f", discount)
	}
}

func TestCoupon_CalculateDiscount_MaxUsesExceeded(t *testing.T) {
	maxUses := 5
	coupon := &Coupon{
		ID:           "c6",
		Code:         "MAXUSES",
		DiscountType: DiscountTypeFixed,
		DiscountValue: 10,
		MaxUses:      &maxUses,
		CurrentUses:  5,
		IsActive:     true,
	}

	discount := coupon.CalculateDiscount(100)
	if discount != 0 {
		t.Errorf("expected 0 when max uses exceeded, got %f", discount)
	}
}

func TestCoupon_CalculateDiscount_Expired(t *testing.T) {
	past := time.Now().Add(-24 * time.Hour)
	coupon := &Coupon{
		ID:           "c7",
		Code:         "EXPIRED",
		DiscountType: DiscountTypePercentage,
		DiscountValue: 10,
		ExpiresAt:    &past,
		IsActive:     true,
	}

	discount := coupon.CalculateDiscount(100)
	if discount != 0 {
		t.Errorf("expected 0 for expired coupon, got %f", discount)
	}
}

func TestCoupon_CalculateDiscount_NotYetStarted(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	coupon := &Coupon{
		ID:           "c8",
		Code:         "FUTURE",
		DiscountType: DiscountTypePercentage,
		DiscountValue: 10,
		StartsAt:     &future,
		IsActive:     true,
	}

	discount := coupon.CalculateDiscount(100)
	if discount != 0 {
		t.Errorf("expected 0 for not-yet-started coupon, got %f", discount)
	}
}

func TestCoupon_CalculateDiscount_FixedExceedsOrderAmount(t *testing.T) {
	coupon := &Coupon{
		ID:           "c9",
		Code:         "BIGFIX",
		DiscountType: DiscountTypeFixed,
		DiscountValue: 500,
		IsActive:     true,
	}

	// Fixed discount larger than order = capped at order amount
	discount := coupon.CalculateDiscount(100)
	if discount != 100 {
		t.Errorf("expected 100 (capped at order amount), got %f", discount)
	}
}

func TestCoupon_CalculateDiscount_WithinValidPeriod(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	future := time.Now().Add(24 * time.Hour)
	coupon := &Coupon{
		ID:           "c10",
		Code:         "VALID",
		DiscountType: DiscountTypePercentage,
		DiscountValue: 20,
		StartsAt:     &past,
		ExpiresAt:    &future,
		IsActive:     true,
	}

	discount := coupon.CalculateDiscount(100)
	if discount != 20 {
		t.Errorf("expected 20, got %f", discount)
	}
}

func TestCoupon_CalculateDiscount_ZeroOrder(t *testing.T) {
	coupon := &Coupon{
		ID:           "c11",
		Code:         "ZERO",
		DiscountType: DiscountTypePercentage,
		DiscountValue: 10,
		IsActive:     true,
	}

	discount := coupon.CalculateDiscount(0)
	if discount != 0 {
		t.Errorf("expected 0 for zero order amount, got %f", discount)
	}
}
