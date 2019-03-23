package gamecrud

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	model "diplomacy/model"
)

// NewSQLPostRepo retunrs implement of game db interface
func NewEngine(Conn *sql.DB) Engine {
	e := Engine{
		Conn: Conn,
	}
	return e
}

type Engine struct {
	Conn *sql.DB
}

func (e *Engine) Create(ctx context.Context, in *model.GameInput) (int64, error) {
	query := "Insert games SET title=?, game_year=?"

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

	err = e.initTerritoryRecords(ctx, game_id)

	if err != nil {
		fmt.Printf("**** ExecContext %v", err)
		return -1, err
	}

	err = e.initGamePieceRecords(ctx, game_id)

	if err != nil {
		fmt.Printf("**** ExecContext %v", err)
		return -1, err
	}

	return game_id, err
}

func (e *Engine) fetch(ctx context.Context, query string, args ...interface{}) ([]*model.Game, error) {
	rows, err := e.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payload := make([]*model.Game, 0)
	for rows.Next() {
		data := new(model.Game)

		err := rows.Scan(
			&data.Id,
			&data.GameYear,
			&data.Phase,
			&data.PhaseEnd,
			&data.Title,
			&data.Processed,
		)
		if err != nil {
			return nil, err
		}
		payload = append(payload, data)
	}
	return payload, nil
}

func (e *Engine) fetchTerritories(ctx context.Context, args ...interface{}) ([]*model.TerritoryRow, error) {
	query := "select id, game_id, owner, country from territory where game_id=?"
	rows, err := e.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payload := make([]*model.TerritoryRow, 0)
	for rows.Next() {
		data := new(model.TerritoryRow)

		err := rows.Scan(
			&data.Id,
			&data.GameId,
			&data.Owner,
			&data.Country,
		)
		if err != nil {
			fmt.Printf("error \n", err)
			return nil, err
		}
		payload = append(payload, data)
	}
	return payload, nil
}

func (e *Engine) fetchPieces(ctx context.Context, args ...interface{}) ([]*model.PieceRow, error) {
	query := "select id, game_id, owner, type, is_active, location, dislodged from pieces where game_id=?"
	rows, err := e.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payload := make([]*model.PieceRow, 0)
	for rows.Next() {
		data := new(model.PieceRow)

		err := rows.Scan(
			&data.Id,
			&data.GameId,
			&data.Owner,
			&data.UnitType,
			&data.IsActive,
			&data.Country,
			&data.Dislodged,
		)
		if err != nil {
			fmt.Printf("error \n", err)
			return nil, err
		}
		payload = append(payload, data)
	}
	return payload, nil
}

func (e *Engine) ProcessMoves(ctx context.Context, gameId int64, phase int) error {
	moves, err := e.GetMovesByIdAndPhase(ctx, gameId, phase)
	if err != nil {
		return err
	}
	moves.ProcessMoves()
	e.save(ctx, moves)
	e.updateGameToProcessed(ctx, gameId)

	return err
}

func (e *Engine) save(ctx context.Context, moves model.Moves) (err error) {
	for _, move := range moves {
		query := "UPDATE pieces SET location=?, dislodged=? WHERE id=?"
		stmt, err := e.Conn.PrepareContext(ctx, query)
		if err != nil {
			fmt.Printf("err %v \n", err)
			return err
		}

		_, err = stmt.ExecContext(
			ctx,
			move.LocationResolved,
			move.Dislodged,
			move.PieceId,
		)

		if err != nil {
			fmt.Printf("err %v \n", err)
			return err
		}
		stmt.Close()
	}
	return err
}

func (e *Engine) updateGameToProcessed(ctx context.Context, gameId int64) (err error) {
	query := "UPDATE games SET processed=? WHERE id=?"
	stmt, err := e.Conn.PrepareContext(ctx, query)
	if err != nil {
		fmt.Printf("err %v \n", err)
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
		true,
		gameId,
	)

	if err != nil {
		fmt.Printf("err %v \n", err)
		return err
	}
	stmt.Close()
	return err
}

