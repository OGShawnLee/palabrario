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

	"palabrario/core-service/internal/api"
	"palabrario/core-service/internal/config"
	"palabrario/core-service/internal/database"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("🔌 Conectando a PostgreSQL en %s...", cfg.DatabaseURL)
	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Error crítico de base de datos: %v", err)
	}
	defer db.Close()
	log.Println("✅ Conexión con PostgreSQL 17 establecida correctamente.")

	router := api.NewRouter(cfg, db)
	handler := router.Setup()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Manejo de apagado gracioso (Graceful Shutdown)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 [core-service] corriendo en http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("❌ Error en servidor HTTP: %v", err)
		}
	}()

	<-stop
	log.Println("🛑 Apagando core-service limpiamente...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️ Error durante el apagado: %v", err)
	}

	log.Println("👋 core-service detenido.")
}
