package handler

import (
	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/usecase/task"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// TaskHandler handles task-related HTTP requests
type TaskHandler struct {
	usecase task.UseCase
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(uc task.UseCase) *TaskHandler {
	return &TaskHandler{usecase: uc}
}

// ListTasks handles GET /api/v1/tasks
func (h *TaskHandler) ListTasks(c *gin.Context) {
	ctx := c.Request.Context()

	// Build filter from query parameters
	filter := &entity.TaskFilter{}

	if contractID := c.Query("contract_id"); contractID != "" {
		filter.ContractID = &contractID
	}
	if assignedTo := c.Query("assigned_to"); assignedTo != "" {
		filter.AssignedTo = &assignedTo
	}
	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}
	if priority := c.Query("priority"); priority != "" {
		filter.Priority = &priority
	}
	if createdBy := c.Query("created_by"); createdBy != "" {
		filter.CreatedBy = &createdBy
	}

	tasks, err := h.usecase.ListTasks(ctx, filter)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch tasks", err)
		return
	}

	response.Success(c, tasks)
}

// GetTaskByID handles GET /api/v1/tasks/:id
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	task, err := h.usecase.GetTaskByID(ctx, id)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch task", err)
		return
	}

	if task == nil {
		response.NotFound(c, "Task not found")
		return
	}

	response.Success(c, task)
}

// CreateTask handles POST /api/v1/tasks
func (h *TaskHandler) CreateTask(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	task, err := h.usecase.CreateTask(ctx, &req)
	if err != nil {
		if err.Error() == "contract not found" || err.Error() == "creator not found" || err.Error() == "assignee not found" {
			response.BadRequest(c, err.Error())
			return
		}
		if err.Error() == "invalid status" || err.Error() == "invalid priority" {
			response.BadRequest(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to create task", err)
		return
	}

	response.Created(c, task)
}

// UpdateTask handles PUT /api/v1/tasks/:id
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req entity.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	task, err := h.usecase.UpdateTask(ctx, id, &req)
	if err != nil {
		if err.Error() == "task not found" {
			response.NotFound(c, "Task not found")
			return
		}
		if err.Error() == "contract not found" || err.Error() == "assignee not found" {
			response.BadRequest(c, err.Error())
			return
		}
		if err.Error() == "invalid status" || err.Error() == "invalid priority" {
			response.BadRequest(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to update task", err)
		return
	}

	response.Success(c, task)
}

// UpdateTaskStatus handles PATCH /api/v1/tasks/:id/status
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req entity.UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	task, err := h.usecase.UpdateTaskStatus(ctx, id, &req)
	if err != nil {
		if err.Error() == "task not found" {
			response.NotFound(c, "Task not found")
			return
		}
		if err.Error() == "invalid status" {
			response.BadRequest(c, err.Error())
			return
		}
		response.SafeInternalError(c, "Failed to update task status", err)
		return
	}

	response.Success(c, task)
}

// DeleteTask handles DELETE /api/v1/tasks/:id
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	err := h.usecase.DeleteTask(ctx, id)
	if err != nil {
		if err.Error() == "task not found" {
			response.NotFound(c, "Task not found")
			return
		}
		response.SafeInternalError(c, "Failed to delete task", err)
		return
	}

	response.Success(c, map[string]string{"message": "Task deleted successfully"})
}

// GetOverdueTasks handles GET /api/v1/tasks/overdue
func (h *TaskHandler) GetOverdueTasks(c *gin.Context) {
	ctx := c.Request.Context()

	tasks, err := h.usecase.GetOverdueTasks(ctx)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch overdue tasks", err)
		return
	}

	response.Success(c, tasks)
}

// GetTasksByContract handles GET /api/v1/tasks/contract/:id
func (h *TaskHandler) GetTasksByContract(c *gin.Context) {
	ctx := c.Request.Context()
	contractID := c.Param("id")

	tasks, err := h.usecase.GetTasksByContract(ctx, contractID)
	if err != nil {
		if err.Error() == "contract not found" {
			response.NotFound(c, "Contract not found")
			return
		}
		response.SafeInternalError(c, "Failed to fetch tasks", err)
		return
	}

	response.Success(c, tasks)
}

// GetTasksByAssignee handles GET /api/v1/tasks/assignee/:id
func (h *TaskHandler) GetTasksByAssignee(c *gin.Context) {
	ctx := c.Request.Context()
	assigneeID := c.Param("id")

	tasks, err := h.usecase.GetTasksByAssignee(ctx, assigneeID)
	if err != nil {
		if err.Error() == "assignee not found" {
			response.NotFound(c, "Assignee not found")
			return
		}
		response.SafeInternalError(c, "Failed to fetch tasks", err)
		return
	}

	response.Success(c, tasks)
}
