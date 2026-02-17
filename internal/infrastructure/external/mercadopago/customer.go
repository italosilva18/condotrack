package mercadopago

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// CreateCustomer creates a customer on Mercado Pago.
func (c *Client) CreateCustomer(ctx context.Context, customer *MPCustomer) (*MPCustomer, error) {
	respBody, err := c.post(ctx, "/v1/customers", customer)
	if err != nil {
		return nil, fmt.Errorf("failed to create MP customer: %w", err)
	}

	var resp MPCustomer
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse MP customer response: %w", err)
	}

	return &resp, nil
}

// FindCustomerByEmail searches for a customer by email on Mercado Pago.
func (c *Client) FindCustomerByEmail(ctx context.Context, email string) (*MPCustomer, error) {
	path := fmt.Sprintf("/v1/customers/search?email=%s", url.QueryEscape(email))
	respBody, err := c.get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to search MP customer: %w", err)
	}

	var result MPCustomerSearchResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse MP customer search: %w", err)
	}

	if len(result.Results) == 0 {
		return nil, nil
	}

	return &result.Results[0], nil
}

// FindOrCreateCustomer finds an existing customer by email or creates one.
func (c *Client) FindOrCreateCustomer(ctx context.Context, customer *MPCustomer) (*MPCustomer, error) {
	existing, err := c.FindCustomerByEmail(ctx, customer.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	return c.CreateCustomer(ctx, customer)
}
