CREATE TABLE IF NOT EXISTS planning_center_organizations(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pc_id INTEGER NOT NULL,
    name TEXT NOT NULL
);