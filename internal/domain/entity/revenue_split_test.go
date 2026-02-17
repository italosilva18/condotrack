package entity

import (
	"math"
	"testing"
)

func almostEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestCalculatePaymentFee_PIX(t *testing.T) {
	fees := DefaultPaymentFees()
	fee := CalculatePaymentFee(100, "pix", fees)
	expected := 100 * (0.99 / 100) // 0.99%
	if !almostEqual(fee, expected, 0.01) {
		t.Errorf("PIX fee: expected %f, got %f", expected, fee)
	}
}

func TestCalculatePaymentFee_Boleto(t *testing.T) {
	fees := DefaultPaymentFees()
	fee := CalculatePaymentFee(100, "boleto", fees)
	if fee != 2.99 {
		t.Errorf("Boleto fee: expected 2.99, got %f", fee)
	}
}

func TestCalculatePaymentFee_Card(t *testing.T) {
	fees := DefaultPaymentFees()
	fee := CalculatePaymentFee(100, "card", fees)
	expected := (100 * (2.99 / 100)) + 0.49 // 2.99% + R$0.49
	if !almostEqual(fee, expected, 0.01) {
		t.Errorf("Card fee: expected %f, got %f", expected, fee)
	}
}

func TestCalculatePaymentFee_Unknown(t *testing.T) {
	fees := DefaultPaymentFees()
	fee := CalculatePaymentFee(100, "crypto", fees)
	if fee != 0 {
		t.Errorf("Unknown method fee: expected 0, got %f", fee)
	}
}

func TestCalculateRevenueSplit_PIX(t *testing.T) {
	result := CalculateRevenueSplit(100, "pix", 70, 30)

	if result.GrossAmount != 100 {
		t.Errorf("GrossAmount: expected 100, got %f", result.GrossAmount)
	}

	expectedFee := 100 * (0.99 / 100)
	if !almostEqual(result.PaymentFee, expectedFee, 0.01) {
		t.Errorf("PaymentFee: expected %f, got %f", expectedFee, result.PaymentFee)
	}

	expectedNet := 100 - expectedFee
	if !almostEqual(result.NetAmount, expectedNet, 0.01) {
		t.Errorf("NetAmount: expected %f, got %f", expectedNet, result.NetAmount)
	}

	expectedInstructor := expectedNet * 0.7
	if !almostEqual(result.InstructorAmount, expectedInstructor, 0.01) {
		t.Errorf("InstructorAmount: expected %f, got %f", expectedInstructor, result.InstructorAmount)
	}

	expectedPlatform := expectedNet * 0.3
	if !almostEqual(result.PlatformAmount, expectedPlatform, 0.01) {
		t.Errorf("PlatformAmount: expected %f, got %f", expectedPlatform, result.PlatformAmount)
	}

	if result.InstructorPercent != 70 {
		t.Errorf("InstructorPercent: expected 70, got %f", result.InstructorPercent)
	}
	if result.PlatformPercent != 30 {
		t.Errorf("PlatformPercent: expected 30, got %f", result.PlatformPercent)
	}
}

func TestCalculateRevenueSplit_Boleto(t *testing.T) {
	result := CalculateRevenueSplit(100, "boleto", 70, 30)

	if result.PaymentFee != 2.99 {
		t.Errorf("Boleto PaymentFee: expected 2.99, got %f", result.PaymentFee)
	}

	expectedNet := 100 - 2.99
	if !almostEqual(result.NetAmount, expectedNet, 0.01) {
		t.Errorf("NetAmount: expected %f, got %f", expectedNet, result.NetAmount)
	}
}

func TestCalculateRevenueSplit_Card(t *testing.T) {
	result := CalculateRevenueSplit(200, "card", 60, 40)

	expectedFee := (200 * (2.99 / 100)) + 0.49
	if !almostEqual(result.PaymentFee, expectedFee, 0.01) {
		t.Errorf("Card PaymentFee: expected %f, got %f", expectedFee, result.PaymentFee)
	}

	expectedNet := 200 - expectedFee
	expectedInstructor := expectedNet * 0.6
	expectedPlatform := expectedNet * 0.4

	if !almostEqual(result.InstructorAmount, expectedInstructor, 0.01) {
		t.Errorf("InstructorAmount: expected %f, got %f", expectedInstructor, result.InstructorAmount)
	}
	if !almostEqual(result.PlatformAmount, expectedPlatform, 0.01) {
		t.Errorf("PlatformAmount: expected %f, got %f", expectedPlatform, result.PlatformAmount)
	}
}

func TestCalculateRevenueSplit_FeeDescription(t *testing.T) {
	tests := []struct {
		method   string
		expected string
	}{
		{"pix", "PIX: 0.99%"},
		{"boleto", "Boleto: R$ 2.99 fixo"},
		{"card", "Cartão: 2.99% + R$ 0.49"},
		{"unknown", "Método desconhecido"},
	}

	for _, tt := range tests {
		result := CalculateRevenueSplit(100, tt.method, 70, 30)
		if result.PaymentFeeDesc != tt.expected {
			t.Errorf("method %s: expected desc '%s', got '%s'", tt.method, tt.expected, result.PaymentFeeDesc)
		}
	}
}

func TestCalculateRevenueSplit_CustomPercents(t *testing.T) {
	result := CalculateRevenueSplit(100, "pix", 80, 20)

	if result.InstructorPercent != 80 {
		t.Errorf("InstructorPercent: expected 80, got %f", result.InstructorPercent)
	}
	if result.PlatformPercent != 20 {
		t.Errorf("PlatformPercent: expected 20, got %f", result.PlatformPercent)
	}
}
