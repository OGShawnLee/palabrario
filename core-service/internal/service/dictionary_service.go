package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"palabrario/core-service/internal/database"
	"palabrario/core-service/internal/db"
	"palabrario/pkg/models"
)

// DictionaryService define el contrato SOA del Módulo 3 (Léxico y Búsqueda Dialectal — CU05, CU06)
// según la sección 2.3 de .context/INTERFACES.md.
type DictionaryService interface {
	SearchDropdown(ctx context.Context, query string, userDialect string, limit int) ([]models.WordSearchResult, error)
	GetWordDetail(ctx context.Context, wordID int) (*models.WordDetail, error)
	ListRegions(ctx context.Context) ([]models.Region, error)
	GetWordOfTheDay(ctx context.Context, dialect string) (*models.WordDetail, error)
	ExploreWords(ctx context.Context, filter models.WordExploreFilter) (*models.PaginatedWords, error)
	CreateWordProposal(ctx context.Context, actorID string, actorRole models.AccountRole, req models.CreateWordProposalRequest) (*models.WordDetail, error)
	UpdateWord(ctx context.Context, actorID string, actorRole models.AccountRole, wordID int, req models.UpdateWordRequest) (*models.WordDetail, error)
	CreateSense(ctx context.Context, actorID string, actorRole models.AccountRole, wordID int, req models.CreateSenseRequest) (*models.SenseDetail, error)
	UpdateSense(ctx context.Context, actorID string, actorRole models.AccountRole, senseID int, req models.UpdateSenseRequest) (*models.SenseDetail, error)
	DeleteSense(ctx context.Context, actorID string, actorRole models.AccountRole, senseID int) error
	CreateExample(ctx context.Context, actorID string, actorRole models.AccountRole, senseID int, req models.CreateExampleInput) (*models.Example, error)
	UpdateExample(ctx context.Context, actorID string, actorRole models.AccountRole, exampleID int, req models.CreateExampleInput) (*models.Example, error)
	DeleteExample(ctx context.Context, actorID string, actorRole models.AccountRole, exampleID int) error
	ListSemanticConcepts(ctx context.Context, search string) ([]models.SemanticConcept, error)
	CreateSemanticConcept(ctx context.Context, req models.CreateConceptRequest) (*models.SemanticConcept, error)
	LinkSenseToConcept(ctx context.Context, senseID int, conceptID int) error
	UnlinkSenseFromConcept(ctx context.Context, senseID int, conceptID int) error
}

type dictionaryService struct {
	db      *database.DB
	queries *db.Queries
}

func NewDictionaryService(databaseInstance *database.DB) DictionaryService {
	return &dictionaryService{
		db:      databaseInstance,
		queries: db.New(databaseInstance.Pool),
	}
}

