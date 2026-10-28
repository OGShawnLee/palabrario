package main

import (
	"bytes"
	"encoding/json"
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

// UploadImage envía un archivo multipart/form-data a image-service (:8002).
func (a *App) UploadImage(category string, filename string, fileBytes []byte) (*models.ImageUploadResponse, error) {
	category = strings.TrimSpace(category)
	switch category {
	case "avatar", "credentials", "word", "sense":
	default:
		return nil, errors.New("categoría inválida; debe ser estrictamente 'avatar', 'credentials', 'word' o 'sense'")
	}

	if filename == "" {
		filename = "image.webp"
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("category", category); err != nil {
		return nil, fmt.Errorf("error al escribir campo category: %w", err)
	}

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("error al crear campo de archivo multipart: %w", err)
	}

	if _, err := io.Copy(part, bytes.NewReader(fileBytes)); err != nil {
		return nil, fmt.Errorf("error al escribir bytes del archivo: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("error al finalizar formulario multipart: %w", err)
	}

	req, err := http.NewRequestWithContext(a.context(), http.MethodPost, a.imageURL+"/api/v1/images/upload", body)
	if err != nil {
		return nil, fmt.Errorf("error al crear petición HTTP multipart: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	token := a.getSessionToken()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error de red al conectar con image-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error string `json:"error"`
		}
		respBody, _ := io.ReadAll(resp.Body)
		if len(respBody) > 0 {
			if jsonErr := json.Unmarshal(respBody, &errResp); jsonErr == nil && errResp.Error != "" {
				return nil, errors.New(errResp.Error)
			}
		}
		return nil, errors.New(http.StatusText(resp.StatusCode))
	}

	var uploadResp models.ImageUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta JSON de imagen: %w", err)
	}

	return &uploadResp, nil
}
