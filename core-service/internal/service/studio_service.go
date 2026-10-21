package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"palabrario/core-service/internal/database"
	"palabrario/core-service/internal/db"
	"palabrario/pkg/models"
)

type StudioService interface {
	// Barajas (CU10)
	ListPublicDecks(ctx context.Context, query *string, officialOnly *bool) ([]models.Deck, error)
	ListStudioDecks(ctx context.Context, userID string) ([]models.StudioDeckSummary, error)
	GetDeckDetail(ctx context.Context, viewerID *string, viewerRole *models.AccountRole, deckID int) (*models.Deck, error)
	CreateDeck(ctx context.Context, userID string, userRole models.AccountRole, req models.CreateDeckRequest) (*models.Deck, error)
	UpdateDeck(ctx context.Context, userID string, userRole models.AccountRole, deckID int, req models.UpdateDeckRequest) (*models.Deck, error)
	DeleteDeck(ctx context.Context, userID string, userRole models.AccountRole, deckID int) error
	AddDeckItem(ctx context.Context, userID string, userRole models.AccountRole, deckID int, req models.AddDeckItemRequest) (*models.DeckItem, error)
	RemoveDeckItem(ctx context.Context, userID string, userRole models.AccountRole, deckID int, senseID int) error
	SaveDeck(ctx context.Context, userID string, userRole models.AccountRole, deckID int) error
	UnsaveDeck(ctx context.Context, userID string, deckID int) error

	// Colecciones (CU11)
	ListPublicCollections(ctx context.Context, featuredOnly *bool) ([]models.Collection, error)
	ListStudioCollections(ctx context.Context, userID string) ([]models.Collection, error)
	GetCollectionDetail(ctx context.Context, viewerID *string, viewerRole *models.AccountRole, collectionID int) (*models.Collection, error)
	CreateCollection(ctx context.Context, userID string, userRole models.AccountRole, req models.CreateCollectionRequest) (*models.Collection, error)
	UpdateCollection(ctx context.Context, userID string, userRole models.AccountRole, collectionID int, req models.UpdateCollectionRequest) (*models.Collection, error)
	DeleteCollection(ctx context.Context, userID string, userRole models.AccountRole, collectionID int) error
	SaveCollection(ctx context.Context, userID string, userRole models.AccountRole, collectionID int) error
	UnsaveCollection(ctx context.Context, userID string, collectionID int) error
}

type studioService struct {
	databaseInstance *database.DB
	queries          *db.Queries
}

func NewStudioService(databaseInstance *database.DB) StudioService {
	return &studioService{
		databaseInstance: databaseInstance,
		queries:          db.New(databaseInstance.Pool),
	}
}

// -----------------------------------------------------------------------------
// BARAJAS (CU10)
// -----------------------------------------------------------------------------

func (s *studioService) ListPublicDecks(ctx context.Context, query *string, officialOnly *bool) ([]models.Deck, error) {
	var qText pgtype.Text
	if query != nil && strings.TrimSpace(*query) != "" {
		qText = pgtype.Text{String: strings.TrimSpace(*query), Valid: true}
	}
	var offBool pgtype.Bool
	if officialOnly != nil {
		offBool = pgtype.Bool{Bool: *officialOnly, Valid: true}
	}

	rows, err := s.queries.ListPublicDecks(ctx, db.ListPublicDecksParams{
		Query:        qText,
		OfficialOnly: offBool,
	})
	if err != nil {
		return nil, fmt.Errorf("error listando barajas públicas: %w", err)
	}

	decks := make([]models.Deck, len(rows))
	for i, r := range rows {
		decks[i] = models.Deck{
			IDDeck:      int(r.IDDeck),
			IDUser:      r.IDUser.String(),
			Title:       r.Title,
			Description: TextToPtr(r.Description),
			IsPublic:    r.IsPublic,
			IsOfficial:  r.IsOfficial,
			ItemCount:   int(r.ItemCount),
			CreatedAt:   r.CreatedAt.Time,
			UpdatedAt:   r.UpdatedAt.Time,
		}
	}
	return decks, nil
}

