package main

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/derjabineli/lyricslink/internal/database"
	"github.com/google/uuid"
)

// Root struct representing the API response
type PlanningCenterArrangement struct {
	Type          string       								`json:"type"`
	ID            string       								`json:"id"`
	Attributes    PlanningCenterArrangementAttributes   	`json:"attributes"`
	Relationships PlanningCenterArrangementRelationships	`json:"relationships"`
	Links         PlanningCenterArrangementLinks        	`json:"links"`
}

// Attributes holds the properties of the arrangement
type PlanningCenterArrangementAttributes struct {
	ArchivedAt             *string `json:"archived_at"` // Nullable
	BPM                    *int    `json:"bpm"`         // Nullable
	ChordChart             string  `json:"chord_chart"`
	ChordChartChordColor   int     `json:"chord_chart_chord_color"`
	ChordChartColumns      int     `json:"chord_chart_columns"`
	ChordChartFont         string  `json:"chord_chart_font"`
	ChordChartFontSize     int     `json:"chord_chart_font_size"`
	ChordChartKey          string  `json:"chord_chart_key"`
	CreatedAt              string  `json:"created_at"`
	HasChordChart          bool    `json:"has_chord_chart"`
	HasChords              bool    `json:"has_chords"`
	Length                 int     `json:"length"`
	Lyrics                 string  `json:"lyrics"`
	LyricsEnabled          bool    `json:"lyrics_enabled"`
	Meter                  *string `json:"meter"`   // Nullable
	Name                   string  `json:"name"`
	Notes                  *string `json:"notes"`   // Nullable
	NumberChartEnabled     bool    `json:"number_chart_enabled"`
	NumeralChartEnabled    bool    `json:"numeral_chart_enabled"`
	PrintMargin            string  `json:"print_margin"`
	PrintOrientation       string  `json:"print_orientation"`
	PrintPageSize          string  `json:"print_page_size"`
	UpdatedAt              string  `json:"updated_at"`
}

// Relationships define links to related entities
type PlanningCenterArrangementRelationships struct {
	UpdatedBy PlanningCenterArrangementRelationshipData `json:"updated_by"`
	CreatedBy PlanningCenterArrangementRelationshipData `json:"created_by"`
	Song      PlanningCenterArrangementRelationshipData `json:"song"`
}

// RelationshipData wraps the related entity reference
type PlanningCenterArrangementRelationshipData struct {
	Data PlanningCenterArrangementRelationshipInfo `json:"data"`
}

// RelationshipInfo represents a single relationship reference
type PlanningCenterArrangementRelationshipInfo struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Links contains API endpoint URLs
type PlanningCenterArrangementLinks struct {
	Self string `json:"self"`
}

func  (cfg *config) savePCArrangementToDB(arrangement PlanningCenterArrangement, songID uuid.UUID) {
	arrangementID, err := strconv.Atoi(arrangement.ID)
	if err != nil {
		return
	}

	cfg.db.AddPCArrangement(context.Background(), database.AddPCArrangementParams{
		Name: arrangement.Attributes.Name,
    	Lyrics: arrangement.Attributes.Lyrics,
    	ChordChart: sql.NullString{String: arrangement.Attributes.ChordChart, Valid: true},
    	ChordChartKey: sql.NullString{String: arrangement.Attributes.ChordChart, Valid: true},
    	HasChordChart: sql.NullBool{Bool: arrangement.Attributes.HasChordChart, Valid: true},
    	HasChords: sql.NullBool{Bool: arrangement.Attributes.HasChords, Valid: true},
    	PcID: sql.NullInt32{Int32: int32(arrangementID), Valid: true},
    	SongID: songID,
	})
}
