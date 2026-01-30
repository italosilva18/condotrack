package asaas

import (
	"context"
	"encoding/json"
	"fmt"
)

// CreatePayment creates a new payment (PIX or Boleto)
func (c *Client) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*PaymentResponse, error) {
	respBody, err := c.post(ctx, "/payments", req)
	if err != nil {
		return nil, err
	}

	var payment PaymentResponse
	if err := json.Unmarshal(respBody, &payment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payment response: %w", err)
	}

	// Get PIX QR code if it's a PIX payment
	if req.BillingType == "PIX" {
		qrCode, err := c.GetPixQRCode(ctx, payment.ID)
		if err == nil {
			payment.PixQRCode = qrCode
		}
	}

	return &payment, nil
}

// CreateCardPayment creates a credit card payment
func (c *Client) CreateCardPayment(ctx context.Context, req *CreateCardPaymentRequest) (*PaymentResponse, error) {
	respBody, err := c.post(ctx, "/payments", req)
	if err != nil {
		return nil, err
	}

	var payment PaymentResponse
	if err := json.Unmarshal(respBody, &payment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payment response: %w", err)
	}

	return &payment, nil
}

// GetPayment retrieves a payment by ID
func (c *Client) GetPayment(ctx context.Context, paymentID string) (*PaymentResponse, error) {
	respBody, err := c.get(ctx, "/payments/"+paymentID)
	if err != nil {
		return nil, err
	}

	var payment PaymentResponse
	if err := json.Unmarshal(respBody, &payment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payment response: %w", err)
	}

	return &payment, nil
}

// GetPixQRCode retrieves the PIX QR code for a payment
func (c *Client) GetPixQRCode(ctx context.Context, paymentID string) (*PixQRCode, error) {
	respBody, err := c.get(ctx, "/payments/"+paymentID+"/pixQrCode")
	if err != nil {
		return nil, err
	}

	var qrCode PixQRCode
	if err := json.Unmarshal(respBody, &qrCode); err != nil {
		return nil, fmt.Errorf("failed to unmarshal QR code response: %w", err)
	}

	return &qrCode, nil
}

// GetBoletoIdentificationField retrieves the boleto barcode
func (c *Client) GetBoletoIdentificationField(ctx context.Context, paymentID string) (string, error) {
	respBody, err := c.get(ctx, "/payments/"+paymentID+"/identificationField")
	if err != nil {
		return "", err
	}

	var response struct {
		IdentificationField string `json:"identificationField"`
		BarCode             string `json:"barCode"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal barcode response: %w", err)
	}

	return response.IdentificationField, nil
}

// RefundPayment refunds a payment
func (c *Client) RefundPayment(ctx context.Context, paymentID string, value float64) (*PaymentResponse, error) {
	body := map[string]interface{}{
		"value": value,
	}

	respBody, err := c.post(ctx, "/payments/"+paymentID+"/refund", body)
	if err != nil {
		return nil, err
	}

	var payment PaymentResponse
	if err := json.Unmarshal(respBody, &payment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payment response: %w", err)
	}

	return &payment, nil
}

// DeletePayment deletes/cancels a payment
func (c *Client) DeletePayment(ctx context.Context, paymentID string) error {
	_, err := c.delete(ctx, "/payments/"+paymentID)
	return err
}

// ReceivePaymentInCash marks a payment as received in cash
func (c *Client) ReceivePaymentInCash(ctx context.Context, paymentID string, paymentDate string, value float64) (*PaymentResponse, error) {
	body := map[string]interface{}{
		"paymentDate": paymentDate,
		"value":       value,
	}

	respBody, err := c.post(ctx, "/payments/"+paymentID+"/receiveInCash", body)
	if err != nil {
		return nil, err
	}

	var payment PaymentResponse
	if err := json.Unmarshal(respBody, &payment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payment response: %w", err)
	}

	return &payment, nil
}

// UndoReceivedInCash reverts a payment marked as received in cash
func (c *Client) UndoReceivedInCash(ctx context.Context, paymentID string) (*PaymentResponse, error) {
	respBody, err := c.post(ctx, "/payments/"+paymentID+"/undoReceivedInCash", nil)
	if err != nil {
		return nil, err
	}

	var payment PaymentResponse
	if err := json.Unmarshal(respBody, &payment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payment response: %w", err)
	}

	return &payment, nil
}
