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

func (e *Engine) countActivePiecesByPlayer(ctx context.Context, args ...interface{}) (map[model.Country]int, error) {
	query := "select count(*), owner from pieces where game_id=? and is_active=true group by owner;"
	rows, err := e.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make(map[model.Country]int)
	for rows.Next() {
		var country model.Country
		var count int
		err := rows.Scan(
			&count,
			&country,
		)
		if err != nil {
			fmt.Printf("error \n", err)
			return nil, err
		}
		data[country] = count
	}
	fmt.Printf("Unit counts by country: %+v \n:", data)
	return data, nil
}

func (e *Engine) fetchPieces(ctx context.Context, args ...interface{}) ([]*model.PieceRow, error) {
	query := "select id, game_id, owner, type, is_active, location, dislodged, dislodged_from from pieces where game_id=?"
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
			&data.DislodgedFrom,
		)
		if err != nil {
			fmt.Printf("error \n", err)
			return nil, err
		}
		payload = append(payload, data)
	}
	return payload, nil
}

func (e *Engine) ProcessBuildPhase(
	ctx context.Context,
	game model.Game,
	pr []*model.PieceRow,
	countryToVictoryCenter, countryToUnitCount map[model.Country]int,
) error {
	builds := make(map[model.Country]int)
	for key, _ := range countryToVictoryCenter {
		builds[key] = countryToVictoryCenter[key] - countryToUnitCount[key]
	}
	fmt.Printf("builds! %+v \n", builds)
	_, err := e.GetMovesByIdAndPhase(ctx, game.Id, game.Phase, game.GameYear)
	return err
}

func (e *Engine) ProcessPhaseMoves(ctx context.Context, game model.Game) error {
	moves, err := e.GetMovesByIdAndPhase(ctx, game.Id, game.Phase, game.GameYear)
	if err != nil {
		return err
	}

	moves, err = e.ProcessPiecesNotMoved(ctx, moves, game.Id, game.Phase)

	if err != nil {
		return err
	}

	moves.ProcessMoves()
	e.updatePieces(ctx, moves)
	e.updateTerritories(ctx, moves, game)
	e.updateResolvedMoves(ctx, moves)

	return err
}

// This must happen prior to process moves to create Hold moves for unmoved pieces
func (e *Engine) ProcessPiecesNotMoved(ctx context.Context, moves model.Moves, gameId int64, phaseId model.GamePhase) (newMoves model.Moves, err error) {
	pieces, err := e.GetPiecesByGameId(ctx, gameId)
	if err != nil {
		return nil, err
	}

	newMoves = moves.HoldUnmovedPieces(pieces)

	for _, move := range newMoves {
		id, err := e.CreateMove(ctx, move, gameId, phaseId)
		if err != nil {
			return nil, err
		}
		move.Id = id
	}

	for _, move := range moves {
		newMoves = append(newMoves, move)
	}

	return newMoves, err
}

func (e *Engine) CreateMove(ctx context.Context, in *model.Move, gameId int64, phaseId model.GamePhase) (int64, error) {
	insert := "Insert moves SET location_start=?, location_submitted=?, phase=?, game_id=?, type=?, piece_owner=?, game_year=?, piece_id=?"

	stmt, err := e.Conn.PrepareContext(ctx, insert)

	if err != nil {
		return 0, err
	}
	res, err := stmt.ExecContext(ctx, in.LocationStart, in.LocationSubmitted, phaseId, gameId, in.OrderType, in.PieceOwner, in.GameYear, in.PieceId)
	defer stmt.Close()

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}
	// return 0
	return id, err
}

func (e *Engine) updateResolvedMoves(ctx context.Context, moves model.Moves) (err error) {
	for _, move := range moves {
		update := "UPDATE moves SET location_resolved=? WHERE id=?"
		stmt, err := e.Conn.PrepareContext(ctx, update)
		if err != nil {
			fmt.Printf("err %v \n", err)
			return err
		}

		_, err = stmt.ExecContext(
			ctx,
			move.LocationResolved,
			move.Id,
		)

		if err != nil {
			fmt.Printf("err %v \n", err)
			return err
		}
		stmt.Close()
	}
	return err
}

// returns a map of Players mapped to a Victory Center count.
// Measured by Territory Owner saved in db
func (e *Engine) countVictoryCenters(ctx context.Context, territories []*model.TerritoryRow, game model.Game) (vc map[model.Country]int) {
	if game.Phase == model.Waiting {
		return
	}
	vc = make(map[model.Country]int)
	for _, t := range territories {
		if t.Owner == "" {
			continue
		}

		// account for edge cases
		if t.Country == "SPN" || t.Country == "SPS" || t.Country == "SNC" || t.Country == "SSC" || t.Country == "BUE" || t.Country == "BUS" {
			continue
		}
		if t.Country.IsVictoryCenter() {
			vc[t.Owner] += 1
		}
	}
	return vc

}

