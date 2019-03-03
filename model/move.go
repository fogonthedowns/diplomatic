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
type Moves []*Move

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

// Maps The Turn End Location Territory to each move
// The Key represents where the piece is ordered to, when the turn resolves.
// from the user perspective so The key depends on the order type
func (moves Moves) CategorizeMovesByTerritory() TerritoryMoves {
	tm := make(TerritoryMoves, 0)
	for _, move := range moves {
		moveType := move.LocationStart.ValidMovement(move.LocationSubmitted)
		// NOTE Do not save these modifications - keep these changes in memory
		if moveType == INVALID {
			move.OrderType = HOLD
			move.LocationSubmitted = move.LocationStart
		}

		// Determine if the destination is contested.
		// The contested territory depends on the type of Order
		// this is done by counting either LocationSubmitted or
		// LocationStart
		if move.OrderType == MOVE {
			tm[move.LocationSubmitted] = append(tm[move.LocationSubmitted], move)
		}

		if move.OrderType == SUPPORT {
			tm[move.LocationStart] = append(tm[move.LocationStart], move)
		}

		if move.OrderType == HOLD {
			tm[move.LocationStart] = append(tm[move.LocationStart], move)
		}
	}
	return tm
}
