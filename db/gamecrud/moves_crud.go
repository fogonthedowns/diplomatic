package gamecrud

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	model "github.com/fogonthedowns/diplomatic/model"
)

// NewSQLPostRepo retunrs implement of game db interface
func NewMovesEngine(Conn *sql.DB) MovesEngine {
	me := MovesEngine{
		Conn: Conn,
	}
	return me
}

type MovesEngine struct {
	Conn *sql.DB
}

// This could be update or create but this conforms to the interface
// A Player can send moves
// Validate a player is part of the game
// Validate they own the country
// Validate time/phase
// Count the pieces_moves table, when the records are complete or when time expires update the pieces_moves.location_resolved
// Update the game phase, year and phase_end based on the orders_interval
// todi introduce finalize moves

func (e *MovesEngine) CreateOrUpdate(ctx context.Context, in *model.Move) (int64, error) {
	query := "Select user_id, game_id, country from users_games where user_id=? and game_id=?"
	gameQuery := "Select id, game_year, phase, phase_end, title From games where id=?"
	doesPieceMoveExist := "select id from moves where game_id=? AND phase=? AND piece_id=? AND game_year=? LIMIT 1"
	pieceMoveInsert := "Insert moves SET location_start=?, location_submitted=?, second_location_submitted=?, phase=?, game_id=?, type=?, piece_owner=?, game_year=?, piece_id=?"
	pieceMoveUpdate := "Update moves SET location_start=?, location_submitted=?, second_location_submitted=?, phase=?, game_id=?, type=?, piece_owner=?, game_year=? WHERE piece_id=?"

	gameUser, err := e.fetchGameUser(ctx, query, in.UserId, in.GameId)
	fmt.Printf("GAME ID? %+v", in)
	// Is the User part of this game?
	if gameUser == nil {
		return 403, errors.New("The User is not a member of this game")
	}

	if err != nil {
		fmt.Printf("error fetchGameUser(): %v \n", err)
		return 500, err
	}

	game, err := e.fetchGame(ctx, gameQuery, in.GameId)
	if game == nil {
		return 500, errors.New("The Game can not be loaded")
	}

	// Does your User control the piece you are trying to move?
	err = in.ValidateCountry(gameUser)
	if err != nil {
		return 403, err
	}

	// Has the game started?
	// Has the phase ended?
	err = game.ValidPhase()

	if err != nil {
		return 403, err
	}

	if in.Phase != game.Phase {
		return 400, errors.New("The phase is invalid")
	}

	if in.GameYear != game.GameYear {
		return 400, errors.New("The year is invalid")
	}

	if game.Phase == model.FallBuild {
		return e.CreateOrUpdateBuildPhaseMove(ctx, in, *game)
	}

	// Does the submitted move exist for this piece?
	if in.PieceId == 0 {
		return 400, errors.New("must include piece id")
	}

	err = e.ValidPiece(ctx, in.PieceId, in.LocationStart, in.GameId, in.PieceOwner, game.Phase)
	if err != nil {
		return 400, err
	}

	// Note, this only checks to see if the piece_id has created an order
	// it does not verify the order is at all valid
	// this accepts location_start where the piece.location does not match
	// to accomindate the game design requirement to accept invalid orders
	// TODO in ProcessMoves() add a lookup to validate piece.location == move.location_start
	move, err := e.fetchMove(ctx, doesPieceMoveExist, in.GameId, game.Phase, in.PieceId, game.GameYear)
	if err != nil {
		return 500, err
	}

	// if the submitted move does not exist (based on piece id) create it; otherwise update it
	var insertType string
	switch move {
	case nil:
		insertType = pieceMoveInsert
	default:
		insertType = pieceMoveUpdate
	}
	stmt, err := e.Conn.PrepareContext(ctx, insertType)

	if err != nil {
		return 500, err
	}
	_, err = stmt.ExecContext(ctx, in.LocationStart, in.LocationSubmitted, in.SecondLocationSubmitted, in.Phase, in.GameId, in.OrderType, in.PieceOwner, in.GameYear, in.PieceId)
	defer stmt.Close()
	if err != nil {
		return 500, err
	}

	// return 0
	return 200, err
}

