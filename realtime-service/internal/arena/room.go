package arena

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"palabrario/pkg/auth"
	"palabrario/pkg/models"
)

const MaxParticipants = 8

// wsClient envuelve una conexión WebSocket con un mutex para envíos concurrentes seguros.
type wsClient struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func newWSClient(conn *websocket.Conn) *wsClient {
	return &wsClient{conn: conn}
}

func (c *wsClient) SendJSON(v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_ = c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	return c.conn.WriteJSON(v)
}

func (c *wsClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.Close()
}

// roomParticipant representa el estado en RAM de un jugador dentro de ActiveRoom.
type roomParticipant struct {
	IDUser           string
	UserName         string
	Name             string
	AvatarURL        string
	DialectCode      string
	FlagEmoji        string
	Score            int
	Rank             int
	IsHost           bool
	IsReady          bool
	Connected        bool
	JoinedAt         time.Time
	TotalTimeSpentMs int
	CorrectSenseIDs  []int
	WrongSenseIDs    []int
	RoundAttempts    []models.ParticipantRoundAttempt

	client          *wsClient
	disconnectTimer *time.Timer
}

type joinRequest struct {
	claims *auth.Claims
	client *wsClient
	resCh  chan error
}

type leaveRequest struct {
	userID string
	client *wsClient
}

type actionMessage struct {
	userID string
	event  models.WSEvent
}

type timerMsgKind int

const (
	timerReconnectExpired timerMsgKind = iota
	timerCountdownFinished
	timerRoundExpired
	timerStartNextRound
	timerFinishMatch
)

type roomTimerMsg struct {
	kind   timerMsgKind
	userID string
	round  int
}

// ActiveRoom gobierna el ciclo de vida y los eventos de una sala de Arena en RAM.
type ActiveRoom struct {
	roomCode          string
	hub               *ArenaHub
	coreClient        *CoreClient
	pkg               *models.MatchPackage
	HostUserID        string
	Status            models.MatchStatus
	TotalRounds       int
	RoundTimeSec      int
	CurrentRound      int
	roundStartedAt    time.Time
	roundFinished     bool
	roundTimer        *time.Timer
	roundAnswered     map[string]bool
	currentQuestion   *models.InternalArenaQuestion
	participants      map[string]*roomParticipant
	ReconnectWindow   time.Duration
	CountdownDuration time.Duration
	InterRoundPause   time.Duration

	joinCh   chan *joinRequest
	leaveCh  chan *leaveRequest
	actionCh chan *actionMessage
	timerCh  chan roomTimerMsg
	stopCh   chan struct{}
	stopOnce sync.Once
}

func newActiveRoom(
	roomCode string,
	hub *ArenaHub,
	coreClient *CoreClient,
	pkg *models.MatchPackage,
	roundTimeSec int,
	totalRounds int,
) *ActiveRoom {
	return &ActiveRoom{
		roomCode:          roomCode,
		hub:               hub,
		coreClient:        coreClient,
		pkg:               pkg,
		HostUserID:        pkg.HostUserID,
		Status:            models.MatchStatusWaiting,
		TotalRounds:       totalRounds,
		RoundTimeSec:      roundTimeSec,
		ReconnectWindow:   hub.ReconnectWindow,
		CountdownDuration: hub.CountdownDuration,
		InterRoundPause:   hub.InterRoundPause,
		participants:      make(map[string]*roomParticipant),
		roundAnswered:     make(map[string]bool),
		joinCh:            make(chan *joinRequest),
		leaveCh:           make(chan *leaveRequest),
		actionCh:          make(chan *actionMessage, 64),
		timerCh:           make(chan roomTimerMsg, 64),
		stopCh:            make(chan struct{}),
	}
}

// RunLoop coordina todos los eventos de la sala en una única goroutine, eliminando data races.
func (r *ActiveRoom) RunLoop() {
	for {
		select {
		case <-r.stopCh:
			r.cleanup()
			return

		case req := <-r.joinCh:
			r.handleJoin(req)

		case req := <-r.leaveCh:
			r.handleLeave(req)

		case act := <-r.actionCh:
			r.handleAction(act)

		case msg := <-r.timerCh:
			r.handleTimer(msg)
		}
	}
}

