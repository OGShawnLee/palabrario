package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"palabrario/pkg/auth"
	"palabrario/realtime-service/internal/api/handlers"
	"palabrario/realtime-service/internal/config"
	"palabrario/realtime-service/internal/service"
)

// Router centraliza la configuración de rutas, middlewares y handlers de realtime-service.
type Router struct {
	cfg          *config.Config
	hub          *service.ArenaHub
	tokenService *auth.TokenService
}

// NewRouter instancia un nuevo Router con la configuración y servicios correspondientes.
func NewRouter(cfg *config.Config, hub *service.ArenaHub, tokenService *auth.TokenService) *Router {
	return &Router{
		cfg:          cfg,
		hub:          hub,
		tokenService: tokenService,
	}
}

// Setup configura los middlewares globales y registra todas las rutas del servicio.
func (rt *Router) Setup() http.Handler {
	r := chi.NewRouter()

	// Middlewares estándar de observabilidad y resiliencia
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	allowedOrigins := rt.cfg.AllowedOrigins
	if len(allowedOrigins) == 0 {
		allowedOrigins = auth.ParseCORSAllowedOrigins("")
	}
	rt.hub.SetAllowedOrigins(allowedOrigins)

	// CORS configured from environment allowlist without credentials (Bearer auth only)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	healthHandler := handlers.NewHealthHandler(rt.hub)
	arenaHandler := handlers.NewArenaHandler(rt.hub, rt.tokenService)

	// Probe de salud operacional y bienvenida
	r.Get("/health", healthHandler.HealthCheck)
	r.Get("/", healthHandler.Root)

	// Conexión WebSocket bidireccional para la Arena
	r.Get("/ws/arena/{roomCode}", arenaHandler.HandleWS)

	return r
}
