ALTER TABLE events_arrangements DROP CONSTRAINT fk_arrangement;
ALTER TABLE events_arrangements DROP COLUMN arrangement_id;

ALTER TABLE arrangements DROP CONSTRAINT arrangements_pkey;
ALTER TABLE arrangements DROP COLUMN id;
ALTER TABLE arrangements ADD COLUMN id UUID PRIMARY KEY;

ALTER TABLE events_arrangements ADD COLUMN arrangement_id UUID;
ALTER TABLE events_arrangements 
ADD CONSTRAINT fk_arrangement
FOREIGN KEY (arrangement_id) REFERENCES arrangements(id);