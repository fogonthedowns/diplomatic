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

	// convoy rules:
	// use movement of land unit follows ValidMovement()
	// do not ever resolve uncontested convoy rules unless
	// there is a valid path.

	for _, move := range moves {
		var moveType MoveType
		// Vallid support moves are determined by the start location bordering the end location
		// TODO(:3/4/19) pass unit type to ValidMovement() return bool
		// Serious confusing fleets and army units
		// TODO(:3/14/19) Switch on move.OrderType, implement valid movements for Convoy
		// Explore introducing a new type MOVEVIACONVY
		if move.OrderType == SUPPORT {
			moveType = move.LocationStart.ValidMovement(move.SecondLocationSubmitted)
		} else {
			moveType = move.LocationStart.ValidMovement(move.LocationSubmitted)
		}
		// NOTE Do not save these modifications - keep these changes in memory
		if moveType == INVALID {
			move.OrderType = HOLD
			move.LocationSubmitted = move.LocationStart
		}

		// Map Contested Territory to Units moving into Contested Territory
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

// TODO before this determine if support is cut
// This may require a function to -+ the MovePower
// addSupportPointsToMove() This will add up the number of times a unit is supported
func (moves Moves) AddSupportPointsToMove(supportedFrom Territory, supportedTo Territory) {
	for _, move := range moves {
		// if uncontested resolve the move
		// remember above LocationSubmitted was edited in the case of invalid moves in memory
		if move.OrderType == MOVE {
			// if the support order matches the order increment the move power counter
			if move.LocationStart == supportedFrom && move.LocationSubmitted == supportedTo {
				if !moves.CalculateIfSupportIsCut(supportedFrom) {
					move.MovePower += 1
				}
			}
		}
	}
}

// Loop through all moves to determine if there is a Valid attack that cuts support
// Determined by any move - successful or not - to the origin of the support order
func (moves Moves) CalculateIfSupportIsCut(supportedFrom Territory) bool {
	for _, move := range moves {
		if move.OrderType == MOVE && move.LocationStart.ValidMovement(supportedFrom) != INVALID {
			return true
		}
	}
	return false
}
