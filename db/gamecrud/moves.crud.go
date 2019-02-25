package gamecrud

import (
	"context"
	"database/sql"
	// "fmt"

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
	query := "Insert pieces_moves SET title=?, game_year=?"

	// TODO IMPLEMENT
	err = e.Validate(in.Id, in.UserId, in.Country, in.Phase, in.PhaseEnd)

	stmt, err := e.Conn.PrepareContext(ctx, query)

	if err != nil {
		return -1, err
	}

	res, err := stmt.ExecContext(ctx, in.Title, "1901-04-01")
	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	game_id, err := res.LastInsertId()

	if err != nil {
		return -1, err
	}

	// err = e.setTerritoryRecords(ctx, game_id)

	// if err != nil {
	// 	return -1, err
	// }

	// err = e.setGamePieceRecords(ctx, game_id)

	// if err != nil {
	// 	return -1, err
	// }

	return game_id, err
}

// TODO IMPLEMENT
func (m *movesEngine) Fetch(ctx context.Context, num int64) ([]*model.Game, error) {
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
