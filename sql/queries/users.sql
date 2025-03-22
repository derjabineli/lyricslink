-- name: CreateUser :one
INSERT INTO users (id, first_name, last_name, email, hashed_password, created_at, updated_at, pc_authorized)
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdatePlanningCenterUser :exec
UPDATE users 
SET avatar = $1, pc_id = $2, pc_authorized = TRUE
WHERE id = $3;

-- name: GetUserByPCID :one
SELECT * FROM users
WHERE pc_id = $1;