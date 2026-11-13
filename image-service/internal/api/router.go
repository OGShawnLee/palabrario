package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"palabrario/image-service/internal/api/handlers"
	"palabrario/image-service/internal/api/middleware"
	"palabrario/image-service/internal/config"
	"palabrario/image-service/internal/processor"
	"palabrario/image-service/internal/safety"
	"palabrario/image-service/internal/storage"
	"palabrario/pkg/auth"
)

// HealthResponse estructura la respuesta del probe de salud.
type HealthResponse struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
}

// Router centraliza la construcción del enrutador HTTP y sus dependencias.
type Router struct {
	cfg       *config.Config
	processor processor.ImageProcessor
	validator safety.SafetyValidator
	storage   storage.StorageProvider
}

// NewRouter crea una nueva instancia del enrutador con las implementaciones estándar.
func NewRouter(cfg *config.Config) *Router {
	proc := processor.NewStandardProcessor()
	val := safety.NewHeuristicSkinEntropyValidator()
	store := storage.NewLocalDiskStorage(cfg.UploadDir, cfg.PublicBaseURL)

	return &Router{
		cfg:       cfg,
		processor: proc,
		validator: val,
		storage:   store,
	}
}

// NewRouterWithDependencies permite inyectar dependencias personalizadas o mocks para pruebas.
func NewRouterWithDependencies(
	cfg *config.Config,
	proc processor.ImageProcessor,
	val safety.SafetyValidator,
	store storage.StorageProvider,
) *Router {
	return &Router{
		cfg:       cfg,
		processor: proc,
		validator: val,
		storage:   store,
	}
}

// Setup configura los middlewares globales, rutas de salud, servidor de archivos y endpoints protegidos.
func (rt *Router) Setup() http.Handler {
	// Asegurar existencia del directorio de uploads
	_ = os.MkdirAll(rt.cfg.UploadDir, 0755)

	ttl := rt.cfg.AccessTokenTTL
	if ttl <= 0 {
		ttl = rt.cfg.JWTExpirationHours
	}
	tokenService := auth.NewTokenService(rt.cfg.JWTSecret, ttl)
	authMiddleware := middleware.NewAuthMiddleware(tokenService)
	imageHandler := handlers.NewImageHandler(rt.processor, rt.validator, rt.storage, rt.cfg.MaxFileSizeMB)

	allowedOrigins := rt.cfg.AllowedOrigins
	if len(allowedOrigins) == 0 {
		allowedOrigins = auth.ParseCORSAllowedOrigins("")
	}

	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := HealthResponse{
			Status:    "ok",
			Service:   "image-service",
			Timestamp: time.Now().UTC(),
		}
		_ = json.NewEncoder(w).Encode(resp)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"Bienvenido a Palabrario Image Service (Golang)"}`))
	})

	// Servidor de archivos estáticos para descargas de imágenes saneadas (RNF-SEG-04)
	// Inyecta 'X-Content-Type-Options: nosniff' para prevenir MIME-sniffing en clientes nativos/web
	fileServer := http.StripPrefix("/uploads/", http.FileServer(http.Dir(rt.cfg.UploadDir)))
	r.Handle("/uploads/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		fileServer.ServeHTTP(w, r)
	}))

	// Endpoints protegidos de carga y gestión de imágenes
	r.Route("/api/v1/images", func(r chi.Router) {
		r.With(authMiddleware.RequireAuth).Post("/upload", imageHandler.Upload)
	})

	return r
}
