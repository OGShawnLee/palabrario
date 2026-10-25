package arena

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

func readNextEvent(t *testing.T, ws *websocket.Conn, timeout time.Duration) models.WSEvent {
	t.Helper()
	_ = ws.SetReadDeadline(time.Now().Add(timeout))
	var ev models.WSEvent
	err := ws.ReadJSON(&ev)
	if err != nil {
		t.Fatalf("error leyendo evento de websocket: %v", err)
	}
	return ev
}

func waitForEventType(t *testing.T, ws *websocket.Conn, targetType string, timeout time.Duration) models.WSEvent {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		_ = ws.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
		var ev models.WSEvent
		if err := ws.ReadJSON(&ev); err == nil {
			if ev.Type == targetType {
				return ev
			}
		}
	}
	t.Fatalf("tiempo de espera agotado esperando evento de tipo: %s", targetType)
	return models.WSEvent{}
}

func waitForRoomStateCondition(t *testing.T, ws *websocket.Conn, predicate func(state models.ArenaRoomState) bool, timeout time.Duration) models.ArenaRoomState {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		_ = ws.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
		var ev models.WSEvent
		if err := ws.ReadJSON(&ev); err == nil {
			if ev.Type == "ROOM_STATE" {
				stateBytes, _ := json.Marshal(ev.Payload)
				var state models.ArenaRoomState
				if err := json.Unmarshal(stateBytes, &state); err == nil {
					if predicate(state) {
						return state
					}
				}
			}
		}
	}
	t.Fatalf("tiempo de espera agotado esperando condición en ROOM_STATE")
	return models.ArenaRoomState{}
}

// 1. TestCalculateRoundPoints_RN10
func TestCalculateRoundPoints_RN10(t *testing.T) {
	roundTime := 15 // 15000 ms

	// 0 ms => 1000 pts
	if pts := CalculateRoundPoints(true, 0, roundTime); pts != 1000 {
		t.Errorf("esperado 1000 pts a 0ms, obtenido %d", pts)
	}

	// ms negativos => acotado a 0ms => 1000 pts
	if pts := CalculateRoundPoints(true, -100, roundTime); pts != 1000 {
		t.Errorf("esperado 1000 pts con ms negativos, obtenido %d", pts)
	}

	// Mitad de tiempo (7500 ms) => 750 pts
	if pts := CalculateRoundPoints(true, 7500, roundTime); pts != 750 {
		t.Errorf("esperado 750 pts a mitad de tiempo, obtenido %d", pts)
	}

	// Límite de tiempo exacto (15000 ms) => 500 pts
	if pts := CalculateRoundPoints(true, 15000, roundTime); pts != 500 {
		t.Errorf("esperado 500 pts al límite de tiempo, obtenido %d", pts)
	}

	// Pasado el límite (15001 ms) => 0 pts (timeout)
	if pts := CalculateRoundPoints(true, 15001, roundTime); pts != 0 {
		t.Errorf("esperado 0 pts en timeout, obtenido %d", pts)
	}

	// Respuesta incorrecta => 0 pts
	if pts := CalculateRoundPoints(false, 500, roundTime); pts != 0 {
		t.Errorf("esperado 0 pts con respuesta incorrecta, obtenido %d", pts)
	}

	// roundTimeSec inválido => 0 pts
	if pts := CalculateRoundPoints(true, 500, 0); pts != 0 {
		t.Errorf("esperado 0 pts con tiempo de ronda <= 0, obtenido %d", pts)
	}
}

