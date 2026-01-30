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

type taskMySQLRepository struct {
	db *sqlx.DB
}

// NewTaskMySQLRepository creates a new MySQL implementation of TaskRepository
func NewTaskMySQLRepository(db *sqlx.DB) repository.TaskRepository {
	return &taskMySQLRepository{db: db}
}

func (r *taskMySQLRepository) FindAll(ctx context.Context, filter *entity.TaskFilter) ([]entity.Task, error) {
	var tasks []entity.Task

	query := `SELECT t.id, t.title, t.description, t.status, t.priority, t.due_date,
			  t.contract_id, c.nome as contract_name,
			  t.assigned_to, ga.nome as assigned_to_name,
			  t.created_by, gc.nome as created_by_name,
			  t.completed_at, t.created_at, t.updated_at
			  FROM tasks t
			  LEFT JOIN contratos c ON c.id = t.contract_id
			  LEFT JOIN gestores ga ON ga.id = t.assigned_to
			  LEFT JOIN gestores gc ON gc.id = t.created_by`

	var conditions []string
	var args []interface{}

	if filter != nil {
		if filter.ContractID != nil && *filter.ContractID != "" {
			conditions = append(conditions, "t.contract_id = ?")
			args = append(args, *filter.ContractID)
		}
		if filter.AssignedTo != nil && *filter.AssignedTo != "" {
			conditions = append(conditions, "t.assigned_to = ?")
			args = append(args, *filter.AssignedTo)
		}
		if filter.Status != nil && *filter.Status != "" {
			conditions = append(conditions, "t.status = ?")
			args = append(args, *filter.Status)
		}
		if filter.Priority != nil && *filter.Priority != "" {
			conditions = append(conditions, "t.priority = ?")
			args = append(args, *filter.Priority)
		}
		if filter.CreatedBy != nil && *filter.CreatedBy != "" {
			conditions = append(conditions, "t.created_by = ?")
			args = append(args, *filter.CreatedBy)
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY CASE t.priority WHEN 'urgent' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 END, t.due_date ASC, t.created_at DESC"

	err := r.db.SelectContext(ctx, &tasks, query, args...)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskMySQLRepository) FindByID(ctx context.Context, id string) (*entity.Task, error) {
	var task entity.Task
	query := `SELECT t.id, t.title, t.description, t.status, t.priority, t.due_date,
			  t.contract_id, c.nome as contract_name,
			  t.assigned_to, ga.nome as assigned_to_name,
			  t.created_by, gc.nome as created_by_name,
			  t.completed_at, t.created_at, t.updated_at
			  FROM tasks t
			  LEFT JOIN contratos c ON c.id = t.contract_id
			  LEFT JOIN gestores ga ON ga.id = t.assigned_to
			  LEFT JOIN gestores gc ON gc.id = t.created_by
			  WHERE t.id = ?`
	err := r.db.GetContext(ctx, &task, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *taskMySQLRepository) FindByContract(ctx context.Context, contractID string) ([]entity.Task, error) {
	var tasks []entity.Task
	query := `SELECT t.id, t.title, t.description, t.status, t.priority, t.due_date,
			  t.contract_id, c.nome as contract_name,
			  t.assigned_to, ga.nome as assigned_to_name,
			  t.created_by, gc.nome as created_by_name,
			  t.completed_at, t.created_at, t.updated_at
			  FROM tasks t
			  LEFT JOIN contratos c ON c.id = t.contract_id
			  LEFT JOIN gestores ga ON ga.id = t.assigned_to
			  LEFT JOIN gestores gc ON gc.id = t.created_by
			  WHERE t.contract_id = ?
			  ORDER BY CASE t.priority WHEN 'urgent' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 END, t.due_date ASC`
	err := r.db.SelectContext(ctx, &tasks, query, contractID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskMySQLRepository) FindByAssignee(ctx context.Context, assigneeID string) ([]entity.Task, error) {
	var tasks []entity.Task
	query := `SELECT t.id, t.title, t.description, t.status, t.priority, t.due_date,
			  t.contract_id, c.nome as contract_name,
			  t.assigned_to, ga.nome as assigned_to_name,
			  t.created_by, gc.nome as created_by_name,
			  t.completed_at, t.created_at, t.updated_at
			  FROM tasks t
			  LEFT JOIN contratos c ON c.id = t.contract_id
			  LEFT JOIN gestores ga ON ga.id = t.assigned_to
			  LEFT JOIN gestores gc ON gc.id = t.created_by
			  WHERE t.assigned_to = ?
			  ORDER BY CASE t.priority WHEN 'urgent' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 END, t.due_date ASC`
	err := r.db.SelectContext(ctx, &tasks, query, assigneeID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskMySQLRepository) FindByCreator(ctx context.Context, creatorID string) ([]entity.Task, error) {
	var tasks []entity.Task
	query := `SELECT t.id, t.title, t.description, t.status, t.priority, t.due_date,
			  t.contract_id, c.nome as contract_name,
			  t.assigned_to, ga.nome as assigned_to_name,
			  t.created_by, gc.nome as created_by_name,
			  t.completed_at, t.created_at, t.updated_at
			  FROM tasks t
			  LEFT JOIN contratos c ON c.id = t.contract_id
			  LEFT JOIN gestores ga ON ga.id = t.assigned_to
			  LEFT JOIN gestores gc ON gc.id = t.created_by
			  WHERE t.created_by = ?
			  ORDER BY CASE t.priority WHEN 'urgent' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 END, t.due_date ASC`
	err := r.db.SelectContext(ctx, &tasks, query, creatorID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskMySQLRepository) FindOverdue(ctx context.Context) ([]entity.Task, error) {
	var tasks []entity.Task
	query := `SELECT t.id, t.title, t.description, t.status, t.priority, t.due_date,
			  t.contract_id, c.nome as contract_name,
			  t.assigned_to, ga.nome as assigned_to_name,
			  t.created_by, gc.nome as created_by_name,
			  t.completed_at, t.created_at, t.updated_at
			  FROM tasks t
			  LEFT JOIN contratos c ON c.id = t.contract_id
			  LEFT JOIN gestores ga ON ga.id = t.assigned_to
			  LEFT JOIN gestores gc ON gc.id = t.created_by
			  WHERE t.due_date < NOW()
			  AND t.status NOT IN ('completed', 'cancelled')
			  ORDER BY t.due_date ASC, CASE t.priority WHEN 'urgent' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 END`
	err := r.db.SelectContext(ctx, &tasks, query)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskMySQLRepository) Create(ctx context.Context, task *entity.Task) error {
	query := `INSERT INTO tasks (id, title, description, status, priority, due_date,
			  contract_id, assigned_to, created_by, completed_at, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query,
		task.ID, task.Title, task.Description, task.Status, task.Priority, task.DueDate,
		task.ContractID, task.AssignedTo, task.CreatedBy, task.CompletedAt)
	return err
}

func (r *taskMySQLRepository) Update(ctx context.Context, task *entity.Task) error {
	query := `UPDATE tasks
			  SET title = ?, description = ?, status = ?, priority = ?, due_date = ?,
			  contract_id = ?, assigned_to = ?, completed_at = ?, updated_at = NOW()
			  WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		task.Title, task.Description, task.Status, task.Priority, task.DueDate,
		task.ContractID, task.AssignedTo, task.CompletedAt, task.ID)
	return err
}

func (r *taskMySQLRepository) UpdateStatus(ctx context.Context, id string, status string, completedAt *interface{}) error {
	var query string
	var args []interface{}

	if status == entity.TaskStatusCompleted {
		now := time.Now()
		query = `UPDATE tasks SET status = ?, completed_at = ?, updated_at = NOW() WHERE id = ?`
		args = []interface{}{status, now, id}
	} else if status == entity.TaskStatusPending || status == entity.TaskStatusInProgress {
		// Clear completed_at when moving back to pending or in_progress
		query = `UPDATE tasks SET status = ?, completed_at = NULL, updated_at = NOW() WHERE id = ?`
		args = []interface{}{status, id}
	} else {
		query = `UPDATE tasks SET status = ?, updated_at = NOW() WHERE id = ?`
		args = []interface{}{status, id}
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *taskMySQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *taskMySQLRepository) CountByContract(ctx context.Context, contractID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM tasks WHERE contract_id = ?`
	err := r.db.GetContext(ctx, &count, query, contractID)
	return count, err
}

func (r *taskMySQLRepository) CountByAssignee(ctx context.Context, assigneeID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM tasks WHERE assigned_to = ?`
	err := r.db.GetContext(ctx, &count, query, assigneeID)
	return count, err
}
