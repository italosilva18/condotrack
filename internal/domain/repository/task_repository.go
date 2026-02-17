package repository

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
)

// TaskRepository defines the interface for task data access
type TaskRepository interface {
	// FindAll returns all tasks with optional filters
	FindAll(ctx context.Context, filter *entity.TaskFilter) ([]entity.Task, error)

	// FindByID returns a task by ID
	FindByID(ctx context.Context, id string) (*entity.Task, error)

	// FindByContract returns all tasks for a specific contract
	FindByContract(ctx context.Context, contractID string) ([]entity.Task, error)

	// FindByAssignee returns all tasks assigned to a specific user
	FindByAssignee(ctx context.Context, assigneeID string) ([]entity.Task, error)

	// FindByCreator returns all tasks created by a specific user
	FindByCreator(ctx context.Context, creatorID string) ([]entity.Task, error)

	// FindOverdue returns all tasks that are overdue (past due date and not completed/cancelled)
	FindOverdue(ctx context.Context) ([]entity.Task, error)

	// Create creates a new task
	Create(ctx context.Context, task *entity.Task) error

	// Update updates an existing task
	Update(ctx context.Context, task *entity.Task) error

	// UpdateStatus updates only the status of a task
	UpdateStatus(ctx context.Context, id string, status string, completedAt *interface{}) error

	// Delete deletes a task by ID
	Delete(ctx context.Context, id string) error

	// CountByContract returns the number of tasks for a contract
	CountByContract(ctx context.Context, contractID string) (int, error)

	// CountByAssignee returns the number of tasks for an assignee
	CountByAssignee(ctx context.Context, assigneeID string) (int, error)
}
