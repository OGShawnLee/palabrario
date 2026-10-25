package processor

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"

	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

var (
	// ErrUnsupportedMimeType se retorna cuando los magic bytes no corresponden a formatos permitidos.
	ErrUnsupportedMimeType = errors.New("tipo de archivo no soportado; se admiten únicamente image/jpeg, image/png o image/webp")
	// ErrEmptyData se retorna cuando el payload está vacío.
	ErrEmptyData = errors.New("los datos del archivo están vacíos")
	// ErrDecodeFailed se retorna cuando falla la decodificación del bitmap en RAM.
	ErrDecodeFailed = errors.New("no fue posible decodificar la imagen")
)

// ProcessedImage contiene el bitmap decodificado en memoria y los bytes limpios listos para persistir.
type ProcessedImage struct {
	Image     image.Image
	Data      []byte
	Extension string
	MIMEType  string
}

// ImageProcessor define la interfaz para el procesamiento, escalado y saneamiento de imágenes.
type ImageProcessor interface {
	Process(data []byte) (*ProcessedImage, error)
	ProcessWithCategory(data []byte, category string) (*ProcessedImage, error)
}

// StandardProcessor implementa ImageProcessor validando magic bytes, redimensionando por bounding-box
// y re-codificando el bitmap puro desde cero para eliminación total de metadatos EXIF (RNF-SEG-04).
type StandardProcessor struct{}

// NewStandardProcessor crea un nuevo procesador estándar.
func NewStandardProcessor() *StandardProcessor {
	return &StandardProcessor{}
}

func (p *StandardProcessor) Process(data []byte) (*ProcessedImage, error) {
	return p.ProcessWithCategory(data, "default")
}

func (p *StandardProcessor) ProcessWithCategory(data []byte, category string) (*ProcessedImage, error) {
	if len(data) == 0 {
		return nil, ErrEmptyData
	}

	headerLen := len(data)
	if headerLen > 512 {
		headerLen = 512
	}

	detectedMime := http.DetectContentType(data[:headerLen])
	if detectedMime != "image/jpeg" && detectedMime != "image/png" && detectedMime != "image/webp" {
		return nil, fmt.Errorf("%w: recibido %s", ErrUnsupportedMimeType, detectedMime)
	}

	var img image.Image
	var err error

	switch detectedMime {
	case "image/jpeg":
		img, err = jpeg.Decode(bytes.NewReader(data))
	case "image/png":
		img, err = png.Decode(bytes.NewReader(data))
	case "image/webp":
		img, err = webp.Decode(bytes.NewReader(data))
	default:
		return nil, ErrUnsupportedMimeType
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecodeFailed, err)
	}

	// Redimensionado restrictivo según categoría (RNF-SEG-04):
	// - avatar: máx 400x400
	// - credentials, word, sense y general: máx 800x800
	maxDim := 800
	if category == "avatar" {
		maxDim = 400
	}

	resizedImg := resizeIfNeeded(img, maxDim)

	// Re-codificación desde el bitmap puro para destrucción total de metadatos EXIF/GPS/IPTC (RNF-SEG-04)
	var buf bytes.Buffer
	var ext string
	var outMime string

	if detectedMime == "image/png" {
		ext = "png"
		outMime = "image/png"
		if err := png.Encode(&buf, resizedImg); err != nil {
			return nil, fmt.Errorf("error al codificar PNG: %w", err)
		}
	} else {
		// JPEG y WebP se normalizan a JPEG con calidad 85
		ext = "jpg"
		outMime = "image/jpeg"
		if err := jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 85}); err != nil {
			return nil, fmt.Errorf("error al codificar JPEG: %w", err)
		}
	}

	return &ProcessedImage{
		Image:     resizedImg,
		Data:      buf.Bytes(),
		Extension: ext,
		MIMEType:  outMime,
	}, nil
}

// resizeIfNeeded escala proporcionalmente la imagen si excede maxDim en ancho o alto usando Go puro.
func resizeIfNeeded(src image.Image, maxDim int) image.Image {
	bounds := src.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w <= maxDim && h <= maxDim {
		return src
	}

	var newW, newH int
	if w > h {
		newW = maxDim
		newH = (h * maxDim) / w
	} else {
		newH = maxDim
		newW = (w * maxDim) / h
	}
	if newW < 1 {
		newW = 1
	}
	if newH < 1 {
		newH = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)
	return dst
}