func (s *studioService) ListStudioDecks(ctx context.Context, userID string) ([]models.StudioDeckSummary, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, NewDomainError(ErrInvalidInput, "Identificador de usuario inválido")
	}

	rows, err := s.queries.ListStudioDecks(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("error listando barajas de Studio: %w", err)
	}

	summaries := make([]models.StudioDeckSummary, len(rows))
	for i, r := range rows {
		summaries[i] = models.StudioDeckSummary{
			Deck: models.Deck{
				IDDeck:      int(r.IDDeck),
				IDUser:      r.IDUser.String(),
				Title:       r.Title,
				Description: TextToPtr(r.Description),
				IsPublic:    r.IsPublic,
				IsOfficial:  r.IsOfficial,
				ItemCount:   int(r.ItemCount),
				CreatedAt:   r.CreatedAt.Time,
				UpdatedAt:   r.UpdatedAt.Time,
			},
			OwnerUserName:    r.OwnerUserName,
			IsSavedFromComm:  r.IsSavedFromCommunity,
			DueCardsToday:    int(r.DueCardsToday),
			RetentionPercent: r.RetentionPercent,
		}
	}
	return summaries, nil
}

func (s *studioService) GetDeckDetail(ctx context.Context, viewerID *string, viewerRole *models.AccountRole, deckID int) (*models.Deck, error) {
	if deckID <= 0 {
		return nil, NewDomainError(ErrInvalidInput, "El identificador de la baraja debe ser positivo")
	}

	deckRow, err := s.queries.GetDeckByID(ctx, int32(deckID))
	if err != nil {
		return nil, MapDBError(err, "Baraja no encontrada")
	}

	// RN04: Privacidad de recursos. Si la baraja es privada, solo el dueño o un ADMIN pueden verla
	if !deckRow.IsPublic {
		isAllowed := false
		if viewerID != nil && *viewerID == deckRow.IDUser.String() {
			isAllowed = true
		} else if viewerRole != nil && *viewerRole == models.RoleAdmin {
			isAllowed = true
		}
		if !isAllowed {
			return nil, NewDomainError(ErrForbidden, "Violación de regla RN04: No tienes acceso a esta baraja privada")
		}
	}

	// Obtener los sentidos asociados
	itemsRows, err := s.queries.ListDeckItemsByDeckID(ctx, int32(deckID))
	if err != nil {
		return nil, fmt.Errorf("error obteniendo tarjetas de la baraja: %w", err)
	}

	// Cargar ejemplos para las acepciones
	senseIDs := make([]int32, len(itemsRows))
	for i, r := range itemsRows {
		senseIDs[i] = r.IDSense
	}
	examplesRows, _ := s.queries.GetExamplesForSenses(ctx, senseIDs)
	examplesMap := make(map[int32][]models.Example)
	for _, ex := range examplesRows {
		examplesMap[ex.IDSense] = append(examplesMap[ex.IDSense], models.Example{
			IDExample: int(ex.IDExample),
			IDSense:   int(ex.IDSense),
			Sentence:  ex.Sentence,
			Source:    TextToPtr(ex.Source),
		})
	}

	deckItems := make([]models.DeckItem, len(itemsRows))
	for i, r := range itemsRows {
		countryCode := ""
		if r.CountryCode.Valid {
			countryCode = r.CountryCode.String
		}

		cefrStr := "B1"
		if r.CefrLevel != nil {
			cefrStr = fmt.Sprint(r.CefrLevel)
		}
		freqStr := "COMMON"
		if r.UsageFrequency != nil {
			freqStr = fmt.Sprint(r.UsageFrequency)
		}

		deckItems[i] = models.DeckItem{
			IDDeck:    int(r.IDDeck),
			IDSense:   int(r.IDSense),
			SortOrder: int(r.SortOrder),
			Sense: &models.SenseDetail{
				IDSense:         int(r.IDSense),
				IDWord:          int(r.IDWord),
				RegionCode:      r.RegionCode,
				RegionName:      r.RegionName,
				CountryCode:     countryCode,
				FlagEmoji:       r.FlagEmoji,
				ShortDefinition: r.ShortDefinition,
				Definition:      r.Definition,
				IsPrimary:       r.IsPrimary,
				IsTabooOrVulgar: r.IsTabooOrVulgar,
				CulturalNote:    TextToPtr(r.CulturalNote),
				CEFRLevel:       cefrStr,
				UsageFrequency:  freqStr,
				State:           r.State,
				Examples:        examplesMap[r.IDSense],
			},
		}
	}

	return &models.Deck{
		IDDeck:      int(deckRow.IDDeck),
		IDUser:      deckRow.IDUser.String(),
		Title:       deckRow.Title,
		Description: TextToPtr(deckRow.Description),
		IsPublic:    deckRow.IsPublic,
		IsOfficial:  deckRow.IsOfficial,
		ItemCount:   len(deckItems),
		Items:       deckItems,
		CreatedAt:   deckRow.CreatedAt.Time,
		UpdatedAt:   deckRow.UpdatedAt.Time,
	}, nil
}

