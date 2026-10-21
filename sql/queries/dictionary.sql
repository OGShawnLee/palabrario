-- ==============================================================================
-- DICCIONARIO Y BÚSQUEDA DIALECTAL (sqlc)
-- ==============================================================================

-- 1. Búsqueda Rápida para Dropdown / Autocompletado
-- name: SearchVocabularyDropdown :many
WITH matched_words AS (
  SELECT
    w.id_word,
    w.lemma,
    w.normalized_lemma,
    similarity(w.normalized_lemma, unaccent(lower(sqlc.arg('search_query')::TEXT))) AS sim
  FROM words w
  WHERE
    sqlc.arg('search_query')::TEXT IS NOT NULL 
    AND length(trim(sqlc.arg('search_query')::TEXT)) > 0 
    AND (
      w.normalized_lemma ILIKE unaccent(lower(sqlc.arg('search_query')::TEXT)) || '%'
      OR w.search_vector @@ plainto_tsquery('spanish_unaccent', sqlc.arg('search_query')::TEXT)
      OR similarity(w.normalized_lemma, unaccent(lower(sqlc.arg('search_query')::TEXT))) > 0.25
    )
  ORDER BY
    w.normalized_lemma LIKE unaccent(lower(sqlc.arg('search_query')::TEXT)) || '%' DESC,
    sim DESC,
    w.lemma ASC
  LIMIT 15
),
distinct_words AS (
  SELECT DISTINCT ON (mw.id_word)
    mw.id_word,
    mw.lemma,
    mw.normalized_lemma,
    coalesce(s.short_definition, 'Sin definición corta') AS short_definition,
    coalesce(r.flag_emoji, '🌐') AS flag_emoji,
    coalesce(r.country_code, '') AS country_code,
    mw.sim AS similarity
  FROM matched_words mw
  LEFT JOIN senses s ON s.id_word = mw.id_word AND s.state = 'APPROVED'
  LEFT JOIN regions r ON r.id_region = s.id_region
  ORDER BY
    mw.id_word,
    (r.code = sqlc.arg('user_dialect')::VARCHAR) DESC NULLS LAST,
    s.is_primary DESC
)
SELECT
  dw.id_word,
  dw.lemma,
  dw.short_definition,
  dw.flag_emoji,
  dw.country_code,
  dw.similarity
FROM distinct_words dw
ORDER BY
  dw.normalized_lemma LIKE unaccent(lower(sqlc.arg('search_query')::TEXT)) || '%' DESC,
  dw.similarity DESC,
  dw.lemma ASC
LIMIT sqlc.arg('max_results')::INT;

-- 2. Detalle de Palabra Canónica con su Región de Origen
-- name: GetWordDetailByID :one
SELECT 
    w.id_word,
    w.lemma,
    w.normalized_lemma,
    w.part_of_speech,
    w.gender,
    w.etymology,
    w.friendly_blurb,
    w.created_at,
    w.updated_at,
    r.id_region AS origin_region_id,
    r.code AS origin_region_code,
    r.name AS origin_region_name,
    r.parent_zone AS origin_region_parent_zone,
    r.country_code AS origin_region_country_code,
    r.flag_emoji AS origin_region_flag_emoji
FROM words w
LEFT JOIN regions r ON r.id_region = w.id_origin_region
WHERE w.id_word = $1;

-- 3. Acepciones Aprobadas con su Región Dialectal
-- name: GetApprovedSensesByWordID :many
SELECT 
    s.id_sense,
    s.id_word,
    s.short_definition,
    s.definition,
    s.is_primary,
    s.is_taboo_or_vulgar,
    s.cultural_note,
    s.cefr_level,
    s.usage_frequency,
    s.state,
    r.code AS region_code,
    r.name AS region_name,
    r.country_code,
    r.flag_emoji
FROM senses s
JOIN regions r ON r.id_region = s.id_region
WHERE s.id_word = $1 AND s.state = 'APPROVED'
ORDER BY s.is_primary DESC, r.name ASC;

-- 4. Ejemplos de Uso asociados a la Palabra
-- name: GetExamplesByWordID :many
SELECT 
    e.id_example, 
    e.id_sense, 
    e.sentence, 
    e.source
