package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"palabrario/core-service/internal/database"
	"palabrario/core-service/internal/db"
	"palabrario/pkg/models"
)

// CurationService define el contrato SOA del Módulo 4 (Curaduría por Quórum y Moderación — CU07, CU08, CU09)
// según la sección 2.4 de .context/INTERFACES.md.
type CurationService interface {
	ListPendingQueue(ctx context.Context, regionFilter *string) ([]models.SenseQuorumItem, error)
	ListMyProposals(ctx context.Context, editorID string, stateFilter *string) ([]models.SenseQuorumItem, error)
	GetProposalDetail(ctx context.Context, senseID int) (*models.SenseQuorumDetail, error)
	SubmitVote(ctx context.Context, editorID string, senseID int, req models.SubmitVoteRequest) (*models.SenseQuorumDetail, error)
	AdminOverrideState(ctx context.Context, adminID string, senseID int, req models.AdminOverrideStateRequest) (*models.SenseDetail, error)

	CreateReport(ctx context.Context, reporterID string, req models.CreateReportRequest) (*models.ContentReport, error)
	ListOpenReports(ctx context.Context, statusFilter *string) ([]models.ContentReport, error)
	ResolveReport(ctx context.Context, adminID string, reportID string, req models.ResolveReportRequest) (*models.ContentReport, error)
}

type curationService struct {
	databaseInstance *database.DB
	queries          *db.Queries
}

func NewCurationService(databaseInstance *database.DB) CurationService {
	return &curationService{
		databaseInstance: databaseInstance,
		queries:          db.New(databaseInstance.Pool),
	}
}

func (s *curationService) ListPendingQueue(ctx context.Context, regionFilter *string) ([]models.SenseQuorumItem, error) {
	var regText pgtype.Text
	if regionFilter != nil {
		regText = StringToText(strings.ToUpper(*regionFilter))
	}

	rows, err := s.queries.ListPendingQuorumQueue(ctx, regText)
	if err != nil {
		return nil, fmt.Errorf("error listando cola de curaduría: %w", err)
	}

	results := make([]models.SenseQuorumItem, len(rows))
	for i, r := range rows {
		results[i] = models.SenseQuorumItem{
			IDSense:         int(r.IDSense),
			IDWord:          int(r.IDWord),
			Lemma:           r.Lemma,
			PartOfSpeech:    r.PartOfSpeech,
			RegionCode:      r.RegionCode,
			RegionName:      r.RegionName,
			FlagEmoji:       r.FlagEmoji,
			ShortDefinition: r.ShortDefinition,
			Definition:      r.Definition,
			CulturalNote:    TextToPtr(r.CulturalNote),
			CEFRLevel:       InterfaceToString(r.CefrLevel, "B1"),
			UsageFrequency:  InterfaceToString(r.UsageFrequency, "COMMON"),
			IsTabooOrVulgar: r.IsTabooOrVulgar,
			State:           r.State,
			AuthorID:        r.AuthorID.String(),
			AuthorUserName:  r.AuthorUserName,
			AuthorFullName:  r.AuthorFullName,
			ApprovalsCount:  int(r.ApprovalsCount),
			RequiredQuorum:  int(r.RequiredQuorum),
			CreatedAt:       r.CreatedAt.Time,
		}
	}

	return results, nil
}

func (s *curationService) ListMyProposals(ctx context.Context, editorID string, stateFilter *string) ([]models.SenseQuorumItem, error) {
	editorUUID, err := uuid.Parse(strings.TrimSpace(editorID))
	if err != nil {
		return nil, NewDomainError(ErrUnauthorized, "Identificador de editor no válido")
	}

	var stText pgtype.Text
	if stateFilter != nil {
		stText = StringToText(strings.ToUpper(*stateFilter))
	}

	rows, err := s.queries.ListEditorProposals(ctx, db.ListEditorProposalsParams{
		IDEditor: editorUUID,
		State:    stText,
	})
	if err != nil {
		return nil, fmt.Errorf("error listando propuestas del editor: %w", err)
	}

	results := make([]models.SenseQuorumItem, len(rows))
	for i, r := range rows {
		results[i] = models.SenseQuorumItem{
			IDSense:         int(r.IDSense),
			IDWord:          int(r.IDWord),
			Lemma:           r.Lemma,
			PartOfSpeech:    r.PartOfSpeech,
			RegionCode:      r.RegionCode,
			RegionName:      r.RegionName,
			FlagEmoji:       r.FlagEmoji,
			ShortDefinition: r.ShortDefinition,
			Definition:      r.Definition,
			CulturalNote:    TextToPtr(r.CulturalNote),
			CEFRLevel:       InterfaceToString(r.CefrLevel, "B1"),
			UsageFrequency:  InterfaceToString(r.UsageFrequency, "COMMON"),
			IsTabooOrVulgar: r.IsTabooOrVulgar,
			State:           r.State,
			AuthorID:        r.AuthorID.String(),
			AuthorUserName:  r.AuthorUserName,
			AuthorFullName:  r.AuthorFullName,
			ApprovalsCount:  int(r.ApprovalsCount),
			RequiredQuorum:  int(r.RequiredQuorum),
			CreatedAt:       r.CreatedAt.Time,
		}
	}

	return results, nil
}

