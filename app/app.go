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
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v3/pkg/application"

	"palabrario/pkg/models"
)

// APIError represents a structured HTTP or transport error returned by the BFF.
type APIError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string { return e.Message }

type persistedSession struct {
	RefreshToken string `json:"refresh_token"`
}

// App is the desktop Backend-For-Frontend service for Palabrario in Wails v3.
type App struct {
	ctx             context.Context
	httpClient      *http.Client
	coreURL         string
	realtimeURL     string
	imageURL        string
	sessionFilePath string

	mu           sync.RWMutex
	accessToken  string
	refreshToken string

	refreshMu sync.Mutex

	wsMu   sync.Mutex
	wsConn *websocket.Conn
	wsStop chan struct{}
}

// NewApp instantiates the App service reading environment defaults and loading any persisted refresh token.
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

	cfgDir, err := os.UserConfigDir()
	if err != nil || strings.TrimSpace(cfgDir) == "" {
		cfgDir = os.TempDir()
	}
	sessionFilePath := filepath.Join(cfgDir, "palabrario", "session.json")

	return newAppWithConfig(coreURL, realtimeURL, imageURL, sessionFilePath)
}

func newAppWithConfig(coreURL, realtimeURL, imageURL, sessionFilePath string) *App {
	a := &App{
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
		coreURL:         strings.TrimRight(coreURL, "/"),
		realtimeURL:     strings.TrimRight(realtimeURL, "/"),
		imageURL:        strings.TrimRight(imageURL, "/"),
		sessionFilePath: sessionFilePath,
	}
	a.loadPersistedSession()
	return a
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

func (a *App) getAccessToken() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.accessToken
}

func (a *App) getRefreshToken() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.refreshToken
}

func (a *App) getTokens() (string, string) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.accessToken, a.refreshToken
}

func (a *App) setTokens(access, refresh string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.accessToken = access
	a.refreshToken = refresh
	if refresh != "" {
		_ = a.writePersistedSessionLocked(refresh)
	} else {
		_ = a.removePersistedSessionLocked()
	}
}

func (a *App) clearSession() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.accessToken = ""
	a.refreshToken = ""
	_ = a.removePersistedSessionLocked()
}

func (a *App) loadPersistedSession() {
	if strings.TrimSpace(a.sessionFilePath) == "" {
		return
	}
	data, err := os.ReadFile(a.sessionFilePath)
	if err != nil {
		return
	}
	var sess persistedSession
	if err := json.Unmarshal(data, &sess); err != nil {
		return
	}
	token := strings.TrimSpace(sess.RefreshToken)
	if token == "" {
		return
	}
	a.mu.Lock()
	a.refreshToken = token
	a.mu.Unlock()
}

