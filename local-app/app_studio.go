package main

import (
	"net/http"
	"net/url"
	"strconv"

	"palabrario/pkg/models"
)

// =============================================================================
// MÓDULO 6: Studio Personal, Repetición Espaciada SM-2 y Maestría (CU10–CU13)
// =============================================================================

// ExplorePublicCollections lista colecciones públicas con filtro de destacadas.
func (a *App) ExplorePublicCollections(featured bool) ([]models.Collection, error) {
	return a.ListCollections(featured)
}

// ListCollections lista colecciones temáticas públicas.
func (a *App) ListCollections(featured bool) ([]models.Collection, error) {
	targetURL := a.coreURL + "/api/v1/collections"
	if featured {
		targetURL += "?featured=true"
	}
	var list []models.Collection
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.Collection{}
	}
	return list, nil
}

// GetCollectionDetail obtiene la ficha detallada de una colección.
func (a *App) GetCollectionDetail(id int) (*models.Collection, error) {
	return a.GetCollection(id)
}

// GetCollection obtiene colección por id.
func (a *App) GetCollection(id int) (*models.Collection, error) {
	var col models.Collection
	targetURL := a.coreURL + "/api/v1/collections/" + strconv.Itoa(id)
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &col); err != nil {
		return nil, err
	}
	return &col, nil
}

// ExplorePublicDecks busca y lista barajas públicas.
func (a *App) ExplorePublicDecks(q string, officialOnly bool) ([]models.Deck, error) {
	return a.ListDecks(q, officialOnly)
}

// ListDecks busca barajas públicas.
func (a *App) ListDecks(q string, officialOnly bool) ([]models.Deck, error) {
	params := url.Values{}
	if q != "" {
		params.Set("q", q)
	}
	if officialOnly {
		params.Set("official_only", "true")
	}
	targetURL := a.coreURL + "/api/v1/decks"
	if len(params) > 0 {
		targetURL += "?" + params.Encode()
	}
	var list []models.Deck
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.Deck{}
	}
	return list, nil
}

// GetDeckDetail obtiene el detalle y tarjetas de una baraja.
func (a *App) GetDeckDetail(id int) (*models.Deck, error) {
	return a.GetDeck(id)
}

// GetDeck obtiene baraja por id.
func (a *App) GetDeck(id int) (*models.Deck, error) {
	var deck models.Deck
	targetURL := a.coreURL + "/api/v1/decks/" + strconv.Itoa(id)
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &deck); err != nil {
		return nil, err
	}
	return &deck, nil
}

// ClonePublicDeck guarda una baraja comunitaria en el Studio personal.
func (a *App) ClonePublicDeck(deckId int) error {
	return a.SaveDeck(deckId)
}

// SaveDeck guarda baraja en Studio.
func (a *App) SaveDeck(deckId int) error {
	targetURL := a.coreURL + "/api/v1/decks/" + strconv.Itoa(deckId) + "/save"
	return a.doRequest(a.context(), http.MethodPost, targetURL, nil, nil)
}

// UnsaveDeck retira la baraja guardada del Studio.
func (a *App) UnsaveDeck(deckId int) error {
	targetURL := a.coreURL + "/api/v1/decks/" + strconv.Itoa(deckId) + "/save"
	return a.doRequest(a.context(), http.MethodDelete, targetURL, nil, nil)
}

// ListMyDecks lista las barajas del Studio del usuario logueado.
func (a *App) ListMyDecks() ([]models.StudioDeckSummary, error) {
	return a.ListStudioDecks()
}

// ListStudioDecks lista barajas del Studio.
func (a *App) ListStudioDecks() ([]models.StudioDeckSummary, error) {
	var list []models.StudioDeckSummary
	if err := a.doRequest(a.context(), http.MethodGet, a.coreURL+"/api/v1/studio/decks", nil, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.StudioDeckSummary{}
	}
	return list, nil
}

// CreateDeck crea una nueva baraja de estudio.
func (a *App) CreateDeck(req models.CreateDeckRequest) (*models.Deck, error) {
	var deck models.Deck
	if err := a.doRequest(a.context(), http.MethodPost, a.coreURL+"/api/v1/decks", req, &deck); err != nil {
		return nil, err
	}
	return &deck, nil
}

// UpdateDeck actualiza título o visibilidad de una baraja.
func (a *App) UpdateDeck(id int, req models.UpdateDeckRequest) (*models.Deck, error) {
	var deck models.Deck
	targetURL := a.coreURL + "/api/v1/decks/" + strconv.Itoa(id)
	if err := a.doRequest(a.context(), http.MethodPut, targetURL, req, &deck); err != nil {
		return nil, err
	}
	return &deck, nil
}