func (s *curationService) GetProposalDetail(ctx context.Context, senseID int) (*models.SenseQuorumDetail, error) {
	if senseID <= 0 {
		return nil, NewDomainError(ErrInvalidInput, "El identificador de acepción debe ser positivo")
	}

	itemRow, err := s.queries.GetSenseQuorumItemByID(ctx, int32(senseID))
	if err != nil {
		return nil, MapDBError(err, "Propuesta no encontrada en la cola de curaduría")
	}

	examplesRows, err := s.queries.GetExamplesBySenseID(ctx, int32(senseID))
	if err != nil {
		return nil, fmt.Errorf("error obteniendo ejemplos de la propuesta: %w", err)
	}

	examples := make([]models.Example, len(examplesRows))
	for i, ex := range examplesRows {
		examples[i] = models.Example{
			IDExample: int(ex.IDExample),
			IDSense:   int(ex.IDSense),
			Sentence:  ex.Sentence,
			Source:    TextToPtr(ex.Source),
		}
	}

	votesRows, err := s.queries.ListVotesBySenseID(ctx, int32(senseID))
	if err != nil {
		return nil, fmt.Errorf("error obteniendo votos de la propuesta: %w", err)
	}

	votes := make([]models.SenseApprovalVote, len(votesRows))
	for i, v := range votesRows {
		votes[i] = models.SenseApprovalVote{
			IDSense:        int(v.IDSense),
			IDEditor:       v.IDEditor.String(),
			EditorUserName: v.EditorUserName,
			EditorFullName: v.EditorFullName,
			EditorDialect:  v.EditorDialect,
			Decision:       models.ReviewDecision(v.Decision),
			FeedbackNotes:  TextToPtr(v.FeedbackNotes),
			ApprovedAt:     v.ApprovedAt.Time,
		}
	}

	proposalItem := models.SenseQuorumItem{
		IDSense:         int(itemRow.IDSense),
		IDWord:          int(itemRow.IDWord),
		Lemma:           itemRow.Lemma,
		PartOfSpeech:    itemRow.PartOfSpeech,
		RegionCode:      itemRow.RegionCode,
		RegionName:      itemRow.RegionName,
		FlagEmoji:       itemRow.FlagEmoji,
		ShortDefinition: itemRow.ShortDefinition,
		Definition:      itemRow.Definition,
		CulturalNote:    TextToPtr(itemRow.CulturalNote),
		CEFRLevel:       InterfaceToString(itemRow.CefrLevel, "B1"),
		UsageFrequency:  InterfaceToString(itemRow.UsageFrequency, "COMMON"),
		IsTabooOrVulgar: itemRow.IsTabooOrVulgar,
		State:           itemRow.State,
		AuthorID:        itemRow.AuthorID.String(),
		AuthorUserName:  itemRow.AuthorUserName,
		AuthorFullName:  itemRow.AuthorFullName,
		ApprovalsCount:  int(itemRow.ApprovalsCount),
		RequiredQuorum:  int(itemRow.RequiredQuorum),
		CreatedAt:       itemRow.CreatedAt.Time,
	}

	return &models.SenseQuorumDetail{
		Proposal: proposalItem,
		Examples: examples,
		Votes:    votes,
	}, nil
}

