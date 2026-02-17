package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type certificadoMySQLRepository struct {
	db *sqlx.DB
}

// NewCertificadoMySQLRepository creates a new MySQL implementation of CertificadoRepository
func NewCertificadoMySQLRepository(db *sqlx.DB) repository.CertificadoRepository {
	return &certificadoMySQLRepository{db: db}
}

func (r *certificadoMySQLRepository) FindByID(ctx context.Context, id string) (*entity.Certificado, error) {
	var cert entity.Certificado
	query := `SELECT id, enrollment_id, student_id, student_name, student_cpf, course_id, course_name,
			  course_hours, instructor_name, completion_date, issue_date, validation_code, status,
			  download_url, created_at
			  FROM certificates
			  WHERE id = ?`
	err := r.db.GetContext(ctx, &cert, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &cert, nil
}

func (r *certificadoMySQLRepository) FindByValidationCode(ctx context.Context, code string) (*entity.Certificado, error) {
	var cert entity.Certificado
	query := `SELECT id, enrollment_id, student_id, student_name, student_cpf, course_id, course_name,
			  course_hours, instructor_name, completion_date, issue_date, validation_code, status,
			  download_url, created_at
			  FROM certificates
			  WHERE validation_code = ?`
	err := r.db.GetContext(ctx, &cert, query, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &cert, nil
}

func (r *certificadoMySQLRepository) FindByStudentID(ctx context.Context, studentID string) ([]entity.Certificado, error) {
	var certs []entity.Certificado
	query := `SELECT id, enrollment_id, student_id, student_name, student_cpf, course_id, course_name,
			  course_hours, instructor_name, completion_date, issue_date, validation_code, status,
			  download_url, created_at
			  FROM certificates
			  WHERE student_id = ?
			  ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &certs, query, studentID)
	return certs, err
}

func (r *certificadoMySQLRepository) FindByEnrollmentID(ctx context.Context, enrollmentID string) (*entity.Certificado, error) {
	var cert entity.Certificado
	query := `SELECT id, enrollment_id, student_id, student_name, student_cpf, course_id, course_name,
			  course_hours, instructor_name, completion_date, issue_date, validation_code, status,
			  download_url, created_at
			  FROM certificates
			  WHERE enrollment_id = ?`
	err := r.db.GetContext(ctx, &cert, query, enrollmentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &cert, nil
}

func (r *certificadoMySQLRepository) Create(ctx context.Context, cert *entity.Certificado) error {
	query := `INSERT INTO certificates (id, enrollment_id, student_id, student_name, student_cpf,
			  course_id, course_name, course_hours, instructor_name, completion_date, issue_date,
			  validation_code, status, download_url, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query,
		cert.ID, cert.EnrollmentID, cert.StudentID, cert.StudentName, cert.StudentCPF,
		cert.CourseID, cert.CourseName, cert.CourseHours, cert.InstructorName,
		cert.CompletionDate, cert.IssueDate, cert.ValidationCode, cert.Status, cert.DownloadURL)
	return err
}

func (r *certificadoMySQLRepository) Update(ctx context.Context, cert *entity.Certificado) error {
	query := `UPDATE certificates SET
			  student_name = ?, student_cpf = ?, course_name = ?, course_hours = ?,
			  instructor_name = ?, completion_date = ?, issue_date = ?, status = ?, download_url = ?
			  WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		cert.StudentName, cert.StudentCPF, cert.CourseName, cert.CourseHours,
		cert.InstructorName, cert.CompletionDate, cert.IssueDate, cert.Status, cert.DownloadURL, cert.ID)
	return err
}

func (r *certificadoMySQLRepository) UpdateStatus(ctx context.Context, id, status string) error {
	query := `UPDATE certificates SET status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *certificadoMySQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM certificates WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *certificadoMySQLRepository) Exists(ctx context.Context, enrollmentID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM certificates WHERE enrollment_id = ?`
	err := r.db.GetContext(ctx, &count, query, enrollmentID)
	return count > 0, err
}

// NotificacaoMySQLRepository implementation
type notificacaoMySQLRepository struct {
	db *sqlx.DB
}

// NewNotificacaoMySQLRepository creates a new MySQL implementation of NotificacaoRepository
func NewNotificacaoMySQLRepository(db *sqlx.DB) repository.NotificacaoRepository {
	return &notificacaoMySQLRepository{db: db}
}

func (r *notificacaoMySQLRepository) FindByUserID(ctx context.Context, userID string) ([]entity.Notificacao, error) {
	var notifs []entity.Notificacao
	query := `SELECT id, user_id, type, title, message, data, is_read, read_at, created_at
			  FROM notifications
			  WHERE user_id = ?
			  ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &notifs, query, userID)
	return notifs, err
}

func (r *notificacaoMySQLRepository) FindUnreadByUserID(ctx context.Context, userID string) ([]entity.Notificacao, error) {
	var notifs []entity.Notificacao
	query := `SELECT id, user_id, type, title, message, data, is_read, read_at, created_at
			  FROM notifications
			  WHERE user_id = ? AND is_read = 0
			  ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &notifs, query, userID)
	return notifs, err
}

func (r *notificacaoMySQLRepository) Create(ctx context.Context, notif *entity.Notificacao) error {
	query := `INSERT INTO notifications (id, user_id, type, title, message, data, is_read, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, 0, NOW())`
	_, err := r.db.ExecContext(ctx, query, notif.ID, notif.UserID, notif.Type, notif.Title, notif.Message, notif.Data)
	return err
}

func (r *notificacaoMySQLRepository) MarkAsRead(ctx context.Context, id string) error {
	query := `UPDATE notifications SET is_read = 1, read_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *notificacaoMySQLRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	query := `UPDATE notifications SET is_read = 1, read_at = NOW() WHERE user_id = ? AND is_read = 0`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *notificacaoMySQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM notifications WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *notificacaoMySQLRepository) CountUnread(ctx context.Context, userID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = 0`
	err := r.db.GetContext(ctx, &count, query, userID)
	return count, err
}

// RevenueSplitMySQLRepository implementation
type revenueSplitMySQLRepository struct {
	db *sqlx.DB
}

// NewRevenueSplitMySQLRepository creates a new MySQL implementation of RevenueSplitRepository
func NewRevenueSplitMySQLRepository(db *sqlx.DB) repository.RevenueSplitRepository {
	return &revenueSplitMySQLRepository{db: db}
}

func (r *revenueSplitMySQLRepository) FindByID(ctx context.Context, id string) (*entity.RevenueSplit, error) {
	var split entity.RevenueSplit
	query := `SELECT id, enrollment_id, payment_id, gross_amount, net_amount, platform_fee,
			  payment_fee, instructor_amount, platform_amount, instructor_id, payment_method,
			  status, processed_at, created_at
			  FROM revenue_splits
			  WHERE id = ?`
	err := r.db.GetContext(ctx, &split, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &split, nil
}

func (r *revenueSplitMySQLRepository) FindAll(ctx context.Context, status string) ([]entity.RevenueSplit, error) {
	var splits []entity.RevenueSplit
	var err error

	if status != "" {
		query := `SELECT id, enrollment_id, payment_id, gross_amount, net_amount, platform_fee,
				  payment_fee, instructor_amount, platform_amount, instructor_id, payment_method,
				  status, processed_at, created_at
				  FROM revenue_splits
				  WHERE status = ?
				  ORDER BY created_at DESC`
		err = r.db.SelectContext(ctx, &splits, query, status)
	} else {
		query := `SELECT id, enrollment_id, payment_id, gross_amount, net_amount, platform_fee,
				  payment_fee, instructor_amount, platform_amount, instructor_id, payment_method,
				  status, processed_at, created_at
				  FROM revenue_splits
				  ORDER BY created_at DESC`
		err = r.db.SelectContext(ctx, &splits, query)
	}

	if err != nil {
		return nil, err
	}
	return splits, nil
}

func (r *revenueSplitMySQLRepository) FindByEnrollmentID(ctx context.Context, enrollmentID string) (*entity.RevenueSplit, error) {
	var split entity.RevenueSplit
	query := `SELECT id, enrollment_id, payment_id, gross_amount, net_amount, platform_fee,
			  payment_fee, instructor_amount, platform_amount, instructor_id, payment_method,
			  status, processed_at, created_at
			  FROM revenue_splits
			  WHERE enrollment_id = ?`
	err := r.db.GetContext(ctx, &split, query, enrollmentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &split, nil
}

func (r *revenueSplitMySQLRepository) FindByPaymentID(ctx context.Context, paymentID string) (*entity.RevenueSplit, error) {
	var split entity.RevenueSplit
	query := `SELECT id, enrollment_id, payment_id, gross_amount, net_amount, platform_fee,
			  payment_fee, instructor_amount, platform_amount, instructor_id, payment_method,
			  status, processed_at, created_at
			  FROM revenue_splits
			  WHERE payment_id = ?`
	err := r.db.GetContext(ctx, &split, query, paymentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &split, nil
}

func (r *revenueSplitMySQLRepository) FindByInstructorID(ctx context.Context, instructorID string) ([]entity.RevenueSplit, error) {
	var splits []entity.RevenueSplit
	query := `SELECT id, enrollment_id, payment_id, gross_amount, net_amount, platform_fee,
			  payment_fee, instructor_amount, platform_amount, instructor_id, payment_method,
			  status, processed_at, created_at
			  FROM revenue_splits
			  WHERE instructor_id = ?
			  ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &splits, query, instructorID)
	return splits, err
}

