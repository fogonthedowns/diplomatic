package model

import (
	"errors"
	"fmt"
)

type Move struct {
	Id                      int64     `json:"id"`
	GameId                  int64     `json:"game_id"`
	UserId                  int64     `json:"user_id"` // must be hard coded in request, based on logged in user_id
	GameYear                string    `json:"game_year"`
	Phase                   GamePhase `json:"phase"`
	LocationStart           Territory `json:"location_start"`
	LocationSubmitted       Territory `json:"location_submitted"`
	SecondLocationSubmitted Territory `json:"second_location_submitted"`
	LocationResolved        Territory `json:"location_resolved"`
	PieceOwner              Country   `json:"piece_owner"`
	PieceId                 int64     `json:"piece_id"`
	OrderType               OrderType `json:"move_type"`
	SupportUnitType         UnitType  `json:"support_unit_type"`
	MovePower               int
	UnitType                UnitType
	Dislodged               bool `json:"dislodged"`
	DislodgedFrom           Territory
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

// This is the Starting Point.
// TODO moves.save()
// TODO prevent a country from attacking itself - units will bounce, support will not be cut (verify)
// TODO location_resolved is unsaved!
func (moves *Moves) ProcessMoves() {
	tm := moves.CategorizeMovesByTerritory()
	moves.ResolveUncontestedMoves(tm)
	moves.CalculateSupport()
	tm.ResolveConflicts(moves)
	moves.MoveConvoysForward()
	moves.logProgress()
}

// returns Moves that do not exist in the current moves
// This is based on a set of active pieces that are not included in moves
// The order is set to Hold
func (moves *Moves) HoldUnmovedPieces(pieces []*PieceRow) Moves {
	// pices [p.1,p.2,p.3,p.4]
	// moves [m.1, m.2]
	m := make(map[int64]bool, 0)
	var gameYear string
	for _, move := range *moves {
		gameYear = move.GameYear
		m[move.PieceId] = true
	}

	newMoves := make(Moves, 0)
	for _, row := range pieces {
		if m[row.Id] != true && row.IsActive == true {
			newMove := &Move{
				PieceId:           row.Id,
				LocationStart:     row.Country,
				LocationSubmitted: row.Country,
				OrderType:         HOLD,
				PieceOwner:        row.Owner,
				UnitType:          row.UnitType,
				GameYear:          gameYear,
			}
			newMoves = append(newMoves, newMove)
		}
	}
	return newMoves
}

func (moves *Moves) logProgress() {
	for _, move := range *moves {
		if move.OrderType == SUPPORT {
			fmt.Printf("********%v %vs  %v -> %v from %v. resolved: %+v (%v)\n", move.PieceId, move.OrderType, move.LocationSubmitted, move.SecondLocationSubmitted, move.LocationStart, move.LocationResolved, move.MovePower)
		} else {
			fmt.Printf("********%v %v (%v -> %v) resolved: %+v (%v)\n", move.PieceId, move.OrderType, move.LocationStart, move.LocationSubmitted, move.LocationResolved, move.MovePower)
		}
	}
	fmt.Print("\n\nOrders:\n")
}

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
			if move.Dislodged {
				move.LocationResolved = move.LocationStart
			} else {
				move.LocationResolved = move.LocationSubmitted
			}
		}
	}
}

// This function is responsible for Setting move.Dislodged
// The Diplomacy Rules state that rules are adjudicated simultaniously
// A successful moving attack piece can not be dislodged
func (move *Move) DislodgeIfHold(moves *Moves) {
	if move.OrderType == HOLD || move.OrderType == SUPPORT {
		move.Dislodged = true

	}
	if move.OrderType == CONVOY {
		move.Dislodged = true
		move.ProcessConvoyDislodge(moves)
	}
}

