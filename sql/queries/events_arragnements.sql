-- name: GetArrangementsWithEventId :many
SELECT 
    es.*, 
    a.name as arrangement_name,
    a.lyrics,
    a.chord_chart
FROM events_songs es
JOIN arrangements a ON es.arrangement_id = a.id
WHERE es.event_id = $1
ORDER BY es.created_at ASC;


-- name: AddArrangementToEvent :one
INSERT INTO events_songs (id, event_id, song_id, arrangement_id)
VALUES (gen_random_uuid(), $1, $2, $3)
RETURNING *;

-- name: GetArrangementsAndSongsWithEventId :many
SELECT 
    s.*,
    a.*
FROM events_songs es
JOIN arrangements a 
    ON a.id =  es.arrangement_id
JOIN songs s
    ON s.id = a.song_id
WHERE es.event_id = $1
ORDER BY es.created_at ASC;

-- name: DeleteEventArrangement :exec
DELETE FROM events_songs WHERE id = $1;

-- name: UpdateArrangement :one
UPDATE events_songs
SET arrangement_id = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;