ALTER TABLE users
ADD CONSTRAINT users_pc_id_unique UNIQUE (pc_id);

ALTER TABLE planning_center_organizations
ADD CONSTRAINT organizations_pc_id_unique UNIQUE (pc_id);