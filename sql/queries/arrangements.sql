-- name: GetArrangementWithSongId :many
SELECT * FROM arrangements 
WHERE song_id = $1;