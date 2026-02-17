package mercadopago

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	prodBaseURL    = "https://api.mercadopago.com"
	sandboxBaseURL = "https://api.mercadopago.com"
)

// Client represents the Mercado Pago API client.
type Client struct {
	accessToken string
	baseURL     string
	httpClient  *http.Client
}

// NewClient creates a new Mercado Pago API client.
func NewClient(accessToken, env string) *Client {
	baseURL := sandboxBaseURL
	if env == "production" {
		baseURL = prodBaseURL
	}

	return &Client{
		accessToken: accessToken,
		baseURL:     baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs an HTTP request to the Mercado Pago API.
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("X-Idempotency-Key", fmt.Sprintf("%d", time.Now().UnixNano()))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr MPAPIError
		if err := json.Unmarshal(respBody, &apiErr); err == nil && apiErr.Message != "" {
			return nil, &apiErr
		}
		return nil, fmt.Errorf("MP API error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// get performs a GET request.
func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	return c.doRequest(ctx, http.MethodGet, path, nil)
}

// post performs a POST request.
func (c *Client) post(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.doRequest(ctx, http.MethodPost, path, body)
}

// put performs a PUT request.
func (c *Client) put(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.doRequest(ctx, http.MethodPut, path, body)
}

// MPAPIError represents a Mercado Pago API error response.
type MPAPIError struct {
	Message string `json:"message"`
	ErrorCode string `json:"error"`
	Status  int    `json:"status"`
	Cause   []struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"cause,omitempty"`
}

// Error returns the error message.
func (e *MPAPIError) Error() string {
	if len(e.Cause) > 0 {
		return fmt.Sprintf("MP API error: %s - %s", e.Cause[0].Code, e.Cause[0].Description)
	}
	return fmt.Sprintf("MP API error: %s", e.Message)
}
