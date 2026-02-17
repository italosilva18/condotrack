package handler

import (
	"net/http"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/condotrack/api/internal/usecase/setting"
	"github.com/gin-gonic/gin"
)

// SettingHandler handles HTTP requests for settings
type SettingHandler struct {
	usecase *setting.UseCase
}

// NewSettingHandler creates a new SettingHandler
func NewSettingHandler(usecase *setting.UseCase) *SettingHandler {
	return &SettingHandler{usecase: usecase}
}

// ListSettings returns all settings grouped by category
// GET /api/v1/settings
func (h *SettingHandler) ListSettings(c *gin.Context) {
	settings, err := h.usecase.GetSettingsByCategory(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
	})
}

// GetAllSettings returns all settings as a flat list
// GET /api/v1/settings/all
func (h *SettingHandler) GetAllSettings(c *gin.Context) {
	settings, err := h.usecase.GetAllSettings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
	})
}

// GetSettingByKey returns a specific setting by key
// GET /api/v1/settings/:key
func (h *SettingHandler) GetSettingByKey(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "setting key is required",
		})
		return
	}

	setting, err := h.usecase.GetSettingByKey(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    setting,
	})
}

// UpdateSetting updates a single setting by key
// PUT /api/v1/settings/:key
func (h *SettingHandler) UpdateSetting(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "setting key is required",
		})
		return
	}

	var req entity.UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := h.usecase.UpdateSetting(c.Request.Context(), key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Setting updated successfully",
	})
}

// BulkUpdateSettings updates multiple settings at once
// PUT /api/v1/settings
func (h *SettingHandler) BulkUpdateSettings(c *gin.Context) {
	var req entity.BulkUpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if len(req.Settings) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "no settings provided",
		})
		return
	}

	if err := h.usecase.BulkUpdateSettings(c.Request.Context(), req.Settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Settings updated successfully",
		"count":   len(req.Settings),
	})
}

// GetCategories returns all available setting categories
// GET /api/v1/settings/categories
func (h *SettingHandler) GetCategories(c *gin.Context) {
	categories, err := h.usecase.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Add labels to categories
	result := make([]map[string]string, len(categories))
	for i, cat := range categories {
		label, ok := entity.CategoryLabels[entity.SettingCategory(cat)]
		if !ok {
			label = cat
		}
		result[i] = map[string]string{
			"key":   cat,
			"label": label,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}
