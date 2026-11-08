package main

import (
	"net/http"
	"net/url"
	"strconv"

	"palabrario/pkg/models"
)

// =============================================================================
// MÓDULO 4: Diccionario, Regiones Dialectales y Conceptos Semánticos (CU05, CU06)
// =============================================================================

// ListRegions lista todas las regiones panhispánicas soportadas.
func (a *App) ListRegions() ([]models.Region, error) {
	var list []models.Region
	if err := a.doRequest(a.context(), http.MethodGet, a.coreURL+"/api/v1/regions", nil, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.Region{}
	}
	return list, nil
}

// GetWordOfTheDay obtiene la palabra del día según anclaje dialectal.
func (a *App) GetWordOfTheDay(dialect string) (*models.WordDetail, error) {
	targetURL := a.coreURL + "/api/v1/words/word-of-the-day"
	if dialect != "" {
		targetURL += "?dialect=" + url.QueryEscape(dialect)
	}
	var detail models.WordDetail
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

// SearchWords ejecuta búsqueda difusa FTS de palabras.
func (a *App) SearchWords(q string, dialect string, limit int) ([]models.WordSearchResult, error) {
	params := url.Values{}
	params.Set("q", q)
	if dialect != "" {
		params.Set("dialect", dialect)
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	var list []models.WordSearchResult
	if err := a.doRequest(a.context(), http.MethodGet, a.coreURL+"/api/v1/words/search?"+params.Encode(), nil, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.WordSearchResult{}
	}
	return list, nil
}

// ExploreWords explora el catálogo léxico con filtros multifacéticos.
func (a *App) ExploreWords(filter models.WordExploreFilter) (*models.PaginatedWords, error) {
	params := url.Values{}
	if filter.Query != "" {
		params.Set("q", filter.Query)
	}
	if filter.RegionCode != "" {
		params.Set("region", filter.RegionCode)
	}
	if filter.PartOfSpeech != "" {
		params.Set("pos", filter.PartOfSpeech)
	}
	if filter.CEFRLevel != "" {
		params.Set("cefr", filter.CEFRLevel)
	}
	if filter.UsageFrequency != "" {
		params.Set("frequency", filter.UsageFrequency)
	}
	params.Set("include_taboo", strconv.FormatBool(filter.IncludeTaboo))

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	params.Set("page", strconv.Itoa(page))

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	params.Set("limit", strconv.Itoa(limit))

	var paginated models.PaginatedWords
	if err := a.doRequest(a.context(), http.MethodGet, a.coreURL+"/api/v1/words?"+params.Encode(), nil, &paginated); err != nil {
		return nil, err
	}
	return &paginated, nil
}

// GetWordDetail obtiene la ficha léxica completa de un vocablo.
func (a *App) GetWordDetail(id int) (*models.WordDetail, error) {
	return a.GetWordByID(id)
}

// GetWordByID obtiene una palabra por su ID numérico.
func (a *App) GetWordByID(id int) (*models.WordDetail, error) {
	var detail models.WordDetail
	targetURL := a.coreURL + "/api/v1/words/" + strconv.Itoa(id)
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

// CreateWordProposal registra una nueva propuesta editorial de lema.
func (a *App) CreateWordProposal(req models.CreateWordProposalRequest) (*models.WordDetail, error) {
	return a.CreateWord(req)
}

// CreateWord registra una propuesta de lema.
func (a *App) CreateWord(req models.CreateWordProposalRequest) (*models.WordDetail, error) {
	var detail models.WordDetail
	if err := a.doRequest(a.context(), http.MethodPost, a.coreURL+"/api/v1/words", req, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

// UpdateWordProposal actualiza los metadatos de un lema.
func (a *App) UpdateWordProposal(id int, req models.UpdateWordRequest) (*models.WordDetail, error) {
	return a.UpdateWord(id, req)
}

// UpdateWord actualiza metadatos de lema existente.
func (a *App) UpdateWord(id int, req models.UpdateWordRequest) (*models.WordDetail, error) {
	var detail models.WordDetail
	targetURL := a.coreURL + "/api/v1/words/" + strconv.Itoa(id)
	if err := a.doRequest(a.context(), http.MethodPut, targetURL, req, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

// CreateSenseProposal propone una nueva acepción dentro de un lema.
func (a *App) CreateSenseProposal(wordId int, req models.CreateSenseRequest) (*models.SenseDetail, error) {
	return a.CreateSense(wordId, req)
}

// CreateSense registra una nueva acepción.
func (a *App) CreateSense(wordId int, req models.CreateSenseRequest) (*models.SenseDetail, error) {
	var sense models.SenseDetail
	targetURL := a.coreURL + "/api/v1/words/" + strconv.Itoa(wordId) + "/senses"
	if err := a.doRequest(a.context(), http.MethodPost, targetURL, req, &sense); err != nil {
		return nil, err
	}
	return &sense, nil
}

// UpdateSenseProposal actualiza una acepción existente.
func (a *App) UpdateSenseProposal(senseId int, req models.UpdateSenseRequest) (*models.SenseDetail, error) {
	return a.UpdateSense(senseId, req)
}

// UpdateSense actualiza la información de una acepción.
func (a *App) UpdateSense(senseId int, req models.UpdateSenseRequest) (*models.SenseDetail, error) {
	var sense models.SenseDetail
	targetURL := a.coreURL + "/api/v1/senses/" + strconv.Itoa(senseId)
	if err := a.doRequest(a.context(), http.MethodPut, targetURL, req, &sense); err != nil {
		return nil, err
	}
	return &sense, nil
}

// DeleteSense elimina una acepción si se cuenta con privilegios editoriales.
func (a *App) DeleteSense(senseId int) error {
	return a.doRequest(a.context(), http.MethodDelete, a.coreURL+"/api/v1/senses/"+strconv.Itoa(senseId), nil, nil)
}

// AddSenseExample añade un ejemplo de uso a una acepción.
func (a *App) AddSenseExample(senseId int, req models.CreateExampleInput) (*models.Example, error) {
	return a.CreateExample(senseId, req)
}

// CreateExample registra un ejemplo en una acepción.
func (a *App) CreateExample(senseId int, req models.CreateExampleInput) (*models.Example, error) {
	var ex models.Example
	targetURL := a.coreURL + "/api/v1/senses/" + strconv.Itoa(senseId) + "/examples"
	if err := a.doRequest(a.context(), http.MethodPost, targetURL, req, &ex); err != nil {
		return nil, err
	}
	return &ex, nil
}

// UpdateExample actualiza un ejemplo de uso.
func (a *App) UpdateExample(exampleId int, req models.CreateExampleInput) (*models.Example, error) {
	var ex models.Example
	targetURL := a.coreURL + "/api/v1/examples/" + strconv.Itoa(exampleId)
	if err := a.doRequest(a.context(), http.MethodPut, targetURL, req, &ex); err != nil {
		return nil, err
	}
	return &ex, nil
}

// DeleteExample elimina un ejemplo.
func (a *App) DeleteExample(exampleId int) error {
	return a.doRequest(a.context(), http.MethodDelete, a.coreURL+"/api/v1/examples/"+strconv.Itoa(exampleId), nil, nil)
}

// ListSemanticConcepts lista conceptos semánticos transdialectales.
func (a *App) ListSemanticConcepts(q string) ([]models.SemanticConcept, error) {
	return a.ListConcepts(q)
}

// ListConcepts busca conceptos semánticos.
func (a *App) ListConcepts(q string) ([]models.SemanticConcept, error) {
	targetURL := a.coreURL + "/api/v1/concepts"
	if q != "" {
		targetURL += "?q=" + url.QueryEscape(q)
	}
	var list []models.SemanticConcept
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.SemanticConcept{}
	}
	return list, nil
}

// CreateSemanticConcept crea un nuevo concepto semántico común.
func (a *App) CreateSemanticConcept(req models.CreateConceptRequest) (*models.SemanticConcept, error) {
	return a.CreateConcept(req)
}

// CreateConcept registra un concepto semántico.
func (a *App) CreateConcept(req models.CreateConceptRequest) (*models.SemanticConcept, error) {
	var concept models.SemanticConcept
	if err := a.doRequest(a.context(), http.MethodPost, a.coreURL+"/api/v1/concepts", req, &concept); err != nil {
		return nil, err
	}
	return &concept, nil
}

// LinkSenseToConcept vincula un sentido a un concepto (heteronimia).
func (a *App) LinkSenseToConcept(senseId int, conceptId int) error {
	targetURL := a.coreURL + "/api/v1/senses/" + strconv.Itoa(senseId) + "/concepts/" + strconv.Itoa(conceptId)
	return a.doRequest(a.context(), http.MethodPost, targetURL, nil, nil)
}

// UnlinkSenseFromConcept desvincula un sentido de un concepto.
func (a *App) UnlinkSenseFromConcept(senseId int, conceptId int) error {
	targetURL := a.coreURL + "/api/v1/senses/" + strconv.Itoa(senseId) + "/concepts/" + strconv.Itoa(conceptId)
	return a.doRequest(a.context(), http.MethodDelete, targetURL, nil, nil)
}

// =============================================================================
// MÓDULO 5: Curaduría Editorial por Quórum (CU07, CU08)
// =============================================================================

// ListPendingQuorumQueue lista la cola de revisión tripartita para editores.
func (a *App) ListPendingQuorumQueue(region string) ([]models.SenseQuorumItem, error) {
	return a.ListPendingQueue(region)
}

// ListPendingQueue lista propuestas pendientes de quórum.
func (a *App) ListPendingQueue(region string) ([]models.SenseQuorumItem, error) {
	targetURL := a.coreURL + "/api/v1/curation/queue"
	if region != "" {
		targetURL += "?region=" + url.QueryEscape(region)
	}
	var list []models.SenseQuorumItem
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.SenseQuorumItem{}
	}
	return list, nil
}

// ListMyLexicalProposals lista las propuestas redactadas por el editor logueado.
func (a *App) ListMyLexicalProposals(state string) ([]models.SenseQuorumItem, error) {
	return a.ListMyProposals(state)
}

// ListMyProposals lista propuestas propias.
func (a *App) ListMyProposals(state string) ([]models.SenseQuorumItem, error) {
	targetURL := a.coreURL + "/api/v1/curation/my-proposals"
	if state != "" {
		targetURL += "?state=" + url.QueryEscape(state)
	}
	var list []models.SenseQuorumItem
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &list); err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.SenseQuorumItem{}
	}
	return list, nil
}

// GetProposalDetail obtiene el estado de votación y detalle de una propuesta.
func (a *App) GetProposalDetail(senseId int) (*models.SenseQuorumDetail, error) {
	var detail models.SenseQuorumDetail
	targetURL := a.coreURL + "/api/v1/curation/senses/" + strconv.Itoa(senseId)
	if err := a.doRequest(a.context(), http.MethodGet, targetURL, nil, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

// SubmitProposalReview emite voto de revisión tripartita sobre una propuesta.
func (a *App) SubmitProposalReview(senseId int, req models.SubmitVoteRequest) (*models.SenseQuorumDetail, error) {
	return a.SubmitVote(senseId, req)
}

// SubmitVote emite voto en el quórum.
func (a *App) SubmitVote(senseId int, req models.SubmitVoteRequest) (*models.SenseQuorumDetail, error) {
	var detail models.SenseQuorumDetail
	targetURL := a.coreURL + "/api/v1/curation/senses/" + strconv.Itoa(senseId) + "/vote"
	if err := a.doRequest(a.context(), http.MethodPost, targetURL, req, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

// AdminForceResolveProposal resolución ejecutiva administrativa directa (RN08).
func (a *App) AdminForceResolveProposal(senseId int, req models.AdminOverrideStateRequest) (*models.SenseDetail, error) {
	return a.AdminOverrideState(senseId, req)
}

// AdminOverrideState resolución administrativa directa.
func (a *App) AdminOverrideState(senseId int, req models.AdminOverrideStateRequest) (*models.SenseDetail, error) {
	var sense models.SenseDetail
	targetURL := a.coreURL + "/api/v1/curation/senses/" + strconv.Itoa(senseId) + "/state"
	if err := a.doRequest(a.context(), http.MethodPut, targetURL, req, &sense); err != nil {
		return nil, err
	}
	return &sense, nil
}
