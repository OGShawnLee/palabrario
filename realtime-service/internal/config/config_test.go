package config

import (
	"os"
	"testing"
	"time"
)

func TestConfigLoad(t *testing.T) {
	t.Run("default values loaded cleanly", func(t *testing.T) {
		_ = os.Unsetenv("PORT")
		_ = os.Unsetenv("ENVIRONMENT")
		_ = os.Unsetenv("NODE_ENV")
		_ = os.Unsetenv("JWT_SECRET")
		_ = os.Unsetenv("CORE_SERVICE_URL")
		_ = os.Unsetenv("JWT_EXPIRATION_HOURS")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("error inesperado al cargar config por defecto: %v", err)
		}

		if cfg.Port != "8001" {
			t.Errorf("esperaba puerto 8001, obtuvo %s", cfg.Port)
		}
		if cfg.Environment != "development" {
			t.Errorf("esperaba environment development, obtuvo %s", cfg.Environment)
		}
		if cfg.CoreServiceURL != "http://localhost:8000" {
			t.Errorf("esperaba coreServiceURL http://localhost:8000, obtuvo %s", cfg.CoreServiceURL)
		}
		if cfg.JWTExpirationHours != 72*time.Hour {
			t.Errorf("esperaba JWTExpirationHours 72h, obtuvo %v", cfg.JWTExpirationHours)
		}
		if len(cfg.JWTSecret) < 32 {
			t.Errorf("JWT_SECRET por defecto debe tener al menos 32 caracteres")
		}
	})

	t.Run("invalid port returns format error", func(t *testing.T) {
		t.Setenv("PORT", "not-a-port")
		_, err := Load()
		if err == nil {
			t.Errorf("esperaba error por puerto no numérico")
		}
	})

	t.Run("short jwt secret returns security error", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "short-secret")
		_, err := Load()
		if err == nil {
			t.Errorf("esperaba error de seguridad por JWT_SECRET menor a 32 caracteres")
		}
	})

	t.Run("invalid jwt expiration hours returns format error", func(t *testing.T) {
		t.Setenv("JWT_EXPIRATION_HOURS", "invalid")
		_, err := Load()
		if err == nil {
			t.Errorf("esperaba error por JWT_EXPIRATION_HOURS no numérico")
		}
	})

	t.Run("custom valid environment variables parsed successfully", func(t *testing.T) {
		t.Setenv("PORT", "9001")
		t.Setenv("ENVIRONMENT", "staging")
		t.Setenv("JWT_SECRET", "this-is-a-very-long-secret-key-for-testing-purposes-123")
		t.Setenv("CORE_SERVICE_URL", "http://core:8000")
		t.Setenv("JWT_EXPIRATION_HOURS", "24")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("error inesperado con variables válidas: %v", err)
		}

		if cfg.Port != "9001" || cfg.Environment != "staging" || cfg.CoreServiceURL != "http://core:8000" || cfg.JWTExpirationHours != 24*time.Hour {
			t.Errorf("variables personalizadas no coinciden con la configuración cargada: %+v", cfg)
		}
	})
}
