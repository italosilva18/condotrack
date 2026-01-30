package asaas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents the Asaas API client
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Asaas API client
func NewClient(apiKey, baseURL string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs an HTTP request to the Asaas API
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
	req.Header.Set("access_token", c.apiKey)

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
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err == nil && len(apiErr.Errors) > 0 {
			return nil, &apiErr
		}
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// get performs a GET request
func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	return c.doRequest(ctx, http.MethodGet, path, nil)
}

// post performs a POST request
func (c *Client) post(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.doRequest(ctx, http.MethodPost, path, body)
}

// put performs a PUT request
func (c *Client) put(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.doRequest(ctx, http.MethodPut, path, body)
}

// delete performs a DELETE request
func (c *Client) delete(ctx context.Context, path string) ([]byte, error) {
	return c.doRequest(ctx, http.MethodDelete, path, nil)
}
