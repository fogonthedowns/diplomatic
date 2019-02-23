package model

type GameSquareData struct {
	Owner Country `json:"owner,omitempty"`
	Units []Unit  `json:"game_squares,omitempty"`
}

type GameBoard map[Territory]GameSquareData
