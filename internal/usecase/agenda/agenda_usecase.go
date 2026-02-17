package agenda

import (
	"context"
	"errors"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/google/uuid"
)

// UseCase defines the agenda use case interface
type UseCase interface {
	ListEvents(ctx context.Context, filter *entity.AgendaFilter) ([]entity.AgendaEvent, error)
	GetEventByID(ctx context.Context, id string) (*entity.AgendaEvent, error)
	CreateEvent(ctx context.Context, req *entity.CreateEventRequest) (*entity.AgendaEvent, error)
	UpdateEvent(ctx context.Context, id string, req *entity.UpdateEventRequest) (*entity.AgendaEvent, error)
	DeleteEvent(ctx context.Context, id string) error
	GetEventsByContract(ctx context.Context, contractID string) ([]entity.AgendaEvent, error)
	GetEventsByUser(ctx context.Context, userID string) ([]entity.AgendaEvent, error)
	GetEventsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entity.AgendaEvent, error)
}

type agendaUseCase struct {
	repo         repository.AgendaRepository
	contratoRepo repository.ContratoRepository
	gestorRepo   repository.GestorRepository
}

// NewUseCase creates a new agenda use case
func NewUseCase(repo repository.AgendaRepository, contratoRepo repository.ContratoRepository, gestorRepo repository.GestorRepository) UseCase {
	return &agendaUseCase{
		repo:         repo,
		contratoRepo: contratoRepo,
		gestorRepo:   gestorRepo,
	}
}

// ListEvents returns all events with optional filters
func (uc *agendaUseCase) ListEvents(ctx context.Context, filter *entity.AgendaFilter) ([]entity.AgendaEvent, error) {
	if filter != nil && (filter.StartDate != nil || filter.EndDate != nil || filter.ContractID != nil || filter.UserID != nil || filter.EventType != nil) {
		return uc.repo.FindWithFilters(ctx, filter)
	}
	return uc.repo.FindAll(ctx)
}

// GetEventByID returns a specific event by ID
func (uc *agendaUseCase) GetEventByID(ctx context.Context, id string) (*entity.AgendaEvent, error) {
	return uc.repo.FindByID(ctx, id)
}

// CreateEvent creates a new calendar event
func (uc *agendaUseCase) CreateEvent(ctx context.Context, req *entity.CreateEventRequest) (*entity.AgendaEvent, error) {
	// Validate event type
	if !entity.IsValidEventType(req.EventType) {
		return nil, errors.New("invalid event type")
	}

	// Validate date range
	if req.EndDatetime.Before(req.StartDatetime) {
		return nil, errors.New("end datetime must be after start datetime")
	}

	// Validate contract if provided
	if req.ContractID != nil && *req.ContractID != "" {
		contrato, err := uc.contratoRepo.FindByID(ctx, *req.ContractID)
		if err != nil {
			return nil, err
		}
		if contrato == nil {
			return nil, errors.New("contract not found")
		}
	}

	// Validate user if provided
	if req.UserID != nil && *req.UserID != "" {
		gestor, err := uc.gestorRepo.FindByID(ctx, *req.UserID)
		if err != nil {
			return nil, err
		}
		if gestor == nil {
			return nil, errors.New("user not found")
		}
	}

	// Create event entity
	event := &entity.AgendaEvent{
		ID:             uuid.New().String(),
		Title:          req.Title,
		Description:    req.Description,
		EventType:      req.EventType,
		StartDatetime:  req.StartDatetime,
		EndDatetime:    req.EndDatetime,
		AllDay:         req.AllDay,
		Location:       req.Location,
		ContractID:     req.ContractID,
		UserID:         req.UserID,
		RecurrenceRule: req.RecurrenceRule,
		Color:          req.Color,
		CreatedAt:      time.Now(),
	}

	// Set default color based on event type if not provided
	if event.Color == nil {
		defaultColor := getDefaultColorForEventType(event.EventType)
		event.Color = &defaultColor
	}

	if err := uc.repo.Create(ctx, event); err != nil {
		return nil, err
	}

	// Fetch the created event to get contract and user names
	return uc.repo.FindByID(ctx, event.ID)
}

