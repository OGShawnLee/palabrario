package service

import (
	"math"
)

// CalculateRoundPoints calcula los puntos obtenidos en una ronda de Arena según la regla RN10
// de la sección 2.7.4 de INTERFACES.md:
//   - Si !isCorrect o elapsedMs > roundTimeSec * 1000 => 0 puntos.
//   - Si elapsedMs < 0 => se ajusta a 0 ms.
//   - Sea tMax = float64(roundTimeSec * 1000) y clamped = math.Min(float64(elapsedMs), tMax).
//   - Puntos = int(math.Floor(500.0 + 500.0 * (1.0 - clamped / tMax))) acotado en el intervalo [500, 1000].
func CalculateRoundPoints(isCorrect bool, elapsedMs int64, roundTimeSec int) int {
	if !isCorrect || roundTimeSec <= 0 {
		return 0
	}

	tMaxMs := int64(roundTimeSec) * 1000
	if elapsedMs > tMaxMs {
		return 0
	}

	if elapsedMs < 0 {
		elapsedMs = 0
	}

	tMax := float64(tMaxMs)
	clamped := math.Min(float64(elapsedMs), tMax)

	pts := int(math.Floor(500.0 + 500.0*(1.0-clamped/tMax)))
	if pts < 500 {
		pts = 500
	}
	if pts > 1000 {
		pts = 1000
	}

	return pts
}
