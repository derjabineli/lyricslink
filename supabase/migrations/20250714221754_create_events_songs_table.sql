CREATE TABLE IF NOT EXISTS events_songs (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL, 
    song_id UUID NOT NULL,
    arrangement_id UUID NOT NULL, 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_event FOREIGN KEY (event_id) 
        REFERENCES events(id) 
        ON DELETE CASCADE, 
    CONSTRAINT fk_song FOREIGN KEY (song_id)
        REFERENCES songs(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_arrangement FOREIGN KEY (arrangement_id) 
        REFERENCES arrangements(id) 
        ON DELETE CASCADE 
);