func (s *curationService) SubmitVote(ctx context.Context, editorID string, senseID int, req models.SubmitVoteRequest) (*models.SenseQuorumDetail, error) {
	if req.Decision != models.DecisionApprove && req.Decision != models.DecisionRequestChanges && req.Decision != models.DecisionReject {
		return nil, NewDomainError(ErrInvalidInput, "Decisión de revisión no válida (debe ser APPROVE, REQUEST_CHANGES o REJECT)")
	}

	if req.Decision != models.DecisionApprove && (req.FeedbackNotes == nil || strings.TrimSpace(*req.FeedbackNotes) == "") {
		return nil, NewDomainError(ErrInvalidInput, "Debes incluir notas de retroalimentación al solicitar correcciones o rechazar una propuesta")
	}

	editorUUID, err := uuid.Parse(strings.TrimSpace(editorID))
	if err != nil {
		return nil, NewDomainError(ErrUnauthorized, "Identificador de editor no válido")
	}

	// 1. Obtener la acepción para validar existencia y regla RN07
	senseRow, err := s.queries.GetSenseByID(ctx, int32(senseID))
	if err != nil {
		return nil, MapDBError(err, "Sentido no encontrado")
	}

	// RN07: Prohibición de auto-aprobación en código Go (también protegida por trigger DB)
	if senseRow.IDEditor == editorUUID {
		return nil, NewDomainError(ErrForbidden, "Violación de integridad (RN07): Un editor no puede votar ni aprobar su propia contribución")
	}

	// 2. Insertar o actualizar voto en sense_approvals
	// El trigger 'trg_auto_approve_sense' en PostgreSQL se disparará automáticamente tras el INSERT/UPDATE:
	// si se alcanzan 3 votos con decision = 'APPROVE', actualizará senses.state = 'APPROVED'
	_, err = s.queries.UpsertSenseApproval(ctx, db.UpsertSenseApprovalParams{
		IDSense:       int32(senseID),
		IDEditor:      editorUUID,
		Decision:      db.ReviewDecision(req.Decision),
		FeedbackNotes: PtrToText(req.FeedbackNotes),
	})
	if err != nil {
		return nil, MapDBError(err, "Error al registrar voto de curaduría")
	}

	// 3. Devolver la propuesta actualizada con votos y estado (si cambió a APPROVED por el trigger RN06)
	return s.GetProposalDetail(ctx, senseID)
}

func (s *curationService) AdminOverrideState(ctx context.Context, adminID string, senseID int, req models.AdminOverrideStateRequest) (*models.SenseDetail, error) {
	st := strings.ToUpper(strings.TrimSpace(req.State))
	if st != "APPROVED" && st != "HIDDEN" {
		return nil, NewDomainError(ErrInvalidInput, "El estado debe ser 'APPROVED' o 'HIDDEN'")
	}

	_, err := s.queries.AdminUpdateSenseState(ctx, db.AdminUpdateSenseStateParams{
		IDSense: int32(senseID),
		State:   db.ReviewState(st),
	})
	if err != nil {
		return nil, MapDBError(err, "Error al actualizar estado del sentido")
	}

	sr, err := s.queries.GetSenseByID(ctx, int32(senseID))
	if err != nil {
		return nil, MapDBError(err, "Sentido no encontrado")
	}

	examplesRows, err := s.queries.GetExamplesBySenseID(ctx, int32(senseID))
	if err != nil {
		return nil, fmt.Errorf("error consultando ejemplos del sentido: %w", err)
	}

	examples := make([]models.Example, len(examplesRows))
	for i, ex := range examplesRows {
		examples[i] = models.Example{
			IDExample: int(ex.IDExample),
			IDSense:   int(ex.IDSense),
			Sentence:  ex.Sentence,
			Source:    TextToPtr(ex.Source),
		}
	}

	countryCode := ""
	if sr.CountryCode.Valid {
		countryCode = sr.CountryCode.String
	}

	cefrStr := ""
	if sr.CefrLevel.Valid {
		cefrStr = string(sr.CefrLevel.CefrLevel)
	}

	freqStr := ""
	if sr.UsageFrequency.Valid {
		freqStr = string(sr.UsageFrequency.UsageFrequency)
	}

	return &models.SenseDetail{
		IDSense:         int(sr.IDSense),
		IDWord:          int(sr.IDWord),
		RegionCode:      sr.RegionCode,
		RegionName:      sr.RegionName,
		CountryCode:     countryCode,
		FlagEmoji:       sr.FlagEmoji,
		ShortDefinition: sr.ShortDefinition,
		Definition:      sr.Definition,
		IsPrimary:       sr.IsPrimary,
		IsTabooOrVulgar: sr.IsTabooOrVulgar,
		CulturalNote:    TextToPtr(sr.CulturalNote),
		CEFRLevel:       cefrStr,
		UsageFrequency:  freqStr,
		State:           string(sr.State),
		Examples:        examples,
	}, nil
}

