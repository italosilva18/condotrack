package handler

import (
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/usecase/inspection"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// InspectionHandler handles inspection-related HTTP requests
type InspectionHandler struct {
	usecase inspection.UseCase
}

// NewInspectionHandler creates a new inspection handler
func NewInspectionHandler(uc inspection.UseCase) *InspectionHandler {
	return &InspectionHandler{usecase: uc}
}

// ListInspections handles GET /api/v1/inspections
func (h *InspectionHandler) ListInspections(c *gin.Context) {
	ctx := c.Request.Context()

	// Build filter from query parameters
	filter := &entity.InspectionFilter{}
	hasFilter := false

	if contractID := c.Query("contract_id"); contractID != "" {
		filter.ContractID = contractID
		hasFilter = true
	}

	if inspectorID := c.Query("inspector_id"); inspectorID != "" {
		filter.InspectorID = inspectorID
		hasFilter = true
	}

	if status := c.Query("status"); status != "" {
		if !entity.IsValidInspectionStatus(status) {
			response.BadRequest(c, "Invalid status value")
			return
		}
		filter.Status = status
		hasFilter = true
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			response.BadRequest(c, "Invalid start_date format. Use YYYY-MM-DD")
			return
		}
		filter.StartDate = &startDate
		hasFilter = true
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			response.BadRequest(c, "Invalid end_date format. Use YYYY-MM-DD")
			return
		}
		// Set end of day for end date
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		filter.EndDate = &endDate
		hasFilter = true
	}

	var filterPtr *entity.InspectionFilter
	if hasFilter {
		filterPtr = filter
	}

	inspections, err := h.usecase.ListInspections(ctx, filterPtr)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch inspections", err)
		return
	}

	response.Success(c, inspections)
}

// GetInspectionByID handles GET /api/v1/inspections/:id
func (h *InspectionHandler) GetInspectionByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	insp, err := h.usecase.GetInspectionByID(ctx, id)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch inspection", err)
		return
	}

	if insp == nil {
		response.NotFound(c, "Inspection not found")
		return
	}

	response.Success(c, insp)
}

// CreateInspection handles POST /api/v1/inspections
func (h *InspectionHandler) CreateInspection(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.CreateInspectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Validate inspection type
	if !entity.IsValidInspectionType(req.InspectionType) {
		response.BadRequest(c, "Invalid inspection_type. Valid values: routine, preventive, corrective, emergency")
		return
	}

	// Validate status if provided
	if req.Status != "" && !entity.IsValidInspectionStatus(req.Status) {
		response.BadRequest(c, "Invalid status. Valid values: scheduled, in_progress, completed, cancelled")
		return
	}

	insp, err := h.usecase.CreateInspection(ctx, &req)
	if err != nil {
		if err.Error() == "contract not found" || err.Error() == "inspector not found" {
			response.NotFound(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to create inspection", err)
		return
	}

	response.Created(c, insp)
}

// UpdateInspection handles PUT /api/v1/inspections/:id
func (h *InspectionHandler) UpdateInspection(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req entity.UpdateInspectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Validate inspection type if provided
	if req.InspectionType != nil && !entity.IsValidInspectionType(*req.InspectionType) {
		response.BadRequest(c, "Invalid inspection_type. Valid values: routine, preventive, corrective, emergency")
		return
	}

	// Validate status if provided
	if req.Status != nil && !entity.IsValidInspectionStatus(*req.Status) {
		response.BadRequest(c, "Invalid status. Valid values: scheduled, in_progress, completed, cancelled")
		return
	}

	insp, err := h.usecase.UpdateInspection(ctx, id, &req)
	if err != nil {
		if err.Error() == "inspection not found" {
			response.NotFound(c, "Inspection not found")
			return
		}
		if err.Error() == "contract not found" || err.Error() == "inspector not found" {
			response.NotFound(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to update inspection", err)
		return
	}

	response.Success(c, insp)
}

// DeleteInspection handles DELETE /api/v1/inspections/:id
func (h *InspectionHandler) DeleteInspection(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	err := h.usecase.DeleteInspection(ctx, id)
	if err != nil {
		if err.Error() == "inspection not found" {
			response.NotFound(c, "Inspection not found")
			return
		}
		response.SafeInternalError(c, "Failed to delete inspection", err)
		return
	}

	response.Success(c, map[string]string{"message": "Inspection deleted successfully"})
}

// GetScheduledInspections handles GET /api/v1/inspections/scheduled
func (h *InspectionHandler) GetScheduledInspections(c *gin.Context) {
	ctx := c.Request.Context()

	inspections, err := h.usecase.GetScheduledInspections(ctx)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch scheduled inspections", err)
		return
	}

	response.Success(c, inspections)
}

// GetInspectionsByContract handles GET /api/v1/contracts/:id/inspections
func (h *InspectionHandler) GetInspectionsByContract(c *gin.Context) {
	ctx := c.Request.Context()
	contractID := c.Param("id")

	inspections, err := h.usecase.GetInspectionsByContract(ctx, contractID)
	if err != nil {
		if err.Error() == "contract not found" {
			response.NotFound(c, "Contract not found")
			return
		}
		response.SafeInternalError(c, "Failed to fetch inspections", err)
		return
	}

	response.Success(c, inspections)
}

// GetInspectionsByInspector handles GET /api/v1/inspectors/:id/inspections
func (h *InspectionHandler) GetInspectionsByInspector(c *gin.Context) {
	ctx := c.Request.Context()
	inspectorID := c.Param("id")

	inspections, err := h.usecase.GetInspectionsByInspector(ctx, inspectorID)
	if err != nil {
		if err.Error() == "inspector not found" {
			response.NotFound(c, "Inspector not found")
			return
		}
		response.SafeInternalError(c, "Failed to fetch inspections", err)
		return
	}

	response.Success(c, inspections)
}
