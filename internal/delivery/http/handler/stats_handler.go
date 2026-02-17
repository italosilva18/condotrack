package handler

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// StatsHandler handles statistics-related HTTP requests
type StatsHandler struct {
	db            *sqlx.DB
	matriculaRepo repository.MatriculaRepository
	auditRepo     repository.AuditRepository
	contratoRepo  repository.ContratoRepository
	gestorRepo    repository.GestorRepository
}

// NewStatsHandler creates a new stats handler
func NewStatsHandler(
	db *sqlx.DB,
	matriculaRepo repository.MatriculaRepository,
	auditRepo repository.AuditRepository,
	contratoRepo repository.ContratoRepository,
	gestorRepo repository.GestorRepository,
) *StatsHandler {
	return &StatsHandler{
		db:            db,
		matriculaRepo: matriculaRepo,
		auditRepo:     auditRepo,
		contratoRepo:  contratoRepo,
		gestorRepo:    gestorRepo,
	}
}

// GetOverview handles GET /api/v1/stats/overview
func (h *StatsHandler) GetOverview(c *gin.Context) {
	ctx := c.Request.Context()

	overview, err := h.buildOverview(ctx)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch system overview", err)
		return
	}

	response.Success(c, overview)
}

// GetEnrollmentStats handles GET /api/v1/stats/enrollments
func (h *StatsHandler) GetEnrollmentStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.buildEnrollmentStats(ctx)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch enrollment statistics", err)
		return
	}

	response.Success(c, stats)
}

// GetPaymentStats handles GET /api/v1/stats/payments
func (h *StatsHandler) GetPaymentStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.buildPaymentStats(ctx)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch payment statistics", err)
		return
	}

	response.Success(c, stats)
}

// GetAuditStats handles GET /api/v1/stats/audits
func (h *StatsHandler) GetAuditStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.buildAuditStats(ctx)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch audit statistics", err)
		return
	}

	response.Success(c, stats)
}

// buildOverview builds the system overview
func (h *StatsHandler) buildOverview(ctx context.Context) (*entity.SystemOverview, error) {
	overview := &entity.SystemOverview{}

	// Get enrollment counts
	var enrollmentStats struct {
		Total  int `db:"total"`
		Active int `db:"active"`
	}
	enrollmentQuery := `
		SELECT
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END), 0) as active
		FROM enrollments
	`
	if err := h.db.GetContext(ctx, &enrollmentStats, enrollmentQuery); err != nil {
		return nil, err
	}
	overview.TotalEnrollments = enrollmentStats.Total
	overview.ActiveEnrollments = enrollmentStats.Active

	// Get revenue
	var revenueStats struct {
		Total     float64 `db:"total"`
		Confirmed float64 `db:"confirmed"`
	}
	revenueQuery := `
		SELECT
			COALESCE(SUM(final_amount), 0) as total,
			COALESCE(SUM(CASE WHEN payment_status = 'confirmed' THEN final_amount ELSE 0 END), 0) as confirmed
		FROM enrollments
	`
	if err := h.db.GetContext(ctx, &revenueStats, revenueQuery); err != nil {
		return nil, err
	}
	overview.TotalRevenue = revenueStats.Total
	overview.ConfirmedRevenue = revenueStats.Confirmed

	// Get audit stats
	var auditStats struct {
		Total    int     `db:"total"`
		AvgScore float64 `db:"avg_score"`
	}
	auditQuery := `
		SELECT
			COUNT(*) as total,
			COALESCE(AVG(score), 0) as avg_score
		FROM audits
	`
	if err := h.db.GetContext(ctx, &auditStats, auditQuery); err != nil {
		return nil, err
	}
	overview.TotalAudits = auditStats.Total
	overview.AverageAuditScore = auditStats.AvgScore

	// Get contract count
	var contractCount int
	contractQuery := `SELECT COUNT(*) FROM contratos WHERE ativo = 1`
	if err := h.db.GetContext(ctx, &contractCount, contractQuery); err != nil {
		return nil, err
	}
	overview.TotalContracts = contractCount

	// Get gestor count
	var gestorCount int
	gestorQuery := `SELECT COUNT(*) FROM gestores WHERE ativo = 1`
	if err := h.db.GetContext(ctx, &gestorCount, gestorQuery); err != nil {
		return nil, err
	}
	overview.TotalGestores = gestorCount

	// Get recent enrollments (last 30 days)
	var recentEnrollments int
	recentEnrollmentsQuery := `
		SELECT COUNT(*) FROM enrollments
		WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
	`
	if err := h.db.GetContext(ctx, &recentEnrollments, recentEnrollmentsQuery); err != nil {
		return nil, err
	}
	overview.RecentEnrollments = recentEnrollments

	// Get recent audits (last 30 days)
	var recentAudits int
	recentAuditsQuery := `
		SELECT COUNT(*) FROM audits
		WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
	`
	if err := h.db.GetContext(ctx, &recentAudits, recentAuditsQuery); err != nil {
		return nil, err
	}
	overview.RecentAudits = recentAudits

	return overview, nil
}

