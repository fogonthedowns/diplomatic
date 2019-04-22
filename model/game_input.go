package model

import (
	"time"
)

type GameInput struct {
	Id                int64     `json:"id"`
	Title             string    `json:"title"`
	StartedAt         time.Time `json:"started_at"`
	GameYear          string    `json:"game_year"`
	Phase             int       `json:"phase"`
	PhaseEnd          string    `json:"phase_end"`
	OrdersInterval    int       `json:"orders_interval"`
	GameBoard         GameBoard `json:"game_squares,omitempty"`
	Country           Country   `json:"country"`
	UserId            int       `json:"user_id"` // must be hard coded in request, based on logged in user_id
	LocationStart     Territory `json:"location_start"`
	LocationSubmitted Territory `json:"location_submitted"`
	PieceId           int       `json:"piece_id"`
	MoveType          OrderType `json:"move_type"`
}
