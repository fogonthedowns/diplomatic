package model

import (
	"time"
)

type Game struct {
	Id             int       `json: "id"`
	Title          string    `json: "title"`
	StartedAt      time.Time `json: "started_at"`
	GameYear       time.Time `json: "game_year"`
	Phase          int       `json: "phase"`
	PhaseEnd       time.Time `json: "phase_end"`
	OrdersInterval int       `json: "orders_interval"`
}