func (e *Engine) GetMovesByIdAndPhase(ctx context.Context, gameId int64, phase int) (model.Moves, error) {

	query := "select moves.id, moves.location_start, moves.location_submitted, moves.second_location_submitted, moves.type, moves.piece_owner, pieces.type, moves.piece_id from moves INNER JOIN pieces ON pieces.id=moves.id where moves.game_id=? and moves.phase=?"

	rows, err := e.Conn.QueryContext(ctx, query, gameId, phase)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payload := make([]*model.Move, 0)
	for rows.Next() {
		data := new(model.Move)

		err := rows.Scan(
			&data.Id,
			&data.LocationStart,
			&data.LocationSubmitted,
			&data.SecondLocationSubmitted,
			&data.OrderType,
			&data.PieceOwner,
			&data.UnitType,
			&data.PieceId,
		)

		if err != nil {
			fmt.Printf("error \n", err)
			return nil, err
		}
		payload = append(payload, data)
	}
	return payload, nil
}

// TODO SORT BY DATE
// WHERE game phase is 0
// Search by user_games
func (m *Engine) Fetch(ctx context.Context, num int64) ([]*model.Game, error) {
	query := "Select id, game_year, phase, phase_end, title From games limit ?"

	return m.fetch(ctx, query, num)
}

// TODO add game.IsActive; modify game query.
func (e *Engine) GetByID(ctx context.Context, gameId int64) (*model.Game, error) {
	query := "Select id, game_year, phase, phase_end, title, processed From games where id=?"

	// Get the Game by gameId
	rows, err := e.fetch(ctx, query, gameId)

	if err != nil {
		fmt.Printf("err %v \n", err)
		return nil, err
	}

	// Make the Game model
	game := &model.Game{}
	if len(rows) > 0 {
		game = rows[0]
	} else {
		return nil, nil //model.ErrNotFound
	}

	// Has the current phase ended?
	//   yes: process the moves
	//   yes: update the game.Phase and update game.PhaseEnd
	phaseOver := game.HasPhaseEnded()
	fmt.Printf("has this phase ended? %v\n", phaseOver)
	fmt.Printf("is this phase processed? %v\n", game.Processed)

	// TODO(:3/1) Should this be on moves.ProcessMoves()?
	// It could return piecesRows and we could switch on PhaseOver
	if phaseOver && !game.Processed {
		e.ProcessMoves(ctx, gameId, game.Phase)
	}

	// Get the Pieces of this game
	piecesRows, err := e.fetchPieces(ctx, gameId)
	if err != nil {
		fmt.Printf("err %v \n", err)
		return nil, err
	}

	// Get the Territories of this game
	territoryRows, err := e.fetchTerritories(ctx, gameId)
	if err != nil {
		fmt.Printf("err %v \n", err)
		return nil, err
	}

	// Make an array of Piece Models
	pm := &model.PieceRow{}
	pieces := make([]model.PieceRow, 0)
	for index, _ := range piecesRows {
		pm = piecesRows[index]
		pieces = append(pieces, *pm)
	}

	// Make an array of Territory Models
	tm := &model.TerritoryRow{}
	territories := make([]model.TerritoryRow, 0)
	for index, _ := range territoryRows {
		tm = territoryRows[index]
		territories = append(territories, *tm)
	}

	game.DrawGameBoard(territories, pieces)
	// return the Game
	return game, nil
}

func (e *Engine) getGameByIdOnly(ctx context.Context, id int64) (*model.Game, error) {
	query := "Select id, game_year, phase, phase_end, title From games where id=?"
	// fetch the Game
	rows, err := e.fetch(ctx, query, id)

	if err != nil {
		fmt.Printf("err %v \n", err)
		return nil, err
	}

	// Make the Game model
	game := &model.Game{}
	if len(rows) > 0 {
		game = rows[0]
	} else {
		return nil, nil //model.ErrNotFound
	}
	// return the Game
	return game, nil
}

