package storage

import (
	"encoding/json"
	"github.com/samber/lo"
	"time"
)

type Item struct {
	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`

	Contents map[string]any `db:"contents"`
}

func (i *Item) ToJSON(includeMeta bool) ([]byte, error) {
	m := i.ToMap(includeMeta)
	result, err := json.Marshal(m)
	return result, err
}

func (i *Item) ToMap(includeMeta bool) map[string]any {
	m := make(map[string]any)

	for k, v := range i.Contents {
		m[k] = v
	}

	if includeMeta {
		m["id"] = i.ID
		m["created_at"] = i.CreatedAt
	}
	return m
}

func ItemFromJSON(j []byte) (*Item, error) {
	m := make(map[string]any)
	i := Item{}

	err := json.Unmarshal(j, &m)
	if err != nil {
		return nil, err
	}

	delete(m, "id")
	delete(m, "created_at")

	i.Contents = m
	return &i, nil
}

func ItemsToJSON(items []*Item, includeMeta bool) ([]byte, error) {
	lst := lo.Map(items, func(i *Item, _ int) map[string]any {
		return i.ToMap(includeMeta)
	})

	return json.Marshal(lst)
}
