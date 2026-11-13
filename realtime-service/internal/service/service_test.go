package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/response"
)

func generateTestUserToken(ts *auth.TokenService, userID, userName, dialect string) string {
	user := &models.UserProfile{
		IDUser:        userID,
		UserName:      userName,
		Email:         userName + "@test.com",
		AnchorDialect: dialect,
		Role:          models.RoleLearner,
	}
	tok, _ := ts.GenerateToken(user)
	return tok
}

type wsReadResult struct {
	event models.WSEvent
	err   error
}

// wsReader owns a single dedicated reader goroutine per WebSocket connection so
// read timeouts never put gorilla/websocket into a permanently failed state.
type wsReader struct {
	conn *websocket.Conn
	ch   chan wsReadResult
}

func newWSReader(ws *websocket.Conn) *wsReader {
	r := &wsReader{
		conn: ws,
		ch:   make(chan wsReadResult, 64),
	}
	go func() {
		defer close(r.ch)
		for {
			var ev models.WSEvent
			err := ws.ReadJSON(&ev)
			r.ch <- wsReadResult{event: ev, err: err}
			if err != nil {
				return
			}
		}
	}()
	return r
}

func readNextEvent(t *testing.T, r *wsReader, timeout time.Duration) models.WSEvent {
	t.Helper()
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case res, ok := <-r.ch:
		if !ok {
			t.Fatal("websocket reader channel closed unexpectedly")
		}
		if res.err != nil {
			t.Fatalf("error reading websocket event: %v", res.err)
		}
		return res.event
	case <-timer.C:
		t.Fatalf("timed out after %v waiting for next websocket event", timeout)
		return models.WSEvent{}
	}
}

func waitForEventType(t *testing.T, r *wsReader, targetType string, timeout time.Duration) models.WSEvent {
	t.Helper()
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case res, ok := <-r.ch:
			if !ok {
				t.Fatalf("websocket reader closed while waiting for event %s", targetType)
			}
			if res.err != nil {
				t.Fatalf("websocket read error while waiting for event %s: %v", targetType, res.err)
			}
			if res.event.Type == targetType {
				return res.event
			}
		case <-timer.C:
			t.Fatalf("timed out after %v waiting for event type: %s", timeout, targetType)
			return models.WSEvent{}
		}
	}
}

func waitForRoomStateCondition(t *testing.T, r *wsReader, predicate func(state models.ArenaRoomState) bool, timeout time.Duration) models.ArenaRoomState {
	t.Helper()
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case res, ok := <-r.ch:
			if !ok {
				t.Fatal("websocket reader closed while waiting for ROOM_STATE condition")
			}
			if res.err != nil {
				t.Fatalf("websocket read error while waiting for ROOM_STATE condition: %v", res.err)
			}
			if res.event.Type == "ROOM_STATE" {
				stateBytes, _ := json.Marshal(res.event.Payload)
				var state models.ArenaRoomState
				if err := json.Unmarshal(stateBytes, &state); err == nil {
					if predicate(state) {
						return state
					}
				}
			}
		case <-timer.C:
			t.Fatalf("timed out after %v waiting for ROOM_STATE condition", timeout)
			return models.ArenaRoomState{}
		}
	}
}

func eventually(t *testing.T, timeout, interval time.Duration, condition func() bool, failureMsg string) {
	t.Helper()
	if condition() {
		return
	}
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if condition() {
				return
			}
		case <-deadline.C:
			t.Fatal(failureMsg)
			return
		}
	}
}

