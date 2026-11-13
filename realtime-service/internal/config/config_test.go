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
		t.Setenv("CORE_SERVICE_URL", "")
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
		t.Setenv("CORE_SERVICE_URL", "")
		t.Setenv("ACCESS_TOKEN_TTL_MINUTES", "")
		t.Setenv("JWT_EXPIRATION_HOURS", "")
		t.Setenv("CORS_ALLOWED_ORIGINS", "")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("unexpected error loading default config: %v", err)
		}

		if cfg.Port != "8001" {
			t.Errorf("expected port 8001, got %s", cfg.Port)
		}
		if cfg.Environment != "development" {
			t.Errorf("expected environment development, got %s", cfg.Environment)
		}
		if cfg.CoreServiceURL != "http://localhost:8000" {
			t.Errorf("expected coreServiceURL http://localhost:8000, got %s", cfg.CoreServiceURL)
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

	t.Run("invalid access token ttl returns format error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", testJWTSecret)
		for _, bad := range []string{"invalid", "0", "-5"} {
			t.Setenv("ACCESS_TOKEN_TTL_MINUTES", bad)
			if _, err := Load(); err == nil {
				t.Errorf("expected error for ACCESS_TOKEN_TTL_MINUTES=%q", bad)
			}
		}
	})

	t.Run("invalid jwt expiration hours returns format error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", testJWTSecret)
		for _, bad := range []string{"invalid", "0", "-1"} {
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
		t.Setenv("PORT", "9001")
		t.Setenv("ENVIRONMENT", "staging")
		t.Setenv("JWT_SECRET", "this-is-a-very-long-secret-key-for-testing-purposes-123")
		t.Setenv("CORE_SERVICE_URL", "http://core:8000")
		t.Setenv("ACCESS_TOKEN_TTL_MINUTES", "30")
		t.Setenv("JWT_EXPIRATION_HOURS", "24")
		t.Setenv("CORS_ALLOWED_ORIGINS", "https://app.palabrario.org, https://play.palabrario.org")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("unexpected error with valid env vars: %v", err)
		}

		if cfg.Port != "9001" || cfg.Environment != "staging" || cfg.CoreServiceURL != "http://core:8000" ||
			cfg.AccessTokenTTL != 30*time.Minute || cfg.JWTExpirationHours != 24*time.Hour ||
			len(cfg.AllowedOrigins) != 2 || cfg.AllowedOrigins[0] != "https://app.palabrario.org" {
			t.Errorf("custom env vars mismatch: %+v", cfg)
		}
	})
}
