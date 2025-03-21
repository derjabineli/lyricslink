ALTER TABLE users_songs ADD COLUMN song_id UUID;
ALTER TABLE users_songs 
ADD CONSTRAINT fk_song 
FOREIGN KEY (song_id) REFERENCES songs(id);

ALTER TABLE arrangements ADD COLUMN song_id UUID;
ALTER TABLE arrangements 
ADD CONSTRAINT fk_song 
FOREIGN KEY (song_id) REFERENCES songs(id);