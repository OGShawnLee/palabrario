package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"palabrario/core-service/internal/service"
	"palabrario/core-service/internal/testutil"
	"palabrario/pkg/auth"
	"palabrario/pkg/models"
	"palabrario/pkg/utils"
)

const (
	adminUUIDStr = "b1111111-0000-0000-0000-000000000001"
	silasUUIDStr = "b1111111-0000-0000-0000-000000000002"
	gaelUUIDStr  = "b1111111-0000-0000-0000-000000000003"
)

func TestIdentityService_ProfileManagement_CU01(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	tokenSvc := auth.NewTokenService("test-secret-key-minimum-32-chars", 24*time.Hour)
	identitySvc := service.NewIdentityService(testDB, tokenSvc)
	ctx := context.Background()

	// 1. Actualizar perfil propio de Gael
	avatar := "https://img.palabrario.org/avatars/gael.webp"
	updated, err := identitySvc.UpdateProfile(ctx, gaelUUIDStr, models.UpdateProfileRequest{
		Name:          "Gael Actualizado",
		AnchorDialect: "ES",
		AvatarURL:     utils.ToPtr(avatar),
	})
	if err != nil {
		t.Fatalf("UpdateProfile falló: %v", err)
	}
	if updated.Name != "Gael Actualizado" || updated.AnchorDialect != "ES" || updated.AvatarURL == nil || *updated.AvatarURL != avatar {
		t.Errorf("datos de perfil no coinciden: %+v", updated)
	}

	// Dialecto inválido debe retornar ErrInvalidInput
	_, err = identitySvc.UpdateProfile(ctx, gaelUUIDStr, models.UpdateProfileRequest{
		Name:          "Gael",
		AnchorDialect: "ZZ",
	})
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Errorf("se esperaba ErrInvalidInput para dialecto inválido, obtenido: %v", err)
	}

	// 2. Cambio de contraseña
	// Contraseña actual errónea -> ErrUnauthorized
	err = identitySvc.ChangePassword(ctx, gaelUUIDStr, models.ChangePasswordRequest{
		CurrentPassword: "wrongpassword",
		NewPassword:     "nuevaClave123!",
	})
	if !errors.Is(err, service.ErrUnauthorized) {
		t.Errorf("se esperaba ErrUnauthorized con contraseña errónea, obtenido: %v", err)
	}

	// Nueva contraseña muy corta -> ErrInvalidInput
	err = identitySvc.ChangePassword(ctx, gaelUUIDStr, models.ChangePasswordRequest{
		CurrentPassword: "password123",
		NewPassword:     "corta",
	})
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Errorf("se esperaba ErrInvalidInput con contraseña corta, obtenido: %v", err)
	}

	// Cambio de contraseña exitoso
	err = identitySvc.ChangePassword(ctx, gaelUUIDStr, models.ChangePasswordRequest{
		CurrentPassword: "password123",
		NewPassword:     "nuevaClave123!",
	})
	if err != nil {
		t.Fatalf("ChangePassword falló con datos válidos: %v", err)
	}

	// Verificar login con la nueva contraseña
	_, err = identitySvc.SignIn(ctx, models.SignInRequest{
		Email:    "gael@palabrario.org",
		Password: "nuevaClave123!",
	})
	if err != nil {
		t.Fatalf("SignIn con nueva contraseña falló: %v", err)
	}

	// 3. Consultar perfil público (RN01: Dueño vs Tercero)
	// Como dueño
	ownerProfile, err := identitySvc.GetPublicProfileByUsername(ctx, utils.ToPtr(gaelUUIDStr), "gael_learner")
	if err != nil {
		t.Fatalf("GetPublicProfileByUsername como dueño falló: %v", err)
	}
	if !ownerProfile.IsOwner {
		t.Error("se esperaba IsOwner = true")
	}
	if ownerProfile.Email == nil || *ownerProfile.Email != "gael@palabrario.org" {
		t.Errorf("se esperaba Email visible para el dueño, obtenido: %v", ownerProfile.Email)
	}
	if ownerProfile.PrivateDecksCount == nil || ownerProfile.DueReviewsToday == nil {
		t.Error("se esperaban conteos privados poblados para el dueño")
	}

	// Como tercero (anónimo o con otro ID)
	publicProfile, err := identitySvc.GetPublicProfileByUsername(ctx, utils.ToPtr(silasUUIDStr), "gael_learner")
	if err != nil {
		t.Fatalf("GetPublicProfileByUsername como tercero falló: %v", err)
	}
	if publicProfile.IsOwner {
		t.Error("se esperaba IsOwner = false para tercero")
	}
	if publicProfile.Email != nil {
		t.Errorf("el Email NO debe ser visible para terceros (RN01), obtenido: %v", *publicProfile.Email)
	}
	if publicProfile.PrivateDecksCount != nil || publicProfile.DueReviewsToday != nil {
		t.Error("conteos privados NO deben ser visibles para terceros (RN01)")
	}
}

