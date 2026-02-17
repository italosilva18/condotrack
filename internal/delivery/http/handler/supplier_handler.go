package handler

import (
	"strings"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/usecase/supplier"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// SupplierHandler handles supplier-related HTTP requests
type SupplierHandler struct {
	usecase supplier.UseCase
}

// NewSupplierHandler creates a new supplier handler
func NewSupplierHandler(uc supplier.UseCase) *SupplierHandler {
	return &SupplierHandler{usecase: uc}
}

// ListSuppliers handles GET /api/v1/suppliers
// Query params: category (string), is_active (bool)
func (h *SupplierHandler) ListSuppliers(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse query parameters
	var category *string
	var isActive *bool

	if cat := c.Query("category"); cat != "" {
		category = &cat
	}

	if active := c.Query("is_active"); active != "" {
		isActiveVal := active == "true" || active == "1"
		isActive = &isActiveVal
	}

	suppliers, err := h.usecase.ListSuppliers(ctx, category, isActive)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch suppliers", err)
		return
	}

	response.Success(c, suppliers)
}

// GetSupplierByID handles GET /api/v1/suppliers/:id
func (h *SupplierHandler) GetSupplierByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	supplier, err := h.usecase.GetSupplierByID(ctx, id)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch supplier", err)
		return
	}

	if supplier == nil {
		response.NotFound(c, "Supplier not found")
		return
	}

	response.Success(c, supplier)
}

// CreateSupplier handles POST /api/v1/suppliers
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	newSupplier, err := h.usecase.CreateSupplier(ctx, &req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			response.BadRequest(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to create supplier", err)
		return
	}

	response.Created(c, newSupplier)
}

// UpdateSupplier handles PUT /api/v1/suppliers/:id
func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req entity.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	updatedSupplier, err := h.usecase.UpdateSupplier(ctx, id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, err.Error())
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			response.BadRequest(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to update supplier", err)
		return
	}

	response.Success(c, updatedSupplier)
}

// DeleteSupplier handles DELETE /api/v1/suppliers/:id
func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	err := h.usecase.DeleteSupplier(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to delete supplier", err)
		return
	}

	response.SuccessWithMessage(c, "Supplier deleted successfully", nil)
}
