package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"palabrario/pkg/auth"
)

// Config centralizes and types all environment variables for image-service.
// Follows the 12-Factor App methodology with strict fail-fast validation.
type Config struct {
	Port               string
	Environment        string
	JWTSecret          string
	AccessTokenTTL     time.Duration
	JWTExpirationHours time.Duration
	UploadDir          string
	PublicBaseURL      string
	MaxFileSizeMB      int64
	AllowedOrigins     []string
}

// Load reads, strictly validates, and structures the environment variables for image-service.
func Load() (*Config, error) {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8002"
	}
	if _, err := strconv.Atoi(port); err != nil {
		return nil, fmt.Errorf("configuration format error: PORT %q is not a valid TCP port", port)
	}

	env := strings.TrimSpace(os.Getenv("ENVIRONMENT"))
	if env == "" {
		env = strings.TrimSpace(os.Getenv("NODE_ENV"))
	}
	if env == "" {
		env = "development"
	}

	jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if jwtSecret == "" {
		return nil, fmt.Errorf("security configuration error: JWT_SECRET is required")
	}
	if len(jwtSecret) < 32 {
		return nil, fmt.Errorf("security configuration error: JWT_SECRET must be at least 32 characters long (got %d)", len(jwtSecret))
	}

	uploadDir := strings.TrimSpace(os.Getenv("UPLOAD_DIR"))
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	publicBaseURL := strings.TrimSpace(os.Getenv("PUBLIC_BASE_URL"))
	if publicBaseURL == "" {
		publicBaseURL = fmt.Sprintf("http://localhost:%s", port)
	}

	var maxFileSizeMB int64 = 5
	if maxStr := strings.TrimSpace(os.Getenv("MAX_FILE_SIZE_MB")); maxStr != "" {
		val, err := strconv.ParseInt(maxStr, 10, 64)
		if err != nil || val <= 0 {
			return nil, fmt.Errorf("configuration format error: MAX_FILE_SIZE_MB %q must be a positive integer", maxStr)
		}
		maxFileSizeMB = val
	}

	accessTTL := auth.DefaultAccessTokenTTL
	if ttlStr := strings.TrimSpace(os.Getenv("ACCESS_TOKEN_TTL_MINUTES")); ttlStr != "" {
		val, err := strconv.Atoi(ttlStr)
		if err != nil || val <= 0 {
			return nil, fmt.Errorf("configuration format error: ACCESS_TOKEN_TTL_MINUTES %q must be a positive integer", ttlStr)
		}
		accessTTL = time.Duration(val) * time.Minute
	}

	expHours := 72
	if expStr := strings.TrimSpace(os.Getenv("JWT_EXPIRATION_HOURS")); expStr != "" {
		val, err := strconv.Atoi(expStr)
		if err != nil || val <= 0 {
			return nil, fmt.Errorf("configuration format error: JWT_EXPIRATION_HOURS %q must be a positive integer", expStr)
		}
		expHours = val
	}

	corsRaw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if corsRaw == "" && strings.EqualFold(env, "production") {
		return nil, fmt.Errorf("security configuration error: CORS_ALLOWED_ORIGINS must be set in production")
	}
	allowedOrigins := auth.ParseCORSAllowedOrigins(corsRaw)

	return &Config{
		Port:               port,
		Environment:        env,
		JWTSecret:          jwtSecret,
		AccessTokenTTL:     accessTTL,
		JWTExpirationHours: time.Duration(expHours) * time.Hour,
		UploadDir:          uploadDir,
		PublicBaseURL:      publicBaseURL,
		MaxFileSizeMB:      maxFileSizeMB,
		AllowedOrigins:     allowedOrigins,
	}, nil
}
