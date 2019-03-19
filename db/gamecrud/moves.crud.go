package gamecrud

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	model "diplomacy/model"
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
	doesPieceMoveExist := "select id from moves where game_id=? AND phase=? AND piece_id=?"
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
	err = e.ValidateCountry(in, gameUser)

	if err != nil {
		return 403, err
	}

	// Has the game started?
	// Has the phase ended?
	err = game.ValidPhase()

	if err != nil {
		return 403, err
	}

	// Does the submitted move exist for this piece?
	move, err := e.fetchMove(ctx, doesPieceMoveExist, in.GameId, in.Phase, in.PieceId)

	if err != nil {
		fmt.Printf("error fetchMove(): %v \n", err)
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

// This validation does not make sense
func (e *MovesEngine) ValidateCountry(move *model.Move, gameUser *model.GameUser) (err error) {
	if move.PieceOwner != gameUser.Country {
		return errors.New("The User does not control this country")
	}
	return err
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
