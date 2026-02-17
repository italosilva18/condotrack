package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/condotrack/api/internal/config"
	"github.com/condotrack/api/internal/infrastructure/storage"
	"github.com/condotrack/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PortalHandler handles portal-specific HTTP requests
type PortalHandler struct {
	storage *storage.StorageService
	cfg     *config.Config
}

// NewPortalHandler creates a new portal handler
func NewPortalHandler(storage *storage.StorageService, cfg *config.Config) *PortalHandler {
	return &PortalHandler{
		storage: storage,
		cfg:     cfg,
	}
}

// System image IDs that map to specific filenames
var systemImageMap = map[string]string{
	"sys_logo":      "logo-condotrack.png",
	"sys_prof":      "professor-sergio.png",
	"sys_prof_mini": "prof-sergio.png",
	"pos_1":         "step1.png",
	"pos_2":         "step2.png",
	"pos_3":         "step3.png",
	"pos_4":         "step4.png",
	"pos_5":         "step5.png",
	"pos_6":         "step6.png",
	"pos_7":         "step7.png",
	"padlet_bg":     "padlet_bg.png",
	"padlet_cover":  "padlet-cover-v2.png",
	"padlet_comm":   "padlet-comunidade.png",
}

// PortalImageResponse represents the response for portal images
type PortalImageResponse struct {
	Name             string   `json:"name"`
	URL              string   `json:"url"`
	Size             int64    `json:"size"`
	LastModified     int64    `json:"mtime"`
	Tags             []string `json:"tags"`
	IsSystemOverwrite bool    `json:"isSystemOverwrite"`
}

// UploadPortalImageRequest represents the request to upload a portal image
type UploadPortalImageRequest struct {
	Image string `json:"image" binding:"required"` // Base64 encoded image
	Name  string `json:"name" binding:"required"`  // Original filename
	ID    string `json:"id,omitempty"`             // System ID for overwrite (e.g., "sys_logo", "pos_1")
}

// checkStorage verifies storage is initialized, returning false and writing an error if not.
func (h *PortalHandler) checkStorage(c *gin.Context) bool {
	if h.storage == nil {
		response.InternalError(c, "Storage service is not available")
		return false
	}
	return true
}

// ListPortalImages handles GET /api/v1/portal/images
func (h *PortalHandler) ListPortalImages(c *gin.Context) {
	if !h.checkStorage(c) {
		return
	}
	ctx := c.Request.Context()

	files, err := h.storage.ListFiles(ctx, h.cfg.MinioBucketPortal, "")
	if err != nil {
		response.SafeInternalError(c, "Failed to list images", err)
		return
	}

	// Build response with portal-specific metadata
	var images []PortalImageResponse
	for _, f := range files {
		tags := []string{"server"}
		isSystem := false

		// Check if this is a system image
		for id, filename := range systemImageMap {
			if f.Name == filename {
				tags = append(tags, "system", id)
				isSystem = true
				break
			}
		}

		images = append(images, PortalImageResponse{
			Name:             f.Name,
			URL:              f.URL,
			Size:             f.Size,
			LastModified:     f.LastModified.Unix(),
			Tags:             tags,
			IsSystemOverwrite: isSystem,
		})
	}

	response.Success(c, images)
}