FROM examples e
JOIN senses s ON s.id_sense = e.id_sense
WHERE s.id_word = $1;

-- 5. Heterónimos Transdialectales (¿Cómo se dice en otros países?)
-- name: GetHeteronymsByWordID :many
SELECT DISTINCT
    r.code AS region_code,
    r.name AS region_name,
    coalesce(r.country_code, '') AS country_code,
    r.flag_emoji,
    w_target.lemma AS term
FROM concept_mappings cm_source
JOIN senses s_source ON s_source.id_sense = cm_source.id_sense
JOIN concept_mappings cm_target ON cm_target.id_concept = cm_source.id_concept
JOIN senses s_target ON s_target.id_sense = cm_target.id_sense AND s_target.state = 'APPROVED'
JOIN words w_target ON w_target.id_word = s_target.id_word
JOIN regions r ON r.id_region = s_target.id_region
WHERE s_source.id_word = $1 AND w_target.id_word != $1
ORDER BY r.name ASC;

-- 6. Listado de Regiones
-- name: ListRegions :many
SELECT id_region, code, name, parent_zone, country_code, flag_emoji
FROM regions
ORDER BY name ASC;

-- 7. Obtener Región por Código
-- name: GetRegionByCode :one
SELECT id_region, code, name, parent_zone, country_code, flag_emoji
FROM regions
WHERE code = $1;

-- 8. Palabra del Día
-- name: GetWordOfTheDayID :one
WITH candidate_words AS (
  SELECT DISTINCT w.id_word
  FROM words w
  JOIN senses s ON s.id_word = w.id_word AND s.state = 'APPROVED'
  JOIN regions r ON r.id_region = s.id_region
  WHERE (sqlc.narg('region_code')::VARCHAR IS NULL OR sqlc.narg('region_code')::VARCHAR = '' OR r.code = sqlc.narg('region_code')::VARCHAR)
  ORDER BY w.id_word ASC
),
indexed_candidates AS (
  SELECT id_word, ROW_NUMBER() OVER (ORDER BY id_word ASC) - 1 AS row_num, COUNT(*) OVER () AS total_count
  FROM candidate_words
)
SELECT id_word
FROM indexed_candidates
WHERE row_num = MOD(EXTRACT(DOY FROM CURRENT_DATE)::INT, total_count)
LIMIT 1;

-- 9. Exploración Paginada de Palabras
-- name: ExploreWordsPaginated :many
WITH filtered_words AS (
  SELECT
    w.id_word,
    w.lemma,
    w.part_of_speech,
    w.gender,
    w.friendly_blurb,
    s.id_sense AS primary_sense_id,
    s.short_definition,
    r.code AS region_code,
    r.name AS region_name,
    r.flag_emoji,
    coalesce(s.cefr_level::TEXT, 'B1') AS cefr_level,
    coalesce(s.usage_frequency::TEXT, 'COMMON') AS usage_frequency,
    s.is_taboo_or_vulgar,
    (SELECT count(*)::INT FROM senses sa WHERE sa.id_word = w.id_word AND sa.state = 'APPROVED') AS approved_senses_count,
    ROW_NUMBER() OVER (
      PARTITION BY w.id_word 
      ORDER BY 
        (r.code = sqlc.narg('region_code')::VARCHAR) DESC NULLS LAST,
        s.is_primary DESC, 
        s.id_sense ASC
    ) AS rn
  FROM words w
  JOIN senses s ON s.id_word = w.id_word AND s.state = 'APPROVED'
  JOIN regions r ON r.id_region = s.id_region
  WHERE
    (sqlc.narg('query')::TEXT IS NULL OR sqlc.narg('query')::TEXT = '' OR 
      w.normalized_lemma ILIKE unaccent(lower(sqlc.narg('query')::TEXT)) || '%' OR
      w.search_vector @@ plainto_tsquery('spanish_unaccent', sqlc.narg('query')::TEXT) OR
      s.sense_search_vector @@ plainto_tsquery('spanish_unaccent', sqlc.narg('query')::TEXT)
    )
    AND (sqlc.narg('region_code')::VARCHAR IS NULL OR sqlc.narg('region_code')::VARCHAR = '' OR r.code = sqlc.narg('region_code')::VARCHAR)
    AND (sqlc.narg('part_of_speech')::VARCHAR IS NULL OR sqlc.narg('part_of_speech')::VARCHAR = '' OR w.part_of_speech = sqlc.narg('part_of_speech')::VARCHAR)
    AND (sqlc.narg('cefr_level')::VARCHAR IS NULL OR sqlc.narg('cefr_level')::VARCHAR = '' OR s.cefr_level::TEXT = sqlc.narg('cefr_level')::VARCHAR)
    AND (sqlc.narg('usage_frequency')::VARCHAR IS NULL OR sqlc.narg('usage_frequency')::VARCHAR = '' OR s.usage_frequency::TEXT = sqlc.narg('usage_frequency')::VARCHAR)
    AND (sqlc.arg('include_taboo')::BOOLEAN = true OR s.is_taboo_or_vulgar = false)
)
SELECT
  id_word,
  lemma,
  part_of_speech,
  gender,
  friendly_blurb,
  primary_sense_id,
  short_definition,
  region_code,
  region_name,
  flag_emoji,
  cefr_level,
  usage_frequency,
  is_taboo_or_vulgar,
  approved_senses_count
