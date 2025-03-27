-- name: GetArrangementsWithEventId :many
SELECT 
    a.*, 
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
ORDER BY ea.created_at DESC,  is_selected DESC;

-- name: AddArrangementToEvent :one
INSERT INTO events_arrangements (id, event_id, arrangement_id, created_at, updated_at)
VALUES (gen_random_uuid(), $1, $2, NOW(), NOW())
RETURNING *;

-- name: GetArrangementsAndSongsWithEventId :many
SELECT 
    s.*,
    a.*, 
    CASE 
        WHEN a.id = ea.arrangement_id THEN TRUE 
        ELSE FALSE 
    END AS is_selected
FROM events_arrangements ea
JOIN arrangements a 
    ON a.song_id = (
        SELECT song_id FROM arrangements WHERE id = ea.arrangement_id
    )
JOIN songs s
    ON s.id = (
        SELECT id FROM songs WHERE id = a.song_id
    )
WHERE ea.event_id = $1
ORDER BY ea.created_at DESC,  is_selected DESC;