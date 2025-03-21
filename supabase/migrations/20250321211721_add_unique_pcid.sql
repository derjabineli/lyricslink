CREATE UNIQUE INDEX unique_song_pc_id ON songs(pc_id) WHERE pc_id IS NOT NULL;
CREATE UNIQUE INDEX unique_arrangement_pc_id ON arrangements(pc_id) WHERE pc_id IS NOT NULL;
