package task

import (
	"context"
	"errors"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/google/uuid"
)

// UseCase defines the task use case interface
type UseCase interface {
	ListTasks(ctx context.Context, filter *entity.TaskFilter) ([]entity.Task, error)
	GetTaskByID(ctx context.Context, id string) (*entity.Task, error)
	CreateTask(ctx context.Context, req *entity.CreateTaskRequest) (*entity.Task, error)
	UpdateTask(ctx context.Context, id string, req *entity.UpdateTaskRequest) (*entity.Task, error)
	UpdateTaskStatus(ctx context.Context, id string, req *entity.UpdateTaskStatusRequest) (*entity.Task, error)
	DeleteTask(ctx context.Context, id string) error
	GetTasksByContract(ctx context.Context, contractID string) ([]entity.Task, error)
	GetTasksByAssignee(ctx context.Context, assigneeID string) ([]entity.Task, error)
	GetOverdueTasks(ctx context.Context) ([]entity.Task, error)
}

type taskUseCase struct {
	repo         repository.TaskRepository
	contratoRepo repository.ContratoRepository
	gestorRepo   repository.GestorRepository
}

// NewUseCase creates a new task use case
func NewUseCase(
	repo repository.TaskRepository,
	contratoRepo repository.ContratoRepository,
	gestorRepo repository.GestorRepository,
) UseCase {
	return &taskUseCase{
		repo:         repo,
		contratoRepo: contratoRepo,
		gestorRepo:   gestorRepo,
	}
}

// ListTasks returns all tasks with optional filters
func (uc *taskUseCase) ListTasks(ctx context.Context, filter *entity.TaskFilter) ([]entity.Task, error) {
	return uc.repo.FindAll(ctx, filter)
}

// GetTaskByID returns a specific task by ID
func (uc *taskUseCase) GetTaskByID(ctx context.Context, id string) (*entity.Task, error) {
	return uc.repo.FindByID(ctx, id)
}

// CreateTask creates a new task
func (uc *taskUseCase) CreateTask(ctx context.Context, req *entity.CreateTaskRequest) (*entity.Task, error) {
	// Validate contract exists if provided
	if req.ContractID != nil && *req.ContractID != "" {
		contrato, err := uc.contratoRepo.FindByID(ctx, *req.ContractID)
		if err != nil {
			return nil, err
		}
		if contrato == nil {
			return nil, errors.New("contract not found")
		}
	}

	// Validate creator exists
	creator, err := uc.gestorRepo.FindByID(ctx, req.CreatedBy)
	if err != nil {
		return nil, err
	}
	if creator == nil {
		return nil, errors.New("creator not found")
	}

	// Validate assignee exists if provided
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		assignee, err := uc.gestorRepo.FindByID(ctx, *req.AssignedTo)
		if err != nil {
			return nil, err
		}
		if assignee == nil {
			return nil, errors.New("assignee not found")
		}
	}

	// Set default status and priority
	status := req.Status
	if status == "" {
		status = entity.TaskStatusPending
	}
	if !entity.ValidTaskStatus(status) {
		return nil, errors.New("invalid status")
	}

	priority := req.Priority
	if priority == "" {
		priority = entity.TaskPriorityMedium
	}
	if !entity.ValidTaskPriority(priority) {
		return nil, errors.New("invalid priority")
	}

	// Create task entity
	task := &entity.Task{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Priority:    priority,
		DueDate:     req.DueDate,
		ContractID:  req.ContractID,
		AssignedTo:  req.AssignedTo,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   time.Now(),
	}

	if err := uc.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	// Fetch the full task with joined names
	return uc.repo.FindByID(ctx, task.ID)
}

// UpdateTask updates an existing task
func (uc *taskUseCase) UpdateTask(ctx context.Context, id string, req *entity.UpdateTaskRequest) (*entity.Task, error) {
	// Verify task exists
	task, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}

	// Validate contract if being updated
	if req.ContractID != nil {
		if *req.ContractID != "" {
			contrato, err := uc.contratoRepo.FindByID(ctx, *req.ContractID)
			if err != nil {
				return nil, err
			}
			if contrato == nil {
				return nil, errors.New("contract not found")
			}
		}
		task.ContractID = req.ContractID
	}

	// Validate assignee if being updated
	if req.AssignedTo != nil {
		if *req.AssignedTo != "" {
			assignee, err := uc.gestorRepo.FindByID(ctx, *req.AssignedTo)
			if err != nil {
				return nil, err
			}
			if assignee == nil {
				return nil, errors.New("assignee not found")
			}
		}
		task.AssignedTo = req.AssignedTo
	}

	// Update fields if provided
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = req.Description
	}
	if req.Status != nil {
		if !entity.ValidTaskStatus(*req.Status) {
			return nil, errors.New("invalid status")
		}
		task.Status = *req.Status

		// Set completed_at if completing
		if *req.Status == entity.TaskStatusCompleted {
			now := time.Now()
			task.CompletedAt = &now
		} else if *req.Status == entity.TaskStatusPending || *req.Status == entity.TaskStatusInProgress {
			task.CompletedAt = nil
		}
	}
	if req.Priority != nil {
		if !entity.ValidTaskPriority(*req.Priority) {
			return nil, errors.New("invalid priority")
		}
		task.Priority = *req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	// Set updated timestamp
	now := time.Now()
	task.UpdatedAt = &now

	if err := uc.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	// Fetch the full task with joined names
	return uc.repo.FindByID(ctx, task.ID)
}

// UpdateTaskStatus updates only the status of a task
func (uc *taskUseCase) UpdateTaskStatus(ctx context.Context, id string, req *entity.UpdateTaskStatusRequest) (*entity.Task, error) {
	// Verify task exists
	task, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}

	// Validate status
	if !entity.ValidTaskStatus(req.Status) {
		return nil, errors.New("invalid status")
	}

	if err := uc.repo.UpdateStatus(ctx, id, req.Status, nil); err != nil {
		return nil, err
	}

	// Fetch the updated task
	return uc.repo.FindByID(ctx, id)
}

// DeleteTask deletes a task by ID
func (uc *taskUseCase) DeleteTask(ctx context.Context, id string) error {
	// Verify task exists
	task, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}

	return uc.repo.Delete(ctx, id)
}

// GetTasksByContract returns all tasks for a specific contract
func (uc *taskUseCase) GetTasksByContract(ctx context.Context, contractID string) ([]entity.Task, error) {
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

// GetTasksByAssignee returns all tasks assigned to a specific user
func (uc *taskUseCase) GetTasksByAssignee(ctx context.Context, assigneeID string) ([]entity.Task, error) {
	// Verify assignee exists
	assignee, err := uc.gestorRepo.FindByID(ctx, assigneeID)
	if err != nil {
		return nil, err
	}
	if assignee == nil {
		return nil, errors.New("assignee not found")
	}

	return uc.repo.FindByAssignee(ctx, assigneeID)
}

// GetOverdueTasks returns all tasks that are overdue
func (uc *taskUseCase) GetOverdueTasks(ctx context.Context) ([]entity.Task, error) {
	return uc.repo.FindOverdue(ctx)
}
