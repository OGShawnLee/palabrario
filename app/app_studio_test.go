package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"palabrario/pkg/models"
	"palabrario/pkg/utils"
)

func TestListMasteryItemsDecodesStudyCardItemsWithMasteryAndExamples(t *testing.T) {
	sessionFile := filepath.Join(t.TempDir(), "palabrario", "session.json")
	now := time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)
	nextReview := now.Add(21 * 24 * time.Hour)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/studio/mastery/items" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer access-studio" {
			t.Fatalf("expected Authorization header 'Bearer access-studio', got %q", got)
		}

		stage := r.URL.Query().Get("stage")
		switch stage {
		case "MASTERED":
			writeJSON(t, w, http.StatusOK, []models.StudyCardItem{
				{
					IDSense:         201,
					IDWord:          101,
					Lemma:           "Guagua",
					PartOfSpeech:    "sustantivo femenino",
					RegionCode:      "CU",
					RegionName:      "Cuba",
					FlagEmoji:       "",
					ShortDefinition: "Autobús de transporte público",
					Definition:      "Ómnibus de servicio urbano o interurbano destinado al transporte masivo de pasajeros.",
					CulturalNote:    utils.ToPtr("En Cuba, tomar la guagua forma parte de la vida cotidiana."),
					Examples: []models.Example{
						{
							IDExample: 301,
							IDSense:   201,
							Sentence:  "Ayer la guagua pasó llena y tuve que caminar.",
							Source:    utils.ToPtr("Habla urbana · Cuba"),
						},
					},
					CurrentMastery: &models.UserMastery{
						IDUser:              "user-123",
						IDSense:             201,
						ProgressStage:       models.StageMastered,
						MasteryPercentage:   100,
						ConsecutiveCorrect:  5,
						TimesReviewed:       6,
						EasinessFactor:      2.6,
						CurrentIntervalDays: 21,
						NextReviewAt:        nextReview,
						LastReviewedAt:      now,
					},
				},
			})
		case "REVIEWING":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("null"))
		default:
			t.Fatalf("unexpected stage query: %q", stage)
		}
	}))
	defer srv.Close()

	app := newAppWithConfig(srv.URL, srv.URL, srv.URL, sessionFile)
	app.setTokens("access-studio", "refresh-studio")

	items, err := app.ListMasteryItems("MASTERED")
	if err != nil {
		t.Fatalf("ListMasteryItems(MASTERED) failed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 StudyCardItem, got %d", len(items))
	}

	card := items[0]
	if card.IDSense != 201 || card.IDWord != 101 || card.Lemma != "Guagua" || card.RegionName != "Cuba" {
		t.Fatalf("unexpected StudyCardItem identity fields: %+v", card)
	}
	if card.CulturalNote == nil || *card.CulturalNote == "" {
		t.Fatalf("expected CulturalNote decoded, got %v", card.CulturalNote)
	}
	if len(card.Examples) != 1 || card.Examples[0].IDExample != 301 || card.Examples[0].Sentence == "" {
		t.Fatalf("expected nested Examples decoded, got %+v", card.Examples)
	}
	if card.CurrentMastery == nil {
		t.Fatalf("expected nested CurrentMastery decoded, got nil")
	}
	if card.CurrentMastery.ProgressStage != models.StageMastered ||
		card.CurrentMastery.MasteryPercentage != 100 ||
		card.CurrentMastery.ConsecutiveCorrect != 5 ||
		card.CurrentMastery.CurrentIntervalDays != 21 {
		t.Fatalf("unexpected CurrentMastery fields: %+v", card.CurrentMastery)
	}

	emptyItems, err := app.ListMasteryItems("REVIEWING")
	if err != nil {
		t.Fatalf("ListMasteryItems(REVIEWING) failed: %v", err)
	}
	if emptyItems == nil || len(emptyItems) != 0 {
		t.Fatalf("expected non-nil empty []models.StudyCardItem, got %#v", emptyItems)
	}
}

func TestGetDeckPreserves404And403APIErrors(t *testing.T) {
	sessionFile := filepath.Join(t.TempDir(), "palabrario", "session.json")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/decks/999":
			writeJSON(t, w, http.StatusNotFound, map[string]string{
				"code":  "NOT_FOUND",
				"error": "Baraja no encontrada",
			})
		case "/api/v1/decks/42":
			writeJSON(t, w, http.StatusForbidden, map[string]string{
				"code":  "FORBIDDEN",
				"error": "Violación de regla RN04: No tienes acceso a esta baraja privada",
			})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	app := newAppWithConfig(srv.URL, srv.URL, srv.URL, sessionFile)

	_, err := app.GetDeck(999)
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusNotFound {
		t.Fatalf("expected 404 *APIError, got %T %v", err, err)
	}

	_, err = app.GetDeck(42)
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusForbidden {
		t.Fatalf("expected 403 *APIError, got %T %v", err, err)
	}
}
