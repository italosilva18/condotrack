package entity

import "time"

// PaymentRecord represents a payment record derived from enrollments
type PaymentRecord struct {
	ID              string     `json:"id"`
	EnrollmentID    string     `json:"enrollment_id"`
	CustomerID      *string    `json:"customer_id,omitempty"`
	AsaasPaymentID  *string    `json:"asaas_payment_id,omitempty"`
	Amount          float64    `json:"amount"`
	OriginalAmount  float64    `json:"original_amount,omitempty"`
	DiscountAmount  float64    `json:"discount_amount,omitempty"`
	NetValue        *float64   `json:"net_value,omitempty"`
	Status          string     `json:"status"`
	AsaasStatus     *string    `json:"asaas_status,omitempty"`
	PaymentMethod   *string    `json:"payment_method,omitempty"`
	StudentName     string     `json:"student_name"`
	StudentEmail    string     `json:"student_email"`
	CourseName      string     `json:"course_name"`
	PaymentDate     *time.Time `json:"payment_date,omitempty"`
	ConfirmedDate   *string    `json:"confirmed_date,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// PaymentStats represents aggregated payment statistics
type PaymentStats struct {
	TotalPayments      int     `json:"total_payments"`
	ConfirmedPayments  int     `json:"confirmed_payments"`
	PendingPayments    int     `json:"pending_payments"`
	FailedPayments     int     `json:"failed_payments"`
	RefundedPayments   int     `json:"refunded_payments"`
	TotalRevenue       float64 `json:"total_revenue"`
	ConfirmedRevenue   float64 `json:"confirmed_revenue"`
	PendingRevenue     float64 `json:"pending_revenue"`
	AverageTicket      float64 `json:"average_ticket"`
	ConversionRate     float64 `json:"conversion_rate"`
	PaymentsByMethod   map[string]int     `json:"payments_by_method"`
	RevenueByMethod    map[string]float64 `json:"revenue_by_method"`
}

// EnrollmentStats represents aggregated enrollment statistics
type EnrollmentStats struct {
	TotalEnrollments    int     `json:"total_enrollments"`
	ActiveEnrollments   int     `json:"active_enrollments"`
	PendingEnrollments  int     `json:"pending_enrollments"`
	CompletedEnrollments int    `json:"completed_enrollments"`
	CancelledEnrollments int    `json:"cancelled_enrollments"`
	ExpiredEnrollments  int     `json:"expired_enrollments"`
	AverageProgress     float64 `json:"average_progress"`
	CompletionRate      float64 `json:"completion_rate"`
	EnrollmentsByCourse map[string]int `json:"enrollments_by_course,omitempty"`
}

// AuditStats represents aggregated audit statistics
type AuditStats struct {
	TotalAudits       int     `json:"total_audits"`
	ApprovedAudits    int     `json:"approved_audits"`
	RejectedAudits    int     `json:"rejected_audits"`
	PendingAudits     int     `json:"pending_audits"`
	AverageScore      float64 `json:"average_score"`
	HighestScore      float64 `json:"highest_score"`
	LowestScore       float64 `json:"lowest_score"`
	ApprovalRate      float64 `json:"approval_rate"`
	AuditsByContract  map[string]int `json:"audits_by_contract,omitempty"`
}

// SystemOverview represents an overview of the entire system
type SystemOverview struct {
	TotalEnrollments    int     `json:"total_enrollments"`
	ActiveEnrollments   int     `json:"active_enrollments"`
	TotalRevenue        float64 `json:"total_revenue"`
	ConfirmedRevenue    float64 `json:"confirmed_revenue"`
	TotalAudits         int     `json:"total_audits"`
	AverageAuditScore   float64 `json:"average_audit_score"`
	TotalContracts      int     `json:"total_contracts"`
	TotalGestores       int     `json:"total_gestores"`
	RecentEnrollments   int     `json:"recent_enrollments"`
	RecentAudits        int     `json:"recent_audits"`
}
