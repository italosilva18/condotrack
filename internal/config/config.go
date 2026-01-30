package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Server
	ServerPort string
	AppEnv     string

	// Database
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string

	// JWT Authentication
	JWTSecret     string
	JWTExpiration int // hours

	// Asaas
	AsaasAPIKey       string
	AsaasAPIURL       string
	AsaasWebhookToken string
	AsaasEnv          string

	// Revenue Split
	RevenueInstructorPercent float64
	RevenuePlatformPercent   float64

	// Upload
	UploadDir     string
	MaxUploadSize int64

	// MinIO
	MinioEndpoint        string
	MinioAccessKey       string
	MinioSecretKey       string
	MinioUseSSL          bool
	MinioPublicURL       string // External URL for accessing files
	MinioBucketUploads   string
	MinioBucketPortal    string
	MinioBucketEvidence  string
	MinioBucketCerts     string

	// AI (Gemini)
	GeminiAPIKey string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{
		// Server
		ServerPort: getEnv("SERVER_PORT", "8000"),
		AppEnv:     getEnv("APP_ENV", "development"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "condotrack"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASS", ""),

		// JWT Authentication
		JWTSecret:     getEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
		JWTExpiration: getEnvInt("JWT_EXPIRATION_HOURS", 24),

		// Asaas
		AsaasAPIKey:       getEnv("ASAAS_API_KEY", ""),
		AsaasAPIURL:       getEnv("ASAAS_API_URL", "https://sandbox.asaas.com/api/v3"),
		AsaasWebhookToken: getEnv("ASAAS_WEBHOOK_TOKEN", ""),
		AsaasEnv:          getEnv("ASAAS_ENV", "sandbox"),

		// Revenue Split
		RevenueInstructorPercent: getEnvFloat("REVENUE_INSTRUCTOR_PERCENT", 70.0),
		RevenuePlatformPercent:   getEnvFloat("REVENUE_PLATFORM_PERCENT", 30.0),

		// Upload
		UploadDir:     getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadSize: getEnvInt64("MAX_UPLOAD_SIZE", 50*1024*1024), // 50MB default

		// MinIO
		MinioEndpoint:       getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey:      getEnv("MINIO_ACCESS_KEY", "condotrack"),
		MinioSecretKey:      getEnv("MINIO_SECRET_KEY", "Condo@2024Minio"),
		MinioUseSSL:         getEnvBool("MINIO_USE_SSL", false),
		MinioPublicURL:      getEnv("MINIO_PUBLIC_URL", ""), // e.g., http://localhost:9002
		MinioBucketUploads:  getEnv("MINIO_BUCKET_UPLOADS", "uploads"),
		MinioBucketPortal:   getEnv("MINIO_BUCKET_PORTAL", "portal-images"),
		MinioBucketEvidence: getEnv("MINIO_BUCKET_EVIDENCE", "evidence"),
		MinioBucketCerts:    getEnv("MINIO_BUCKET_CERTIFICATES", "certificates"),

		// AI (Gemini)
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
	}

	return cfg, nil
}

// getEnv returns environment variable value or default
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvFloat returns environment variable as float64 or default
func getEnvFloat(key string, defaultValue float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	}
	return defaultValue
}

// getEnvInt64 returns environment variable as int64 or default
func getEnvInt64(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}
	}
	return defaultValue
}

// getEnvInt returns environment variable as int or default
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

// getEnvBool returns environment variable as bool or default
func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

// GetDSN returns the MySQL connection string
func (c *Config) GetDSN() string {
	return c.DBUser + ":" + c.DBPassword + "@tcp(" + c.DBHost + ":" + c.DBPort + ")/" + c.DBName + "?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}