FROM filtered_words
WHERE rn = 1
ORDER BY lemma ASC
LIMIT $1 OFFSET $2;

-- 10. Conteo para Exploración Paginada
-- name: CountExploreWords :one
SELECT count(DISTINCT w.id_word)::INT
FROM words w
JOIN senses s ON s.id_word = w.id_word AND s.state = 'APPROVED'
JOIN regions r ON r.id_region = s.id_region
WHERE
  (sqlc.narg('query')::TEXT IS NULL OR sqlc.narg('query')::TEXT = '' OR 
    w.normalized_lemma ILIKE unaccent(lower(sqlc.narg('query')::TEXT)) || '%' OR
    w.search_vector @@ plainto_tsquery('spanish_unaccent', sqlc.narg('query')::TEXT) OR
    s.sense_search_vector @@ plainto_tsquery('spanish_unaccent', sqlc.narg('query')::TEXT)
  )
  AND (sqlc.narg('region_code')::VARCHAR IS NULL OR sqlc.narg('region_code')::VARCHAR = '' OR r.code = sqlc.narg('region_code')::VARCHAR)
  AND (sqlc.narg('part_of_speech')::VARCHAR IS NULL OR sqlc.narg('part_of_speech')::VARCHAR = '' OR w.part_of_speech = sqlc.narg('part_of_speech')::VARCHAR)
  AND (sqlc.narg('cefr_level')::VARCHAR IS NULL OR sqlc.narg('cefr_level')::VARCHAR = '' OR s.cefr_level::TEXT = sqlc.narg('cefr_level')::VARCHAR)
  AND (sqlc.narg('usage_frequency')::VARCHAR IS NULL OR sqlc.narg('usage_frequency')::VARCHAR = '' OR s.usage_frequency::TEXT = sqlc.narg('usage_frequency')::VARCHAR)
  AND (sqlc.arg('include_taboo')::BOOLEAN = true OR s.is_taboo_or_vulgar = false);

-- 11. Búsqueda por Lema Normalizado y Categoría Gramatical
-- name: GetWordByNormalizedLemmaAndPOS :one
SELECT 
    id_word,
    lemma,
    normalized_lemma,
    part_of_speech,
    gender,
    etymology,
    friendly_blurb,
    id_origin_region,
    created_at,
    updated_at
FROM words
WHERE normalized_lemma = $1 AND part_of_speech = $2;

