package storage

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalDiskStorage(t *testing.T) {
	uploadDir := t.TempDir()
	baseURL := "http://localhost:8002"
	s := NewLocalDiskStorage(uploadDir, baseURL)

	ctx := context.Background()
	data := []byte("fake-image-bytes-content")

	url, err := s.Save(ctx, "avatar", "png", data)
	if err != nil {
		t.Fatalf("error guardando archivo: %v", err)
	}

	expectedPrefix := "http://localhost:8002/uploads/avatar/"
	if !strings.HasPrefix(url, expectedPrefix) {
		t.Errorf("URL esperada con prefijo %s, obtuvo: %s", expectedPrefix, url)
	}
	if !strings.HasSuffix(url, ".png") {
		t.Errorf("URL esperada con extensión .png, obtuvo: %s", url)
	}

	// Verificar que el archivo realmente existe en el disco
	filename := strings.TrimPrefix(url, expectedPrefix)
	savedPath := filepath.Join(uploadDir, "avatar", filename)

	readData, err := os.ReadFile(savedPath)
	if err != nil {
		t.Fatalf("error leyendo archivo guardado en disco: %v", err)
	}
	if string(readData) != string(data) {
		t.Errorf("contenido leído no coincide con los datos guardados")
	}
}
