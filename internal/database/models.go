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
	ID         uuid.UUID
	Name       string
	Lyrics     string
	ChordChart sql.NullString
	SongID     uuid.UUID
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
	ArrangementID uuid.UUID
	ID            uuid.UUID
}

type EventsSong struct {
	EventID uuid.UUID
	SongID  uuid.UUID
}

type Song struct {
	ID         uuid.UUID
	PcID       sql.NullInt32
	Admin      sql.NullString
	Author     sql.NullString
	CcliNumber sql.NullInt32
	CopyRight  sql.NullString
	Themes     sql.NullString
	Title      string
	UserID     uuid.UUID
}

type User struct {
	ID             uuid.UUID
	FirstName      string
	LastName       string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	PcAuthorized   bool
}
