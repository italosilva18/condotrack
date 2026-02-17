package entity

import "time"

// TaskStatus constants define the possible states of a task
const (
	TaskStatusPending    = "pending"
	TaskStatusInProgress = "in_progress"
	TaskStatusCompleted  = "completed"
	TaskStatusCancelled  = "cancelled"
)

// TaskPriority constants define the priority levels of a task
const (
	TaskPriorityLow    = "low"
	TaskPriorityMedium = "medium"
	TaskPriorityHigh   = "high"
	TaskPriorityUrgent = "urgent"
)

// Task represents a task entity
type Task struct {
	ID             string     `db:"id" json:"id"`
	Title          string     `db:"title" json:"title"`
	Description    *string    `db:"description" json:"description,omitempty"`
	Status         string     `db:"status" json:"status"`
	Priority       string     `db:"priority" json:"priority"`
	DueDate        *time.Time `db:"due_date" json:"due_date,omitempty"`
	ContractID     *string    `db:"contract_id" json:"contract_id,omitempty"`
	ContractName   *string    `db:"contract_name" json:"contract_name,omitempty"`
	AssignedTo     *string    `db:"assigned_to" json:"assigned_to,omitempty"`
	AssignedToName *string    `db:"assigned_to_name" json:"assigned_to_name,omitempty"`
	CreatedBy      string     `db:"created_by" json:"created_by"`
	CreatedByName  *string    `db:"created_by_name" json:"created_by_name,omitempty"`
	CompletedAt    *time.Time `db:"completed_at" json:"completed_at,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// TaskWithDetails represents a task with related entity names
type TaskWithDetails struct {
	Task
}

// CreateTaskRequest represents the request to create a task
type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description *string    `json:"description,omitempty"`
	Status      string     `json:"status,omitempty"`
	Priority    string     `json:"priority,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	ContractID  *string    `json:"contract_id,omitempty"`
	AssignedTo  *string    `json:"assigned_to,omitempty"`
	CreatedBy   string     `json:"created_by" binding:"required"`
}

// UpdateTaskRequest represents the request to update a task
type UpdateTaskRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Priority    *string    `json:"priority,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	ContractID  *string    `json:"contract_id,omitempty"`
	AssignedTo  *string    `json:"assigned_to,omitempty"`
}

// UpdateTaskStatusRequest represents the request to update task status only
type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// TaskFilter represents filter options for listing tasks
type TaskFilter struct {
	ContractID *string
	AssignedTo *string
	Status     *string
	Priority   *string
	CreatedBy  *string
}

// ValidTaskStatus checks if the given status is valid
func ValidTaskStatus(status string) bool {
	switch status {
	case TaskStatusPending, TaskStatusInProgress, TaskStatusCompleted, TaskStatusCancelled:
		return true
	}
	return false
}

// ValidTaskPriority checks if the given priority is valid
func ValidTaskPriority(priority string) bool {
	switch priority {
	case TaskPriorityLow, TaskPriorityMedium, TaskPriorityHigh, TaskPriorityUrgent:
		return true
	}
	return false
}
