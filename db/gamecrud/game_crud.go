package gamecrud

import (
	"context"
	"database/sql"
	"fmt"

	model "diplomacy/model"
	db "diplomacy/db"
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

func (e *engine) Create(ctx context.Context, p *model.Game) (int64, error) {
	query := "Insert games SET title=?"

	stmt, err := e.Conn.PrepareContext(ctx, query)
	fmt.Printf("err %+v \n", err)
	if err != nil {
		return -1, err
	}

	res, err := stmt.ExecContext(ctx, p.Title)
	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
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

// func (m *mysqlPostRepo) GetByID(ctx context.Context, id int64) (*models.Post, error) {
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



// func (m *mysqlPostRepo) Update(ctx context.Context, p *models.Post) (*models.Post, error) {
// 	query := "Update posts set title=?, content=? where id=?"

// 	stmt, err := m.Conn.PrepareContext(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	_, err = stmt.ExecContext(
// 		ctx,
// 		p.Title,
// 		p.Content,
// 		p.ID,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer stmt.Close()

// 	return p, nil
// }

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