func (a *App) writePersistedSessionLocked(refresh string) error {
	if strings.TrimSpace(a.sessionFilePath) == "" {
		return nil
	}
	dir := filepath.Dir(a.sessionFilePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	payload, err := json.Marshal(persistedSession{RefreshToken: refresh})
	if err != nil {
		return err
	}
	tmpFile, err := os.CreateTemp(dir, "session-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	cleanupTmp := true
	defer func() {
		if cleanupTmp {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := tmpFile.Write(payload); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpPath, 0600); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, a.sessionFilePath); err != nil {
		return err
	}
	cleanupTmp = false
	return nil
}

func (a *App) removePersistedSessionLocked() error {
	if strings.TrimSpace(a.sessionFilePath) == "" {
		return nil
	}
	err := os.Remove(a.sessionFilePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (a *App) emitAuthExpired(code, reason string) {
	if app := application.Get(); app != nil && app.Event != nil {
		app.Event.Emit("auth:expired", map[string]string{
			"code":   code,
			"reason": reason,
		})
	}
}

func parseAPIError(resp *http.Response) *APIError {
	respBody, _ := io.ReadAll(resp.Body)
	var errResp struct {
		Code  string `json:"code"`
		Error string `json:"error"`
	}
	if len(respBody) > 0 {
		_ = json.Unmarshal(respBody, &errResp)
	}
	msg := strings.TrimSpace(errResp.Error)
	if msg == "" {
		msg = http.StatusText(resp.StatusCode)
	}
	return &APIError{
		Status:  resp.StatusCode,
		Code:    strings.TrimSpace(errResp.Code),
		Message: msg,
	}
}

// refreshSession exchanges the persisted refresh token for a fresh access/refresh token pair.
func (a *App) refreshSession(ctx context.Context) (*models.UserProfile, error) {
	a.refreshMu.Lock()
	defer a.refreshMu.Unlock()
	return a.doRefreshLocked(ctx)
}

// refreshSessionIfStale performs a single-flight refresh only if the access token has not already been rotated by another goroutine.
func (a *App) refreshSessionIfStale(ctx context.Context, staleAccessToken string) (*models.UserProfile, bool, error) {
	a.refreshMu.Lock()
	defer a.refreshMu.Unlock()

	curAccess := a.getAccessToken()
	if staleAccessToken == "" && curAccess != "" {
		return nil, false, nil
	}
	if staleAccessToken != "" && curAccess != "" && curAccess != staleAccessToken {
		return nil, false, nil
	}

	profile, err := a.doRefreshLocked(ctx)
	return profile, true, err
}

func (a *App) doRefreshLocked(ctx context.Context) (*models.UserProfile, error) {
	refresh := a.getRefreshToken()
	if refresh == "" {
		apiErr := &APIError{
			Status:  http.StatusUnauthorized,
			Code:    "UNAUTHORIZED",
			Message: "Your session has expired. Please sign in again.",
		}
		a.clearSession()
		a.emitAuthExpired(apiErr.Code, apiErr.Message)
		return nil, apiErr
	}

	payload, err := json.Marshal(models.RefreshTokenRequest{RefreshToken: refresh})
	if err != nil {
		return nil, fmt.Errorf("failed to serialize refresh request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.coreURL+"/api/v1/auth/refresh", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, &APIError{
			Status:  0,
			Code:    "NETWORK_ERROR",
			Message: "Unable to reach the Palabrario server. Please check your connection.",
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		apiErr := parseAPIError(resp)
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			a.clearSession()
			a.emitAuthExpired(apiErr.Code, apiErr.Message)
		}
		return nil, apiErr
	}

	var authResp models.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, fmt.Errorf("failed to decode refresh response: %w", err)
	}

	a.setTokens(authResp.Token, authResp.RefreshToken)
	return &authResp.User, nil
}

// executeRequest executes a request built by buildReq, injects Bearer auth, handles single-flight 401 refresh,
// and evicts expired/suspended sessions cleanly.
func (a *App) executeRequest(ctx context.Context, buildReq func() (*http.Request, error), target any, allowRetry bool) error {
	req, err := buildReq()
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	token := a.getAccessToken()
	tokenWasSent := false
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
		tokenWasSent = true
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return &APIError{
			Status:  0,
			Code:    "NETWORK_ERROR",
			Message: "Unable to reach the Palabrario server. Please check your connection.",
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		apiErr := parseAPIError(resp)

		if resp.StatusCode == http.StatusUnauthorized {
			if allowRetry && tokenWasSent && a.getRefreshToken() != "" {
				if _, _, refreshErr := a.refreshSessionIfStale(ctx, token); refreshErr == nil {
					return a.executeRequest(ctx, buildReq, target, false)
				} else {
					return refreshErr
				}
			}
			if tokenWasSent {
				a.clearSession()
				a.emitAuthExpired(apiErr.Code, apiErr.Message)
			}
		} else if resp.StatusCode == http.StatusForbidden && apiErr.Code == "ACCOUNT_SUSPENDED" {
			a.clearSession()
			a.emitAuthExpired(apiErr.Code, apiErr.Message)
		}

		return apiErr
	}

	if resp.StatusCode == http.StatusNoContent || target == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode server response: %w", err)
	}
	return nil
}

// doRequest serializes JSON payloads and delegates execution to executeRequest.
func (a *App) doRequest(ctx context.Context, method, targetURL string, body any, target any) error {
	var data []byte
	if body != nil {
		var err error
		data, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to serialize JSON body: %w", err)
		}
	}

	return a.executeRequest(ctx, func() (*http.Request, error) {
		var reqBody io.Reader
		if data != nil {
			reqBody = bytes.NewReader(data)
		}
		req, err := http.NewRequestWithContext(ctx, method, targetURL, reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to create HTTP request: %w", err)
		}
		if data != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		return req, nil
	}, target, true)
}

// HealthCheck returns connectivity status for local client.
func (a *App) HealthCheck() (string, error) {
	return "ok", nil
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
