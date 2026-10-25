package config

import (
	"os"
	"testing"
)

func TestConfigLoad(t *testing.T) {
	t.Run("default values loaded cleanly", func(t *testing.T) {
		_ = os.Unsetenv("PORT")
		_ = os.Unsetenv("ENVIRONMENT")
		_ = os.Unsetenv("JWT_SECRET")
		_ = os.Unsetenv("UPLOAD_DIR")
		_ = os.Unsetenv("PUBLIC_BASE_URL")
		_ = os.Unsetenv("MAX_FILE_SIZE_MB")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("error inesperado al cargar config por defecto: %v", err)
		}

		if cfg.Port != "8002" {
			t.Errorf("esperaba puerto 8002, obtuvo %s", cfg.Port)
		}
		if cfg.Environment != "development" {
			t.Errorf("esperaba environment development, obtuvo %s", cfg.Environment)
		}
		if cfg.MaxFileSizeMB != 5 {
			t.Errorf("esperaba maxFileSizeMB 5, obtuvo %d", cfg.MaxFileSizeMB)
		}
	})

	t.Run("invalid port returns format error", func(t *testing.T) {
		t.Setenv("PORT", "not-a-port")
		_, err := Load()
		if err == nil {
			t.Errorf("esperaba error por puerto no numérico")
		}
	})

	t.Run("invalid max file size returns error", func(t *testing.T) {
		t.Setenv("PORT", "8002")
		t.Setenv("MAX_FILE_SIZE_MB", "invalid-size")
		_, err := Load()
		if err == nil {
			t.Errorf("esperaba error por tamaño no numérico")
		}
	})
}
