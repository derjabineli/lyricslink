CREATE TABLE IF NOT EXISTS arrangements (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    lyrics TEXT NOT NULL,
    chord_chart TEXT,
    song_id UUID NOT NULL,
    CONSTRAINT fk_song FOREIGN KEY (song_id)
        REFERENCES songs(id)
        ON DELETE CASCADE
);