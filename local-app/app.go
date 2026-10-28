package main

import (
	"context"
	"net/http"
	"os"
	"sync"
	"time"

	"palabrario/pkg/models"
)

// App is the desktop Backend-For-Frontend service for Palabrario in Wails v3.
type App struct {
	ctx          context.Context
	httpClient   *http.Client
	coreURL      string
	realtimeURL  string
	imageURL     string
	mu           sync.RWMutex
	sessionToken string
}

// NewApp instantiates the App service reading environment defaults.
func NewApp() *App {
	coreURL := os.Getenv("PALABRARIO_CORE_URL")
	if coreURL == "" {
		coreURL = "http://localhost:8000"
	}

	realtimeURL := os.Getenv("PALABRARIO_REALTIME_URL")
	if realtimeURL == "" {
		realtimeURL = "ws://localhost:8001"
	}

	imageURL := os.Getenv("PALABRARIO_IMAGE_URL")
	if imageURL == "" {
		imageURL = "http://localhost:8002"
	}

	return &App{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		coreURL:     coreURL,
		realtimeURL: realtimeURL,
		imageURL:    imageURL,
	}
}

// SetContext sets the context from Wails runtime.
func (a *App) SetContext(ctx context.Context) {
	a.ctx = ctx
}

// HasActiveSession checks if there is an active session token in memory.
func (a *App) HasActiveSession() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.sessionToken != ""
}

// SignOut clears the in-memory session token.
func (a *App) SignOut() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.sessionToken = ""
	return nil
}

// HealthCheck returns connectivity status for local client.
func (a *App) HealthCheck() (string, error) {
	return "ok", nil
}

// ModelBindingsBridge references models from pkg/models to enable Wails CLI binding generation.
func (a *App) ModelBindingsBridge(
	_ *models.SignInRequest,
	_ *models.SignUpRequest,
	_ *models.AuthResponse,
	_ *models.UserProfile,
	_ *models.WordSummaryCard,
	_ *models.WordDetail,
	_ *models.SenseDetail,
	_ *models.HeteronymItem,
	_ *models.StudioDeckSummary,
	_ *models.Deck,
	_ *models.Collection,
	_ *models.StudyCardItem,
	_ *models.ArenaRoomState,
	_ *models.WSEvent,
	_ *models.ImageUploadResponse,
) error {
	return nil
}
