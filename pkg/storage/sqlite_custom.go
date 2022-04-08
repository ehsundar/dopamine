package storage

import (
	"context"
	"github.com/doug-martin/goqu/v9/exp"
)

func (s *sqlite) QueryOne(ctx context.Context, table string, queryBuilder QueryBuilder) (*Item, error) {
	i := Item{}
	_, err := queryBuilder(s.qu.From(table)).ScanStructContext(ctx, &i)
	return &i, err
}

func (s *sqlite) QueryMany(ctx context.Context, table string, queryBuilder QueryBuilder) ([]*Item, error) {
	var result []*Item
	err := queryBuilder(s.qu.From(table)).ScanStructsContext(ctx, &result)
	return result, err
}

func (s *sqlite) Exec(ctx context.Context, expr exp.SQLExpression) error {
	sql, args, err := expr.ToSQL()
	if err != nil {
		return err
	}
	_, err = s.qu.Db.ExecContext(ctx, sql, args...)
	return err
}