func (s *studioService) CreateDeck(ctx context.Context, userID string, userRole models.AccountRole, req models.CreateDeckRequest) (*models.Deck, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, NewDomainError(ErrInvalidInput, "Identificador de usuario inválido")
	}

	title := strings.TrimSpace(req.Title)
	if len(title) < 3 || len(title) > 150 {
		return nil, NewDomainError(ErrInvalidInput, "El título de la baraja debe tener entre 3 y 150 caracteres")
	}

	// RN05: is_official solo es permitido para EDITOR o ADMIN; de lo contrario se fuerza a false
	isOfficial := req.IsOfficial
	if userRole != models.RoleEditor && userRole != models.RoleAdmin {
		isOfficial = false
	}

	created, err := s.queries.CreateDeck(ctx, db.CreateDeckParams{
		IDUser:      userUUID,
		Title:       title,
		Description: stringToText(req.Description),
		IsPublic:    req.IsPublic,
		IsOfficial:  isOfficial,
	})
	if err != nil {
		return nil, fmt.Errorf("error al crear baraja: %w", err)
	}

	// Insertar tarjetas iniciales si fueron enviadas
	if len(req.SenseIDs) > 0 {
		for idx, senseID := range req.SenseIDs {
			_, _ = s.queries.AddDeckItem(ctx, db.AddDeckItemParams{
				IDDeck:    created.IDDeck,
				IDSense:   int32(senseID),
				SortOrder: int32(idx + 1),
			})
		}
	}

	return &models.Deck{
		IDDeck:      int(created.IDDeck),
		IDUser:      created.IDUser.String(),
		Title:       created.Title,
		Description: TextToPtr(created.Description),
		IsPublic:    created.IsPublic,
		IsOfficial:  created.IsOfficial,
		ItemCount:   len(req.SenseIDs),
		CreatedAt:   created.CreatedAt.Time,
		UpdatedAt:   created.UpdatedAt.Time,
	}, nil
}

func (s *studioService) UpdateDeck(ctx context.Context, userID string, userRole models.AccountRole, deckID int, req models.UpdateDeckRequest) (*models.Deck, error) {
	if deckID <= 0 {
		return nil, NewDomainError(ErrInvalidInput, "Identificador de baraja inválido")
	}

	deckRow, err := s.queries.GetDeckByID(ctx, int32(deckID))
	if err != nil {
		return nil, MapDBError(err, "Baraja no encontrada")
	}

	// RN03: Solo el propietario o un ADMIN pueden editar la baraja
	if deckRow.IDUser.String() != userID && userRole != models.RoleAdmin {
		return nil, NewDomainError(ErrForbidden, "Violación de regla RN03: Solo el propietario o un Administrador pueden modificar esta baraja")
	}

	title := strings.TrimSpace(req.Title)
	if len(title) < 3 || len(title) > 150 {
		return nil, NewDomainError(ErrInvalidInput, "El título de la baraja debe tener entre 3 y 150 caracteres")
	}

	// RN05: is_official exclusivo para EDITOR o ADMIN
	isOfficial := req.IsOfficial
	if userRole != models.RoleEditor && userRole != models.RoleAdmin {
		isOfficial = false
	}

	updated, err := s.queries.UpdateDeck(ctx, db.UpdateDeckParams{
		IDDeck:      int32(deckID),
		Title:       title,
		Description: stringToText(req.Description),
		IsPublic:    req.IsPublic,
		IsOfficial:  isOfficial,
	})
	if err != nil {
		return nil, fmt.Errorf("error al actualizar baraja: %w", err)
	}

	return &models.Deck{
		IDDeck:      int(updated.IDDeck),
		IDUser:      updated.IDUser.String(),
		Title:       updated.Title,
		Description: TextToPtr(updated.Description),
		IsPublic:    updated.IsPublic,
		IsOfficial:  updated.IsOfficial,
		ItemCount:   int(deckRow.ItemCount),
		CreatedAt:   updated.CreatedAt.Time,
		UpdatedAt:   updated.UpdatedAt.Time,
	}, nil
}

