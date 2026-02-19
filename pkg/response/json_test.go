package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func parseResponse(t *testing.T, w *httptest.ResponseRecorder) Response {
	t.Helper()
	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
	return resp
}

func TestSuccess(t *testing.T) {
	c, w := newTestContext()
	Success(c, map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	resp := parseResponse(t, w)
	if !resp.Success {
		t.Error("Success should be true")
	}
	if resp.Error != "" {
		t.Errorf("Error should be empty, got %q", resp.Error)
	}
}

func TestCreated(t *testing.T) {
	c, w := newTestContext()
	Created(c, map[string]string{"id": "new-1"})

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}
	resp := parseResponse(t, w)
	if !resp.Success {
		t.Error("Success should be true")
	}
}

func TestBadRequest(t *testing.T) {
	c, w := newTestContext()
	BadRequest(c, "invalid input")

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
	resp := parseResponse(t, w)
	if resp.Success {
		t.Error("Success should be false")
	}
	if resp.Error != "invalid input" {
		t.Errorf("Error = %q, want %q", resp.Error, "invalid input")
	}
}

func TestNotFound(t *testing.T) {
	c, w := newTestContext()
	NotFound(c, "resource not found")

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
	resp := parseResponse(t, w)
	if resp.Success {
		t.Error("Success should be false")
	}
	if resp.Error != "resource not found" {
		t.Errorf("Error = %q, want %q", resp.Error, "resource not found")
	}
}

func TestInternalError(t *testing.T) {
	c, w := newTestContext()
	InternalError(c, "something went wrong")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
	resp := parseResponse(t, w)
	if resp.Success {
		t.Error("Success should be false")
	}
	if resp.Error != "something went wrong" {
		t.Errorf("Error = %q, want %q", resp.Error, "something went wrong")
	}
}

func TestSafeInternalError(t *testing.T) {
	c, w := newTestContext()
	SafeInternalError(c, "DB query failed", errors.New("connection refused"))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
	resp := parseResponse(t, w)
	if resp.Success {
		t.Error("Success should be false")
	}
	// Should NOT expose the real error
	if resp.Error != "Internal server error" {
		t.Errorf("Error = %q, want generic message %q", resp.Error, "Internal server error")
	}
}

func TestUnauthorized(t *testing.T) {
	c, w := newTestContext()
	Unauthorized(c, "invalid token")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
	resp := parseResponse(t, w)
	if resp.Success {
		t.Error("Success should be false")
	}
	if resp.Error != "invalid token" {
		t.Errorf("Error = %q, want %q", resp.Error, "invalid token")
	}
}

func TestForbidden(t *testing.T) {
	c, w := newTestContext()
	Forbidden(c, "access denied")

	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
	resp := parseResponse(t, w)
	if resp.Success {
		t.Error("Success should be false")
	}
	if resp.Error != "access denied" {
		t.Errorf("Error = %q, want %q", resp.Error, "access denied")
	}
}

func TestValidationError(t *testing.T) {
	c, w := newTestContext()
	ValidationError(c, "email is required")

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnprocessableEntity)
	}
	resp := parseResponse(t, w)
	if resp.Success {
		t.Error("Success should be false")
	}
	if resp.Error != "email is required" {
		t.Errorf("Error = %q, want %q", resp.Error, "email is required")
	}
}

func TestSuccessWithMessage(t *testing.T) {
	c, w := newTestContext()
	SuccessWithMessage(c, "operation completed", map[string]string{"result": "ok"})

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	resp := parseResponse(t, w)
	if !resp.Success {
		t.Error("Success should be true")
	}
	if resp.Message != "operation completed" {
		t.Errorf("Message = %q, want %q", resp.Message, "operation completed")
	}
}
