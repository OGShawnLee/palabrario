package safety

import (
	"context"
	"image"
	"image/color"
	"testing"
)

func TestHeuristicSkinEntropyValidator(t *testing.T) {
	v := NewHeuristicSkinEntropyValidator()
	ctx := context.Background()

	t.Run("nil image returns error", func(t *testing.T) {
		safe, _, err := v.Validate(ctx, nil)
		if err == nil || safe {
			t.Errorf("esperaba error para imagen nil, safe=%v err=%v", safe, err)
		}
	})

	t.Run("solid peach/cream background is safe because entropy is low", func(t *testing.T) {
		// Píxeles de tono piel pero totalmente uniformes (entropía = 0)
		img := image.NewRGBA(image.Rect(0, 0, 50, 50))
		r, g, b := color.YCbCrToRGB(150, 105, 150)
		for y := 0; y < 50; y++ {
			for x := 0; x < 50; x++ {
				img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
			}
		}

		safe, reason, err := v.Validate(ctx, img)
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if !safe {
			t.Errorf("un fondo sólido o ícono plano no debe rechazarse como NSFW, motivo: %s", reason)
		}
	})

	t.Run("landscape / non-skin image is safe", func(t *testing.T) {
		// Imagen azul / verde (cielo / campo)
		img := image.NewRGBA(image.Rect(0, 0, 50, 50))
		for y := 0; y < 50; y++ {
			for x := 0; x < 50; x++ {
				img.Set(x, y, color.RGBA{R: 20, G: 120, B: 200, A: 255})
			}
		}

		safe, _, err := v.Validate(ctx, img)
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if !safe {
			t.Errorf("una imagen sin tonos dérmicos debe considerarse safe")
		}
	})

	t.Run("synthetic image with high skin ratio and high entropy is rejected as NSFW", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 64, 64))
		luminances := []uint8{50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160, 170, 180, 190, 200}
		for y := 0; y < 64; y++ {
			lum := luminances[y%len(luminances)]
			r, g, b := color.YCbCrToRGB(lum, 105, 150)
			for x := 0; x < 64; x++ {
				img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
			}
		}

		safe, reason, err := v.Validate(ctx, img)
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if safe {
			t.Errorf("esperaba rechazo NSFW para imagen con alto skinRatio y entropía")
		}
		if reason != "imagen rechazada por el filtro de moderación visual (NSFW)" {
			t.Errorf("motivo inesperado: %s", reason)
		}
	})
}
