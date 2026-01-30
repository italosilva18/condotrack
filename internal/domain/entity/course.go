package entity

import "time"

// Course represents a course entity
type Course struct {
	ID             string     `db:"id" json:"id"`
	Name           string     `db:"name" json:"name"`
	Description    *string    `db:"description" json:"description,omitempty"`
	InstructorID   *string    `db:"instructor_id" json:"instructor_id,omitempty"`
	InstructorName *string    `db:"instructor_name" json:"instructor_name,omitempty"`
	DurationHours  int        `db:"duration_hours" json:"duration_hours"`
	Price          float64    `db:"price" json:"price"`
	DiscountPrice  *float64   `db:"discount_price" json:"discount_price,omitempty"`
	ThumbnailURL   *string    `db:"thumbnail_url" json:"thumbnail_url,omitempty"`
	IsActive       bool       `db:"is_active" json:"is_active"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// CreateCourseRequest represents the request to create a course
type CreateCourseRequest struct {
	Name          string   `json:"name" binding:"required"`
	Description   *string  `json:"description,omitempty"`
	InstructorID  *string  `json:"instructor_id,omitempty"`
	DurationHours int      `json:"duration_hours" binding:"required,gt=0"`
	Price         float64  `json:"price" binding:"required,gte=0"`
	DiscountPrice *float64 `json:"discount_price,omitempty"`
	ThumbnailURL  *string  `json:"thumbnail_url,omitempty"`
	IsActive      *bool    `json:"is_active,omitempty"`
}

// UpdateCourseRequest represents the request to update a course
type UpdateCourseRequest struct {
	Name          *string  `json:"name,omitempty"`
	Description   *string  `json:"description,omitempty"`
	InstructorID  *string  `json:"instructor_id,omitempty"`
	DurationHours *int     `json:"duration_hours,omitempty"`
	Price         *float64 `json:"price,omitempty"`
	DiscountPrice *float64 `json:"discount_price,omitempty"`
	ThumbnailURL  *string  `json:"thumbnail_url,omitempty"`
	IsActive      *bool    `json:"is_active,omitempty"`
}

// CourseListResponse represents the response for listing courses
type CourseListResponse struct {
	Courses []Course `json:"courses"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
}