func (r *ActiveRoom) cleanup() {
	if r.roundTimer != nil {
		r.roundTimer.Stop()
	}
	for _, p := range r.participants {
		if p.disconnectTimer != nil {
			p.disconnectTimer.Stop()
		}
		if p.client != nil {
			_ = p.client.Close()
		}
	}
}

func (r *ActiveRoom) stop() {
	r.stopOnce.Do(func() {
		close(r.stopCh)
	})
}

func (r *ActiveRoom) handleJoin(req *joinRequest) {
	userID := req.claims.UserID

	// Caso 1: Reconexión de un participante existente (RNF-CON-02)
	if p, exists := r.participants[userID]; exists {
		if p.disconnectTimer != nil {
			p.disconnectTimer.Stop()
			p.disconnectTimer = nil
		}
		if p.client != nil && p.client != req.client {
			_ = p.client.Close()
		}
		p.client = req.client
		p.Connected = true
		p.IsReady = true

		req.resCh <- nil

		r.broadcastRoomState()

		// Si la partida está en curso y la pregunta actual está activa, enviar QUESTION inmediatamente
		if r.Status == models.MatchStatusInProgress && r.currentQuestion != nil && !r.roundFinished {
			_ = req.client.SendJSON(models.WSEvent{
				Type: "QUESTION",
				Payload: models.ArenaQuestion{
					RoundNumber:   r.currentQuestion.RoundNumber,
					Modality:      r.currentQuestion.Modality,
					Prompt:        r.currentQuestion.Prompt,
					Options:       r.currentQuestion.Options,
					Explanation:   r.currentQuestion.Explanation,
					TimeLimitSecs: r.RoundTimeSec,
				},
			})
		}
		return
	}

	// Caso 2: Nuevo participante
	if r.Status != models.MatchStatusWaiting {
		req.resCh <- fmt.Errorf("La partida ya está en progreso")
		return
	}

	if len(r.participants) >= MaxParticipants {
		req.resCh <- fmt.Errorf("La sala está llena (máximo 8 jugadores)")
		return
	}

	// Si no hay host configurado, el primero en unirse es Host
	if r.HostUserID == "" {
		r.HostUserID = userID
	}

	name := req.claims.UserName
	if name == "" {
		name = req.claims.Email
	}

	isHost := (userID == r.HostUserID)

	p := &roomParticipant{
		IDUser:          userID,
		UserName:        req.claims.UserName,
		Name:            name,
		AvatarURL:       "",
		DialectCode:     req.claims.AnchorDialect,
		FlagEmoji:       DialectToFlag(req.claims.AnchorDialect),
		Score:           0,
		Rank:            len(r.participants) + 1,
		IsHost:          isHost,
		IsReady:         true,
		Connected:       true,
		JoinedAt:        time.Now().UTC(),
		CorrectSenseIDs: make([]int, 0),
		WrongSenseIDs:   make([]int, 0),
		RoundAttempts:   make([]models.ParticipantRoundAttempt, 0),
		client:          req.client,
	}

	r.participants[userID] = p
	r.recalculateRanks()

	req.resCh <- nil

	r.broadcastRoomState()
}

func (r *ActiveRoom) handleLeave(req *leaveRequest) {
	p, exists := r.participants[req.userID]
	if !exists {
		return
	}

	// Si el socket que se cerró no es el activo, ignorar
	if p.client != req.client {
		return
	}

	p.Connected = false
	p.IsReady = false

	r.broadcastRoomState()

	// Programar ventana de reconexión (10s por default)
	uid := req.userID
	p.disconnectTimer = time.AfterFunc(r.ReconnectWindow, func() {
		select {
		case r.timerCh <- roomTimerMsg{kind: timerReconnectExpired, userID: uid}:
		case <-r.stopCh:
		}
	})
}

