package matricula

import (
	"context"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/google/uuid"
)

// UseCase defines the matricula use case interface
type UseCase interface {
	ListEnrollments(ctx context.Context, page, perPage int) (*entity.MatriculaListResponse, error)
	ListEnrollmentsByStudent(ctx context.Context, studentID string) ([]entity.Matricula, error)
	GetEnrollmentByID(ctx context.Context, id string) (*entity.Matricula, error)
	CreateEnrollment(ctx context.Context, req *entity.CreateMatriculaRequest) (*entity.Matricula, error)
	UpdatePaymentStatus(ctx context.Context, id, status string) error
	UpdateProgress(ctx context.Context, id string, progress float64) error
}

type matriculaUseCase struct {
	repo repository.MatriculaRepository
}

// NewUseCase creates a new matricula use case
func NewUseCase(repo repository.MatriculaRepository) UseCase {
	return &matriculaUseCase{repo: repo}
}

// ListEnrollments returns all enrollments with pagination
func (uc *matriculaUseCase) ListEnrollments(ctx context.Context, page, perPage int) (*entity.MatriculaListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	enrollments, total, err := uc.repo.FindAll(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	return &entity.MatriculaListResponse{
		Enrollments: enrollments,
		Total:       total,
		Page:        page,
		PerPage:     perPage,
	}, nil
}

// ListEnrollmentsByStudent returns all enrollments for a specific student
func (uc *matriculaUseCase) ListEnrollmentsByStudent(ctx context.Context, studentID string) ([]entity.Matricula, error) {
	return uc.repo.FindByStudentID(ctx, studentID)
}

// GetEnrollmentByID returns a specific enrollment by ID
func (uc *matriculaUseCase) GetEnrollmentByID(ctx context.Context, id string) (*entity.Matricula, error) {
	return uc.repo.FindByID(ctx, id)
}

// CreateEnrollment creates a new enrollment
func (uc *matriculaUseCase) CreateEnrollment(ctx context.Context, req *entity.CreateMatriculaRequest) (*entity.Matricula, error) {
	// Calculate final amount
	finalAmount := req.Amount - req.DiscountAmount
	if finalAmount < 0 {
		finalAmount = 0
	}

	// Create enrollment entity
	enrollment := &entity.Matricula{
		ID:             uuid.New().String(),
		StudentID:      req.StudentID,
		StudentName:    req.StudentName,
		StudentEmail:   req.StudentEmail,
		StudentCPF:     req.StudentCPF,
		StudentPhone:   req.StudentPhone,
		CourseID:       req.CourseID,
		CourseName:     req.CourseName,
		InstructorID:   req.InstructorID,
		InstructorName: req.InstructorName,
		PaymentStatus:  entity.PaymentStatusPending,
		Amount:         req.Amount,
		DiscountAmount: req.DiscountAmount,
		FinalAmount:    finalAmount,
		PaymentMethod:  &req.PaymentMethod,
		EnrollmentDate: time.Now(),
		Status:         entity.EnrollmentStatusPending,
		Progress:       0,
		CreatedAt:      time.Now(),
	}

	if err := uc.repo.Create(ctx, enrollment); err != nil {
		return nil, err
	}

	return enrollment, nil
}

// UpdatePaymentStatus updates the payment status of an enrollment
func (uc *matriculaUseCase) UpdatePaymentStatus(ctx context.Context, id, status string) error {
	return uc.repo.UpdatePaymentStatus(ctx, id, status)
}

// UpdateProgress updates the progress of an enrollment
func (uc *matriculaUseCase) UpdateProgress(ctx context.Context, id string, progress float64) error {
	return uc.repo.UpdateProgress(ctx, id, progress)
}