func TestGovernanceService_EditorRequests_CU02_CU03(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	govSvc := service.NewGovernanceService(testDB)
	ctx := context.Background()

	// 1. Gael (LEARNER) solicita ascenso a Editor
	req, err := govSvc.CreateEditorRequest(ctx, gaelUUIDStr, models.CreateEditorRequestDTO{
		InstitutionName: "Universidad de Guadalajara",
		Justification:   "Soy profesor de lingüística con 10 años de experiencia en dialectología.",
		CredentialsURL:  utils.ToPtr("https://institucion.edu/credenciales_gael.pdf"),
	})
	if err != nil {
		t.Fatalf("CreateEditorRequest falló: %v", err)
	}
	if req.Status != models.RequestPending {
		t.Errorf("se esperaba estado PENDING, obtenido: %s", req.Status)
	}
	if req.UserName != "gael_learner" {
		t.Errorf("se esperaba user_name gael_learner, obtenido: %s", req.UserName)
	}

	// 2. RN02: Intentar crear una segunda solicitud mientras una está PENDING debe fallar con ErrConflict
	_, err = govSvc.CreateEditorRequest(ctx, gaelUUIDStr, models.CreateEditorRequestDTO{
		InstitutionName: "Otra Institución",
		Justification:   "Esta es otra justificación de más de veinte caracteres requeridos.",
	})
	if !errors.Is(err, service.ErrConflict) {
		t.Errorf("se esperaba ErrConflict por RN02 (1 solicitud PENDING a la vez), obtenido: %v", err)
	}

	// 3. Usuario que ya es EDITOR (Silas) no puede solicitar ascenso -> ErrConflict
	_, err = govSvc.CreateEditorRequest(ctx, silasUUIDStr, models.CreateEditorRequestDTO{
		InstitutionName: "Universidad de Buenos Aires",
		Justification:   "Ya soy editor pero intento solicitar de nuevo de forma redundante.",
	})
	if !errors.Is(err, service.ErrConflict) {
		t.Errorf("se esperaba ErrConflict para usuario que ya es EDITOR, obtenido: %v", err)
	}

	// 4. Admin lista solicitudes pendientes
	list, err := govSvc.ListEditorRequests(ctx, utils.ToPtr("PENDING"))
	if err != nil {
		t.Fatalf("ListEditorRequests falló: %v", err)
	}
	if len(list) == 0 {
		t.Fatal("se esperaba al menos 1 solicitud pendiente")
	}

	// 5. Admin resuelve y aprueba la solicitud -> Promoción atómica a EDITOR
	resolved, err := govSvc.ResolveEditorRequest(ctx, adminUUIDStr, req.IDRequest, models.ResolveEditorRequestDTO{
		Status:      models.RequestApproved,
		ReviewNotes: utils.ToPtr("Credenciales académicas verificadas y aprobadas."),
	})
	if err != nil {
		t.Fatalf("ResolveEditorRequest falló: %v", err)
	}
	if resolved.Status != models.RequestApproved {
		t.Errorf("se esperaba estado APPROVED, obtenido: %s", resolved.Status)
	}

	// Verificar que en accounts el rol de Gael fue promovido a EDITOR
	var newRole string
	err = testDB.Pool.QueryRow(ctx, "SELECT role::TEXT FROM accounts WHERE id_account = $1", gaelUUIDStr).Scan(&newRole)
	if err != nil {
		t.Fatalf("error consultando rol promovido: %v", err)
	}
	if newRole != "EDITOR" {
		t.Errorf("se esperaba rol EDITOR tras aprobación transaccional, obtenido: %s", newRole)
	}

	// 6. Intentar resolver de nuevo una solicitud ya resuelta -> ErrConflict
	_, err = govSvc.ResolveEditorRequest(ctx, adminUUIDStr, req.IDRequest, models.ResolveEditorRequestDTO{
		Status: models.RequestRejected,
	})
	if !errors.Is(err, service.ErrConflict) {
		t.Errorf("se esperaba ErrConflict al re-resolver solicitud, obtenido: %v", err)
	}
}

