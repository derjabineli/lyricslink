ALTER TABLE arrangements DROP COLUMN has_chords;
ALTER TABLE arrangements DROP COLUMN has_chord_chart;
ALTER TABLE arrangements ADD COLUMN has_chords BOOLEAN NOT NULL;
ALTER TABLE arrangements ADD COLUMN has_chord_chart BOOLEAN NOT NULL;