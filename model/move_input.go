package model

type Move struct {
	Id                int64     `json:"id"`
	GameId            int64     `json:"game_id"`
	UserId            int       `json:"user_id"` // must be hard coded in request, based on logged in user_id
	GameYear          string    `json:"game_year"`
	Phase             int       `json:"phase"`
	LocationStart     Territory `json:"location_start"`
	LocationSubmitted Territory `json:"location_submitted"`
	PieceOwner        Country   `json:"piece_owner"`
	PieceId           int       `json:"piece_id"`
	OrderType         OrderType `json:"move_type"`
}
