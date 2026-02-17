package payment

import (
	"context"
	"testing"

	"github.com/condotrack/api/internal/config"
	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/gateway"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/condotrack/api/internal/testutil"
)

func newTestUseCase() (UseCase, *testutil.MockGateway, *testutil.MockPaymentRepository) {
	mockGw := &testutil.MockGateway{}
	mockRepo := testutil.NewMockPaymentRepository()
	cfg := &config.Config{
		RevenueInstructorPercent: 70,
		RevenuePlatformPercent:   30,
	}
	uc := NewUseCase(mockGw, mockRepo, cfg)
	return uc, mockGw, mockRepo
}

func TestCreateCustomer_Success(t *testing.T) {
	uc, mockGw, _ := newTestUseCase()
	mockGw.CreateCustomerFunc = func(ctx context.Context, req gateway.CreateCustomerRequest) (*gateway.CustomerResponse, error) {
		return &gateway.CustomerResponse{
			GatewayID: "cust_123",
			Name:      req.Name,
			Email:     req.Email,
		}, nil
	}

	resp, err := uc.CreateCustomer(context.Background(), &CreateCustomerRequest{
		Name:     "Test User",
		Email:    "test@test.com",
		Document: "12345678900",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.GatewayID != "cust_123" {
		t.Errorf("expected GatewayID cust_123, got %s", resp.GatewayID)
	}
}

func TestCreateCustomer_MissingName(t *testing.T) {
	uc, _, _ := newTestUseCase()
	_, err := uc.CreateCustomer(context.Background(), &CreateCustomerRequest{
		Document: "12345678900",
	})
	if err == nil {
		t.Error("expected error for missing name")
	}
}

func TestCreateCustomer_MissingDocument(t *testing.T) {
	uc, _, _ := newTestUseCase()
	_, err := uc.CreateCustomer(context.Background(), &CreateCustomerRequest{
		Name: "Test User",
	})
	if err == nil {
		t.Error("expected error for missing document")
	}
}

func TestCreatePixPayment_Success(t *testing.T) {
	uc, mockGw, _ := newTestUseCase()
	mockGw.CreatePixPaymentFunc = func(ctx context.Context, req gateway.CreatePaymentRequest) (*gateway.PaymentResponse, error) {
		return &gateway.PaymentResponse{
			GatewayPaymentID: "pay_pix_123",
			Status:           gateway.StatusPending,
			Amount:           req.Amount,
			PixQRCodeBase64:  "base64qrcode",
			PixCopyPaste:     "copypaste",
		}, nil
	}

	resp, err := uc.CreatePixPayment(context.Background(), &CreatePaymentRequest{
		CustomerGatewayID: "cust_123",
		Amount:            100,
		DueDate:           "2026-03-01",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.GatewayPaymentID != "pay_pix_123" {
		t.Errorf("expected pay_pix_123, got %s", resp.GatewayPaymentID)
	}
	if resp.PixQRCodeBase64 != "base64qrcode" {
		t.Error("expected PIX QR code in response")
	}
}

func TestCreatePixPayment_MissingCustomerID(t *testing.T) {
	uc, _, _ := newTestUseCase()
	_, err := uc.CreatePixPayment(context.Background(), &CreatePaymentRequest{
		Amount: 100,
	})
	if err == nil {
		t.Error("expected error for missing customer ID")
	}
}

func TestCreatePixPayment_ZeroAmount(t *testing.T) {
	uc, _, _ := newTestUseCase()
	_, err := uc.CreatePixPayment(context.Background(), &CreatePaymentRequest{
		CustomerGatewayID: "cust_123",
		Amount:            0,
	})
	if err == nil {
		t.Error("expected error for zero amount")
	}
}

func TestGetPaymentStatus_FromLocalDB(t *testing.T) {
	uc, _, mockRepo := newTestUseCase()
	gwPaymentID := "pay_gw_123"
	mockRepo.Payments["p1"] = &entity.Payment{
		ID:               "p1",
		EnrollmentID:     "e1",
		GrossAmount:      100,
		NetAmount:        90,
		Status:           entity.FinPaymentStatusConfirmed,
		Gateway:          entity.GatewayAsaas,
		GatewayPaymentID: &gwPaymentID,
	}

	resp, err := uc.GetPaymentStatus(context.Background(), "p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "p1" {
		t.Errorf("expected ID p1, got %s", resp.ID)
	}
	if resp.Status != entity.FinPaymentStatusConfirmed {
		t.Errorf("expected status confirmed, got %s", resp.Status)
	}
	if resp.Gateway != entity.GatewayAsaas {
		t.Errorf("expected gateway asaas, got %s", resp.Gateway)
	}
}

func TestGetPaymentStatus_FallbackToGateway(t *testing.T) {
	uc, mockGw, _ := newTestUseCase()
	mockGw.GetPaymentFunc = func(ctx context.Context, gatewayPaymentID string) (*gateway.PaymentResponse, error) {
		return &gateway.PaymentResponse{
			GatewayPaymentID: gatewayPaymentID,
			Status:           gateway.StatusConfirmed,
			Amount:           200,
			NetAmount:        195,
			BillingType:      gateway.BillingPIX,
		}, nil
	}

	resp, err := uc.GetPaymentStatus(context.Background(), "pay_gw_456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.GatewayPaymentID != "pay_gw_456" {
		t.Errorf("expected gateway payment ID pay_gw_456, got %s", resp.GatewayPaymentID)
	}
	if resp.Status != gateway.StatusConfirmed {
		t.Errorf("expected status confirmed, got %s", resp.Status)
	}
}

func TestGetPaymentStatus_EmptyID(t *testing.T) {
	uc, _, _ := newTestUseCase()
	_, err := uc.GetPaymentStatus(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty payment ID")
	}
}

func TestSimulateRevenueSplit_PIX(t *testing.T) {
	uc, _, _ := newTestUseCase()
	result := uc.SimulateRevenueSplit(&entity.CalculateSplitRequest{
		GrossAmount:   100,
		PaymentMethod: "pix",
	})

	if result.GrossAmount != 100 {
		t.Errorf("expected GrossAmount 100, got %f", result.GrossAmount)
	}
	if result.InstructorPercent != 70 {
		t.Errorf("expected InstructorPercent 70, got %f", result.InstructorPercent)
	}
	if result.PlatformPercent != 30 {
		t.Errorf("expected PlatformPercent 30, got %f", result.PlatformPercent)
	}
}

func TestSimulateRevenueSplit_CustomPercents(t *testing.T) {
	uc, _, _ := newTestUseCase()
	result := uc.SimulateRevenueSplit(&entity.CalculateSplitRequest{
		GrossAmount:       200,
		PaymentMethod:     "card",
		InstructorPercent: 60,
		PlatformPercent:   40,
	})

	if result.InstructorPercent != 60 {
		t.Errorf("expected InstructorPercent 60, got %f", result.InstructorPercent)
	}
	if result.PlatformPercent != 40 {
		t.Errorf("expected PlatformPercent 40, got %f", result.PlatformPercent)
	}
}

func TestListPayments(t *testing.T) {
	uc, _, mockRepo := newTestUseCase()
	mockRepo.Payments["p1"] = &entity.Payment{ID: "p1", EnrollmentID: "e1", Status: "pending"}
	mockRepo.Payments["p2"] = &entity.Payment{ID: "p2", EnrollmentID: "e2", Status: "confirmed"}

	payments, total, err := uc.ListPayments(context.Background(), repository.PaymentFilters{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(payments) != 2 {
		t.Errorf("expected 2 payments, got %d", len(payments))
	}
}

func TestGetPaymentsByEnrollment(t *testing.T) {
	uc, _, mockRepo := newTestUseCase()
	mockRepo.Payments["p1"] = &entity.Payment{ID: "p1", EnrollmentID: "e1", Status: "pending"}
	mockRepo.Payments["p2"] = &entity.Payment{ID: "p2", EnrollmentID: "e1", Status: "confirmed"}
	mockRepo.Payments["p3"] = &entity.Payment{ID: "p3", EnrollmentID: "e2", Status: "pending"}

	payments, err := uc.GetPaymentsByEnrollment(context.Background(), "e1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 2 {
		t.Errorf("expected 2 payments for e1, got %d", len(payments))
	}
}
