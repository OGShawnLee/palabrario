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
	client := newTestHTTPClient()
	payload := models.SignInRequest{
		Email:    email,
		Password: password,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal sign-in payload: %v", err)
	}
	resp, err := client.Post(serverURL+"/api/v1/auth/sign-in", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("sign-in failed for %s: %v", email, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("sign-in failed with status %d for %s", resp.StatusCode, email)
	}

	var authResp models.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatalf("failed to decode AuthResponse for %s: %v", email, err)
	}
	return authResp.Token
}

func TestUsersAPI_Lifecycle(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	gaelToken := loginHelper(t, server.URL, "gael@palabrario.org", "password123")
	client := newTestHTTPClient()

	// 1. PUT /api/v1/users/me
	updatePayload := models.UpdateProfileRequest{
		Name:          "Gael Modificado HTTP",
		AnchorDialect: "AR",
	}
	body, err := json.Marshal(updatePayload)
	if err != nil {
		t.Fatalf("failed to marshal update profile payload: %v", err)
	}
	req, err := http.NewRequest(http.MethodPut, server.URL+"/api/v1/users/me", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to build PUT /api/v1/users/me request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+gaelToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("PUT /api/v1/users/me failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on PUT /api/v1/users/me, got: %d", resp.StatusCode)
	}

	var updated models.UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode UserProfile: %v", err)
	}
	if updated.Name != "Gael Modificado HTTP" || updated.AnchorDialect != "AR" {
		t.Errorf("profile mismatch on PUT /api/v1/users/me: %+v", updated)
	}

	// 2. GET /api/v1/users/u/gael_learner (with owner token -> IsOwner == true)
	reqGet, err := http.NewRequest(http.MethodGet, server.URL+"/api/v1/users/u/gael_learner", nil)
	if err != nil {
		t.Fatalf("failed to build GET /api/v1/users/u/gael_learner request: %v", err)
	}
	reqGet.Header.Set("Authorization", "Bearer "+gaelToken)

	respGet, err := client.Do(reqGet)
	if err != nil {
		t.Fatalf("GET /api/v1/users/u/gael_learner failed: %v", err)
	}
	defer respGet.Body.Close()

	if respGet.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on GET /api/v1/users/u/gael_learner, got: %d", respGet.StatusCode)
	}

	var pubProfile models.PublicUserProfile
	if err := json.NewDecoder(respGet.Body).Decode(&pubProfile); err != nil {
		t.Fatalf("failed to decode PublicUserProfile: %v", err)
	}
	if !pubProfile.IsOwner {
		t.Error("expected IsOwner = true for profile owner")
	}
	if pubProfile.Email == nil {
		t.Error("expected Email visible for profile owner")
	}
}

func TestGovernanceAPI_Lifecycle(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	adminToken := loginHelper(t, server.URL, "admin@palabrario.org", "password123")
	gaelToken := loginHelper(t, server.URL, "gael@palabrario.org", "password123")
	client := newTestHTTPClient()

	// 1. Gael requests promotion: POST /api/v1/governance/editor-requests
	reqPayload := models.CreateEditorRequestDTO{
		InstitutionName: "Instituto de Investigaciones Filológicas UNAM",
		Justification:   "Investigador posdoctoral en lexicografía diacrónica y sincrónica.",
	}
	body, err := json.Marshal(reqPayload)
	if err != nil {
		t.Fatalf("failed to marshal editor request payload: %v", err)
	}
	httpReq, err := http.NewRequest(http.MethodPost, server.URL+"/api/v1/governance/editor-requests", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to build POST /api/v1/governance/editor-requests: %v", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+gaelToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		t.Fatalf("POST /api/v1/governance/editor-requests failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created on editor-requests, got: %d", resp.StatusCode)
	}

	var createdReq models.EditorRequest
	if err := json.NewDecoder(resp.Body).Decode(&createdReq); err != nil {
		t.Fatalf("failed to decode EditorRequest: %v", err)
	}
	if createdReq.IDRequest == "" {
		t.Error("expected non-empty IDRequest")
	}

	// 2. Admin lists requests: GET /api/v1/governance/editor-requests
	adminGetReq, err := http.NewRequest(http.MethodGet, server.URL+"/api/v1/governance/editor-requests?status=PENDING", nil)
	if err != nil {
		t.Fatalf("failed to build GET /api/v1/governance/editor-requests: %v", err)
	}
	adminGetReq.Header.Set("Authorization", "Bearer "+adminToken)

	respList, err := client.Do(adminGetReq)
	if err != nil {
		t.Fatalf("GET /api/v1/governance/editor-requests failed: %v", err)
	}
	defer respList.Body.Close()

	if respList.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK listing editor requests, got: %d", respList.StatusCode)
	}

	var reqList []models.EditorRequest
	if err := json.NewDecoder(respList.Body).Decode(&reqList); err != nil {
		t.Fatalf("failed to decode EditorRequest slice: %v", err)
	}
	if len(reqList) == 0 {
		t.Fatal("expected at least 1 editor request")
	}

	// 3. Admin creates ticket: POST /api/v1/governance/tickets
	ticketPayload := models.CreateTicketRequest{
		GrantedRole: models.RoleEditor,
		MaxUses:     2,
	}
	bodyTicket, err := json.Marshal(ticketPayload)
	if err != nil {
		t.Fatalf("failed to marshal CreateTicketRequest: %v", err)
	}
	reqTicket, err := http.NewRequest(http.MethodPost, server.URL+"/api/v1/governance/tickets", bytes.NewBuffer(bodyTicket))
	if err != nil {
		t.Fatalf("failed to build POST /api/v1/governance/tickets: %v", err)
	}
	reqTicket.Header.Set("Authorization", "Bearer "+adminToken)
	reqTicket.Header.Set("Content-Type", "application/json")

	respTicket, err := client.Do(reqTicket)
	if err != nil {
		t.Fatalf("POST /api/v1/governance/tickets failed: %v", err)
	}
	defer respTicket.Body.Close()

	if respTicket.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created creating ticket, got: %d", respTicket.StatusCode)
	}
}