func (s *studioService) DeleteDeck(ctx context.Context, userID string, userRole models.AccountRole, deckID int) error {
	if deckID <= 0 {
		return NewDomainError(ErrInvalidInput, "Identificador de baraja inválido")
	}

	deckRow, err := s.queries.GetDeckByID(ctx, int32(deckID))
	if err != nil {
		return MapDBError(err, "Baraja no encontrada")
	}

	// RN03: Solo el dueño o un ADMIN pueden eliminarla
	if deckRow.IDUser.String() != userID && userRole != models.RoleAdmin {
		return NewDomainError(ErrForbidden, "Violación de regla RN03: Solo el propietario o un Administrador pueden eliminar esta baraja")
	}

	return s.queries.DeleteDeck(ctx, int32(deckID))
}

func (s *studioService) AddDeckItem(ctx context.Context, userID string, userRole models.AccountRole, deckID int, req models.AddDeckItemRequest) (*models.DeckItem, error) {
	if deckID <= 0 || req.IDSense <= 0 {
		return nil, NewDomainError(ErrInvalidInput, "Identificadores de baraja y sentido requeridos")
	}

	deckRow, err := s.queries.GetDeckByID(ctx, int32(deckID))
	if err != nil {
		return nil, MapDBError(err, "Baraja no encontrada")
	}

	// RN03: Solo el dueño o un ADMIN pueden agregar ítems
	if deckRow.IDUser.String() != userID && userRole != models.RoleAdmin {
		return nil, NewDomainError(ErrForbidden, "Violación de regla RN03: Solo el propietario o un Administrador pueden modificar esta baraja")
	}

	sortOrder := int32(0)
	if req.SortOrder != nil {
		sortOrder = int32(*req.SortOrder)
	} else {
		sortOrder = deckRow.ItemCount + 1
	}

	item, err := s.queries.AddDeckItem(ctx, db.AddDeckItemParams{
		IDDeck:    int32(deckID),
		IDSense:   int32(req.IDSense),
		SortOrder: sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("error al agregar tarjeta a la baraja: %w", err)
	}

	return &models.DeckItem{
		IDDeck:    int(item.IDDeck),
		IDSense:   int(item.IDSense),
		SortOrder: int(item.SortOrder),
	}, nil
}

func (s *studioService) RemoveDeckItem(ctx context.Context, userID string, userRole models.AccountRole, deckID int, senseID int) error {
	if deckID <= 0 || senseID <= 0 {
		return NewDomainError(ErrInvalidInput, "Identificadores de baraja y tarjeta requeridos")
	}

	deckRow, err := s.queries.GetDeckByID(ctx, int32(deckID))
	if err != nil {
		return MapDBError(err, "Baraja no encontrada")
	}

	// RN03: Solo el dueño o un ADMIN pueden eliminar ítems
	if deckRow.IDUser.String() != userID && userRole != models.RoleAdmin {
		return NewDomainError(ErrForbidden, "Violación de regla RN03: Solo el propietario o un Administrador pueden modificar esta baraja")
	}

	return s.queries.RemoveDeckItem(ctx, db.RemoveDeckItemParams{
		IDDeck:  int32(deckID),
		IDSense: int32(senseID),
	})
}

func (s *studioService) SaveDeck(ctx context.Context, userID string, userRole models.AccountRole, deckID int) error {
	if deckID <= 0 {
		return NewDomainError(ErrInvalidInput, "Identificador de baraja inválido")
	}

	deckRow, err := s.queries.GetDeckByID(ctx, int32(deckID))
	if err != nil {
		return MapDBError(err, "Baraja no encontrada")
	}

	// RN04: No se puede guardar una baraja privada ajena
	if !deckRow.IsPublic && deckRow.IDUser.String() != userID && userRole != models.RoleAdmin {
		return NewDomainError(ErrForbidden, "Violación de regla RN04: No puedes guardar una baraja privada ajena")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return NewDomainError(ErrInvalidInput, "Identificador de usuario inválido")
	}

	return s.queries.SaveDeck(ctx, db.SaveDeckParams{
		IDUser: userUUID,
		IDDeck: int32(deckID),
	})
}

func (s *studioService) UnsaveDeck(ctx context.Context, userID string, deckID int) error {
	if deckID <= 0 {
		return NewDomainError(ErrInvalidInput, "Identificador de baraja inválido")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return NewDomainError(ErrInvalidInput, "Identificador de usuario inválido")
	}

	return s.queries.UnsaveDeck(ctx, db.UnsaveDeckParams{
		IDUser: userUUID,
		IDDeck: int32(deckID),
	})
}

// -----------------------------------------------------------------------------
// COLECCIONES TEMÁTICAS (CU11)
// -----------------------------------------------------------------------------

func (s *studioService) ListPublicCollections(ctx context.Context, featuredOnly *bool) ([]models.Collection, error) {
	var featBool pgtype.Bool
	if featuredOnly != nil {
		featBool = pgtype.Bool{Bool: *featuredOnly, Valid: true}
	}

	rows, err := s.queries.ListPublicCollections(ctx, featBool)
	if err != nil {
		return nil, fmt.Errorf("error listando colecciones públicas: %w", err)
	}

	cols := make([]models.Collection, len(rows))
	for i, r := range rows {
		cols[i] = models.Collection{
			IDCollection: int(r.IDCollection),
			IDUser:       r.IDUser.String(),
			Title:        r.Title,
			Description:  TextToPtr(r.Description),
			IsPublic:     r.IsPublic,
			IsFeatured:   r.IsFeatured,
			IsOfficial:   r.IsOfficial,
			CreatedAt:    r.CreatedAt.Time,
			UpdatedAt:    r.UpdatedAt.Time,
		}
	}
	return cols, nil
}

func (s *studioService) ListStudioCollections(ctx context.Context, userID string) ([]models.Collection, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, NewDomainError(ErrInvalidInput, "Identificador de usuario inválido")
	}

	rows, err := s.queries.ListStudioCollections(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("error listando colecciones de Studio: %w", err)
	}

	cols := make([]models.Collection, len(rows))
	for i, r := range rows {
		cols[i] = models.Collection{
			IDCollection: int(r.IDCollection),
			IDUser:       r.IDUser.String(),
			Title:        r.Title,
			Description:  TextToPtr(r.Description),
			IsPublic:     r.IsPublic,
			IsFeatured:   r.IsFeatured,
			IsOfficial:   r.IsOfficial,
			CreatedAt:    r.CreatedAt.Time,
			UpdatedAt:    r.UpdatedAt.Time,
		}
	}
	return cols, nil
}

