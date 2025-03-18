-- name: GetSongById :one
SELECT * FROM songs
WHERE id = $1;