// buildEnrollmentStats builds enrollment statistics
func (h *StatsHandler) buildEnrollmentStats(ctx context.Context) (*entity.EnrollmentStats, error) {
	stats := &entity.EnrollmentStats{
		EnrollmentsByCourse: make(map[string]int),
	}

	// Get status counts
	var statusStats struct {
		Total     int `db:"total"`
		Active    int `db:"active"`
		Pending   int `db:"pending"`
		Completed int `db:"completed"`
		Cancelled int `db:"cancelled"`
		Expired   int `db:"expired"`
	}
	statusQuery := `
		SELECT
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END), 0) as active,
			COALESCE(SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END), 0) as pending,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END), 0) as completed,
			COALESCE(SUM(CASE WHEN status = 'cancelled' THEN 1 ELSE 0 END), 0) as cancelled,
			COALESCE(SUM(CASE WHEN status = 'expired' THEN 1 ELSE 0 END), 0) as expired
		FROM enrollments
	`
	if err := h.db.GetContext(ctx, &statusStats, statusQuery); err != nil {
		return nil, err
	}
	stats.TotalEnrollments = statusStats.Total
	stats.ActiveEnrollments = statusStats.Active
	stats.PendingEnrollments = statusStats.Pending
	stats.CompletedEnrollments = statusStats.Completed
	stats.CancelledEnrollments = statusStats.Cancelled
	stats.ExpiredEnrollments = statusStats.Expired

	// Get average progress
	var avgProgress float64
	progressQuery := `SELECT COALESCE(AVG(progress), 0) FROM enrollments WHERE status = 'active'`
	if err := h.db.GetContext(ctx, &avgProgress, progressQuery); err != nil {
		return nil, err
	}
	stats.AverageProgress = avgProgress

	// Calculate completion rate
	if stats.TotalEnrollments > 0 {
		stats.CompletionRate = float64(stats.CompletedEnrollments) / float64(stats.TotalEnrollments) * 100
	}

	// Get enrollments by course
	type courseCount struct {
		CourseName string `db:"course_name"`
		Count      int    `db:"count"`
	}
	var courseCounts []courseCount
	courseQuery := `
		SELECT course_name, COUNT(*) as count
		FROM enrollments
		GROUP BY course_name
		ORDER BY count DESC
		LIMIT 10
	`
	if err := h.db.SelectContext(ctx, &courseCounts, courseQuery); err != nil {
		return nil, err
	}
	for _, cc := range courseCounts {
		stats.EnrollmentsByCourse[cc.CourseName] = cc.Count
	}

	return stats, nil
}

