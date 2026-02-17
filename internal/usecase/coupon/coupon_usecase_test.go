package coupon

import (
	"context"
	"testing"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/testutil"
)

func newTestUseCase() (UseCase, *testutil.MockCouponRepository) {
	mockRepo := testutil.NewMockCouponRepository()
	uc := NewUseCase(mockRepo)
	return uc, mockRepo
}

func TestCreate_Success(t *testing.T) {
	uc, _ := newTestUseCase()
	coupon, err := uc.Create(context.Background(), &CreateCouponRequest{
		Code:          "SAVE20",
		Description:   "20% off",
		DiscountType:  entity.DiscountTypePercentage,
		DiscountValue: 20,
		IsActive:      true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if coupon.Code != "SAVE20" {
		t.Errorf("expected code SAVE20, got %s", coupon.Code)
	}
	if coupon.DiscountValue != 20 {
		t.Errorf("expected discount value 20, got %f", coupon.DiscountValue)
	}
	if !coupon.IsActive {
		t.Error("expected coupon to be active")
	}
}

func TestCreate_CodeUppercase(t *testing.T) {
	uc, _ := newTestUseCase()
	coupon, err := uc.Create(context.Background(), &CreateCouponRequest{
		Code:          "lowercase",
		DiscountType:  entity.DiscountTypeFixed,
		DiscountValue: 10,
		IsActive:      true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if coupon.Code != "LOWERCASE" {
		t.Errorf("expected uppercase code LOWERCASE, got %s", coupon.Code)
	}
}

func TestCreate_InvalidDiscountType(t *testing.T) {
	uc, _ := newTestUseCase()
	_, err := uc.Create(context.Background(), &CreateCouponRequest{
		Code:          "INVALID",
		DiscountType:  "invalid_type",
		DiscountValue: 10,
		IsActive:      true,
	})

	if err == nil {
		t.Error("expected error for invalid discount type")
	}
}

func TestCreate_PercentageOver100(t *testing.T) {
	uc, _ := newTestUseCase()
	_, err := uc.Create(context.Background(), &CreateCouponRequest{
		Code:          "OVER100",
		DiscountType:  entity.DiscountTypePercentage,
		DiscountValue: 150,
		IsActive:      true,
	})

	if err == nil {
		t.Error("expected error for percentage > 100")
	}
}

func TestCreate_DuplicateCode(t *testing.T) {
	uc, mockRepo := newTestUseCase()

	// Pre-create a coupon
	mockRepo.Coupons["existing"] = &entity.Coupon{
		ID:   "existing",
		Code: "DUP",
	}
	mockRepo.CouponCodes["DUP"] = "existing"

	_, err := uc.Create(context.Background(), &CreateCouponRequest{
		Code:          "DUP",
		DiscountType:  entity.DiscountTypeFixed,
		DiscountValue: 10,
		IsActive:      true,
	})

	if err == nil {
		t.Error("expected error for duplicate code")
	}
}

func TestUpdate_Success(t *testing.T) {
	uc, mockRepo := newTestUseCase()
	mockRepo.Coupons["c1"] = &entity.Coupon{
		ID:           "c1",
		Code:         "OLD",
		DiscountType: entity.DiscountTypePercentage,
		DiscountValue: 10,
		IsActive:     true,
		AppliesTo:    entity.CouponAppliesToAll,
	}
	mockRepo.CouponCodes["OLD"] = "c1"

	newCode := "NEW"
	newValue := 25.0
	updated, err := uc.Update(context.Background(), "c1", &UpdateCouponRequest{
		Code:          &newCode,
		DiscountValue: &newValue,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Code != "NEW" {
		t.Errorf("expected code NEW, got %s", updated.Code)
	}
	if updated.DiscountValue != 25 {
		t.Errorf("expected discount value 25, got %f", updated.DiscountValue)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	uc, _ := newTestUseCase()
	_, err := uc.Update(context.Background(), "nonexistent", &UpdateCouponRequest{})
	if err == nil {
		t.Error("expected error for non-existent coupon")
	}
}

func TestDelete_Success(t *testing.T) {
	uc, mockRepo := newTestUseCase()
	mockRepo.Coupons["c1"] = &entity.Coupon{ID: "c1", Code: "DEL"}
	mockRepo.CouponCodes["DEL"] = "c1"

	err := uc.Delete(context.Background(), "c1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mockRepo.Coupons) != 0 {
		t.Error("expected coupon to be deleted")
	}
}

func TestDelete_NotFound(t *testing.T) {
	uc, _ := newTestUseCase()
	err := uc.Delete(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error for non-existent coupon")
	}
}

func TestValidateCoupon_Valid(t *testing.T) {
	uc, mockRepo := newTestUseCase()
	mockRepo.Coupons["c1"] = &entity.Coupon{
		ID:           "c1",
		Code:         "VALID10",
		DiscountType: entity.DiscountTypePercentage,
		DiscountValue: 10,
		IsActive:     true,
	}
	mockRepo.CouponCodes["VALID10"] = "c1"

	result, err := uc.ValidateCoupon(context.Background(), &ValidateCouponRequest{
		Code:   "VALID10",
		Amount: 100,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Valid {
		t.Error("expected coupon to be valid")
	}
	if result.DiscountAmount != 10 {
		t.Errorf("expected discount 10, got %f", result.DiscountAmount)
	}
	if result.FinalAmount != 90 {
		t.Errorf("expected final amount 90, got %f", result.FinalAmount)
	}
}

func TestValidateCoupon_NotFound(t *testing.T) {
	uc, _ := newTestUseCase()
	result, err := uc.ValidateCoupon(context.Background(), &ValidateCouponRequest{
		Code:   "NONEXISTENT",
		Amount: 100,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Valid {
		t.Error("expected coupon to be invalid")
	}
}

func TestValidateCoupon_Expired(t *testing.T) {
	uc, mockRepo := newTestUseCase()
	past := time.Now().Add(-24 * time.Hour)
	mockRepo.Coupons["c1"] = &entity.Coupon{
		ID:           "c1",
		Code:         "EXPIRED",
		DiscountType: entity.DiscountTypePercentage,
		DiscountValue: 10,
		ExpiresAt:    &past,
		IsActive:     true,
	}
	mockRepo.CouponCodes["EXPIRED"] = "c1"

	result, err := uc.ValidateCoupon(context.Background(), &ValidateCouponRequest{
		Code:   "EXPIRED",
		Amount: 100,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Valid {
		t.Error("expected expired coupon to be invalid")
	}
}

func TestValidateCoupon_PerUserLimit(t *testing.T) {
	uc, mockRepo := newTestUseCase()
	maxPerUser := 1
	mockRepo.Coupons["c1"] = &entity.Coupon{
		ID:             "c1",
		Code:           "ONCE",
		DiscountType:   entity.DiscountTypeFixed,
		DiscountValue:  10,
		MaxUsesPerUser: &maxPerUser,
		IsActive:       true,
	}
	mockRepo.CouponCodes["ONCE"] = "c1"

	// Simulate the user already used it
	mockRepo.Usages["c1:user1"] = 1

	result, err := uc.ValidateCoupon(context.Background(), &ValidateCouponRequest{
		Code:   "ONCE",
		Amount: 100,
		UserID: "user1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Valid {
		t.Error("expected coupon to be invalid when per-user limit exceeded")
	}
}

func TestValidateCoupon_PerUserLimit_OtherUser(t *testing.T) {
	uc, mockRepo := newTestUseCase()
	maxPerUser := 1
	mockRepo.Coupons["c1"] = &entity.Coupon{
		ID:             "c1",
		Code:           "ONCE",
		DiscountType:   entity.DiscountTypeFixed,
		DiscountValue:  10,
		MaxUsesPerUser: &maxPerUser,
		IsActive:       true,
	}
	mockRepo.CouponCodes["ONCE"] = "c1"

	// user1 already used it, but user2 hasn't
	mockRepo.Usages["c1:user1"] = 1

	result, err := uc.ValidateCoupon(context.Background(), &ValidateCouponRequest{
		Code:   "ONCE",
		Amount: 100,
		UserID: "user2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Valid {
		t.Error("expected coupon to be valid for different user")
	}
}

func TestFindAll_ActiveOnly(t *testing.T) {
	uc, mockRepo := newTestUseCase()
	mockRepo.Coupons["c1"] = &entity.Coupon{ID: "c1", Code: "ACTIVE", IsActive: true}
	mockRepo.Coupons["c2"] = &entity.Coupon{ID: "c2", Code: "INACTIVE", IsActive: false}

	coupons, total, err := uc.FindAll(context.Background(), true, 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 active coupon, got %d", total)
	}
	if len(coupons) != 1 {
		t.Errorf("expected 1 coupon, got %d", len(coupons))
	}
}

func TestFindAll_All(t *testing.T) {
	uc, mockRepo := newTestUseCase()
	mockRepo.Coupons["c1"] = &entity.Coupon{ID: "c1", Code: "A", IsActive: true}
	mockRepo.Coupons["c2"] = &entity.Coupon{ID: "c2", Code: "B", IsActive: false}

	coupons, total, err := uc.FindAll(context.Background(), false, 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 total coupons, got %d", total)
	}
	if len(coupons) != 2 {
		t.Errorf("expected 2 coupons, got %d", len(coupons))
	}
}