-- 12. Creación de Palabra
-- name: CreateWord :one
INSERT INTO words (
  lemma,
  normalized_lemma,
  part_of_speech,
  gender,
  etymology,
  friendly_blurb,
  id_origin_region
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING 
  id_word,
  lemma,
  normalized_lemma,
  part_of_speech,
  gender,
  etymology,
  friendly_blurb,
  id_origin_region,
  created_at,
  updated_at;

-- 13. Actualización de Metadatos de Palabra
-- name: UpdateWordMetadata :one
UPDATE words SET
  lemma = $2,
  normalized_lemma = $3,
  part_of_speech = $4,
  gender = $5,
  etymology = $6,
  friendly_blurb = $7,
  id_origin_region = $8,
  updated_at = now()
WHERE id_word = $1
RETURNING 
  id_word,
  lemma,
  normalized_lemma,
  part_of_speech,
  gender,
  etymology,
  friendly_blurb,
  id_origin_region,
  created_at,
  updated_at;

-- 14. Creación de Sentido
-- name: CreateSense :one
INSERT INTO senses (
  id_word,
  id_region,
  id_editor,
  short_definition,
  definition,
  is_primary,
  is_taboo_or_vulgar,
  cultural_note,
  cefr_level,
  usage_frequency,
  state
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING 
  id_sense,
  id_word,
  id_region,
  id_editor,
  short_definition,
  definition,
  is_primary,
  is_taboo_or_vulgar,
  cultural_note,
  cefr_level,
  usage_frequency,
  state,
  created_at,
  updated_at;

-- 15. Obtener Sentido por ID con Región
-- name: GetSenseByID :one
SELECT 
    s.id_sense,
    s.id_word,
    s.id_region,
    s.id_editor,
    s.short_definition,
    s.definition,
    s.is_primary,
    s.is_taboo_or_vulgar,
    s.cultural_note,
    s.cefr_level,
    s.usage_frequency,
    s.state,
    s.created_at,
    s.updated_at,
    r.code AS region_code,
    r.name AS region_name,
    r.country_code,
    r.flag_emoji
FROM senses s
JOIN regions r ON r.id_region = s.id_region
WHERE s.id_sense = $1;

-- 16. Actualizar Sentido
-- name: UpdateSense :one
UPDATE senses SET
  id_region = $2,
  short_definition = $3,
  definition = $4,
  is_primary = $5,
  is_taboo_or_vulgar = $6,
  cultural_note = $7,
  cefr_level = $8,
  usage_frequency = $9,
  updated_at = now()
WHERE id_sense = $1
RETURNING 
  id_sense,
  id_word,
  id_region,
  id_editor,
  short_definition,
  definition,
  is_primary,
  is_taboo_or_vulgar,
  cultural_note,
  cefr_level,
  usage_frequency,
  state,
  created_at,
  updated_at;

-- 17. Eliminar Sentido
-- name: DeleteSense :exec
DELETE FROM senses
WHERE id_sense = $1;

-- 18. CRUD de Ejemplos
-- name: CreateExample :one
INSERT INTO examples (
  id_sense,
  id_editor,
  sentence,
  source
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetExampleByID :one
SELECT * FROM examples WHERE id_example = $1;

-- name: UpdateExample :one
UPDATE examples SET
  sentence = $2,
  source = $3,
  updated_at = now()
WHERE id_example = $1
RETURNING *;

-- name: DeleteExample :exec
DELETE FROM examples
WHERE id_example = $1;

-- name: GetExamplesBySenseID :many
SELECT id_example, id_sense, sentence, source
FROM examples
WHERE id_sense = $1
ORDER BY id_example ASC;

-- 19. Conceptos Semánticos y Heteronimia
-- name: ListSemanticConcepts :many
SELECT 
    sc.id_concept,
    sc.concept_key,
    sc.description AS universal_description,
    count(cm.id_sense)::INT AS linked_senses_count
FROM semantic_concepts sc
LEFT JOIN concept_mappings cm ON cm.id_concept = sc.id_concept
WHERE (
    sqlc.narg('query')::TEXT IS NULL 
    OR sqlc.narg('query')::TEXT = ''
    OR sc.concept_key ILIKE '%' || sqlc.narg('query')::TEXT || '%'
    OR sc.description ILIKE '%' || sqlc.narg('query')::TEXT || '%'
)
GROUP BY sc.id_concept, sc.concept_key, sc.description
ORDER BY sc.concept_key ASC;

-- name: CreateSemanticConcept :one
INSERT INTO semantic_concepts (concept_key, description)
VALUES ($1, $2)
RETURNING id_concept, concept_key, description AS universal_description;

-- name: LinkSenseToConcept :exec
INSERT INTO concept_mappings (id_concept, id_sense)
VALUES ($1, $2)
ON CONFLICT (id_concept, id_sense) DO NOTHING;

-- name: UnlinkSenseFromConcept :exec
DELETE FROM concept_mappings
WHERE id_concept = $1 AND id_sense = $2;
