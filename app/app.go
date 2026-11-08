package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"

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

	wsMu   sync.Mutex
	wsConn *websocket.Conn
	wsStop chan struct{}
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
			Timeout: 20 * time.Second,
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

func (a *App) context() context.Context {
	if a.ctx != nil {
		return a.ctx
	}
	return context.Background()
}

func (a *App) getSessionToken() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.sessionToken
}

func (a *App) setSessionToken(token string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.sessionToken = token
}

// doRequest serializa JSON, inyecta Bearer JWT y maneja errores con descripciones humanas exactas.
func (a *App) doRequest(ctx context.Context, method, targetURL string, body any, target any) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error al serializar cuerpo JSON: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, targetURL, reqBody)
	if err != nil {
		return fmt.Errorf("error al crear petición HTTP: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	token := a.getSessionToken()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error de red al conectar con el servidor: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error string `json:"error"`
		}
		respBody, _ := io.ReadAll(resp.Body)
		if len(respBody) > 0 {
			if jsonErr := json.Unmarshal(respBody, &errResp); jsonErr == nil && errResp.Error != "" {
				return errors.New(errResp.Error)
			}
		}
		return errors.New(http.StatusText(resp.StatusCode))
	}

	if resp.StatusCode == http.StatusNoContent || target == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("error al decodificar respuesta JSON: %w", err)
	}
	return nil
}

// HealthCheck returns connectivity status for local client.
func (a *App) HealthCheck() (string, error) {
	return "ok", nil
}

// HasActiveSession checks if there is an active session token in memory.
func (a *App) HasActiveSession() bool {
	return a.getSessionToken() != ""
}

// =============================================================================
// Puente de Modelos para Generación de Bindings en Wails CLI
// =============================================================================

// ModelBindingsBridge garantiza que Wails exporte todos los structs a TypeScript.
func (a *App) ModelBindingsBridge(
	_ *models.SignInRequest,
	_ *models.SignUpRequest,
	_ *models.AuthResponse,
	_ *models.UserProfile,
	_ *models.UpdateProfileRequest,
	_ *models.ChangePasswordRequest,
	_ *models.PublicUserProfile,
	_ *models.CreateEditorRequestDTO,
	_ *models.ResolveEditorRequestDTO,
	_ *models.EditorRequest,
	_ *models.UpdateUserRoleOrStatusRequest,
	_ *models.CreateTicketRequest,
	_ *models.InvitationTicket,
	_ *models.RedeemTicketRequest,
	_ *models.WordSearchResult,
	_ *models.Example,
	_ *models.HeteronymItem,
	_ *models.SenseDetail,
	_ *models.WordDetail,
	_ *models.WordExploreFilter,
	_ *models.WordSummaryCard,
	_ *models.PaginatedWords,
	_ *models.CreateExampleInput,
	_ *models.CreateSenseRequest,
	_ *models.CreateWordProposalRequest,
	_ *models.UpdateWordRequest,
	_ *models.UpdateSenseRequest,
	_ *models.SemanticConcept,
	_ *models.CreateConceptRequest,
	_ *models.SenseApprovalVote,
	_ *models.SenseQuorumItem,
	_ *models.SenseQuorumDetail,
	_ *models.SubmitVoteRequest,
	_ *models.AdminOverrideStateRequest,
	_ *models.ContentReport,
	_ *models.CreateReportRequest,
	_ *models.ResolveReportRequest,
	_ *models.DeckItem,
	_ *models.Deck,
	_ *models.Collection,
	_ *models.UserMastery,
	_ *models.StudioDeckSummary,
	_ *models.CreateDeckRequest,
	_ *models.UpdateDeckRequest,
	_ *models.AddDeckItemRequest,
	_ *models.CreateCollectionRequest,
	_ *models.UpdateCollectionRequest,
	_ *models.StudyCardItem,
	_ *models.StudySessionPayload,
	_ *models.SubmitReviewRequest,
	_ *models.CardReviewResult,
	_ *models.XPTransactionItem,
	_ *models.AchievementStatus,
	_ *models.MasteryDashboard,
	_ *models.ArenaParticipant,
	_ *models.ArenaQuestion,
	_ *models.ArenaRoomState,
	_ *models.WSEvent,
	_ *models.CreateMatchRequest,
	_ *models.InternalArenaQuestion,
	_ *models.MatchPackage,
	_ *models.SubmitAnswerPayload,
	_ *models.RoundResultPayload,
	_ *models.ParticipantRoundAttempt,
	_ *models.ParticipantFinalResult,
	_ *models.ConsolidateMatchRequest,
	_ *models.MatchSummary,
	_ *models.ParticipantRoundAnswer,
	_ *models.MatchRoundReplay,
	_ *models.MatchReplayDetail,
	_ *models.ImageUploadResponse,
) error {
	return nil
}
