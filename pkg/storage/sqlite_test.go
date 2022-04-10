package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlite_GetAll(t *testing.T) {
	s := assert.New(t)

	storage := NewSqliteStorage("dopamine_test.sqlite3")

	item, err := storage.InsertOne(context.Background(), "test", &Item{
		ContentsMap: map[string]any{
			"key1": "value1",
		},
	})
	s.Nil(err)
	s.NotNil(item)

	items, err := storage.GetAll(context.Background(), "test")
	s.Nil(err)
	s.NotNil(items)
	s.Len(items, 1)
	s.Equal("{\"key1\":\"value1\"}", items[0].Contents)

	teardownDatabase(s)
}