func (r *ActiveRoom) handleTimer(msg roomTimerMsg) {
	switch msg.kind {
	case timerReconnectExpired:
		p, exists := r.participants[msg.userID]
		if !exists || p.Connected {
			return
		}

		if r.Status == models.MatchStatusWaiting {
			delete(r.participants, msg.userID)

			connectedCount := 0
			for _, remaining := range r.participants {
				if remaining.Connected {
					connectedCount++
				}
			}

			// Si no queda ningún participante conectado, cerrar la sala
			if connectedCount == 0 {
				r.hub.RemoveRoom(r.roomCode)
				r.stop()
				return
			}

			// Sucesión Automática del Host si el eliminado era el Host
			if msg.userID == r.HostUserID {
				var oldest *roomParticipant
				for _, cand := range r.participants {
					if cand.Connected {
						if oldest == nil || cand.JoinedAt.Before(oldest.JoinedAt) {
							oldest = cand
						}
					}
				}
				if oldest != nil {
					r.HostUserID = oldest.IDUser
					for _, part := range r.participants {
						part.IsHost = (part.IDUser == r.HostUserID)
					}
				}
			}

			r.recalculateRanks()
			r.broadcastRoomState()

		} else if r.Status == models.MatchStatusInProgress {
			// Si todos los jugadores que quedan conectados ya respondieron, cerrar anticipadamente
			if r.allConnectedAnswered() && !r.roundFinished {
				r.finishCurrentRound()
			}
		}

	case timerCountdownFinished:
		if r.Status == models.MatchStatusInProgress {
			r.startRound(1)
		}

	case timerRoundExpired:
		if r.Status == models.MatchStatusInProgress && r.CurrentRound == msg.round && !r.roundFinished {
			r.finishCurrentRound()
		}

	case timerStartNextRound:
		if r.Status == models.MatchStatusInProgress && msg.round <= r.TotalRounds {
			r.startRound(msg.round)
		}

	case timerFinishMatch:
		r.finishMatch()
	}
}

func (r *ActiveRoom) handleAction(act *actionMessage) {
	switch act.event.Type {
	case "START_MATCH":
		if act.userID != r.HostUserID || r.Status != models.MatchStatusWaiting {
			return
		}

		r.Status = models.MatchStatusInProgress
		r.broadcastRoomState()

		countdownSecs := int(r.CountdownDuration.Seconds())
		if countdownSecs <= 0 {
			countdownSecs = 3
		}

		r.broadcast(models.WSEvent{
			Type:    "COUNTDOWN",
			Payload: map[string]int{"seconds": countdownSecs},
		})

		time.AfterFunc(r.CountdownDuration, func() {
			select {
			case r.timerCh <- roomTimerMsg{kind: timerCountdownFinished, round: 1}:
			case <-r.stopCh:
			}
		})

	case "SUBMIT_ANSWER":
		if r.Status != models.MatchStatusInProgress || r.currentQuestion == nil || r.roundFinished {
			return
		}

		var payload models.SubmitAnswerPayload
		payloadBytes, err := json.Marshal(act.event.Payload)
		if err != nil {
			return
		}
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return
		}

		if payload.RoundNumber != r.CurrentRound {
			return
		}

		if r.roundAnswered[act.userID] {
			return // Ya respondió en esta ronda
		}

		r.roundAnswered[act.userID] = true
		elapsedMs := time.Since(r.roundStartedAt).Milliseconds()

		p, exists := r.participants[act.userID]
		if !exists {
			return
		}

		isCorrect := (strings.TrimSpace(payload.SelectedOption) == strings.TrimSpace(r.currentQuestion.CorrectOption))
		pts := CalculateRoundPoints(isCorrect, elapsedMs, r.RoundTimeSec)

		p.Score += pts
		p.TotalTimeSpentMs += int(elapsedMs)
		if isCorrect {
			p.CorrectSenseIDs = append(p.CorrectSenseIDs, r.currentQuestion.IDSense)
		} else {
			p.WrongSenseIDs = append(p.WrongSenseIDs, r.currentQuestion.IDSense)
		}

		p.RoundAttempts = append(p.RoundAttempts, models.ParticipantRoundAttempt{
			RoundNumber:    r.CurrentRound,
			IDSense:        r.currentQuestion.IDSense,
			IDQuestion:     r.currentQuestion.IDQuestion,
			SelectedOption: payload.SelectedOption,
			IsCorrect:      isCorrect,
			TimeSpentMs:    int(elapsedMs),
			PointsEarned:   pts,
		})

		if r.allConnectedAnswered() {
			r.finishCurrentRound()
		}
	}
}

func (r *ActiveRoom) allConnectedAnswered() bool {
	hasConnected := false
	for _, p := range r.participants {
		if p.Connected {
			hasConnected = true
			if !r.roundAnswered[p.IDUser] {
				return false
			}
		}
	}
	return hasConnected
}

