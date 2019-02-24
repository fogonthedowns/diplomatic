package model

type TerritoryRow struct {
	Id      int       `json: "id"`
	GameId  int64     `json: "game_id"`
	Owner   Country   `json: "owner"`
	Country Territory `json: "country"`
}
