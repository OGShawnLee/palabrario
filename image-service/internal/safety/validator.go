package safety

import (
	"context"
	"errors"
	"image"
	"image/color"
	"math"
)

// SafetyValidator define la estrategia para evaluar si una imagen es apta para la plataforma.
type SafetyValidator interface {
	Validate(ctx context.Context, img image.Image) (safe bool, reason string, err error)
}

// HeuristicSkinEntropyValidator implementa SafetyValidator en 100% Go puro.
// Analiza la proporción de tonos de piel en espacios YCbCr y RGB/HSV junto a la entropía de Shannon
// de la luminancia para descartar fotos con alto contenido NSFW sin falsos positivos en íconos o fondos planos.
type HeuristicSkinEntropyValidator struct{}

// NewHeuristicSkinEntropyValidator crea una nueva instancia del validador heurístico.
func NewHeuristicSkinEntropyValidator() *HeuristicSkinEntropyValidator {
	return &HeuristicSkinEntropyValidator{}
}

func (v *HeuristicSkinEntropyValidator) Validate(ctx context.Context, img image.Image) (bool, string, error) {
	if img == nil {
		return false, "imagen inválida", errors.New("la imagen proporcionada es nil")
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	if w <= 0 || h <= 0 {
		return false, "dimensiones de imagen inválidas", errors.New("dimensiones de imagen vacías o inválidas")
	}

	// Cuadrícula de muestreo de hasta 128x128
	sampleW := w
	if sampleW > 128 {
		sampleW = 128
	}
	sampleH := h
	if sampleH > 128 {
		sampleH = 128
	}

	totalPixels := sampleW * sampleH
	if totalPixels == 0 {
		return true, "", nil
	}

	skinPixels := 0
	var hist [32]int

	for iy := 0; iy < sampleH; iy++ {
		yCoord := bounds.Min.Y + (iy * h / sampleH)
		for ix := 0; ix < sampleW; ix++ {
			xCoord := bounds.Min.X + (ix * w / sampleW)

			r16, g16, b16, _ := img.At(xCoord, yCoord).RGBA()
			r := uint8(r16 >> 8)
			g := uint8(g16 >> 8)
			b := uint8(b16 >> 8)

			y, cb, cr := color.RGBToYCbCr(r, g, b)

			// Histograma de 32 bins sobre luminancia Y (0..255 dividido entre 8)
			bin := int(y) / 8
			if bin > 31 {
				bin = 31
			}
			hist[bin]++

			// a) Espacio YCbCr: tono dérmico
			ycbcrSkin := y > 40 && cb >= 77 && cb <= 127 && cr >= 133 && cr <= 173

			// b) Espacio HSV/RGB adicional
			maxVal := r
			if g > maxVal {
				maxVal = g
			}
			if b > maxVal {
				maxVal = b
			}
			minVal := r
			if g < minVal {
				minVal = g
			}
			if b < minVal {
				minVal = b
			}

			diffRG := int(r) - int(g)
			if diffRG < 0 {
				diffRG = -diffRG
			}

			rgbSkin := r > 95 && g > 40 && b > 20 && (int(maxVal)-int(minVal)) > 15 && diffRG > 15 && r > g && r > b

			if ycbcrSkin && rgbSkin {
				skinPixels++
			}
		}
	}

	// Entropía de Shannon de la luminancia Y
	var entropy float64
	for _, count := range hist {
		if count > 0 {
			p := float64(count) / float64(totalPixels)
			entropy -= p * math.Log2(p)
		}
	}

	skinRatio := float64(skinPixels) / float64(totalPixels)

	// Regla de descarte NSFW
	if skinRatio > 0.45 && entropy > 2.5 {
		return false, "imagen rechazada por el filtro de moderación visual (NSFW)", nil
	}

	return true, "", nil
}
