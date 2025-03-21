ALTER TABLE events_arrangements DROP CONSTRAINT fk_arrangement;
ALTER TABLE events_arrangements DROP COLUMN arrangement_id;
ALTER TABLE events_arrangements ADD COLUMN arrangement_id UUID NOT NULL;
ALTER TABLE events_arrangements 
ADD CONSTRAINT fk_arrangement
FOREIGN KEY (arrangement_id) REFERENCES arrangements(id);