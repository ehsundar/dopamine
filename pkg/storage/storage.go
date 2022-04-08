package storage

import (
	"context"
	"errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

var (
	ErrTableNotExist = errors.New("table does not exist")
)

type QueryBuilder func(dataset *goqu.SelectDataset) (filteredDataset *goqu.SelectDataset)

type Storage interface {
	GetAll(ctx context.Context, table string) ([]*Item, error)
	GetOne(ctx context.Context, table string, id int) (*Item, error)
	InsertOne(ctx context.Context, table string, item *Item) (*Item, error)
	UpdateOne(ctx context.Context, table string, item *Item) (*Item, error)
	DeleteOne(ctx context.Context, table string, id int) error
}

type CustomQueryStorage interface {
	QueryOne(ctx context.Context, table string, queryBuilder QueryBuilder) (*Item, error)
	QueryMany(ctx context.Context, table string, queryBuilder QueryBuilder) ([]*Item, error)
	Exec(ctx context.Context, expr exp.SQLExpression) error
}
