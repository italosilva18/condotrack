package handler

import (
	"strconv"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/usecase/course"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// CourseHandler handles course-related HTTP requests
type CourseHandler struct {
	usecase course.UseCase
}

// NewCourseHandler creates a new course handler
func NewCourseHandler(uc course.UseCase) *CourseHandler {
	return &CourseHandler{usecase: uc}
}

// ListCourses handles GET /api/v1/courses
func (h *CourseHandler) ListCourses(c *gin.Context) {
	ctx := c.Request.Context()

	// Check for instructor_id filter
	instructorID := c.Query("instructor_id")
	if instructorID != "" {
		courses, err := h.usecase.ListCoursesByInstructor(ctx, instructorID)
		if err != nil {
			response.SafeInternalError(c, "Failed to fetch courses", err)
			return
		}
		response.Success(c, courses)
		return
	}

	// Check for active filter
	activeOnly := c.Query("active")

	// Parse pagination parameters
	page := 1
	perPage := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if pp := c.Query("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 {
			if parsed > 100 {
				parsed = 100
			}
			perPage = parsed
		}
	}

	var result *entity.CourseListResponse
	var err error

	if activeOnly == "true" || activeOnly == "1" {
		result, err = h.usecase.ListActiveCourses(ctx, page, perPage)
	} else {
		result, err = h.usecase.ListCourses(ctx, page, perPage)
	}

	if err != nil {
		response.SafeInternalError(c, "Failed to fetch courses", err)
		return
	}

	response.Success(c, result)
}

// GetCourseByID handles GET /api/v1/courses/:id
func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	courseData, err := h.usecase.GetCourseByID(ctx, id)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch course", err)
		return
	}

	if courseData == nil {
		response.NotFound(c, "Course not found")
		return
	}

	response.Success(c, courseData)
}

// CreateCourse handles POST /api/v1/courses
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	courseData, err := h.usecase.CreateCourse(ctx, &req)
	if err != nil {
		response.SafeInternalError(c, "Failed to create course", err)
		return
	}

	response.Created(c, courseData)
}

// UpdateCourse handles PUT /api/v1/courses/:id
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req entity.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	courseData, err := h.usecase.UpdateCourse(ctx, id, &req)
	if err != nil {
		response.SafeInternalError(c, "Failed to update course", err)
		return
	}

	if courseData == nil {
		response.NotFound(c, "Course not found")
		return
	}

	response.Success(c, courseData)
}

// DeleteCourse handles DELETE /api/v1/courses/:id
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	// Check if course exists
	courseData, err := h.usecase.GetCourseByID(ctx, id)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch course", err)
		return
	}

	if courseData == nil {
		response.NotFound(c, "Course not found")
		return
	}

	if err := h.usecase.DeleteCourse(ctx, id); err != nil {
		response.SafeInternalError(c, "Failed to delete course", err)
		return
	}

	response.Success(c, gin.H{
		"message": "Course deleted successfully",
	})
}
