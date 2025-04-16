-- name: CreateUser :one
INSERT INTO users (id, first_name, last_name, email, created_at, updated_at, avatar, pc_id, administrator)
VALUES(gen_random_uuid(), $1, $2, $3, NOW(), NOW(), $4, $5, $6)
ON CONFLICT (pc_id) DO UPDATE 
    SET email = EXCLUDED.email,
        updated_at = NOW(),
        avatar = EXCLUDED.avatar,
        administrator = EXCLUDED.administrator
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdatePlanningCenterUser :exec
UPDATE users 
SET avatar = $1, pc_id = $2
WHERE id = $3;

-- name: GetUserByPCID :one
SELECT * FROM users
WHERE pc_id = $1;