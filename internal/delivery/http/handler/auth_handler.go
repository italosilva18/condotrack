package handler

import (
	"strings"

	infraAuth "github.com/condotrack/api/internal/infrastructure/auth"

	"github.com/condotrack/api/internal/delivery/http/middleware"
	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/usecase/auth"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	usecase    auth.UseCase
	jwtManager *infraAuth.JWTManager
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(uc auth.UseCase, jwtManager *infraAuth.JWTManager) *AuthHandler {
	return &AuthHandler{usecase: uc, jwtManager: jwtManager}
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	result, err := h.usecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials, auth.ErrUserInactive:
			// Return same message for both to prevent user enumeration
			response.Unauthorized(c, "Invalid email or password")
		default:
			response.SafeInternalError(c, "Login failed", err)
		}
		return
	}

	response.Success(c, result)
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	user, err := h.usecase.Register(ctx, req)
	if err != nil {
		if err == auth.ErrEmailAlreadyExists {
			response.BadRequest(c, "Email already registered")
			return
		}
		response.SafeInternalError(c, "Registration failed", err)
		return
	}

	response.Created(c, user.ToPublic())
}

// GetCurrentUser handles GET /api/v1/auth/me
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	user, err := h.usecase.GetCurrentUser(ctx, userID)
	if err != nil {
		if err == auth.ErrUserNotFound {
			response.NotFound(c, "User not found")
			return
		}
		response.SafeInternalError(c, "Failed to get user", err)
		return
	}

	response.Success(c, user)
}

// ChangePassword handles POST /api/v1/auth/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req entity.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	err := h.usecase.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword)
	if err != nil {
		switch err {
		case auth.ErrUserNotFound:
			response.NotFound(c, "User not found")
		case auth.ErrInvalidOldPassword:
			response.BadRequest(c, "Invalid old password")
		case auth.ErrSamePassword:
			response.BadRequest(c, "New password must be different from old password")
		default:
			response.SafeInternalError(c, "Failed to change password", err)
		}
		return
	}

	response.SuccessWithMessage(c, "Password changed successfully", nil)
}

// Logout handles POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if claims, err := h.jwtManager.ValidateToken(tokenString); err == nil {
			h.jwtManager.BlacklistToken(tokenString, claims.ExpiresAt.Time)
		}
	}
	response.SuccessWithMessage(c, "Logged out successfully", nil)
}

// ListUsers handles GET /api/v1/auth/users (admin only)
func (h *AuthHandler) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := h.usecase.ListUsers(ctx)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch users", err)
		return
	}

	response.Success(c, users)
}

// GetUserByID handles GET /api/v1/auth/users/:id
func (h *AuthHandler) GetUserByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	user, err := h.usecase.GetUserByID(ctx, id)
	if err != nil {
		response.SafeInternalError(c, "Failed to fetch user", err)
		return
	}

	if user == nil {
		response.NotFound(c, "User not found")
		return
	}

	response.Success(c, user.ToPublic())
}

// UpdateUser handles PUT /api/v1/auth/me
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req entity.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	user, err := h.usecase.UpdateUser(ctx, userID, req)
	if err != nil {
		if err == auth.ErrUserNotFound {
			response.NotFound(c, "User not found")
			return
		}
		response.SafeInternalError(c, "Failed to update user", err)
		return
	}

	response.Success(c, user)
}

// AdminUpdateUser handles PUT /api/v1/auth/users/:id (admin only)
func (h *AuthHandler) AdminUpdateUser(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req entity.AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	user, err := h.usecase.AdminUpdateUser(ctx, id, req)
	if err != nil {
		if err == auth.ErrUserNotFound {
			response.NotFound(c, "User not found")
			return
		}
		if err == auth.ErrEmailAlreadyExists {
			response.BadRequest(c, "Email already in use")
			return
		}
		response.SafeInternalError(c, "Failed to update user", err)
		return
	}

	response.Success(c, user)
}

// DeleteUser handles DELETE /api/v1/auth/users/:id (admin only)
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	// Prevent self-deletion
	currentUserID, ok := middleware.GetUserID(c)
	if ok && currentUserID == id {
		response.BadRequest(c, "Cannot delete your own account")
		return
	}

	err := h.usecase.DeleteUser(ctx, id)
	if err != nil {
		if err == auth.ErrUserNotFound {
			response.NotFound(c, "User not found")
			return
		}
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "User not found")
			return
		}
		response.SafeInternalError(c, "Failed to delete user", err)
		return
	}

	response.SuccessWithMessage(c, "User deleted successfully", nil)
}
