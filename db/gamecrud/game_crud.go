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

	err = e.setTerritoryRecords(ctx, game_id, in.Country)

	if err != nil {
		return -1, err
	}

	return game_id, err
}

// func (m *mysqlPostRepo) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Post, error) {
// 	rows, err := m.Conn.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	payload := make([]*models.Post, 0)
// 	for rows.Next() {
// 		data := new(models.Post)

// 		err := rows.Scan(
// 			&data.ID,
// 			&data.Title,
// 			&data.Content,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		payload = append(payload, data)
// 	}
// 	return payload, nil
// }

// func (m *mysqlPostRepo) Fetch(ctx context.Context, num int64) ([]*models.Post, error) {
// 	query := "Select id, title, content From posts limit ?"

// 	return m.fetch(ctx, query, num)
// }

// func (m *mysqlPostRepo) FindByID(ctx context.Context, id int64) (*model.Game, error) {
// 	query := "Select id, title, content From posts where id=?"

// 	rows, err := m.fetch(ctx, query, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	payload := &models.Post{}
// 	if len(rows) > 0 {
// 		payload = rows[0]
// 	} else {
// 		return nil, models.ErrNotFound
// 	}

// 	return payload, nil
// }

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

func (e *engine) setTerritoryRecords(ctx context.Context, game_id int64, country model.Country) error {
	g := model.Game{}
	g.BuildGameBoard()
	query := "Insert INTO territory(game_id, country, owner) VALUES "
	// vals := []interface{}{}

	for key, territory := range g.GameBoard {
		// query += "(?, ?, ?),"
		query += "(" + strconv.FormatInt(game_id, 10) + ", " + fmt.Sprintf("%#v", key) + ", " + fmt.Sprintf("%#v", territory.Owner) + "),"
		// vals = append(vals, game_id, key, territory.Owner)
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
	//format all vals at once
	_, err = stmt.Exec()
	fmt.Printf("err %v", err)
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
