package repository

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
)

// CourseRepository defines the interface for course data access
type CourseRepository interface {
	// FindAll returns all courses with pagination
	FindAll(ctx context.Context, page, perPage int) ([]entity.Course, int, error)

	// FindByID returns a course by ID
	FindByID(ctx context.Context, id string) (*entity.Course, error)

	// FindByInstructor returns all courses for a specific instructor
	FindByInstructor(ctx context.Context, instructorID string) ([]entity.Course, error)

	// FindActive returns all active courses with pagination
	FindActive(ctx context.Context, page, perPage int) ([]entity.Course, int, error)

	// Create creates a new course
	Create(ctx context.Context, course *entity.Course) error

	// Update updates an existing course
	Update(ctx context.Context, course *entity.Course) error

	// Delete deletes a course by ID
	Delete(ctx context.Context, id string) error

	// CountByInstructor returns the number of courses for an instructor
	CountByInstructor(ctx context.Context, instructorID string) (int, error)
}
