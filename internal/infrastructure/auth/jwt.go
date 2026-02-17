package auth

import (
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken is returned when the token is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken is returned when the token has expired
	ErrExpiredToken = errors.New("token has expired")
	// ErrInvalidClaims is returned when the token claims are invalid
	ErrInvalidClaims = errors.New("invalid token claims")
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token operations
type JWTManager struct {
	secretKey     []byte
	tokenDuration time.Duration
	blacklist     sync.Map // token string -> expiry time
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string, tokenDurationHours int) *JWTManager {
	return &JWTManager{
		secretKey:     []byte(secretKey),
		tokenDuration: time.Duration(tokenDurationHours) * time.Hour,
	}
}

// GenerateToken generates a new JWT token for a user
func (m *JWTManager) GenerateToken(userID, email, role string) (string, error) {
	now := time.Now()
	expirationTime := now.Add(m.tokenDuration)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "condotrack-api",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// RefreshToken generates a new token from an existing valid token
func (m *JWTManager) RefreshToken(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	return m.GenerateToken(claims.UserID, claims.Email, claims.Role)
}

// GetTokenDuration returns the token duration
func (m *JWTManager) GetTokenDuration() time.Duration {
	return m.tokenDuration
}

// BlacklistToken adds a token to the blacklist until its expiry time.
func (m *JWTManager) BlacklistToken(tokenString string, expiry time.Time) {
	m.blacklist.Store(tokenString, expiry)
}

// IsBlacklisted checks whether a token has been blacklisted.
func (m *JWTManager) IsBlacklisted(tokenString string) bool {
	val, ok := m.blacklist.Load(tokenString)
	if !ok {
		return false
	}
	expiry := val.(time.Time)
	if time.Now().After(expiry) {
		m.blacklist.Delete(tokenString)
		return false
	}
	return true
}

// StartBlacklistCleanup starts a background goroutine that periodically
// removes expired entries from the token blacklist.
func (m *JWTManager) StartBlacklistCleanup(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			now := time.Now()
			m.blacklist.Range(func(key, value interface{}) bool {
				if now.After(value.(time.Time)) {
					m.blacklist.Delete(key)
				}
				return true
			})
		}
	}()
}
