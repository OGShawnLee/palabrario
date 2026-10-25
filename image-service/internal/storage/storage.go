package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// StorageProvider abstrae el almacenamiento de artefactos de imagen.
type StorageProvider interface {
	Save(ctx context.Context, category string, ext string, data []byte) (publicURL string, err error)
}

// LocalDiskStorage implementa StorageProvider guardando archivos en el disco local y exponiéndolos vía URL estática.
type LocalDiskStorage struct {
	uploadDir     string
	publicBaseURL string
}

// NewLocalDiskStorage construye un proveedor de almacenamiento en disco.
func NewLocalDiskStorage(uploadDir, publicBaseURL string) *LocalDiskStorage {
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	if publicBaseURL == "" {
		publicBaseURL = "http://localhost:8002"
	}
	return &LocalDiskStorage{
		uploadDir:     uploadDir,
		publicBaseURL: publicBaseURL,
	}
}

func (s *LocalDiskStorage) Save(ctx context.Context, category string, ext string, data []byte) (string, error) {
	categoryDir := filepath.Join(s.uploadDir, category)
	if err := os.MkdirAll(categoryDir, 0755); err != nil {
		return "", fmt.Errorf("error creando directorio de almacenamiento para categoría %s: %w", category, err)
	}

	id := uuid.NewString()
	cleanExt := strings.TrimPrefix(ext, ".")
	fileName := fmt.Sprintf("%s.%s", id, cleanExt)
	targetPath := filepath.Join(categoryDir, fileName)

	if err := os.WriteFile(targetPath, data, 0644); err != nil {
		return "", fmt.Errorf("error escribiendo imagen en disco: %w", err)
	}

	publicURL := fmt.Sprintf("%s/uploads/%s/%s.%s", strings.TrimRight(s.publicBaseURL, "/"), category, id, cleanExt)
	return publicURL, nil
}
