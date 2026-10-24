-- ==============================================================================
-- GOBERNANZA INSTITUCIONAL, ROLES Y BOLETOS (CU02, CU03, CU04)
-- ==============================================================================

-- 1. Crear Solicitud de Editor (CU02)
-- name: CreateEditorRequest :one
INSERT INTO editor_requests (
    id_user,
    institution_name,
    justification,
    credentials_url,
    status
) VALUES (
    $1, $2, $3, $4, 'PENDING'
) RETURNING 
    id_request,
    id_user,
    institution_name,
    justification,
    credentials_url,
    status,
    reviewed_by,
    review_notes,
    created_at,
    resolved_at;

-- 2. Listar Mis Solicitudes de Editor (CU02)
-- name: ListMyEditorRequests :many
SELECT 
    er.id_request,
    er.id_user,
    u.user_name,
    u.name AS user_full_name,
    a.email AS user_email,
    er.institution_name,
    er.justification,
    er.credentials_url,
    er.status::TEXT AS status,
    coalesce(er.reviewed_by::TEXT, '')::TEXT AS reviewed_by,
    er.review_notes,
    er.created_at,
    er.resolved_at
FROM editor_requests er
JOIN users u ON u.id_user = er.id_user
JOIN accounts a ON a.id_account = er.id_user
WHERE er.id_user = $1
ORDER BY er.created_at DESC;

-- 3. Listar Todas las Solicitudes de Editor para Admin (CU03)
-- name: ListEditorRequests :many
SELECT 
    er.id_request,
    er.id_user,
    u.user_name,
    u.name AS user_full_name,
    a.email AS user_email,
    er.institution_name,
    er.justification,
    er.credentials_url,
    er.status::TEXT AS status,
    coalesce(er.reviewed_by::TEXT, '')::TEXT AS reviewed_by,
    er.review_notes,
    er.created_at,
    er.resolved_at
FROM editor_requests er
JOIN users u ON u.id_user = er.id_user
JOIN accounts a ON a.id_account = er.id_user
WHERE (sqlc.narg('status')::VARCHAR IS NULL OR sqlc.narg('status')::VARCHAR = '' OR er.status::TEXT = sqlc.narg('status')::VARCHAR)
ORDER BY er.created_at DESC;

-- 4. Obtener Solicitud por ID
-- name: GetEditorRequestByID :one
SELECT 
    er.id_request,
    er.id_user,
    u.user_name,
    u.name AS user_full_name,
    a.email AS user_email,
    er.institution_name,
    er.justification,
    er.credentials_url,
    er.status::TEXT AS status,
    coalesce(er.reviewed_by::TEXT, '')::TEXT AS reviewed_by,
    er.review_notes,
    er.created_at,
    er.resolved_at
FROM editor_requests er
JOIN users u ON u.id_user = er.id_user
JOIN accounts a ON a.id_account = er.id_user
WHERE er.id_request = $1;

-- 5. Resolver Solicitud de Editor (CU03)
-- name: ResolveEditorRequest :one
UPDATE editor_requests
SET 
    status = $2,
    reviewed_by = $3,
    review_notes = $4,
    resolved_at = now()
WHERE id_request = $1
RETURNING 
    id_request,
    id_user,
    institution_name,
    justification,
    credentials_url,
    status,
    reviewed_by,
    review_notes,
    created_at,
    resolved_at;

-- 6. Actualizar Rol de Cuenta
-- name: UpdateAccountRole :exec
UPDATE accounts
SET role = $2, updated_at = now()
WHERE id_account = $1;

-- 7. Contar Administradores Activos (Protección del último admin CU03 4a)
-- name: CountActiveAdmins :one
SELECT count(*)::INT
FROM accounts
WHERE role = 'ADMIN' AND is_active = true;

