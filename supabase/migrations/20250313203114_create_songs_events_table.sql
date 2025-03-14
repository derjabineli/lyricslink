CREATE TABLE IF NOT EXISTS events_songs ( 
    event_id UUID NOT NULL, 
    song_id UUID NOT NULL, 
    CONSTRAINT fk_event FOREIGN KEY (event_id) 
        REFERENCES events(id) 
        ON DELETE CASCADE, 
    CONSTRAINT fk_song FOREIGN KEY (song_id) 
        REFERENCES songs(id) 
        ON DELETE CASCADE 
);