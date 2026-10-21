-- ==============================================================================
-- CURADURÍA EDITORIAL POR QUÓRUM Y MODERACIÓN (sqlc)
-- ==============================================================================

-- 1. Cola de Revisión de Quórum Pendiente (state = 'DRAFT')
-- name: ListPendingQuorumQueue :many
SELECT 
    s.id_sense,
    s.id_word,
    w.lemma,
    w.part_of_speech,
    r.code AS region_code,
    r.name AS region_name,
    r.flag_emoji,
    s.short_definition,
    s.definition,
    s.cultural_note,
    coalesce(s.cefr_level::TEXT, 'B1') AS cefr_level,
    coalesce(s.usage_frequency::TEXT, 'COMMON') AS usage_frequency,
    s.is_taboo_or_vulgar,
    s.state::TEXT AS state,
    u.id_user AS author_id,
    u.user_name AS author_user_name,
    u.name AS author_full_name,
    (SELECT count(*)::INT FROM sense_approvals sa WHERE sa.id_sense = s.id_sense AND sa.decision = 'APPROVE') AS approvals_count,
    3::INT AS required_quorum,
    s.created_at
FROM senses s
JOIN words w ON w.id_word = s.id_word
JOIN regions r ON r.id_region = s.id_region
JOIN users u ON u.id_user = s.id_editor
WHERE s.state = 'DRAFT'
  AND (sqlc.narg('region_code')::VARCHAR IS NULL OR sqlc.narg('region_code')::VARCHAR = '' OR r.code = sqlc.narg('region_code')::VARCHAR)
ORDER BY s.created_at ASC;

-- 2. Propuestas Creadas por un Editor
-- name: ListEditorProposals :many
SELECT 
    s.id_sense,
    s.id_word,
    w.lemma,
    w.part_of_speech,
    r.code AS region_code,
    r.name AS region_name,
    r.flag_emoji,
    s.short_definition,
    s.definition,
    s.cultural_note,
    coalesce(s.cefr_level::TEXT, 'B1') AS cefr_level,
    coalesce(s.usage_frequency::TEXT, 'COMMON') AS usage_frequency,
    s.is_taboo_or_vulgar,
    s.state::TEXT AS state,
    u.id_user AS author_id,
    u.user_name AS author_user_name,
    u.name AS author_full_name,
    (SELECT count(*)::INT FROM sense_approvals sa WHERE sa.id_sense = s.id_sense AND sa.decision = 'APPROVE') AS approvals_count,
    3::INT AS required_quorum,
    s.created_at
FROM senses s
JOIN words w ON w.id_word = s.id_word
JOIN regions r ON r.id_region = s.id_region
JOIN users u ON u.id_user = s.id_editor
WHERE s.id_editor = $1
  AND (sqlc.narg('state')::VARCHAR IS NULL OR sqlc.narg('state')::VARCHAR = '' OR s.state::TEXT = sqlc.narg('state')::VARCHAR)
ORDER BY s.created_at DESC;

-- 3. Detalle de Ítem de Quórum por ID de Sentido
-- name: GetSenseQuorumItemByID :one
SELECT 
    s.id_sense,
    s.id_word,
    w.lemma,
    w.part_of_speech,
    r.code AS region_code,
    r.name AS region_name,
    r.flag_emoji,
    s.short_definition,
    s.definition,
    s.cultural_note,
    coalesce(s.cefr_level::TEXT, 'B1') AS cefr_level,
    coalesce(s.usage_frequency::TEXT, 'COMMON') AS usage_frequency,
    s.is_taboo_or_vulgar,
    s.state::TEXT AS state,
    u.id_user AS author_id,
    u.user_name AS author_user_name,
    u.name AS author_full_name,
    (SELECT count(*)::INT FROM sense_approvals sa WHERE sa.id_sense = s.id_sense AND sa.decision = 'APPROVE') AS approvals_count,
    3::INT AS required_quorum,
    s.created_at
FROM senses s
JOIN words w ON w.id_word = s.id_word
JOIN regions r ON r.id_region = s.id_region
JOIN users u ON u.id_user = s.id_editor
WHERE s.id_sense = $1;

-- 4. Votos Emitidos sobre un Sentido
-- name: ListVotesBySenseID :many
SELECT 
    sa.id_sense,
    sa.id_editor,
    u.user_name AS editor_user_name,
    u.name AS editor_full_name,
    u.anchor_dialect AS editor_dialect,
    sa.decision,
    sa.feedback_notes,
    sa.approved_at
FROM sense_approvals sa
JOIN users u ON u.id_user = sa.id_editor
WHERE sa.id_sense = $1
ORDER BY sa.approved_at ASC;

-- 5. Emitir o Actualizar Voto de Aprobación
-- name: UpsertSenseApproval :one
INSERT INTO sense_approvals (
  id_sense,
  id_editor,
  decision,
  feedback_notes,
  approved_at
) VALUES (
  $1, $2, $3, $4, now()
)
ON CONFLICT (id_sense, id_editor) DO UPDATE SET
  decision = EXCLUDED.decision,
  feedback_notes = EXCLUDED.feedback_notes,
  approved_at = now()
RETURNING id_sense, id_editor, decision, feedback_notes, approved_at;

-- 6. Administrador: Forzar Estado de un Sentido (RN08)
-- name: AdminUpdateSenseState :one
UPDATE senses SET
  state = $2,
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
