-- ==============================================================================
-- MOTOR DE APRENDIZAJE SM-2, MAESTRÍA Y GAMIFICACIÓN (CU12, CU13)
-- ==============================================================================

-- 1. Recuperar Tarjetas de Sesión de Estudio por Baraja
-- name: GetStudySessionCardsByDeck :many
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
    um.id_user AS mastery_user_id,
    coalesce(um.progress_stage::TEXT, 'NEW') AS progress_stage,
    coalesce(um.mastery_percentage, 0.0)::REAL AS mastery_percentage,
    coalesce(um.consecutive_correct, 0)::INT AS consecutive_correct,
    coalesce(um.times_reviewed, 0)::INT AS times_reviewed,
    coalesce(um.easiness_factor, 2.5)::REAL AS easiness_factor,
    coalesce(um.current_interval_days, 1)::INT AS current_interval_days,
    coalesce(um.next_review_at, now()) AS next_review_at,
    coalesce(um.last_reviewed_at, now()) AS last_reviewed_at
FROM deck_items di
JOIN senses s ON s.id_sense = di.id_sense
JOIN words w ON w.id_word = s.id_word
JOIN regions r ON r.id_region = s.id_region
LEFT JOIN user_mastery um ON um.id_sense = s.id_sense AND um.id_user = @user_id
WHERE di.id_deck = @deck_id
  AND s.state != 'HIDDEN'
  AND (um.id_user IS NULL OR um.next_review_at <= now())
ORDER BY COALESCE(um.next_review_at, '1970-01-01'::timestamptz) ASC, di.sort_order ASC
LIMIT @limit_val;

-- 2. Conteo de Tarjetas Pendientes de Estudio por Baraja
-- name: CountDueStudySessionCardsByDeck :one
SELECT count(*)::INT
FROM deck_items di
JOIN senses s ON s.id_sense = di.id_sense
LEFT JOIN user_mastery um ON um.id_sense = s.id_sense AND um.id_user = @user_id
WHERE di.id_deck = @deck_id
  AND s.state != 'HIDDEN'
  AND (um.id_user IS NULL OR um.next_review_at <= now());

-- 3. Recuperar Tarjetas de Sesión de Estudio de Todo el Studio (Propias + Guardadas)
-- name: GetStudySessionCardsAll :many
SELECT DISTINCT ON (s.id_sense)
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
    um.id_user AS mastery_user_id,
    coalesce(um.progress_stage::TEXT, 'NEW') AS progress_stage,
    coalesce(um.mastery_percentage, 0.0)::REAL AS mastery_percentage,
    coalesce(um.consecutive_correct, 0)::INT AS consecutive_correct,
    coalesce(um.times_reviewed, 0)::INT AS times_reviewed,
    coalesce(um.easiness_factor, 2.5)::REAL AS easiness_factor,
    coalesce(um.current_interval_days, 1)::INT AS current_interval_days,
    coalesce(um.next_review_at, now()) AS next_review_at,
    coalesce(um.last_reviewed_at, now()) AS last_reviewed_at
FROM decks d
JOIN deck_items di ON di.id_deck = d.id_deck
JOIN senses s ON s.id_sense = di.id_sense
JOIN words w ON w.id_word = s.id_word
JOIN regions r ON r.id_region = s.id_region
LEFT JOIN saved_decks sd ON sd.id_deck = d.id_deck AND sd.id_user = @user_id
LEFT JOIN user_mastery um ON um.id_sense = s.id_sense AND um.id_user = @user_id
WHERE (d.id_user = @user_id OR sd.id_user IS NOT NULL)
  AND s.state != 'HIDDEN'
  AND (um.id_user IS NULL OR um.next_review_at <= now())
ORDER BY s.id_sense, COALESCE(um.next_review_at, '1970-01-01'::timestamptz) ASC
LIMIT @limit_val;