// 2. TestActiveRoom_Max8Participants
func TestActiveRoom_Max8Participants(t *testing.T) {
	jwtSecret := "test-secret-key-12345"
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

	// Mock core-service
	mockCore := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/package") {
			authHeader := r.Header.Get("Authorization")
			token := strings.TrimPrefix(authHeader, "Bearer ")
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

	// Router para realtime-service
	r := chi.NewRouter()
	r.Get("/ws/arena/{roomCode}", func(w http.ResponseWriter, req *http.Request) {
		tokenStr := req.URL.Query().Get("token")
		claims, err := tokenService.VerifyToken(tokenStr)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Token inválido")
			return
		}
		hub.HandleWS(w, req, chi.URLParam(req, "roomCode"), claims)
	})

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/arena/MX8888"

	conns := make([]*websocket.Conn, 8)
	// Conectar 8 jugadores
	for i := 1; i <= 8; i++ {
		userID := "user-" + string(rune('0'+i))
		token := generateTestUserToken(tokenService, userID, "Player"+string(rune('0'+i)), "MX")
		ws, _, err := websocket.DefaultDialer.Dial(wsURL+"?token="+token, nil)
		if err != nil {
			t.Fatalf("error conectando jugador %d: %v", i, err)
		}
		conns[i-1] = ws
		_ = readNextEvent(t, ws, 1*time.Second) // ROOM_STATE inicial
	}

	defer func() {
		for _, ws := range conns {
			if ws != nil {
				_ = ws.Close()
			}
		}
	}()

	// Intentar conectar el 9° jugador
	token9 := generateTestUserToken(tokenService, "user-9", "Player9", "MX")
	ws9, _, err := websocket.DefaultDialer.Dial(wsURL+"?token="+token9, nil)
	if err == nil {
		defer ws9.Close()
		ev := readNextEvent(t, ws9, 1*time.Second)
		if ev.Type != "ERROR" {
			t.Fatalf("esperado evento ERROR para el 9° jugador, recibido: %s", ev.Type)
		}
		payloadBytes, _ := json.Marshal(ev.Payload)
		if !strings.Contains(string(payloadBytes), "La sala está llena (máximo 8 jugadores)") {
			t.Fatalf("mensaje de error inesperado: %s", string(payloadBytes))
		}
	}
}

// 3. TestActiveRoom_ReconnectWithinWindow_And_HostSuccession
func TestActiveRoom_ReconnectWithinWindow_And_HostSuccession(t *testing.T) {
	jwtSecret := "test-secret-key-succession"
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
	hub.ReconnectWindow = 80 * time.Millisecond
	hub.CountdownDuration = 20 * time.Millisecond

	r := chi.NewRouter()
	r.Get("/ws/arena/{roomCode}", func(w http.ResponseWriter, req *http.Request) {
		tokenStr := req.URL.Query().Get("token")
		claims, err := tokenService.VerifyToken(tokenStr)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Token inválido")
			return
		}
		hub.HandleWS(w, req, chi.URLParam(req, "roomCode"), claims)
	})

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/arena/MX7777"

	// Conectar Host
	tokenHost := generateTestUserToken(tokenService, "user-host", "HostUser", "MX")
	wsHost, _, err := websocket.DefaultDialer.Dial(wsURL+"?token="+tokenHost, nil)
	if err != nil {
		t.Fatalf("error conectando host: %v", err)
	}
	_ = readNextEvent(t, wsHost, 1*time.Second)

	time.Sleep(10 * time.Millisecond)

	// Conectar Jugador 2
	tokenP2 := generateTestUserToken(tokenService, "user-p2", "Player2", "ES")
	wsP2, _, err := websocket.DefaultDialer.Dial(wsURL+"?token="+tokenP2, nil)
	if err != nil {
		t.Fatalf("error conectando p2: %v", err)
	}
	_ = readNextEvent(t, wsP2, 1*time.Second)

	// Parte A: Jugador 2 se desconecta y reconecta antes de ReconnectWindow (80ms)
	_ = wsP2.Close()
	time.Sleep(20 * time.Millisecond) // 20ms < 80ms

	wsP2Reconnected, _, err := websocket.DefaultDialer.Dial(wsURL+"?token="+tokenP2, nil)
	if err != nil {
		t.Fatalf("error reconectando p2: %v", err)
	}
	defer wsP2Reconnected.Close()

	// Verifica que al reconectar recibe ROOM_STATE con Connected = true
	reconnectState := waitForRoomStateCondition(t, wsP2Reconnected, func(state models.ArenaRoomState) bool {
		for _, p := range state.Participants {
			if p.IDUser == "user-p2" && p.IsReady {
				return true
			}
		}
		return false
	}, 1*time.Second)

	if len(reconnectState.Participants) < 2 {
		t.Fatalf("esperados al menos 2 participantes tras reconexión")
	}

	// Parte B: El Host se desconecta y expira ReconnectWindow (80ms)
	_ = wsHost.Close()

	// Esperamos a que Player 2 reciba la actualización de ROOM_STATE con IsHost = true
	successionState := waitForRoomStateCondition(t, wsP2Reconnected, func(state models.ArenaRoomState) bool {
		for _, p := range state.Participants {
			if p.IDUser == "user-p2" && p.IsHost {
				return true
			}
		}
		return false
	}, 2*time.Second)

	// Verificar que el host antiguo ya no está en los participantes
	for _, p := range successionState.Participants {
		if p.IDUser == "user-host" {
			t.Fatalf("el host desconectado no debió figurar en participantes tras expirar la ventana")
		}
	}

	// Player 2 como nuevo Host ahora dispara START_MATCH
	startEvent := models.WSEvent{
		Type: "START_MATCH",
	}
	if err := wsP2Reconnected.WriteJSON(startEvent); err != nil {
		t.Fatalf("error enviando START_MATCH como nuevo host: %v", err)
	}

	// Debe recibir COUNTDOWN
	evCountdown := waitForEventType(t, wsP2Reconnected, "COUNTDOWN", 1*time.Second)
	if evCountdown.Type != "COUNTDOWN" {
		t.Fatalf("esperado evento COUNTDOWN tras START_MATCH, recibido: %s", evCountdown.Type)
	}
}