func (e *Engine) updateTerritories(ctx context.Context, moves model.Moves, game model.Game) (err error) {
	if game.Phase != model.FallRetreat {
		return
	}

	for _, move := range moves {
		query := "UPDATE territory SET owner=? WHERE country=? AND game_id=?"
		stmt, err := e.Conn.PrepareContext(ctx, query)
		if err != nil {
			fmt.Printf("err %v \n", err)
			return err
		}

		_, err = stmt.ExecContext(
			ctx,
			move.PieceOwner,
			move.LocationResolved,
			game.Id,
		)

		if err != nil {
			fmt.Printf("err %v \n", err)
			return err
		}
		stmt.Close()
	}
	return err
}

func (e *Engine) updatePieces(ctx context.Context, moves model.Moves) (err error) {
	for _, move := range moves {
		query := "UPDATE pieces SET location=?, dislodged=?, dislodged_from=? WHERE id=?"
		stmt, err := e.Conn.PrepareContext(ctx, query)
		if err != nil {
			fmt.Printf("err %v \n", err)
			return err
		}

		_, err = stmt.ExecContext(
			ctx,
			move.LocationResolved,
			move.Dislodged,
			move.DislodgedFrom,
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

func (e *Engine) updateGameToProcessed(ctx context.Context, game model.Game) (err error) {
	query := "UPDATE games SET phase=?, game_year=?, phase_end=? WHERE id=?"
	stmt, err := e.Conn.PrepareContext(ctx, query)
	if err != nil {
		fmt.Printf("err %v \n", err)
		return err
	}

	phase := model.NewPhase(game.Phase)
	year := newYear(game)
	phaseOver := time.Now().Add(time.Hour * time.Duration(12)).Unix()

	_, err = stmt.ExecContext(
		ctx,
		phase,
		year,
		phaseOver,
		game.Id,
	)

	if err != nil {
		fmt.Printf("err %v \n", err)
		return err
	}
	stmt.Close()
	return err
}

func newYear(game model.Game) string {
	if game.Phase == 5 {
		n, err := strconv.ParseInt(game.GameYear, 10, 64)
		if err == nil {
			fmt.Printf("%d of type %T", n, n)
		}
		n = n + 1
		return strconv.FormatInt(n, 10)
	} else {
		return game.GameYear
	}
}

func (e *Engine) GetMovesByIdAndPhase(ctx context.Context, gameId int64, phase model.GamePhase, year string) (model.Moves, error) {

	query := "select moves.id, moves.location_start, moves.location_submitted, moves.second_location_submitted, moves.type, moves.piece_owner, moves.game_year, pieces.type, moves.piece_id from moves INNER JOIN pieces ON pieces.id=moves.piece_id where moves.game_id=? and moves.phase=? and moves.game_year=?"

	rows, err := e.Conn.QueryContext(ctx, query, gameId, phase, year)
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
			&data.GameYear,
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

func (e *Engine) GetPiecesByGame(ctx context.Context, gameId int64) ([]*model.PieceRow, error) {

	query := "select id, location from pieces where game_id=?"
	rows, err := e.Conn.QueryContext(ctx, query, gameId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payload := make([]*model.PieceRow, 0)
	for rows.Next() {
		data := new(model.PieceRow)

		err := rows.Scan(
			&data.Id,
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

func (e *Engine) GetPiecesByGameId(ctx context.Context, gameId int64) ([]*model.PieceRow, error) {

	query := "select id, owner, location, type, is_active, dislodged_from from pieces where game_id=?"

	rows, err := e.Conn.QueryContext(ctx, query, gameId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payload := make([]*model.PieceRow, 0)
	for rows.Next() {
		data := new(model.PieceRow)

		err := rows.Scan(
			&data.Id,
			&data.Owner,
			&data.Country,
			&data.UnitType,
			&data.IsActive,
			&data.DislodgedFrom,
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
	query := "Select id, game_year, phase, phase_end, title From games where id=?"

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
	//   if so: process the move
	fmt.Printf("%v over??? \n", game.Phase.HasPhaseEnded(game.PhaseEnd))
	if game.Phase.HasPhaseEnded(game.PhaseEnd) && game.Phase != model.FallBuild {
		e.ProcessPhaseMoves(ctx, *game)
		e.updateGameToProcessed(ctx, *game)
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

	// Count the Victory Centers
	vc := e.countVictoryCenters(ctx, territoryRows, *game)

	// Count active units by player
	au, err := e.countActivePiecesByPlayer(ctx, game.Id)

	// Has the build phase ended?
	if game.Phase.HasPhaseEnded(game.PhaseEnd) && game.Phase == model.FallBuild {
		e.ProcessBuildPhase(ctx, *game, piecesRows, vc, au)
		e.updateGameToProcessed(ctx, *game)
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

	game.DrawGameBoard(territories, pieces, vc)
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
		err := e.updateGamePhase(ctx, in.Id, model.Spring)
		return nil, 500, err
	}

	return in, 200, nil
}

// Create Piece records, setting the user.id
// Create Territory records, setting the user.id
func (e *Engine) updateGamePhase(ctx context.Context, game_id int64, phase model.GamePhase) error {
	game, err := e.getGameByIdOnly(ctx, game_id)
	if err != nil {
		fmt.Printf("**** getGameByIdOnly %v\n", err)
		return err
	}
	err = game.Phase.ValidatePhaseUpdate(phase)
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
