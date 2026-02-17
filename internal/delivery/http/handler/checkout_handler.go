package handler

import (
	"github.com/condotrack/api/internal/usecase/checkout"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// CheckoutHandler handles checkout-related HTTP requests
type CheckoutHandler struct {
	usecase checkout.UseCase
}

// NewCheckoutHandler creates a new checkout handler
func NewCheckoutHandler(uc checkout.UseCase) *CheckoutHandler {
	return &CheckoutHandler{usecase: uc}
}

// CreateCheckout handles POST /api/v1/checkout
func (h *CheckoutHandler) CreateCheckout(c *gin.Context) {
	ctx := c.Request.Context()

	var req checkout.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	result, err := h.usecase.CreateCheckout(ctx, &req)
	if err != nil {
		response.SafeInternalError(c, "Failed to create checkout", err)
		return
	}

	response.Created(c, result)
}

// GetCheckoutStatus handles GET /api/v1/checkout/:id/status
func (h *CheckoutHandler) GetCheckoutStatus(c *gin.Context) {
	ctx := c.Request.Context()
	enrollmentID := c.Param("id")

	result, err := h.usecase.GetCheckoutStatus(ctx, enrollmentID)
	if err != nil {
		response.SafeInternalError(c, "Failed to get checkout status", err)
		return
	}

	response.Success(c, result)
}