// UploadPortalImage handles POST /api/v1/portal/images
func (h *PortalHandler) UploadPortalImage(c *gin.Context) {
	if !h.checkStorage(c) {
		return
	}
	ctx := c.Request.Context()

	var req UploadPortalImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Validate base64 data
	if !strings.HasPrefix(req.Image, "data:image/") {
		response.BadRequest(c, "Invalid image format. Must be base64 encoded image with data URI.")
		return
	}

	// Validate base64 image size (rough estimate: base64 is ~4/3 of original)
	// Limit to ~10MB decoded (which is ~13.3MB base64)
	if len(req.Image) > 14*1024*1024 {
		response.BadRequest(c, "Image too large. Maximum size is approximately 10MB.")
		return
	}

	// Determine target filename
	targetFilename := req.Name
	isSystemOverwrite := false

	if req.ID != "" {
		// Check if it's a system ID
		if sysFilename, ok := systemImageMap[req.ID]; ok {
			targetFilename = sysFilename
			isSystemOverwrite = true
		}
	}

	// If not a system overwrite, generate unique filename
	if !isSystemOverwrite {
		ext := filepath.Ext(req.Name)
		if ext == "" {
			ext = ".png"
		}
		targetFilename = fmt.Sprintf("upload_%d_%s%s", time.Now().Unix(), uuid.New().String()[:8], ext)
	}

	// Upload to MinIO
	result, err := h.storage.UploadBase64(ctx, h.cfg.MinioBucketPortal, targetFilename, req.Image)
	if err != nil {
		response.SafeInternalError(c, "Failed to upload image", err)
		return
	}

	// Add cache buster to URL
	urlWithCache := fmt.Sprintf("%s?t=%d", result.URL, time.Now().Unix())

	response.Created(c, gin.H{
		"success": true,
		"url":     urlWithCache,
		"name":    result.Filename,
		"path":    result.URL,
		"message": "Imagem salva e sincronizada no servidor",
	})
}

// DeletePortalImage handles DELETE /api/v1/portal/images/:filename
func (h *PortalHandler) DeletePortalImage(c *gin.Context) {
	if !h.checkStorage(c) {
		return
	}
	ctx := c.Request.Context()
	filename := c.Param("filename")

	if filename == "" {
		response.BadRequest(c, "Filename is required")
		return
	}

	// Prevent path traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		response.BadRequest(c, "Invalid filename")
		return
	}

	// Check if it's a system image (cannot delete system images)
	for _, sysFilename := range systemImageMap {
		if filename == sysFilename {
			response.BadRequest(c, "Cannot delete system images. You can only overwrite them.")
			return
		}
	}

	// Check if it's a user upload (starts with "upload_")
	if !strings.HasPrefix(filename, "upload_") {
		response.BadRequest(c, "Cannot delete this file. Only user uploads can be deleted.")
		return
	}

	// Delete from MinIO
	err := h.storage.DeleteFile(ctx, h.cfg.MinioBucketPortal, filename)
	if err != nil {
		response.SafeInternalError(c, "Failed to delete image", err)
		return
	}

	response.Success(c, gin.H{
		"success": true,
		"message": "Arquivo removido",
	})
}

// AIProxyRequest represents the request to proxy to Gemini AI
type AIProxyRequest struct {
	Message string `json:"message" binding:"required"`
	Context string `json:"context,omitempty"`
}

// AIProxyResponse represents the response from Gemini AI
type AIProxyResponse struct {
	Success  bool   `json:"success"`
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}