func (s *studioService) GetCollectionDetail(ctx context.Context, viewerID *string, viewerRole *models.AccountRole, collectionID int) (*models.Collection, error) {
	if collectionID <= 0 {
		return nil, NewDomainError(ErrInvalidInput, "Identificador de colección inválido")
	}

	colRow, err := s.queries.GetCollectionByID(ctx, int32(collectionID))
	if err != nil {
		return nil, MapDBError(err, "Colección no encontrada")
	}

	// RN04: Privacidad de recursos. Si la colección es privada, solo dueño o ADMIN pueden verla
	if !colRow.IsPublic {
		isAllowed := false
		if viewerID != nil && *viewerID == colRow.IDUser.String() {
			isAllowed = true
		} else if viewerRole != nil && *viewerRole == models.RoleAdmin {
			isAllowed = true
		}
		if !isAllowed {
			return nil, NewDomainError(ErrForbidden, "Violación de regla RN04: No tienes acceso a esta colección privada")
		}
	}

	// Cargar barajas vinculadas
	deckRows, err := s.queries.ListCollectionDecks(ctx, int32(collectionID))
	if err != nil {
		return nil, fmt.Errorf("error obteniendo barajas de la colección: %w", err)
	}

	decks := make([]models.Deck, len(deckRows))
	for i, r := range deckRows {
		decks[i] = models.Deck{
			IDDeck:      int(r.IDDeck),
			IDUser:      r.IDUser.String(),
			Title:       r.Title,
			Description: TextToPtr(r.Description),
			IsPublic:    r.IsPublic,
			IsOfficial:  r.IsOfficial,
			ItemCount:   int(r.ItemCount),
			CreatedAt:   r.CreatedAt.Time,
			UpdatedAt:   r.UpdatedAt.Time,
		}
	}

	return &models.Collection{
		IDCollection: int(colRow.IDCollection),
		IDUser:       colRow.IDUser.String(),
		Title:        colRow.Title,
		Description:  TextToPtr(colRow.Description),
		IsPublic:     colRow.IsPublic,
		IsFeatured:   colRow.IsFeatured,
		IsOfficial:   colRow.IsOfficial,
		Decks:        decks,
		CreatedAt:    colRow.CreatedAt.Time,
		UpdatedAt:    colRow.UpdatedAt.Time,
	}, nil
}

