package handler

import (
	"strings"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/usecase/team"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// TeamHandler handles team-related HTTP requests
type TeamHandler struct {
	usecase team.UseCase
}

// NewTeamHandler creates a new team handler
func NewTeamHandler(uc team.UseCase) *TeamHandler {
	return &TeamHandler{usecase: uc}
}

// ListTeamMembers handles GET /api/v1/team
func (h *TeamHandler) ListTeamMembers(c *gin.Context) {
	ctx := c.Request.Context()

	// Build filter from query parameters
	filter := &entity.TeamMemberFilter{}

	if contractID := c.Query("contract_id"); contractID != "" {
		filter.ContractID = &contractID
	}

	if userID := c.Query("user_id"); userID != "" {
		filter.UserID = &userID
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true" || isActiveStr == "1"
		filter.IsActive = &isActive
	}

	members, err := h.usecase.ListTeamMembers(ctx, filter)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch team members", err)
		return
	}

	response.Success(c, members)
}

// GetTeamMemberByID handles GET /api/v1/team/:id
func (h *TeamHandler) GetTeamMemberByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	member, err := h.usecase.GetTeamMemberByID(ctx, id)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch team member", err)
		return
	}

	if member == nil {
		response.NotFound(c, "Team member not found")
		return
	}

	response.Success(c, member)
}

// CreateTeamMember handles POST /api/v1/team
func (h *TeamHandler) CreateTeamMember(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.CreateTeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	member, err := h.usecase.CreateTeamMember(ctx, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, err.Error())
			return
		}
		if strings.Contains(err.Error(), "already assigned") {
			response.BadRequest(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to create team member", err)
		return
	}

	response.Created(c, member)
}

// UpdateTeamMember handles PUT /api/v1/team/:id
func (h *TeamHandler) UpdateTeamMember(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req entity.UpdateTeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	member, err := h.usecase.UpdateTeamMember(ctx, id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to update team member", err)
		return
	}

	response.Success(c, member)
}

// DeleteTeamMember handles DELETE /api/v1/team/:id
func (h *TeamHandler) DeleteTeamMember(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	err := h.usecase.DeleteTeamMember(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to delete team member", err)
		return
	}

	response.SuccessWithMessage(c, "Team member removed successfully", nil)
}

// GetTeamByContract handles GET /api/v1/team/contract/:id
func (h *TeamHandler) GetTeamByContract(c *gin.Context) {
	ctx := c.Request.Context()
	contractID := c.Param("id")

	members, err := h.usecase.GetTeamByContract(ctx, contractID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to fetch team for contract", err)
		return
	}

	response.Success(c, members)
}

// GetContractsByUser handles GET /api/v1/team/user/:id
func (h *TeamHandler) GetContractsByUser(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("id")

	members, err := h.usecase.GetContractsByUser(ctx, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to fetch contracts for user", err)
		return
	}

	response.Success(c, members)
}
