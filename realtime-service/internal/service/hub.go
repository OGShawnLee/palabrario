package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/response"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Permite conexiones desde cualquier origen en desarrollo y Wails
	},
}

// ArenaHub administra todas las salas activas de la Arena en memoria RAM.
type ArenaHub struct {
	mu                sync.RWMutex
	rooms             map[string]*ActiveRoom
	coreClient        *CoreClient
	tokenService      *auth.TokenService
	ReconnectWindow   time.Duration
	CountdownDuration time.Duration
	InterRoundPause   time.Duration
}

// NewArenaHub inicializa el hub central de Arena con configuraciones por defecto.
func NewArenaHub(coreClient *CoreClient, tokenService *auth.TokenService) *ArenaHub {
	return &ArenaHub{
		rooms:             make(map[string]*ActiveRoom),
		coreClient:        coreClient,
		tokenService:      tokenService,
		ReconnectWindow:   10 * time.Second,
		CountdownDuration: 3 * time.Second,
		InterRoundPause:   3 * time.Second,
	}
}

// ActiveRoomsCount retorna el número de salas activas en memoria RAM de forma concurrente y segura.
func (h *ArenaHub) ActiveRoomsCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms)
}

// GetRoom retorna una sala activa si existe en RAM.
func (h *ArenaHub) GetRoom(roomCode string) *ActiveRoom {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.rooms[roomCode]
}

// RemoveRoom retira una sala de la memoria RAM de forma segura.
func (h *ArenaHub) RemoveRoom(roomCode string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.rooms, roomCode)
}

// GetOrCreateRoom busca la sala en RAM o la hidrata desde core-service vía S2S.
func (h *ArenaHub) GetOrCreateRoom(ctx context.Context, roomCode string) (*ActiveRoom, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, exists := h.rooms[roomCode]; exists {
		return room, nil
	}

	pkg, err := h.coreClient.FetchMatchPackage(ctx, roomCode)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el paquete de partida de core-service: %w", err)
	}

	roundTimeSec := pkg.RoundTimeSec
	if roundTimeSec <= 0 {
		roundTimeSec = 15
	}
	totalRounds := pkg.TotalRounds
	if totalRounds <= 0 {
		totalRounds = len(pkg.Questions)
	}

	room := newActiveRoom(roomCode, h, h.coreClient, pkg, roundTimeSec, totalRounds)
	h.rooms[roomCode] = room

	go room.RunLoop()

	return room, nil
}

// HandleWS gestiona el upgrade de WebSocket y canaliza el ciclo de vida del jugador en la sala.
func (h *ArenaHub) HandleWS(w http.ResponseWriter, r *http.Request, roomCode string, claims *auth.Claims) {
	room, err := h.GetOrCreateRoom(r.Context(), roomCode)
	if err != nil {
		response.Error(w, http.StatusNotFound, fmt.Sprintf("Sala %s no disponible: %v", roomCode, err))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("⚠️ Error actualizando a WebSocket para sala %s: %v", roomCode, err)
		return
	}

	client := newWSClient(conn)

	resCh := make(chan error, 1)
	room.joinCh <- &joinRequest{
		claims: claims,
		client: client,
		resCh:  resCh,
	}

	if err := <-resCh; err != nil {
		_ = client.SendJSON(models.WSEvent{
			Type: "ERROR",
			Payload: map[string]string{
				"error": err.Error(),
			},
		})
		_ = client.Close()
		return
	}

	// Ciclo de lectura
	defer func() {
		room.leaveCh <- &leaveRequest{
			userID: claims.UserID,
			client: client,
		}
		_ = client.Close()
	}()

	conn.SetReadLimit(4096)
	for {
		var event models.WSEvent
		if err := conn.ReadJSON(&event); err != nil {
			break
		}

		select {
		case room.actionCh <- &actionMessage{
			userID: claims.UserID,
			event:  event,
		}:
		case <-room.stopCh:
			return
		}
	}
}
