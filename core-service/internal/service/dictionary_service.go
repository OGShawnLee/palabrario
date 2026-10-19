package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"palabrario/core-service/internal/database"
	"palabrario/core-service/internal/db"
	"palabrario/pkg/models"
)

// DictionaryService define el contrato SOA del Módulo 3 (Léxico y Búsqueda Dialectal — CU05, CU06)
// según la sección 2.3 de .context/INTERFACES.md.
// Los métodos adicionales de exploración y autoría (ListRegions, GetWordOfTheDay, ExploreWords,
// CreateWordProposal, CreateSenseProposal, etc.) se añadirán en la Fase 3.3 siguiendo este molde.
type DictionaryService interface {
	SearchDropdown(ctx context.Context, query string, userDialect string, limit int) ([]models.WordSearchResult, error)
	GetWordDetail(ctx context.Context, wordID int) (*models.WordDetail, error)
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