func newTestWSHandler(tokenService *auth.TokenService, hub *ArenaHub) http.Handler {
	r := chi.NewRouter()
	r.Get("/ws/arena/{roomCode}", func(w http.ResponseWriter, req *http.Request) {
		tokenStr, present, err := auth.ExtractBearerToken(req)
		if present {
			if err != nil {
				response.ErrorWithCode(w, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Invalid authorization header format")
				return
			}
		} else {
			tokenStr = strings.TrimSpace(req.URL.Query().Get("token"))
		}
		claims, err := tokenService.VerifyToken(tokenStr)
		if err != nil {
			response.ErrorWithCode(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
			return
		}
		hub.HandleWS(w, req, chi.URLParam(req, "roomCode"), claims)
	})
	return r
}

// 1. TestCalculateRoundPoints_RN10
func TestCalculateRoundPoints_RN10(t *testing.T) {
	roundTime := 15 // 15000 ms

	// 0 ms => 1000 pts
	if pts := CalculateRoundPoints(true, 0, roundTime); pts != 1000 {
		t.Errorf("expected 1000 pts at 0ms, got %d", pts)
	}

	// negative ms => clamped to 0ms => 1000 pts
	if pts := CalculateRoundPoints(true, -100, roundTime); pts != 1000 {
		t.Errorf("expected 1000 pts with negative ms, got %d", pts)
	}

	// half time (7500 ms) => 750 pts
	if pts := CalculateRoundPoints(true, 7500, roundTime); pts != 750 {
		t.Errorf("expected 750 pts at half time, got %d", pts)
	}

	// exact deadline (15000 ms) => 500 pts
	if pts := CalculateRoundPoints(true, 15000, roundTime); pts != 500 {
		t.Errorf("expected 500 pts at time limit, got %d", pts)
	}

	// past limit (15001 ms) => 0 pts (timeout)
	if pts := CalculateRoundPoints(true, 15001, roundTime); pts != 0 {
		t.Errorf("expected 0 pts on timeout, got %d", pts)
	}

	// wrong answer => 0 pts
	if pts := CalculateRoundPoints(false, 500, roundTime); pts != 0 {
		t.Errorf("expected 0 pts for wrong answer, got %d", pts)
	}

	// invalid roundTimeSec <= 0 => 0 pts
	if pts := CalculateRoundPoints(true, 500, 0); pts != 0 {
		t.Errorf("expected 0 pts when roundTimeSec <= 0, got %d", pts)
	}
}

// 2. TestActiveRoom_Max8Participants
func TestActiveRoom_Max8Participants(t *testing.T) {
	jwtSecret := "test-jwt-secret-key-palabrario-minimum-32-chars"
	tokenService := auth.NewTokenService(jwtSecret, 1*time.Hour)

	mockMatchPkg := &models.MatchPackage{
		IDMatch:      "match-max-8",
		RoomCode:     "MX8888",
		HostUserID:   "user-1",
		IDDeck:       1,
		DeckTitle:    "Baraja Prueba 8",
		TotalRounds:  2,
		RoundTimeSec: 10,
		Questions: []models.InternalArenaQuestion{
			{
				RoundNumber:   1,
				IDSense:       101,
				Prompt:        "Pregunta 1",
				Options:       []string{"A", "B", "C", "D"},
				CorrectOption: "A",
				TimeLimitSecs: 10,
			},
		},
	}

	mockCore := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/package") {
			token, _, err := auth.ExtractBearerToken(r)
			if err != nil {
				http.Error(w, "Unauthorized S2S", http.StatusUnauthorized)
				return
			}
			if _, err := tokenService.VerifyInternalToken(token); err != nil {
				http.Error(w, "Unauthorized S2S", http.StatusUnauthorized)
				return
			}
			response.JSON(w, http.StatusOK, mockMatchPkg)
			return
		}
		http.NotFound(w, r)
	}))
	defer mockCore.Close()

	coreClient := NewCoreClient(mockCore.URL, mockCore.Client(), tokenService)
	hub := NewArenaHub(coreClient, tokenService)

	server := httptest.NewServer(newTestWSHandler(tokenService, hub))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/arena/MX8888"

	conns := make([]*websocket.Conn, 8)
	for i := 1; i <= 8; i++ {
		userID := "user-" + string(rune('0'+i))
		token := generateTestUserToken(tokenService, userID, "Player"+string(rune('0'+i)), "MX")
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, http.Header{"Authorization": {"Bearer " + token}})
		if err != nil {
			t.Fatalf("error connecting player %d: %v", i, err)
		}
		conns[i-1] = ws
		r := newWSReader(ws)
		_ = readNextEvent(t, r, 2*time.Second) // initial ROOM_STATE
	}

	defer func() {
		for _, ws := range conns {
			if ws != nil {
				_ = ws.Close()
			}
		}
	}()

	// Attempt to connect the 9th player
	token9 := generateTestUserToken(tokenService, "user-9", "Player9", "MX")
	ws9, _, err := websocket.DefaultDialer.Dial(wsURL, http.Header{"Authorization": {"Bearer " + token9}})
	if err == nil {
		defer ws9.Close()
		r9 := newWSReader(ws9)
		ev := readNextEvent(t, r9, 2*time.Second)
		if ev.Type != "ERROR" {
			t.Fatalf("expected ERROR event for 9th player, got: %s", ev.Type)
		}
		payloadBytes, _ := json.Marshal(ev.Payload)
		if !strings.Contains(string(payloadBytes), "La sala está llena (máximo 8 jugadores)") {
			t.Fatalf("unexpected error payload: %s", string(payloadBytes))
		}
	}
}

