package auth

import (
	"testing"
)

func TestHashPassword_Success(t *testing.T) {
	hash, err := HashPassword("mypassword123")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == "" {
		t.Error("HashPassword() returned empty hash")
	}
}

func TestHashPassword_DifferentFromOriginal(t *testing.T) {
	password := "mypassword123"
	hash, _ := HashPassword(password)
	if hash == password {
		t.Error("hash should differ from original password")
	}
}

func TestHashPassword_DifferentHashesForSamePassword(t *testing.T) {
	hash1, _ := HashPassword("samepassword")
	hash2, _ := HashPassword("samepassword")
	if hash1 == hash2 {
		t.Error("bcrypt should generate different hashes for same password (different salts)")
	}
}

func TestCheckPassword_Correct(t *testing.T) {
	password := "correctpassword"
	hash, _ := HashPassword(password)
	if !CheckPassword(password, hash) {
		t.Error("CheckPassword() should return true for correct password")
	}
}

func TestCheckPassword_Wrong(t *testing.T) {
	hash, _ := HashPassword("originalpassword")
	if CheckPassword("wrongpassword", hash) {
		t.Error("CheckPassword() should return false for wrong password")
	}
}

func TestNewPasswordHasher_ValidCost(t *testing.T) {
	h := NewPasswordHasher(12)
	if h == nil {
		t.Fatal("NewPasswordHasher() returned nil")
	}
	if h.cost != 12 {
		t.Errorf("cost = %d, want 12", h.cost)
	}
}

func TestNewPasswordHasher_InvalidCostUsesDefault(t *testing.T) {
	// bcrypt.MinCost = 4, bcrypt.MaxCost = 31
	h := NewPasswordHasher(0)
	if h.cost != DefaultBcryptCost {
		t.Errorf("cost with invalid value = %d, want %d", h.cost, DefaultBcryptCost)
	}

	h2 := NewPasswordHasher(100)
	if h2.cost != DefaultBcryptCost {
		t.Errorf("cost with too-high value = %d, want %d", h2.cost, DefaultBcryptCost)
	}
}

func TestPasswordHasher_HashAndVerify(t *testing.T) {
	h := NewPasswordHasher(10)
	password := "hasherpassword"

	hash, err := h.Hash(password)
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	if hash == "" {
		t.Error("Hash() returned empty hash")
	}

	if !h.Verify(password, hash) {
		t.Error("Verify() should return true for correct password")
	}
	if h.Verify("wrongpassword", hash) {
		t.Error("Verify() should return false for wrong password")
	}
}