// buildPaymentStats builds payment statistics
func (h *StatsHandler) buildPaymentStats(ctx context.Context) (*entity.PaymentStats, error) {
	stats := &entity.PaymentStats{
		PaymentsByMethod: make(map[string]int),
		RevenueByMethod:  make(map[string]float64),
	}

	// Get payment status counts and revenue
	var paymentStats struct {
		Total            int     `db:"total"`
		Confirmed        int     `db:"confirmed"`
		Pending          int     `db:"pending"`
		Failed           int     `db:"failed"`
		Refunded         int     `db:"refunded"`
		TotalRevenue     float64 `db:"total_revenue"`
		ConfirmedRevenue float64 `db:"confirmed_revenue"`
		PendingRevenue   float64 `db:"pending_revenue"`
	}
	paymentQuery := `
		SELECT
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN payment_status = 'confirmed' THEN 1 ELSE 0 END), 0) as confirmed,
			COALESCE(SUM(CASE WHEN payment_status = 'pending' THEN 1 ELSE 0 END), 0) as pending,
			COALESCE(SUM(CASE WHEN payment_status = 'failed' THEN 1 ELSE 0 END), 0) as failed,
			COALESCE(SUM(CASE WHEN payment_status = 'refunded' THEN 1 ELSE 0 END), 0) as refunded,
			COALESCE(SUM(final_amount), 0) as total_revenue,
			COALESCE(SUM(CASE WHEN payment_status = 'confirmed' THEN final_amount ELSE 0 END), 0) as confirmed_revenue,
			COALESCE(SUM(CASE WHEN payment_status = 'pending' THEN final_amount ELSE 0 END), 0) as pending_revenue
		FROM enrollments
	`
	if err := h.db.GetContext(ctx, &paymentStats, paymentQuery); err != nil {
		return nil, err
	}
	stats.TotalPayments = paymentStats.Total
	stats.ConfirmedPayments = paymentStats.Confirmed
	stats.PendingPayments = paymentStats.Pending
	stats.FailedPayments = paymentStats.Failed
	stats.RefundedPayments = paymentStats.Refunded
	stats.TotalRevenue = paymentStats.TotalRevenue
	stats.ConfirmedRevenue = paymentStats.ConfirmedRevenue
	stats.PendingRevenue = paymentStats.PendingRevenue

	// Calculate average ticket
	if stats.ConfirmedPayments > 0 {
		stats.AverageTicket = stats.ConfirmedRevenue / float64(stats.ConfirmedPayments)
	}

	// Calculate conversion rate
	if stats.TotalPayments > 0 {
		stats.ConversionRate = float64(stats.ConfirmedPayments) / float64(stats.TotalPayments) * 100
	}

	// Get payments by method
	type methodCount struct {
		Method  string  `db:"payment_method"`
		Count   int     `db:"count"`
		Revenue float64 `db:"revenue"`
	}
	var methodCounts []methodCount
	methodQuery := `
		SELECT
			COALESCE(payment_method, 'unknown') as payment_method,
			COUNT(*) as count,
			COALESCE(SUM(final_amount), 0) as revenue
		FROM enrollments
		WHERE payment_status = 'confirmed'
		GROUP BY payment_method
	`
	if err := h.db.SelectContext(ctx, &methodCounts, methodQuery); err != nil {
		return nil, err
	}
	for _, mc := range methodCounts {
		stats.PaymentsByMethod[mc.Method] = mc.Count
		stats.RevenueByMethod[mc.Method] = mc.Revenue
	}

	return stats, nil
}

// buildAuditStats builds audit statistics
func (h *StatsHandler) buildAuditStats(ctx context.Context) (*entity.AuditStats, error) {
	stats := &entity.AuditStats{
		AuditsByContract: make(map[string]int),
	}

	// Get audit status counts and scores
	var auditStats struct {
		Total        int     `db:"total"`
		Approved     int     `db:"approved"`
		Rejected     int     `db:"rejected"`
		Pending      int     `db:"pending"`
		AvgScore     float64 `db:"avg_score"`
		HighestScore float64 `db:"highest_score"`
		LowestScore  float64 `db:"lowest_score"`
	}
	auditQuery := `
		SELECT
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN status = 'approved' THEN 1 ELSE 0 END), 0) as approved,
			COALESCE(SUM(CASE WHEN status = 'rejected' THEN 1 ELSE 0 END), 0) as rejected,
			COALESCE(SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END), 0) as pending,
			COALESCE(AVG(score), 0) as avg_score,
			COALESCE(MAX(score), 0) as highest_score,
			COALESCE(MIN(score), 0) as lowest_score
		FROM audits
	`
	if err := h.db.GetContext(ctx, &auditStats, auditQuery); err != nil {
		return nil, err
	}
	stats.TotalAudits = auditStats.Total
	stats.ApprovedAudits = auditStats.Approved
	stats.RejectedAudits = auditStats.Rejected
	stats.PendingAudits = auditStats.Pending
	stats.AverageScore = auditStats.AvgScore
	stats.HighestScore = auditStats.HighestScore
	stats.LowestScore = auditStats.LowestScore

	// Calculate approval rate
	if stats.TotalAudits > 0 {
		stats.ApprovalRate = float64(stats.ApprovedAudits) / float64(stats.TotalAudits) * 100
	}

	// Get audits by contract
	type contractCount struct {
		ContractName string `db:"contract_name"`
		Count        int    `db:"count"`
	}
	var contractCounts []contractCount
	contractQuery := `
		SELECT c.nome as contract_name, COUNT(*) as count
		FROM audits a
		INNER JOIN contratos c ON c.id = a.contract_id
		GROUP BY c.id, c.nome
		ORDER BY count DESC
		LIMIT 10
	`
	if err := h.db.SelectContext(ctx, &contractCounts, contractQuery); err != nil {
		return nil, err
	}
	for _, cc := range contractCounts {
		stats.AuditsByContract[cc.ContractName] = cc.Count
	}

	return stats, nil
}