func (e *MovesEngine) CreateOrUpdateBuildPhaseMove(ctx context.Context, in *model.Move, game model.Game) (int64, error) {
	doesPieceMoveExist := "select id from moves where game_id=? AND phase=? AND location_start=? AND game_year=? LIMIT 1"
	// TODO piece_id was removed this could cause a problem when pieces are saved/updated
	// TODO create a new piece when the move succeeds
	// moveInsert will fail because moves table piece_id is NOT NULL, it must exist prior to this being set.
	moveInsert := "Insert moves SET location_start=?, phase=?, game_id=?, type=?, piece_owner=?, game_year=?, location_submitted=?"
	moveUpdate := "Update moves SET location_start=?, phase=?, game_id=?, type=?, piece_owner=?, game_year=? WHERE location_submitted=?"
	move, err := e.fetchMove(ctx, doesPieceMoveExist, in.GameId, game.Phase, in.LocationStart, game.GameYear)
	if err != nil {
		fmt.Printf("error fetching move %v \n", err)
		return 500, err
	}

	// if the submitted move does not exist (based on piece id) create it; otherwise update it
	var insertType string
	switch move {
	case nil:
		insertType = moveInsert
	default:
		insertType = moveUpdate
	}
	stmt, err := e.Conn.PrepareContext(ctx, insertType)

	if err != nil {
		fmt.Printf("error inserting move %v \n", err)
		return 500, err
	}
	_, err = stmt.ExecContext(ctx, in.LocationStart, in.Phase, in.GameId, in.OrderType, in.PieceOwner, in.GameYear, in.LocationSubmitted)
	defer stmt.Close()
	if err != nil {
		fmt.Printf("\n error inserting move error %v with insert type: %v \n", err, insertType)
		return 500, err
	}

	// return 0
	return 200, err
}

func (e *MovesEngine) CreateBlankPiece(ctx context.Context, in *model.GameInput) (int64, error) {
	query := "Insert pieces SET location?, is_active=false"

	stmt, err := e.Conn.PrepareContext(ctx, query)

	if err != nil {
		fmt.Printf("**** PrepareContext %v", err)
		return -1, err
	}

	res, err := stmt.ExecContext(ctx, in.Title, "1901")
	defer stmt.Close()

	if err != nil {
		fmt.Printf("**** ExecContext %v", err)
		return -1, err
	}

	game_id, err := res.LastInsertId()

	if err != nil {
		fmt.Printf("**** ExecContext %v", err)
		return -1, err
	}

	return game_id, err
}

func (e *MovesEngine) fetchGame(ctx context.Context, query string, args ...interface{}) (*model.Game, error) {
	rows, err := e.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payload *model.Game
	for rows.Next() {
		data := &model.Game{}

		err := rows.Scan(
			&data.Id,
			&data.GameYear,
			&data.Phase,
			&data.PhaseEnd,
			&data.Title,
		)
		if err != nil {
			return nil, err
		}
		payload = data
	}
	return payload, nil
}

func (e *MovesEngine) fetchMove(ctx context.Context, query string, args ...interface{}) (*model.Move, error) {
	rows, err := e.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payload *model.Move
	for rows.Next() {
		data := &model.Move{}

		err := rows.Scan(
			&data.Id,
		)
		if err != nil {
			return nil, err
		}
		payload = data
	}
	return payload, nil
}

func (e *MovesEngine) ValidPiece(ctx context.Context, pieceId int64, locationStart model.Territory, gameId int64, moveCreatedByCountry model.Country, phase model.GamePhase) (err error) {
	query := "SELECT is_active, location, game_id, owner, dislodged, dislodged_from from pieces where id=?"
	rows, err := e.Conn.QueryContext(ctx, query, pieceId)
	if err != nil {
		return err
	}
	defer rows.Close()
	var payload *model.PieceRow
	for rows.Next() {
		data := &model.PieceRow{}
		err := rows.Scan(
			&data.IsActive,
			&data.Country,
			&data.GameId,
			&data.Owner,
			&data.Dislodged,
			&data.DislodgedFrom,
		)
		if err != nil {
			return err
		}
		payload = data
	}

	if phase == model.SpringRetreat || phase == model.FallRetreat {
		if !payload.Dislodged {
			return errors.New("Only dislodged pieces may move in Retreat Phase")
		}
	}
	if payload == nil {
		return errors.New("Piece could not be loaded")
	}
	if payload.Country != locationStart {
		return errors.New("Piece is not at start location")
	}
	if payload.GameId != gameId {
		return errors.New("Piece does not exist")
	}
	if !payload.IsActive {
		return errors.New("Piece is not active")
	}
	if payload.Owner != moveCreatedByCountry {
		return errors.New("Piece is not controlled by Country who issued the order")
	}
	return nil
}

func (e *MovesEngine) fetchGameUser(ctx context.Context, query string, args ...interface{}) (*model.GameUser, error) {
	rows, err := e.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payload *model.GameUser
	for rows.Next() {
		data := &model.GameUser{}

		err := rows.Scan(
			&data.UserId,
			&data.GameId,
			&data.Country,
		)
		if err != nil {
			return nil, err
		}
		payload = data
	}
	return payload, nil
}

// TODO IMPLEMENT
func (e *MovesEngine) Fetch(ctx context.Context, num int64) ([]*model.Game, error) {
	game := make([]*model.Game, 0)
	return game, nil
}

// TODO IMPLEMENT
func (e *MovesEngine) GetByID(ctx context.Context, id int64) (*model.Game, error) {
	game := &model.Game{}
	return game, nil
}

// TODO IMPLEMENT
func (e *MovesEngine) Update(ctx context.Context, in *model.Move) (*model.Move, int, error) {
	in = &model.Move{}
	return in, 200, nil
}
