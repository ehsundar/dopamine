package storage

import (
	"context"
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const dsn = "dopamine_test.sqlite3"

func TestSqlite_QueryOne(t *testing.T) {
	s := assert.New(t)

	storage := NewSqliteStorage(dsn)
	item, err := storage.InsertOne(context.Background(), "test", &Item{
		ContentsMap: map[string]any{
			"key1": "value1",
		},
	})
	s.Nil(err)
	s.NotNil(item)

	cqs, ok := storage.(CustomQueryStorage)
	s.True(ok)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}
	dbqu := goqu.New("sqlite3", db)
	expr := dbqu.From("test").Where(goqu.Ex{
		"id": item.ID,
	})
	i, err := cqs.QueryOne(context.Background(), expr)
	s.Nil(err)

	s.Equal("{\"key1\":\"value1\"}", i.Contents)

	teardownDatabase(s)
}

func TestSqlite_QueryMany(t *testing.T) {
	s := assert.New(t)

	storage := NewSqliteStorage(dsn)
	_, err := storage.InsertOne(context.Background(), "test", &Item{
		ContentsMap: map[string]any{
			"key1": "value1",
		},
	})
	s.Nil(err)

	_, err = storage.InsertOne(context.Background(), "test", &Item{
		ContentsMap: map[string]any{
			"key2": "value2",
		},
	})
	s.Nil(err)

	cqs, ok := storage.(CustomQueryStorage)
	s.True(ok)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}
	dbqu := goqu.New("sqlite3", db)
	expr := dbqu.From("test")
	items, err := cqs.QueryMany(context.Background(), expr)
	s.Nil(err)

	s.Len(items, 2)

	teardownDatabase(s)
}

func teardownDatabase(s *assert.Assertions) {
	err := os.Remove(dsn)
	s.Nil(err)
}
