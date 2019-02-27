package gamecrud

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	db "diplomacy/db"
	model "diplomacy/model"
)

// NewSQLPostRepo retunrs implement of game db interface
func NewMovesEngine(Conn *sql.DB) db.Crud {
	return &movesEngine{
		Conn: Conn,
	}
}

type movesEngine struct {
	Conn *sql.DB
}

// This could be update or create but this conforms to the interface
// A Player can send moves
// Validate a player is part of the game
// Validate they own the country
// Validate time/phase
// Count the pieces_moves table, when the records are complete or when time expires update the pieces_moves.location_resolved
// Update the game phase, year and phase_end based on the orders_interval

// TODO USE A DIFFERNET STRUCT where COUNTRY IS NOT AMBIGOUS
// ITS REALLY THE PIECE COUNTRY BUT I LOST 20 minutes trying to determine if it was user country!
// IN THIS STRUCT INTRODUCE FINALIZED CONCEPT

func (e *movesEngine) Create(ctx context.Context, in *model.GameInput) (int64, error) {
	query := "Select user_id, game_id, country from users_games where user_id=? and game_id=?"
	gameQuery := "Select id, game_year, phase, phase_end, title From games where id=?"
	doesPieceMoveExist := "select id from moves where game_id=? AND phase=? AND piece_id=?"
	pieceMoveInsert := "Insert moves SET location_start=?, location_submitted=?, phase=?, game_id=?, type=?, piece_id=?"
	pieceMoveUpdate := "Update moves SET location_start=?, location_submitted=?, phase=?, game_id=?, type=? WHERE piece_id=?"

	gameUser, err := e.fetchGameUser(ctx, query, in.UserId, in.Id)

	if gameUser == nil {
		return -1, errors.New("The User is not a member of this game")
	}

	if err != nil {
		fmt.Printf("error fetchGameUser(): %v \n", err)
		return -1, err
	}

	game, err := e.fetchGame(ctx, gameQuery, in.Id)

	if game == nil {
		return -1, errors.New("The Game can not be loaded")
	}

	// validate
	err = e.ValidateCountry(in, gameUser)

	if err != nil {
		return -1, err
	}

	err = e.ValidPhase(in, game)

	if err != nil {
		return -1, err
	}

	// Does the submitted move exist for this piece?
	move, err := e.fetchMove(ctx, doesPieceMoveExist, in.Id, in.Phase, in.PieceId)

	if err != nil {
		fmt.Printf("error fetchMove(): %v \n", err)
		return -1, err
	}

	// if the submitted move does not exist create it; otherwise update it
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

	_, err = stmt.ExecContext(ctx, in.LocationStart, in.LocationSubmitted, in.Phase, in.Id, in.MoveType, in.PieceId)
	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	// res_id, err := res.LastInsertId()

	if err != nil {
		return -1, err
	}

	return 0, err
}

// This validation does not make sense
func (e *movesEngine) ValidateCountry(in *model.GameInput, res *model.GameUser) (err error) {
	if in.Country != res.Country {
		return errors.New("The User does not control this country")
	}
	return err
}

// TODO determine when to move game from phase 0 -> phase 1
// TODO determine where to set the phase time.Time when the above occurs
func (e *movesEngine) ValidPhase(in *model.GameInput, game *model.Game) (err error) {
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

func (e *movesEngine) fetchGame(ctx context.Context, query string, args ...interface{}) (*model.Game, error) {
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

func (e *movesEngine) fetchMove(ctx context.Context, query string, args ...interface{}) (*model.GameInput, error) {
	rows, err := e.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payload *model.GameInput
	for rows.Next() {
		data := &model.GameInput{}

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

func (e *movesEngine) fetchGameUser(ctx context.Context, query string, args ...interface{}) (*model.GameUser, error) {
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
func (e *movesEngine) Fetch(ctx context.Context, num int64) ([]*model.Game, error) {
	game := make([]*model.Game, 0)
	return game, nil
}

// TODO IMPLEMENT
func (e *movesEngine) GetByID(ctx context.Context, id int64) (*model.Game, error) {
	game := &model.Game{}
	return game, nil
}

// TODO IMPLEMENT
func (e *movesEngine) Update(ctx context.Context, in *model.GameInput) (*model.GameInput, int, error) {
	in = &model.GameInput{}
	return in, 200, nil
}
