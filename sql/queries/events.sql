-- name: CreateEvent :one
INSERT INTO events (id, created_at, updated_at, name, date, user_id)
VALUES(gen_random_uuid(), NOW(), NOW(), $1, $2, $3)
RETURNING *;

-- name: GetEventsByUserId :many
SELECT * FROM events
WHERE user_id = $1;

-- name: GetEventById :one
SELECT * FROM events 
WHERE id = $1;

-- name: UpdateEventDate :one
UPDATE events
SET date = $1
WHERE id = $2
RETURNING *;