func TestGovernanceService_ProtectLastAdmin_CU03(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	govSvc := service.NewGovernanceService(testDB)
	ctx := context.Background()

	// 1. Intentar degradar o suspender al único ADMIN activo (admin@palabrario.org) debe fallar (CU03 4a)
	_, err := govSvc.UpdateUserRoleOrStatus(ctx, adminUUIDStr, adminUUIDStr, models.UpdateUserRoleOrStatusRequest{
		Role: utils.ToPtr(models.RoleLearner),
	})
	if !errors.Is(err, service.ErrConflict) {
		t.Errorf("se esperaba ErrConflict protegiendo al último ADMIN activo, obtenido: %v", err)
	}

	_, err = govSvc.UpdateUserRoleOrStatus(ctx, adminUUIDStr, adminUUIDStr, models.UpdateUserRoleOrStatusRequest{
		IsActive: utils.ToPtr(false),
	})
	if !errors.Is(err, service.ErrConflict) {
		t.Errorf("se esperaba ErrConflict protegiendo contra desactivar al último ADMIN activo, obtenido: %v", err)
	}

	// 2. Ascender a Silas a ADMIN -> Ahora hay 2 administradores activos
	_, err = govSvc.UpdateUserRoleOrStatus(ctx, adminUUIDStr, silasUUIDStr, models.UpdateUserRoleOrStatusRequest{
		Role: utils.ToPtr(models.RoleAdmin),
	})
	if err != nil {
		t.Fatalf("error ascendiendo a Silas a ADMIN: %v", err)
	}

	// 3. Ahora que hay 2 administradores, degradar al admin original sí debe permitirse
	_, err = govSvc.UpdateUserRoleOrStatus(ctx, silasUUIDStr, adminUUIDStr, models.UpdateUserRoleOrStatusRequest{
		Role: utils.ToPtr(models.RoleLearner),
	})
	if err != nil {
		t.Fatalf("degradación debió tener éxito habiendo otro admin activo: %v", err)
	}
}

func TestGovernanceService_InvitationTickets_CU04(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	govSvc := service.NewGovernanceService(testDB)
	ctx := context.Background()

	// 1. Admin crea un boleto para otorgar rol EDITOR con 1 solo uso
	code := "TICKET-UNAM-2026-TEST"
	ticket, err := govSvc.CreateInvitationTicket(ctx, adminUUIDStr, models.CreateTicketRequest{
		Code:            utils.ToPtr(code),
		GrantedRole:     models.RoleEditor,
		InstitutionName: utils.ToPtr("Universidad Nacional Autónoma de México"),
		MaxUses:         1,
	})
	if err != nil {
		t.Fatalf("CreateInvitationTicket falló: %v", err)
	}
	if ticket.Code != code || ticket.GrantedRole != models.RoleEditor || ticket.MaxUses != 1 {
		t.Errorf("boleto creado con datos incorrectos: %+v", ticket)
	}

	// 2. Aprendiz Gael redime el boleto
	redeemedProfile, err := govSvc.RedeemInvitationTicket(ctx, gaelUUIDStr, code)
	if err != nil {
		t.Fatalf("RedeemInvitationTicket falló: %v", err)
	}
	if redeemedProfile.Role != models.RoleEditor {
		t.Errorf("se esperaba que Gael fuera promovido a EDITOR tras redimir boleto, obtenido: %s", redeemedProfile.Role)
	}

	// 3. Gael intenta redimir el mismo boleto de nuevo -> ErrConflict (ya redimido por el usuario)
	_, err = govSvc.RedeemInvitationTicket(ctx, gaelUUIDStr, code)
	if !errors.Is(err, service.ErrConflict) {
		t.Errorf("se esperaba ErrConflict por redención duplicada, obtenido: %v", err)
	}

	// 4. Crear otro aprendiz que intente usar el boleto agotado (max_uses: 1) -> ErrConflict
	tokenSvc := auth.NewTokenService("test-secret-key-minimum-32-chars", 24*time.Hour)
	identitySvc := service.NewIdentityService(testDB, tokenSvc)
	user2, err := identitySvc.SignUp(ctx, models.SignUpRequest{
		Email:         "usuario2@palabrario.org",
		Password:      "password123!",
		Name:          "Usuario Dos",
		UserName:      "usuario_dos",
		AnchorDialect: "AR",
	})
	if err != nil {
		t.Fatalf("error creando usuario 2: %v", err)
	}

	_, err = govSvc.RedeemInvitationTicket(ctx, user2.User.IDUser, code)
	if !errors.Is(err, service.ErrConflict) {
		t.Errorf("se esperaba ErrConflict por boleto agotado en max_uses, obtenido: %v", err)
	}

	// 5. Boleto revocado
	codeRevoked := "TICKET-REVOKED-TEST"
	tRev, err := govSvc.CreateInvitationTicket(ctx, adminUUIDStr, models.CreateTicketRequest{
		Code:        utils.ToPtr(codeRevoked),
		GrantedRole: models.RoleEditor,
		MaxUses:     5,
	})
	if err != nil {
		t.Fatalf("error creando boleto revocable: %v", err)
	}

	err = govSvc.RevokeInvitationTicket(ctx, adminUUIDStr, tRev.IDTicket)
	if err != nil {
		t.Fatalf("RevokeInvitationTicket falló: %v", err)
	}

	_, err = govSvc.RedeemInvitationTicket(ctx, user2.User.IDUser, codeRevoked)
	if !errors.Is(err, service.ErrConflict) {
		t.Errorf("se esperaba ErrConflict por boleto inactivo/revocado, obtenido: %v", err)
	}
}

