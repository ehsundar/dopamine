package storage

import (
	"context"
	"github.com/doug-martin/goqu/v9"
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

	i, err := cqs.QueryOne(context.Background(), "test", func(dataset *goqu.SelectDataset) *goqu.SelectDataset {
		return dataset.Where(goqu.Ex{
			"id": item.ID,
		})
	})
	s.Nil(err)

	s.Equal("{\"key1\":\"value1\"}", i.Contents)

	teardownDatabase(s)
}

func TestSqlite_QueryMany(t *testing.T) {
	s := assert.New(t)

	storage := NewSqliteStorage(dsn)
	_, err := storage.InsertOne(
		context.Background(),
		"test",
		&Item{
			ContentsMap: map[string]any{
				"key1": "value1",
			},
		},
	)
	s.Nil(err)

	_, err = storage.InsertOne(context.Background(), "test", &Item{
		ContentsMap: map[string]any{
			"key2": "value2",
		},
	})
	s.Nil(err)

	cqs, ok := storage.(CustomQueryStorage)
	s.True(ok)

	items, err := cqs.QueryMany(
		context.Background(),
		"test",
		func(dataset *goqu.SelectDataset) *goqu.SelectDataset {
			return dataset
		},
	)
	s.Nil(err)

	s.Len(items, 2)

	teardownDatabase(s)
}

func teardownDatabase(s *assert.Assertions) {
	err := os.Remove(dsn)
	s.Nil(err)
}
