package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/condotrack/api/internal/config"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// StorageService handles file storage operations
type StorageService struct {
	client   *minio.Client
	cfg      *config.Config
	endpoint string
}

// FileInfo represents information about a stored file
type FileInfo struct {
	Name         string    `json:"name"`
	URL          string    `json:"url"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	LastModified time.Time `json:"last_modified"`
	Tags         []string  `json:"tags,omitempty"`
	IsSystem     bool      `json:"is_system,omitempty"`
}

// UploadResult represents the result of an upload operation
type UploadResult struct {
	Filename    string `json:"filename"`
	URL         string `json:"url"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	Bucket      string `json:"bucket"`
}

// NewStorageService creates a new MinIO storage service
func NewStorageService(cfg *config.Config) (*StorageService, error) {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &StorageService{
		client:   client,
		cfg:      cfg,
		endpoint: cfg.MinioEndpoint,
	}, nil
}

// UploadFile uploads a file to the specified bucket
func (s *StorageService) UploadFile(ctx context.Context, bucket string, filename string, reader io.Reader, size int64, contentType string) (*UploadResult, error) {
	// Generate unique filename if not provided
	if filename == "" {
		filename = fmt.Sprintf("%d_%s", time.Now().Unix(), uuid.New().String())
	}

	// Upload to MinIO
	info, err := s.client.PutObject(ctx, bucket, filename, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return &UploadResult{
		Filename:    filename,
		URL:         s.GetPublicURL(bucket, filename),
		Size:        info.Size,
		ContentType: contentType,
		Bucket:      bucket,
	}, nil
}

// UploadBase64 uploads a base64-encoded file
func (s *StorageService) UploadBase64(ctx context.Context, bucket string, filename string, base64Data string) (*UploadResult, error) {
	// Remove data URI prefix if present
	data := base64Data
	contentType := "application/octet-stream"

	if strings.HasPrefix(data, "data:") {
		parts := strings.SplitN(data, ",", 2)
		if len(parts) == 2 {
			// Extract content type from data URI
			header := parts[0]
			data = parts[1]

			// Parse content type (e.g., "data:image/png;base64")
			if strings.Contains(header, ":") && strings.Contains(header, ";") {
				ctParts := strings.SplitN(header, ":", 2)
				if len(ctParts) == 2 {
					ctEnd := strings.Index(ctParts[1], ";")
					if ctEnd > 0 {
						contentType = ctParts[1][:ctEnd]
					}
				}
			}
		}
	}

	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Determine extension from content type
	ext := getExtensionFromContentType(contentType)
	if ext != "" && !strings.HasSuffix(filename, ext) {
		filename = strings.TrimSuffix(filename, filepath.Ext(filename)) + ext
	}

	reader := bytes.NewReader(decoded)
	return s.UploadFile(ctx, bucket, filename, reader, int64(len(decoded)), contentType)
}

// ListFiles lists all files in a bucket
func (s *StorageService) ListFiles(ctx context.Context, bucket string, prefix string) ([]FileInfo, error) {
	var files []FileInfo

	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}

	for object := range s.client.ListObjects(ctx, bucket, opts) {
		if object.Err != nil {
			return nil, fmt.Errorf("error listing objects: %w", object.Err)
		}

		files = append(files, FileInfo{
			Name:         object.Key,
			URL:          s.GetPublicURL(bucket, object.Key),
			Size:         object.Size,
			ContentType:  object.ContentType,
			LastModified: object.LastModified,
		})
	}

	return files, nil
}

// DeleteFile deletes a file from a bucket
func (s *StorageService) DeleteFile(ctx context.Context, bucket string, filename string) error {
	err := s.client.RemoveObject(ctx, bucket, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetFile retrieves a file from a bucket
func (s *StorageService) GetFile(ctx context.Context, bucket string, filename string) (io.ReadCloser, *minio.ObjectInfo, error) {
	obj, err := s.client.GetObject(ctx, bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file: %w", err)
	}

	info, err := obj.Stat()
	if err != nil {
		obj.Close()
		return nil, nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return obj, &info, nil
}

// FileExists checks if a file exists in a bucket
func (s *StorageService) FileExists(ctx context.Context, bucket string, filename string) bool {
	_, err := s.client.StatObject(ctx, bucket, filename, minio.StatObjectOptions{})
	return err == nil
}

// GetPublicURL returns the public URL for a file
func (s *StorageService) GetPublicURL(bucket string, filename string) string {
	// Use configured public URL if available
	if s.cfg.MinioPublicURL != "" {
		return fmt.Sprintf("%s/%s/%s", s.cfg.MinioPublicURL, bucket, url.PathEscape(filename))
	}

	// Fallback to internal endpoint
	protocol := "http"
	if s.cfg.MinioUseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, s.endpoint, bucket, url.PathEscape(filename))
}

// GetPresignedURL generates a presigned URL for temporary access
func (s *StorageService) GetPresignedURL(ctx context.Context, bucket string, filename string, expiry time.Duration) (string, error) {
	presignedURL, err := s.client.PresignedGetObject(ctx, bucket, filename, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return presignedURL.String(), nil
}

// CopyFile copies a file within or between buckets
func (s *StorageService) CopyFile(ctx context.Context, srcBucket, srcFile, dstBucket, dstFile string) error {
	src := minio.CopySrcOptions{
		Bucket: srcBucket,
		Object: srcFile,
	}
	dst := minio.CopyDestOptions{
		Bucket: dstBucket,
		Object: dstFile,
	}

	_, err := s.client.CopyObject(ctx, dst, src)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	return nil
}

// Helper function to get file extension from content type
func getExtensionFromContentType(contentType string) string {
	extensions := map[string]string{
		"image/jpeg":      ".jpg",
		"image/png":       ".png",
		"image/gif":       ".gif",
		"image/webp":      ".webp",
		"image/svg+xml":   ".svg",
		"application/pdf": ".pdf",
		"video/mp4":       ".mp4",
		"video/webm":      ".webm",
		"audio/mpeg":      ".mp3",
		"audio/wav":       ".wav",
	}

	if ext, ok := extensions[contentType]; ok {
		return ext
	}
	return ""
}

// GetContentTypeFromExtension returns content type based on file extension
func GetContentTypeFromExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
		".pdf":  "application/pdf",
		".mp4":  "video/mp4",
		".webm": "video/webm",
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	}

	if ct, ok := contentTypes[ext]; ok {
		return ct
	}
	return "application/octet-stream"
}

// IsAllowedImageType checks if the content type is an allowed image type
func IsAllowedImageType(contentType string) bool {
	allowed := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/svg+xml",
	}

	for _, a := range allowed {
		if contentType == a {
			return true
		}
	}
	return false
}

// IsAllowedEvidenceType checks if the content type is allowed for evidence files
func IsAllowedEvidenceType(contentType string) bool {
	allowed := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	}

	for _, a := range allowed {
		if contentType == a {
			return true
		}
	}
	return false
}
