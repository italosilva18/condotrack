package certificado

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/google/uuid"
)

// UseCase defines the certificado use case interface
type UseCase interface {
	GetCertificatesByStudent(ctx context.Context, studentID string) ([]entity.Certificado, error)
	GetCertificateByID(ctx context.Context, id string) (*entity.Certificado, error)
	ValidateCertificate(ctx context.Context, code string) (*entity.CertificateResponse, error)
	GenerateCertificate(ctx context.Context, enrollmentID string) (*entity.Certificado, error)
}

type certificadoUseCase struct {
	certRepo      repository.CertificadoRepository
	matriculaRepo repository.MatriculaRepository
}

// NewUseCase creates a new certificado use case
func NewUseCase(certRepo repository.CertificadoRepository, matriculaRepo repository.MatriculaRepository) UseCase {
	return &certificadoUseCase{
		certRepo:      certRepo,
		matriculaRepo: matriculaRepo,
	}
}

// GetCertificatesByStudent returns all certificates for a student
func (uc *certificadoUseCase) GetCertificatesByStudent(ctx context.Context, studentID string) ([]entity.Certificado, error) {
	return uc.certRepo.FindByStudentID(ctx, studentID)
}

// GetCertificateByID returns a specific certificate by ID
func (uc *certificadoUseCase) GetCertificateByID(ctx context.Context, id string) (*entity.Certificado, error) {
	return uc.certRepo.FindByID(ctx, id)
}

// ValidateCertificate validates a certificate by its code
func (uc *certificadoUseCase) ValidateCertificate(ctx context.Context, code string) (*entity.CertificateResponse, error) {
	cert, err := uc.certRepo.FindByValidationCode(ctx, strings.ToUpper(code))
	if err != nil {
		return nil, err
	}
	if cert == nil {
		return &entity.CertificateResponse{
			Valid: false,
		}, nil
	}

	valid := cert.Status == entity.CertificateStatusActive

	response := &entity.CertificateResponse{
		ID:             cert.ID,
		StudentName:    cert.StudentName,
		CourseName:     cert.CourseName,
		CourseHours:    cert.CourseHours,
		CompletionDate: cert.CompletionDate.Format("02/01/2006"),
		IssueDate:      cert.IssueDate.Format("02/01/2006"),
		ValidationCode: cert.ValidationCode,
		ValidationURL:  fmt.Sprintf("/api/v1/certificados/validate/%s", cert.ValidationCode),
		Valid:          valid,
	}

	if cert.DownloadURL != nil {
		response.DownloadURL = *cert.DownloadURL
	}

	return response, nil
}

// GenerateCertificate generates a certificate for a completed enrollment
func (uc *certificadoUseCase) GenerateCertificate(ctx context.Context, enrollmentID string) (*entity.Certificado, error) {
	// Check if certificate already exists
	existing, err := uc.certRepo.FindByEnrollmentID(ctx, enrollmentID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	// Get enrollment
	enrollment, err := uc.matriculaRepo.FindByID(ctx, enrollmentID)
	if err != nil {
		return nil, err
	}
	if enrollment == nil {
		return nil, errors.New("enrollment not found")
	}

	// Verify enrollment is completed
	if enrollment.Status != entity.EnrollmentStatusCompleted {
		return nil, errors.New("enrollment is not completed")
	}

	// Verify payment is confirmed
	if enrollment.PaymentStatus != entity.PaymentStatusConfirmed {
		return nil, errors.New("payment is not confirmed")
	}

	// Generate validation code
	validationCode := generateValidationCode()

	// Set completion date
	completionDate := time.Now()
	if enrollment.CompletionDate != nil {
		completionDate = *enrollment.CompletionDate
	}

	// Create certificate
	cert := &entity.Certificado{
		ID:             uuid.New().String(),
		EnrollmentID:   enrollmentID,
		StudentID:      enrollment.StudentID,
		StudentName:    enrollment.StudentName,
		StudentCPF:     enrollment.StudentCPF,
		CourseID:       enrollment.CourseID,
		CourseName:     enrollment.CourseName,
		CourseHours:    40, // Default course hours, could be from course entity
		InstructorName: enrollment.InstructorName,
		CompletionDate: completionDate,
		IssueDate:      time.Now(),
		ValidationCode: validationCode,
		Status:         entity.CertificateStatusActive,
		CreatedAt:      time.Now(),
	}

	if err := uc.certRepo.Create(ctx, cert); err != nil {
		return nil, err
	}

	// Update enrollment with certificate ID
	enrollment.CertificateID = &cert.ID
	if err := uc.matriculaRepo.Update(ctx, enrollment); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to update enrollment with certificate ID: %v\n", err)
	}

	return cert, nil
}

// generateValidationCode generates a unique validation code
func generateValidationCode() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	code := hex.EncodeToString(bytes)
	return strings.ToUpper(code)
}