func (r *ActiveRoom) startRound(roundNum int) {
	r.CurrentRound = roundNum
	r.roundFinished = false
	r.roundStartedAt = time.Now()
	r.roundAnswered = make(map[string]bool)

	var currentQ *models.InternalArenaQuestion
	for i := range r.pkg.Questions {
		if r.pkg.Questions[i].RoundNumber == roundNum {
			currentQ = &r.pkg.Questions[i]
			break
		}
	}
	if currentQ == nil && len(r.pkg.Questions) >= roundNum {
		currentQ = &r.pkg.Questions[roundNum-1]
	}
	if currentQ == nil {
		log.Printf("⚠️ No se encontró reactivo para ronda %d en sala %s", roundNum, r.roomCode)
		return
	}
	r.currentQuestion = currentQ

	// NUNCA difundir CorrectOption en el evento QUESTION
	r.broadcast(models.WSEvent{
		Type: "QUESTION",
		Payload: models.ArenaQuestion{
			RoundNumber:   currentQ.RoundNumber,
			Modality:      currentQ.Modality,
			Prompt:        currentQ.Prompt,
			Options:       currentQ.Options,
			Explanation:   currentQ.Explanation,
			TimeLimitSecs: r.RoundTimeSec,
		},
	})

	rRound := roundNum
	r.roundTimer = time.AfterFunc(time.Duration(r.RoundTimeSec)*time.Second, func() {
		select {
		case r.timerCh <- roomTimerMsg{kind: timerRoundExpired, round: rRound}:
		case <-r.stopCh:
		}
	})
}

func (r *ActiveRoom) finishCurrentRound() {
	if r.roundFinished {
		return
	}
	r.roundFinished = true

	if r.roundTimer != nil {
		r.roundTimer.Stop()
		r.roundTimer = nil
	}

	timeoutMs := r.RoundTimeSec * 1000
	for _, p := range r.participants {
		if !r.roundAnswered[p.IDUser] {
			p.TotalTimeSpentMs += timeoutMs
			p.WrongSenseIDs = append(p.WrongSenseIDs, r.currentQuestion.IDSense)
			p.RoundAttempts = append(p.RoundAttempts, models.ParticipantRoundAttempt{
				RoundNumber:    r.CurrentRound,
				IDSense:        r.currentQuestion.IDSense,
				IDQuestion:     r.currentQuestion.IDQuestion,
				SelectedOption: "",
				IsCorrect:      false,
				TimeSpentMs:    timeoutMs,
				PointsEarned:   0,
			})
		}
	}

	r.recalculateRanks()
	leaderboard := r.buildLeaderboard()

	r.broadcast(models.WSEvent{
		Type: "ROUND_RESULT",
		Payload: models.RoundResultPayload{
			RoundNumber:   r.CurrentRound,
			IDSense:       r.currentQuestion.IDSense,
			CorrectOption: r.currentQuestion.CorrectOption,
			Explanation:   r.currentQuestion.Explanation,
			Leaderboard:   leaderboard,
		},
	})

	if r.CurrentRound >= r.TotalRounds {
		time.AfterFunc(r.InterRoundPause, func() {
			select {
			case r.timerCh <- roomTimerMsg{kind: timerFinishMatch}:
			case <-r.stopCh:
			}
		})
	} else {
		nextRound := r.CurrentRound + 1
		time.AfterFunc(r.InterRoundPause, func() {
			select {
			case r.timerCh <- roomTimerMsg{kind: timerStartNextRound, round: nextRound}:
			case <-r.stopCh:
			}
		})
	}
}

