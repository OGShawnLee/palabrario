package processor

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"
)

func TestStandardProcessor(t *testing.T) {
	proc := NewStandardProcessor()

	t.Run("empty payload returns error", func(t *testing.T) {
		_, err := proc.Process(nil)
		if err == nil {
			t.Errorf("esperaba error para payload vacío")
		}
	})

	t.Run("non-image payload returns unsupported MIME error", func(t *testing.T) {
		_, err := proc.Process([]byte("hello world plain text"))
		if err == nil {
			t.Errorf("esperaba error para texto plano")
		}
	})

	t.Run("valid PNG is processed and re-encoded", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 20, 20))
		img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		var buf bytes.Buffer
		_ = png.Encode(&buf, img)

		processed, err := proc.Process(buf.Bytes())
		if err != nil {
			t.Fatalf("error procesando PNG válido: %v", err)
		}
		if processed.Extension != "png" {
			t.Errorf("extensión esperada 'png', obtuvo %s", processed.Extension)
		}
		if processed.MIMEType != "image/png" {
			t.Errorf("MIME esperado 'image/png', obtuvo %s", processed.MIMEType)
		}
		if len(processed.Data) == 0 {
			t.Errorf("los bytes procesados no deben estar vacíos")
		}

		// Debe poder decodificarse nuevamente
		_, err = png.Decode(bytes.NewReader(processed.Data))
		if err != nil {
			t.Errorf("la imagen re-codificada no es PNG válido: %v", err)
		}
	})

	t.Run("valid JPEG is processed and re-encoded", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 20, 20))
		img.Set(0, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})
		var buf bytes.Buffer
		_ = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})

		processed, err := proc.Process(buf.Bytes())
		if err != nil {
			t.Fatalf("error procesando JPEG válido: %v", err)
		}
		if processed.Extension != "jpg" {
			t.Errorf("extensión esperada 'jpg', obtuvo %s", processed.Extension)
		}
		if processed.MIMEType != "image/jpeg" {
			t.Errorf("MIME esperado 'image/jpeg', obtuvo %s", processed.MIMEType)
		}
		if len(processed.Data) == 0 {
			t.Errorf("los bytes procesados no deben estar vacíos")
		}

		// Debe poder decodificarse nuevamente
		_, err = jpeg.Decode(bytes.NewReader(processed.Data))
		if err != nil {
			t.Errorf("la imagen re-codificada no es JPEG válido: %v", err)
		}
	})
}
