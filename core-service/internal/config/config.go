package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config centraliza y tipa todas las variables de entorno de core-service.
// Cumple con la metodología 12-Factor App y el principio Fail-Fast:
// si alguna variable crítica no está configurada, el servicio se niega a iniciar.
type Config struct {
	Port               string
	Environment        string
	DatabaseURL        string
	JWTSecret          string
	JWTExpirationHours time.Duration
	ImageServiceURL    string
	RealtimeServiceURL string
}

// Load lee, valida rigurosamente y estructura las variables de entorno.
// Arroja un error grave con el detalle exacto de las variables faltantes
// en lugar de usar valores predeterminados inseguros en memoria.
func Load() (*Config, error) {
	var missingVars []string

	// Función auxiliar para validar variables estrictamente obligatorias
	requireEnv := func(name string) string {
		val := strings.TrimSpace(os.Getenv(name))
		if val == "" {
			missingVars = append(missingVars, name)
		}
		return val
	}

	port := requireEnv("PORT")
	dbURL := requireEnv("DATABASE_URL")
	jwtSecret := requireEnv("JWT_SECRET")

	env := strings.TrimSpace(os.Getenv("ENVIRONMENT"))
	if env == "" {
		// Validar si viene como NODE_ENV por compatibilidad previa
		env = strings.TrimSpace(os.Getenv("NODE_ENV"))
	}
	if env == "" {
		missingVars = append(missingVars, "ENVIRONMENT")
	}

	// Si hay variables obligatorias ausentes, fallar de inmediato
	if len(missingVars) > 0 {
		return nil, fmt.Errorf("error grave de configuración: las siguientes variables de entorno son obligatorias y no están definidas: [%s]", strings.Join(missingVars, ", "))
	}

	// Validación semántica de seguridad para JWT_SECRET (mínimo 32 caracteres)
	if len(jwtSecret) < 32 {
		return nil, fmt.Errorf("error de seguridad en configuración: JWT_SECRET debe tener una longitud mínima de 32 caracteres (obtenidos: %d caracteres)", len(jwtSecret))
	}

	// Validación de puerto numérico
	if _, err := strconv.Atoi(port); err != nil {
		return nil, fmt.Errorf("error de formato en configuración: PORT '%s' no es un puerto TCP válido", port)
	}

	// Expiración de JWT configurable (opcional, por defecto 72 horas)
	expHours := 72
	if expStr := strings.TrimSpace(os.Getenv("JWT_EXPIRATION_HOURS")); expStr != "" {
		if parsedExp, err := strconv.Atoi(expStr); err == nil && parsedExp > 0 {
			expHours = parsedExp
		} else {
			return nil, fmt.Errorf("error de formato en configuración: JWT_EXPIRATION_HOURS '%s' debe ser un entero positivo", expStr)
		}
	}

	return &Config{
		Port:               port,
		Environment:        env,
		DatabaseURL:        dbURL,
		JWTSecret:          jwtSecret,
		JWTExpirationHours: time.Duration(expHours) * time.Hour,
		ImageServiceURL:    strings.TrimSpace(os.Getenv("IMAGE_SERVICE_URL")),
		RealtimeServiceURL: strings.TrimSpace(os.Getenv("REALTIME_SERVICE_URL")),
	}, nil
}