// NormalizeLemma convierte un lema a minúsculas y elimina los signos diacríticos (tildes/diéresis).
func NormalizeLemma(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	var b strings.Builder
	for _, r := range s {
		switch r {
		case 'á', 'à', 'ä', 'â':
			b.WriteRune('a')
		case 'é', 'è', 'ë', 'ê':
			b.WriteRune('e')
		case 'í', 'ì', 'ï', 'î':
			b.WriteRune('i')
		case 'ó', 'ò', 'ö', 'ô':
			b.WriteRune('o')
		case 'ú', 'ù', 'ü', 'û':
			b.WriteRune('u')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (s *dictionaryService) SearchDropdown(ctx context.Context, query string, userDialect string, limit int) ([]models.WordSearchResult, error) {
	if userDialect == "" {
		userDialect = "MX"
	}
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	rows, err := s.queries.SearchVocabularyDropdown(ctx, db.SearchVocabularyDropdownParams{
		SearchQuery: query,
		UserDialect: userDialect,
		MaxResults:  int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("error al buscar palabras: %w", err)
	}

	results := make([]models.WordSearchResult, len(rows))
	for i, row := range rows {
		results[i] = models.WordSearchResult{
			IDWord:          int(row.IDWord),
			Lemma:           row.Lemma,
			ShortDefinition: row.ShortDefinition,
			FlagEmoji:       row.FlagEmoji,
			CountryCode:     row.CountryCode,
			Similarity:      row.Similarity,
		}
	}

	return results, nil
}

func (s *dictionaryService) GetWordDetail(ctx context.Context, wordID int) (*models.WordDetail, error) {
	if wordID <= 0 {
		return nil, NewDomainError(ErrInvalidInput, "El identificador de palabra debe ser un número entero positivo")
	}

	// 1. Obtener la palabra canónica
	wordRow, err := s.queries.GetWordDetailByID(ctx, int32(wordID))
	if err != nil {
		return nil, MapDBError(err, "Palabra no encontrada en el diccionario")
	}

	// 2. Obtener acepciones aprobadas
	sensesRows, err := s.queries.GetApprovedSensesByWordID(ctx, int32(wordID))
	if err != nil {
		return nil, fmt.Errorf("error consultando acepciones: %w", err)
	}

	// 3. Obtener ejemplos de uso
	examplesRows, err := s.queries.GetExamplesByWordID(ctx, int32(wordID))
	if err != nil {
		return nil, fmt.Errorf("error consultando ejemplos: %w", err)
	}

	// Indexar ejemplos por id_sense
	examplesMap := make(map[int32][]models.Example)
	for _, ex := range examplesRows {
		examplesMap[ex.IDSense] = append(examplesMap[ex.IDSense], models.Example{
			IDExample: int(ex.IDExample),
			IDSense:   int(ex.IDSense),
			Sentence:  ex.Sentence,
			Source:    TextToPtr(ex.Source),
		})
	}

	// Construir acepciones con sus respectivos ejemplos
	senses := make([]models.SenseDetail, len(sensesRows))
	for i, sr := range sensesRows {
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

		senses[i] = models.SenseDetail{
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
			Examples:        examplesMap[sr.IDSense],
		}
		if senses[i].Examples == nil {
			senses[i].Examples = []models.Example{}
		}
	}

	// 4. Obtener heterónimos (¿Cómo se dice en otros países?)
	heteronymsRows, err := s.queries.GetHeteronymsByWordID(ctx, int32(wordID))
	if err != nil {
		return nil, fmt.Errorf("error consultando heterónimos: %w", err)
	}

	heteronyms := make([]models.HeteronymItem, len(heteronymsRows))
	for i, h := range heteronymsRows {
		heteronyms[i] = models.HeteronymItem{
			RegionCode:  h.RegionCode,
			RegionName:  h.RegionName,
			CountryCode: h.CountryCode,
			FlagEmoji:   h.FlagEmoji,
			Term:        h.Term,
		}
	}

	// 5. Construir objeto de región de origen si existe
	var originRegion *models.Region
	if wordRow.OriginRegionID.Valid {
		originUUID, _ := uuid.FromBytes(wordRow.OriginRegionID.Bytes[:])
		var countryPtr *string
		if wordRow.OriginRegionCountryCode.Valid {
			c := wordRow.OriginRegionCountryCode.String
			countryPtr = &c
		}

		originRegion = &models.Region{
			IDRegion:    originUUID.String(),
			Code:        wordRow.OriginRegionCode.String,
			Name:        wordRow.OriginRegionName.String,
			ParentZone:  wordRow.OriginRegionParentZone.String,
			CountryCode: countryPtr,
			FlagEmoji:   wordRow.OriginRegionFlagEmoji.String,
		}
	}

	genderStr := models.GenderNeutral
	if wordRow.Gender.Valid {
		genderStr = models.WordGender(wordRow.Gender.WordGender)
	}

	return &models.WordDetail{
		IDWord:          int(wordRow.IDWord),
		Lemma:           wordRow.Lemma,
		NormalizedLemma: wordRow.NormalizedLemma,
		PartOfSpeech:    wordRow.PartOfSpeech,
		Gender:          genderStr,
		Etymology:       TextToPtr(wordRow.Etymology),
		FriendlyBlurb:   TextToPtr(wordRow.FriendlyBlurb),
		OriginRegion:    originRegion,
		Senses:          senses,
		Heteronyms:      heteronyms,
		CreatedAt:       wordRow.CreatedAt.Time,
		UpdatedAt:       wordRow.UpdatedAt.Time,
	}, nil
}

func (s *dictionaryService) ListRegions(ctx context.Context) ([]models.Region, error) {
	rows, err := s.queries.ListRegions(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listando regiones: %w", err)
	}

	results := make([]models.Region, len(rows))
	for i, r := range rows {
		var cc *string
		if r.CountryCode.Valid {
			c := r.CountryCode.String
			cc = &c
		}
		results[i] = models.Region{
			IDRegion:    r.IDRegion.String(),
			Code:        r.Code,
			Name:        r.Name,
			ParentZone:  r.ParentZone,
			CountryCode: cc,
			FlagEmoji:   r.FlagEmoji,
		}
	}
	return results, nil
}

func (s *dictionaryService) GetWordOfTheDay(ctx context.Context, dialect string) (*models.WordDetail, error) {
	dialect = strings.TrimSpace(strings.ToUpper(dialect))
	var regText pgtype.Text
	if dialect != "" {
		regText = pgtype.Text{String: dialect, Valid: true}
	}

	wordID, err := s.queries.GetWordOfTheDayID(ctx, regText)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) && regText.Valid {
			// Fallback sin filtro de región
			wordID, err = s.queries.GetWordOfTheDayID(ctx, pgtype.Text{Valid: false})
		}
		if err != nil {
			return nil, MapDBError(err, "No hay palabras disponibles para la palabra del día")
		}
	}

	return s.GetWordDetail(ctx, int(wordID))
}

func (s *dictionaryService) ExploreWords(ctx context.Context, filter models.WordExploreFilter) (*models.PaginatedWords, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	offset := (page - 1) * limit

	var qText, regText, posText, cefrText, freqText pgtype.Text
	if strings.TrimSpace(filter.Query) != "" {
		qText = pgtype.Text{String: strings.TrimSpace(filter.Query), Valid: true}
	}
	if strings.TrimSpace(filter.RegionCode) != "" {
		regText = pgtype.Text{String: strings.TrimSpace(strings.ToUpper(filter.RegionCode)), Valid: true}
	}
	if strings.TrimSpace(filter.PartOfSpeech) != "" {
		posText = pgtype.Text{String: strings.TrimSpace(strings.ToLower(filter.PartOfSpeech)), Valid: true}
	}
	if strings.TrimSpace(filter.CEFRLevel) != "" {
		cefrText = pgtype.Text{String: strings.TrimSpace(strings.ToUpper(filter.CEFRLevel)), Valid: true}
	}
	if strings.TrimSpace(filter.UsageFrequency) != "" {
		freqText = pgtype.Text{String: strings.TrimSpace(strings.ToUpper(filter.UsageFrequency)), Valid: true}
	}

	rows, err := s.queries.ExploreWordsPaginated(ctx, db.ExploreWordsPaginatedParams{
		Limit:          int32(limit),
		Offset:         int32(offset),
		RegionCode:     regText,
		Query:          qText,
		PartOfSpeech:   posText,
		CefrLevel:      cefrText,
		UsageFrequency: freqText,
		IncludeTaboo:   filter.IncludeTaboo,
	})
	if err != nil {
		return nil, fmt.Errorf("error explorando palabras: %w", err)
	}

	totalCount, err := s.queries.CountExploreWords(ctx, db.CountExploreWordsParams{
		RegionCode:     regText,
		Query:          qText,
		PartOfSpeech:   posText,
		CefrLevel:      cefrText,
		UsageFrequency: freqText,
		IncludeTaboo:   filter.IncludeTaboo,
	})
	if err != nil {
		return nil, fmt.Errorf("error contando palabras de exploración: %w", err)
	}

	items := make([]models.WordSummaryCard, len(rows))
	for i, r := range rows {
		genderStr := models.GenderNeutral
		if r.Gender.Valid {
			genderStr = models.WordGender(r.Gender.WordGender)
		}

		items[i] = models.WordSummaryCard{
			IDWord:          int(r.IDWord),
			Lemma:           r.Lemma,
			PartOfSpeech:    r.PartOfSpeech,
			Gender:          genderStr,
			FriendlyBlurb:   TextToPtr(r.FriendlyBlurb),
			PrimarySenseID:  int(r.PrimarySenseID),
			ShortDefinition: r.ShortDefinition,
			RegionCode:      r.RegionCode,
			RegionName:      r.RegionName,
			FlagEmoji:       r.FlagEmoji,
			CEFRLevel:       interfaceToString(r.CefrLevel, "B1"),
			UsageFrequency:  interfaceToString(r.UsageFrequency, "COMMON"),
			IsTabooOrVulgar: r.IsTabooOrVulgar,
			ApprovedSenses:  int(r.ApprovedSensesCount),
		}
	}

	return &models.PaginatedWords{
		Items:      items,
		TotalCount: int(totalCount),
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *dictionaryService) CreateWordProposal(ctx context.Context, actorID string, actorRole models.AccountRole, req models.CreateWordProposalRequest) (*models.WordDetail, error) {
	lemma := strings.TrimSpace(req.Lemma)
	if lemma == "" {
		return nil, NewDomainError(ErrInvalidInput, "El lema de la palabra es obligatorio")
	}
	pos := strings.TrimSpace(strings.ToLower(req.PartOfSpeech))
	if pos == "" {
		return nil, NewDomainError(ErrInvalidInput, "La categoría gramatical (part_of_speech) es obligatoria")
	}

	initSense := req.InitialSense
	if strings.TrimSpace(initSense.RegionCode) == "" {
		return nil, NewDomainError(ErrInvalidInput, "El código de región para la acepción inicial es obligatorio")
	}
	if strings.TrimSpace(initSense.ShortDefinition) == "" {
		return nil, NewDomainError(ErrInvalidInput, "La definición corta es obligatoria")
	}
	if len(initSense.ShortDefinition) > 255 {
		return nil, NewDomainError(ErrInvalidInput, "La definición corta no debe exceder 255 caracteres")
	}
	if strings.TrimSpace(initSense.Definition) == "" {
		return nil, NewDomainError(ErrInvalidInput, "La definición detallada es obligatoria")
	}

	actorUUID, err := uuid.Parse(actorID)
	if err != nil {
		return nil, NewDomainError(ErrUnauthorized, "Identificador de usuario no válido")
	}

	// 1. Validar región del sentido
	senseRegion, err := s.queries.GetRegionByCode(ctx, strings.ToUpper(strings.TrimSpace(initSense.RegionCode)))
	if err != nil {
		return nil, NewDomainError(ErrInvalidInput, "Región dialectal de la acepción no válida")
	}

	// 2. Resolver región de origen opcional
	var originRegionUUID pgtype.UUID
	if req.OriginRegionCode != nil && strings.TrimSpace(*req.OriginRegionCode) != "" {
		orRow, err := s.queries.GetRegionByCode(ctx, strings.ToUpper(strings.TrimSpace(*req.OriginRegionCode)))
		if err == nil {
			originRegionUUID = pgtype.UUID{Bytes: orRow.IDRegion, Valid: true}
		}
	}

	normalized := NormalizeLemma(lemma)

	// 3. Buscar o crear la palabra
	var wordID int32
	existingWord, err := s.queries.GetWordByNormalizedLemmaAndPOS(ctx, db.GetWordByNormalizedLemmaAndPOSParams{
		NormalizedLemma: normalized,
		PartOfSpeech:    pos,
	})
	if err == nil {
		wordID = existingWord.IDWord
	} else if errors.Is(err, pgx.ErrNoRows) {
		var etymText, blurbText pgtype.Text
		if req.Etymology != nil && strings.TrimSpace(*req.Etymology) != "" {
			etymText = pgtype.Text{String: strings.TrimSpace(*req.Etymology), Valid: true}
		}
		if req.FriendlyBlurb != nil && strings.TrimSpace(*req.FriendlyBlurb) != "" {
			blurbText = pgtype.Text{String: strings.TrimSpace(*req.FriendlyBlurb), Valid: true}
		}
		genderVal := db.WordGenderNEUTRAL
		if req.Gender != "" {
			genderVal = db.WordGender(req.Gender)
		}

		createdWord, err := s.queries.CreateWord(ctx, db.CreateWordParams{
			Lemma:           lemma,
			NormalizedLemma: normalized,
			PartOfSpeech:    pos,
			Gender:          db.NullWordGender{WordGender: genderVal, Valid: true},
			Etymology:       etymText,
			FriendlyBlurb:   blurbText,
			IDOriginRegion:  originRegionUUID,
		})
		if err != nil {
			return nil, MapDBError(err, "Error al registrar la palabra en el diccionario")
		}
		wordID = createdWord.IDWord
	} else {
		return nil, MapDBError(err, "Error al comprobar existencia de la palabra")
	}

	// 4. Determinar estado inicial (RN08: ADMIN puede publicar directamente; EDITOR crea DRAFT)
	initialState := db.ReviewStateDRAFT
	if actorRole == models.RoleAdmin {
		initialState = db.ReviewStateAPPROVED
	}

	var cultNote pgtype.Text
	if initSense.CulturalNote != nil && strings.TrimSpace(*initSense.CulturalNote) != "" {
		cultNote = pgtype.Text{String: strings.TrimSpace(*initSense.CulturalNote), Valid: true}
	}

	cefrVal := db.CefrLevelB1
	if initSense.CEFRLevel != "" {
		cefrVal = db.CefrLevel(initSense.CEFRLevel)
	}
	freqVal := db.UsageFrequencyCOMMON
	if initSense.UsageFrequency != "" {
		freqVal = db.UsageFrequency(initSense.UsageFrequency)
	}

	createdSense, err := s.queries.CreateSense(ctx, db.CreateSenseParams{
		IDWord:          wordID,
		IDRegion:        senseRegion.IDRegion,
		IDEditor:        actorUUID,
		ShortDefinition: strings.TrimSpace(initSense.ShortDefinition),
		Definition:      strings.TrimSpace(initSense.Definition),
		IsPrimary:       initSense.IsPrimary,
		IsTabooOrVulgar: initSense.IsTabooOrVulgar,
		CulturalNote:    cultNote,
		CefrLevel:       db.NullCefrLevel{CefrLevel: cefrVal, Valid: true},
		UsageFrequency:  db.NullUsageFrequency{UsageFrequency: freqVal, Valid: true},
		State:           initialState,
	})
	if err != nil {
		return nil, MapDBError(err, "Error al crear la acepción inicial")
	}

	// 5. Insertar ejemplos iniciales si se suministraron
	for _, ex := range initSense.Examples {
		sent := strings.TrimSpace(ex.Sentence)
		if len(sent) >= 5 {
			var src pgtype.Text
			if ex.Source != nil && strings.TrimSpace(*ex.Source) != "" {
				src = pgtype.Text{String: strings.TrimSpace(*ex.Source), Valid: true}
			}
			_, _ = s.queries.CreateExample(ctx, db.CreateExampleParams{
				IDSense:  createdSense.IDSense,
				IDEditor: actorUUID,
				Sentence: sent,
				Source:   src,
			})
		}
	}

	// 6. Vincular conceptos semánticos si se especificaron
	for _, cID := range initSense.ConceptIDs {
		_ = s.queries.LinkSenseToConcept(ctx, db.LinkSenseToConceptParams{
			IDConcept: int32(cID),
			IDSense:   createdSense.IDSense,
		})
	}

	return s.GetWordDetail(ctx, int(wordID))
}

func (s *dictionaryService) UpdateWord(ctx context.Context, actorID string, actorRole models.AccountRole, wordID int, req models.UpdateWordRequest) (*models.WordDetail, error) {
	lemma := strings.TrimSpace(req.Lemma)
	if lemma == "" {
		return nil, NewDomainError(ErrInvalidInput, "El lema de la palabra es obligatorio")
	}
	pos := strings.TrimSpace(strings.ToLower(req.PartOfSpeech))
	if pos == "" {
		return nil, NewDomainError(ErrInvalidInput, "La categoría gramatical es obligatoria")
	}

	var originRegionUUID pgtype.UUID
	if req.OriginRegionCode != nil && strings.TrimSpace(*req.OriginRegionCode) != "" {
		orRow, err := s.queries.GetRegionByCode(ctx, strings.ToUpper(strings.TrimSpace(*req.OriginRegionCode)))
		if err == nil {
			originRegionUUID = pgtype.UUID{Bytes: orRow.IDRegion, Valid: true}
		}
	}

	var etymText, blurbText pgtype.Text
	if req.Etymology != nil && strings.TrimSpace(*req.Etymology) != "" {
		etymText = pgtype.Text{String: strings.TrimSpace(*req.Etymology), Valid: true}
	}
	if req.FriendlyBlurb != nil && strings.TrimSpace(*req.FriendlyBlurb) != "" {
		blurbText = pgtype.Text{String: strings.TrimSpace(*req.FriendlyBlurb), Valid: true}
	}
	genderVal := db.WordGenderNEUTRAL
	if req.Gender != "" {
		genderVal = db.WordGender(req.Gender)
	}

	_, err := s.queries.UpdateWordMetadata(ctx, db.UpdateWordMetadataParams{
		IDWord:          int32(wordID),
		Lemma:           lemma,
		NormalizedLemma: NormalizeLemma(lemma),
		PartOfSpeech:    pos,
		Gender:          db.NullWordGender{WordGender: genderVal, Valid: true},
		Etymology:       etymText,
		FriendlyBlurb:   blurbText,
		IDOriginRegion:  originRegionUUID,
	})
	if err != nil {
		return nil, MapDBError(err, "Error al actualizar metadatos de la palabra")
	}

	return s.GetWordDetail(ctx, wordID)
}

func (s *dictionaryService) CreateSense(ctx context.Context, actorID string, actorRole models.AccountRole, wordID int, req models.CreateSenseRequest) (*models.SenseDetail, error) {
	if strings.TrimSpace(req.RegionCode) == "" {
		return nil, NewDomainError(ErrInvalidInput, "El código de región es obligatorio")
	}
	if strings.TrimSpace(req.ShortDefinition) == "" {
		return nil, NewDomainError(ErrInvalidInput, "La definición corta es obligatoria")
	}
	if len(req.ShortDefinition) > 255 {
		return nil, NewDomainError(ErrInvalidInput, "La definición corta no debe exceder 255 caracteres")
	}
	if strings.TrimSpace(req.Definition) == "" {
		return nil, NewDomainError(ErrInvalidInput, "La definición detallada es obligatoria")
	}

	actorUUID, err := uuid.Parse(actorID)
	if err != nil {
		return nil, NewDomainError(ErrUnauthorized, "Identificador de usuario no válido")
	}

	// Comprobar que la palabra exista
	_, err = s.queries.GetWordDetailByID(ctx, int32(wordID))
	if err != nil {
		return nil, MapDBError(err, "Palabra no encontrada")
	}

	senseRegion, err := s.queries.GetRegionByCode(ctx, strings.ToUpper(strings.TrimSpace(req.RegionCode)))
	if err != nil {
		return nil, NewDomainError(ErrInvalidInput, "Región dialectal no válida")
	}

	initialState := db.ReviewStateDRAFT
	if actorRole == models.RoleAdmin {
		initialState = db.ReviewStateAPPROVED
	}

	var cultNote pgtype.Text
	if req.CulturalNote != nil && strings.TrimSpace(*req.CulturalNote) != "" {
		cultNote = pgtype.Text{String: strings.TrimSpace(*req.CulturalNote), Valid: true}
	}

	cefrVal := db.CefrLevelB1
	if req.CEFRLevel != "" {
		cefrVal = db.CefrLevel(req.CEFRLevel)
	}
	freqVal := db.UsageFrequencyCOMMON
	if req.UsageFrequency != "" {
		freqVal = db.UsageFrequency(req.UsageFrequency)
	}

	createdSense, err := s.queries.CreateSense(ctx, db.CreateSenseParams{
		IDWord:          int32(wordID),
		IDRegion:        senseRegion.IDRegion,
		IDEditor:        actorUUID,
		ShortDefinition: strings.TrimSpace(req.ShortDefinition),
		Definition:      strings.TrimSpace(req.Definition),
		IsPrimary:       req.IsPrimary,
		IsTabooOrVulgar: req.IsTabooOrVulgar,
		CulturalNote:    cultNote,
		CefrLevel:       db.NullCefrLevel{CefrLevel: cefrVal, Valid: true},
		UsageFrequency:  db.NullUsageFrequency{UsageFrequency: freqVal, Valid: true},
		State:           initialState,
	})
	if err != nil {
		return nil, MapDBError(err, "Error al crear la acepción")
	}

	// Insertar ejemplos iniciales
	for _, ex := range req.Examples {
		sent := strings.TrimSpace(ex.Sentence)
		if len(sent) >= 5 {
			var src pgtype.Text
			if ex.Source != nil && strings.TrimSpace(*ex.Source) != "" {
				src = pgtype.Text{String: strings.TrimSpace(*ex.Source), Valid: true}
			}
			_, _ = s.queries.CreateExample(ctx, db.CreateExampleParams{
				IDSense:  createdSense.IDSense,
				IDEditor: actorUUID,
				Sentence: sent,
				Source:   src,
			})
		}
	}

	for _, cID := range req.ConceptIDs {
		_ = s.queries.LinkSenseToConcept(ctx, db.LinkSenseToConceptParams{
			IDConcept: int32(cID),
			IDSense:   createdSense.IDSense,
		})
	}

	return s.getSenseDetailByID(ctx, int(createdSense.IDSense))
}

func (s *dictionaryService) UpdateSense(ctx context.Context, actorID string, actorRole models.AccountRole, senseID int, req models.UpdateSenseRequest) (*models.SenseDetail, error) {
	senseRow, err := s.queries.GetSenseByID(ctx, int32(senseID))
	if err != nil {
		return nil, MapDBError(err, "Acepción no encontrada")
	}

	actorUUID, err := uuid.Parse(actorID)
	if err != nil {
		return nil, NewDomainError(ErrUnauthorized, "Identificador de usuario no válido")
	}

	// CU06 2a: Si el actor es EDITOR (no ADMIN), solo puede modificar acepciones propias en DRAFT
	if actorRole != models.RoleAdmin {
		if senseRow.IDEditor != actorUUID || senseRow.State != db.ReviewStateDRAFT {
			return nil, NewDomainError(ErrForbidden, "Solo puedes modificar acepciones en borrador de tu propia autoría")
		}
	}

	senseRegion, err := s.queries.GetRegionByCode(ctx, strings.ToUpper(strings.TrimSpace(req.RegionCode)))
	if err != nil {
		return nil, NewDomainError(ErrInvalidInput, "Región dialectal no válida")
	}

	var cultNote pgtype.Text
	if req.CulturalNote != nil && strings.TrimSpace(*req.CulturalNote) != "" {
		cultNote = pgtype.Text{String: strings.TrimSpace(*req.CulturalNote), Valid: true}
	}

	cefrVal := db.CefrLevelB1
	if req.CEFRLevel != "" {
		cefrVal = db.CefrLevel(req.CEFRLevel)
	}
	freqVal := db.UsageFrequencyCOMMON
	if req.UsageFrequency != "" {
		freqVal = db.UsageFrequency(req.UsageFrequency)
	}

	_, err = s.queries.UpdateSense(ctx, db.UpdateSenseParams{
		IDSense:         int32(senseID),
		IDRegion:        senseRegion.IDRegion,
		ShortDefinition: strings.TrimSpace(req.ShortDefinition),
		Definition:      strings.TrimSpace(req.Definition),
		IsPrimary:       req.IsPrimary,
		IsTabooOrVulgar: req.IsTabooOrVulgar,
		CulturalNote:    cultNote,
		CefrLevel:       db.NullCefrLevel{CefrLevel: cefrVal, Valid: true},
		UsageFrequency:  db.NullUsageFrequency{UsageFrequency: freqVal, Valid: true},
	})
	if err != nil {
		return nil, MapDBError(err, "Error al actualizar la acepción")
	}

	return s.getSenseDetailByID(ctx, senseID)
}

func (s *dictionaryService) DeleteSense(ctx context.Context, actorID string, actorRole models.AccountRole, senseID int) error {
	senseRow, err := s.queries.GetSenseByID(ctx, int32(senseID))
	if err != nil {
		return MapDBError(err, "Acepción no encontrada")
	}

	actorUUID, err := uuid.Parse(actorID)
	if err != nil {
		return NewDomainError(ErrUnauthorized, "Identificador de usuario no válido")
	}

	if actorRole != models.RoleAdmin {
		if senseRow.IDEditor != actorUUID || senseRow.State != db.ReviewStateDRAFT {
			return NewDomainError(ErrForbidden, "Solo puedes eliminar acepciones en borrador de tu propia autoría")
		}
	}

	if err := s.queries.DeleteSense(ctx, int32(senseID)); err != nil {
		return MapDBError(err, "Error al eliminar la acepción")
	}
	return nil
}

func (s *dictionaryService) CreateExample(ctx context.Context, actorID string, actorRole models.AccountRole, senseID int, req models.CreateExampleInput) (*models.Example, error) {
	sent := strings.TrimSpace(req.Sentence)
	if len(sent) < 5 {
		return nil, NewDomainError(ErrInvalidInput, "La oración de ejemplo debe tener al menos 5 caracteres")
	}

	actorUUID, err := uuid.Parse(actorID)
	if err != nil {
		return nil, NewDomainError(ErrUnauthorized, "Identificador de usuario no válido")
	}

	senseRow, err := s.queries.GetSenseByID(ctx, int32(senseID))
	if err != nil {
		return nil, MapDBError(err, "Acepción no encontrada")
	}

	if actorRole != models.RoleAdmin {
		if senseRow.IDEditor != actorUUID || senseRow.State != db.ReviewStateDRAFT {
			return nil, NewDomainError(ErrForbidden, "No tienes permiso para agregar ejemplos a esta acepción")
		}
	}

	var src pgtype.Text
	if req.Source != nil && strings.TrimSpace(*req.Source) != "" {
		src = pgtype.Text{String: strings.TrimSpace(*req.Source), Valid: true}
	}

	exRow, err := s.queries.CreateExample(ctx, db.CreateExampleParams{
		IDSense:  int32(senseID),
		IDEditor: actorUUID,
		Sentence: sent,
		Source:   src,
	})
	if err != nil {
		return nil, MapDBError(err, "Error al registrar ejemplo")
	}

	return &models.Example{
		IDExample: int(exRow.IDExample),
		IDSense:   int(exRow.IDSense),
		Sentence:  exRow.Sentence,
		Source:    TextToPtr(exRow.Source),
	}, nil
}

func (s *dictionaryService) UpdateExample(ctx context.Context, actorID string, actorRole models.AccountRole, exampleID int, req models.CreateExampleInput) (*models.Example, error) {
	sent := strings.TrimSpace(req.Sentence)
	if len(sent) < 5 {
		return nil, NewDomainError(ErrInvalidInput, "La oración de ejemplo debe tener al menos 5 caracteres")
	}

	actorUUID, err := uuid.Parse(actorID)
	if err != nil {
		return nil, NewDomainError(ErrUnauthorized, "Identificador de usuario no válido")
	}

	exRow, err := s.queries.GetExampleByID(ctx, int32(exampleID))
	if err != nil {
		return nil, MapDBError(err, "Ejemplo no encontrado")
	}

	if actorRole != models.RoleAdmin {
		if exRow.IDEditor != actorUUID {
			return nil, NewDomainError(ErrForbidden, "Solo puedes modificar ejemplos de tu propia autoría")
		}
	}

	var src pgtype.Text
	if req.Source != nil && strings.TrimSpace(*req.Source) != "" {
		src = pgtype.Text{String: strings.TrimSpace(*req.Source), Valid: true}
	}

	updated, err := s.queries.UpdateExample(ctx, db.UpdateExampleParams{
		IDExample: int32(exampleID),
		Sentence:  sent,
		Source:    src,
	})
	if err != nil {
		return nil, MapDBError(err, "Error al actualizar ejemplo")
	}

	return &models.Example{
		IDExample: int(updated.IDExample),
		IDSense:   int(updated.IDSense),
		Sentence:  updated.Sentence,
		Source:    TextToPtr(updated.Source),
	}, nil
}

func (s *dictionaryService) DeleteExample(ctx context.Context, actorID string, actorRole models.AccountRole, exampleID int) error {
	actorUUID, err := uuid.Parse(actorID)
	if err != nil {
		return NewDomainError(ErrUnauthorized, "Identificador de usuario no válido")
	}

	exRow, err := s.queries.GetExampleByID(ctx, int32(exampleID))
	if err != nil {
		return MapDBError(err, "Ejemplo no encontrado")
	}

	if actorRole != models.RoleAdmin {
		if exRow.IDEditor != actorUUID {
			return NewDomainError(ErrForbidden, "Solo puedes eliminar ejemplos de tu propia autoría")
		}
	}

	if err := s.queries.DeleteExample(ctx, int32(exampleID)); err != nil {
		return MapDBError(err, "Error al eliminar ejemplo")
	}
	return nil
}

func (s *dictionaryService) ListSemanticConcepts(ctx context.Context, search string) ([]models.SemanticConcept, error) {
	var qText pgtype.Text
	if strings.TrimSpace(search) != "" {
		qText = pgtype.Text{String: strings.TrimSpace(search), Valid: true}
	}

	rows, err := s.queries.ListSemanticConcepts(ctx, qText)
	if err != nil {
		return nil, fmt.Errorf("error listando conceptos: %w", err)
	}

	results := make([]models.SemanticConcept, len(rows))
	for i, r := range rows {
		results[i] = models.SemanticConcept{
			IDConcept:            int(r.IDConcept),
			ConceptKey:           r.ConceptKey,
			UniversalDescription: r.UniversalDescription,
			LinkedSensesCount:    int(r.LinkedSensesCount),
		}
	}
	return results, nil
}

func (s *dictionaryService) CreateSemanticConcept(ctx context.Context, req models.CreateConceptRequest) (*models.SemanticConcept, error) {
	key := strings.TrimSpace(req.ConceptKey)
	desc := strings.TrimSpace(req.UniversalDescription)
	if key == "" || desc == "" {
		return nil, NewDomainError(ErrInvalidInput, "concept_key y universal_description son obligatorios")
	}

	row, err := s.queries.CreateSemanticConcept(ctx, db.CreateSemanticConceptParams{
		ConceptKey:  key,
		Description: desc,
	})
	if err != nil {
		return nil, MapDBError(err, "Error creando concepto semántico")
	}

	return &models.SemanticConcept{
		IDConcept:            int(row.IDConcept),
		ConceptKey:           row.ConceptKey,
		UniversalDescription: row.UniversalDescription,
		LinkedSensesCount:    0,
	}, nil
}

func (s *dictionaryService) LinkSenseToConcept(ctx context.Context, senseID int, conceptID int) error {
	err := s.queries.LinkSenseToConcept(ctx, db.LinkSenseToConceptParams{
		IDConcept: int32(conceptID),
		IDSense:   int32(senseID),
	})
	if err != nil {
		return MapDBError(err, "Error vinculando concepto a la acepción")
	}
	return nil
}

func (s *dictionaryService) UnlinkSenseFromConcept(ctx context.Context, senseID int, conceptID int) error {
	err := s.queries.UnlinkSenseFromConcept(ctx, db.UnlinkSenseFromConceptParams{
		IDConcept: int32(conceptID),
		IDSense:   int32(senseID),
	})
	if err != nil {
		return MapDBError(err, "Error desvinculando concepto de la acepción")
	}
	return nil
}

// Helper interno para obtener un SenseDetail completo
func (s *dictionaryService) getSenseDetailByID(ctx context.Context, senseID int) (*models.SenseDetail, error) {
	sr, err := s.queries.GetSenseByID(ctx, int32(senseID))
	if err != nil {
		return nil, MapDBError(err, "Acepción no encontrada")
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

func interfaceToString(val interface{}, defaultVal string) string {
	if val == nil {
		return defaultVal
	}
	switch v := val.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
