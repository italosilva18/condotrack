package course

import (
	"context"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/google/uuid"
)

// UseCase defines the course use case interface
type UseCase interface {
	ListCourses(ctx context.Context, page, perPage int) (*entity.CourseListResponse, error)
	ListCoursesByInstructor(ctx context.Context, instructorID string) ([]entity.Course, error)
	ListActiveCourses(ctx context.Context, page, perPage int) (*entity.CourseListResponse, error)
	GetCourseByID(ctx context.Context, id string) (*entity.Course, error)
	CreateCourse(ctx context.Context, req *entity.CreateCourseRequest) (*entity.Course, error)
	UpdateCourse(ctx context.Context, id string, req *entity.UpdateCourseRequest) (*entity.Course, error)
	DeleteCourse(ctx context.Context, id string) error
}

type courseUseCase struct {
	repo repository.CourseRepository
}

// NewUseCase creates a new course use case
func NewUseCase(repo repository.CourseRepository) UseCase {
	return &courseUseCase{repo: repo}
}

// ListCourses returns all courses with pagination
func (uc *courseUseCase) ListCourses(ctx context.Context, page, perPage int) (*entity.CourseListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	courses, total, err := uc.repo.FindAll(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	return &entity.CourseListResponse{
		Courses: courses,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}, nil
}

// ListCoursesByInstructor returns all courses for a specific instructor
func (uc *courseUseCase) ListCoursesByInstructor(ctx context.Context, instructorID string) ([]entity.Course, error) {
	return uc.repo.FindByInstructor(ctx, instructorID)
}

// ListActiveCourses returns all active courses with pagination
func (uc *courseUseCase) ListActiveCourses(ctx context.Context, page, perPage int) (*entity.CourseListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	courses, total, err := uc.repo.FindActive(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	return &entity.CourseListResponse{
		Courses: courses,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}, nil
}

// GetCourseByID returns a specific course by ID
func (uc *courseUseCase) GetCourseByID(ctx context.Context, id string) (*entity.Course, error) {
	return uc.repo.FindByID(ctx, id)
}

// CreateCourse creates a new course
func (uc *courseUseCase) CreateCourse(ctx context.Context, req *entity.CreateCourseRequest) (*entity.Course, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	course := &entity.Course{
		ID:            uuid.New().String(),
		Name:          req.Name,
		Description:   req.Description,
		InstructorID:  req.InstructorID,
		DurationHours: req.DurationHours,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		ThumbnailURL:  req.ThumbnailURL,
		IsActive:      isActive,
		CreatedAt:     time.Now(),
	}

	if err := uc.repo.Create(ctx, course); err != nil {
		return nil, err
	}

	// Fetch the complete course with instructor name
	return uc.repo.FindByID(ctx, course.ID)
}

// UpdateCourse updates an existing course
func (uc *courseUseCase) UpdateCourse(ctx context.Context, id string, req *entity.UpdateCourseRequest) (*entity.Course, error) {
	// Get existing course
	course, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, nil
	}

	// Update fields if provided
	if req.Name != nil {
		course.Name = *req.Name
	}
	if req.Description != nil {
		course.Description = req.Description
	}
	if req.InstructorID != nil {
		course.InstructorID = req.InstructorID
	}
	if req.DurationHours != nil {
		course.DurationHours = *req.DurationHours
	}
	if req.Price != nil {
		course.Price = *req.Price
	}
	if req.DiscountPrice != nil {
		course.DiscountPrice = req.DiscountPrice
	}
	if req.ThumbnailURL != nil {
		course.ThumbnailURL = req.ThumbnailURL
	}
	if req.IsActive != nil {
		course.IsActive = *req.IsActive
	}

	now := time.Now()
	course.UpdatedAt = &now

	if err := uc.repo.Update(ctx, course); err != nil {
		return nil, err
	}

	// Fetch the complete course with instructor name
	return uc.repo.FindByID(ctx, course.ID)
}

// DeleteCourse deletes a course by ID
func (uc *courseUseCase) DeleteCourse(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
