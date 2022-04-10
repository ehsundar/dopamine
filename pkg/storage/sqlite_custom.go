package storage

import (
	"context"
	"github.com/doug-martin/goqu/v9/exp"
)

func (s *sqlite) QueryOne(ctx context.Context, table string, queryBuilder QueryBuilder) (*Item, error) {
	i := Item{}
	_, err := queryBuilder(s.qu.From(table)).ScanStructContext(ctx, &i)
	if err != nil {
		return nil, err
	}

	err = i.LoadContentsMap(false)
	return &i, err
}

func (s *sqlite) QueryMany(ctx context.Context, table string, queryBuilder QueryBuilder) ([]*Item, error) {
	var result []*Item
	err := queryBuilder(s.qu.From(table)).ScanStructsContext(ctx, &result)
	if err != nil {
		return nil, err
	}

	for _, r := range result {
		err = r.LoadContentsMap(false)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *sqlite) Exec(ctx context.Context, expr exp.SQLExpression) error {
	sql, args, err := expr.ToSQL()
	if err != nil {
		return err
	}
	_, err = s.qu.Db.ExecContext(ctx, sql, args...)
	return err
}
