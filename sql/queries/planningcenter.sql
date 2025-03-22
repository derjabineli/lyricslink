-- name: SavePCSong :one
INSERT INTO songs (id, created_at, updated_at, title, themes, copy_right, ccli_number, author, admin, pc_id)
VALUES(gen_random_uuid(), NOW(), NOW(), $1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (pc_id) 
DO UPDATE SET
    updated_at = NOW(),
    title = EXCLUDED.title,
    themes = EXCLUDED.themes,
    copy_right = EXCLUDED.copy_right, 
    ccli_number = EXCLUDED.ccli_number, 
    author = EXCLUDED.author, 
    admin = EXCLUDED.admin
RETURNING id;


-- name: SavePCArrangement :one
INSERT INTO arrangements (id, created_at, updated_at, name, lyrics, chord_chart, chord_chart_key, has_chord_chart, has_chords, pc_id, song_id)
VALUES(gen_random_uuid(), NOW(), NOW(), $1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (pc_id) 
DO UPDATE SET
    updated_at = NOW(),
    name = EXCLUDED.name,
    lyrics = EXCLUDED.lyrics,
    chord_chart = EXCLUDED.chord_chart, 
    chord_chart_key = EXCLUDED.chord_chart_key, 
    has_chord_chart = EXCLUDED.has_chord_chart, 
    has_chords = EXCLUDED.has_chords
RETURNING id;

-- name: CreateUserSongRelation :exec
INSERT INTO users_songs (id, user_id, song_id)
VALUES(gen_random_uuid(), $1, $2);