-- name: GetSongById :one
SELECT * FROM songs
WHERE id = $1;

-- name: SearchSongs :many
SELECT * FROM songs
WHERE title LIKE $1
AND user_id = $2;