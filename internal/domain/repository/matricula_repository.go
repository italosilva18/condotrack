package repository

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

// MatriculaRepository defines the interface for matricula data access
type MatriculaRepository interface {
	// FindAll returns all matriculas with pagination
	FindAll(ctx context.Context, page, perPage int) ([]entity.Matricula, int, error)

	// FindByID returns a matricula by ID
	FindByID(ctx context.Context, id string) (*entity.Matricula, error)

	// FindByStudentID returns all matriculas for a specific student
	FindByStudentID(ctx context.Context, studentID string) ([]entity.Matricula, error)

	// FindByCourseID returns all matriculas for a specific course
	FindByCourseID(ctx context.Context, courseID string) ([]entity.Matricula, error)

	// FindByAsaasPaymentID returns a matricula by Asaas payment ID
	FindByAsaasPaymentID(ctx context.Context, asaasPaymentID string) (*entity.Matricula, error)

	// Create creates a new matricula
	Create(ctx context.Context, matricula *entity.Matricula) error

	// CreateWithTx creates a new matricula within a transaction
	CreateWithTx(ctx context.Context, tx *sqlx.Tx, matricula *entity.Matricula) error

	// Update updates an existing matricula
	Update(ctx context.Context, matricula *entity.Matricula) error

	// UpdateWithTx updates an existing matricula within a transaction
	UpdateWithTx(ctx context.Context, tx *sqlx.Tx, matricula *entity.Matricula) error

	// UpdatePaymentStatus updates only the payment status
	UpdatePaymentStatus(ctx context.Context, id, paymentStatus string) error

	// UpdatePaymentStatusWithTx updates payment status within a transaction
	UpdatePaymentStatusWithTx(ctx context.Context, tx *sqlx.Tx, id, paymentStatus string) error

	// UpdateStatus updates the enrollment status
	UpdateStatus(ctx context.Context, id, status string) error

	// UpdateProgress updates the progress percentage
	UpdateProgress(ctx context.Context, id string, progress float64) error

	// Delete deletes a matricula by ID
	Delete(ctx context.Context, id string) error

	// CountByStudentID returns the number of enrollments for a student
	CountByStudentID(ctx context.Context, studentID string) (int, error)

	// CountByCourseID returns the number of enrollments for a course
	CountByCourseID(ctx context.Context, courseID string) (int, error)
}
