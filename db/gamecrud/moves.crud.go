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
func (e *movesEngine) Create(ctx context.Context, in *model.GameInput) (int64, error) {
	query := "Select user_id, game_id, country from users_games where user_id=? and game_id=?"
	gameQuery := "Select id, game_year, phase, phase_end, title From games where id=?"
	// moveQuery := "Select id, game_year, phase, piece_id where id=?"
	// type corresponds to move type (hold, move, support)
	pieceMoveInsert := "Insert moves SET location_start=?, location_submitted=?, phase=?, piece_id=?, game_id=?, type=?"

	res, err := e.fetchGameUser(ctx, query, in.UserId, in.Id)

	if res == nil {
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

	if err != nil {
		fmt.Printf("error fetchGameUser(): %v \n", err)
		return -1, err
	}
	// TODO IMPLEMENT
	_, err = e.ValidateCountry(in, res)
	_, err = e.ValidPhase(in, game)

	stmt, err := e.Conn.PrepareContext(ctx, pieceMoveInsert)

	if err != nil {
		return -1, err
	}

	fmt.Printf("********* %+v \n", in)

	_, err = stmt.ExecContext(ctx, in.LocationStart, in.LocationSubmitted, in.Phase, in.PieceId, in.Id, in.MoveType)
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

func (e *movesEngine) ValidateCountry(in *model.GameInput, res *model.GameUser) (valid bool, err error) {
	if in.Country != res.Country {
		return false, errors.New("The User does not control this country")
	}

	return true, err
}

// TODO determine when to move game from phase 0 -> phase 1
// TODO determine where to set the phase time.Time when the above occurs
func (e *movesEngine) ValidPhase(in *model.GameInput, game *model.Game) (valid bool, err error) {
	// fetch the Game
	if game.Phase < 1 {
		return false, errors.New("The Game has not started yet")
	}
	now := time.Now()
	fmt.Printf("*** before *game.PhaseEnd%v\n", *game)
	timestamp, err := strconv.ParseInt(game.PhaseEnd, 10, 64)
	tm := time.Unix(timestamp, 0)
	valid = now.Before(tm)
	if err != nil {
		panic(err)
	}
	fmt.Printf("***** valid %v \n", valid)
	return valid, err
	// TODO compare the time.Before()
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
