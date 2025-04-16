// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: organizations.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createPlanningCenterOrganization = `-- name: CreatePlanningCenterOrganization :one
INSERT INTO planning_center_organizations (id, pc_id, name)
    VALUES(gen_random_uuid(), $1, $2)
ON CONFLICT (pc_id) DO UPDATE 
    SET name = EXCLUDED.name
RETURNING id, name, pc_id
`

type CreatePlanningCenterOrganizationParams struct {
	PcID string
	Name string
}

func (q *Queries) CreatePlanningCenterOrganization(ctx context.Context, arg CreatePlanningCenterOrganizationParams) (PlanningCenterOrganization, error) {
	row := q.db.QueryRowContext(ctx, createPlanningCenterOrganization, arg.PcID, arg.Name)
	var i PlanningCenterOrganization
	err := row.Scan(&i.ID, &i.Name, &i.PcID)
	return i, err
}

const createUserOrgRelation = `-- name: CreateUserOrgRelation :exec
INSERT INTO organizations_users (id, user_id, organization_id)
VALUES (gen_random_uuid(), $1, $2)
ON CONFLICT ON CONSTRAINT unique_pc_user_pair DO NOTHING
`

type CreateUserOrgRelationParams struct {
	UserID         uuid.UUID
	OrganizationID uuid.UUID
}

func (q *Queries) CreateUserOrgRelation(ctx context.Context, arg CreateUserOrgRelationParams) error {
	_, err := q.db.ExecContext(ctx, createUserOrgRelation, arg.UserID, arg.OrganizationID)
	return err
}

const getOrganizationByPCId = `-- name: GetOrganizationByPCId :one
SELECT id, name, pc_id FROM planning_center_organizations
WHERE pc_id = $1
`

func (q *Queries) GetOrganizationByPCId(ctx context.Context, pcID string) (PlanningCenterOrganization, error) {
	row := q.db.QueryRowContext(ctx, getOrganizationByPCId, pcID)
	var i PlanningCenterOrganization
	err := row.Scan(&i.ID, &i.Name, &i.PcID)
	return i, err
}

const getUserOrgRelation = `-- name: GetUserOrgRelation :one
SELECT id, organization_id, user_id FROM organizations_users
WHERE user_id = $1
`

func (q *Queries) GetUserOrgRelation(ctx context.Context, userID uuid.UUID) (OrganizationsUser, error) {
	row := q.db.QueryRowContext(ctx, getUserOrgRelation, userID)
	var i OrganizationsUser
	err := row.Scan(&i.ID, &i.OrganizationID, &i.UserID)
	return i, err
}
