CREATE TABLE IF NOT EXISTS organizations_users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT fk_organization
        FOREIGN KEY (organization_id) REFERENCES planning_center_organizations(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE
)