// DeleteDeck elimina una baraja de autoría propia.
func (a *App) DeleteDeck(id int) error {
	return a.doRequest(a.context(), http.MethodDelete, a.coreURL+"/api/v1/decks/"+strconv.Itoa(id), nil, nil)
}

// AddSenseToDeck vincula un sentido léxico a una baraja.
func (a *App) AddSenseToDeck(deckId int, req models.AddDeckItemRequest) (*models.DeckItem, error) {
	return a.AddDeckItem(deckId, req)
}

// AddDeckItem añade tarjeta a baraja.
func (a *App) AddDeckItem(deckId int, req models.AddDeckItemRequest) (*models.DeckItem, error) {
	var item models.DeckItem
	targetURL := a.coreURL + "/api/v1/decks/" + strconv.Itoa(deckId) + "/items"
	if err := a.doRequest(a.context(), http.MethodPost, targetURL, req, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

// RemoveSenseFromDeck retira un sentido de una baraja.
func (a *App) RemoveSenseFromDeck(deckId int, senseId int) error {
	return a.RemoveDeckItem(deckId, senseId)
}

// RemoveDeckItem retira tarjeta de baraja.
func (a *App) RemoveDeckItem(deckId int, senseId int) error {
	targetURL := a.coreURL + "/api/v1/decks/" + strconv.Itoa(deckId) + "/items/" + strconv.Itoa(senseId)
	return a.doRequest(a.context(), http.MethodDelete, targetURL, nil, nil)
}

// ListMyCollections lista las colecciones temáticas propias del usuario.
func (a *App) ListMyCollections() ([]models.Collection, error) {
	return a.ListStudioCollections()
}

// ListStudioCollections lista colecciones del Studio.
func (a *App) ListStudioCollections() ([]models.Collection, error) {
	var list []models.Collection
	if err := a.doRequest(a.context(), http.MethodGet, a.coreURL+"/api/v1/studio/collections", nil, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.Collection{}
	}
	return list, nil
}

// CreateCollection crea una nueva colección temática.
func (a *App) CreateCollection(req models.CreateCollectionRequest) (*models.Collection, error) {
	var col models.Collection
	if err := a.doRequest(a.context(), http.MethodPost, a.coreURL+"/api/v1/collections", req, &col); err != nil {
		return nil, err
	}
	return &col, nil
}

// UpdateCollection actualiza una colección temática.
func (a *App) UpdateCollection(id int, req models.UpdateCollectionRequest) (*models.Collection, error) {
	var col models.Collection
	targetURL := a.coreURL + "/api/v1/collections/" + strconv.Itoa(id)
	if err := a.doRequest(a.context(), http.MethodPut, targetURL, req, &col); err != nil {
		return nil, err
	}
	return &col, nil
}

// DeleteCollection elimina una colección temática.
func (a *App) DeleteCollection(id int) error {
	return a.doRequest(a.context(), http.MethodDelete, a.coreURL+"/api/v1/collections/"+strconv.Itoa(id), nil, nil)
}

// SaveCollection guarda una colección comunitaria en el Studio personal.
func (a *App) SaveCollection(id int) error {
	targetURL := a.coreURL + "/api/v1/collections/" + strconv.Itoa(id) + "/save"
	return a.doRequest(a.context(), http.MethodPost, targetURL, nil, nil)
}

// UnsaveCollection retira la colección guardada del Studio.
func (a *App) UnsaveCollection(id int) error {
	targetURL := a.coreURL + "/api/v1/collections/" + strconv.Itoa(id) + "/save"
	return a.doRequest(a.context(), http.MethodDelete, targetURL, nil, nil)
}

// AddDeckToCollection añade una baraja a una colección temática.
func (a *App) AddDeckToCollection(collectionID int, deckID int) (*models.Collection, error) {
	col, err := a.GetCollection(collectionID)
	if err != nil {
		return nil, err
	}

	deckIDs := make([]int, 0, len(col.Decks)+1)
	alreadyExists := false
	for _, d := range col.Decks {
		deckIDs = append(deckIDs, d.IDDeck)
		if d.IDDeck == deckID {
			alreadyExists = true
		}
	}
	if !alreadyExists {
		deckIDs = append(deckIDs, deckID)
	}

	return a.UpdateCollection(collectionID, models.UpdateCollectionRequest{
		Title:       col.Title,
		Description: col.Description,
		IsPublic:    col.IsPublic,
		IsFeatured:  col.IsFeatured,
		IsOfficial:  col.IsOfficial,
		DeckIDs:     deckIDs,
	})
}

// RemoveDeckFromCollection remueve una baraja de una colección temática.
func (a *App) RemoveDeckFromCollection(collectionID int, deckID int) (*models.Collection, error) {
	col, err := a.GetCollection(collectionID)
	if err != nil {
		return nil, err
	}

	deckIDs := make([]int, 0, len(col.Decks))
	for _, d := range col.Decks {
		if d.IDDeck != deckID {
			deckIDs = append(deckIDs, d.IDDeck)
		}
	}

	return a.UpdateCollection(collectionID, models.UpdateCollectionRequest{
		Title:       col.Title,
		Description: col.Description,
		IsPublic:    col.IsPublic,
		IsFeatured:  col.IsFeatured,
		IsOfficial:  col.IsOfficial,
		DeckIDs:     deckIDs,
	})
}

// GetStudySessionQueue genera una sesión de repaso SM-2.
func (a *App) GetStudySessionQueue(deckId int, limit int) (*models.StudySessionPayload, error) {
	return a.GetStudySession(deckId, limit)
}

// GetStudySession obtiene tarjetas programadas para repaso hoy.
func (a *App) GetStudySession(deckId int, limit int) (*models.StudySessionPayload, error) {
	params := url.Values{}
	if deckId > 0 {
		params.Set("deck_id", strconv.Itoa(deckId))
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	targetURL := a.coreURL + "/api/v1/studio/study-session"
	if len(params) > 0 {
		targetURL += "?" + params.Encode()
	}
	var payload models.StudySessionPayload
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

// SubmitCardReview procesa la calificación de repaso aplicando el algoritmo SM-2.
func (a *App) SubmitCardReview(req models.SubmitReviewRequest) (*models.CardReviewResult, error) {
	return a.SubmitReview(req)
}

// SubmitReview califica una tarjeta y actualiza la maestría.
func (a *App) SubmitReview(req models.SubmitReviewRequest) (*models.CardReviewResult, error) {
	var res models.CardReviewResult
	if err := a.doRequest(a.context(), http.MethodPost, a.coreURL+"/api/v1/studio/study-session/review", req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// SubmitStudyReview alias para calificación de tarjeta.
func (a *App) SubmitStudyReview(req models.SubmitReviewRequest) (*models.CardReviewResult, error) {
	return a.SubmitReview(req)
}

// GetMasterySummary obtiene el resumen integral de maestría, racha y logros.
func (a *App) GetMasterySummary() (*models.MasteryDashboard, error) {
	return a.GetMasteryDashboard()
}

// GetMasteryDashboard obtiene métricas y estado del tablero de maestría.
func (a *App) GetMasteryDashboard() (*models.MasteryDashboard, error) {
	var dash models.MasteryDashboard
	if err := a.doRequest(a.context(), http.MethodGet, a.coreURL+"/api/v1/studio/mastery", nil, &dash); err != nil {
		return nil, err
	}
	return &dash, nil
}

// ListMasteryItems lista tarjetas estudiadas con filtro opcional de etapa SM-2.
func (a *App) ListMasteryItems(stage string) ([]models.UserMastery, error) {
	targetURL := a.coreURL + "/api/v1/studio/mastery/items"
	if stage != "" {
		targetURL += "?stage=" + url.QueryEscape(stage)
	}
	var list []models.UserMastery
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.UserMastery{}
	}
	return list, nil
}

// GetMyXPTransactions lista transacciones recientes del libro mayor de XP.
func (a *App) GetMyXPTransactions(limit int) ([]models.XPTransactionItem, error) {
	return a.ListXPLedger(limit)
}

// ListXPLedger obtiene historial de XP.
func (a *App) ListXPLedger(limit int) ([]models.XPTransactionItem, error) {
	targetURL := a.coreURL + "/api/v1/studio/xp-ledger"
	if limit > 0 {
		targetURL += "?limit=" + strconv.Itoa(limit)
	}
	var list []models.XPTransactionItem
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.XPTransactionItem{}
	}
	return list, nil
}

// GetMyAchievements obtiene el catálogo de logros con estado de desbloqueo.
func (a *App) GetMyAchievements() ([]models.AchievementStatus, error) {
	dash, err := a.GetMasteryDashboard()
	if err != nil {
		return nil, err
	}
	if dash.Achievements == nil {
		return []models.AchievementStatus{}, nil
	}
	return dash.Achievements, nil
}
