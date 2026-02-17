package handler

import (
	"strings"

	"github.com/condotrack/api/internal/usecase/gestor"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// GestorHandler handles gestor-related HTTP requests
type GestorHandler struct {
	usecase gestor.UseCase
}

// NewGestorHandler creates a new gestor handler
func NewGestorHandler(uc gestor.UseCase) *GestorHandler {
	return &GestorHandler{usecase: uc}
}

// ListGestores handles GET /api/v1/gestores
func (h *GestorHandler) ListGestores(c *gin.Context) {
	ctx := c.Request.Context()

	// Check if we should include contract counts
	includeContracts := c.Query("include_contracts") == "true"

	if includeContracts {
		gestores, err := h.usecase.ListGestoresWithContracts(ctx)
		if err != nil {
			response.SafeInternalError(c, "Failed to fetch gestores", err)
			return
		}
		response.Success(c, gestores)
		return
	}

	gestores, err := h.usecase.ListGestores(ctx)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch gestores", err)
		return
	}

	response.Success(c, gestores)
}

// GetGestorByID handles GET /api/v1/gestores/:id
func (h *GestorHandler) GetGestorByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	gestor, err := h.usecase.GetGestorByID(ctx, id)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch gestor", err)
		return
	}

	if gestor == nil {
		response.NotFound(c, "Gestor not found")
		return
	}

	response.Success(c, gestor)
}

// CreateGestor handles POST /api/v1/gestores
func (h *GestorHandler) CreateGestor(c *gin.Context) {
	ctx := c.Request.Context()

	var req gestor.CreateGestorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	newGestor, err := h.usecase.CreateGestor(ctx, &req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			response.BadRequest(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to create gestor", err)
		return
	}

	response.Created(c, newGestor)
}

// UpdateGestor handles PUT /api/v1/gestores/:id
func (h *GestorHandler) UpdateGestor(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req gestor.UpdateGestorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	updatedGestor, err := h.usecase.UpdateGestor(ctx, id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, err.Error())
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			response.BadRequest(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to update gestor", err)
		return
	}

	response.Success(c, updatedGestor)
}

// DeleteGestor handles DELETE /api/v1/gestores/:id
func (h *GestorHandler) DeleteGestor(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	err := h.usecase.DeleteGestor(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to delete gestor", err)
		return
	}

	response.SuccessWithMessage(c, "Gestor deleted successfully", nil)
}
