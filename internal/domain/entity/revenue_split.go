package entity

import "time"

// RevenueSplit represents a revenue split calculation
type RevenueSplit struct {
	ID               string     `db:"id" json:"id"`
	EnrollmentID     string     `db:"enrollment_id" json:"enrollment_id"`
	PaymentID        string     `db:"payment_id" json:"payment_id"`
	GrossAmount      float64    `db:"gross_amount" json:"gross_amount"`
	NetAmount        float64    `db:"net_amount" json:"net_amount"`
	PlatformFee      float64    `db:"platform_fee" json:"platform_fee"`
	PaymentFee       float64    `db:"payment_fee" json:"payment_fee"`
	InstructorAmount float64    `db:"instructor_amount" json:"instructor_amount"`
	PlatformAmount   float64    `db:"platform_amount" json:"platform_amount"`
	InstructorID     *string    `db:"instructor_id" json:"instructor_id,omitempty"`
	PaymentMethod    string     `db:"payment_method" json:"payment_method"`
	Status           string     `db:"status" json:"status"`
	ProcessedAt      *time.Time `db:"processed_at" json:"processed_at,omitempty"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
}

// Revenue split status constants
const (
	RevenueSplitStatusPending   = "pending"
	RevenueSplitStatusProcessed = "processed"
	RevenueSplitStatusFailed    = "failed"
)

// PaymentFees holds the fee configuration for different payment methods
type PaymentFees struct {
	PixPercentage    float64 // 0.99%
	BoletoFixed      float64 // R$ 2.99
	CardPercentage   float64 // 2.99%
	CardFixed        float64 // R$ 0.49
}

// DefaultPaymentFees returns the default Asaas payment fees
func DefaultPaymentFees() PaymentFees {
	return PaymentFees{
		PixPercentage:  0.99,
		BoletoFixed:    2.99,
		CardPercentage: 2.99,
		CardFixed:      0.49,
	}
}

// CalculateSplitRequest represents the request to calculate revenue split
type CalculateSplitRequest struct {
	GrossAmount           float64 `json:"gross_amount" binding:"required,gt=0"`
	PaymentMethod         string  `json:"payment_method" binding:"required"`
	InstructorPercent     float64 `json:"instructor_percent,omitempty"`
	PlatformPercent       float64 `json:"platform_percent,omitempty"`
}

// CalculateSplitResponse represents the response of revenue split calculation
type CalculateSplitResponse struct {
	GrossAmount      float64 `json:"gross_amount"`
	PaymentFee       float64 `json:"payment_fee"`
	PaymentFeeDesc   string  `json:"payment_fee_description"`
	NetAmount        float64 `json:"net_amount"`
	InstructorAmount float64 `json:"instructor_amount"`
	PlatformAmount   float64 `json:"platform_amount"`
	InstructorPercent float64 `json:"instructor_percent"`
	PlatformPercent   float64 `json:"platform_percent"`
}

// CalculatePaymentFee calculates the payment gateway fee based on method
func CalculatePaymentFee(amount float64, method string, fees PaymentFees) float64 {
	switch method {
	case "pix", "PIX":
		return amount * (fees.PixPercentage / 100)
	case "boleto", "BOLETO":
		return fees.BoletoFixed
	case "card", "CREDIT_CARD", "credit_card":
		return (amount * (fees.CardPercentage / 100)) + fees.CardFixed
	default:
		return 0
	}
}

// CalculateRevenueSplit calculates how revenue is split between instructor and platform
func CalculateRevenueSplit(grossAmount float64, paymentMethod string, instructorPercent, platformPercent float64) CalculateSplitResponse {
	fees := DefaultPaymentFees()
	paymentFee := CalculatePaymentFee(grossAmount, paymentMethod, fees)
	netAmount := grossAmount - paymentFee

	instructorAmount := netAmount * (instructorPercent / 100)
	platformAmount := netAmount * (platformPercent / 100)

	var feeDesc string
	switch paymentMethod {
	case "pix", "PIX":
		feeDesc = "PIX: 0.99%"
	case "boleto", "BOLETO":
		feeDesc = "Boleto: R$ 2.99 fixo"
	case "card", "CREDIT_CARD", "credit_card":
		feeDesc = "Cartão: 2.99% + R$ 0.49"
	default:
		feeDesc = "Método desconhecido"
	}

	return CalculateSplitResponse{
		GrossAmount:       grossAmount,
		PaymentFee:        paymentFee,
		PaymentFeeDesc:    feeDesc,
		NetAmount:         netAmount,
		InstructorAmount:  instructorAmount,
		PlatformAmount:    platformAmount,
		InstructorPercent: instructorPercent,
		PlatformPercent:   platformPercent,
	}
}