-- 8. Actualizar Rol y/o Estado Activo de un Usuario (CU03)
-- name: UpdateUserRoleOrStatus :one
WITH updated_acc AS (
  UPDATE accounts
  SET 
    role = COALESCE(sqlc.narg('role')::account_role, role),
    is_active = COALESCE(sqlc.narg('is_active')::BOOLEAN, is_active),
    updated_at = now()
  WHERE id_account = @id_user
  RETURNING id_account, email, role, is_active, updated_at
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
FROM users u
JOIN updated_acc a ON a.id_account = u.id_user;

-- 9. Listar Usuarios Paginados con Filtros (CU03)
-- name: ListUsers :many
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
WHERE (
    sqlc.narg('search')::VARCHAR IS NULL OR sqlc.narg('search')::VARCHAR = ''
    OR u.name ILIKE '%' || sqlc.narg('search')::VARCHAR || '%'
    OR u.user_name ILIKE '%' || sqlc.narg('search')::VARCHAR || '%'
    OR a.email ILIKE '%' || sqlc.narg('search')::VARCHAR || '%'
)
AND (
    sqlc.narg('role')::VARCHAR IS NULL OR sqlc.narg('role')::VARCHAR = ''
    OR a.role::TEXT = sqlc.narg('role')::VARCHAR
)
ORDER BY u.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- 10. Contar Usuarios con Filtros (CU03)
-- name: CountUsers :one
SELECT count(*)::INT
FROM users u
JOIN accounts a ON a.id_account = u.id_user
WHERE (
    sqlc.narg('search')::VARCHAR IS NULL OR sqlc.narg('search')::VARCHAR = ''
    OR u.name ILIKE '%' || sqlc.narg('search')::VARCHAR || '%'
    OR u.user_name ILIKE '%' || sqlc.narg('search')::VARCHAR || '%'
    OR a.email ILIKE '%' || sqlc.narg('search')::VARCHAR || '%'
)
AND (
    sqlc.narg('role')::VARCHAR IS NULL OR sqlc.narg('role')::VARCHAR = ''
    OR a.role::TEXT = sqlc.narg('role')::VARCHAR
);

-- 11. Crear Boleto de Invitación (CU04)
-- name: CreateInvitationTicket :one
INSERT INTO invitation_tickets (
    code,
    created_by,
    granted_role,
    institution_name,
    max_uses,
    current_uses,
    is_active,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, 0, true, $6
) RETURNING 
    id_ticket,
    code,
    created_by,
    granted_role,
    institution_name,
    max_uses,
    current_uses,
    is_active,
    expires_at,
    created_at;

-- 12. Listar Boletos de Invitación (CU04)
-- name: ListInvitationTickets :many
SELECT 
    id_ticket,
    code,
    created_by,
    granted_role,
    institution_name,
    max_uses,
    current_uses,
    is_active,
    expires_at,
    created_at
FROM invitation_tickets
WHERE (
    sqlc.narg('active_only')::BOOLEAN IS NULL 
    OR sqlc.narg('active_only')::BOOLEAN = false 
    OR (is_active = true AND expires_at > now() AND current_uses < max_uses)
)
ORDER BY created_at DESC;

-- 13. Obtener Boleto de Invitación por ID
-- name: GetInvitationTicketByID :one
SELECT 
    id_ticket,
    code,
    created_by,
    granted_role,
    institution_name,
    max_uses,
    current_uses,
    is_active,
    expires_at,
    created_at
FROM invitation_tickets
WHERE id_ticket = $1;

-- 14. Revocar Boleto de Invitación (CU04)
-- name: RevokeInvitationTicket :exec
UPDATE invitation_tickets
SET is_active = false
WHERE id_ticket = $1;

-- 15. Bloquear Boleto por Código para Redención Concurrente Segura (CU04 FOR UPDATE)
-- name: GetInvitationTicketByCodeForUpdate :one
SELECT 
    id_ticket,
    code,
    created_by,
    granted_role,
    institution_name,
    max_uses,
    current_uses,
    is_active,
    expires_at,
    created_at
FROM invitation_tickets
WHERE code = $1
FOR UPDATE;

-- 16. Verificar si un Usuario ya Redimió un Boleto
-- name: HasUserRedeemedTicket :one
SELECT EXISTS(
    SELECT 1 FROM ticket_redemptions
    WHERE id_ticket = $1 AND id_user = $2
)::BOOLEAN;

-- 17. Registrar Redención de Boleto (CU04)
-- name: CreateTicketRedemption :exec
INSERT INTO ticket_redemptions (
    id_ticket,
    id_user,
    redeemed_at
) VALUES (
    $1, $2, now()
);

-- 18. Incrementar Usos de Boleto (CU04)
-- name: IncrementTicketUses :one
UPDATE invitation_tickets
SET current_uses = current_uses + 1
WHERE id_ticket = $1
RETURNING current_uses, max_uses;
