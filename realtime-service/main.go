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

	"palabrario/pkg/auth"
	"palabrario/realtime-service/internal/api"
	"palabrario/realtime-service/internal/config"
	"palabrario/realtime-service/internal/service"
)

func main() {
	// 1. Carga y validación estricta de variables de entorno (Fail-Fast / 12-Factor App)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("🚨 ERROR FATAL DE INICIALIZACIÓN: %v", err)
	}

	// 2. Inicialización de servicios y clientes de dominio
	tokenService := auth.NewTokenService(cfg.JWTSecret, cfg.JWTExpirationHours)
	coreClient := service.NewCoreClient(cfg.CoreServiceURL, nil, tokenService)
	hub := service.NewArenaHub(coreClient, tokenService)

	// 3. Inicialización del enrutador HTTP y WebSockets
	router := api.NewRouter(cfg, hub, tokenService)
	handler := router.Setup()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 4. Manejo de apagado gracioso (Graceful Shutdown)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("⚡ [realtime-service] [%s] corriendo en http://localhost:%s", cfg.Environment, cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("❌ Error en realtime-service: %v", err)
		}
	}()

	<-stop
	log.Println("🛑 Apagando realtime-service limpiamente...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️ Error durante el apagado: %v", err)
	}

	log.Println("👋 realtime-service detenido.")
}
