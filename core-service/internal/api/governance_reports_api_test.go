package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"palabrario/pkg/models"
)

func loginHelper(t *testing.T, serverURL, email, password string) string {
	t.Helper()
	payload := models.SignInRequest{
		Email:    email,
		Password: password,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(serverURL+"/api/v1/auth/sign-in", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("login falló para %s: %v", email, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login falló con status %d para %s", resp.StatusCode, email)
	}

	var authResp models.AuthResponse
	_ = json.NewDecoder(resp.Body).Decode(&authResp)
	return authResp.Token
}

func TestUsersAPI_Lifecycle(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	gaelToken := loginHelper(t, server.URL, "gael@palabrario.org", "password123")

	// 1. PUT /api/v1/users/me
	updatePayload := models.UpdateProfileRequest{
		Name:          "Gael Modificado HTTP",
		AnchorDialect: "AR",
	}
	body, _ := json.Marshal(updatePayload)
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/v1/users/me", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+gaelToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error en PUT /api/v1/users/me: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("se esperaba 200 OK en PUT /api/v1/users/me, obtenido: %d", resp.StatusCode)
	}

	var updated models.UserProfile
	_ = json.NewDecoder(resp.Body).Decode(&updated)
	if updated.Name != "Gael Modificado HTTP" || updated.AnchorDialect != "AR" {
		t.Errorf("datos no coinciden en PUT /api/v1/users/me: %+v", updated)
	}

	// 2. GET /api/v1/users/u/gael_learner (con token propio -> IsOwner == true)
	reqGet, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/users/u/gael_learner", nil)
	reqGet.Header.Set("Authorization", "Bearer "+gaelToken)

	respGet, err := client.Do(reqGet)
	if err != nil {
		t.Fatalf("error en GET /api/v1/users/u/gael_learner: %v", err)
	}
	defer respGet.Body.Close()

	if respGet.StatusCode != http.StatusOK {
		t.Fatalf("se esperaba 200 OK en GET /api/v1/users/u/gael_learner, obtenido: %d", respGet.StatusCode)
	}

	var pubProfile models.PublicUserProfile
	_ = json.NewDecoder(respGet.Body).Decode(&pubProfile)
	if !pubProfile.IsOwner {
		t.Error("se esperaba IsOwner = true para el dueño")
	}
	if pubProfile.Email == nil {
		t.Error("se esperaba Email visible para el dueño")
	}
}

func TestGovernanceAPI_Lifecycle(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	adminToken := loginHelper(t, server.URL, "admin@palabrario.org", "password123")
	gaelToken := loginHelper(t, server.URL, "gael@palabrario.org", "password123")
	client := &http.Client{}

	// 1. Gael solicita ascenso: POST /api/v1/governance/editor-requests
	reqPayload := models.CreateEditorRequestDTO{
		InstitutionName: "Instituto de Investigaciones Filológicas UNAM",
		Justification:   "Investigador posdoctoral en lexicografía diacrónica y sincrónica.",
	}
	body, _ := json.Marshal(reqPayload)
	httpReq, _ := http.NewRequest(http.MethodPost, server.URL+"/api/v1/governance/editor-requests", bytes.NewBuffer(body))
	httpReq.Header.Set("Authorization", "Bearer "+gaelToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		t.Fatalf("error en POST /api/v1/governance/editor-requests: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("se esperaba 201 Created en editor-requests, obtenido: %d", resp.StatusCode)
	}

	var createdReq models.EditorRequest
	_ = json.NewDecoder(resp.Body).Decode(&createdReq)

	// 2. Admin lista solicitudes: GET /api/v1/governance/editor-requests
	adminGetReq, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/governance/editor-requests?status=PENDING", nil)
	adminGetReq.Header.Set("Authorization", "Bearer "+adminToken)

	respList, err := client.Do(adminGetReq)
	if err != nil {
		t.Fatalf("error en GET /api/v1/governance/editor-requests: %v", err)
	}
	defer respList.Body.Close()

	if respList.StatusCode != http.StatusOK {
		t.Fatalf("se esperaba 200 OK para admin listando solicitudes, obtenido: %d", respList.StatusCode)
	}

	var reqList []models.EditorRequest
	_ = json.NewDecoder(respList.Body).Decode(&reqList)
	if len(reqList) == 0 {
		t.Fatal("se esperaba al menos 1 solicitud listada")
	}

	// 3. Admin crea boleto: POST /api/v1/governance/tickets
	ticketPayload := models.CreateTicketRequest{
		GrantedRole: models.RoleEditor,
		MaxUses:     2,
	}
	bodyTicket, _ := json.Marshal(ticketPayload)
	reqTicket, _ := http.NewRequest(http.MethodPost, server.URL+"/api/v1/governance/tickets", bytes.NewBuffer(bodyTicket))
	reqTicket.Header.Set("Authorization", "Bearer "+adminToken)
	reqTicket.Header.Set("Content-Type", "application/json")

	respTicket, err := client.Do(reqTicket)
	if err != nil {
		t.Fatalf("error en POST /api/v1/governance/tickets: %v", err)
	}
	defer respTicket.Body.Close()

	if respTicket.StatusCode != http.StatusCreated {
		t.Fatalf("se esperaba 201 Created para admin creando ticket, obtenido: %d", respTicket.StatusCode)
	}
}

func TestReportsAPI_Lifecycle(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	adminToken := loginHelper(t, server.URL, "admin@palabrario.org", "password123")
	gaelToken := loginHelper(t, server.URL, "gael@palabrario.org", "password123")
	client := &http.Client{}

	// 1. Gael reporta una acepción: POST /api/v1/reports
	senseID := 101
	repPayload := models.CreateReportRequest{
		TargetType: "SENSE",
		IDSense:    &senseID,
		Reason:     "INACCURATE",
		Details:    "La definición contiene un error en el uso del término en el Caribe.",
	}
	body, _ := json.Marshal(repPayload)
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/api/v1/reports", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+gaelToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error en POST /api/v1/reports: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("se esperaba 201 Created en POST /api/v1/reports, obtenido: %d", resp.StatusCode)
	}

	var rep models.ContentReport
	_ = json.NewDecoder(resp.Body).Decode(&rep)

	// 2. Admin lista reportes: GET /api/v1/reports
	reqList, _ := http.NewRequest(http.MethodGet, server.URL+"/api/v1/reports?status=PENDING", nil)
	reqList.Header.Set("Authorization", "Bearer "+adminToken)

	respList, err := client.Do(reqList)
	if err != nil {
		t.Fatalf("error en GET /api/v1/reports: %v", err)
	}
	defer respList.Body.Close()

	if respList.StatusCode != http.StatusOK {
		t.Fatalf("se esperaba 200 OK en GET /api/v1/reports, obtenido: %d", respList.StatusCode)
	}

	// 3. Admin resuelve reporte: PUT /api/v1/reports/{reportId}/resolve
	resolvePayload := models.ResolveReportRequest{
		Status: "DISMISSED",
	}
	bodyRes, _ := json.Marshal(resolvePayload)
	reqRes, _ := http.NewRequest(http.MethodPut, server.URL+"/api/v1/reports/"+rep.IDReport+"/resolve", bytes.NewBuffer(bodyRes))
	reqRes.Header.Set("Authorization", "Bearer "+adminToken)
	reqRes.Header.Set("Content-Type", "application/json")

	respRes, err := client.Do(reqRes)
	if err != nil {
		t.Fatalf("error en PUT /api/v1/reports/{id}/resolve: %v", err)
	}
	defer respRes.Body.Close()

	if respRes.StatusCode != http.StatusOK {
		t.Fatalf("se esperaba 200 OK en PUT /api/v1/reports/{id}/resolve, obtenido: %d", respRes.StatusCode)
	}
}
