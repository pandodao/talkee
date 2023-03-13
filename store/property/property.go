package property

import (
	"context"
	_ "embed"
	"talkee/core"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) core.PropertyStore {
	return &store{
		db: db,
	}
}

type store struct {
	db *sqlx.DB
}

func (s *store) Get(ctx context.Context, key string) (core.PropertyValue, error) {
	query, args, err := s.db.BindNamed(stmtGet, map[string]interface{}{
		"key": key,
	})

	if err != nil {
		return "", err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return "", err
	}

	pp := &core.Property{}

	if err := scanRow(rows, pp); err != nil {
		return "", err
	}

	return pp.Value, nil
}

func (s *store) Set(ctx context.Context, key string, value interface{}) error {
	query, args, err := s.db.BindNamed(stmtSet, map[string]interface{}{
		"key":   key,
		"value": value,
	})

	if err != nil {
		return err
	}

	if _, err = s.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

//go:embed sql/get.sql
var stmtGet string

//go:embed sql/set.sql
var stmtSet string
