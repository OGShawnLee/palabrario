package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"palabrario/image-service/internal/api"
	"palabrario/image-service/internal/config"
)

// NewServer construye el handler HTTP a partir de parámetros explícitos para facilitar pruebas unitarias y de integración.
func NewServer(jwtSecret, uploadDir, publicBaseURL string, maxFileSizeMB int64) http.Handler {
	cfg := &config.Config{
		JWTSecret:          jwtSecret,
		JWTExpirationHours: 72 * time.Hour,
		UploadDir:          uploadDir,
		PublicBaseURL:      publicBaseURL,
		MaxFileSizeMB:      maxFileSizeMB,
	}
	router := api.NewRouter(cfg)
	return router.Setup()
}

func main() {
	// 1. Carga y validación estricta de variables de entorno (Fail-Fast)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("🚨 ERROR FATAL DE INICIALIZACIÓN: %v", err)
	}

	// 2. Inicialización de enrutador y dependencias
	router := api.NewRouter(cfg)
	handler := router.Setup()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 3. Manejo de apagado gracioso (Graceful Shutdown)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🖼️ [image-service] [%s] corriendo en http://localhost:%s", cfg.Environment, cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("❌ Error en image-service: %v", err)
		}
	}()

	<-stop
	log.Println("🛑 Apagando image-service limpiamente...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️ Error durante el apagado: %v", err)
	}

	log.Println("👋 image-service detenido.")
}
