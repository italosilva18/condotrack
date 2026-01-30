package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type agendaMySQLRepository struct {
	db *sqlx.DB
}

// NewAgendaMySQLRepository creates a new MySQL implementation of AgendaRepository
func NewAgendaMySQLRepository(db *sqlx.DB) repository.AgendaRepository {
	return &agendaMySQLRepository{db: db}
}

func (r *agendaMySQLRepository) FindAll(ctx context.Context) ([]entity.AgendaEvent, error) {
	var events []entity.AgendaEvent
	query := `SELECT
			  e.id, e.title, e.description, e.event_type, e.start_datetime, e.end_datetime,
			  e.all_day, e.location, e.contract_id, e.user_id, e.recurrence_rule, e.color,
			  e.created_at, e.updated_at,
			  c.nome as contract_name,
			  g.nome as user_name
			  FROM agenda e
			  LEFT JOIN contratos c ON c.id = e.contract_id
			  LEFT JOIN gestores g ON g.id = e.user_id
			  ORDER BY e.start_datetime ASC`
	err := r.db.SelectContext(ctx, &events, query)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *agendaMySQLRepository) FindByID(ctx context.Context, id string) (*entity.AgendaEvent, error) {
	var event entity.AgendaEvent
	query := `SELECT
			  e.id, e.title, e.description, e.event_type, e.start_datetime, e.end_datetime,
			  e.all_day, e.location, e.contract_id, e.user_id, e.recurrence_rule, e.color,
			  e.created_at, e.updated_at,
			  c.nome as contract_name,
			  g.nome as user_name
			  FROM agenda e
			  LEFT JOIN contratos c ON c.id = e.contract_id
			  LEFT JOIN gestores g ON g.id = e.user_id
			  WHERE e.id = ?`
	err := r.db.GetContext(ctx, &event, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (r *agendaMySQLRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entity.AgendaEvent, error) {
	var events []entity.AgendaEvent
	query := `SELECT
			  e.id, e.title, e.description, e.event_type, e.start_datetime, e.end_datetime,
			  e.all_day, e.location, e.contract_id, e.user_id, e.recurrence_rule, e.color,
			  e.created_at, e.updated_at,
			  c.nome as contract_name,
			  g.nome as user_name
			  FROM agenda e
			  LEFT JOIN contratos c ON c.id = e.contract_id
			  LEFT JOIN gestores g ON g.id = e.user_id
			  WHERE (e.start_datetime BETWEEN ? AND ?) OR (e.end_datetime BETWEEN ? AND ?)
			  OR (e.start_datetime <= ? AND e.end_datetime >= ?)
			  ORDER BY e.start_datetime ASC`
	err := r.db.SelectContext(ctx, &events, query, startDate, endDate, startDate, endDate, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *agendaMySQLRepository) FindByContract(ctx context.Context, contractID string) ([]entity.AgendaEvent, error) {
	var events []entity.AgendaEvent
	query := `SELECT
			  e.id, e.title, e.description, e.event_type, e.start_datetime, e.end_datetime,
			  e.all_day, e.location, e.contract_id, e.user_id, e.recurrence_rule, e.color,
			  e.created_at, e.updated_at,
			  c.nome as contract_name,
			  g.nome as user_name
			  FROM agenda e
			  LEFT JOIN contratos c ON c.id = e.contract_id
			  LEFT JOIN gestores g ON g.id = e.user_id
			  WHERE e.contract_id = ?
			  ORDER BY e.start_datetime ASC`
	err := r.db.SelectContext(ctx, &events, query, contractID)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *agendaMySQLRepository) FindByUser(ctx context.Context, userID string) ([]entity.AgendaEvent, error) {
	var events []entity.AgendaEvent
	query := `SELECT
			  e.id, e.title, e.description, e.event_type, e.start_datetime, e.end_datetime,
			  e.all_day, e.location, e.contract_id, e.user_id, e.recurrence_rule, e.color,
			  e.created_at, e.updated_at,
			  c.nome as contract_name,
			  g.nome as user_name
			  FROM agenda e
			  LEFT JOIN contratos c ON c.id = e.contract_id
			  LEFT JOIN gestores g ON g.id = e.user_id
			  WHERE e.user_id = ?
			  ORDER BY e.start_datetime ASC`
	err := r.db.SelectContext(ctx, &events, query, userID)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *agendaMySQLRepository) FindWithFilters(ctx context.Context, filter *entity.AgendaFilter) ([]entity.AgendaEvent, error) {
	var events []entity.AgendaEvent
	var conditions []string
	var args []interface{}

	baseQuery := `SELECT
			  e.id, e.title, e.description, e.event_type, e.start_datetime, e.end_datetime,
			  e.all_day, e.location, e.contract_id, e.user_id, e.recurrence_rule, e.color,
			  e.created_at, e.updated_at,
			  c.nome as contract_name,
			  g.nome as user_name
			  FROM agenda e
			  LEFT JOIN contratos c ON c.id = e.contract_id
			  LEFT JOIN gestores g ON g.id = e.user_id`

	if filter != nil {
		if filter.StartDate != nil && filter.EndDate != nil {
			conditions = append(conditions, "((e.start_datetime BETWEEN ? AND ?) OR (e.end_datetime BETWEEN ? AND ?) OR (e.start_datetime <= ? AND e.end_datetime >= ?))")
			args = append(args, *filter.StartDate, *filter.EndDate, *filter.StartDate, *filter.EndDate, *filter.StartDate, *filter.EndDate)
		} else if filter.StartDate != nil {
			conditions = append(conditions, "e.start_datetime >= ?")
			args = append(args, *filter.StartDate)
		} else if filter.EndDate != nil {
			conditions = append(conditions, "e.end_datetime <= ?")
			args = append(args, *filter.EndDate)
		}

		if filter.ContractID != nil {
			conditions = append(conditions, "e.contract_id = ?")
			args = append(args, *filter.ContractID)
		}

		if filter.UserID != nil {
			conditions = append(conditions, "e.user_id = ?")
			args = append(args, *filter.UserID)
		}

		if filter.EventType != nil {
			conditions = append(conditions, "e.event_type = ?")
			args = append(args, *filter.EventType)
		}
	}

	query := baseQuery
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY e.start_datetime ASC"

	err := r.db.SelectContext(ctx, &events, query, args...)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *agendaMySQLRepository) Create(ctx context.Context, event *entity.AgendaEvent) error {
	query := `INSERT INTO agenda (id, title, description, event_type, start_datetime, end_datetime,
			  all_day, location, contract_id, user_id, recurrence_rule, color, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query,
		event.ID, event.Title, event.Description, event.EventType, event.StartDatetime, event.EndDatetime,
		event.AllDay, event.Location, event.ContractID, event.UserID, event.RecurrenceRule, event.Color)
	return err
}

func (r *agendaMySQLRepository) Update(ctx context.Context, event *entity.AgendaEvent) error {
	query := `UPDATE agenda
			  SET title = ?, description = ?, event_type = ?, start_datetime = ?, end_datetime = ?,
			  all_day = ?, location = ?, contract_id = ?, user_id = ?, recurrence_rule = ?, color = ?, updated_at = NOW()
			  WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		event.Title, event.Description, event.EventType, event.StartDatetime, event.EndDatetime,
		event.AllDay, event.Location, event.ContractID, event.UserID, event.RecurrenceRule, event.Color, event.ID)
	return err
}

func (r *agendaMySQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM agenda WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