-- 4. Conteo de Tarjetas Pendientes en Todo el Studio
-- name: CountDueStudySessionCardsAll :one
SELECT count(DISTINCT s.id_sense)::INT
FROM decks d
JOIN deck_items di ON di.id_deck = d.id_deck
JOIN senses s ON s.id_sense = di.id_sense
LEFT JOIN saved_decks sd ON sd.id_deck = d.id_deck AND sd.id_user = @user_id
LEFT JOIN user_mastery um ON um.id_sense = s.id_sense AND um.id_user = @user_id
WHERE (d.id_user = @user_id OR sd.id_user IS NOT NULL)
  AND s.state != 'HIDDEN'
  AND (um.id_user IS NULL OR um.next_review_at <= now());

-- 5. Obtener Estado de Maestría de un Sentido Específico
-- name: GetUserMasteryForSense :one
SELECT 
    id_user,
    id_sense,
    progress_stage::TEXT AS progress_stage,
    mastery_percentage,
    consecutive_correct,
    times_reviewed,
    easiness_factor,
    current_interval_days,
    next_review_at,
    last_reviewed_at
FROM user_mastery
WHERE id_user = $1 AND id_sense = $2;

-- 6. Upsert Estado de Maestría SM-2
-- name: UpsertUserMastery :one
INSERT INTO user_mastery (
    id_user,
    id_sense,
    progress_stage,
    mastery_percentage,
    consecutive_correct,
    times_reviewed,
    easiness_factor,
    current_interval_days,
    next_review_at,
    last_reviewed_at
) VALUES (
    @id_user,
    @id_sense,
    @progress_stage::mastery_stage,
    @mastery_percentage,
    @consecutive_correct,
    @times_reviewed,
    @easiness_factor,
    @current_interval_days,
    @next_review_at,
    @last_reviewed_at
)
ON CONFLICT (id_user, id_sense) DO UPDATE SET
    progress_stage = EXCLUDED.progress_stage,
    mastery_percentage = EXCLUDED.mastery_percentage,
    consecutive_correct = EXCLUDED.consecutive_correct,
    times_reviewed = EXCLUDED.times_reviewed,
    easiness_factor = EXCLUDED.easiness_factor,
    current_interval_days = EXCLUDED.current_interval_days,
    next_review_at = EXCLUDED.next_review_at,
    last_reviewed_at = EXCLUDED.last_reviewed_at
RETURNING 
    id_user,
    id_sense,
    progress_stage::TEXT AS progress_stage,
    mastery_percentage,
    consecutive_correct,
    times_reviewed,
    easiness_factor,
    current_interval_days,
    next_review_at,
    last_reviewed_at;

-- 7. Registrar Transacción en Ledger de XP (Inmutable)
-- name: InsertXPTransaction :one
INSERT INTO xp_transactions (
    id_user,
    amount_xp,
    reason,
    metadata_id,
    created_at
) VALUES (
    @id_user,
    @amount_xp,
    @reason::xp_reason,
    @metadata_id,
    now()
) RETURNING 
    id_transaction,
    id_user,
    amount_xp,
    reason::TEXT AS reason,
    metadata_id,
    created_at;

-- 8. Sumar XP y Actualizar Racha y Nivel Atómicamente
-- name: AddUserXPAndUpdateStreak :one
UPDATE users
SET 
  total_xp = total_xp + @amount_xp,
  current_level = ((total_xp + @amount_xp) / 1000) + 1,
  current_streak = CASE
    WHEN current_streak = 0 THEN 1
    WHEN last_activity_at::DATE = CURRENT_DATE THEN current_streak
    WHEN last_activity_at::DATE = CURRENT_DATE - INTERVAL '1 day' THEN current_streak + 1
    ELSE 1
  END,
  longest_streak = GREATEST(
    longest_streak,
    CASE
      WHEN current_streak = 0 THEN 1
      WHEN last_activity_at::DATE = CURRENT_DATE THEN current_streak
      WHEN last_activity_at::DATE = CURRENT_DATE - INTERVAL '1 day' THEN current_streak + 1
      ELSE 1
    END
  ),
  last_activity_at = now(),
  updated_at = now()
