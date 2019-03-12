package model

import "fmt"

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
	UnitType                UnitType
	Dislodged               bool `json:"dislodged"`
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
	if move.OrderType == MOVE {
		move.LocationResolved = move.LocationSubmitted
	} else if move.OrderType == HOLD {
		move.LocationResolved = move.LocationSubmitted
	} else if move.OrderType == SUPPORT {
		move.LocationResolved = move.LocationStart
	} else if move.OrderType == MOVEVIACONVOY {
		move.LocationResolved = move.LocationSubmitted
	} else if move.OrderType == CONVOY {
		move.LocationResolved = move.LocationStart
	}
}

func (move *Move) DislodgeIfHold() {
	if move.OrderType == HOLD || move.OrderType == SUPPORT {
		move.Dislodged = true
	}
}

func (move *Move) BouncePiece() {
	move.LocationResolved = move.LocationStart
}

func (moves *Moves) ProcessMoves() {
	tm := moves.CategorizeMovesByTerritory()
	moves.ResolveUncontestedMoves(tm)
	moves.CalculateSupport()
	tm.ResolveConflicts()

	for _, move := range *moves {
		if move.OrderType == SUPPORT {
			fmt.Printf("******** %vs  %v -> %v from %v. resolved: %+v (%v)\n", move.OrderType, move.LocationSubmitted, move.SecondLocationSubmitted, move.LocationStart, move.LocationResolved, move.MovePower)

		} else {
			fmt.Printf("******** %v (%v -> %v) resolved: %+v (%v)\n", move.OrderType, move.LocationStart, move.LocationSubmitted, move.LocationResolved, move.MovePower)
		}
	}
}

