package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type matriculaMySQLRepository struct {
	db *sqlx.DB
}

// NewMatriculaMySQLRepository creates a new MySQL implementation of MatriculaRepository
func NewMatriculaMySQLRepository(db *sqlx.DB) repository.MatriculaRepository {
	return &matriculaMySQLRepository{db: db}
}

func (r *matriculaMySQLRepository) FindAll(ctx context.Context, page, perPage int) ([]entity.Matricula, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM enrollments`
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	var matriculas []entity.Matricula
	query := `SELECT id, student_id, student_name, student_email, student_cpf, student_phone,
			  course_id, course_name, instructor_id, instructor_name, payment_id, payment_status,
			  amount, discount_amount, final_amount, payment_method, enrollment_date, completion_date,
			  expiration_date, status, progress, certificate_id, asaas_customer_id, asaas_payment_id,
			  created_at, updated_at
			  FROM enrollments
			  ORDER BY created_at DESC
			  LIMIT ? OFFSET ?`
	err := r.db.SelectContext(ctx, &matriculas, query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	return matriculas, total, nil
}

func (r *matriculaMySQLRepository) FindByID(ctx context.Context, id string) (*entity.Matricula, error) {
	var matricula entity.Matricula
	query := `SELECT id, student_id, student_name, student_email, student_cpf, student_phone,
			  course_id, course_name, instructor_id, instructor_name, payment_id, payment_status,
			  amount, discount_amount, final_amount, payment_method, enrollment_date, completion_date,
			  expiration_date, status, progress, certificate_id, asaas_customer_id, asaas_payment_id,
			  created_at, updated_at
			  FROM enrollments
			  WHERE id = ?`
	err := r.db.GetContext(ctx, &matricula, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &matricula, nil
}

func (r *matriculaMySQLRepository) FindByStudentID(ctx context.Context, studentID string) ([]entity.Matricula, error) {
	var matriculas []entity.Matricula
	query := `SELECT id, student_id, student_name, student_email, student_cpf, student_phone,
			  course_id, course_name, instructor_id, instructor_name, payment_id, payment_status,
			  amount, discount_amount, final_amount, payment_method, enrollment_date, completion_date,
			  expiration_date, status, progress, certificate_id, asaas_customer_id, asaas_payment_id,
			  created_at, updated_at
			  FROM enrollments
			  WHERE student_id = ?
			  ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &matriculas, query, studentID)
	return matriculas, err
}

func (r *matriculaMySQLRepository) FindByCourseID(ctx context.Context, courseID string) ([]entity.Matricula, error) {
	var matriculas []entity.Matricula
	query := `SELECT id, student_id, student_name, student_email, student_cpf, student_phone,
			  course_id, course_name, instructor_id, instructor_name, payment_id, payment_status,
			  amount, discount_amount, final_amount, payment_method, enrollment_date, completion_date,
			  expiration_date, status, progress, certificate_id, asaas_customer_id, asaas_payment_id,
			  created_at, updated_at
			  FROM enrollments
			  WHERE course_id = ?
			  ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &matriculas, query, courseID)
	return matriculas, err
}

func (r *matriculaMySQLRepository) FindByAsaasPaymentID(ctx context.Context, asaasPaymentID string) (*entity.Matricula, error) {
	var matricula entity.Matricula
	query := `SELECT id, student_id, student_name, student_email, student_cpf, student_phone,
			  course_id, course_name, instructor_id, instructor_name, payment_id, payment_status,
			  amount, discount_amount, final_amount, payment_method, enrollment_date, completion_date,
			  expiration_date, status, progress, certificate_id, asaas_customer_id, asaas_payment_id,
			  created_at, updated_at
			  FROM enrollments
			  WHERE asaas_payment_id = ?`
	err := r.db.GetContext(ctx, &matricula, query, asaasPaymentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &matricula, nil
}

func (r *matriculaMySQLRepository) Create(ctx context.Context, m *entity.Matricula) error {
	query := `INSERT INTO enrollments (id, student_id, student_name, student_email, student_cpf, student_phone,
			  course_id, course_name, instructor_id, instructor_name, payment_id, payment_status,
			  amount, discount_amount, final_amount, payment_method, enrollment_date, completion_date,
			  expiration_date, status, progress, certificate_id, asaas_customer_id, asaas_payment_id, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query,
		m.ID, m.StudentID, m.StudentName, m.StudentEmail, m.StudentCPF, m.StudentPhone,
		m.CourseID, m.CourseName, m.InstructorID, m.InstructorName, m.PaymentID, m.PaymentStatus,
		m.Amount, m.DiscountAmount, m.FinalAmount, m.PaymentMethod, m.EnrollmentDate, m.CompletionDate,
		m.ExpirationDate, m.Status, m.Progress, m.CertificateID, m.AsaasCustomerID, m.AsaasPaymentID)
	return err
}

