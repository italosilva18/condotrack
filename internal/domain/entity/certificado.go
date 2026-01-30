package entity

import "time"

// Certificado represents a certificate entity
type Certificado struct {
	ID             string     `db:"id" json:"id"`
	EnrollmentID   string     `db:"enrollment_id" json:"enrollment_id"`
	StudentID      string     `db:"student_id" json:"student_id"`
	StudentName    string     `db:"student_name" json:"student_name"`
	StudentCPF     *string    `db:"student_cpf" json:"student_cpf,omitempty"`
	CourseID       string     `db:"course_id" json:"course_id"`
	CourseName     string     `db:"course_name" json:"course_name"`
	CourseHours    int        `db:"course_hours" json:"course_hours"`
	InstructorName *string    `db:"instructor_name" json:"instructor_name,omitempty"`
	CompletionDate time.Time  `db:"completion_date" json:"completion_date"`
	IssueDate      time.Time  `db:"issue_date" json:"issue_date"`
	ValidationCode string     `db:"validation_code" json:"validation_code"`
	Status         string     `db:"status" json:"status"`
	DownloadURL    *string    `db:"download_url" json:"download_url,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
}

// Certificate status constants
const (
	CertificateStatusActive   = "active"
	CertificateStatusRevoked  = "revoked"
	CertificateStatusExpired  = "expired"
)

// CertificateData represents the data needed to generate a certificate
type CertificateData struct {
	StudentName    string    `json:"student_name"`
	StudentCPF     string    `json:"student_cpf,omitempty"`
	CourseName     string    `json:"course_name"`
	CourseHours    int       `json:"course_hours"`
	InstructorName string    `json:"instructor_name,omitempty"`
	CompletionDate time.Time `json:"completion_date"`
	ValidationCode string    `json:"validation_code"`
}

// GenerateCertificateRequest represents the request to generate a certificate
type GenerateCertificateRequest struct {
	EnrollmentID string `json:"enrollment_id" binding:"required"`
}

// ValidateCertificateRequest represents the request to validate a certificate
type ValidateCertificateRequest struct {
	ValidationCode string `json:"validation_code" binding:"required"`
}

// CertificateResponse represents the certificate response
type CertificateResponse struct {
	ID             string `json:"id"`
	StudentName    string `json:"student_name"`
	CourseName     string `json:"course_name"`
	CourseHours    int    `json:"course_hours"`
	CompletionDate string `json:"completion_date"`
	IssueDate      string `json:"issue_date"`
	ValidationCode string `json:"validation_code"`
	ValidationURL  string `json:"validation_url"`
	DownloadURL    string `json:"download_url,omitempty"`
	Valid          bool   `json:"valid"`
}
