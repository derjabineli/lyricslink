ALTER TABLE organizations_users
ADD CONSTRAINT unique_pc_user_pair UNIQUE (organization_id, user_id);