package handler

import (
	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/usecase/audit"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuditCategoryHandler handles audit category-related HTTP requests
type AuditCategoryHandler struct {
	usecase audit.CategoryUseCase
}

// NewAuditCategoryHandler creates a new audit category handler
func NewAuditCategoryHandler(uc audit.CategoryUseCase) *AuditCategoryHandler {
	return &AuditCategoryHandler{usecase: uc}
}

// ListCategories handles GET /api/v1/audit-categories
func (h *AuditCategoryHandler) ListCategories(c *gin.Context) {
	ctx := c.Request.Context()

	categories, err := h.usecase.ListCategories(ctx)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch audit categories", err)
		return
	}

	response.Success(c, categories)
}

// GetCategoryByID handles GET /api/v1/audit-categories/:id
func (h *AuditCategoryHandler) GetCategoryByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	category, err := h.usecase.GetCategoryByID(ctx, id)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch audit category", err)
		return
	}

	if category == nil {
		response.NotFound(c, "Audit category not found")
		return
	}

	response.Success(c, category)
}

// CreateCategory handles POST /api/v1/audit-categories
func (h *AuditCategoryHandler) CreateCategory(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.CreateAuditCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	category, err := h.usecase.CreateCategory(ctx, &req)
	if err != nil {
		if err.Error() == "category with this name already exists" {
			response.BadRequest(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to create audit category", err)
		return
	}

	response.Created(c, category)
}

// UpdateCategory handles PUT /api/v1/audit-categories/:id
func (h *AuditCategoryHandler) UpdateCategory(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req entity.UpdateAuditCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	category, err := h.usecase.UpdateCategory(ctx, id, &req)
	if err != nil {
		if err.Error() == "category not found" {
			response.NotFound(c, err.Error())
			return
		}
		if err.Error() == "category with this name already exists" {
			response.BadRequest(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to update audit category", err)
		return
	}

	response.Success(c, category)
}

// DeleteCategory handles DELETE /api/v1/audit-categories/:id
func (h *AuditCategoryHandler) DeleteCategory(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	err := h.usecase.DeleteCategory(ctx, id)
	if err != nil {
		if err.Error() == "category not found" {
			response.NotFound(c, err.Error())
			return
		}
		if err.Error() == "cannot delete category: it is being used by audit items" {
			response.BadRequest(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to delete audit category", err)
		return
	}

	response.SuccessWithMessage(c, "Audit category deleted successfully", nil)
}
