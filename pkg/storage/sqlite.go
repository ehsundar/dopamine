package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/samber/lo"
	"log"
	"time"
)

type sqlite struct {
	db *sql.DB
}

func NewSqliteStorage(dsn string) Storage {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}

	return &sqlite{
		db: db,
	}
}

func (s *sqlite) GetAll(ctx context.Context, table string) ([]*Item, error) {
	query := getQuery("items/retrieve-many", table)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var items []*Item
	for rows.Next() {
		i := Item{}
		var c string
		err = rows.Scan(&i.ID, &c, &i.CreatedAt)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(c), &i.Contents)
		if err != nil {
			return nil, err
		}

		items = append(items, &i)
	}

	return items, nil
}

func (s *sqlite) GetOne(ctx context.Context, table string, id int) (*Item, error) {
	exists, err := s.checkTableExists(ctx, table)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrTableNotExist
	}

	query := getQuery("items/retrieve-one", table)
	row := s.db.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	i := Item{}
	var c string
	err = row.Scan(&i.ID, &c, &i.CreatedAt)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(c), &i.Contents)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (s *sqlite) InsertOne(ctx context.Context, table string, item *Item) (*Item, error) {
	err := s.createTable(ctx, table)
	if err != nil {
		return nil, err
	}

	contents, err := item.ToJSON(false)

	query := getQuery("items/insert-one", table)
	row := s.db.QueryRowContext(ctx, query, string(contents), time.Now())
	if row.Err() != nil {
		return nil, row.Err()
	}

	i := Item{}
	var c string
	err = row.Scan(&i.ID, &c, &i.CreatedAt)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(c), &i.Contents)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (s *sqlite) UpdateOne(ctx context.Context, table string, item *Item) (*Item, error) {
	contents, err := item.ToJSON(false)

	query := getQuery("items/update-one", table)
	row := s.db.QueryRowContext(ctx, query, string(contents), item.ID)
	if row.Err() != nil {
		return nil, row.Err()
	}

	i := Item{}
	var c string
	err = row.Scan(&i.ID, &c, &i.CreatedAt)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(c), &i.Contents)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (s *sqlite) DeleteOne(ctx context.Context, table string, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := getQuery("items/delete-one", table)
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *sqlite) createTable(ctx context.Context, table string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	createTableQuery := getQuery("items/create-table", table)
	stmt, err := s.db.PrepareContext(ctx, createTableQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *sqlite) checkTableExists(ctx context.Context, table string) (bool, error) {
	tables, err := s.tables(ctx)
	if err != nil {
		return false, err
	}
	return lo.Contains(tables, table), nil
}

func (s *sqlite) tables(ctx context.Context) ([]string, error) {
	query := getQuery("list-tables")
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return []string{}, err
	}

	var names []string
	var n string
	for rows.Next() {
		err = rows.Scan(&n)
		if err != nil {
			return []string{}, err
		}
		names = append(names, n)
	}
	return names, nil
}
