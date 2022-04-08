package storage

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

func (s *sqlite) QueryOne(ctx context.Context, expr *goqu.SelectDataset) (*Item, error) {
	i := Item{}
	_, err := expr.ScanStructContext(ctx, &i)
	return &i, err
}

func (s *sqlite) QueryMany(ctx context.Context, expr *goqu.SelectDataset) ([]*Item, error) {
	var result []*Item
	err := expr.ScanStructsContext(ctx, &result)
	return result, err
}

func (s *sqlite) Exec(ctx context.Context, expr exp.SQLExpression) error {
	sql, args, err := expr.ToSQL()
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, sql, args...)
	return err
}
