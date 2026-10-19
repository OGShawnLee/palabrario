package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"palabrario/core-service/internal/api/handlers"
	"palabrario/core-service/internal/api/middleware"
	"palabrario/core-service/internal/config"
	"palabrario/core-service/internal/database"
	"palabrario/core-service/internal/service"
	"palabrario/pkg/auth"
	"palabrario/pkg/response"
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
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// CORS permisivo para desarrollo y Wails/Frontend
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// 1. Capa de Servicios SOA (Contratos de .context/INTERFACES.md)
	tokenService := auth.NewTokenService(rt.cfg.JWTSecret, rt.cfg.JWTExpirationHours)
	identityService := service.NewIdentityService(rt.db, tokenService)
	dictionaryService := service.NewDictionaryService(rt.db)

	// 2. Capa de Middlewares y Handlers HTTP
	authMiddleware := middleware.NewAuthMiddleware(tokenService)
	authHandler := handlers.NewAuthHandler(identityService)
	dictionaryHandler := handlers.NewDictionaryHandler(dictionaryService)
	healthHandler := handlers.NewHealthHandler(rt.db)

	// Health check raíz
	r.Get("/health", healthHandler.HealthCheck)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{
			"message": "Bienvenido a Palabrario Core API (Golang + PostgreSQL 17)",
		})
	})

	// Rutas API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Módulo 1: Autenticación e Identidad (CU01)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", authHandler.SignUp)
			r.Post("/sign-in", authHandler.SignIn)

			// Rutas protegidas de Auth
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Get("/me", authHandler.GetMe)
			})
		})

		// Módulo 3: Diccionario y Búsqueda Dialectal (CU05)
		r.Route("/words", func(r chi.Router) {
			r.With(authMiddleware.OptionalAuth).Get("/search", dictionaryHandler.Search)
			r.Get("/{id}", dictionaryHandler.GetByID)
		})
	})

	return r
}
