package entity

import "time"

// Matricula represents an enrollment entity
type Matricula struct {
	ID               string     `db:"id" json:"id"`
	StudentID        string     `db:"student_id" json:"student_id"`
	StudentName      string     `db:"student_name" json:"student_name"`
	StudentEmail     string     `db:"student_email" json:"student_email"`
	StudentCPF       *string    `db:"student_cpf" json:"student_cpf,omitempty"`
	StudentPhone     *string    `db:"student_phone" json:"student_phone,omitempty"`
	CourseID         string     `db:"course_id" json:"course_id"`
	CourseName       string     `db:"course_name" json:"course_name"`
	InstructorID     *string    `db:"instructor_id" json:"instructor_id,omitempty"`
	InstructorName   *string    `db:"instructor_name" json:"instructor_name,omitempty"`
	PaymentID        *string    `db:"payment_id" json:"payment_id,omitempty"`
	PaymentStatus    string     `db:"payment_status" json:"payment_status"`
	Amount           float64    `db:"amount" json:"amount"`
	DiscountAmount   float64    `db:"discount_amount" json:"discount_amount"`
	FinalAmount      float64    `db:"final_amount" json:"final_amount"`
	PaymentMethod    *string    `db:"payment_method" json:"payment_method,omitempty"`
	EnrollmentDate   time.Time  `db:"enrollment_date" json:"enrollment_date"`
	CompletionDate   *time.Time `db:"completion_date" json:"completion_date,omitempty"`
	ExpirationDate   *time.Time `db:"expiration_date" json:"expiration_date,omitempty"`
	Status           string     `db:"status" json:"status"`
	Progress         float64    `db:"progress" json:"progress"`
	CertificateID    *string    `db:"certificate_id" json:"certificate_id,omitempty"`
	AsaasCustomerID  *string    `db:"asaas_customer_id" json:"asaas_customer_id,omitempty"`
	AsaasPaymentID   *string    `db:"asaas_payment_id" json:"asaas_payment_id,omitempty"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// Enrollment status constants
const (
	EnrollmentStatusPending   = "pending"
	EnrollmentStatusActive    = "active"
	EnrollmentStatusCompleted = "completed"
	EnrollmentStatusCancelled = "cancelled"
	EnrollmentStatusExpired   = "expired"
)

// Payment status constants
const (
	PaymentStatusPending    = "pending"
	PaymentStatusConfirmed  = "confirmed"
	PaymentStatusFailed     = "failed"
	PaymentStatusRefunded   = "refunded"
	PaymentStatusOverdue    = "overdue"
	PaymentStatusChargeback = "chargeback"
)

// CreateMatriculaRequest represents the request to create an enrollment
type CreateMatriculaRequest struct {
	StudentID      string   `json:"student_id" binding:"required"`
	StudentName    string   `json:"student_name" binding:"required"`
	StudentEmail   string   `json:"student_email" binding:"required,email"`
	StudentCPF     *string  `json:"student_cpf,omitempty"`
	StudentPhone   *string  `json:"student_phone,omitempty"`
	CourseID       string   `json:"course_id" binding:"required"`
	CourseName     string   `json:"course_name" binding:"required"`
	InstructorID   *string  `json:"instructor_id,omitempty"`
	InstructorName *string  `json:"instructor_name,omitempty"`
	Amount         float64  `json:"amount" binding:"required,gt=0"`
	DiscountAmount float64  `json:"discount_amount,omitempty"`
	PaymentMethod  string   `json:"payment_method" binding:"required"`
}

// UpdateMatriculaRequest represents the request to update an enrollment
type UpdateMatriculaRequest struct {
	Status         *string  `json:"status,omitempty"`
	PaymentStatus  *string  `json:"payment_status,omitempty"`
	Progress       *float64 `json:"progress,omitempty"`
	CertificateID  *string  `json:"certificate_id,omitempty"`
	PaymentID      *string  `json:"payment_id,omitempty"`
	AsaasPaymentID *string  `json:"asaas_payment_id,omitempty"`
}

// MatriculaListResponse represents the response for listing enrollments
type MatriculaListResponse struct {
	Enrollments []Matricula `json:"enrollments"`
	Total       int         `json:"total"`
	Page        int         `json:"page"`
	PerPage     int         `json:"per_page"`
}
