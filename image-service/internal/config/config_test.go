package config

import (
	"testing"
	"time"
)

const testJWTSecret = "test-jwt-secret-key-palabrario-minimum-32-chars"

func TestConfigLoad(t *testing.T) {
	t.Run("missing jwt secret returns security error", func(t *testing.T) {
		t.Setenv("PORT", "")
		t.Setenv("ENVIRONMENT", "")
		t.Setenv("NODE_ENV", "")
		t.Setenv("JWT_SECRET", "")
		t.Setenv("UPLOAD_DIR", "")
		t.Setenv("PUBLIC_BASE_URL", "")
		t.Setenv("MAX_FILE_SIZE_MB", "")
		t.Setenv("ACCESS_TOKEN_TTL_MINUTES", "")
		t.Setenv("JWT_EXPIRATION_HOURS", "")
		t.Setenv("CORS_ALLOWED_ORIGINS", "")

		if _, err := Load(); err == nil {
			t.Fatal("expected error when JWT_SECRET is missing")
		}
	})

	t.Run("default values loaded cleanly when JWT_SECRET is set", func(t *testing.T) {
		t.Setenv("PORT", "")
		t.Setenv("ENVIRONMENT", "")
		t.Setenv("NODE_ENV", "")
		t.Setenv("JWT_SECRET", testJWTSecret)
		t.Setenv("UPLOAD_DIR", "")
		t.Setenv("PUBLIC_BASE_URL", "")
		t.Setenv("MAX_FILE_SIZE_MB", "")
		t.Setenv("ACCESS_TOKEN_TTL_MINUTES", "")
		t.Setenv("JWT_EXPIRATION_HOURS", "")
		t.Setenv("CORS_ALLOWED_ORIGINS", "")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("unexpected error loading default config: %v", err)
		}

		if cfg.Port != "8002" {
			t.Errorf("expected port 8002, got %s", cfg.Port)
		}
		if cfg.Environment != "development" {
			t.Errorf("expected environment development, got %s", cfg.Environment)
		}
		if cfg.MaxFileSizeMB != 5 {
			t.Errorf("expected maxFileSizeMB 5, got %d", cfg.MaxFileSizeMB)
		}
		if cfg.AccessTokenTTL != 15*time.Minute {
			t.Errorf("expected AccessTokenTTL 15m, got %v", cfg.AccessTokenTTL)
		}
		if cfg.JWTExpirationHours != 72*time.Hour {
			t.Errorf("expected JWTExpirationHours 72h, got %v", cfg.JWTExpirationHours)
		}
		if len(cfg.AllowedOrigins) == 0 {
			t.Error("expected default dev AllowedOrigins to be populated")
		}
	})

	t.Run("invalid port returns format error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", testJWTSecret)
		t.Setenv("PORT", "not-a-port")
		if _, err := Load(); err == nil {
			t.Error("expected error for non-numeric port")
		}
	})

	t.Run("short jwt secret returns security error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "short-secret")
		if _, err := Load(); err == nil {
			t.Error("expected security error for JWT_SECRET shorter than 32 chars")
		}
	})

	t.Run("invalid max file size returns error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", testJWTSecret)
		t.Setenv("PORT", "8002")
		t.Setenv("MAX_FILE_SIZE_MB", "invalid-size")
		if _, err := Load(); err == nil {
			t.Error("expected error for non-numeric MAX_FILE_SIZE_MB")
		}
	})

	t.Run("invalid access token ttl returns format error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", testJWTSecret)
		for _, bad := range []string{"invalid", "0", "-10"} {
			t.Setenv("ACCESS_TOKEN_TTL_MINUTES", bad)
			if _, err := Load(); err == nil {
				t.Errorf("expected error for ACCESS_TOKEN_TTL_MINUTES=%q", bad)
			}
		}
	})

	t.Run("invalid jwt expiration hours returns format error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", testJWTSecret)
		for _, bad := range []string{"invalid", "0", "-3"} {
			t.Setenv("JWT_EXPIRATION_HOURS", bad)
			if _, err := Load(); err == nil {
				t.Errorf("expected error for JWT_EXPIRATION_HOURS=%q", bad)
			}
		}
	})

	t.Run("production requires explicit CORS_ALLOWED_ORIGINS", func(t *testing.T) {
		t.Setenv("JWT_SECRET", testJWTSecret)
		t.Setenv("ENVIRONMENT", "production")
		t.Setenv("CORS_ALLOWED_ORIGINS", "")
		if _, err := Load(); err == nil {
			t.Error("expected error when CORS_ALLOWED_ORIGINS is empty in production")
		}
	})

	t.Run("custom valid environment variables parsed successfully", func(t *testing.T) {
		t.Setenv("PORT", "9002")
		t.Setenv("ENVIRONMENT", "staging")
		t.Setenv("JWT_SECRET", "this-is-a-very-long-secret-key-for-testing-purposes-123")
		t.Setenv("UPLOAD_DIR", "/tmp/uploads")
		t.Setenv("PUBLIC_BASE_URL", "https://media.palabrario.org")
		t.Setenv("MAX_FILE_SIZE_MB", "10")
		t.Setenv("ACCESS_TOKEN_TTL_MINUTES", "20")
		t.Setenv("JWT_EXPIRATION_HOURS", "12")
		t.Setenv("CORS_ALLOWED_ORIGINS", "https://app.palabrario.org")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("unexpected error with valid env vars: %v", err)
		}

		if cfg.Port != "9002" || cfg.Environment != "staging" || cfg.UploadDir != "/tmp/uploads" ||
			cfg.PublicBaseURL != "https://media.palabrario.org" || cfg.MaxFileSizeMB != 10 ||
			cfg.AccessTokenTTL != 20*time.Minute || cfg.JWTExpirationHours != 12*time.Hour ||
			len(cfg.AllowedOrigins) != 1 || cfg.AllowedOrigins[0] != "https://app.palabrario.org" {
			t.Errorf("custom env vars mismatch: %+v", cfg)
		}
	})
}
