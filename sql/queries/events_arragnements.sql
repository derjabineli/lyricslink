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
ORDER BY is_selected DESC;