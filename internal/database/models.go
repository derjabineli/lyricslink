// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Arrangement struct {
	Name          string
	Lyrics        string
	ChordChart    sql.NullString
	ID            uuid.UUID
	PcID          sql.NullInt32
	ChordChartKey sql.NullString
	SongID        uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	HasChords     bool
	HasChordChart bool
}

type Event struct {
	ID        uuid.UUID
	Name      string
	Date      time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
}

type EventsArrangement struct {
	EventID       uuid.UUID
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ArrangementID uuid.UUID
}

type OrganizationsSong struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	SongID         uuid.UUID
}

type OrganizationsUser struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
}

type PlanningCenterOrganization struct {
	ID   uuid.UUID
	Name string
	PcID string
}

type PlanningCenterToken struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	AccessToken  string
	TokenType    string
	ExpiresIn    int32
	RefreshToken string
	Scope        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Revoked      bool
}

type Song struct {
	PcID       sql.NullInt32
	Admin      sql.NullString
	Author     sql.NullString
	CcliNumber sql.NullInt32
	CopyRight  sql.NullString
	Themes     sql.NullString
	Title      string
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type User struct {
	ID            uuid.UUID
	FirstName     string
	LastName      string
	Email         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Administrator bool
	Avatar        string
	PcID          string
}

type UserSession struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	AccessToken  string
	TokenType    string
	ExpiresIn    int32
	RefreshToken string
	Scope        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Revoked      bool
}

type UsersSong struct {
	ID     uuid.UUID
	UserID uuid.UUID
	SongID uuid.UUID
}