func (s *studioService) CreateCollection(ctx context.Context, userID string, userRole models.AccountRole, req models.CreateCollectionRequest) (*models.Collection, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, NewDomainError(ErrInvalidInput, "Identificador de usuario inválido")
	}

	title := strings.TrimSpace(req.Title)
	if len(title) < 3 || len(title) > 150 {
		return nil, NewDomainError(ErrInvalidInput, "El título de la colección debe tener entre 3 y 150 caracteres")
	}

	// RN05: is_official e is_featured exclusivos para EDITOR o ADMIN
	isOfficial := req.IsOfficial
	isFeatured := req.IsFeatured
	if userRole != models.RoleEditor && userRole != models.RoleAdmin {
		isOfficial = false
		isFeatured = false
	}

	created, err := s.queries.CreateCollection(ctx, db.CreateCollectionParams{
		IDUser:      userUUID,
		Title:       title,
		Description: stringToText(req.Description),
		IsPublic:    req.IsPublic,
		IsFeatured:  isFeatured,
		IsOfficial:  isOfficial,
	})
	if err != nil {
		return nil, fmt.Errorf("error al crear colección: %w", err)
	}

	// Vincular barajas
	if len(req.DeckIDs) > 0 {
		for idx, deckID := range req.DeckIDs {
			_ = s.queries.AddCollectionDeck(ctx, db.AddCollectionDeckParams{
				IDCollection: created.IDCollection,
				IDDeck:       int32(deckID),
				SortOrder:    int32(idx + 1),
			})
		}
	}

	return &models.Collection{
		IDCollection: int(created.IDCollection),
		IDUser:       created.IDUser.String(),
		Title:        created.Title,
		Description:  TextToPtr(created.Description),
		IsPublic:     created.IsPublic,
		IsFeatured:   created.IsFeatured,
		IsOfficial:   created.IsOfficial,
		CreatedAt:    created.CreatedAt.Time,
		UpdatedAt:    created.UpdatedAt.Time,
	}, nil
}

