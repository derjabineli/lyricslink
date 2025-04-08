// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: events_arragnements.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const addArrangementToEvent = `-- name: AddArrangementToEvent :one
INSERT INTO events_arrangements (id, event_id, arrangement_id, created_at, updated_at)
VALUES (gen_random_uuid(), $1, $2, NOW(), NOW())
RETURNING event_id, id, created_at, updated_at, arrangement_id
`

type AddArrangementToEventParams struct {
	EventID       uuid.UUID
	ArrangementID uuid.UUID
}

func (q *Queries) AddArrangementToEvent(ctx context.Context, arg AddArrangementToEventParams) (EventsArrangement, error) {
	row := q.db.QueryRowContext(ctx, addArrangementToEvent, arg.EventID, arg.ArrangementID)
	var i EventsArrangement
	err := row.Scan(
		&i.EventID,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ArrangementID,
	)
	return i, err
}

const deleteEventArrangement = `-- name: DeleteEventArrangement :exec
DELETE FROM events_arrangements WHERE id = $1
`

func (q *Queries) DeleteEventArrangement(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteEventArrangement, id)
	return err
}

const getArrangementsAndSongsWithEventId = `-- name: GetArrangementsAndSongsWithEventId :many
SELECT 
    s.pc_id, s.admin, s.author, s.ccli_number, s.copy_right, s.themes, s.title, s.id, s.created_at, s.updated_at,
    a.name, a.lyrics, a.chord_chart, a.id, a.pc_id, a.chord_chart_key, a.song_id, a.created_at, a.updated_at, a.has_chords, a.has_chord_chart
FROM events_arrangements ea
JOIN arrangements a 
    ON a.id =  ea.arrangement_id
JOIN songs s
    ON s.id = a.song_id
WHERE ea.event_id = $1
ORDER BY ea.created_at ASC
`

type GetArrangementsAndSongsWithEventIdRow struct {
	PcID          sql.NullInt32
	Admin         sql.NullString
	Author        sql.NullString
	CcliNumber    sql.NullInt32
	CopyRight     sql.NullString
	Themes        sql.NullString
	Title         string
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Lyrics        string
	ChordChart    sql.NullString
	ID_2          uuid.UUID
	PcID_2        sql.NullInt32
	ChordChartKey sql.NullString
	SongID        uuid.UUID
	CreatedAt_2   time.Time
	UpdatedAt_2   time.Time
	HasChords     bool
	HasChordChart bool
}

func (q *Queries) GetArrangementsAndSongsWithEventId(ctx context.Context, eventID uuid.UUID) ([]GetArrangementsAndSongsWithEventIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getArrangementsAndSongsWithEventId, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetArrangementsAndSongsWithEventIdRow
	for rows.Next() {
		var i GetArrangementsAndSongsWithEventIdRow
		if err := rows.Scan(
			&i.PcID,
			&i.Admin,
			&i.Author,
			&i.CcliNumber,
			&i.CopyRight,
			&i.Themes,
			&i.Title,
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Lyrics,
			&i.ChordChart,
			&i.ID_2,
			&i.PcID_2,
			&i.ChordChartKey,
			&i.SongID,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
			&i.HasChords,
			&i.HasChordChart,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getArrangementsWithEventId = `-- name: GetArrangementsWithEventId :many
SELECT 
    a.name, a.lyrics, a.chord_chart, a.id, a.pc_id, a.chord_chart_key, a.song_id, a.created_at, a.updated_at, a.has_chords, a.has_chord_chart,
    ea.id AS event_arrangement_id, 
    CASE 
        WHEN a.id = ea.arrangement_id THEN TRUE 
        ELSE FALSE 
    END AS is_selected
FROM events_arrangements ea
JOIN arrangements a 
    ON a.song_id = (
        SELECT song_id FROM arrangements WHERE id = ea.arrangement_id
    )
WHERE ea.event_id = $1
ORDER BY ea.created_at ASC
`

type GetArrangementsWithEventIdRow struct {
	Name               string
	Lyrics             string
	ChordChart         sql.NullString
	ID                 uuid.UUID
	PcID               sql.NullInt32
	ChordChartKey      sql.NullString
	SongID             uuid.UUID
	CreatedAt          time.Time
	UpdatedAt          time.Time
	HasChords          bool
	HasChordChart      bool
	EventArrangementID uuid.UUID
	IsSelected         bool
}

func (q *Queries) GetArrangementsWithEventId(ctx context.Context, eventID uuid.UUID) ([]GetArrangementsWithEventIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getArrangementsWithEventId, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetArrangementsWithEventIdRow
	for rows.Next() {
		var i GetArrangementsWithEventIdRow
		if err := rows.Scan(
			&i.Name,
			&i.Lyrics,
			&i.ChordChart,
			&i.ID,
			&i.PcID,
			&i.ChordChartKey,
			&i.SongID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.HasChords,
			&i.HasChordChart,
			&i.EventArrangementID,
			&i.IsSelected,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateArrangement = `-- name: UpdateArrangement :one
UPDATE events_arrangements
SET arrangement_id = $1
WHERE id = $2
RETURNING event_id, id, created_at, updated_at, arrangement_id
`

type UpdateArrangementParams struct {
	ArrangementID uuid.UUID
	ID            uuid.UUID
}

func (q *Queries) UpdateArrangement(ctx context.Context, arg UpdateArrangementParams) (EventsArrangement, error) {
	row := q.db.QueryRowContext(ctx, updateArrangement, arg.ArrangementID, arg.ID)
	var i EventsArrangement
	err := row.Scan(
		&i.EventID,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ArrangementID,
	)
	return i, err
}
