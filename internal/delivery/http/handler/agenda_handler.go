package handler

import (
	"fmt"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/usecase/agenda"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// AgendaHandler handles agenda/calendar-related HTTP requests
type AgendaHandler struct {
	usecase agenda.UseCase
}

// NewAgendaHandler creates a new agenda handler
func NewAgendaHandler(uc agenda.UseCase) *AgendaHandler {
	return &AgendaHandler{usecase: uc}
}

// ListEvents handles GET /api/v1/agenda
// Query parameters: start_date, end_date, contract_id, user_id, event_type
func (h *AgendaHandler) ListEvents(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse query parameters
	filter := &entity.AgendaFilter{}
	hasFilter := false

	// Parse start_date
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := parseDateTime(startDateStr)
		if err != nil {
			response.BadRequest(c, "Invalid start_date format. Use ISO 8601 format (YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS)")
			return
		}
		filter.StartDate = &startDate
		hasFilter = true
	}

	// Parse end_date
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := parseDateTime(endDateStr)
		if err != nil {
			response.BadRequest(c, "Invalid end_date format. Use ISO 8601 format (YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS)")
			return
		}
		// If only date is provided, set to end of day
		if len(endDateStr) == 10 {
			endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
		filter.EndDate = &endDate
		hasFilter = true
	}

	// Parse contract_id
	if contractID := c.Query("contract_id"); contractID != "" {
		filter.ContractID = &contractID
		hasFilter = true
	}

	// Parse user_id
	if userID := c.Query("user_id"); userID != "" {
		filter.UserID = &userID
		hasFilter = true
	}

	// Parse event_type
	if eventTypeStr := c.Query("event_type"); eventTypeStr != "" {
		eventType := entity.EventType(eventTypeStr)
		if !entity.IsValidEventType(eventType) {
			response.BadRequest(c, "Invalid event_type. Must be one of: audit, inspection, meeting, task, other")
			return
		}
		filter.EventType = &eventType
		hasFilter = true
	}

	var events []entity.AgendaEvent
	var err error

	if hasFilter {
		events, err = h.usecase.ListEvents(ctx, filter)
	} else {
		events, err = h.usecase.ListEvents(ctx, nil)
	}

	if err != nil {
		response.SafeInternalError(c, "Failed to fetch events", err)
		return
	}

	response.Success(c, events)
}

// GetEventByID handles GET /api/v1/agenda/:id
func (h *AgendaHandler) GetEventByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	event, err := h.usecase.GetEventByID(ctx, id)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch event", err)
		return
	}

	if event == nil {
		response.NotFound(c, "Event not found")
		return
	}

	response.Success(c, event)
}

// CreateEvent handles POST /api/v1/agenda
func (h *AgendaHandler) CreateEvent(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	event, err := h.usecase.CreateEvent(ctx, &req)
	if err != nil {
		if err.Error() == "invalid event type" {
			response.BadRequest(c, "Invalid event_type. Must be one of: audit, inspection, meeting, task, other")
			return
		}
		if err.Error() == "contract not found" {
			response.BadRequest(c, "Contract not found")
			return
		}
		if err.Error() == "user not found" {
			response.BadRequest(c, "User not found")
			return
		}
		if err.Error() == "end datetime must be after start datetime" {
			response.BadRequest(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to create event", err)
		return
	}

	response.Created(c, event)
}

// UpdateEvent handles PUT /api/v1/agenda/:id
func (h *AgendaHandler) UpdateEvent(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req entity.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	event, err := h.usecase.UpdateEvent(ctx, id, &req)
	if err != nil {
		if err.Error() == "event not found" {
			response.NotFound(c, "Event not found")
			return
		}
		if err.Error() == "invalid event type" {
			response.BadRequest(c, "Invalid event_type. Must be one of: audit, inspection, meeting, task, other")
			return
		}
		if err.Error() == "contract not found" {
			response.BadRequest(c, "Contract not found")
			return
		}
		if err.Error() == "user not found" {
			response.BadRequest(c, "User not found")
			return
		}
		if err.Error() == "end datetime must be after start datetime" {
			response.BadRequest(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to update event", err)
		return
	}

	response.Success(c, event)
}

// DeleteEvent handles DELETE /api/v1/agenda/:id
func (h *AgendaHandler) DeleteEvent(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	err := h.usecase.DeleteEvent(ctx, id)
	if err != nil {
		if err.Error() == "event not found" {
			response.NotFound(c, "Event not found")
			return
		}
		response.SafeInternalError(c, "Failed to delete event", err)
		return
	}

	response.Success(c, map[string]string{"message": "Event deleted successfully"})
}

// parseDateTime parses a date/time string in various formats
func parseDateTime(s string) (time.Time, error) {
	// Try different formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("could not parse date: %s", s)
}