func (r *revenueSplitMySQLRepository) Create(ctx context.Context, split *entity.RevenueSplit) error {
	query := `INSERT INTO revenue_splits (id, enrollment_id, payment_id, gross_amount, net_amount,
			  platform_fee, payment_fee, instructor_amount, platform_amount, instructor_id,
			  payment_method, status, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query,
		split.ID, split.EnrollmentID, split.PaymentID, split.GrossAmount, split.NetAmount,
		split.PlatformFee, split.PaymentFee, split.InstructorAmount, split.PlatformAmount,
		split.InstructorID, split.PaymentMethod, split.Status)
	return err
}

func (r *revenueSplitMySQLRepository) CreateWithTx(ctx context.Context, tx *sqlx.Tx, split *entity.RevenueSplit) error {
	query := `INSERT INTO revenue_splits (id, enrollment_id, payment_id, gross_amount, net_amount,
			  platform_fee, payment_fee, instructor_amount, platform_amount, instructor_id,
			  payment_method, status, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := tx.ExecContext(ctx, query,
		split.ID, split.EnrollmentID, split.PaymentID, split.GrossAmount, split.NetAmount,
		split.PlatformFee, split.PaymentFee, split.InstructorAmount, split.PlatformAmount,
		split.InstructorID, split.PaymentMethod, split.Status)
	return err
}

func (r *revenueSplitMySQLRepository) UpdateStatus(ctx context.Context, id, status string) error {
	query := `UPDATE revenue_splits SET status = ?, processed_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *revenueSplitMySQLRepository) GetTotalByInstructor(ctx context.Context, instructorID string) (float64, error) {
	var total float64
	query := `SELECT COALESCE(SUM(instructor_amount), 0) FROM revenue_splits
			  WHERE instructor_id = ? AND status = 'processed'`
	err := r.db.GetContext(ctx, &total, query, instructorID)
	return total, err
}
