package entity

import (
	"testing"
	"time"
)

func TestUserRole_String(t *testing.T) {
	tests := []struct {
		role     UserRole
		expected string
	}{
		{RoleAdmin, "admin"},
		{RoleManager, "manager"},
		{RoleInstructor, "instructor"},
		{RoleStudent, "student"},
		{RoleUser, "user"},
	}

	for _, tt := range tests {
		if got := tt.role.String(); got != tt.expected {
			t.Errorf("UserRole(%q).String() = %q, want %q", tt.role, got, tt.expected)
		}
	}
}

func TestUserRole_IsValid_ValidRoles(t *testing.T) {
	validRoles := []UserRole{RoleAdmin, RoleManager, RoleInstructor, RoleStudent, RoleUser}
	for _, role := range validRoles {
		if !role.IsValid() {
			t.Errorf("UserRole(%q).IsValid() = false, want true", role)
		}
	}
}

func TestUserRole_IsValid_InvalidRoles(t *testing.T) {
	invalidRoles := []UserRole{"superadmin", "guest", "", "ADMIN", "Admin"}
	for _, role := range invalidRoles {
		if role.IsValid() {
			t.Errorf("UserRole(%q).IsValid() = true, want false", role)
		}
	}
}

func TestUser_ToPublic(t *testing.T) {
	phone := "11999999999"
	avatar := "https://example.com/avatar.jpg"
	now := time.Now()
	lastLogin := now.Add(-1 * time.Hour)

	user := &User{
		ID:           "user-123",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		Nome:         "Test User",
		Role:         RoleAdmin,
		IsActive:     true,
		Phone:        &phone,
		AvatarURL:    &avatar,
		LastLoginAt:  &lastLogin,
		CreatedAt:    now,
	}

	pub := user.ToPublic()

	if pub.ID != user.ID {
		t.Errorf("ToPublic().ID = %q, want %q", pub.ID, user.ID)
	}
	if pub.Email != user.Email {
		t.Errorf("ToPublic().Email = %q, want %q", pub.Email, user.Email)
	}
	if pub.Nome != user.Nome {
		t.Errorf("ToPublic().Nome = %q, want %q", pub.Nome, user.Nome)
	}
	if pub.Role != user.Role {
		t.Errorf("ToPublic().Role = %q, want %q", pub.Role, user.Role)
	}
	if pub.IsActive != user.IsActive {
		t.Errorf("ToPublic().IsActive = %v, want %v", pub.IsActive, user.IsActive)
	}
	if pub.Phone == nil || *pub.Phone != phone {
		t.Error("ToPublic().Phone should match original")
	}
	if pub.AvatarURL == nil || *pub.AvatarURL != avatar {
		t.Error("ToPublic().AvatarURL should match original")
	}
	if pub.LastLoginAt == nil || !pub.LastLoginAt.Equal(lastLogin) {
		t.Error("ToPublic().LastLoginAt should match original")
	}
	if !pub.CreatedAt.Equal(now) {
		t.Error("ToPublic().CreatedAt should match original")
	}
}

func TestUser_ToPublic_NilOptionalFields(t *testing.T) {
	user := &User{
		ID:       "user-456",
		Email:    "minimal@example.com",
		Nome:     "Minimal",
		Role:     RoleStudent,
		IsActive: true,
	}

	pub := user.ToPublic()

	if pub.Phone != nil {
		t.Error("ToPublic().Phone should be nil")
	}
	if pub.AvatarURL != nil {
		t.Error("ToPublic().AvatarURL should be nil")
	}
	if pub.LastLoginAt != nil {
		t.Error("ToPublic().LastLoginAt should be nil")
	}
}