// This endpoint is used to Join an existing Game
// Therefore it is an Update action
// Create Piece records, setting the user.id
// Create Territory records, setting the user.id
func (e *Engine) Update(ctx context.Context, in *model.GameInput) (*model.GameInput, int, error) {
	query := "Insert users_games SET user_id=?, country=?, game_id=?"
	stmt, err := e.Conn.PrepareContext(ctx, query)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return nil, 500, err
	}

	gameusers, err := e.getGameUsers(ctx, in.Id)

	err = model.ValidateCountryAndGameIsOpen(gameusers, in.Country)

	if err != nil {
		fmt.Printf("err1: %v\n", err)
		return nil, 409, err
	}

	_, err = stmt.ExecContext(
		ctx,
		in.UserId,
		in.Country,
		in.Id,
	)

	if err != nil {
		fmt.Printf("err2: %v\n", err)
		return nil, 500, err
	}
	defer stmt.Close()

	// The last user was added to the game with success of ExecContext()
	// and a user count of 6, update the phase!
	// TODO(:2/28) updateGamePhase should depend on interval
	if len(gameusers) == 6 {
		err := e.updateGamePhase(ctx, in.Id, 1)
		return nil, 500, err
	}

	return in, 200, nil
}

// Create Piece records, setting the user.id
// Create Territory records, setting the user.id
func (e *Engine) updateGamePhase(ctx context.Context, game_id int64, phase int) error {
	game, err := e.getGameByIdOnly(ctx, game_id)
	if err != nil {
		fmt.Printf("**** getGameByIdOnly %v\n", err)
		return err
	}
	err = game.ValidatePhaseUpdate(phase)
	if err != nil {
		fmt.Printf("**** error %v\n", err)
		return err
	}

	query := "UPDATE games SET phase = ?, phase_end=? WHERE id=?"
	stmt, err := e.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
		phase,
		time.Now().Add(time.Hour*time.Duration(12)).Unix(),
		game_id,
	)

	if err != nil {
		return err
	}
	defer stmt.Close()

	return nil
}

func (e *Engine) getGameUsers(ctx context.Context, game_id int64) ([]model.GameUser, error) {
	gameusers := []model.GameUser{}
	var err error
	var rows *sql.Rows

	rows, err = e.Conn.Query(`
          SELECT user_id, game_id, country
          FROM users_games
          WHERE game_id = ?`,
		game_id)
	if err != nil {
		return nil, fmt.Errorf("%v; inputs %#v", err, game_id)
	}
	defer rows.Close()

	for rows.Next() {
		var record model.GameUser
		if err = rows.Scan(&record.UserId, &record.GameId, &record.Country); err != nil {
			return nil, err
		}
		gameusers = append(gameusers, record)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return gameusers, err
}

func (e *Engine) initTerritoryRecords(ctx context.Context, game_id int64) error {
	g := model.Game{}
	g.NewGameBoard()
	query := "Insert INTO territory(game_id, country, owner) VALUES "

	for key, territory := range g.GameBoard {
		query += "(" + strconv.FormatInt(game_id, 10) + ", " + fmt.Sprintf("%#v", key) + ", " + fmt.Sprintf("%#v", territory.Owner) + "),"
	}

	//trim the last ,
	query = query[0 : len(query)-1]
	//prepare the statement
	fmt.Printf("query %v", query)
	stmt, err := e.Conn.PrepareContext(ctx, query)

	if err != nil {
		fmt.Printf("err %v", err)
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Printf("err %v", err)
	}

	return err
}

func (e *Engine) initGamePieceRecords(ctx context.Context, game_id int64) error {
	g := model.Game{}
	g.NewGameBoard()
	query := "Insert INTO pieces(game_id, type, location, owner) VALUES "

	for key, territory := range g.GameBoard {
		if len(territory.Units) > 0 {
			for _, piece := range territory.Units {
				query += "(" + strconv.FormatInt(game_id, 10) + ", " + fmt.Sprintf("%#v", piece.UnitType) + ", " + fmt.Sprintf("%#v", key) + ", " + fmt.Sprintf("%#v", piece.Owner) + "),"
			}
		}
	}

	//trim the last ,
	query = query[0 : len(query)-1]
	//prepare the statement
	fmt.Printf("query %v", query)
	stmt, err := e.Conn.PrepareContext(ctx, query)

	if err != nil {
		fmt.Printf("err %v", err)
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Printf("err %v", err)
	}

	return err
}
