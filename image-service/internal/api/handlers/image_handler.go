package handlers

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"palabrario/image-service/internal/processor"
	"palabrario/image-service/internal/safety"
	"palabrario/image-service/internal/storage"
	"palabrario/pkg/models"
	"palabrario/pkg/response"
)

// ImageHandler encapsula la lógica HTTP para la carga y gestión de imágenes.
type ImageHandler struct {
	processor     processor.ImageProcessor
	validator     safety.SafetyValidator
	storage       storage.StorageProvider
	maxFileSizeMB int64
}

// NewImageHandler inicializa un nuevo ImageHandler con sus dependencias inyectadas.
func NewImageHandler(
	proc processor.ImageProcessor,
	val safety.SafetyValidator,
	store storage.StorageProvider,
	maxFileSizeMB int64,
) *ImageHandler {
	if maxFileSizeMB <= 0 {
		maxFileSizeMB = 5
	}
	return &ImageHandler{
		processor:     proc,
		validator:     val,
		storage:       store,
		maxFileSizeMB: maxFileSizeMB,
	}
}

// Upload maneja la subida, validación de formato/NSFW, saneamiento y persistencia en disco.
func (h *ImageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	maxBytes := h.maxFileSizeMB * 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	if err := r.ParseMultipartForm(maxBytes); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) || strings.Contains(strings.ToLower(err.Error()), "too large") {
			response.Error(w, http.StatusRequestEntityTooLarge, "el archivo supera el tamaño máximo permitido")
			return
		}
		response.Error(w, http.StatusBadRequest, "error al procesar formulario multipart: "+err.Error())
		return
	}

	category := strings.TrimSpace(r.FormValue("category"))
	switch category {
	case "avatar", "credentials", "word", "sense":
	default:
		response.Error(w, http.StatusBadRequest, "categoría inválida; debe ser estrictamente 'avatar', 'credentials', 'word' o 'sense'")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) || strings.Contains(strings.ToLower(err.Error()), "too large") {
			response.Error(w, http.StatusRequestEntityTooLarge, "el archivo supera el tamaño máximo permitido")
			return
		}
		response.Error(w, http.StatusBadRequest, "campo 'file' no encontrado o inválido: "+err.Error())
		return
	}
	defer file.Close()

	if header.Size > maxBytes {
		response.Error(w, http.StatusRequestEntityTooLarge, "el archivo supera el tamaño máximo permitido")
		return
	}

	data, err := io.ReadAll(io.LimitReader(file, maxBytes+1))
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) || strings.Contains(strings.ToLower(err.Error()), "too large") {
			response.Error(w, http.StatusRequestEntityTooLarge, "el archivo supera el tamaño máximo permitido")
			return
		}
		response.Error(w, http.StatusBadRequest, "error al leer los bytes del archivo: "+err.Error())
		return
	}

	if int64(len(data)) > maxBytes {
		response.Error(w, http.StatusRequestEntityTooLarge, "el archivo supera el tamaño máximo permitido")
		return
	}
	if len(data) == 0 {
		response.Error(w, http.StatusBadRequest, "el archivo proporcionado está vacío")
		return
	}

	processed, err := h.processor.ProcessWithCategory(data, category)
	if err != nil {
		if errors.Is(err, processor.ErrUnsupportedMimeType) {
			response.Error(w, http.StatusUnsupportedMediaType, "formato multimedia no soportado; se admiten únicamente image/jpeg, image/png o image/webp")
			return
		}
		response.Error(w, http.StatusBadRequest, "error en procesamiento de imagen: "+err.Error())
		return
	}

	safe, reason, err := h.validator.Validate(r.Context(), processed.Image)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "error al moderar imagen: "+err.Error())
		return
	}
	if !safe {
		response.Error(w, http.StatusUnprocessableEntity, reason)
		return
	}

	publicURL, err := h.storage.Save(r.Context(), category, processed.Extension, processed.Data)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "error al almacenar imagen en disco: "+err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, models.ImageUploadResponse{
		URL:      publicURL,
		Category: category,
		MIMEType: processed.MIMEType,
		Size:     int64(len(processed.Data)),
		Bytes:    len(processed.Data),
	})
}
