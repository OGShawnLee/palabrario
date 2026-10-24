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

-- 7. Crear Reporte de Contenido (CU09)
-- name: CreateContentReport :one
INSERT INTO content_reports (
  id_reporter,
  target_type,
  id_word,
  id_sense,
  id_example,
  reason,
  details
) VALUES (
  @id_reporter,
  @target_type::report_target_type,
  @id_word,
  @id_sense,
  @id_example,
  @reason::report_reason,
  @details
) RETURNING 
  id_report,
  id_reporter,
  target_type::TEXT AS target_type,
  id_word,
  id_sense,
  id_example,
  reason::TEXT AS reason,
  details,
  status::TEXT AS status,
  resolved_by,
  resolution_notes,
  created_at,
  resolved_at;

-- 8. Obtener Reporte por ID (CU09)
-- name: GetContentReportByID :one
SELECT 
    cr.id_report,
    cr.id_reporter,
    u.user_name AS reporter_user_name,
    cr.target_type::TEXT AS target_type,
    cr.id_word,
    cr.id_sense,
    cr.id_example,
    COALESCE(
      CASE cr.target_type
        WHEN 'WORD' THEN w.lemma
        WHEN 'SENSE' THEN COALESCE(ws.lemma || ': ' || s.short_definition, s.short_definition)
        WHEN 'EXAMPLE' THEN e.sentence
      END,
      ''
    ) AS target_preview,
    cr.reason::TEXT AS reason,
    cr.details,
    cr.status::TEXT AS status,
    coalesce(cr.resolved_by::TEXT, '')::TEXT AS resolved_by,
    cr.resolution_notes,
    cr.created_at,
    cr.resolved_at
FROM content_reports cr
JOIN users u ON u.id_user = cr.id_reporter
LEFT JOIN words w ON w.id_word = cr.id_word
LEFT JOIN senses s ON s.id_sense = cr.id_sense
LEFT JOIN words ws ON ws.id_word = s.id_word
LEFT JOIN examples e ON e.id_example = cr.id_example
WHERE cr.id_report = $1;

-- 9. Listar Reportes de Contenido (CU09)
-- name: ListContentReports :many
SELECT 
    cr.id_report,
    cr.id_reporter,
    u.user_name AS reporter_user_name,
    cr.target_type::TEXT AS target_type,
    cr.id_word,
    cr.id_sense,
    cr.id_example,
    COALESCE(
      CASE cr.target_type
        WHEN 'WORD' THEN w.lemma
        WHEN 'SENSE' THEN COALESCE(ws.lemma || ': ' || s.short_definition, s.short_definition)
        WHEN 'EXAMPLE' THEN e.sentence
      END,
      ''
    ) AS target_preview,
    cr.reason::TEXT AS reason,
    cr.details,
    cr.status::TEXT AS status,
    coalesce(cr.resolved_by::TEXT, '')::TEXT AS resolved_by,
    cr.resolution_notes,
    cr.created_at,
    cr.resolved_at
FROM content_reports cr
JOIN users u ON u.id_user = cr.id_reporter
LEFT JOIN words w ON w.id_word = cr.id_word
LEFT JOIN senses s ON s.id_sense = cr.id_sense
LEFT JOIN words ws ON ws.id_word = s.id_word
LEFT JOIN examples e ON e.id_example = cr.id_example
WHERE (sqlc.narg('status')::VARCHAR IS NULL OR sqlc.narg('status')::VARCHAR = '' OR cr.status::TEXT = sqlc.narg('status')::VARCHAR)
ORDER BY cr.created_at DESC;

-- 10. Resolver Reporte de Contenido (CU09)
-- name: ResolveContentReport :one
UPDATE content_reports
SET 
  status = @status::report_status,
  resolved_by = @resolved_by,
  resolution_notes = @resolution_notes,
  resolved_at = now()
WHERE id_report = @id_report
RETURNING 
  id_report,
  id_reporter,
  target_type::TEXT AS target_type,
  id_word,
  id_sense,
  id_example,
  reason::TEXT AS reason,
  details,
  status::TEXT AS status,
  resolved_by,
  resolution_notes,
  created_at,
  resolved_at;
