package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"palabrario/pkg/auth"
)

// Config centralizes and types all environment variables for core-service.
// Follows the 12-Factor App methodology with strict fail-fast validation.
type Config struct {
	Port               string
	Environment        string
	DatabaseURL        string
	JWTSecret          string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
	JWTExpirationHours time.Duration
	ImageServiceURL    string
	RealtimeServiceURL string
	AllowedOrigins     []string
}

// Load reads, strictly validates, and structures environment variables.
// Fails fast with explicit details instead of falling back to insecure defaults.
func Load() (*Config, error) {
	var missingVars []string

	requireEnv := func(name string) string {
		val := strings.TrimSpace(os.Getenv(name))
		if val == "" {
			missingVars = append(missingVars, name)
		}
		return val
	}

	port := requireEnv("PORT")
	dbURL := requireEnv("DATABASE_URL")
	jwtSecret := requireEnv("JWT_SECRET")

	env := strings.TrimSpace(os.Getenv("ENVIRONMENT"))
	if env == "" {
		env = strings.TrimSpace(os.Getenv("NODE_ENV"))
	}
	if env == "" {
		missingVars = append(missingVars, "ENVIRONMENT")
	}

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("fatal configuration error: missing required environment variables: [%s]", strings.Join(missingVars, ", "))
	}

	if len(jwtSecret) < 32 {
		return nil, fmt.Errorf("security configuration error: JWT_SECRET must be at least 32 characters long (got %d)", len(jwtSecret))
	}

	if _, err := strconv.Atoi(port); err != nil {
		return nil, fmt.Errorf("configuration format error: PORT %q is not a valid TCP port", port)
	}

	accessTTL := auth.DefaultAccessTokenTTL
	if ttlStr := strings.TrimSpace(os.Getenv("ACCESS_TOKEN_TTL_MINUTES")); ttlStr != "" {
		val, err := strconv.Atoi(ttlStr)
		if err != nil || val <= 0 {
			return nil, fmt.Errorf("configuration format error: ACCESS_TOKEN_TTL_MINUTES %q must be a positive integer", ttlStr)
		}
		accessTTL = time.Duration(val) * time.Minute
	}

	refreshHours := 168
	if refreshStr := strings.TrimSpace(os.Getenv("REFRESH_TOKEN_TTL_HOURS")); refreshStr != "" {
		val, err := strconv.Atoi(refreshStr)
		if err != nil || val <= 0 {
			return nil, fmt.Errorf("configuration format error: REFRESH_TOKEN_TTL_HOURS %q must be a positive integer", refreshStr)
		}
		refreshHours = val
	} else if expStr := strings.TrimSpace(os.Getenv("JWT_EXPIRATION_HOURS")); expStr != "" {
		val, err := strconv.Atoi(expStr)
		if err != nil || val <= 0 {
			return nil, fmt.Errorf("configuration format error: JWT_EXPIRATION_HOURS %q must be a positive integer", expStr)
		}
		refreshHours = val
	}
	refreshTTL := time.Duration(refreshHours) * time.Hour

	corsRaw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if corsRaw == "" && strings.EqualFold(env, "production") {
		return nil, fmt.Errorf("security configuration error: CORS_ALLOWED_ORIGINS must be set in production")
	}
	allowedOrigins := auth.ParseCORSAllowedOrigins(corsRaw)

	return &Config{
		Port:               port,
		Environment:        env,
		DatabaseURL:        dbURL,
		JWTSecret:          jwtSecret,
		AccessTokenTTL:     accessTTL,
		RefreshTokenTTL:    refreshTTL,
		JWTExpirationHours: refreshTTL,
		ImageServiceURL:    strings.TrimSpace(os.Getenv("IMAGE_SERVICE_URL")),
		RealtimeServiceURL: strings.TrimSpace(os.Getenv("REALTIME_SERVICE_URL")),
		AllowedOrigins:     allowedOrigins,
	}, nil
}
