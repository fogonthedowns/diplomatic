package model

import (
	"time"
)

type GameInput struct {
	Id             int64     `json: "id"`
	Title          string    `json: "title"`
	StartedAt      time.Time `json: "started_at"`
	GameYear       time.Time `json: "game_year"`
	Phase          int       `json: "phase"`
	PhaseEnd       time.Time `json: "phase_end"`
	OrdersInterval int       `json: "orders_interval"`
	GameBoard      GameBoard `json: "game_squares,omitempty"`
	Country        Country   `json:"country"`
	UserId         int       `json:"user_id"`
}