// Maps The Turn End Location Territory to each move
// The Key represents where the piece is ordered to, when the turn resolves.
// from the user perspective so The key depends on the order type
func (moves *Moves) CategorizeMovesByTerritory() TerritoryMoves {
	tm := make(TerritoryMoves, 0)

	for _, move := range *moves {
		var valid bool
		// Vallid support moves are determined by the start location bordering the end location
		// TODO(:3/12/19) Refactor ValidMovement to include ConvoyPathDoesExist()

		if move.OrderType == MOVEVIACONVOY {
			valid = moves.ConvoyPathDoesExist(move.LocationStart, move.LocationSubmitted)
		} else {
			valid = move.LocationStart.ValidMovement(*move)
		}

		// TODO(:3/12/19)
		// Message concept to indicate move coerced to Move
		// Consider Historical move.HistoricalOrder to show past moves
		if !valid {
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

		if move.OrderType == MOVEVIACONVOY {
			tm[move.LocationStart] = append(tm[move.LocationSubmitted], move)
		}

		if move.OrderType == CONVOY {
			tm[move.LocationStart] = append(tm[move.LocationStart], move)
		}

	}
	return tm
}

// TODO before this determine if support is cut
// This may require a function to -+ the MovePower
// addSupportPointsToMove() This will add up the number of times a unit is supported
func (moves Moves) CalculateSupport() {
	for index, move := range moves {
		if move.OrderType == SUPPORT {
			moves.AddSupportPointsToMove(*move)
			moves[index].MovePieceForward()
		}
	}
}

// ConvoyPathDoesExist() loops through all the moves if the move is a Convoy then check the begining and end move
// sent by the Army which is being convoyed. Build up a slice of Territories, and then
// build up a slice of Neighbor Territory's determine if the path
// from begining to end exists from that slice.
func (moves Moves) ConvoyPathDoesExist(begin Territory, end Territory) bool {
	fmt.Printf("****** %v :: %v \n", begin, end)
	convoyPathTerritories := make([]Territory, 0)
	allConnections := make(map[Territory][]Territory)
	// TODO this does not support multiple convoys
	for _, move := range moves {
		if move.OrderType == CONVOY && move.LocationSubmitted == begin && move.SecondLocationSubmitted == end {
			convoyPathTerritories = append(convoyPathTerritories, move.LocationStart)
		}
	}

	for _, orderTerritory := range convoyPathTerritories {
		allConnections[orderTerritory] = validSeaMoves[orderTerritory]
	}

	var beginingConnection, endConnection Territory
	// fmt.Printf("%+v\n", convoyPathTerritories)
	//var beginValid, endValid bool
	for _, convoyTerritory := range convoyPathTerritories {
		for _, t := range allConnections[convoyTerritory] {
			if begin == t {
				//beginValid = true
				beginingConnection = convoyTerritory
			}
			if end == t {
				//endValid = true
				endConnection = convoyTerritory

			}
		}

	}

	// fmt.Printf("%+v\n", beginValid)
	// fmt.Printf("%+v\n", endValid)
	// fmt.Printf("%+v\n", beginingConnection)
	// fmt.Printf("%+v\n", endConnection)

	var path bool
	if len(convoyPathTerritories) > 1 {
		path = DoesBeginToEndPathExist(beginingConnection, endConnection, allConnections)
		fmt.Printf("%v\n", path)
	} else {
		path = DoesBeginToEndPathExist(beginingConnection, end, allConnections)
	}
	fmt.Printf("path %v \n", path)
	return path
}

func DoesBeginToEndPathExist(begin Territory, end Territory, allConnections map[Territory][]Territory) bool {
	// if begin != nil && end != nil {
	fmt.Printf("begin %v \n", begin)
	fmt.Printf("end %v \n", end)
	fmt.Printf("connections %v \n", allConnections)
	for _, territoriesConnectedToBegin := range allConnections[begin] {
		fmt.Printf("territoriesConnectedToBegin %v \n", territoriesConnectedToBegin)

		// TODO why is this a rune?
		// fmt.Printf("%v : %v \n", Territory(territory), end)
		if territoriesConnectedToBegin == end {
			return true
		}

	}
	// }
	return false
}

func (moves Moves) ResolveUncontestedMoves(tm TerritoryMoves) {
	// Resolve Moves
	for index, move := range moves {
		// if uncontested resolve the move
		// remember above LocationSubmitted was edited in the case of invalid moves in memory

		// convoy rules:
		// use movement of land unit follows ValidMovement()
		// do not ever resolve uncontested convoy rules unless
		// there is a valid path.

		// TODO (3/4/19) Uncontested should return false in the case of convoy
		// TODO (3/4/19) Implement a seperate fun to check convoy path and if it is uncontested.
		if tm.Uncontested(move.LocationSubmitted) {
			moves[index].MovePieceForward()
		}
	}
}

func (moves Moves) AddSupportPointsToMove(supportMove Move) {
	// if uncontested resolve the move
	// remember above LocationSubmitted was edited in the case of invalid moves in memory
	for idx, move := range moves {
		if move.OrderType == MOVE || move.OrderType == HOLD {
			// if the support order matches the order increment the move power counter
			if move.LocationStart == supportMove.LocationSubmitted && move.LocationSubmitted == supportMove.SecondLocationSubmitted {
				if !moves.CalculateIfSupportIsCut(supportMove) {
					moves[idx].MovePower += 1
				}
			}
		}
	}
}

// Loop through all moves to determine if there is a Valid attack that cuts support
// Determined by any move - successful or not - to the origin of the support order
func (moves Moves) CalculateIfSupportIsCut(originOfSupportOrder Move) (cut bool) {
	cut = false
	for _, move := range moves {
		// loop through moves, if the submitted move matches originOfSupportOrder
		// check the submitted moves Validity (from its LocationStart)
		if move.OrderType == MOVE && move.LocationSubmitted == originOfSupportOrder.LocationStart {
			// originOfSupportOrder.LocationStart, originOfSupportOrder.UnitType
			cut = move.LocationStart.ValidMovement(originOfSupportOrder)
		}
	}
	return cut
}