func TestCurationService_ContentReports_CU09(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	if testDB == nil {
		return
	}

	curationSvc := service.NewCurationService(testDB)
	ctx := context.Background()

	// 1. Reportar acepción 101 ('popote') por spam/vandalismo
	senseID := 101
	rep, err := curationSvc.CreateReport(ctx, gaelUUIDStr, models.CreateReportRequest{
		TargetType: "SENSE",
		IDSense:    utils.ToPtr(senseID),
		Reason:     "SPAM_OR_VANDALISM",
		Details:    "Este texto contiene información no apropiada para la comunidad panhispánica.",
	})
	if err != nil {
		t.Fatalf("CreateReport falló: %v", err)
	}
	if rep.Status != "PENDING" || rep.TargetType != "SENSE" || rep.TargetPreview == "" {
		t.Errorf("reporte creado incorrecto: %+v", rep)
	}

	// Validaciones: detalle muy corto (< 10 caracteres) -> ErrInvalidInput
	_, err = curationSvc.CreateReport(ctx, gaelUUIDStr, models.CreateReportRequest{
		TargetType: "SENSE",
		IDSense:    utils.ToPtr(senseID),
		Reason:     "SPAM_OR_VANDALISM",
		Details:    "corto",
	})
	if !errors.Is(err, service.ErrInvalidInput) {
		t.Errorf("se esperaba ErrInvalidInput para detalles menores a 10 caracteres, obtenido: %v", err)
	}

	// 2. Admin lista reportes pendientes
	list, err := curationSvc.ListOpenReports(ctx, utils.ToPtr("PENDING"))
	if err != nil {
		t.Fatalf("ListOpenReports falló: %v", err)
	}
	if len(list) == 0 {
		t.Fatal("se esperaba al menos 1 reporte pendiente en la lista")
	}

	// 3. Admin resuelve reporte como RESOLVED_HIDDEN -> La acepción 101 pasa a estado 'HIDDEN'
	resolved, err := curationSvc.ResolveReport(ctx, adminUUIDStr, rep.IDReport, models.ResolveReportRequest{
		Status:          "RESOLVED_HIDDEN",
		ResolutionNotes: utils.ToPtr("Aceptado el reporte; se oculta la acepción para revisión lexicográfica."),
	})
	if err != nil {
		t.Fatalf("ResolveReport falló: %v", err)
	}
	if resolved.Status != "RESOLVED_HIDDEN" {
		t.Errorf("se esperaba estado RESOLVED_HIDDEN, obtenido: %s", resolved.Status)
	}

	// Verificar que el sentido 101 fue ocultado en la base de datos
	var senseState string
	err = testDB.Pool.QueryRow(ctx, "SELECT state::TEXT FROM senses WHERE id_sense = $1", senseID).Scan(&senseState)
	if err != nil {
		t.Fatalf("error consultando estado de acepción reportada: %v", err)
	}
	if senseState != "HIDDEN" {
		t.Errorf("se esperaba que el sentido 101 estuviera en estado 'HIDDEN', obtenido: %s", senseState)
	}

	// 4. Intentar resolver de nuevo -> ErrConflict
	_, err = curationSvc.ResolveReport(ctx, adminUUIDStr, rep.IDReport, models.ResolveReportRequest{
		Status: "DISMISSED",
	})
	if !errors.Is(err, service.ErrConflict) {
		t.Errorf("se esperaba ErrConflict al intentar re-resolver un reporte, obtenido: %v", err)
	}
}