func TestReportsAPI_Lifecycle(t *testing.T) {
	server, _ := setupTestServer(t)
	if server == nil {
		return
	}

	adminToken := loginHelper(t, server.URL, "admin@palabrario.org", "password123")
	gaelToken := loginHelper(t, server.URL, "gael@palabrario.org", "password123")
	client := newTestHTTPClient()

	// 1. Gael creates a report: POST /api/v1/reports
	senseID := 101
	repPayload := models.CreateReportRequest{
		TargetType: "SENSE",
		IDSense:    &senseID,
		Reason:     "INACCURATE",
		Details:    "La definición contiene un error en el uso del término en el Caribe.",
	}
	body, err := json.Marshal(repPayload)
	if err != nil {
		t.Fatalf("failed to marshal CreateReportRequest: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/v1/reports", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to build POST /api/v1/reports: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+gaelToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("POST /api/v1/reports failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created on POST /api/v1/reports, got: %d", resp.StatusCode)
	}

	var rep models.ContentReport
	if err := json.NewDecoder(resp.Body).Decode(&rep); err != nil {
		t.Fatalf("failed to decode ContentReport: %v", err)
	}

	// 2. Admin lists reports: GET /api/v1/reports
	reqList, err := http.NewRequest(http.MethodGet, server.URL+"/api/v1/reports?status=PENDING", nil)
	if err != nil {
		t.Fatalf("failed to build GET /api/v1/reports: %v", err)
	}
	reqList.Header.Set("Authorization", "Bearer "+adminToken)

	respList, err := client.Do(reqList)
	if err != nil {
		t.Fatalf("GET /api/v1/reports failed: %v", err)
	}
	defer respList.Body.Close()

	if respList.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on GET /api/v1/reports, got: %d", respList.StatusCode)
	}

	// 3. Admin resolves report: PUT /api/v1/reports/{reportId}/resolve
	resolvePayload := models.ResolveReportRequest{
		Status: "DISMISSED",
	}
	bodyRes, err := json.Marshal(resolvePayload)
	if err != nil {
		t.Fatalf("failed to marshal ResolveReportRequest: %v", err)
	}
	reqRes, err := http.NewRequest(http.MethodPut, server.URL+"/api/v1/reports/"+rep.IDReport+"/resolve", bytes.NewBuffer(bodyRes))
	if err != nil {
		t.Fatalf("failed to build PUT /api/v1/reports/{id}/resolve: %v", err)
	}
	reqRes.Header.Set("Authorization", "Bearer "+adminToken)
	reqRes.Header.Set("Content-Type", "application/json")

	respRes, err := client.Do(reqRes)
	if err != nil {
		t.Fatalf("PUT /api/v1/reports/{id}/resolve failed: %v", err)
	}
	defer respRes.Body.Close()

	if respRes.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on PUT /api/v1/reports/{id}/resolve, got: %d", respRes.StatusCode)
	}
}
