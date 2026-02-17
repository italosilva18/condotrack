package entity

import "time"

// UserRole represents the role of a user
type UserRole string

const (
	// RoleAdmin represents an administrator user
	RoleAdmin UserRole = "admin"
	// RoleManager represents a manager user
	RoleManager UserRole = "manager"
	// RoleInstructor represents an instructor user
	RoleInstructor UserRole = "instructor"
	// RoleStudent represents a student user
	RoleStudent UserRole = "student"
	// RoleUser represents a regular user
	RoleUser UserRole = "user"
)

// String returns the string representation of the role
func (r UserRole) String() string {
	return string(r)
}

// IsValid checks if the role is valid
func (r UserRole) IsValid() bool {
	switch r {
	case RoleAdmin, RoleManager, RoleInstructor, RoleStudent, RoleUser:
		return true
	}
	return false
}

// User represents a user entity in the system
type User struct {
	ID           string     `db:"id" json:"id"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	Nome         string     `db:"name" json:"nome"`
	Role         UserRole   `db:"role" json:"role"`
	IsActive     bool       `db:"is_active" json:"is_active"`
	Phone        *string    `db:"phone" json:"phone,omitempty"`
	CPF          *string    `db:"cpf" json:"cpf,omitempty"`
	AvatarURL    *string    `db:"avatar_url" json:"avatar_url,omitempty"`
	LastLoginAt  *time.Time `db:"last_login" json:"last_login,omitempty"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// UserPublic represents public user information (without sensitive data)
type UserPublic struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	Nome        string     `json:"nome"`
	Role        UserRole   `json:"role"`
	IsActive    bool       `json:"is_active"`
	Phone       *string    `json:"phone,omitempty"`
	AvatarURL   *string    `json:"avatar_url,omitempty"`
	LastLoginAt *time.Time `json:"last_login,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ToPublic converts a User to UserPublic (removes sensitive data)
func (u *User) ToPublic() *UserPublic {
	return &UserPublic{
		ID:          u.ID,
		Email:       u.Email,
		Nome:        u.Nome,
		Role:        u.Role,
		IsActive:    u.IsActive,
		Phone:       u.Phone,
		AvatarURL:   u.AvatarURL,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Nome     string   `json:"nome" binding:"required,min=2"`
	Role     UserRole `json:"role"`
	Phone    *string  `json:"phone,omitempty"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Token     string      `json:"token"`
	ExpiresIn int64       `json:"expires_in"`
	User      *UserPublic `json:"user"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UpdateUserRequest represents a user update request
type UpdateUserRequest struct {
	Nome      *string `json:"nome,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// AdminUpdateUserRequest represents an admin user update request
type AdminUpdateUserRequest struct {
	Nome      *string   `json:"nome,omitempty"`
	Email     *string   `json:"email,omitempty"`
	Role      *UserRole `json:"role,omitempty"`
	IsActive  *bool     `json:"is_active,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
}
