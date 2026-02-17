package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type courseMySQLRepository struct {
	db *sqlx.DB
}

// NewCourseMySQLRepository creates a new MySQL implementation of CourseRepository
func NewCourseMySQLRepository(db *sqlx.DB) repository.CourseRepository {
	return &courseMySQLRepository{db: db}
}

func (r *courseMySQLRepository) FindAll(ctx context.Context, page, perPage int) ([]entity.Course, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM courses`
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	var courses []entity.Course
	query := `SELECT c.id, c.name, c.description, c.instructor_id,
			  u.name as instructor_name, c.duration_hours, c.price,
			  c.discount_price, c.thumbnail_url, c.is_active,
			  c.created_at, c.updated_at
			  FROM courses c
			  LEFT JOIN users u ON c.instructor_id = u.id
			  ORDER BY c.created_at DESC
			  LIMIT ? OFFSET ?`
	err := r.db.SelectContext(ctx, &courses, query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	return courses, total, nil
}

func (r *courseMySQLRepository) FindByID(ctx context.Context, id string) (*entity.Course, error) {
	var course entity.Course
	query := `SELECT c.id, c.name, c.description, c.instructor_id,
			  u.name as instructor_name, c.duration_hours, c.price,
			  c.discount_price, c.thumbnail_url, c.is_active,
			  c.created_at, c.updated_at
			  FROM courses c
			  LEFT JOIN users u ON c.instructor_id = u.id
			  WHERE c.id = ?`
	err := r.db.GetContext(ctx, &course, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &course, nil
}

func (r *courseMySQLRepository) FindByInstructor(ctx context.Context, instructorID string) ([]entity.Course, error) {
	var courses []entity.Course
	query := `SELECT c.id, c.name, c.description, c.instructor_id,
			  u.name as instructor_name, c.duration_hours, c.price,
			  c.discount_price, c.thumbnail_url, c.is_active,
			  c.created_at, c.updated_at
			  FROM courses c
			  LEFT JOIN users u ON c.instructor_id = u.id
			  WHERE c.instructor_id = ?
			  ORDER BY c.created_at DESC`
	err := r.db.SelectContext(ctx, &courses, query, instructorID)
	return courses, err
}

func (r *courseMySQLRepository) FindActive(ctx context.Context, page, perPage int) ([]entity.Course, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM courses WHERE is_active = true`
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	var courses []entity.Course
	query := `SELECT c.id, c.name, c.description, c.instructor_id,
			  u.name as instructor_name, c.duration_hours, c.price,
			  c.discount_price, c.thumbnail_url, c.is_active,
			  c.created_at, c.updated_at
			  FROM courses c
			  LEFT JOIN users u ON c.instructor_id = u.id
			  WHERE c.is_active = true
			  ORDER BY c.created_at DESC
			  LIMIT ? OFFSET ?`
	err := r.db.SelectContext(ctx, &courses, query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	return courses, total, nil
}

func (r *courseMySQLRepository) Create(ctx context.Context, c *entity.Course) error {
	query := `INSERT INTO courses (id, name, description, instructor_id,
			  duration_hours, price, discount_price, thumbnail_url,
			  is_active, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.Name, c.Description, c.InstructorID,
		c.DurationHours, c.Price, c.DiscountPrice, c.ThumbnailURL,
		c.IsActive)
	return err
}

func (r *courseMySQLRepository) Update(ctx context.Context, c *entity.Course) error {
	query := `UPDATE courses SET
			  name = ?, description = ?, instructor_id = ?,
			  duration_hours = ?, price = ?, discount_price = ?,
			  thumbnail_url = ?, is_active = ?, updated_at = NOW()
			  WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		c.Name, c.Description, c.InstructorID,
		c.DurationHours, c.Price, c.DiscountPrice,
		c.ThumbnailURL, c.IsActive, c.ID)
	return err
}

func (r *courseMySQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM courses WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *courseMySQLRepository) CountByInstructor(ctx context.Context, instructorID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM courses WHERE instructor_id = ?`
	err := r.db.GetContext(ctx, &count, query, instructorID)
	return count, err
}
