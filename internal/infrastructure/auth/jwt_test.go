package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func newTestJWTManager() *JWTManager {
	return NewJWTManager("test-secret-key-for-testing", 24)
}

func TestNewJWTManager(t *testing.T) {
	m := NewJWTManager("my-secret", 48)
	if m == nil {
		t.Fatal("NewJWTManager returned nil")
	}
	if string(m.secretKey) != "my-secret" {
		t.Errorf("secretKey = %q, want %q", string(m.secretKey), "my-secret")
	}
	if m.tokenDuration != 48*time.Hour {
		t.Errorf("tokenDuration = %v, want %v", m.tokenDuration, 48*time.Hour)
	}
}

func TestGenerateToken_Success(t *testing.T) {
	m := newTestJWTManager()
	token, err := m.GenerateToken("user-1", "test@example.com", "admin")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if token == "" {
		t.Error("GenerateToken() returned empty token")
	}
}

func TestGenerateToken_ClaimsCorrect(t *testing.T) {
	m := newTestJWTManager()
	tokenStr, err := m.GenerateToken("user-2", "claims@example.com", "student")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := m.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != "user-2" {
		t.Errorf("claims.UserID = %q, want %q", claims.UserID, "user-2")
	}
	if claims.Email != "claims@example.com" {
		t.Errorf("claims.Email = %q, want %q", claims.Email, "claims@example.com")
	}
	if claims.Role != "student" {
		t.Errorf("claims.Role = %q, want %q", claims.Role, "student")
	}
	if claims.Issuer != "condotrack-api" {
		t.Errorf("claims.Issuer = %q, want %q", claims.Issuer, "condotrack-api")
	}
	if claims.Subject != "user-2" {
		t.Errorf("claims.Subject = %q, want %q", claims.Subject, "user-2")
	}
}

func TestValidateToken_Valid(t *testing.T) {
	m := newTestJWTManager()
	tokenStr, _ := m.GenerateToken("user-3", "valid@example.com", "admin")

	claims, err := m.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}
	if claims.UserID != "user-3" {
		t.Errorf("claims.UserID = %q, want %q", claims.UserID, "user-3")
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	m := newTestJWTManager()
	_, err := m.ValidateToken("invalid.token.string")
	if err == nil {
		t.Error("ValidateToken() with invalid token should return error")
	}
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	// Create a manager with 0-hour duration (expires immediately)
	m := NewJWTManager("test-secret", 0)

	// Manually create an expired token
	claims := &Claims{
		UserID: "user-4",
		Email:  "expired@example.com",
		Role:   "student",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "condotrack-api",
			Subject:   "user-4",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("test-secret"))

	_, err := m.ValidateToken(tokenStr)
	if err == nil {
		t.Error("ValidateToken() with expired token should return error")
	}
	if err != ErrExpiredToken {
		t.Errorf("expected ErrExpiredToken, got %v", err)
	}
}

func TestValidateToken_WrongSigningMethod(t *testing.T) {
	m := newTestJWTManager()

	// Create a token with a different signing method (none)
	claims := &Claims{
		UserID: "user-5",
		Email:  "wrong@example.com",
		Role:   "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenStr, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	_, err := m.ValidateToken(tokenStr)
	if err == nil {
		t.Error("ValidateToken() with wrong signing method should return error")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	m1 := NewJWTManager("secret-one", 24)
	m2 := NewJWTManager("secret-two", 24)

	tokenStr, _ := m1.GenerateToken("user-6", "wrong@example.com", "admin")
	_, err := m2.ValidateToken(tokenStr)
	if err == nil {
		t.Error("ValidateToken() with wrong secret should return error")
	}
}

func TestRefreshToken_Valid(t *testing.T) {
	m := newTestJWTManager()
	originalToken, _ := m.GenerateToken("user-7", "refresh@example.com", "instructor")

	newToken, err := m.RefreshToken(originalToken)
	if err != nil {
		t.Fatalf("RefreshToken() error = %v", err)
	}
	if newToken == "" {
		t.Error("RefreshToken() returned empty token")
	}

	// Validate the new token has same claims
	claims, err := m.ValidateToken(newToken)
	if err != nil {
		t.Fatalf("ValidateToken(refreshed) error = %v", err)
	}
	if claims.UserID != "user-7" {
		t.Errorf("refreshed claims.UserID = %q, want %q", claims.UserID, "user-7")
	}
	if claims.Email != "refresh@example.com" {
		t.Errorf("refreshed claims.Email = %q, want %q", claims.Email, "refresh@example.com")
	}
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	m := newTestJWTManager()
	_, err := m.RefreshToken("invalid-token")
	if err == nil {
		t.Error("RefreshToken() with invalid token should return error")
	}
}

func TestGetTokenDuration(t *testing.T) {
	m := NewJWTManager("secret", 12)
	if m.GetTokenDuration() != 12*time.Hour {
		t.Errorf("GetTokenDuration() = %v, want %v", m.GetTokenDuration(), 12*time.Hour)
	}
}

func TestBlacklistToken_IsBlacklisted(t *testing.T) {
	m := newTestJWTManager()
	token, _ := m.GenerateToken("user-8", "bl@example.com", "admin")

	// Before blacklisting
	if m.IsBlacklisted(token) {
		t.Error("token should not be blacklisted before BlacklistToken()")
	}

	// After blacklisting
	m.BlacklistToken(token, time.Now().Add(1*time.Hour))
	if !m.IsBlacklisted(token) {
		t.Error("token should be blacklisted after BlacklistToken()")
	}
}

func TestBlacklistToken_NotBlacklisted(t *testing.T) {
	m := newTestJWTManager()
	if m.IsBlacklisted("random-token-never-blacklisted") {
		t.Error("random token should not be blacklisted")
	}
}

func TestBlacklistToken_ExpiredEntryRemoved(t *testing.T) {
	m := newTestJWTManager()
	token := "expired-blacklist-token"

	// Blacklist with expiry in the past
	m.BlacklistToken(token, time.Now().Add(-1*time.Second))

	// Should not be blacklisted (expired entry gets cleaned up on check)
	if m.IsBlacklisted(token) {
		t.Error("expired blacklist entry should return false")
	}
}
