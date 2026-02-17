package handler

import (
	"github.com/condotrack/api/internal/usecase/revenue"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// RevenueHandler handles revenue split-related HTTP requests
type RevenueHandler struct {
	usecase revenue.UseCase
}

// NewRevenueHandler creates a new revenue handler
func NewRevenueHandler(uc revenue.UseCase) *RevenueHandler {
	return &RevenueHandler{usecase: uc}
}

// ListRevenueSplits handles GET /api/v1/revenue-splits
// Query parameters: enrollment_id, instructor_id, status
func (h *RevenueHandler) ListRevenueSplits(c *gin.Context) {
	ctx := c.Request.Context()

	filters := revenue.RevenueSplitFilters{
		EnrollmentID: c.Query("enrollment_id"),
		InstructorID: c.Query("instructor_id"),
		Status:       c.Query("status"),
	}

	splits, err := h.usecase.ListRevenueSplits(ctx, filters)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch revenue splits", err)
		return
	}

	response.Success(c, splits)
}

// GetRevenueSplitByID handles GET /api/v1/revenue-splits/:id
func (h *RevenueHandler) GetRevenueSplitByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	split, err := h.usecase.GetRevenueSplitByID(ctx, id)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch revenue split", err)
		return
	}

	if split == nil {
		response.NotFound(c, "Revenue split not found")
		return
	}

	response.Success(c, split)
}

// GetRevenueSplitByEnrollment handles GET /api/v1/revenue-splits/enrollment/:id
func (h *RevenueHandler) GetRevenueSplitByEnrollment(c *gin.Context) {
	ctx := c.Request.Context()
	enrollmentID := c.Param("id")

	split, err := h.usecase.GetRevenueSplitByEnrollment(ctx, enrollmentID)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch revenue split", err)
		return
	}

	if split == nil {
		response.NotFound(c, "Revenue split not found for this enrollment")
		return
	}

	response.Success(c, split)
}

// GetInstructorEarnings handles GET /api/v1/revenue-splits/instructor/:id
func (h *RevenueHandler) GetInstructorEarnings(c *gin.Context) {
	ctx := c.Request.Context()
	instructorID := c.Param("id")

	splits, err := h.usecase.GetInstructorEarnings(ctx, instructorID)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch instructor earnings", err)
		return
	}

	response.Success(c, splits)
}

// GetInstructorTotalEarnings handles GET /api/v1/revenue-splits/instructor/:id/total
func (h *RevenueHandler) GetInstructorTotalEarnings(c *gin.Context) {
	ctx := c.Request.Context()
	instructorID := c.Param("id")

	total, err := h.usecase.GetInstructorTotalEarnings(ctx, instructorID)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch instructor total earnings", err)
		return
	}

	response.Success(c, total)
}

// UpdateStatus handles PATCH /api/v1/revenue-splits/:id/status
func (h *RevenueHandler) UpdateStatus(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.usecase.UpdateStatus(ctx, id, req.Status); err != nil {
		if err.Error() == "revenue split not found" {
			response.NotFound(c, err.Error())
			return
		}
		if err.Error() == "invalid status: must be 'pending', 'processed', or 'failed'" {
			response.BadRequest(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to update status", err)
		return
	}

	response.Success(c, gin.H{
		"message": "Revenue split status updated successfully",
	})
}
