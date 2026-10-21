-- ==============================================================================
-- STUDIO PERSONAL: BARAJAS Y COLECCIONES (CU10, CU11)
-- ==============================================================================

-- 1. Listar Barajas Públicas
-- name: ListPublicDecks :many
SELECT 
    d.id_deck,
    d.id_user,
    d.title,
    d.description,
    d.is_public,
    d.is_official,
    d.created_at,
    d.updated_at,
    (SELECT count(*)::INT FROM deck_items di WHERE di.id_deck = d.id_deck) AS item_count
FROM decks d
WHERE d.is_public = true
  AND (sqlc.narg('query')::VARCHAR IS NULL OR sqlc.narg('query')::VARCHAR = '' OR d.title ILIKE '%' || sqlc.narg('query')::VARCHAR || '%' OR coalesce(d.description, '') ILIKE '%' || sqlc.narg('query')::VARCHAR || '%')
  AND (sqlc.narg('official_only')::BOOLEAN IS NULL OR sqlc.narg('official_only')::BOOLEAN = false OR d.is_official = true)
ORDER BY d.created_at DESC;

-- 2. Listar Barajas del Studio Personal (Propias y Guardadas de la Comunidad)
-- name: ListStudioDecks :many
SELECT 
    d.id_deck,
    d.id_user,
    d.title,
    d.description,
    d.is_public,
    d.is_official,
    d.created_at,
    d.updated_at,
    u.user_name AS owner_user_name,
    (d.id_user != @user_id)::BOOLEAN AS is_saved_from_community,
    (SELECT count(*)::INT FROM deck_items di WHERE di.id_deck = d.id_deck) AS item_count,
    (
        SELECT count(*)::INT 
        FROM deck_items di
        LEFT JOIN user_mastery um ON um.id_sense = di.id_sense AND um.id_user = @user_id
        WHERE di.id_deck = d.id_deck
          AND (um.id_user IS NULL OR um.next_review_at <= now())
    ) AS due_cards_today,
    COALESCE((
        SELECT AVG(um.mastery_percentage)::REAL
        FROM deck_items di
        JOIN user_mastery um ON um.id_sense = di.id_sense AND um.id_user = @user_id
        WHERE di.id_deck = d.id_deck
    ), 0.0)::REAL AS retention_percent
FROM decks d
JOIN users u ON u.id_user = d.id_user
WHERE d.id_user = @user_id 
   OR (d.is_public = true AND EXISTS (SELECT 1 FROM saved_decks sd WHERE sd.id_deck = d.id_deck AND sd.id_user = @user_id))
ORDER BY d.updated_at DESC;

-- 3. Detalle de Baraja
-- name: GetDeckByID :one
SELECT 
    d.id_deck,
    d.id_user,
    d.title,
    d.description,
    d.is_public,
    d.is_official,
    d.created_at,
    d.updated_at,
    (SELECT count(*)::INT FROM deck_items di WHERE di.id_deck = d.id_deck) AS item_count
FROM decks d
WHERE d.id_deck = $1;

-- 4. Listar Elementos (Sentidos) de una Baraja
-- name: ListDeckItemsByDeckID :many
SELECT 
    di.id_deck,
    di.id_sense,
    di.sort_order,
    s.id_word,
    w.lemma,
    w.part_of_speech,
    s.id_region,
    r.code AS region_code,
    r.name AS region_name,
    r.parent_zone,
    r.country_code,
    r.flag_emoji,
    s.short_definition,
    s.definition,
    s.is_primary,
    s.is_taboo_or_vulgar,
    s.cultural_note,
    coalesce(s.cefr_level::TEXT, 'B1') AS cefr_level,
    coalesce(s.usage_frequency::TEXT, 'COMMON') AS usage_frequency,
    s.state::TEXT AS state,
    s.id_editor,
    u.user_name AS editor_user_name,
    s.created_at,
    s.updated_at
FROM deck_items di
JOIN senses s ON s.id_sense = di.id_sense
JOIN words w ON w.id_word = s.id_word
JOIN regions r ON r.id_region = s.id_region
JOIN users u ON u.id_user = s.id_editor
WHERE di.id_deck = $1
ORDER BY di.sort_order ASC, di.id_sense ASC;

-- 5. Crear Baraja
-- name: CreateDeck :one
INSERT INTO decks (
    id_user,
    title,
    description,
    is_public,
    is_official,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, now(), now()
) RETURNING id_deck, id_user, title, description, is_public, is_official, created_at, updated_at;

-- 6. Actualizar Baraja
-- name: UpdateDeck :one
UPDATE decks
SET 
    title = $2,
    description = $3,
    is_public = $4,
    is_official = $5,
    updated_at = now()
WHERE id_deck = $1
RETURNING id_deck, id_user, title, description, is_public, is_official, created_at, updated_at;

-- 7. Eliminar Baraja
-- name: DeleteDeck :exec
DELETE FROM decks WHERE id_deck = $1;

-- 8. Añadir Tarjeta (Sentido) a Baraja
-- name: AddDeckItem :one
INSERT INTO deck_items (
    id_deck,
    id_sense,
    sort_order
) VALUES (
    $1, $2, $3
) ON CONFLICT (id_deck, id_sense) DO UPDATE
SET sort_order = EXCLUDED.sort_order
RETURNING id_deck, id_sense, sort_order;

