ALTER TABLE users 
ADD COLUMN organization_id TEXT NOT NULL,
ADD COLUMN administrator BOOLEAN NOT NULL,
DROP COLUMN pc_authorized;