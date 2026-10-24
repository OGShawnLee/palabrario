-- ==============================================================================
-- CONSULTAS TIPADAS SQLC: ARENA MULTIJUGADOR Y REPLAY (CU14, CU15, CU16)
-- ==============================================================================

-- name: GetApprovedDeckSensesForArena :many
SELECT
  s.id_sense,
  s.id_word,
  w.lemma,
  w.part_of_speech,
  s.short_definition,
  s.definition,
  s.cultural_note,
  r.code AS region_code,
  r.name AS region_name,
  r.flag_emoji
FROM deck_items di
JOIN senses s ON s.id_sense = di.id_sense
JOIN words w ON w.id_word = s.id_word
JOIN regions r ON r.id_region = s.id_region
WHERE di.id_deck = $1
  AND s.state = 'APPROVED'
ORDER BY di.sort_order ASC, s.id_sense ASC;

-- name: GetCuratedQuestionsBySenseIDs :many
SELECT
  id_question,
  id_sense,
  id_editor,
  type,
  prompt,
  custom_distractors,
  explanation,
  created_at
FROM questions
WHERE id_sense = ANY($1::int[])
ORDER BY id_question ASC;

-- name: CreateMatch :one
INSERT INTO matches (
  id_deck,
  host_user_id,
  room_code,
  status,
  is_public,
  total_rounds,
  round_time_sec
) VALUES (
  $1, $2, $3, 'WAITING', $4, $5, $6
)
RETURNING
  id_match,
  id_deck,
  host_user_id,
  room_code,
  status,
  winner_user_id,
  is_public,
  total_rounds,
  round_time_sec,
  created_at,
  finished_at;

-- name: InsertMatchRound :exec
INSERT INTO match_rounds (
  id_match,
  round_number,
  id_sense,
  id_question,
  modality,
  prompt,
  option_a,
  option_b,
  option_c,
  option_d,
  correct_option,
  explanation
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
);

-- name: GetMatchByRoomCode :one
SELECT
  m.id_match,
  m.id_deck,
  d.title AS deck_title,
  m.host_user_id,
  m.room_code,
  m.status,
  m.winner_user_id,
  m.is_public,
  m.total_rounds,
  m.round_time_sec,
  m.created_at,
  m.finished_at
FROM matches m
JOIN decks d ON d.id_deck = m.id_deck
WHERE UPPER(m.room_code) = UPPER($1)
LIMIT 1;

-- name: GetMatchByID :one
SELECT
  m.id_match,
  m.id_deck,
  d.title AS deck_title,
  m.host_user_id,
  m.room_code,
  m.status,
  m.winner_user_id,
  m.is_public,
  m.total_rounds,
  m.round_time_sec,
  m.created_at,
  m.finished_at
FROM matches m
JOIN decks d ON d.id_deck = m.id_deck
WHERE m.id_match = $1
LIMIT 1;

-- name: ListMatchRoundsByMatchID :many
SELECT
  mr.id_match,
  mr.round_number,
  mr.id_sense,
  w.lemma,
  mr.id_question,
  mr.modality,
  mr.prompt,
  mr.option_a,
  mr.option_b,
  mr.option_c,
  mr.option_d,
  mr.correct_option,
  mr.explanation
FROM match_rounds mr
JOIN senses s ON s.id_sense = mr.id_sense
JOIN words w ON w.id_word = s.id_word
WHERE mr.id_match = $1
ORDER BY mr.round_number ASC;

-- name: UpdateMatchStatus :exec
UPDATE matches
SET status = $2
WHERE id_match = $1;

-- name: MarkMatchFinished :one
UPDATE matches
SET
  status = 'FINISHED',
  winner_user_id = $2,
  finished_at = now()
WHERE id_match = $1
RETURNING
  id_match,
  id_deck,
  host_user_id,
  room_code,
  status,
  winner_user_id,
  is_public,
  total_rounds,
  round_time_sec,
  created_at,
  finished_at;

-- name: UpsertMatchParticipant :exec
INSERT INTO match_participants (
  id_match,
  id_user,
  score,
  final_rank,
  correct_answers,
  total_time_spent_ms,
  xp_earned
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (id_match, id_user) DO UPDATE SET
  score = EXCLUDED.score,
  final_rank = EXCLUDED.final_rank,
  correct_answers = EXCLUDED.correct_answers,
  total_time_spent_ms = EXCLUDED.total_time_spent_ms,
  xp_earned = EXCLUDED.xp_earned;

-- name: UpsertMatchRoundQuizAttempt :exec
INSERT INTO quiz_attempts (
  id_user,
  id_sense,
  id_question,
  id_match,
  round_number,
  selected_option,
  is_correct,
  quality_rating,
  time_spent_ms,
  points_earned
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
ON CONFLICT (id_match, round_number, id_user) DO UPDATE SET
  selected_option = EXCLUDED.selected_option,
  is_correct = EXCLUDED.is_correct,
  quality_rating = EXCLUDED.quality_rating,
  time_spent_ms = EXCLUDED.time_spent_ms,
  points_earned = EXCLUDED.points_earned,
  answered_at = now();

-- name: ListUserMatchHistory :many
SELECT
  m.id_match,
  m.room_code,
  m.id_deck,
  d.title AS deck_title,
  m.winner_user_id,
  m.status,
  m.total_rounds,
  m.round_time_sec,
  m.created_at,
  m.finished_at
FROM matches m
JOIN decks d ON d.id_deck = m.id_deck
JOIN match_participants mp ON mp.id_match = m.id_match
WHERE mp.id_user = $1
  AND m.status = 'FINISHED'
ORDER BY m.finished_at DESC NULLS LAST, m.created_at DESC
LIMIT $2;

-- name: ListMatchParticipantsByMatchIDs :many
SELECT
  mp.id_match,
  mp.id_user,
  u.user_name,
  u.name,
  coalesce(u.avatar_url, '') AS avatar_url,
  u.anchor_dialect,
  coalesce(r.flag_emoji, '🌐') AS flag_emoji,
  mp.score,
  mp.final_rank,
  mp.correct_answers,
  mp.total_time_spent_ms,
  mp.xp_earned,
  (m.host_user_id = mp.id_user) AS is_host
FROM match_participants mp
JOIN matches m ON m.id_match = mp.id_match
JOIN users u ON u.id_user = mp.id_user
LEFT JOIN regions r ON r.code = u.anchor_dialect
WHERE mp.id_match = ANY($1::uuid[])
ORDER BY mp.id_match, mp.final_rank ASC, mp.score DESC, mp.total_time_spent_ms ASC;

-- name: ListMatchQuizAttemptsByMatchID :many
SELECT
  qa.id_attempt,
  qa.id_match,
  qa.round_number,
  qa.id_user,
  u.user_name,
  u.name,
  coalesce(u.avatar_url, '') AS avatar_url,
  u.anchor_dialect,
  coalesce(r.flag_emoji, '🌐') AS flag_emoji,
  qa.id_sense,
  qa.selected_option,
  qa.is_correct,
  qa.time_spent_ms,
  qa.points_earned,
  qa.answered_at
FROM quiz_attempts qa
JOIN users u ON u.id_user = qa.id_user
LEFT JOIN regions r ON r.code = u.anchor_dialect
WHERE qa.id_match = $1
ORDER BY qa.round_number ASC, qa.points_earned DESC, qa.time_spent_ms ASC;
