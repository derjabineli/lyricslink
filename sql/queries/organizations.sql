-- name: CreatePlanningCenterOrganization :one
INSERT INTO planning_center_organizations (id, pc_id, name)
    VALUES(gen_random_uuid(), $1, $2)
ON CONFLICT (pc_id) DO UPDATE 
    SET name = EXCLUDED.name
RETURNING *;

-- name: GetOrganizationByPCId :one
SELECT * FROM planning_center_organizations
WHERE pc_id = $1;

-- name: GetUserOrgRelation :one
SELECT * FROM organizations_users
WHERE user_id = $1;

-- name: CreateUserOrgRelation :exec
INSERT INTO organizations_users (id, user_id, organization_id)
VALUES (gen_random_uuid(), $1, $2)
ON CONFLICT ON CONSTRAINT unique_pc_user_pair DO NOTHING;