func (r *matriculaMySQLRepository) CreateWithTx(ctx context.Context, tx *sqlx.Tx, m *entity.Matricula) error {
	query := `INSERT INTO enrollments (id, student_id, student_name, student_email, student_cpf, student_phone,
			  course_id, course_name, instructor_id, instructor_name, payment_id, payment_status,
			  amount, discount_amount, final_amount, payment_method, enrollment_date, completion_date,
			  expiration_date, status, progress, certificate_id, asaas_customer_id, asaas_payment_id, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := tx.ExecContext(ctx, query,
		m.ID, m.StudentID, m.StudentName, m.StudentEmail, m.StudentCPF, m.StudentPhone,
		m.CourseID, m.CourseName, m.InstructorID, m.InstructorName, m.PaymentID, m.PaymentStatus,
		m.Amount, m.DiscountAmount, m.FinalAmount, m.PaymentMethod, m.EnrollmentDate, m.CompletionDate,
		m.ExpirationDate, m.Status, m.Progress, m.CertificateID, m.AsaasCustomerID, m.AsaasPaymentID)
	return err
}

func (r *matriculaMySQLRepository) Update(ctx context.Context, m *entity.Matricula) error {
	query := `UPDATE enrollments SET
			  student_name = ?, student_email = ?, student_cpf = ?, student_phone = ?,
			  course_name = ?, instructor_id = ?, instructor_name = ?, payment_id = ?, payment_status = ?,
			  amount = ?, discount_amount = ?, final_amount = ?, payment_method = ?, completion_date = ?,
			  expiration_date = ?, status = ?, progress = ?, certificate_id = ?,
			  asaas_customer_id = ?, asaas_payment_id = ?, updated_at = NOW()
			  WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		m.StudentName, m.StudentEmail, m.StudentCPF, m.StudentPhone,
		m.CourseName, m.InstructorID, m.InstructorName, m.PaymentID, m.PaymentStatus,
		m.Amount, m.DiscountAmount, m.FinalAmount, m.PaymentMethod, m.CompletionDate,
		m.ExpirationDate, m.Status, m.Progress, m.CertificateID,
		m.AsaasCustomerID, m.AsaasPaymentID, m.ID)
	return err
}

func (r *matriculaMySQLRepository) UpdateWithTx(ctx context.Context, tx *sqlx.Tx, m *entity.Matricula) error {
	query := `UPDATE enrollments SET
			  student_name = ?, student_email = ?, student_cpf = ?, student_phone = ?,
			  course_name = ?, instructor_id = ?, instructor_name = ?, payment_id = ?, payment_status = ?,
			  amount = ?, discount_amount = ?, final_amount = ?, payment_method = ?, completion_date = ?,
			  expiration_date = ?, status = ?, progress = ?, certificate_id = ?,
			  asaas_customer_id = ?, asaas_payment_id = ?, updated_at = NOW()
			  WHERE id = ?`
	_, err := tx.ExecContext(ctx, query,
		m.StudentName, m.StudentEmail, m.StudentCPF, m.StudentPhone,
		m.CourseName, m.InstructorID, m.InstructorName, m.PaymentID, m.PaymentStatus,
		m.Amount, m.DiscountAmount, m.FinalAmount, m.PaymentMethod, m.CompletionDate,
		m.ExpirationDate, m.Status, m.Progress, m.CertificateID,
		m.AsaasCustomerID, m.AsaasPaymentID, m.ID)
	return err
}

func (r *matriculaMySQLRepository) UpdatePaymentStatus(ctx context.Context, id, paymentStatus string) error {
	query := `UPDATE enrollments SET payment_status = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, paymentStatus, id)
	return err
}

func (r *matriculaMySQLRepository) UpdatePaymentStatusWithTx(ctx context.Context, tx *sqlx.Tx, id, paymentStatus string) error {
	query := `UPDATE enrollments SET payment_status = ?, updated_at = NOW() WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, paymentStatus, id)
	return err
}

func (r *matriculaMySQLRepository) UpdateStatus(ctx context.Context, id, status string) error {
	query := `UPDATE enrollments SET status = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *matriculaMySQLRepository) UpdateProgress(ctx context.Context, id string, progress float64) error {
	query := `UPDATE enrollments SET progress = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, progress, id)
	return err
}

func (r *matriculaMySQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM enrollments WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *matriculaMySQLRepository) CountByStudentID(ctx context.Context, studentID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM enrollments WHERE student_id = ?`
	err := r.db.GetContext(ctx, &count, query, studentID)
	return count, err
}

func (r *matriculaMySQLRepository) CountByCourseID(ctx context.Context, courseID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM enrollments WHERE course_id = ?`
	err := r.db.GetContext(ctx, &count, query, courseID)
	return count, err
}
