package handler

import (
	"github.com/condotrack/api/internal/usecase/certificado"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// CertificadoHandler handles certificate-related HTTP requests
type CertificadoHandler struct {
	usecase certificado.UseCase
}

// NewCertificadoHandler creates a new certificado handler
func NewCertificadoHandler(uc certificado.UseCase) *CertificadoHandler {
	return &CertificadoHandler{usecase: uc}
}

// GetCertificatesByStudent handles GET /api/v1/certificados/:aluno_id
func (h *CertificadoHandler) GetCertificatesByStudent(c *gin.Context) {
	ctx := c.Request.Context()
	studentID := c.Param("aluno_id")

	certificates, err := h.usecase.GetCertificatesByStudent(ctx, studentID)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch certificates", err)
		return
	}

	response.Success(c, certificates)
}

// GetCertificateByID handles GET /api/v1/certificados/detail/:id
func (h *CertificadoHandler) GetCertificateByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	certificate, err := h.usecase.GetCertificateByID(ctx, id)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch certificate", err)
		return
	}

	if certificate == nil {
		response.NotFound(c, "Certificate not found")
		return
	}

	response.Success(c, certificate)
}

// ValidateCertificate handles GET /api/v1/certificados/validate/:code
func (h *CertificadoHandler) ValidateCertificate(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	result, err := h.usecase.ValidateCertificate(ctx, code)
	if err != nil {
		response.SafeInternalError(c, "Failed to validate certificate", err)
		return
	}

	response.Success(c, result)
}

// GenerateCertificate handles POST /api/v1/certificados/generate
func (h *CertificadoHandler) GenerateCertificate(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		EnrollmentID string `json:"enrollment_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	certificate, err := h.usecase.GenerateCertificate(ctx, req.EnrollmentID)
	if err != nil {
		response.SafeInternalError(c, "Failed to generate certificate", err)
		return
	}

	response.Created(c, certificate)
}
