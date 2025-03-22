-- name: AddPCSong :one
WITH upsert AS (
    UPDATE songs
    SET updated_at = NOW(),
        title = $1,
        themes = $2,
        copy_right = $3, 
        ccli_number = $4, 
        author = $5, 
        admin = $6
    WHERE pc_id = $7
    RETURNING id
)
INSERT INTO songs (id, created_at, updated_at, title, themes, copy_right, ccli_number, author, admin, pc_id)
SELECT gen_random_uuid(), NOW(), NOW(), $1, $2, $3, $4, $5, $6, $7
WHERE NOT EXISTS (SELECT 1 FROM upsert)
RETURNING id;


-- name: AddPCArrangement :one
WITH upsert AS (
    UPDATE arrangements
    SET updated_at = NOW(),
        name = $1,
        lyrics = $2,
        chord_chart = $3, 
        chord_chart_key = $4, 
        has_chord_chart = $5, 
        has_chords = $6,
        song_id = $8
    WHERE pc_id = $7
    RETURNING id
)
INSERT INTO arrangements (id, created_at, updated_at, name, lyrics, chord_chart, chord_chart_key, has_chord_chart, has_chords, pc_id, song_id)
SELECT gen_random_uuid(), NOW(), NOW(), $1, $2, $3, $4, $5, $6, $7, $8
WHERE NOT EXISTS (SELECT 1 FROM upsert)
RETURNING id;

-- name: CreateUserSongRelation :exec
INSERT INTO users_songs (id, user_id, song_id)
VALUES(gen_random_uuid(), $1, $2);