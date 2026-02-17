package handler

import (
	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/usecase/contrato"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// ContratoHandler handles contrato-related HTTP requests
type ContratoHandler struct {
	usecase contrato.UseCase
}

// NewContratoHandler creates a new contrato handler
func NewContratoHandler(uc contrato.UseCase) *ContratoHandler {
	return &ContratoHandler{usecase: uc}
}

// ListContratos handles GET /api/v1/contratos
func (h *ContratoHandler) ListContratos(c *gin.Context) {
	ctx := c.Request.Context()

	// Check for gestor_id filter
	gestorID := c.Query("gestor_id")
	if gestorID != "" {
		contratos, err := h.usecase.ListContratosByGestor(ctx, gestorID)
		if err != nil {
			response.SafeInternalError(c, "Failed to fetch contratos", err)
			return
		}
		response.Success(c, contratos)
		return
	}

	// Check if we should include gestor info
	includeGestor := c.Query("include_gestor") == "true"

	if includeGestor {
		contratos, err := h.usecase.ListContratosWithGestor(ctx)
		if err != nil {
			response.SafeInternalError(c, "Failed to fetch contratos", err)
			return
		}
		response.Success(c, contratos)
		return
	}

	contratos, err := h.usecase.ListContratos(ctx)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch contratos", err)
		return
	}

	response.Success(c, contratos)
}

// GetContratoByID handles GET /api/v1/contratos/:id
func (h *ContratoHandler) GetContratoByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	contrato, err := h.usecase.GetContratoByID(ctx, id)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch contrato", err)
		return
	}

	if contrato == nil {
		response.NotFound(c, "Contrato not found")
		return
	}

	response.Success(c, contrato)
}

// CreateContrato handles POST /api/v1/contratos
func (h *ContratoHandler) CreateContrato(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.CreateContratoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	contrato, err := h.usecase.CreateContrato(ctx, &req)
	if err != nil {
		response.SafeInternalError(c, "Failed to create contrato", err)
		return
	}

	response.Created(c, contrato)
}

// UpdateContrato handles PUT /api/v1/contratos/:id
func (h *ContratoHandler) UpdateContrato(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req entity.UpdateContratoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	contrato, err := h.usecase.UpdateContrato(ctx, id, &req)
	if err != nil {
		if err.Error() == "contrato not found" {
			response.NotFound(c, "Contrato not found")
			return
		}
		if err.Error() == "gestor not found" {
			response.BadRequest(c, "Gestor not found")
			return
		}
		response.SafeInternalError(c, "Failed to update contrato", err)
		return
	}

	response.Success(c, contrato)
}

// DeleteContrato handles DELETE /api/v1/contratos/:id
func (h *ContratoHandler) DeleteContrato(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	err := h.usecase.DeleteContrato(ctx, id)
	if err != nil {
		if err.Error() == "contrato not found" {
			response.NotFound(c, "Contrato not found")
			return
		}
		response.SafeInternalError(c, "Failed to delete contrato", err)
		return
	}

	response.Success(c, map[string]string{"message": "Contrato deleted successfully"})
}
