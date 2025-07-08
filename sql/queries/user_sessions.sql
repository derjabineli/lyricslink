-- name: AddUserSession :one
INSERT INTO user_sessions(id, user_id, access_token, token_type, expires_in, refresh_token, scope, created_at, updated_at, revoked)
VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, NOW(), NOW(), FALSE)
RETURNING *;

-- name: GetSessionRevokedStatus :one
SELECT revoked FROM user_sessions
WHERE id = $1;

-- name: RevokeSession :exec
UPDATE user_sessions
SET revoked = TRUE
WHERE id = $1;

-- name: GetSessionByID :one
SELECT * FROM user_sessions
WHERE id = $1;