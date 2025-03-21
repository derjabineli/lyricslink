ALTER TABLE events_arrangements DROP CONSTRAINT fk_arrangement;

ALTER TABLE events_arrangements DROP COLUMN arrangement_id;

ALTER TABLE events_arrangements ADD COLUMN arrangement_id INTEGER NOT NULL;

ALTER TABLE arrangements DROP CONSTRAINT arrangements_pkey;

ALTER TABLE arrangements DROP COLUMN id;

ALTER TABLE arrangements ADD COLUMN id INTEGER NOT NULL;

ALTER TABLE arrangements ADD PRIMARY KEY (id);

ALTER TABLE events_arrangements 
ADD CONSTRAINT fk_arrangement 
FOREIGN KEY (arrangement_id) REFERENCES arrangements(id);