// UpdateEvent updates an existing calendar event
func (uc *agendaUseCase) UpdateEvent(ctx context.Context, id string, req *entity.UpdateEventRequest) (*entity.AgendaEvent, error) {
	// Verify event exists
	event, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, errors.New("event not found")
	}

	// Update fields if provided
	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = req.Description
	}
	if req.EventType != nil {
		if !entity.IsValidEventType(*req.EventType) {
			return nil, errors.New("invalid event type")
		}
		event.EventType = *req.EventType
	}
	if req.StartDatetime != nil {
		event.StartDatetime = *req.StartDatetime
	}
	if req.EndDatetime != nil {
		event.EndDatetime = *req.EndDatetime
	}

	// Validate date range after updates
	if event.EndDatetime.Before(event.StartDatetime) {
		return nil, errors.New("end datetime must be after start datetime")
	}

	if req.AllDay != nil {
		event.AllDay = *req.AllDay
	}
	if req.Location != nil {
		event.Location = req.Location
	}
	if req.ContractID != nil {
		// Validate contract if provided
		if *req.ContractID != "" {
			contrato, err := uc.contratoRepo.FindByID(ctx, *req.ContractID)
			if err != nil {
				return nil, err
			}
			if contrato == nil {
				return nil, errors.New("contract not found")
			}
		}
		event.ContractID = req.ContractID
	}
	if req.UserID != nil {
		// Validate user if provided
		if *req.UserID != "" {
			gestor, err := uc.gestorRepo.FindByID(ctx, *req.UserID)
			if err != nil {
				return nil, err
			}
			if gestor == nil {
				return nil, errors.New("user not found")
			}
		}
		event.UserID = req.UserID
	}
	if req.RecurrenceRule != nil {
		event.RecurrenceRule = req.RecurrenceRule
	}
	if req.Color != nil {
		event.Color = req.Color
	}

	// Set updated timestamp
	now := time.Now()
	event.UpdatedAt = &now

	if err := uc.repo.Update(ctx, event); err != nil {
		return nil, err
	}

	// Fetch the updated event to get contract and user names
	return uc.repo.FindByID(ctx, id)
}

// DeleteEvent deletes an event by ID
func (uc *agendaUseCase) DeleteEvent(ctx context.Context, id string) error {
	// Verify event exists
	event, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("event not found")
	}

	return uc.repo.Delete(ctx, id)
}

// GetEventsByContract returns all events for a specific contract
func (uc *agendaUseCase) GetEventsByContract(ctx context.Context, contractID string) ([]entity.AgendaEvent, error) {
	// Verify contract exists
	contrato, err := uc.contratoRepo.FindByID(ctx, contractID)
	if err != nil {
		return nil, err
	}
	if contrato == nil {
		return nil, errors.New("contract not found")
	}

	return uc.repo.FindByContract(ctx, contractID)
}

// GetEventsByUser returns all events for a specific user
func (uc *agendaUseCase) GetEventsByUser(ctx context.Context, userID string) ([]entity.AgendaEvent, error) {
	// Verify user exists
	gestor, err := uc.gestorRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if gestor == nil {
		return nil, errors.New("user not found")
	}

	return uc.repo.FindByUser(ctx, userID)
}

// GetEventsByDateRange returns all events within a date range
func (uc *agendaUseCase) GetEventsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entity.AgendaEvent, error) {
	if endDate.Before(startDate) {
		return nil, errors.New("end date must be after start date")
	}

	return uc.repo.FindByDateRange(ctx, startDate, endDate)
}

// getDefaultColorForEventType returns a default color for each event type
func getDefaultColorForEventType(eventType entity.EventType) string {
	switch eventType {
	case entity.EventTypeAudit:
		return "#3B82F6" // Blue
	case entity.EventTypeInspection:
		return "#F59E0B" // Amber
	case entity.EventTypeMeeting:
		return "#8B5CF6" // Purple
	case entity.EventTypeTask:
		return "#10B981" // Emerald
	case entity.EventTypeOther:
		return "#6B7280" // Gray
	default:
		return "#6B7280" // Gray
	}
}