func (s *studioService) UpdateCollection(ctx context.Context, userID string, userRole models.AccountRole, collectionID int, req models.UpdateCollectionRequest) (*models.Collection, error) {
	if collectionID <= 0 {
		return nil, NewDomainError(ErrInvalidInput, "Identificador de colección inválido")
	}

	colRow, err := s.queries.GetCollectionByID(ctx, int32(collectionID))
	if err != nil {
		return nil, MapDBError(err, "Colección no encontrada")
	}

	// RN03: Solo dueño o ADMIN
	if colRow.IDUser.String() != userID && userRole != models.RoleAdmin {
		return nil, NewDomainError(ErrForbidden, "Violación de regla RN03: Solo el propietario o un Administrador pueden modificar esta colección")
	}

	title := strings.TrimSpace(req.Title)
	if len(title) < 3 || len(title) > 150 {
		return nil, NewDomainError(ErrInvalidInput, "El título de la colección debe tener entre 3 y 150 caracteres")
	}

	// RN05: is_official e is_featured exclusivos para EDITOR o ADMIN
	isOfficial := req.IsOfficial
	isFeatured := req.IsFeatured
	if userRole != models.RoleEditor && userRole != models.RoleAdmin {
		isOfficial = false
		isFeatured = false
	}

	updated, err := s.queries.UpdateCollection(ctx, db.UpdateCollectionParams{
		IDCollection: int32(collectionID),
		Title:        title,
		Description:  stringToText(req.Description),
		IsPublic:     req.IsPublic,
		IsFeatured:   isFeatured,
		IsOfficial:   isOfficial,
	})
	if err != nil {
		return nil, fmt.Errorf("error al actualizar colección: %w", err)
	}

	if req.DeckIDs != nil {
		_ = s.queries.ClearCollectionDecks(ctx, int32(collectionID))
		for idx, deckID := range req.DeckIDs {
			_ = s.queries.AddCollectionDeck(ctx, db.AddCollectionDeckParams{
				IDCollection: int32(collectionID),
				IDDeck:       int32(deckID),
				SortOrder:    int32(idx + 1),
			})
		}
	}

	return &models.Collection{
		IDCollection: int(updated.IDCollection),
		IDUser:       updated.IDUser.String(),
		Title:        updated.Title,
		Description:  TextToPtr(updated.Description),
		IsPublic:     updated.IsPublic,
		IsFeatured:   updated.IsFeatured,
		IsOfficial:   updated.IsOfficial,
		CreatedAt:    updated.CreatedAt.Time,
		UpdatedAt:    updated.UpdatedAt.Time,
	}, nil
}

func (s *studioService) DeleteCollection(ctx context.Context, userID string, userRole models.AccountRole, collectionID int) error {
	if collectionID <= 0 {
		return NewDomainError(ErrInvalidInput, "Identificador de colección inválido")
	}

	colRow, err := s.queries.GetCollectionByID(ctx, int32(collectionID))
	if err != nil {
		return MapDBError(err, "Colección no encontrada")
	}

	// RN03: Solo dueño o ADMIN
	if colRow.IDUser.String() != userID && userRole != models.RoleAdmin {
		return NewDomainError(ErrForbidden, "Violación de regla RN03: Solo el propietario o un Administrador pueden eliminar esta colección")
	}

	return s.queries.DeleteCollection(ctx, int32(collectionID))
}

func (s *studioService) SaveCollection(ctx context.Context, userID string, userRole models.AccountRole, collectionID int) error {
	if collectionID <= 0 {
		return NewDomainError(ErrInvalidInput, "Identificador de colección inválido")
	}

	colRow, err := s.queries.GetCollectionByID(ctx, int32(collectionID))
	if err != nil {
		return MapDBError(err, "Colección no encontrada")
	}

	// RN04: No se puede guardar una colección privada ajena
	if !colRow.IsPublic && colRow.IDUser.String() != userID && userRole != models.RoleAdmin {
		return NewDomainError(ErrForbidden, "Violación de regla RN04: No puedes guardar una colección privada ajena")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return NewDomainError(ErrInvalidInput, "Identificador de usuario inválido")
	}

	return s.queries.SaveCollection(ctx, db.SaveCollectionParams{
		IDUser:       userUUID,
		IDCollection: int32(collectionID),
	})
}

func (s *studioService) UnsaveCollection(ctx context.Context, userID string, collectionID int) error {
	if collectionID <= 0 {
		return NewDomainError(ErrInvalidInput, "Identificador de colección inválido")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return NewDomainError(ErrInvalidInput, "Identificador de usuario inválido")
	}

	return s.queries.UnsaveCollection(ctx, db.UnsaveCollectionParams{
		IDUser:       userUUID,
		IDCollection: int32(collectionID),
	})
}

func stringToText(s *string) pgtype.Text {
	if s == nil || strings.TrimSpace(*s) == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: strings.TrimSpace(*s), Valid: true}
}
