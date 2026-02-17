package handler

import (
	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/usecase/audit"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuditHandler handles audit-related HTTP requests
type AuditHandler struct {
	usecase audit.UseCase
}

// NewAuditHandler creates a new audit handler
func NewAuditHandler(uc audit.UseCase) *AuditHandler {
	return &AuditHandler{usecase: uc}
}

// ListAudits handles GET /api/v1/audits
func (h *AuditHandler) ListAudits(c *gin.Context) {
	ctx := c.Request.Context()

	// Check for contract_id filter
	contractID := c.Query("contract_id")
	if contractID != "" {
		audits, err := h.usecase.ListAuditsByContract(ctx, contractID)
		if err != nil {
			response.SafeInternalError(c, "Failed to fetch audits", err)
			return
		}
		response.Success(c, audits)
		return
	}

	// Check if we should include contract info
	includeContract := c.Query("include_contract") == "true"

	if includeContract {
		audits, err := h.usecase.ListAuditsWithContract(ctx)
		if err != nil {
			response.SafeInternalError(c, "Failed to fetch audits", err)
			return
		}
		response.Success(c, audits)
		return
	}

	audits, err := h.usecase.ListAudits(ctx)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch audits", err)
		return
	}

	response.Success(c, audits)
}

// GetAuditByID handles GET /api/v1/audits/:id
func (h *AuditHandler) GetAuditByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	audit, err := h.usecase.GetAuditByID(ctx, id)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch audit", err)
		return
	}

	if audit == nil {
		response.NotFound(c, "Audit not found")
		return
	}

	response.Success(c, audit)
}

// GetAuditMeta handles GET /api/v1/audits/meta
func (h *AuditHandler) GetAuditMeta(c *gin.Context) {
	ctx := c.Request.Context()

	contractID := c.Query("contract_id")
	if contractID == "" {
		response.BadRequest(c, "contract_id is required")
		return
	}

	meta, err := h.usecase.GetAuditMeta(ctx, contractID)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch audit meta", err)
		return
	}

	response.Success(c, meta)
}

// CreateAudit handles POST /api/v1/audits
func (h *AuditHandler) CreateAudit(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.CreateAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	audit, err := h.usecase.CreateAudit(ctx, &req)
	if err != nil {
		response.SafeInternalError(c, "Failed to create audit", err)
		return
	}

	// Return with calculated status
	c.JSON(200, gin.H{
		"success": true,
		"id":      audit.ID,
		"status":  audit.Status,
		"score":   audit.Score,
		"data":    audit,
	})
}

// UpdateAudit handles PUT /api/v1/audits/:id
func (h *AuditHandler) UpdateAudit(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req entity.UpdateAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	audit, err := h.usecase.UpdateAudit(ctx, id, &req)
	if err != nil {
		if err.Error() == "audit not found" {
			response.NotFound(c, "Audit not found")
			return
		}
		response.SafeInternalError(c, "Failed to update audit", err)
		return
	}

	response.Success(c, audit)
}

// DeleteAudit handles DELETE /api/v1/audits/:id
func (h *AuditHandler) DeleteAudit(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	err := h.usecase.DeleteAudit(ctx, id)
	if err != nil {
		if err.Error() == "audit not found" {
			response.NotFound(c, "Audit not found")
			return
		}
		response.SafeInternalError(c, "Failed to delete audit", err)
		return
	}

	response.Success(c, map[string]string{"message": "Audit deleted successfully"})
}
