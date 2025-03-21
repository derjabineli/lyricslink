-- name: GetSongById :one
SELECT * FROM songs
WHERE id = $1;

-- name: SearchSongs :many
SELECT us.song_id, s.* FROM users_songs us
RIGHT JOIN songs s ON songs(id) = song_id
WHERE users_songs(user_id) = $1 AND title LIKE $2;