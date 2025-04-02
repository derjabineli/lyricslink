-- name: GetArrangementsWithEventId :many
SELECT 
    a.*,
    ea.id AS event_arrangement_id, 
    CASE 
        WHEN a.id = ea.arrangement_id THEN TRUE 
        ELSE FALSE 
    END AS is_selected
FROM events_arrangements ea
JOIN arrangements a 
    ON a.song_id = (
        SELECT song_id FROM arrangements WHERE id = ea.arrangement_id
    )
WHERE ea.event_id = $1
ORDER BY ea.created_at ASC;

-- name: AddArrangementToEvent :one
INSERT INTO events_arrangements (id, event_id, arrangement_id, created_at, updated_at)
VALUES (gen_random_uuid(), $1, $2, NOW(), NOW())
RETURNING *;

-- name: GetArrangementsAndSongsWithEventId :many
SELECT 
    s.*,
    a.*
FROM events_arrangements ea
JOIN arrangements a 
    ON a.id =  ea.arrangement_id
JOIN songs s
    ON s.id = a.song_id
WHERE ea.event_id = $1
ORDER BY ea.created_at ASC;

-- name: DeleteEventArrangement :exec
DELETE FROM events_arrangements WHERE id = $1;