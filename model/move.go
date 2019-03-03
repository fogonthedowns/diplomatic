package model

type Move struct {
	Id                      int64     `json:"id"`
	GameId                  int64     `json:"game_id"`
	UserId                  int       `json:"user_id"` // must be hard coded in request, based on logged in user_id
	GameYear                string    `json:"game_year"`
	Phase                   int       `json:"phase"`
	LocationStart           Territory `json:"location_start"`
	LocationSubmitted       Territory `json:"location_submitted"`
	SecondLocationSubmitted Territory `json:"second_location_submitted"`
	LocationResolved        Territory `json:"location_resolved"`
	PieceOwner              Country   `json:"piece_owner"`
	PieceId                 int       `json:"piece_id"`
	OrderType               OrderType `json:"move_type"`
	MovePower               int
}

const (
	BOTH    = MoveType("both")
	SEA     = MoveType("sea")
	LAND    = MoveType("land")
	INVALID = MoveType("invalid")
)

type MoveType string

func (move *Move) MovePieceForward() {
	if move.OrderType == MOVE || move.OrderType == HOLD {
		move.LocationResolved = move.LocationSubmitted
	} else if move.OrderType == SUPPORT {
		move.LocationResolved = move.LocationStart
	}
}

func (move *Move) BouncePiece() {
	move.LocationResolved = move.LocationStart
}