WHERE id_user = @id_user
RETURNING id_user, total_xp, current_level, current_streak, longest_streak;

-- 9. Conteo y Estadísticas del Tablero de Maestría
-- name: GetMasteryStageCounts :one
SELECT 
    count(*)::INT AS total_studied,
    count(*) FILTER (WHERE next_review_at <= now())::INT AS due_today,
    count(*) FILTER (WHERE progress_stage = 'NEW')::INT AS count_new,
    count(*) FILTER (WHERE progress_stage = 'LEARNING')::INT AS count_learning,
    count(*) FILTER (WHERE progress_stage = 'REVIEWING')::INT AS count_reviewing,
    count(*) FILTER (WHERE progress_stage = 'MASTERED')::INT AS count_mastered,
    COALESCE(avg(mastery_percentage), 0.0)::REAL AS global_retention
FROM user_mastery
WHERE id_user = $1;

-- 10. Perfil de Gamificación del Usuario
-- name: GetUserGamificationProfile :one
SELECT 
    id_user,
    total_xp,
    current_level,
    current_streak,
    longest_streak
FROM users
WHERE id_user = $1;

-- 11. Listar Logros con Estado para el Usuario
-- name: ListUserAchievements :many
SELECT 
    a.id_achievement,
    a.slug,
    a.title,
    a.description,
    coalesce(a.badge_icon_url, '') AS badge_icon_url,
    a.xp_reward,
    a.tier::TEXT AS tier,
    a.category::TEXT AS category,
    (ua.id_user IS NOT NULL)::BOOLEAN AS is_unlocked,
    ua.unlocked_at
FROM achievements a
LEFT JOIN user_achievements ua ON ua.id_achievement = a.id_achievement AND ua.id_user = $1
ORDER BY a.id_achievement ASC;

-- 12. Desbloquear Logro
-- name: UnlockUserAchievement :exec
INSERT INTO user_achievements (id_user, id_achievement, unlocked_at)
VALUES ($1, $2, now())
ON CONFLICT (id_user, id_achievement) DO NOTHING;

-- 13. Listar Tarjetas de Maestría por Etapa
-- name: ListMasteryCardsByStage :many
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
    um.id_user AS mastery_user_id,
    um.progress_stage::TEXT AS progress_stage,
    um.mastery_percentage,
    um.consecutive_correct,
    um.times_reviewed,
    um.easiness_factor,
    um.current_interval_days,
    um.next_review_at,
    um.last_reviewed_at
FROM user_mastery um
JOIN senses s ON s.id_sense = um.id_sense
JOIN words w ON w.id_word = s.id_word
JOIN regions r ON r.id_region = s.id_region
WHERE um.id_user = @user_id
  AND (sqlc.narg('stage')::VARCHAR IS NULL OR sqlc.narg('stage')::VARCHAR = '' OR um.progress_stage::TEXT = sqlc.narg('stage')::VARCHAR)
ORDER BY um.last_reviewed_at DESC;

-- 14. Listar Transacciones Recientes de XP
-- name: ListRecentXPTransactions :many
SELECT 
    id_transaction,
    id_user,
    amount_xp,
    reason::TEXT AS reason,
    metadata_id,
    created_at
FROM xp_transactions
WHERE id_user = $1
ORDER BY created_at DESC
LIMIT $2;

-- 15. Ejemplos de Sentidos
-- name: GetExamplesForSenses :many
SELECT 
    e.id_example,
    e.id_sense,
    e.id_editor,
    e.sentence,
    e.source,
    e.created_at,
    e.updated_at
FROM examples e
WHERE e.id_sense = ANY($1::int[])
ORDER BY e.id_example ASC;
