package mercadopago

import (
	"context"
	"encoding/json"
	"fmt"
)

// CreatePayment creates a payment on Mercado Pago.
func (c *Client) CreatePayment(ctx context.Context, req *MPCreatePaymentRequest) (*MPPaymentResponse, error) {
	respBody, err := c.post(ctx, "/v1/payments", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create MP payment: %w", err)
	}

	var resp MPPaymentResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse MP payment response: %w", err)
	}

	return &resp, nil
}

// GetPayment retrieves a payment by ID from Mercado Pago.
func (c *Client) GetPayment(ctx context.Context, paymentID string) (*MPPaymentResponse, error) {
	respBody, err := c.get(ctx, fmt.Sprintf("/v1/payments/%s", paymentID))
	if err != nil {
		return nil, fmt.Errorf("failed to get MP payment: %w", err)
	}

	var resp MPPaymentResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse MP payment response: %w", err)
	}

	return &resp, nil
}

// RefundPayment creates a refund for a payment on Mercado Pago.
func (c *Client) RefundPayment(ctx context.Context, paymentID string, amount float64) (*MPPaymentResponse, error) {
	body := map[string]interface{}{}
	if amount > 0 {
		body["amount"] = amount
	}

	_, err := c.post(ctx, fmt.Sprintf("/v1/payments/%s/refunds", paymentID), body)
	if err != nil {
		return nil, fmt.Errorf("failed to refund MP payment: %w", err)
	}

	// Get updated payment after refund
	return c.GetPayment(ctx, paymentID)
}

// CancelPayment cancels a pending payment on Mercado Pago.
func (c *Client) CancelPayment(ctx context.Context, paymentID string) error {
	body := map[string]string{
		"status": "cancelled",
	}
	_, err := c.put(ctx, fmt.Sprintf("/v1/payments/%s", paymentID), body)
	if err != nil {
		return fmt.Errorf("failed to cancel MP payment: %w", err)
	}
	return nil
}
