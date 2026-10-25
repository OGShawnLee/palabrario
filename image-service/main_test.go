package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/response"
)

func generateTestJWT(t *testing.T, secret string) string {
	t.Helper()
	tokenService := auth.NewTokenService(secret, 24*time.Hour)
	user := &models.UserProfile{
		IDUser:        "test-user-uuid-1",
		Email:         "tester@palabrario.org",
		UserName:      "tester",
		Role:          models.RoleEditor,
		AnchorDialect: "MX",
	}
	token, err := tokenService.GenerateToken(user)
	if err != nil {
		t.Fatalf("error generando token de prueba: %v", err)
	}
	return token
}

func createMultipartBody(category, filename string, fileData []byte) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if category != "" {
		if err := writer.WriteField("category", category); err != nil {
			return nil, "", err
		}
	}

	if filename != "" || fileData != nil {
		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			return nil, "", err
		}
		if _, err := part.Write(fileData); err != nil {
			return nil, "", err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func createSolidPNG(w, h int, c color.Color) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func createNSFWImage(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	// Variamos luminancia Y a través de múltiples bins asegurando tonos de piel en YCbCr y RGB
	luminances := []uint8{50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160, 170, 180, 190, 200}
	for y := 0; y < h; y++ {
		lum := luminances[y%len(luminances)]
		r, g, b := color.YCbCrToRGB(lum, 105, 150)
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func TestImageService(t *testing.T) {
	jwtSecret := "super-test-secret-palabrario-2026"
	uploadDir := t.TempDir()

	// Creamos un servidor de prueba provisional para obtener su dirección URL
	var ts *httptest.Server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Se delega al NewServer pasando ts.URL como publicBaseURL
		srv := NewServer(jwtSecret, uploadDir, ts.URL, 1) // 1 MB límite
		srv.ServeHTTP(w, r)
	})
	ts = httptest.NewServer(handler)
	defer ts.Close()

	validToken := generateTestJWT(t, jwtSecret)

	t.Run("1. 401 Unauthorized cuando falta token o es inválido", func(t *testing.T) {
		validImg := createSolidPNG(10, 10, color.RGBA{R: 50, G: 100, B: 150, A: 255})
		body, contentType, err := createMultipartBody("avatar", "test.png", validImg)
		if err != nil {
			t.Fatalf("error creando multipart: %v", err)
		}

		// Sin token
		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/images/upload", bytes.NewReader(body.Bytes()))
		req.Header.Set("Content-Type", contentType)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("fallo petición HTTP: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("esperaba 401 Unauthorized sin token, obtuvo %d", resp.StatusCode)
		}

		// Con token inválido
		reqInvalid, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/images/upload", bytes.NewReader(body.Bytes()))
		reqInvalid.Header.Set("Content-Type", contentType)
		reqInvalid.Header.Set("Authorization", "Bearer token-completamente-invalido.123")

		respInvalid, err := http.DefaultClient.Do(reqInvalid)
		if err != nil {
			t.Fatalf("fallo petición HTTP: %v", err)
		}
		defer respInvalid.Body.Close()

		if respInvalid.StatusCode != http.StatusUnauthorized {
			t.Errorf("esperaba 401 Unauthorized con token inválido, obtuvo %d", respInvalid.StatusCode)
		}
	})

	t.Run("2. 400 Bad Request cuando category no es una de las 4 permitidas", func(t *testing.T) {
		validImg := createSolidPNG(10, 10, color.RGBA{R: 50, G: 100, B: 150, A: 255})
		body, contentType, err := createMultipartBody("unauthorized_category", "test.png", validImg)
		if err != nil {
			t.Fatalf("error creando multipart: %v", err)
		}

		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/images/upload", body)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", "Bearer "+validToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("fallo petición HTTP: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("esperaba 400 Bad Request por categoría inválida, obtuvo %d", resp.StatusCode)
		}
	})

	t.Run("3. 415 Unsupported Media Type cuando se sube script disfrazado como png", func(t *testing.T) {
		fakePNG := []byte("#!/bin/bash\necho 'malicious script'\nexit 0\n")
		body, contentType, err := createMultipartBody("word", "script.png", fakePNG)
		if err != nil {
			t.Fatalf("error creando multipart: %v", err)
		}

		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/images/upload", body)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", "Bearer "+validToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("fallo petición HTTP: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnsupportedMediaType {
			t.Errorf("esperaba 415 Unsupported Media Type por fallo de magic bytes, obtuvo %d", resp.StatusCode)
		}
	})

	t.Run("4. 413 Payload Too Large cuando el archivo supera el límite de tamaño", func(t *testing.T) {
		// Límite configurado es 1 MB. Creamos payload de 1.5 MB
		largeData := make([]byte, 1500*1024)
		body, contentType, err := createMultipartBody("sense", "large.png", largeData)
		if err != nil {
			t.Fatalf("error creando multipart: %v", err)
		}

		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/images/upload", body)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", "Bearer "+validToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("fallo petición HTTP: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusRequestEntityTooLarge {
			t.Errorf("esperaba 413 Payload Too Large por exceder tamaño, obtuvo %d", resp.StatusCode)
		}
	})

	t.Run("5. 422 Unprocessable Entity cuando se sube imagen sintética NSFW", func(t *testing.T) {
		nsfwData := createNSFWImage(64, 64)
		body, contentType, err := createMultipartBody("avatar", "nsfw.png", nsfwData)
		if err != nil {
			t.Fatalf("error creando multipart: %v", err)
		}

		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/images/upload", body)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", "Bearer "+validToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("fallo petición HTTP: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("esperaba 422 Unprocessable Entity por moderación NSFW, obtuvo %d", resp.StatusCode)
		}

		var errResp response.ErrorBody
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		if !strings.Contains(errResp.Error, "NSFW") {
			t.Errorf("esperaba mención a NSFW en el motivo de rechazo, obtuvo: %s", errResp.Error)
		}
	})

	t.Run("6. 201 Created con imagen válida, metadatos eliminados, escalada y servida con nosniff", func(t *testing.T) {
		// Creamos imagen de 1000x500 para comprobar redimensionamiento proporcional a máx 800 (credentials)
		cleanPNG := createSolidPNG(1000, 500, color.RGBA{R: 40, G: 160, B: 220, A: 255})
		// Simulamos metadatos adicionales/basura al final del archivo PNG original
		dirtyPNG := append(cleanPNG, []byte("EXIF_GPS_LOCATION_LAT_40_LONG_3_METADATA_SECRET")...)

		body, contentType, err := createMultipartBody("credentials", "valid_doc.png", dirtyPNG)
		if err != nil {
			t.Fatalf("error creando multipart: %v", err)
		}

		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/images/upload", body)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", "Bearer "+validToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("fallo petición HTTP: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			tbodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("esperaba 201 Created, obtuvo %d: %s", resp.StatusCode, string(tbodyBytes))
		}

		var uploadResp models.ImageUploadResponse
		if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
			t.Fatalf("error decodificando respuesta exitosa: %v", err)
		}

		if uploadResp.Bytes <= 0 || uploadResp.Size <= 0 {
			t.Errorf("esperaba bytes y size > 0, obtuvo bytes=%d size=%d", uploadResp.Bytes, uploadResp.Size)
		}
		if uploadResp.Category != "credentials" {
			t.Errorf("categoría esperada 'credentials', obtuvo '%s'", uploadResp.Category)
		}
		if uploadResp.MIMEType != "image/png" {
			t.Errorf("MIMEType esperado 'image/png', obtuvo '%s'", uploadResp.MIMEType)
		}
		if !strings.HasPrefix(uploadResp.URL, ts.URL+"/uploads/credentials/") {
			t.Errorf("URL pública inesperada: %s", uploadResp.URL)
		}

		// Verificar que la imagen se sirve correctamente vía GET /uploads/... con X-Content-Type-Options: nosniff
		getResp, err := http.Get(uploadResp.URL)
		if err != nil {
			t.Fatalf("error realizando GET a la URL de imagen: %v", err)
		}
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusOK {
			t.Errorf("esperaba 200 OK al servir la imagen, obtuvo %d", getResp.StatusCode)
		}
		if nosniff := getResp.Header.Get("X-Content-Type-Options"); nosniff != "nosniff" {
			t.Errorf("esperaba encabezado X-Content-Type-Options: nosniff, obtuvo '%s'", nosniff)
		}

		servedBytes, err := io.ReadAll(getResp.Body)
		if err != nil {
			t.Fatalf("error leyendo imagen servida: %v", err)
		}

		// Verificar destrucción total de los metadatos simulados (RNF-SEG-04)
		if bytes.Contains(servedBytes, []byte("EXIF_GPS_LOCATION_LAT_40_LONG_3_METADATA_SECRET")) {
			t.Errorf("los metadatos simulados NO fueron eliminados durante el procesamiento")
		}

		// Verificar que el archivo servido es un PNG válido decodificable y escalado a máx 800
		servedImg, err := png.Decode(bytes.NewReader(servedBytes))
		if err != nil {
			t.Errorf("la imagen servida no es un PNG válido: %v", err)
		} else {
			bounds := servedImg.Bounds()
			if bounds.Dx() > 800 || bounds.Dy() > 800 {
				t.Errorf("la imagen servida supera 800px: dimensiones %dx%d", bounds.Dx(), bounds.Dy())
			}
		}
	})
}
