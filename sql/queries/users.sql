-- name: CreateAccount :one
INSERT INTO accounts (email, password_hash, role)
VALUES ($1, $2, 'LEARNER')
RETURNING id_account, email, password_hash, role, is_active, created_at, updated_at;

-- name: CreateUserProfile :one
INSERT INTO users (id_user, name, user_name, anchor_dialect)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAccountByEmail :one
SELECT id_account, email, password_hash, role, is_active, created_at, updated_at
FROM accounts
WHERE email = $1;

-- name: GetAccountByID :one
SELECT id_account, email, password_hash, role, is_active, created_at, updated_at
FROM accounts
WHERE id_account = $1;

-- name: UpdateAccountPassword :exec
UPDATE accounts
SET password_hash = $2, updated_at = now()
WHERE id_account = $1;

-- name: GetUserProfileByID :one
SELECT 
    u.id_user,
    a.email,
    u.name,
    u.user_name,
    u.avatar_url,
    u.anchor_dialect,
    u.total_xp,
    u.current_level,
    u.current_streak,
    u.longest_streak,
    a.role,
    u.created_at,
    u.updated_at
FROM users u
JOIN accounts a ON a.id_account = u.id_user
WHERE u.id_user = $1;

-- name: GetUserProfileByUsername :one
SELECT 
    u.id_user,
    a.email,
    u.name,
    u.user_name,
    u.avatar_url,
    u.anchor_dialect,
    u.total_xp,
    u.current_level,
    u.current_streak,
    u.longest_streak,
    a.role,
    u.created_at,
    u.updated_at
FROM users u
JOIN accounts a ON a.id_account = u.id_user
WHERE LOWER(u.user_name) = LOWER($1);

-- name: UpdateUserProfile :one
WITH updated AS (
  UPDATE users
  SET 
    name = $2,
    anchor_dialect = $3,
    avatar_url = $4,
    updated_at = now()
  WHERE id_user = $1
  RETURNING *
)
SELECT 
    u.id_user,
    a.email,
    u.name,
    u.user_name,
    u.avatar_url,
    u.anchor_dialect,
    u.total_xp,
    u.current_level,
    u.current_streak,
    u.longest_streak,
    a.role,
    u.created_at,
    u.updated_at
FROM updated u
JOIN accounts a ON a.id_account = u.id_user;

-- name: CountUserPrivateDecks :one
SELECT count(*)::INT
FROM decks
WHERE id_user = $1 AND is_public = false;

-- name: ListUserPublicDecks :many
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
WHERE d.id_user = $1 AND d.is_public = true
ORDER BY d.created_at DESC;

-- name: ListUserPublicCollections :many
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
    (SELECT count(*)::INT FROM collection_decks cd WHERE cd.id_collection = c.id_collection) AS deck_count
FROM collections c
WHERE c.id_user = $1 AND c.is_public = true
ORDER BY c.created_at DESC;

-- name: GetUserArenaStats :one
SELECT
    (SELECT count(*)::INT FROM matches WHERE winner_user_id = $1 AND status = 'FINISHED') AS matches_won,
    (SELECT count(*)::INT FROM match_participants WHERE id_user = $1 AND final_rank BETWEEN 1 AND 3) AS podium_count;
