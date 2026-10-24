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
	"palabrario/pkg/models"
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
		ExposedHeaders:   []string{"Link", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// 1. Capa de Servicios SOA (Contratos de .context/INTERFACES.md)
	tokenService := auth.NewTokenService(rt.cfg.JWTSecret, rt.cfg.JWTExpirationHours)
	identityService := service.NewIdentityService(rt.db, tokenService)
	governanceService := service.NewGovernanceService(rt.db)
	dictionaryService := service.NewDictionaryService(rt.db)
	curationService := service.NewCurationService(rt.db)
	studioService := service.NewStudioService(rt.db)
	masteryService := service.NewMasteryService(rt.db)

	// 2. Capa de Middlewares y Handlers HTTP
	authMiddleware := middleware.NewAuthMiddleware(tokenService)
	authHandler := handlers.NewAuthHandler(identityService)
	userHandler := handlers.NewUserHandler(identityService)
	governanceHandler := handlers.NewGovernanceHandler(governanceService)
	dictionaryHandler := handlers.NewDictionaryHandler(dictionaryService)
	curationHandler := handlers.NewCurationHandler(curationService)
	studioHandler := handlers.NewStudioHandler(studioService)
	masteryHandler := handlers.NewMasteryHandler(masteryService)
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

		// Módulo 1: Perfiles de Usuario (CU01)
		r.Route("/users", func(r chi.Router) {
			r.With(authMiddleware.OptionalAuth).Get("/u/{username}", userHandler.GetPublicProfile)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Put("/me", userHandler.UpdateProfile)
				r.Put("/me/password", userHandler.ChangePassword)
			})
		})

		// Módulo 2: Gobernanza Institucional, Roles y Boletos (CU02, CU03, CU04)
		r.Route("/governance", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)

				// CU02: Solicitudes de ascenso a Editor (Aprendiz)
				r.Post("/editor-requests", governanceHandler.CreateEditorRequest)
				r.Get("/editor-requests/me", governanceHandler.ListMyEditorRequests)

				// CU04: Redención de boleto
				r.Post("/tickets/redeem", governanceHandler.RedeemInvitationTicket)

				// CU03 & CU04: Administración Institucional (Admin)
				r.Group(func(r chi.Router) {
					r.Use(authMiddleware.RequireRole(models.RoleAdmin))

					r.Get("/editor-requests", governanceHandler.ListEditorRequests)
					r.Put("/editor-requests/{id}/resolve", governanceHandler.ResolveEditorRequest)

					r.Get("/users", governanceHandler.ListUsers)
					r.Put("/users/{id}/access", governanceHandler.UpdateUserRoleOrStatus)

					r.Post("/tickets", governanceHandler.CreateInvitationTicket)
					r.Get("/tickets", governanceHandler.ListInvitationTickets)
					r.Delete("/tickets/{id}", governanceHandler.RevokeInvitationTicket)
				})
			})
		})

		// Módulo 3: Regiones Dialectales (CU05)
		r.Get("/regions", dictionaryHandler.ListRegions)

		// Módulo 3: Palabras, Exploración y Búsqueda Dialectal (CU05, CU06)
		r.Route("/words", func(r chi.Router) {
			r.With(authMiddleware.OptionalAuth).Get("/search", dictionaryHandler.Search)
			r.With(authMiddleware.OptionalAuth).Get("/word-of-the-day", dictionaryHandler.GetWordOfTheDay)
			r.Get("/", dictionaryHandler.ExploreWords)
			r.Get("/{id}", dictionaryHandler.GetByID)

			// Autoría editorial de palabras (CU06)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Use(authMiddleware.RequireRole(models.RoleEditor, models.RoleAdmin))
				r.Post("/", dictionaryHandler.CreateWord)
				r.Put("/{id}", dictionaryHandler.UpdateWord)
				r.Post("/{id}/senses", dictionaryHandler.CreateSense)
			})
		})

		// Módulo 3: Acepciones y Ejemplos (CU06)
		r.Route("/senses", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Use(authMiddleware.RequireRole(models.RoleEditor, models.RoleAdmin))
				r.Put("/{senseId}", dictionaryHandler.UpdateSense)
				r.Delete("/{senseId}", dictionaryHandler.DeleteSense)
				r.Post("/{senseId}/examples", dictionaryHandler.CreateExample)
				r.Post("/{senseId}/concepts/{conceptId}", dictionaryHandler.LinkConcept)
				r.Delete("/{senseId}/concepts/{conceptId}", dictionaryHandler.UnlinkConcept)
			})
		})

		r.Route("/examples", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Use(authMiddleware.RequireRole(models.RoleEditor, models.RoleAdmin))
				r.Put("/{exampleId}", dictionaryHandler.UpdateExample)
				r.Delete("/{exampleId}", dictionaryHandler.DeleteExample)
			})
		})

		// Módulo 3: Heteronimia y Conceptos Semánticos (CU05, CU06)
		r.Route("/concepts", func(r chi.Router) {
			r.Get("/", dictionaryHandler.ListConcepts)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Use(authMiddleware.RequireRole(models.RoleEditor, models.RoleAdmin))
				r.Post("/", dictionaryHandler.CreateConcept)
			})
		})

		// Módulo 4: Curaduría Editorial por Quórum (CU07, CU08)
		r.Route("/curation", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Use(authMiddleware.RequireRole(models.RoleEditor, models.RoleAdmin))

				r.Get("/queue", curationHandler.ListPendingQueue)
				r.Get("/my-proposals", curationHandler.ListMyProposals)
				r.Get("/senses/{senseId}", curationHandler.GetProposalDetail)
				r.Post("/senses/{senseId}/vote", curationHandler.SubmitVote)
			})

			// Veto o Aprobación Directa de Administrador (RN08)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Use(authMiddleware.RequireRole(models.RoleAdmin))
				r.Put("/senses/{senseId}/state", curationHandler.AdminOverrideState)
			})
		})

		// Módulo 4: Moderación de Reportes Comunitarios (CU09)
		r.Route("/reports", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Post("/", curationHandler.CreateReport)

				r.Group(func(r chi.Router) {
					r.Use(authMiddleware.RequireRole(models.RoleAdmin))
					r.Get("/", curationHandler.ListReports)
					r.Put("/{reportId}/resolve", curationHandler.ResolveReport)
				})
			})
		})

		// Módulo 5: Barajas de Estudio (CU10)
		r.Route("/decks", func(r chi.Router) {
			r.Get("/", studioHandler.ListDecks)
			r.With(authMiddleware.OptionalAuth).Get("/{id}", studioHandler.GetDeck)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Post("/", studioHandler.CreateDeck)
				r.Put("/{id}", studioHandler.UpdateDeck)
				r.Delete("/{id}", studioHandler.DeleteDeck)
				r.Post("/{id}/items", studioHandler.AddDeckItem)
				r.Delete("/{id}/items/{senseId}", studioHandler.RemoveDeckItem)
				r.Post("/{id}/save", studioHandler.SaveDeck)
				r.Delete("/{id}/save", studioHandler.UnsaveDeck)
			})
		})

		// Módulo 5: Colecciones Temáticas (CU11)
		r.Route("/collections", func(r chi.Router) {
			r.Get("/", studioHandler.ListCollections)
			r.With(authMiddleware.OptionalAuth).Get("/{id}", studioHandler.GetCollection)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Post("/", studioHandler.CreateCollection)
				r.Put("/{id}", studioHandler.UpdateCollection)
				r.Delete("/{id}", studioHandler.DeleteCollection)
				r.Post("/{id}/save", studioHandler.SaveCollection)
				r.Delete("/{id}/save", studioHandler.UnsaveCollection)
			})
		})

		// Módulo 5 & 6: Studio Personal, Repetición Espaciada SM-2 y Maestría (CU10, CU11, CU12, CU13)
		r.Route("/studio", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)

				// Studio Decks & Collections (CU10, CU11)
				r.Get("/decks", studioHandler.ListStudioDecks)
				r.Get("/collections", studioHandler.ListStudioCollections)

				// Repetición Espaciada SM-2 (CU12)
				r.Get("/study-session", masteryHandler.GetStudySession)
				r.Post("/study-session/review", masteryHandler.SubmitReview)

				// Maestría y Ledger de XP (CU13)
				r.Get("/mastery", masteryHandler.GetMasteryDashboard)
				r.Get("/mastery/items", masteryHandler.ListMasteryItems)
				r.Get("/xp-ledger", masteryHandler.ListXPLedger)
			})
		})
	})

	return r
}
