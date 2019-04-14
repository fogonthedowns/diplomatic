package model

type PieceRow struct {
	Id        int64     `json:"id"`
	GameId    int64     `json:"game_id"`
	Owner     Country   `json:"owner"`
	UnitType  UnitType  `json:"type"`
	IsActive  bool      `json:"is_active"`
	Country   Territory `json:"location"`
	Dislodged bool      `json:"dislodged"`
}
