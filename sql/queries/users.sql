-- name: CreateUser :one
INSERT INTO users (id, first_name, last_name, email, hashed_password, created_at, updated_at, pc_authorized)
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1;