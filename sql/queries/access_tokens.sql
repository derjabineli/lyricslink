-- name: AddAccessToken :one
INSERT INTO planning_center_tokens(id, user_id, access_token, token_type, expires_in, refresh_token, scope, created_at, updated_at, revoked)
VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, NOW(), NOW(), FALSE)
ON CONFLICT (user_id) DO UPDATE 
    SET access_token = EXCLUDED.access_token,
        token_type = EXCLUDED.token_type,
        expires_in = EXCLUDED.expires_in,
        refresh_token = EXCLUDED.refresh_token,
        scope = EXCLUDED.scope,
        updated_at = NOW()
RETURNING *;

-- name: GetTokenByUserID :one
SELECT * FROM user_sessions
WHERE user_id = $1 AND revoked = FALSE;

-- name: UpdateUserToken :one
UPDATE user_sessions
SET access_token = $1,
    refresh_token = $2,
    scope = $3,
    updated_at = NOW()
WHERE id = $4
RETURNING *;
