package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config centraliza y tipa todas las variables de entorno de realtime-service.
// Sigue la metodología 12-Factor App con validación estructurada y principio Fail-Fast.
type Config struct {
	Port               string
	Environment        string
	JWTSecret          string
	JWTExpirationHours time.Duration
	CoreServiceURL     string
}

// Load lee, valida rigurosamente y estructura las variables de entorno de realtime-service.
func Load() (*Config, error) {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8001"
	}
	if _, err := strconv.Atoi(port); err != nil {
		return nil, fmt.Errorf("error de formato en configuración: PORT '%s' no es un puerto TCP válido", port)
	}

	env := strings.TrimSpace(os.Getenv("ENVIRONMENT"))
	if env == "" {
		env = strings.TrimSpace(os.Getenv("NODE_ENV"))
	}
	if env == "" {
		env = "development"
	}

	jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if jwtSecret == "" {
		jwtSecret = "super_secret_jwt_key_palabrario_change_in_production_min_32_chars"
	}
	if len(jwtSecret) < 32 {
		return nil, fmt.Errorf("error de seguridad en configuración: JWT_SECRET debe tener una longitud mínima de 32 caracteres (obtenidos: %d caracteres)", len(jwtSecret))
	}

	coreServiceURL := strings.TrimSpace(os.Getenv("CORE_SERVICE_URL"))
	if coreServiceURL == "" {
		coreServiceURL = "http://localhost:8000"
	}

	expHours := 72
	if expStr := strings.TrimSpace(os.Getenv("JWT_EXPIRATION_HOURS")); expStr != "" {
		if val, err := strconv.Atoi(expStr); err == nil && val > 0 {
			expHours = val
		} else {
			return nil, fmt.Errorf("error de formato en configuración: JWT_EXPIRATION_HOURS '%s' debe ser un entero positivo", expStr)
		}
	}

	return &Config{
		Port:               port,
		Environment:        env,
		JWTSecret:          jwtSecret,
		JWTExpirationHours: time.Duration(expHours) * time.Hour,
		CoreServiceURL:     coreServiceURL,
	}, nil
}
