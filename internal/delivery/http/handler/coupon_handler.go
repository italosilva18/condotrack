package handler

import (
	"strconv"

	"github.com/condotrack/api/internal/usecase/coupon"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// CouponHandler handles coupon-related HTTP requests
type CouponHandler struct {
	usecase coupon.UseCase
}

// NewCouponHandler creates a new coupon handler
func NewCouponHandler(uc coupon.UseCase) *CouponHandler {
	return &CouponHandler{usecase: uc}
}

// ListCoupons handles GET /api/v1/coupons
func (h *CouponHandler) ListCoupons(c *gin.Context) {
	ctx := c.Request.Context()

	activeOnly := c.Query("active_only") == "true"
	page := 1
	perPage := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if pp := c.Query("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 && parsed <= 100 {
			perPage = parsed
		}
	}

	coupons, total, err := h.usecase.FindAll(ctx, activeOnly, page, perPage)
	if err != nil {
		response.SafeInternalError(c, "Failed to list coupons", err)
		return
	}

	response.Success(c, gin.H{
		"coupons":  coupons,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

// GetCouponByID handles GET /api/v1/coupons/:id
func (h *CouponHandler) GetCouponByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	cp, err := h.usecase.FindByID(ctx, id)
	if err != nil {
		response.SafeInternalError(c, "Failed to get coupon", err)
		return
	}
	if cp == nil {
		response.NotFound(c, "Coupon not found")
		return
	}

	response.Success(c, cp)
}

// CreateCoupon handles POST /api/v1/coupons
func (h *CouponHandler) CreateCoupon(c *gin.Context) {
	ctx := c.Request.Context()

	var req coupon.CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Set created_by from auth context if available
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok {
			req.CreatedBy = uid
		}
	}

	cp, err := h.usecase.Create(ctx, &req)
	if err != nil {
		response.SafeInternalError(c, "Failed to create coupon", err)
		return
	}

	response.Created(c, cp)
}

// UpdateCoupon handles PUT /api/v1/coupons/:id
func (h *CouponHandler) UpdateCoupon(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req coupon.UpdateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	cp, err := h.usecase.Update(ctx, id, &req)
	if err != nil {
		response.SafeInternalError(c, "Failed to update coupon", err)
		return
	}

	response.Success(c, cp)
}

// DeleteCoupon handles DELETE /api/v1/coupons/:id
func (h *CouponHandler) DeleteCoupon(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if err := h.usecase.Delete(ctx, id); err != nil {
		response.SafeInternalError(c, "Failed to delete coupon", err)
		return
	}

	response.Success(c, gin.H{"message": "Coupon deleted"})
}

// ValidateCoupon handles POST /api/v1/coupons/validate (public endpoint)
func (h *CouponHandler) ValidateCoupon(c *gin.Context) {
	ctx := c.Request.Context()

	var req coupon.ValidateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Set user_id from auth context if available
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok {
			req.UserID = uid
		}
	}

	result, err := h.usecase.ValidateCoupon(ctx, &req)
	if err != nil {
		response.SafeInternalError(c, "Failed to validate coupon", err)
		return
	}

	response.Success(c, result)
}