func (r *ActiveRoom) finishMatch() {
	r.Status = models.MatchStatusFinished
	r.recalculateRanks()

	partsFinal := make([]models.ParticipantFinalResult, 0, len(r.participants))
	for _, p := range r.participants {
		var accuracy float32 = 0
		totalAttempts := len(p.RoundAttempts)
		if totalAttempts > 0 {
			accuracy = float32(len(p.CorrectSenseIDs)) / float32(totalAttempts)
		}
		partsFinal = append(partsFinal, models.ParticipantFinalResult{
			IDUser:           p.IDUser,
			FinalScore:       p.Score,
			AccuracyRate:     accuracy,
			FinalRank:        p.Rank,
			TotalTimeSpentMs: p.TotalTimeSpentMs,
			CorrectSenseIDs:  p.CorrectSenseIDs,
			WrongSenseIDs:    p.WrongSenseIDs,
			RoundAttempts:    p.RoundAttempts,
		})
	}

	consolidateReq := models.ConsolidateMatchRequest{
		IDMatch:      r.pkg.IDMatch,
		RoomCode:     r.roomCode,
		Participants: partsFinal,
	}

	summary, err := r.coreClient.ConsolidateMatch(context.Background(), consolidateReq)
	if err != nil {
		log.Printf("⚠️ Error al consolidar partida %s con core-service: %v", r.roomCode, err)
	}

	if summary == nil {
		podium := r.buildLeaderboard()
		var winnerID *string
		if len(podium) > 0 {
			winnerID = &podium[0].IDUser
		}
		summary = &models.MatchSummary{
			IDMatch:      r.pkg.IDMatch,
			RoomCode:     r.roomCode,
			IDDeck:       r.pkg.IDDeck,
			DeckTitle:    r.pkg.DeckTitle,
			WinnerUserID: winnerID,
			Status:       models.MatchStatusFinished,
			TotalRounds:  r.TotalRounds,
			RoundTimeSec: r.RoundTimeSec,
			Podium:       podium,
			FinishedAt:   time.Now().UTC().Format(time.RFC3339),
		}
	}

	r.broadcast(models.WSEvent{
		Type:    "PODIUM",
		Payload: summary,
	})

	r.hub.RemoveRoom(r.roomCode)
	r.stop()
}

func (r *ActiveRoom) recalculateRanks() {
	parts := make([]*roomParticipant, 0, len(r.participants))
	for _, p := range r.participants {
		parts = append(parts, p)
	}

	sort.Slice(parts, func(i, j int) bool {
		a, b := parts[i], parts[j]
		if a.Score != b.Score {
			return a.Score > b.Score
		}
		if len(a.CorrectSenseIDs) != len(b.CorrectSenseIDs) {
			return len(a.CorrectSenseIDs) > len(b.CorrectSenseIDs)
		}
		if a.TotalTimeSpentMs != b.TotalTimeSpentMs {
			return a.TotalTimeSpentMs < b.TotalTimeSpentMs
		}
		return a.JoinedAt.Before(b.JoinedAt)
	})

	for idx, p := range parts {
		p.Rank = idx + 1
	}
}

func (r *ActiveRoom) buildLeaderboard() []models.ArenaParticipant {
	r.recalculateRanks()

	parts := make([]*roomParticipant, 0, len(r.participants))
	for _, p := range r.participants {
		parts = append(parts, p)
	}

	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Rank < parts[j].Rank
	})

	res := make([]models.ArenaParticipant, len(parts))
	for i, p := range parts {
		res[i] = models.ArenaParticipant{
			IDUser:      p.IDUser,
			UserName:    p.UserName,
			Name:        p.Name,
			AvatarURL:   p.AvatarURL,
			Score:       p.Score,
			Rank:        p.Rank,
			IsHost:      p.IsHost,
			IsReady:     p.Connected,
			DialectCode: p.DialectCode,
			FlagEmoji:   p.FlagEmoji,
		}
	}
	return res
}

func (r *ActiveRoom) buildRoomState() models.ArenaRoomState {
	leaderboard := r.buildLeaderboard()
	return models.ArenaRoomState{
		RoomCode:     r.roomCode,
		DeckID:       r.pkg.IDDeck,
		DeckTitle:    r.pkg.DeckTitle,
		Status:       r.Status,
		TotalRounds:  r.TotalRounds,
		CurrentRound: r.CurrentRound,
		RoundTimeSec: r.RoundTimeSec,
		Participants: leaderboard,
	}
}

func (r *ActiveRoom) broadcastRoomState() {
	state := r.buildRoomState()
	r.broadcast(models.WSEvent{
		Type:    "ROOM_STATE",
		Payload: state,
	})
}

func (r *ActiveRoom) broadcast(event models.WSEvent) {
	for _, p := range r.participants {
		if p.Connected && p.client != nil {
			_ = p.client.SendJSON(event)
		}
	}
}
