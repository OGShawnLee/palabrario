package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"palabrario/pkg/auth"
)

// Config centralizes and types all environment variables for realtime-service.
// Follows the 12-Factor App methodology with strict fail-fast validation.
type Config struct {
	Port               string
	Environment        string
	JWTSecret          string
	AccessTokenTTL     time.Duration
	JWTExpirationHours time.Duration
	CoreServiceURL     string
	AllowedOrigins     []string
}

// Load reads, strictly validates, and structures the environment variables for realtime-service.
func Load() (*Config, error) {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8001"
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

	coreServiceURL := strings.TrimSpace(os.Getenv("CORE_SERVICE_URL"))
	if coreServiceURL == "" {
		coreServiceURL = "http://localhost:8000"
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
		CoreServiceURL:     coreServiceURL,
		AllowedOrigins:     allowedOrigins,
	}, nil
}
