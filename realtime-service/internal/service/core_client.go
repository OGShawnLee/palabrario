package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"palabrario/pkg/auth"
	"palabrario/pkg/models"
)

// CoreClient interactúa mediante HTTP con core-service utilizando autenticación S2S.
type CoreClient struct {
	baseURL      string
	httpClient   *http.Client
	tokenService *auth.TokenService
}

// NewCoreClient crea una nueva instancia de CoreClient hacia core-service.
func NewCoreClient(baseURL string, httpClient *http.Client, tokenService *auth.TokenService) *CoreClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &CoreClient{
		baseURL:      strings.TrimRight(baseURL, "/"),
		httpClient:   httpClient,
		tokenService: tokenService,
	}
}

// FetchMatchPackage obtiene el paquete de partida hidratado desde core-service con token S2S.
func (c *CoreClient) FetchMatchPackage(ctx context.Context, roomCode string) (*models.MatchPackage, error) {
	s2sToken, err := c.tokenService.GenerateInternalToken("realtime-service", roomCode)
	if err != nil {
		return nil, fmt.Errorf("error generando token S2S para FetchMatchPackage: %w", err)
	}

	endpoint := fmt.Sprintf("%s/api/v1/arena/matches/%s/package", c.baseURL, url.PathEscape(roomCode))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error construyendo request para FetchMatchPackage: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s2sToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando HTTP request hacia core-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("core-service respondió con código %d: %s", resp.StatusCode, string(body))
	}

	var matchPkg models.MatchPackage
	if err := json.NewDecoder(resp.Body).Decode(&matchPkg); err != nil {
		return nil, fmt.Errorf("error decodificando MatchPackage: %w", err)
	}

	return &matchPkg, nil
}

// ConsolidateMatch envía el resultado final de la partida a core-service para consolidación ACID con token S2S.
func (c *CoreClient) ConsolidateMatch(ctx context.Context, req models.ConsolidateMatchRequest) (*models.MatchSummary, error) {
	s2sToken, err := c.tokenService.GenerateInternalToken("realtime-service", req.RoomCode)
	if err != nil {
		return nil, fmt.Errorf("error generando token S2S para ConsolidateMatch: %w", err)
	}

	endpoint := fmt.Sprintf("%s/api/v1/arena/matches/consolidate", c.baseURL)
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error serializando ConsolidateMatchRequest: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error construyendo request para ConsolidateMatch: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+s2sToken)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando ConsolidateMatch en core-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("core-service rechazó consolidación con código %d: %s", resp.StatusCode, string(body))
	}

	var summary models.MatchSummary
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		return nil, fmt.Errorf("error decodificando MatchSummary: %w", err)
	}

	return &summary, nil
}