// 3. TestActiveRoom_ReconnectWithinWindow_And_HostSuccession
func TestActiveRoom_ReconnectWithinWindow_And_HostSuccession(t *testing.T) {
	jwtSecret := "test-jwt-secret-key-palabrario-minimum-32-chars"
	tokenService := auth.NewTokenService(jwtSecret, 1*time.Hour)

	mockMatchPkg := &models.MatchPackage{
		IDMatch:      "match-succession",
		RoomCode:     "MX7777",
		HostUserID:   "user-host",
		IDDeck:       1,
		DeckTitle:    "Baraja Sucesión",
		TotalRounds:  2,
		RoundTimeSec: 10,
		Questions: []models.InternalArenaQuestion{
			{
				RoundNumber:   1,
				IDSense:       101,
				Prompt:        "Pregunta 1",
				Options:       []string{"A", "B", "C", "D"},
				CorrectOption: "A",
				TimeLimitSecs: 10,
			},
		},
	}

	mockCore := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/package") {
			response.JSON(w, http.StatusOK, mockMatchPkg)
			return
		}
		http.NotFound(w, r)
	}))
	defer mockCore.Close()

	coreClient := NewCoreClient(mockCore.URL, mockCore.Client(), tokenService)
	hub := NewArenaHub(coreClient, tokenService)
	hub.ReconnectWindow = 250 * time.Millisecond
	hub.CountdownDuration = 20 * time.Millisecond

	server := httptest.NewServer(newTestWSHandler(tokenService, hub))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/arena/MX7777"

	// Connect Host and confirm initial ROOM_STATE
	tokenHost := generateTestUserToken(tokenService, "user-host", "HostUser", "MX")
	wsHost, _, err := websocket.DefaultDialer.Dial(wsURL, http.Header{"Authorization": {"Bearer " + tokenHost}})
	if err != nil {
		t.Fatalf("error connecting host: %v", err)
	}
	defer wsHost.Close()
	rHost := newWSReader(wsHost)
	_ = readNextEvent(t, rHost, 2*time.Second)

	// Connect Player 2 and confirm initial ROOM_STATE
	tokenP2 := generateTestUserToken(tokenService, "user-p2", "Player2", "ES")
	wsP2, _, err := websocket.DefaultDialer.Dial(wsURL, http.Header{"Authorization": {"Bearer " + tokenP2}})
	if err != nil {
		t.Fatalf("error connecting p2: %v", err)
	}
	rP2 := newWSReader(wsP2)
	_ = readNextEvent(t, rP2, 2*time.Second)

	// Part A: Player 2 disconnects; wait deterministically on Host for IsReady=false broadcast,
	// then reconnect immediately within the 250ms ReconnectWindow.
	_ = wsP2.Close()
	_ = waitForRoomStateCondition(t, rHost, func(state models.ArenaRoomState) bool {
		for _, p := range state.Participants {
			if p.IDUser == "user-p2" && !p.IsReady {
				return true
			}
		}
		return false
	}, 2*time.Second)

	wsP2Reconnected, _, err := websocket.DefaultDialer.Dial(wsURL, http.Header{"Authorization": {"Bearer " + tokenP2}})
	if err != nil {
		t.Fatalf("error reconnecting p2: %v", err)
	}
	defer wsP2Reconnected.Close()
	rP2Reconnected := newWSReader(wsP2Reconnected)

	reconnectState := waitForRoomStateCondition(t, rP2Reconnected, func(state models.ArenaRoomState) bool {
		for _, p := range state.Participants {
			if p.IDUser == "user-p2" && p.IsReady {
				return true
			}
		}
		return false
	}, 2*time.Second)

	if len(reconnectState.Participants) < 2 {
		t.Fatalf("expected at least 2 participants after reconnection, got %d", len(reconnectState.Participants))
	}

	// Part B: Host disconnects and ReconnectWindow (250ms) expires
	_ = wsHost.Close()

	// Wait for Player 2 to receive ROOM_STATE promotion with IsHost = true
	successionState := waitForRoomStateCondition(t, rP2Reconnected, func(state models.ArenaRoomState) bool {
		for _, p := range state.Participants {
			if p.IDUser == "user-p2" && p.IsHost {
				return true
			}
		}
		return false
	}, 2*time.Second)

	for _, p := range successionState.Participants {
		if p.IDUser == "user-host" {
			t.Fatalf("disconnected host should not remain in participants after ReconnectWindow expired")
		}
	}

	// Player 2 as the new Host triggers START_MATCH
	if err := wsP2Reconnected.WriteJSON(models.WSEvent{Type: "START_MATCH"}); err != nil {
		t.Fatalf("error sending START_MATCH as new host: %v", err)
	}

	evCountdown := waitForEventType(t, rP2Reconnected, "COUNTDOWN", 2*time.Second)
	if evCountdown.Type != "COUNTDOWN" {
		t.Fatalf("expected COUNTDOWN event after START_MATCH, got: %s", evCountdown.Type)
	}
}

