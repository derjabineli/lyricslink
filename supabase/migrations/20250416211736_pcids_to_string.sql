ALTER TABLE users DROP COLUMN pc_id;
ALTER TABLE users ADD COLUMN pc_id TEXT NOT NULL;

ALTER TABLE planning_center_organizations DROP COLUMN pc_id;
ALTER TABLE planning_center_organizations ADD COLUMN pc_id TEXT NOT NULL;