-- 9. Eliminar Tarjeta de Baraja
-- name: RemoveDeckItem :exec
DELETE FROM deck_items WHERE id_deck = $1 AND id_sense = $2;

-- 10. Guardar Baraja Pública en Studio Personal
-- name: SaveDeck :exec
INSERT INTO saved_decks (id_user, id_deck, saved_at)
VALUES ($1, $2, now())
ON CONFLICT (id_user, id_deck) DO NOTHING;

-- 11. Dejar de Guardar Baraja
-- name: UnsaveDeck :exec
DELETE FROM saved_decks WHERE id_user = $1 AND id_deck = $2;

-- 12. Verificar si una Baraja está Guardada
-- name: IsDeckSavedByUser :one
SELECT EXISTS(
    SELECT 1 FROM saved_decks WHERE id_user = $1 AND id_deck = $2
);

-- ==============================================================================
-- COLECCIONES TEMÁTICAS (CU11)
-- ==============================================================================

-- 13. Listar Colecciones Públicas
-- name: ListPublicCollections :many
SELECT 
    c.id_collection,
    c.id_user,
    c.title,
    c.description,
    c.is_public,
    c.is_featured,
    c.is_official,
    c.created_at,
    c.updated_at
FROM collections c
WHERE c.is_public = true
  AND (sqlc.narg('featured_only')::BOOLEAN IS NULL OR sqlc.narg('featured_only')::BOOLEAN = false OR c.is_featured = true)
ORDER BY c.created_at DESC;

-- 14. Listar Colecciones del Studio Personal
-- name: ListStudioCollections :many
SELECT 
    c.id_collection,
    c.id_user,
    c.title,
    c.description,
    c.is_public,
    c.is_featured,
    c.is_official,
    c.created_at,
    c.updated_at,
    (c.id_user != @user_id)::BOOLEAN AS is_saved_from_community
FROM collections c
WHERE c.id_user = @user_id 
   OR (c.is_public = true AND EXISTS (SELECT 1 FROM saved_collections sc WHERE sc.id_collection = c.id_collection AND sc.id_user = @user_id))
ORDER BY c.updated_at DESC;

-- 15. Detalle de Colección
-- name: GetCollectionByID :one
SELECT 
    c.id_collection,
    c.id_user,
    c.title,
    c.description,
    c.is_public,
    c.is_featured,
    c.is_official,
    c.created_at,
    c.updated_at
FROM collections c
WHERE c.id_collection = $1;

-- 16. Listar Barajas de una Colección
-- name: ListCollectionDecks :many
SELECT 
    cd.id_collection,
    cd.id_deck,
    cd.sort_order,
    d.id_user,
    d.title,
    d.description,
    d.is_public,
    d.is_official,
    d.created_at,
    d.updated_at,
    (SELECT count(*)::INT FROM deck_items di WHERE di.id_deck = d.id_deck) AS item_count
FROM collection_decks cd
JOIN decks d ON d.id_deck = cd.id_deck
WHERE cd.id_collection = $1
ORDER BY cd.sort_order ASC, cd.id_deck ASC;

-- 17. Crear Colección
-- name: CreateCollection :one
INSERT INTO collections (
    id_user,
    title,
    description,
    is_public,
    is_featured,
    is_official,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, now(), now()
) RETURNING id_collection, id_user, title, description, is_public, is_featured, is_official, created_at, updated_at;

-- 18. Actualizar Colección
-- name: UpdateCollection :one
UPDATE collections
SET 
    title = $2,
    description = $3,
    is_public = $4,
    is_featured = $5,
    is_official = $6,
    updated_at = now()
WHERE id_collection = $1
RETURNING id_collection, id_user, title, description, is_public, is_featured, is_official, created_at, updated_at;

-- 19. Eliminar Colección
-- name: DeleteCollection :exec
DELETE FROM collections WHERE id_collection = $1;

-- 20. Añadir Baraja a Colección
-- name: AddCollectionDeck :exec
INSERT INTO collection_decks (
    id_collection,
    id_deck,
    sort_order
) VALUES (
    $1, $2, $3
) ON CONFLICT (id_collection, id_deck) DO UPDATE
SET sort_order = EXCLUDED.sort_order;

-- 21. Limpiar Barajas de Colección
-- name: ClearCollectionDecks :exec
DELETE FROM collection_decks WHERE id_collection = $1;

-- 22. Guardar Colección Pública en Studio Personal
-- name: SaveCollection :exec
INSERT INTO saved_collections (id_user, id_collection, saved_at)
VALUES ($1, $2, now())
ON CONFLICT (id_user, id_collection) DO NOTHING;

-- 23. Dejar de Guardar Colección
-- name: UnsaveCollection :exec
DELETE FROM saved_collections WHERE id_user = $1 AND id_collection = $2;

-- 24. Verificar si una Colección está Guardada
-- name: IsCollectionSavedByUser :one
SELECT EXISTS(
    SELECT 1 FROM saved_collections WHERE id_user = $1 AND id_collection = $2
);
