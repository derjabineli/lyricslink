-- name: GetArrangementsWithEventId :many
SELECT a.*
FROM events_arrangements ea
JOIN arrangements a 
    ON a.id = ea.arrangement_id
WHERE ea.event_id = $1;
