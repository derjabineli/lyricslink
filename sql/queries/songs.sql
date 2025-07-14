-- name: GetSongById :one
SELECT * FROM songs
WHERE id = $1;

-- name: SearchSongs :many
SELECT os.song_id, s.* FROM organizations_songs os
RIGHT JOIN songs s ON s.id = os.song_id
WHERE os.organization_id = $1 AND title LIKE $2;

-- name: GetSongIdByPCId :one
SELECT id FROM songs
WHERE pc_id = $1;

-- name: AddSong :one
INSERT INTO songs (id, created_at, updated_at, title, themes, copy_right, ccli_number, author, admin, pc_id)
    VALUES(gen_random_uuid(), NOW(), NOW(), $1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (pc_id) DO UPDATE 
    SET title = EXCLUDED.title,
        themes = EXCLUDED.themes,
        copy_right = EXCLUDED.copy_right,
        ccli_number = EXCLUDED.ccli_number,
        author = EXCLUDED.author,
        admin = EXCLUDED.admin,
        updated_at = NOW()
RETURNING id;