// ------------------------------------------------------------------------------
// CU09: Moderación y Resolución de Reportes de Contenido
// ------------------------------------------------------------------------------

func (s *curationService) CreateReport(ctx context.Context, reporterID string, req models.CreateReportRequest) (*models.ContentReport, error) {
	reporterUUID, err := ParseDomainUUID(reporterID, "Identificador de usuario inválido")
	if err != nil {
		return nil, err
	}

	tt := strings.ToUpper(strings.TrimSpace(req.TargetType))
	if tt != "WORD" && tt != "SENSE" && tt != "EXAMPLE" {
		return nil, NewDomainError(ErrInvalidInput, "Tipo de objetivo inválido (debe ser 'WORD', 'SENSE' o 'EXAMPLE')")
	}

	var wID, sID, eID pgtype.Int4
	switch tt {
	case "WORD":
		if req.IDWord == nil || *req.IDWord <= 0 {
			return nil, NewDomainError(ErrInvalidInput, "id_word es obligatorio para reportes de palabra")
		}
		wID = pgtype.Int4{Int32: int32(*req.IDWord), Valid: true}
	case "SENSE":
		if req.IDSense == nil || *req.IDSense <= 0 {
			return nil, NewDomainError(ErrInvalidInput, "id_sense es obligatorio para reportes de acepción")
		}
		sID = pgtype.Int4{Int32: int32(*req.IDSense), Valid: true}
	case "EXAMPLE":
		if req.IDExample == nil || *req.IDExample <= 0 {
			return nil, NewDomainError(ErrInvalidInput, "id_example es obligatorio para reportes de ejemplo")
		}
		eID = pgtype.Int4{Int32: int32(*req.IDExample), Valid: true}
	}

	reason := strings.ToUpper(strings.TrimSpace(req.Reason))
	validReasons := map[string]bool{
		"INACCURATE":         true,
		"OFFENSIVE_UNMARKED": true,
		"SPAM_OR_VANDALISM":  true,
		"DUPLICATE":          true,
		"OTHER":              true,
	}
	if !validReasons[reason] {
		return nil, NewDomainError(ErrInvalidInput, "Motivo de reporte no válido")
	}

	details := strings.TrimSpace(req.Details)
	if len(details) < 10 {
		return nil, NewDomainError(ErrInvalidInput, "El detalle del reporte debe tener al menos 10 caracteres")
	}

	created, err := s.queries.CreateContentReport(ctx, db.CreateContentReportParams{
		IDReporter: reporterUUID,
		TargetType: db.ReportTargetType(tt),
		IDWord:     wID,
		IDSense:    sID,
		IDExample:  eID,
		Reason:     db.ReportReason(reason),
		Details:    details,
	})
	if err != nil {
		return nil, MapDBError(err, "Error al crear reporte de contenido")
	}

	return s.getReportDetail(ctx, created.IDReport)
}

func (s *curationService) ListOpenReports(ctx context.Context, statusFilter *string) ([]models.ContentReport, error) {
	var st pgtype.Text
	if statusFilter != nil {
		st = StringToText(strings.ToUpper(*statusFilter))
	}

	rows, err := s.queries.ListContentReports(ctx, st)
	if err != nil {
		return nil, fmt.Errorf("error listando reportes de contenido: %w", err)
	}

	results := make([]models.ContentReport, len(rows))
	for i, r := range rows {
		results[i] = mapContentReportRow(r.IDReport, r.IDReporter, r.ReporterUserName, r.TargetType, r.IDWord, r.IDSense, r.IDExample, r.TargetPreview, r.Reason, r.Details, r.Status, r.ResolvedBy, r.ResolutionNotes, r.CreatedAt, r.ResolvedAt)
	}
	return results, nil
}

