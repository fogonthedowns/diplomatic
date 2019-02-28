package gamecrud

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

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

func (e *MovesEngine) CreateOrUpdate(ctx context.Context, in *model.MoveInput) (int64, error) {
	query := "Select user_id, game_id, country from users_games where user_id=? and game_id=?"
	gameQuery := "Select id, game_year, phase, phase_end, title From games where id=?"
	doesPieceMoveExist := "select id from moves where game_id=? AND phase=? AND piece_id=?"
	pieceMoveInsert := "Insert moves SET location_start=?, location_submitted=?, phase=?, game_id=?, type=?, piece_id=?"
	pieceMoveUpdate := "Update moves SET location_start=?, location_submitted=?, phase=?, game_id=?, type=? WHERE piece_id=?"

	gameUser, err := e.fetchGameUser(ctx, query, in.UserId, in.GameId)
	fmt.Printf("GAME ID? %+v", in)

	// Is the User part of this game?
	if gameUser == nil {
		return -1, errors.New("The User is not a member of this game")
	}

	if err != nil {
		fmt.Printf("error fetchGameUser(): %v \n", err)
		return -1, err
	}

	game, err := e.fetchGame(ctx, gameQuery, in.GameId)

	if game == nil {
		return -1, errors.New("The Game can not be loaded")
	}

	// Does your User control the piece you are trying to move?
	err = e.ValidateCountry(in, gameUser)

	if err != nil {
		return -1, err
	}

	// Has the game started?
	// Has the phase ended?
	err = e.ValidPhase(in, game)

	if err != nil {
		return -1, err
	}

	// Does the submitted move exist for this piece?
	move, err := e.fetchMove(ctx, doesPieceMoveExist, in.GameId, in.Phase, in.PieceId)

	if err != nil {
		fmt.Printf("error fetchMove(): %v \n", err)
		return -1, err
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
		return -1, err
	}
	_, err = stmt.ExecContext(ctx, in.LocationStart, in.LocationSubmitted, in.Phase, in.GameId, in.OrderType, in.PieceId)
	defer stmt.Close()
	if err != nil {
		return -1, err
	}

	// return 0
	return 0, err
}

// This validation does not make sense
func (e *MovesEngine) ValidateCountry(move *model.MoveInput, gameUser *model.GameUser) (err error) {
	if move.PieceOwner != gameUser.Country {
		return errors.New("The User does not control this country")
	}
	return err
}

// TODO determine when to move game from phase 0 -> phase 1
// TODO determine where to set the phase time.Time when the above occurs
func (e *MovesEngine) ValidPhase(in *model.MoveInput, game *model.Game) (err error) {
	// fetch the Game
	if game.Phase < 1 {
		return errors.New("The Game has not started yet")
	}
	now := time.Now()
	timestamp, err := strconv.ParseInt(game.PhaseEnd, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(timestamp, 0)
	valid := now.Before(tm)
	if !valid {
		return errors.New("The phase has ended")
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

func (e *MovesEngine) fetchMove(ctx context.Context, query string, args ...interface{}) (*model.MoveInput, error) {
	rows, err := e.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payload *model.MoveInput
	for rows.Next() {
		data := &model.MoveInput{}

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
func (e *MovesEngine) Update(ctx context.Context, in *model.MoveInput) (*model.MoveInput, int, error) {
	in = &model.MoveInput{}
	return in, 200, nil
}
