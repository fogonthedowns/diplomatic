package model

type GameSquareData struct {
	Owner       Country `json:"owner,omitempty"`
	Units       []Unit  `json:"game_squares,omitempty"`
	TerritoryId int     `json:"territory_id"`
}

type GameBoard map[Territory]GameSquareData
