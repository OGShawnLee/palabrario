package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"palabrario/pkg/models"
)

// =============================================================================
// MÓDULO 8: Multimedia Seguro (image-service :8002)
// =============================================================================

// UploadImage uploads a multipart/form-data file to image-service (:8002).
func (a *App) UploadImage(category string, filename string, fileBytes []byte) (*models.ImageUploadResponse, error) {
	category = strings.TrimSpace(category)
	switch category {
	case "avatar", "credentials", "word", "sense":
	default:
		return nil, errors.New("invalid category; must be 'avatar', 'credentials', 'word', or 'sense'")
	}

	if filename == "" {
		filename = "image.webp"
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("category", category); err != nil {
		return nil, fmt.Errorf("failed to write category field: %w", err)
	}

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create multipart file field: %w", err)
	}

	if _, err := io.Copy(part, bytes.NewReader(fileBytes)); err != nil {
		return nil, fmt.Errorf("failed to write file bytes: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize multipart form: %w", err)
	}

	payload := body.Bytes()
	contentType := writer.FormDataContentType()

	var uploadResp models.ImageUploadResponse
	err = a.executeRequest(a.context(), func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(a.context(), http.MethodPost, a.imageURL+"/api/v1/images/upload", bytes.NewReader(payload))
		if err != nil {
			return nil, fmt.Errorf("failed to create multipart HTTP request: %w", err)
		}
		req.Header.Set("Content-Type", contentType)
		return req, nil
	}, &uploadResp, true)
	if err != nil {
		return nil, err
	}

	return &uploadResp, nil
}
