-- name: GetSongById :one
SELECT * FROM songs
WHERE id = $1;

-- name: SearchSongs :many
SELECT us.song_id, s.* FROM users_songs us
RIGHT JOIN songs s ON s.id = us.song_id
WHERE us.user_id = $1 AND title LIKE $2;

-- name: GetSongIdByPCId :one
SELECT id FROM songs
WHERE pc_id = $1;