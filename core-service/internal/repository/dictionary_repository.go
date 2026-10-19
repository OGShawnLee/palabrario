package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"palabrario/core-service/internal/database"
	"palabrario/core-service/internal/db"
	"palabrario/pkg/models"
)

var ErrWordNotFound = errors.New("palabra no encontrada")

type DictionaryRepository struct {
	db      *database.DB
	queries *db.Queries
}

func NewDictionaryRepository(databaseInstance *database.DB) *DictionaryRepository {
	return &DictionaryRepository{
		db:      databaseInstance,
		queries: db.New(databaseInstance.Pool),
	}
}

func (r *DictionaryRepository) SearchWords(ctx context.Context, query, dialect string, limit int) ([]models.WordSearchResult, error) {
	if dialect == "" {
		dialect = "MX"
	}
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	rows, err := r.queries.SearchVocabularyDropdown(ctx, db.SearchVocabularyDropdownParams{
		SearchQuery: query,
		UserDialect: dialect,
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

func (r *DictionaryRepository) GetWordDetail(ctx context.Context, wordID int) (*models.WordDetail, error) {
	// 1. Obtener la palabra canónica
	wordRow, err := r.queries.GetWordDetailByID(ctx, int32(wordID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWordNotFound
		}
		return nil, fmt.Errorf("error consultando detalle de la palabra: %w", err)
	}

	// 2. Obtener acepciones aprobadas
	sensesRows, err := r.queries.GetApprovedSensesByWordID(ctx, int32(wordID))
	if err != nil {
		return nil, fmt.Errorf("error consultando acepciones: %w", err)
	}

	// 3. Obtener ejemplos de uso
	examplesRows, err := r.queries.GetExamplesByWordID(ctx, int32(wordID))
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
			Source:    textToPtr(ex.Source),
		})
	}

	// Construir acepciones con sus respectivos ejemplos
	senses := make([]models.SenseDetail, len(sensesRows))
	for i, s := range sensesRows {
		countryCode := ""
		if s.CountryCode.Valid {
			countryCode = s.CountryCode.String
		}

		cefrStr := ""
		if s.CefrLevel.Valid {
			cefrStr = string(s.CefrLevel.CefrLevel)
		}

		freqStr := ""
		if s.UsageFrequency.Valid {
			freqStr = string(s.UsageFrequency.UsageFrequency)
		}

		senses[i] = models.SenseDetail{
			IDSense:         int(s.IDSense),
			IDWord:          int(s.IDWord),
			RegionCode:      s.RegionCode,
			RegionName:      s.RegionName,
			CountryCode:     countryCode,
			FlagEmoji:       s.FlagEmoji,
			ShortDefinition: s.ShortDefinition,
			Definition:      s.Definition,
			IsPrimary:       s.IsPrimary,
			IsTabooOrVulgar: s.IsTabooOrVulgar,
			CulturalNote:    textToPtr(s.CulturalNote),
			CEFRLevel:       cefrStr,
			UsageFrequency:  freqStr,
			State:           string(s.State),
			Examples:        examplesMap[s.IDSense],
		}
		if senses[i].Examples == nil {
			senses[i].Examples = []models.Example{}
		}
	}

	// 4. Obtener heterónimos (¿Cómo se dice en otros países?)
	heteronymsRows, err := r.queries.GetHeteronymsByWordID(ctx, int32(wordID))
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
		Etymology:       textToPtr(wordRow.Etymology),
		FriendlyBlurb:   textToPtr(wordRow.FriendlyBlurb),
		OriginRegion:    originRegion,
		Senses:          senses,
		Heteronyms:      heteronyms,
		CreatedAt:       wordRow.CreatedAt.Time,
		UpdatedAt:       wordRow.UpdatedAt.Time,
	}, nil
}
