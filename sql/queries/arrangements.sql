-- name: GetArrangementWithSongId :many
SELECT * FROM arrangements 
WHERE song_id = $1;

-- name: GetAvailableArrangements :many
SELECT * FROM arrangements
WHERE song_id = $1 AND id != $2;