package models

// ImageUploadResponse representa el contrato de respuesta estandarizado
// al subir un artefacto multimedia en image-service (CU01, CU02, CU06, RNF-SEG-04).
type ImageUploadResponse struct {
	URL      string `json:"url"`
	Category string `json:"category"`
	MIMEType string `json:"mime_type"`
	Size     int64  `json:"size"`
	Bytes    int    `json:"bytes,omitempty"` // Alias numérico de compatibilidad
}
