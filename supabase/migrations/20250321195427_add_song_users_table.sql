CREATE TABLE IF NOT EXISTS users_songs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL, 
    song_id INTEGER NOT NULL, 
    CONSTRAINT fk_user FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE, 
    CONSTRAINT fk_song FOREIGN KEY (song_id) 
        REFERENCES songs(id) 
        ON DELETE CASCADE 
);