// 4. TestActiveRoom_FullMatchFlow_And_BatchFlushS2S
func TestActiveRoom_FullMatchFlow_And_BatchFlushS2S(t *testing.T) {
	jwtSecret := "test-secret-key-fullflow"
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
			authHeader := r.Header.Get("Authorization")
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if _, err := tokenService.VerifyInternalToken(token); err != nil {
				http.Error(w, "Unauthorized S2S", http.StatusUnauthorized)
				return
			}
			response.JSON(w, http.StatusOK, mockMatchPkg)
			return
		}

		if strings.HasSuffix(r.URL.Path, "/consolidate") {
			authHeader := r.Header.Get("Authorization")
			token := strings.TrimPrefix(authHeader, "Bearer ")
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
	hub.ReconnectWindow = 100 * time.Millisecond

	r := chi.NewRouter()
	r.Get("/ws/arena/{roomCode}", func(w http.ResponseWriter, req *http.Request) {
		tokenStr := req.URL.Query().Get("token")
		claims, err := tokenService.VerifyToken(tokenStr)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Token inválido")
			return
		}
		hub.HandleWS(w, req, chi.URLParam(req, "roomCode"), claims)
	})

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/arena/MX9999"

	tokenHost := generateTestUserToken(tokenService, "user-host", "HostPlayer", "MX")
	wsHost, _, err := websocket.DefaultDialer.Dial(wsURL+"?token="+tokenHost, nil)
	if err != nil {
		t.Fatalf("error conectando host: %v", err)
	}
	defer wsHost.Close()
	_ = readNextEvent(t, wsHost, 1*time.Second)

	tokenP2 := generateTestUserToken(tokenService, "user-p2", "PlayerDos", "CO")
	wsP2, _, err := websocket.DefaultDialer.Dial(wsURL+"?token="+tokenP2, nil)
	if err != nil {
		t.Fatalf("error conectando p2: %v", err)
	}
	defer wsP2.Close()
	_ = readNextEvent(t, wsP2, 1*time.Second)

	// Host inicia la partida
	if err := wsHost.WriteJSON(models.WSEvent{Type: "START_MATCH"}); err != nil {
		t.Fatalf("error enviando START_MATCH: %v", err)
	}

	// Ambos reciben COUNTDOWN
	_ = waitForEventType(t, wsHost, "COUNTDOWN", 1*time.Second)
	_ = waitForEventType(t, wsP2, "COUNTDOWN", 1*time.Second)

	// Ronda 1: Ambos reciben QUESTION
	q1Host := waitForEventType(t, wsHost, "QUESTION", 1*time.Second)
	_ = waitForEventType(t, wsP2, "QUESTION", 1*time.Second)

	// Verificar que QUESTION no revela CorrectOption
	q1Bytes, _ := json.Marshal(q1Host.Payload)
	if strings.Contains(string(q1Bytes), "correct_option") {
		t.Fatalf("el evento QUESTION no debe contener el campo correct_option: %s", string(q1Bytes))
	}

	// Pequeña pausa para asegurar elapsedMs >= 5ms
	time.Sleep(5 * time.Millisecond)

	// Enviar respuestas para Ronda 1
	_ = wsHost.WriteJSON(models.WSEvent{
		Type: "SUBMIT_ANSWER",
		Payload: models.SubmitAnswerPayload{
			RoundNumber:    1,
			SelectedOption: "Pajilla", // Correcta
		},
	})
	_ = wsP2.WriteJSON(models.WSEvent{
		Type: "SUBMIT_ANSWER",
		Payload: models.SubmitAnswerPayload{
			RoundNumber:    1,
			SelectedOption: "Vaso", // Incorrecta
		},
	})

	// Ronda 1: Ambos reciben ROUND_RESULT
	rr1Host := waitForEventType(t, wsHost, "ROUND_RESULT", 1*time.Second)
	_ = waitForEventType(t, wsP2, "ROUND_RESULT", 1*time.Second)
	rr1Bytes, _ := json.Marshal(rr1Host.Payload)
	if !strings.Contains(string(rr1Bytes), "Pajilla") {
		t.Fatalf("ROUND_RESULT debe contener la opción correcta: %s", string(rr1Bytes))
	}

	// Ronda 2: Ambos reciben QUESTION
	_ = waitForEventType(t, wsHost, "QUESTION", 1*time.Second)
	_ = waitForEventType(t, wsP2, "QUESTION", 1*time.Second)

	time.Sleep(5 * time.Millisecond)

	// Enviar respuestas para Ronda 2
	_ = wsHost.WriteJSON(models.WSEvent{
		Type: "SUBMIT_ANSWER",
		Payload: models.SubmitAnswerPayload{
			RoundNumber:    2,
			SelectedOption: "Chaqueta", // Correcta
		},
	})
	_ = wsP2.WriteJSON(models.WSEvent{
		Type: "SUBMIT_ANSWER",
		Payload: models.SubmitAnswerPayload{
			RoundNumber:    2,
			SelectedOption: "Chaqueta", // Correcta
		},
	})

	// Ronda 2: Ambos reciben ROUND_RESULT
	_ = waitForEventType(t, wsHost, "ROUND_RESULT", 1*time.Second)
	_ = waitForEventType(t, wsP2, "ROUND_RESULT", 1*time.Second)

	// Al terminar la ronda 2, tras InterRoundPause, ambos reciben PODIUM
	podiumHost := waitForEventType(t, wsHost, "PODIUM", 2*time.Second)
	_ = waitForEventType(t, wsP2, "PODIUM", 2*time.Second)

	if podiumHost.Type != "PODIUM" {
		t.Fatalf("esperado evento PODIUM al finalizar partida")
	}

	// Verificar consolidación S2S capturada en el mock
	mu.Lock()
	req := capturedConsolidateReq
	tok := capturedS2SToken
	mu.Unlock()

	if req == nil {
		t.Fatalf("no se recibió la consolidación en core-service")
	}

	if tok == "" {
		t.Fatalf("no se capturó el token S2S en la consolidación")
	}

	if req.RoomCode != "MX9999" || req.IDMatch != "match-flow-123" {
		t.Fatalf("RoomCode o IDMatch inválidos en consolidación: %+v", req)
	}

	if len(req.Participants) != 2 {
		t.Fatalf("esperados 2 participantes en la consolidación, obtenidos %d", len(req.Participants))
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
		t.Fatalf("participantes no encontrados en la consolidación")
	}

	// Host acertó 2/2
	if len(hostPart.CorrectSenseIDs) != 2 || len(hostPart.WrongSenseIDs) != 0 {
		t.Errorf("host debió tener 2 aciertos y 0 fallos, obtenidos: aciertos=%v, fallos=%v", hostPart.CorrectSenseIDs, hostPart.WrongSenseIDs)
	}
	if len(hostPart.RoundAttempts) != 2 {
		t.Errorf("host debió tener 2 RoundAttempts, obtenidos: %d", len(hostPart.RoundAttempts))
	}
	if hostPart.TotalTimeSpentMs <= 0 {
		t.Errorf("host TotalTimeSpentMs debe ser > 0, obtenido %d", hostPart.TotalTimeSpentMs)
	}
	if hostPart.FinalRank != 1 {
		t.Errorf("host debió terminar en Rank 1, obtenido: %d", hostPart.FinalRank)
	}

	// Player 2 acertó 1/2
	if len(p2Part.CorrectSenseIDs) != 1 || len(p2Part.WrongSenseIDs) != 1 {
		t.Errorf("p2 debió tener 1 acierto y 1 fallo, obtenidos: aciertos=%v, fallos=%v", p2Part.CorrectSenseIDs, p2Part.WrongSenseIDs)
	}
	if len(p2Part.RoundAttempts) != 2 {
		t.Errorf("p2 debió tener 2 RoundAttempts, obtenidos: %d", len(p2Part.RoundAttempts))
	}
	if p2Part.FinalRank != 2 {
		t.Errorf("p2 debió terminar en Rank 2, obtenido: %d", p2Part.FinalRank)
	}

	// Sala debe haber sido removida del ArenaHub
	time.Sleep(50 * time.Millisecond)
	if count := hub.ActiveRoomsCount(); count != 0 {
		t.Errorf("esperadas 0 salas activas en RAM tras el PODIUM, encontradas %d", count)
	}
}
