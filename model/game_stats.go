package model

type Stats struct {
	VictoryCenters int `json:"victory_centers,omitempty"`
	UnitCount      int `json:"unit_count,omitempty"`
}

type GameStats map[Country]Stats
