-- name: CreateAccount :one
INSERT INTO accounts (email, password_hash, role)
VALUES ($1, $2, 'LEARNER')
RETURNING id_account, email, password_hash, role, created_at, updated_at;

-- name: CreateUserProfile :one
INSERT INTO users (id_user, name, user_name, anchor_dialect)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAccountByEmail :one
SELECT id_account, email, password_hash, role, created_at, updated_at
FROM accounts
WHERE email = $1;

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
