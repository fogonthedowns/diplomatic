package gamecrud

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	db "diplomacy/db"
	model "diplomacy/model"
)

// NewSQLPostRepo retunrs implement of game db interface
func NewEngine(Conn *sql.DB) db.Crud {
	return &engine{
		Conn: Conn,
	}
}

type engine struct {
	Conn *sql.DB
}

func (e *engine) Create(ctx context.Context, in *model.GameInput) (int64, error) {
	query := "Insert games SET title=?, game_year=?"

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

	err = e.setTerritoryRecords(ctx, game_id)

	if err != nil {
		return -1, err
	}

	err = e.setGamePieceRecords(ctx, game_id)

	if err != nil {
		return -1, err
	}

	return game_id, err
}

func (e *engine) fetch(ctx context.Context, query string, args ...interface{}) ([]*model.Game, error) {
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

func (e *engine) fetchTerritories(ctx context.Context, query string, args ...interface{}) ([]*model.TerritoryRow, error) {
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

// func (m *mysqlPostRepo) Fetch(ctx context.Context, num int64) ([]*models.Post, error) {
// 	query := "Select id, title, content From posts limit ?"

// 	return m.fetch(ctx, query, num)
// }

func (e *engine) GetByID(ctx context.Context, id int64) (*model.Game, error) {
	query := "Select id, game_year, phase, phase_end, title From games where id=?"
	second_query := "select * from territory where game_id=?"

	// fetch the Pieces of this game
	// territory_rows, err := e.fetchPieces(ctx, third_query, id)
	// if err != nil {
	// 	fmt.Printf("err %v \n", err)
	// 	return nil, err
	// }

	// fetch the Territories of this game
	territory_rows, err := e.fetchTerritories(ctx, second_query, id)
	if err != nil {
		fmt.Printf("err %v \n", err)
		return nil, err
	}

	// fetch the Game
	rows, err := e.fetch(ctx, query, id)

	if err != nil {
		fmt.Printf("err %v \n", err)
		return nil, err
	}

	// Make an array of Territory Models
	tm := &model.TerritoryRow{}
	territories := make([]model.TerritoryRow, 0)
	for index, _ := range territory_rows {
		tm = territory_rows[index]
		territories = append(territories, *tm)
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

// Create Piece records, setting the user.id
// Create Territory records, setting the user.id
func (e *engine) Update(ctx context.Context, in *model.GameInput) (*model.GameInput, int, error) {
	query := "Insert users_games SET user_id=?, country=?, game_id=?"
	stmt, err := e.Conn.PrepareContext(ctx, query)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return nil, 500, err
	}

	gameusers, err := e.getGameUsers(ctx, in.Id)
	err = model.Validate(gameusers, in.Country)

	if err != nil {
		fmt.Printf("err %v\n", err)
		return nil, 409, err
	}

	_, err = stmt.ExecContext(
		ctx,
		in.UserId,
		in.Country,
		in.Id,
	)

	if err != nil {
		fmt.Printf("err %v\n", err)
		return nil, 500, err
	}
	defer stmt.Close()

	return in, 200, nil
}

func (e *engine) getGameUsers(ctx context.Context, game_id int64) ([]model.GameUser, error) {
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

func (e *engine) setTerritoryRecords(ctx context.Context, game_id int64) error {
	g := model.Game{}
	g.BuildGameBoard()
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

func (e *engine) setGamePieceRecords(ctx context.Context, game_id int64) error {
	g := model.Game{}
	g.BuildGameBoard()
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

// func (m *mysqlPostRepo) Delete(ctx context.Context, id int64) (bool, error) {
// 	query := "Delete From posts Where id=?"

// 	stmt, err := m.Conn.PrepareContext(ctx, query)
// 	if err != nil {
// 		return false, err
// 	}
// 	_, err = stmt.ExecContext(ctx, id)
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }
