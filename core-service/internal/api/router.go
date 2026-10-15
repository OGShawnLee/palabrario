package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"palabrario/core-service/internal/api/handlers"
	"palabrario/core-service/internal/config"
	"palabrario/core-service/internal/database"
)

type Router struct {
	cfg *config.Config
	db  *database.DB
}

func NewRouter(cfg *config.Config, db *database.DB) *Router {
	return &Router{cfg: cfg, db: db}
}

func (rt *Router) Setup() http.Handler {
	r := chi.NewRouter()

	// Middlewares esenciales
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS permisivo para desarrollo y Wails/Frontend
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	healthHandler := handlers.NewHealthHandler(rt.db)
	r.Get("/health", healthHandler.HealthCheck)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"Bienvenido a Palabrario Core API (Golang + PostgreSQL 17)"}`))
	})

	return r
}
