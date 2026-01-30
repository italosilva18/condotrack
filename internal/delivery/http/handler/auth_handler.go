package handler

import (
	"strings"

	"github.com/condotrack/api/internal/delivery/http/middleware"
	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/usecase/auth"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	usecase auth.UseCase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(uc auth.UseCase) *AuthHandler {
	return &AuthHandler{usecase: uc}
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
		case auth.ErrInvalidCredentials:
			response.Unauthorized(c, "Invalid email or password")
		case auth.ErrUserInactive:
			response.Forbidden(c, "User account is inactive")
		default:
			response.InternalError(c, "Login failed: "+err.Error())
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
		response.InternalError(c, "Registration failed: "+err.Error())
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
		response.InternalError(c, "Failed to get user: "+err.Error())
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
			response.InternalError(c, "Failed to change password: "+err.Error())
		}
		return
	}

	response.SuccessWithMessage(c, "Password changed successfully", nil)
}

// Logout handles POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT tokens are stateless, so logout is handled client-side
	// by removing the token from storage. This endpoint is provided
	// for API completeness and can be used to log the logout event.
	response.SuccessWithMessage(c, "Logged out successfully", nil)
}

// ListUsers handles GET /api/v1/auth/users (admin only)
func (h *AuthHandler) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := h.usecase.ListUsers(ctx)
	if err != nil {
		response.InternalError(c, "Failed to fetch users: "+err.Error())
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
		response.InternalError(c, "Failed to fetch user: "+err.Error())
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
		response.InternalError(c, "Failed to update user: "+err.Error())
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
		response.InternalError(c, "Failed to update user: "+err.Error())
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
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, "Failed to delete user: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "User deleted successfully", nil)
}
