CREATE TABLE IF NOT EXISTS events_arrangements (
    event_id UUID NOT NULL, 
    arrangement_id UUID NOT NULL, 
    CONSTRAINT fk_event FOREIGN KEY (event_id) 
        REFERENCES events(id) 
        ON DELETE CASCADE, 
    CONSTRAINT fk_arrangement FOREIGN KEY (arrangement_id) 
        REFERENCES arrangements(id) 
        ON DELETE CASCADE 
);