// ProxyGeminiAI handles POST /api/v1/portal/ai
func (h *PortalHandler) ProxyGeminiAI(c *gin.Context) {
	var req AIProxyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Limit input length to prevent payload amplification
	if len(req.Message) > 5000 {
		response.BadRequest(c, "Message too long. Maximum 5000 characters.")
		return
	}
	if len(req.Context) > 10000 {
		response.BadRequest(c, "Context too long. Maximum 10000 characters.")
		return
	}

	if h.cfg.GeminiAPIKey == "" {
		response.InternalError(c, "AI service not available")
		return
	}

	// Build Gemini request
	geminiReq := map[string]interface{}{
		"systemInstruction": map[string]interface{}{
			"parts": []map[string]string{
				{"text": "CRÍTICO: TODAS as respostas DEVEM ser em PORTUGUÊS DO BRASIL, independente do idioma da pergunta. Você é um assistente especializado em gestão de condomínios, contratos de facilities, auditorias e ISO 9001."},
			},
		},
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": req.Message},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.7,
			"maxOutputTokens": 1024,
		},
	}

	// Add context if provided
	if req.Context != "" {
		contents := geminiReq["contents"].([]map[string]interface{})
		contents[0]["parts"] = []map[string]string{
			{"text": fmt.Sprintf("Contexto: %s\n\nPergunta: %s", req.Context, req.Message)},
		}
		geminiReq["contents"] = contents
	}

	jsonData, err := json.Marshal(geminiReq)
	if err != nil {
		response.InternalError(c, "Failed to build request")
		return
	}

	// Call Gemini API (API key sent via header, not URL, to avoid logging exposure)
	apiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent"

	httpReq, err := http.NewRequestWithContext(c.Request.Context(), "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		response.InternalError(c, "Failed to create request")
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-goog-api-key", h.cfg.GeminiAPIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		response.SafeInternalError(c, "Failed to call AI service", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.InternalError(c, "Failed to read response")
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] Gemini API error (status %d): %s", resp.StatusCode, string(body))
		response.InternalError(c, "AI service temporarily unavailable")
		return
	}

	// Parse Gemini response
	var geminiResp map[string]interface{}
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		response.InternalError(c, "Failed to parse Gemini response")
		return
	}

	// Extract text from response
	text := ""
	if candidates, ok := geminiResp["candidates"].([]interface{}); ok && len(candidates) > 0 {
		if candidate, ok := candidates[0].(map[string]interface{}); ok {
			if content, ok := candidate["content"].(map[string]interface{}); ok {
				if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
					if part, ok := parts[0].(map[string]interface{}); ok {
						if t, ok := part["text"].(string); ok {
							text = t
						}
					}
				}
			}
		}
	}

	if text == "" {
		response.InternalError(c, "Empty response from Gemini")
		return
	}

	response.Success(c, AIProxyResponse{
		Success:  true,
		Response: text,
	})
}

// UploadEvidence handles POST /api/v1/portal/evidence - Upload de arquivos de evidência para auditorias
func (h *PortalHandler) UploadEvidence(c *gin.Context) {
	if !h.checkStorage(c) {
		return
	}
	ctx := c.Request.Context()

	// Parse multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "No file provided: "+err.Error())
		return
	}
	defer file.Close()

	// Validate file type
	contentType := storage.GetContentTypeFromExtension(header.Filename)
	if !storage.IsAllowedEvidenceType(contentType) {
		response.BadRequest(c, "File type not allowed. Allowed: jpg, jpeg, png, gif, pdf, doc, docx")
		return
	}

	// Validate file size (max 10MB for evidence)
	if header.Size > 10*1024*1024 {
		response.BadRequest(c, "File too large. Maximum size is 10MB")
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("ev_%d_%s%s", time.Now().Unix(), uuid.New().String()[:8], ext)

	// Upload to MinIO
	result, err := h.storage.UploadFile(ctx, h.cfg.MinioBucketEvidence, filename, file, header.Size, contentType)
	if err != nil {
		response.SafeInternalError(c, "Failed to upload evidence", err)
		return
	}

	response.Created(c, gin.H{
		"success":  true,
		"url":      result.URL,
		"filename": result.Filename,
	})
}

// ListEvidence handles GET /api/v1/portal/evidence
func (h *PortalHandler) ListEvidence(c *gin.Context) {
	if !h.checkStorage(c) {
		return
	}
	ctx := c.Request.Context()

	files, err := h.storage.ListFiles(ctx, h.cfg.MinioBucketEvidence, "")
	if err != nil {
		response.SafeInternalError(c, "Failed to list evidence files", err)
		return
	}

	response.Success(c, files)
}

// DeleteEvidence handles DELETE /api/v1/portal/evidence/:filename
func (h *PortalHandler) DeleteEvidence(c *gin.Context) {
	if !h.checkStorage(c) {
		return
	}
	ctx := c.Request.Context()
	filename := c.Param("filename")

	if filename == "" {
		response.BadRequest(c, "Filename is required")
		return
	}

	// Prevent path traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		response.BadRequest(c, "Invalid filename")
		return
	}

	// Delete from MinIO
	err := h.storage.DeleteFile(ctx, h.cfg.MinioBucketEvidence, filename)
	if err != nil {
		response.SafeInternalError(c, "Failed to delete evidence", err)
		return
	}

	response.Success(c, gin.H{
		"success": true,
		"message": "Evidence file removed",
	})
}
