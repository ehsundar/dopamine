package storage

import (
	"context"
	"errors"
)

var (
	ErrTableNotExist = errors.New("table does not exist")
)

type Storage interface {
	GetAll(ctx context.Context, table string) ([]*Item, error)
	GetOne(ctx context.Context, table string, id int) (*Item, error)
	InsertOne(ctx context.Context, table string, item *Item) (*Item, error)
	UpdateOne(ctx context.Context, table string, item *Item) (*Item, error)
	DeleteOne(ctx context.Context, table string, id int) error
}