// this is complex as fuck. It accepts a Dislodged Convoy move and all moves.
// first loop finds the Convoy move ids from the Army being convoy.
// the second ranges over the Convoy move ids and if the convoy is dislodged it sets convoyDislodged true
// the third sets the Unit being convoy as dislodged, if the convoy is dislodged
// Example input:
// Move info:Convoy ALB : [2]
// Move info:Convoy LON : [0 4 5]
// Move info:Convoy LON : [0 4 5]
func (dislodgedConvoyMove *Move) ProcessConvoyDislodge(moves *Moves) {
	var moveIds []int64
	var convoyDislodged bool

	for _, m := range *moves {
		if m.OrderType == MOVEVIACONVOY && dislodgedConvoyMove.SecondLocationSubmitted == m.LocationSubmitted {
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

// Maps The Turn End Location Territory to each move
// The Key represents where the piece is ordered to, when the turn resolves.
// from the user perspective so The key depends on the order type
func (moves *Moves) CategorizeMovesByTerritory() TerritoryMoves {
	tm := make(TerritoryMoves, 0)

	for _, move := range *moves {
		var valid bool
		// Vallid support moves are determined by the start location bordering the end location
		valid = moves.ValidMovement(*move)

		// TODO pass game to this func
		// if the phase is retreat, check if unit is dislodged
		// if not issue HOLD order

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
			if move.LocationSubmitted.SpainEdgeCase() {
				tm[SPAIN] = append(tm[SPAIN], move)
			} else if move.LocationSubmitted.BulgariaEdgeCase() {
				tm[BULGARIA] = append(tm[BULGARIA], move)
			} else if move.LocationSubmitted.RussiaEdgeCase() {
				tm[ST_PETERSBURG] = append(tm[ST_PETERSBURG], move)
			} else {
				tm[move.LocationSubmitted] = append(tm[move.LocationSubmitted], move)
			}
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

// moves is only used to loop through moves to determine if this move is attacking yourself
// and to determine if the convoy path exists.
func (moves *Moves) ValidMovement(move Move) bool {
	t := move.LocationStart
	var check Territory

	switch move.OrderType {
	case SUPPORT:
		check = move.SecondLocationSubmitted
	case CONVOY:
		check = move.SecondLocationSubmitted
	case MOVE:
		check = move.LocationSubmitted
	case MOVEVIACONVOY:
		check = move.LocationSubmitted
	case RETREAT:
		check = move.LocationSubmitted
	default:
		check = move.LocationSubmitted
	}

	// attacking yourself is invalid
	// TODO return a message
	// false, invalid move
	if moves.attackingYourSelf(check, move) {
		return false
	}

	switch move.UnitType {
	case ARMY:
		if move.OrderType == MOVEVIACONVOY {
			return t.ValidConvoyBeginAndEnd(check) && moves.ConvoyPathDoesExist(move.LocationStart, move.LocationSubmitted)
		} else if move.OrderType == RETREAT {
			return t.ValidLandMovement(check, RETREAT) && check != move.DislodgedFrom
		} else {
			return t.ValidLandMovement(check, move.OrderType)
		}
	case NAVY:
		if move.OrderType == CONVOY {
			return move.LocationSubmitted.ValidConvoyBeginAndEnd(check)
		} else {
			return t.ValidSeaMovement(check)
		}
	default:
		return false
	}
}

// attackingYourSelf determins if you are attacking your own piece
func (moves Moves) attackingYourSelf(destinationTerritory Territory, move Move) bool {
	// check the Piece Owner of the destination
	// does not match the Piece Owner of the Deti
	for _, m := range moves {
		switch move.OrderType {
		// SUPPORT HOLD OF YOUR SELF is valid
		case SUPPORT:
			return false
		default:
			if destinationTerritory == m.LocationStart && move.PieceOwner == m.PieceOwner {
				return true
			}
		}

	}
	return false
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
		if move.OrderType == MOVEVIACONVOY && move.LocationSubmitted == end {
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
			if ValidSupportOrder(move, supportMove) {
				if !moves.CalculateIfSupportIsCut(supportMove) {
					moves[idx].MovePower += 1
				}
			}
		}
	}
}

func ValidSupportOrder(move *Move, supportMove Move) bool {
	return move.LocationStart == supportMove.LocationSubmitted && move.LocationSubmitted == supportMove.SecondLocationSubmitted
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
			cut = moves.ValidMovement(*move)
		}
	}
	return cut
}

func (move *Move) ValidateCountry(gameUser *GameUser) (err error) {
	if move.PieceOwner != gameUser.Country {
		return errors.New("The User does not control this country")
	}
	return err
}