// 4. TestActiveRoom_FullMatchFlow_And_BatchFlushS2S
func TestActiveRoom_FullMatchFlow_And_BatchFlushS2S(t *testing.T) {
	jwtSecret := "test-jwt-secret-key-palabrario-minimum-32-chars"
	tokenService := auth.NewTokenService(jwtSecret, 1*time.Hour)

	mockMatchPkg := &models.MatchPackage{
		IDMatch:      "match-flow-123",
		RoomCode:     "MX9999",
		HostUserID:   "user-host",
		IDDeck:       1,
		DeckTitle:    "Baraja Flujo Completo",
		TotalRounds:  2,
		RoundTimeSec: 2,
		Questions: []models.InternalArenaQuestion{
			{
				RoundNumber:   1,
				IDSense:       101,
				Prompt:        "¿Qué es popote?",
				Options:       []string{"Pajilla", "Vaso", "Plato", "Cuchara"},
				CorrectOption: "Pajilla",
				Explanation:   "En México popote significa pajilla.",
				TimeLimitSecs: 2,
			},
			{
				RoundNumber:   2,
				IDSense:       102,
				Prompt:        "¿Qué es chamarra?",
				Options:       []string{"Chaqueta", "Pantalón", "Sombrero", "Zapato"},
				CorrectOption: "Chaqueta",
				Explanation:   "En México chamarra significa chaqueta.",
				TimeLimitSecs: 2,
			},
		},
	}

	var mu sync.Mutex
	var capturedConsolidateReq *models.ConsolidateMatchRequest
	var capturedS2SToken string

	mockCore := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/package") {
			token, _, err := auth.ExtractBearerToken(r)
			if err != nil {
				http.Error(w, "Unauthorized S2S", http.StatusUnauthorized)
				return
			}
			if _, err := tokenService.VerifyInternalToken(token); err != nil {
				http.Error(w, "Unauthorized S2S", http.StatusUnauthorized)
				return
			}
			response.JSON(w, http.StatusOK, mockMatchPkg)
			return
		}

		if strings.HasSuffix(r.URL.Path, "/consolidate") {
			token, _, err := auth.ExtractBearerToken(r)
			if err != nil {
				http.Error(w, "Unauthorized S2S", http.StatusUnauthorized)
				return
			}
			claims, err := tokenService.VerifyInternalToken(token)
			if err != nil {
				http.Error(w, "Unauthorized S2S", http.StatusUnauthorized)
				return
			}
			mu.Lock()
			capturedS2SToken = token
			if claims.ServiceName != "realtime-service" {
				mu.Unlock()
				http.Error(w, "Invalid service name", http.StatusBadRequest)
				return
			}
			var req models.ConsolidateMatchRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				mu.Unlock()
				http.Error(w, "Invalid body", http.StatusBadRequest)
				return
			}
			capturedConsolidateReq = &req
			mu.Unlock()

			summary := models.MatchSummary{
				IDMatch:      req.IDMatch,
				RoomCode:     req.RoomCode,
				IDDeck:       1,
				DeckTitle:    "Baraja Flujo Completo",
				Status:       models.MatchStatusFinished,
				TotalRounds:  2,
				RoundTimeSec: 2,
				FinishedAt:   time.Now().UTC().Format(time.RFC3339),
			}
			response.JSON(w, http.StatusOK, summary)
			return
		}

		http.NotFound(w, r)
	}))
	defer mockCore.Close()

	coreClient := NewCoreClient(mockCore.URL, mockCore.Client(), tokenService)
	hub := NewArenaHub(coreClient, tokenService)
	hub.CountdownDuration = 20 * time.Millisecond
	hub.InterRoundPause = 30 * time.Millisecond
	hub.ReconnectWindow = 250 * time.Millisecond

	server := httptest.NewServer(newTestWSHandler(tokenService, hub))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/arena/MX9999"

	tokenHost := generateTestUserToken(tokenService, "user-host", "HostPlayer", "MX")
	wsHost, _, err := websocket.DefaultDialer.Dial(wsURL, http.Header{"Authorization": {"Bearer " + tokenHost}})
	if err != nil {
		t.Fatalf("error connecting host: %v", err)
	}
	defer wsHost.Close()
	rHost := newWSReader(wsHost)
	_ = readNextEvent(t, rHost, 2*time.Second)

	tokenP2 := generateTestUserToken(tokenService, "user-p2", "PlayerDos", "CO")
	wsP2, _, err := websocket.DefaultDialer.Dial(wsURL, http.Header{"Authorization": {"Bearer " + tokenP2}})
	if err != nil {
		t.Fatalf("error connecting p2: %v", err)
	}
	defer wsP2.Close()
	rP2 := newWSReader(wsP2)
	_ = readNextEvent(t, rP2, 2*time.Second)

	// Host starts match
	if err := wsHost.WriteJSON(models.WSEvent{Type: "START_MATCH"}); err != nil {
		t.Fatalf("error sending START_MATCH: %v", err)
	}

	// Both receive COUNTDOWN
	_ = waitForEventType(t, rHost, "COUNTDOWN", 2*time.Second)
	_ = waitForEventType(t, rP2, "COUNTDOWN", 2*time.Second)

	// Round 1: both receive QUESTION
	q1Host := waitForEventType(t, rHost, "QUESTION", 2*time.Second)
	_ = waitForEventType(t, rP2, "QUESTION", 2*time.Second)

	// QUESTION must never leak correct_option
	q1Bytes, _ := json.Marshal(q1Host.Payload)
	if strings.Contains(string(q1Bytes), "correct_option") {
		t.Fatalf("QUESTION event must not contain correct_option: %s", string(q1Bytes))
	}

	// Submit Round 1 answers
	if err := wsHost.WriteJSON(models.WSEvent{
		Type: "SUBMIT_ANSWER",
		Payload: models.SubmitAnswerPayload{
			RoundNumber:    1,
			SelectedOption: "Pajilla",
		},
	}); err != nil {
		t.Fatalf("host submit round 1 failed: %v", err)
	}
	if err := wsP2.WriteJSON(models.WSEvent{
		Type: "SUBMIT_ANSWER",
		Payload: models.SubmitAnswerPayload{
			RoundNumber:    1,
			SelectedOption: "Vaso",
		},
	}); err != nil {
		t.Fatalf("p2 submit round 1 failed: %v", err)
	}

	// Round 1: both receive ROUND_RESULT
	rr1Host := waitForEventType(t, rHost, "ROUND_RESULT", 2*time.Second)
	_ = waitForEventType(t, rP2, "ROUND_RESULT", 2*time.Second)
	rr1Bytes, _ := json.Marshal(rr1Host.Payload)
	if !strings.Contains(string(rr1Bytes), "Pajilla") {
		t.Fatalf("ROUND_RESULT must include correct option: %s", string(rr1Bytes))
	}

	// Round 2: both receive QUESTION
	_ = waitForEventType(t, rHost, "QUESTION", 2*time.Second)
	_ = waitForEventType(t, rP2, "QUESTION", 2*time.Second)

	// Submit Round 2 answers
	if err := wsHost.WriteJSON(models.WSEvent{
		Type: "SUBMIT_ANSWER",
		Payload: models.SubmitAnswerPayload{
			RoundNumber:    2,
			SelectedOption: "Chaqueta",
		},
	}); err != nil {
		t.Fatalf("host submit round 2 failed: %v", err)
	}
	if err := wsP2.WriteJSON(models.WSEvent{
		Type: "SUBMIT_ANSWER",
		Payload: models.SubmitAnswerPayload{
			RoundNumber:    2,
			SelectedOption: "Chaqueta",
		},
	}); err != nil {
		t.Fatalf("p2 submit round 2 failed: %v", err)
	}

	// Round 2: both receive ROUND_RESULT
	_ = waitForEventType(t, rHost, "ROUND_RESULT", 2*time.Second)
	_ = waitForEventType(t, rP2, "ROUND_RESULT", 2*time.Second)

	// After InterRoundPause, both receive PODIUM
	podiumHost := waitForEventType(t, rHost, "PODIUM", 2*time.Second)
	_ = waitForEventType(t, rP2, "PODIUM", 2*time.Second)

	if podiumHost.Type != "PODIUM" {
		t.Fatal("expected PODIUM event at end of match")
	}

	// Verify captured S2S consolidation payload
	mu.Lock()
	req := capturedConsolidateReq
	tok := capturedS2SToken
	mu.Unlock()

	if req == nil {
		t.Fatal("consolidation request was not received by mock core-service")
	}
	if tok == "" {
		t.Fatal("S2S token was not captured in consolidation request")
	}
	if req.RoomCode != "MX9999" || req.IDMatch != "match-flow-123" {
		t.Fatalf("invalid RoomCode or IDMatch in consolidation: %+v", req)
	}
	if len(req.Participants) != 2 {
		t.Fatalf("expected 2 participants in consolidation, got %d", len(req.Participants))
	}

	var hostPart, p2Part *models.ParticipantFinalResult
	for i := range req.Participants {
		if req.Participants[i].IDUser == "user-host" {
			hostPart = &req.Participants[i]
		} else if req.Participants[i].IDUser == "user-p2" {
			p2Part = &req.Participants[i]
		}
	}
	if hostPart == nil || p2Part == nil {
		t.Fatal("participants missing from consolidation payload")
	}

	maxMatchMs := mockMatchPkg.TotalRounds * mockMatchPkg.RoundTimeSec * 1000
	maxRoundMs := mockMatchPkg.RoundTimeSec * 1000

	// Host answered 2/2 correctly within round window
	if len(hostPart.CorrectSenseIDs) != 2 || len(hostPart.WrongSenseIDs) != 0 {
		t.Errorf("host expected 2 correct and 0 wrong, got correct=%v wrong=%v", hostPart.CorrectSenseIDs, hostPart.WrongSenseIDs)
	}
	if len(hostPart.RoundAttempts) != 2 {
		t.Errorf("host expected 2 RoundAttempts, got %d", len(hostPart.RoundAttempts))
	}
	for _, att := range hostPart.RoundAttempts {
		if att.TimeSpentMs < 0 || att.TimeSpentMs > maxRoundMs {
			t.Errorf("host round %d TimeSpentMs out of bounds [0, %d]: %d", att.RoundNumber, maxRoundMs, att.TimeSpentMs)
		}
		if att.PointsEarned < 500 || att.PointsEarned > 1000 {
			t.Errorf("host round %d PointsEarned out of bounds [500, 1000]: %d", att.RoundNumber, att.PointsEarned)
		}
	}
	if hostPart.TotalTimeSpentMs < 0 || hostPart.TotalTimeSpentMs > maxMatchMs {
		t.Errorf("host TotalTimeSpentMs out of bounds [0, %d]: %d", maxMatchMs, hostPart.TotalTimeSpentMs)
	}
	if hostPart.FinalScore < 1000 || hostPart.FinalScore > 2000 {
		t.Errorf("host FinalScore out of bounds [1000, 2000]: %d", hostPart.FinalScore)
	}
	if hostPart.FinalRank != 1 {
		t.Errorf("host expected FinalRank 1, got %d", hostPart.FinalRank)
	}

	// Player 2 answered 1/2 correctly
	if len(p2Part.CorrectSenseIDs) != 1 || len(p2Part.WrongSenseIDs) != 1 {
		t.Errorf("p2 expected 1 correct and 1 wrong, got correct=%v wrong=%v", p2Part.CorrectSenseIDs, p2Part.WrongSenseIDs)
	}
	if len(p2Part.RoundAttempts) != 2 {
		t.Errorf("p2 expected 2 RoundAttempts, got %d", len(p2Part.RoundAttempts))
	}
	if p2Part.FinalScore < 500 || p2Part.FinalScore > 1000 {
		t.Errorf("p2 FinalScore out of bounds [500, 1000]: %d", p2Part.FinalScore)
	}
	if p2Part.FinalRank != 2 {
		t.Errorf("p2 expected FinalRank 2, got %d", p2Part.FinalRank)
	}

	// Room must be removed from ArenaHub shortly after PODIUM broadcast
	eventually(t, 2*time.Second, 10*time.Millisecond, func() bool {
		return hub.ActiveRoomsCount() == 0
	}, "expected 0 active rooms in RAM after PODIUM")
}
