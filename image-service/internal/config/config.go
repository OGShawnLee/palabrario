package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config centraliza y tipa todas las variables de entorno de image-service.
// Sigue la metodología 12-Factor App con validación estructurada.
type Config struct {
	Port               string
	Environment        string
	JWTSecret          string
	JWTExpirationHours time.Duration
	UploadDir          string
	PublicBaseURL      string
	MaxFileSizeMB      int64
}

// Load lee, valida y estructura las variables de entorno de image-service.
func Load() (*Config, error) {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8002"
	}
	if _, err := strconv.Atoi(port); err != nil {
		return nil, fmt.Errorf("error de formato en configuración: PORT '%s' no es un puerto TCP válido", port)
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
		jwtSecret = "super-secret-jwt-key-palabrario-2026"
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
		if val, err := strconv.ParseInt(maxStr, 10, 64); err == nil && val > 0 {
			maxFileSizeMB = val
		} else {
			return nil, fmt.Errorf("error de formato en configuración: MAX_FILE_SIZE_MB '%s' debe ser un entero positivo", maxStr)
		}
	}

	expHours := 72
	if expStr := strings.TrimSpace(os.Getenv("JWT_EXPIRATION_HOURS")); expStr != "" {
		if val, err := strconv.Atoi(expStr); err == nil && val > 0 {
			expHours = val
		}
	}

	return &Config{
		Port:               port,
		Environment:        env,
		JWTSecret:          jwtSecret,
		JWTExpirationHours: time.Duration(expHours) * time.Hour,
		UploadDir:          uploadDir,
		PublicBaseURL:      publicBaseURL,
		MaxFileSizeMB:      maxFileSizeMB,
	}, nil
}
