ALTER TABLE users_songs DROP CONSTRAINT fk_song;
ALTER TABLE users_songs DROP COLUMN song_id;
ALTER TABLE users_songs ADD COLUMN song_id UUID NOT NULL;
ALTER TABLE users_songs 
ADD CONSTRAINT fk_song
FOREIGN KEY (song_id) REFERENCES songs(id);

ALTER TABLE arrangements DROP CONSTRAINT fk_song;
ALTER TABLE arrangements DROP COLUMN song_id;
ALTER TABLE arrangements ADD COLUMN song_id UUID NOT NULL;
ALTER TABLE arrangements 
ADD CONSTRAINT fk_song
FOREIGN KEY (song_id) REFERENCES songs(id);