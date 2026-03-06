package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ColumnType represents the type of column a card belongs to.
type ColumnType string

const (
	ColumnWentWell    ColumnType = "went_well"
	ColumnDidntGoWell ColumnType = "didnt_go_well"
	ColumnPuzzling    ColumnType = "puzzling"
	ColumnActionItem  ColumnType = "action_item"
)

var ValidColumnTypes = map[ColumnType]bool{
	ColumnWentWell:    true,
	ColumnDidntGoWell: true,
	ColumnPuzzling:    true,
	ColumnActionItem:  true,
}

// --- Request types ---

type CreateRetroRequest struct {
	Name string `json:"name" binding:"required" validate:"required"`
}

type CreateCardRequest struct {
	ColumnType ColumnType `json:"column_type" binding:"required" validate:"required"`
	Content    string     `json:"content" binding:"required" validate:"required"`
	AuthorName string     `json:"author_name"`
}

// --- Response types ---

type RetroResponse struct {
	ID        uuid.UUID      `json:"id" validate:"required"`
	Name      string         `json:"name" validate:"required"`
	Cards     []CardResponse `json:"cards" validate:"required"`
	CreatedAt time.Time      `json:"created_at" validate:"required"`
}

type CardResponse struct {
	ID         uuid.UUID  `json:"id" validate:"required"`
	ColumnType ColumnType `json:"column_type" validate:"required"`
	Content    string     `json:"content" validate:"required"`
	AuthorName string     `json:"author_name"`
	CreatedAt  time.Time  `json:"created_at" validate:"required"`
}

type RetroListItem struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}

type AnalysisResponse struct {
	ID        uuid.UUID       `json:"id" validate:"required"`
	RetroID   uuid.UUID       `json:"retro_id" validate:"required"`
	Result    json.RawMessage `json:"result" validate:"required"`
	CreatedAt time.Time       `json:"created_at" validate:"required"`
}

// --- LLM structured output types ---

type AnalysisResult struct {
	Groups         []AnalysisGroup `json:"groups"`
	OverallSummary string          `json:"overall_summary"`
}

type AnalysisGroup struct {
	Theme              string   `json:"theme"`
	CardIDs            []string `json:"card_ids"`
	Synthesis          string   `json:"synthesis"`
	FiveWhys           []string `json:"five_whys"`
	HypothesisedCauses []string `json:"hypothesised_causes"`
}
