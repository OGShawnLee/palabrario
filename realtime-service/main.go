package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"palabrario/pkg/auth"
	"palabrario/pkg/response"
	"palabrario/realtime-service/internal/arena"
)

type HealthResponse struct {
	Status      string    `json:"status"`
	Service     string    `json:"service"`
	ActiveRooms int       `json:"active_rooms"`
	Timestamp   time.Time `json:"timestamp"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-key-change-in-prod"
	}

	coreServiceURL := os.Getenv("CORE_SERVICE_URL")
	if coreServiceURL == "" {
		coreServiceURL = "http://localhost:8000"
	}

	tokenService := auth.NewTokenService(jwtSecret, 72*time.Hour)
	coreClient := arena.NewCoreClient(coreServiceURL, nil, tokenService)
	hub := arena.NewArenaHub(coreClient, tokenService)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		resp := HealthResponse{
			Status:      "ok",
			Service:     "realtime-service",
			ActiveRooms: hub.ActiveRoomsCount(),
			Timestamp:   time.Now().UTC(),
		}
		_ = json.NewEncoder(w).Encode(resp)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"message":"Bienvenido a Palabrario Realtime Service (Golang)"}`))
	})

	r.Get("/ws/arena/{roomCode}", func(w http.ResponseWriter, r *http.Request) {
		roomCode := strings.TrimSpace(chi.URLParam(r, "roomCode"))
		if roomCode == "" {
			response.Error(w, http.StatusBadRequest, "El código de sala es requerido")
			return
		}

		tokenStr := strings.TrimSpace(r.URL.Query().Get("token"))
		if tokenStr == "" {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenStr == "" {
			response.Error(w, http.StatusUnauthorized, "Token JWT requerido en ?token= o Authorization")
			return
		}

		claims, err := tokenService.VerifyToken(tokenStr)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Token inválido o expirado")
			return
		}

		hub.HandleWS(w, r, roomCode, claims)
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("⚡ [realtime-service] corriendo en http://localhost:%s", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("❌ Error en realtime-service: %v", err)
		}
	}()

	<-stop
	log.Println("🛑 Apagando realtime-service...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️ Error durante el apagado: %v", err)
	}

	log.Println("👋 realtime-service detenido.")
}
