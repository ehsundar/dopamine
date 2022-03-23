package items

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Item struct {
	ID        int       `db:"id"`
	Contents  string    `db:"contents"`
	CreatedAt time.Time `db:"created_at"`
}

func getQuery(name string, unsafeParams ...string) string {
	query, err := ioutil.ReadFile("sql/" + name + ".sql")
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot read query: %s: %s", name, err))
	}
	q := string(query)

	for _, p := range unsafeParams {
		q = strings.Replace(q, "?", p, 1)
	}

	return q
}

func createTable(ctx context.Context, db *sql.DB, namespace string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	createTableQuery := getQuery("items/create-table", namespace)
	stmt, err := db.PrepareContext(ctx, createTableQuery)
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

func insertOneItem(ctx context.Context, db *sql.DB, namespace string, contents string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	createTableQuery := getQuery("items/insert-one", namespace)
	stmt, err := db.PrepareContext(ctx, createTableQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, contents, time.Now())
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func listItems(ctx context.Context, db *sql.DB, namespace string) ([]*Item, error) {
	rows, err := db.QueryContext(ctx, "select id, contents, created_at from "+namespace)
	if err != nil {
		return nil, err
	}

	var items []*Item
	for rows.Next() {
		i := Item{}
		err = rows.Scan(&i.ID, &i.Contents, &i.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, &i)
	}

	return items, nil
}

func getItem(ctx context.Context, db *sql.DB, namespace string, id int) (*Item, error) {
	retriveOneQuery := getQuery("items/retrieve-one", namespace)
	row := db.QueryRowContext(ctx, retriveOneQuery, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	i := Item{}
	err := row.Scan(&i.ID, &i.Contents, &i.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func updateItem(ctx context.Context, db *sql.DB, namespace string, id int, contents string) (*Item, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	query := getQuery("items/update-one", namespace)
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, contents, id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func deleteItem(ctx context.Context, db *sql.DB, namespace string, id int) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := getQuery("items/delete-one", namespace)
	stmt, err := db.PrepareContext(ctx, query)
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