func (s *curationService) ResolveReport(ctx context.Context, adminID string, reportID string, req models.ResolveReportRequest) (*models.ContentReport, error) {
	adminUUID, err := ParseDomainUUID(adminID, "Identificador de administrador inválido")
	if err != nil {
		return nil, err
	}
	reportUUID, err := ParseDomainUUID(reportID, "Identificador de reporte inválido")
	if err != nil {
		return nil, err
	}

	st := strings.ToUpper(strings.TrimSpace(req.Status))
	if st != "RESOLVED_HIDDEN" && st != "RESOLVED_AMENDED" && st != "DISMISSED" {
		return nil, NewDomainError(ErrInvalidInput, "Estado de resolución inválido (debe ser 'RESOLVED_HIDDEN', 'RESOLVED_AMENDED' o 'DISMISSED')")
	}

	tx, err := s.databaseInstance.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error iniciando transacción: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.queries.WithTx(tx)

	currentRep, err := qtx.GetContentReportByID(ctx, reportUUID)
	if err != nil {
		return nil, MapDBError(err, "Reporte de contenido no encontrado")
	}
	if currentRep.Status != "PENDING" {
		return nil, NewDomainError(ErrConflict, "El reporte ya ha sido resuelto previamente")
	}

	_, err = qtx.ResolveContentReport(ctx, db.ResolveContentReportParams{
		IDReport:        reportUUID,
		Status:          db.ReportStatus(st),
		ResolvedBy:      UUIDToPgtype(adminUUID),
		ResolutionNotes: PtrToText(req.ResolutionNotes),
	})
	if err != nil {
		return nil, MapDBError(err, "Error resolviendo reporte de contenido")
	}

	// Si se resuelve como RESOLVED_HIDDEN y el target era un SENSE, ocultar el sentido
	if st == "RESOLVED_HIDDEN" && currentRep.TargetType == "SENSE" && currentRep.IDSense.Valid {
		_, err = qtx.AdminUpdateSenseState(ctx, db.AdminUpdateSenseStateParams{
			IDSense: currentRep.IDSense.Int32,
			State:   db.ReviewStateHIDDEN,
		})
		if err != nil {
			return nil, MapDBError(err, "Error ocultando la acepción reportada")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error confirmando transacción: %w", err)
	}

	return s.getReportDetail(ctx, reportUUID)
}

func (s *curationService) getReportDetail(ctx context.Context, id uuid.UUID) (*models.ContentReport, error) {
	r, err := s.queries.GetContentReportByID(ctx, id)
	if err != nil {
		return nil, MapDBError(err, "Reporte no encontrado")
	}
	res := mapContentReportRow(r.IDReport, r.IDReporter, r.ReporterUserName, r.TargetType, r.IDWord, r.IDSense, r.IDExample, r.TargetPreview, r.Reason, r.Details, r.Status, r.ResolvedBy, r.ResolutionNotes, r.CreatedAt, r.ResolvedAt)
	return &res, nil
}

func mapContentReportRow(idReport, idReporter uuid.UUID, reporterName, targetType string, idWord, idSense, idExample pgtype.Int4, targetPreview interface{}, reason, details, status string, resolvedBy string, resolutionNotes pgtype.Text, createdAt, resolvedAt pgtype.Timestamptz) models.ContentReport {
	var wID, sID, eID *int
	if idWord.Valid {
		v := int(idWord.Int32)
		wID = &v
	}
	if idSense.Valid {
		v := int(idSense.Int32)
		sID = &v
	}
	if idExample.Valid {
		v := int(idExample.Int32)
		eID = &v
	}

	preview := InterfaceToString(targetPreview, "")

	var resByPtr *string
	if resolvedBy != "" {
		resByPtr = &resolvedBy
	}

	var resAtPtr *time.Time
	if resolvedAt.Valid {
		t := resolvedAt.Time
		resAtPtr = &t
	}

	return models.ContentReport{
		IDReport:        idReport.String(),
		IDReporter:      idReporter.String(),
		ReporterName:    reporterName,
		TargetType:      targetType,
		IDWord:          wID,
		IDSense:         sID,
		IDExample:       eID,
		TargetPreview:   preview,
		Reason:          reason,
		Details:         details,
		Status:          status,
		ResolvedBy:      resByPtr,
		ResolutionNotes: TextToPtr(resolutionNotes),
		CreatedAt:       createdAt.Time,
		ResolvedAt:      resAtPtr,
	}
}
