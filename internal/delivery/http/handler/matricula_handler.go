package handler

import (
	"strconv"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/usecase/matricula"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// MatriculaHandler handles matricula-related HTTP requests
type MatriculaHandler struct {
	usecase matricula.UseCase
}

// NewMatriculaHandler creates a new matricula handler
func NewMatriculaHandler(uc matricula.UseCase) *MatriculaHandler {
	return &MatriculaHandler{usecase: uc}
}

// ListEnrollments handles GET /api/v1/enrollments
func (h *MatriculaHandler) ListEnrollments(c *gin.Context) {
	ctx := c.Request.Context()

	// Check for student_id filter
	studentID := c.Query("student_id")
	if studentID != "" {
		enrollments, err := h.usecase.ListEnrollmentsByStudent(ctx, studentID)
		if err != nil {
			response.SafeInternalError(c, "Failed to fetch enrollments", err)
			return
		}
		response.Success(c, enrollments)
		return
	}

	// Parse pagination parameters
	page := 1
	perPage := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if pp := c.Query("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 {
			if parsed > 100 {
				parsed = 100
			}
			perPage = parsed
		}
	}

	result, err := h.usecase.ListEnrollments(ctx, page, perPage)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch enrollments", err)
		return
	}

	response.Success(c, result)
}

// GetEnrollmentByID handles GET /api/v1/enrollments/:id
func (h *MatriculaHandler) GetEnrollmentByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	enrollment, err := h.usecase.GetEnrollmentByID(ctx, id)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch enrollment", err)
		return
	}

	if enrollment == nil {
		response.NotFound(c, "Enrollment not found")
		return
	}

	response.Success(c, enrollment)
}

// CreateEnrollment handles POST /api/v1/enrollments
func (h *MatriculaHandler) CreateEnrollment(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.CreateMatriculaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	enrollment, err := h.usecase.CreateEnrollment(ctx, &req)
	if err != nil {
		response.SafeInternalError(c, "Failed to create enrollment", err)
		return
	}

	response.Created(c, enrollment)
}

// UpdatePaymentStatus handles PATCH /api/v1/enrollments/:id/payment-status
func (h *MatriculaHandler) UpdatePaymentStatus(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.usecase.UpdatePaymentStatus(ctx, id, req.Status); err != nil {
		response.SafeInternalError(c, "Failed to update payment status", err)
		return
	}

	response.Success(c, gin.H{
		"message": "Payment status updated successfully",
	})
}

// UpdateProgress handles PATCH /api/v1/enrollments/:id/progress
func (h *MatriculaHandler) UpdateProgress(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req struct {
		Progress float64 `json:"progress" binding:"required,min=0,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.usecase.UpdateProgress(ctx, id, req.Progress); err != nil {
		response.SafeInternalError(c, "Failed to update progress", err)
		return
	}

	response.Success(c, gin.H{
		"message": "Progress updated successfully",
	})
}
