package asaas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// CreateCustomer creates a new customer in Asaas
func (c *Client) CreateCustomer(ctx context.Context, req *CreateCustomerRequest) (*Customer, error) {
	respBody, err := c.post(ctx, "/customers", req)
	if err != nil {
		return nil, err
	}

	var customer Customer
	if err := json.Unmarshal(respBody, &customer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal customer response: %w", err)
	}

	return &customer, nil
}

// GetCustomer retrieves a customer by ID
func (c *Client) GetCustomer(ctx context.Context, customerID string) (*Customer, error) {
	respBody, err := c.get(ctx, "/customers/"+customerID)
	if err != nil {
		return nil, err
	}

	var customer Customer
	if err := json.Unmarshal(respBody, &customer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal customer response: %w", err)
	}

	return &customer, nil
}

// FindCustomerByCPF finds a customer by CPF/CNPJ
func (c *Client) FindCustomerByCPF(ctx context.Context, cpfCnpj string) (*Customer, error) {
	path := fmt.Sprintf("/customers?cpfCnpj=%s", url.QueryEscape(cpfCnpj))
	respBody, err := c.get(ctx, path)
	if err != nil {
		return nil, err
	}

	var response CustomerListResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal customer list response: %w", err)
	}

	if len(response.Data) == 0 {
		return nil, nil
	}

	return &response.Data[0], nil
}

// FindCustomerByEmail finds a customer by email
func (c *Client) FindCustomerByEmail(ctx context.Context, email string) (*Customer, error) {
	path := fmt.Sprintf("/customers?email=%s", url.QueryEscape(email))
	respBody, err := c.get(ctx, path)
	if err != nil {
		return nil, err
	}

	var response CustomerListResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal customer list response: %w", err)
	}

	if len(response.Data) == 0 {
		return nil, nil
	}

	return &response.Data[0], nil
}

// FindOrCreateCustomer finds an existing customer or creates a new one
func (c *Client) FindOrCreateCustomer(ctx context.Context, req *CreateCustomerRequest) (*Customer, error) {
	// Try to find by CPF/CNPJ first
	if req.CPFCnpj != "" {
		existing, err := c.FindCustomerByCPF(ctx, req.CPFCnpj)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return existing, nil
		}
	}

	// Try to find by email
	if req.Email != "" {
		existing, err := c.FindCustomerByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return existing, nil
		}
	}

	// Create new customer
	return c.CreateCustomer(ctx, req)
}

// UpdateCustomer updates an existing customer
func (c *Client) UpdateCustomer(ctx context.Context, customerID string, req *CreateCustomerRequest) (*Customer, error) {
	respBody, err := c.put(ctx, "/customers/"+customerID, req)
	if err != nil {
		return nil, err
	}

	var customer Customer
	if err := json.Unmarshal(respBody, &customer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal customer response: %w", err)
	}

	return &customer, nil
}

// DeleteCustomer deletes a customer
func (c *Client) DeleteCustomer(ctx context.Context, customerID string) error {
	_, err := c.delete(ctx, "/customers/"+customerID)
	return err
}
