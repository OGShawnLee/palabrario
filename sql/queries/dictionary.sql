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
