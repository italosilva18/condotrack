package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/condotrack/api/internal/config"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ImageHandler handles image upload and management
type ImageHandler struct {
	cfg *config.Config
}

// NewImageHandler creates a new image handler
func NewImageHandler(cfg *config.Config) *ImageHandler {
	return &ImageHandler{cfg: cfg}
}

// ListImages handles GET /api/v1/images
func (h *ImageHandler) ListImages(c *gin.Context) {
	uploadDir := h.cfg.UploadDir

	// Ensure directory exists
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		response.Success(c, []string{})
		return
	}

	files, err := os.ReadDir(uploadDir)
	if err != nil {
		response.SafeInternalError(c, "Failed to list images", err)
		return
	}

	var images []gin.H
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Filter only image files
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if !isImageExtension(ext) {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		images = append(images, gin.H{
			"filename":   file.Name(),
			"url":        fmt.Sprintf("/uploads/%s", file.Name()),
			"size":       info.Size(),
			"created_at": info.ModTime().Format(time.RFC3339),
		})
	}

	response.Success(c, images)
}

// UploadImage handles POST /api/v1/images
func (h *ImageHandler) UploadImage(c *gin.Context) {
	uploadDir := h.cfg.UploadDir

	// Ensure directory exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.SafeInternalError(c, "Failed to create upload directory", err)
		return
	}

	// Get file from request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "No file uploaded: "+err.Error())
		return
	}
	defer file.Close()

	// Check file size
	if header.Size > h.cfg.MaxUploadSize {
		response.BadRequest(c, fmt.Sprintf("File size exceeds maximum allowed size of %d bytes", h.cfg.MaxUploadSize))
		return
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !isImageExtension(ext) {
		response.BadRequest(c, "Invalid file type. Allowed types: jpg, jpeg, png, gif, webp")
		return
	}

	// Generate unique filename
	newFilename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)
	filePath := filepath.Join(uploadDir, newFilename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		response.SafeInternalError(c, "Failed to create file", err)
		return
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		response.SafeInternalError(c, "Failed to save file", err)
		return
	}

	response.Created(c, gin.H{
		"filename":     newFilename,
		"url":          fmt.Sprintf("/uploads/%s", newFilename),
		"original_name": header.Filename,
		"size":         header.Size,
	})
}

// DeleteImage handles DELETE /api/v1/images/:filename
func (h *ImageHandler) DeleteImage(c *gin.Context) {
	filename := c.Param("filename")

	// Validate filename (prevent path traversal)
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		response.BadRequest(c, "Invalid filename")
		return
	}

	filePath := filepath.Join(h.cfg.UploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		response.NotFound(c, "Image not found")
		return
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		response.SafeInternalError(c, "Failed to delete image", err)
		return
	}

	response.Success(c, gin.H{
		"message": "Image deleted successfully",
	})
}

// isImageExtension checks if the extension is a valid image extension
func isImageExtension(ext string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, valid := range validExtensions {
		if ext == valid {
			return true
		}
	}
	return false
}
