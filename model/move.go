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
	ConvoyPathMoveIds       []int64
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
	} else if move.OrderType == CONVOY {
		move.LocationResolved = move.LocationStart
	}
}

// I like the way this works better than the above MovePieceForward()
// because it happens last based on Dislodged concept
func (moves *Moves) MoveConvoysForward() {
	for _, move := range *moves {
		if move.OrderType == MOVEVIACONVOY {
			fmt.Printf("move %v \n", move)
			if move.Dislodged {
				move.LocationResolved = move.LocationStart
			} else {
				move.LocationResolved = move.LocationSubmitted
			}
		}
	}
}

func (move *Move) DislodgeIfHold(moves *Moves) {
	if move.OrderType == HOLD || move.OrderType == SUPPORT {
		move.Dislodged = true

	}
	if move.OrderType == CONVOY {
		move.Dislodged = true
		move.ProcessConvoyDislodge(moves)
	}
}

// this is complex as fuck. And won't support a double convoy attack. fuck.
// first loop finds the move ids of the Convoy
// the second ranges over the move id and if the move is dislodged dislodges the entire convoy
// the third sets the Unit being convoy as dislodged, if the convoy is dislodged
func (move *Move) ProcessConvoyDislodge(moves *Moves) {
	var moveIds []int64
	var convoyDislodged bool
	for _, m := range *moves {
		if m.OrderType == MOVEVIACONVOY {
			moveIds = m.ConvoyPathMoveIds
		}
	}

	for _, m := range *moves {
		for _, id := range moveIds {
			if id == m.Id {
				if m.Dislodged {
					convoyDislodged = true
				}

			}
		}
	}

	for _, m := range *moves {
		if m.OrderType == MOVEVIACONVOY {
			if convoyDislodged {
				m.Dislodged = true
			}
		}
	}
}

func (move *Move) BouncePiece() {
	move.LocationResolved = move.LocationStart
}

func (moves *Moves) ProcessMoves() {
	tm := moves.CategorizeMovesByTerritory()
	moves.ResolveUncontestedMoves(tm)
	moves.CalculateSupport()
	tm.ResolveConflicts(moves)
	moves.MoveConvoysForward()

	for _, move := range *moves {
		if move.OrderType == SUPPORT {
			fmt.Printf("******** %vs  %v -> %v from %v. resolved: %+v (%v)\n", move.OrderType, move.LocationSubmitted, move.SecondLocationSubmitted, move.LocationStart, move.LocationResolved, move.MovePower)

		} else {
			fmt.Printf("******** %v (%v -> %v) resolved: %+v (%v)\n", move.OrderType, move.LocationStart, move.LocationSubmitted, move.LocationResolved, move.MovePower)
		}
	}
	fmt.Print("\n\nOrders:\n")
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
	convoyPathTerritories := make([]Territory, 0)
	convoyPathMoveIds := make([]int64, 0)
	allConnections := make(map[Territory][]Territory)
	// TODO Check if this supports multiple convoys, it could bc this is kicked off by a single move.
	for _, move := range moves {
		if move.OrderType == CONVOY && move.LocationSubmitted == begin && move.SecondLocationSubmitted == end {
			convoyPathTerritories = append(convoyPathTerritories, move.LocationStart)
			convoyPathMoveIds = append(convoyPathMoveIds, move.Id)
		}
	}

	for _, move := range moves {
		if move.OrderType == MOVEVIACONVOY {
			move.ConvoyPathMoveIds = convoyPathMoveIds
		}
	}

	for _, orderTerritory := range convoyPathTerritories {
		allConnections[orderTerritory] = validSeaMoves[orderTerritory]
	}

	var beginingConnection, endConnection *Territory

	// Check to see if the Convoy Begin and End order (the Move issued by Army)
	// exists in the Map of allConnections
	for _, convoyTerritory := range convoyPathTerritories {
		for _, t := range allConnections[convoyTerritory] {
			if begin == t {
				//beginValid = true
				beginingConnection = &convoyTerritory
			}
			if end == t {
				//endValid = true
				endConnection = &convoyTerritory

			}
		}

	}

	// Not Valid if there is no connection to the Begin and End
	if beginingConnection == nil {
		return false
	}
	if endConnection == nil {
		return false
	}

	var path bool
	if len(convoyPathTerritories) > 1 {
		// first clean up allConnections, to only relevant info
		// then determine if the path exists
		reducedConnections := reduceAllConnectionsToRelevant(allConnections, convoyPathTerritories)
		path = doesPathExist(reducedConnections)
	} else {
		// since there is already a connection to begin and end from this territory
		return true
	}
	return path
}

// reduceAllConnectionsToRelevant()
// Eliminates useless information in allConnections Terristory map.
// e.g.
// 1 : [2,3,4,5]
// 2 : [12, 3, 45, 1, 8, 11]
// is reduced to:
// 1: [2]
// 2: [1]
// where the Numbers are Territory elements
func reduceAllConnectionsToRelevant(allConnections map[Territory][]Territory, convoyList []Territory) map[Territory][]Territory {
	// Reduce
	reducedConnections := make(map[Territory][]Territory)
	for key, tArray := range allConnections {
		list := make([]Territory, 0)
		for _, territoryElement := range tArray {
			for _, convoyTerritory := range convoyList {
				if territoryElement == convoyTerritory {
					list = append(list, territoryElement)
				}
			}
		}
		reducedConnections[key] = list
	}

	return reducedConnections
}

// Determins if a valid path exists
func doesPathExist(rc map[Territory][]Territory) bool {
	// Path
	seen := make([]Territory, 0)
	for key, list := range rc {
		if len(list) == 0 {
			return false
		}
		seen = appendIfMissing(seen, key)
		for _, t := range list {
			seen = appendIfMissing(seen, t)
		}
	}
	return len(seen) == len(rc)
}

// appendIfMissing() Appends to Territory Slice, if the element is missing
func appendIfMissing(slice []Territory, i Territory) []Territory {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
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
