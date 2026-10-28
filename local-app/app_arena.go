package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v3/pkg/application"

	"palabrario/pkg/models"
)

// =============================================================================
// MÓDULO 7: Arena Multijugador y WebSockets (CU14–CU16)
// =============================================================================

// CreateArenaMatch prepara una nueva partida competitiva en vivo.
func (a *App) CreateArenaMatch(req models.CreateMatchRequest) (*models.MatchPackage, error) {
	return a.PrepareMatch(req)
}

// PrepareMatch prepara una partida de arena.
func (a *App) PrepareMatch(req models.CreateMatchRequest) (*models.MatchPackage, error) {
	var pkg models.MatchPackage
	if err := a.doRequest(a.context(), http.MethodPost, a.coreURL+"/api/v1/arena/matches/prepare", req, &pkg); err != nil {
		return nil, err
	}
	return &pkg, nil
}

// JoinArenaMatch valida e inicializa la participación en una sala por código PIN.
func (a *App) JoinArenaMatch(roomCode string) (*models.MatchPackage, error) {
	roomCode = strings.TrimSpace(roomCode)
	if roomCode == "" {
		return nil, errors.New("código de sala requerido")
	}
	return &models.MatchPackage{RoomCode: roomCode}, nil
}

// GetMyMatchHistory consulta el historial de partidas del usuario.
func (a *App) GetMyMatchHistory(limit int) ([]models.MatchSummary, error) {
	return a.ListMatchHistory(limit)
}

// ListMatchHistory lista historial de partidas.
func (a *App) ListMatchHistory(limit int) ([]models.MatchSummary, error) {
	targetURL := a.coreURL + "/api/v1/arena/matches/history"
	if limit > 0 {
		targetURL += "?limit=" + strconv.Itoa(limit)
	}
	var list []models.MatchSummary
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.MatchSummary{}
	}
	return list, nil
}

// GetMatchReplayDetail recupera el replay granular 3NF ronda por ronda.
func (a *App) GetMatchReplayDetail(matchId string) (*models.MatchReplayDetail, error) {
	return a.GetMatchReplay(matchId)
}

// GetMatchReplay obtiene replay de partida.
func (a *App) GetMatchReplay(matchId string) (*models.MatchReplayDetail, error) {
	var rep models.MatchReplayDetail
	targetURL := a.coreURL + "/api/v1/arena/matches/" + url.PathEscape(matchId)
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &rep); err != nil {
		return nil, err
	}
	return &rep, nil
}

// ConnectArenaRoom establece conexión WebSocket hacia realtime-service (:8001).
func (a *App) ConnectArenaRoom(roomCode string) error {
	a.wsMu.Lock()
	defer a.wsMu.Unlock()

	// Si hay conexión previa abierta, cerrarla limpiamente
	if a.wsConn != nil {
		if a.wsStop != nil {
			close(a.wsStop)
			a.wsStop = nil
		}
		_ = a.wsConn.Close()
		a.wsConn = nil
	}

	token := a.getSessionToken()
	if token == "" {
		return errors.New("se requiere una sesión activa para ingresar a la arena")
	}

	roomCode = strings.TrimSpace(roomCode)
	if roomCode == "" {
		return errors.New("código de sala requerido")
	}

	wsBase := a.realtimeURL
	if strings.HasPrefix(wsBase, "http://") {
		wsBase = "ws://" + strings.TrimPrefix(wsBase, "http://")
	} else if strings.HasPrefix(wsBase, "https://") {
		wsBase = "wss://" + strings.TrimPrefix(wsBase, "https://")
	}

	wsURL := fmt.Sprintf("%s/ws/arena/%s?token=%s", wsBase, url.PathEscape(roomCode), url.QueryEscape(token))

	conn, _, err := websocket.DefaultDialer.DialContext(a.context(), wsURL, nil)
	if err != nil {
		return fmt.Errorf("error al conectar con la sala de arena: %w", err)
	}

	a.wsConn = conn
	a.wsStop = make(chan struct{})

	// Goroutine receptora de tramas WS
	go func(ws *websocket.Conn, stop <-chan struct{}) {
		defer func() {
			_ = ws.Close()
		}()

		for {
			var evt models.WSEvent
			err := ws.ReadJSON(&evt)
			if err != nil {
				select {
				case <-stop:
					return
				default:
					if app := application.Get(); app != nil && app.Event != nil {
						app.Event.Emit("arena:ws_closed", map[string]string{"error": err.Error()})
					}
					return
				}
			}

			if app := application.Get(); app != nil && app.Event != nil {
				app.Event.Emit("arena:ws_message", evt)
			}
		}
	}(conn, a.wsStop)

	return nil
}

// ConnectArena alias para conexión a sala de arena.
func (a *App) ConnectArena(roomCode string) error {
	return a.ConnectArenaRoom(roomCode)
}

// StartArenaMatch emite el evento para iniciar la partida (Host).
func (a *App) StartArenaMatch() error {
	a.wsMu.Lock()
	defer a.wsMu.Unlock()

	if a.wsConn == nil {
		return errors.New("no hay conexión activa con la sala de arena")
	}

	evt := models.WSEvent{
		Type: "START_MATCH",
	}
	return a.wsConn.WriteJSON(evt)
}

// SubmitArenaAnswer envía la respuesta del jugador en la ronda actual.
func (a *App) SubmitArenaAnswer(roundNumber int, selectedOption string) error {
	a.wsMu.Lock()
	defer a.wsMu.Unlock()

	if a.wsConn == nil {
		return errors.New("no hay conexión activa con la sala de arena")
	}

	evt := models.WSEvent{
		Type: "SUBMIT_ANSWER",
		Payload: models.SubmitAnswerPayload{
			RoundNumber:    roundNumber,
			SelectedOption: selectedOption,
		},
	}
	return a.wsConn.WriteJSON(evt)
}

// LeaveArenaRoom cierra la conexión activa con la sala.
func (a *App) LeaveArenaRoom() error {
	a.wsMu.Lock()
	defer a.wsMu.Unlock()

	if a.wsStop != nil {
		close(a.wsStop)
		a.wsStop = nil
	}
	if a.wsConn != nil {
		_ = a.wsConn.Close()
		a.wsConn = nil
	}
	return nil
}

// LeaveArena alias para abandonar sala.
func (a *App) LeaveArena() error {
	return a.LeaveArenaRoom()
}
