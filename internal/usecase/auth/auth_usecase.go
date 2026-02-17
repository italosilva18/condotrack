package auth

import (
	"context"
	"errors"
	"log"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/domain/repository"
	"github.com/condotrack/api/internal/infrastructure/auth"
	"github.com/google/uuid"
)

var (
	// ErrInvalidCredentials is returned when email or password is incorrect
	ErrInvalidCredentials = errors.New("invalid email or password")
	// ErrUserNotFound is returned when user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrUserInactive is returned when user account is inactive
	ErrUserInactive = errors.New("user account is inactive")
	// ErrEmailAlreadyExists is returned when email is already registered
	ErrEmailAlreadyExists = errors.New("email already registered")
	// ErrInvalidOldPassword is returned when the old password is incorrect
	ErrInvalidOldPassword = errors.New("invalid old password")
	// ErrSamePassword is returned when new password is same as old
	ErrSamePassword = errors.New("new password must be different from old password")
)

// UseCase defines the interface for authentication use cases
type UseCase interface {
	// Login authenticates a user and returns a token
	Login(ctx context.Context, email, password string) (*entity.LoginResponse, error)

	// Register creates a new user account
	Register(ctx context.Context, req entity.RegisterRequest) (*entity.User, error)

	// GetCurrentUser returns the current user by ID
	GetCurrentUser(ctx context.Context, userID string) (*entity.UserPublic, error)

	// ChangePassword changes the user's password
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error

	// GetUserByID returns a user by ID
	GetUserByID(ctx context.Context, id string) (*entity.User, error)

	// ListUsers returns all users
	ListUsers(ctx context.Context) ([]entity.UserPublic, error)

	// UpdateUser updates user information
	UpdateUser(ctx context.Context, userID string, req entity.UpdateUserRequest) (*entity.UserPublic, error)

	// AdminUpdateUser updates user information (admin only)
	AdminUpdateUser(ctx context.Context, userID string, req entity.AdminUpdateUserRequest) (*entity.UserPublic, error)

	// DeleteUser soft deletes a user
	DeleteUser(ctx context.Context, userID string) error
}

type authUseCase struct {
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager
}

// NewUseCase creates a new authentication use case
func NewUseCase(userRepo repository.UserRepository, jwtManager *auth.JWTManager) UseCase {
	return &authUseCase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Login authenticates a user and returns a token
func (uc *authUseCase) Login(ctx context.Context, email, password string) (*entity.LoginResponse, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Verify password
	if !auth.CheckPassword(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, err := uc.jwtManager.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	// Update last login
	if err := uc.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		log.Printf("[WARN] Failed to update last login for user %s (email: %s): %v", user.ID, user.Email, err)
	}

	return &entity.LoginResponse{
		Token:     token,
		ExpiresIn: int64(uc.jwtManager.GetTokenDuration().Seconds()),
		User:      user.ToPublic(),
	}, nil
}

// Register creates a new user account
func (uc *authUseCase) Register(ctx context.Context, req entity.RegisterRequest) (*entity.User, error) {
	// Check if email already exists
	exists, err := uc.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Force role to "student" on self-registration to prevent privilege escalation.
	// Only admins can assign elevated roles via AdminUpdateUser.
	role := entity.RoleStudent

	// Create user
	user := &entity.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: passwordHash,
		Nome:         req.Nome,
		Role:         role,
		IsActive:     true,
		Phone:        req.Phone,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetCurrentUser returns the current user by ID
func (uc *authUseCase) GetCurrentUser(ctx context.Context, userID string) (*entity.UserPublic, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user.ToPublic(), nil
}

// ChangePassword changes the user's password
func (uc *authUseCase) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	// Find user
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Verify old password
	if !auth.CheckPassword(oldPassword, user.PasswordHash) {
		return ErrInvalidOldPassword
	}

	// Check if new password is same as old
	if auth.CheckPassword(newPassword, user.PasswordHash) {
		return ErrSamePassword
	}

	// Hash new password
	newHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	return uc.userRepo.UpdatePassword(ctx, userID, newHash)
}

// GetUserByID returns a user by ID
func (uc *authUseCase) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	return uc.userRepo.FindByID(ctx, id)
}

// ListUsers returns all users
func (uc *authUseCase) ListUsers(ctx context.Context) ([]entity.UserPublic, error) {
	users, err := uc.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	publicUsers := make([]entity.UserPublic, len(users))
	for i, u := range users {
		publicUsers[i] = *u.ToPublic()
	}

	return publicUsers, nil
}

// UpdateUser updates user information
func (uc *authUseCase) UpdateUser(ctx context.Context, userID string, req entity.UpdateUserRequest) (*entity.UserPublic, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Update fields if provided
	if req.Nome != nil {
		user.Nome = *req.Nome
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user.ToPublic(), nil
}

// AdminUpdateUser updates user information (admin only)
func (uc *authUseCase) AdminUpdateUser(ctx context.Context, userID string, req entity.AdminUpdateUserRequest) (*entity.UserPublic, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Update fields if provided
	if req.Nome != nil {
		user.Nome = *req.Nome
	}
	if req.Email != nil {
		// Check if new email already exists
		if *req.Email != user.Email {
			exists, err := uc.userRepo.ExistsByEmail(ctx, *req.Email)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, ErrEmailAlreadyExists
			}
		}
		user.Email = *req.Email
	}
	if req.Role != nil && req.Role.IsValid() {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user.ToPublic(), nil
}

// DeleteUser soft deletes a user
func (uc *authUseCase) DeleteUser(ctx context.Context, userID string) error {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return uc.userRepo.Delete(ctx, userID)
}
