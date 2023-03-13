package site

import (
	"context"
	_ "embed"
	"talkee/core"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *store {
	return &store{db: db}
}

type Config struct {
}

type store struct {
	db *sqlx.DB
}

//go:embed sql/get_site.sql
var stmtGetSite string

//go:embed sql/get_site_by_origin.sql
var stmtGetSiteByOrigin string

//go:embed sql/get_site_by_user.sql
var stmtGetSiteByUser string

//go:embed sql/create_site.sql
var stmtCreateSite string

//go:embed sql/update_site.sql
var stmtUpdateSite string

func (s *store) GetSite(ctx context.Context, id uint64) (*core.Site, error) {
	query, args, err := s.db.BindNamed(stmtGetSite, map[string]interface{}{
		"id": id,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	site, err := scanRow(rows)
	if err != nil {
		return nil, err
	}

	return site, nil
}

func (s *store) GetSiteByOrigin(ctx context.Context, origin string) (*core.Site, error) {
	query, args, err := s.db.BindNamed(stmtGetSiteByOrigin, map[string]interface{}{
		"origin": origin,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	site, err := scanRow(rows)
	if err != nil {
		return nil, err
	}

	return site, nil
}

func (s *store) GetSitesByUser(ctx context.Context, userID uint64) ([]*core.Site, error) {
	query, args, err := s.db.BindNamed(stmtGetSiteByUser, map[string]interface{}{
		"user_id": userID,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	sites, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	return sites, nil
}

func (s *store) CreateSite(ctx context.Context, userID uint64, name string, origins []string, useArweave bool) (uint64, error) {
	query, args, err := s.db.BindNamed(stmtCreateSite, map[string]interface{}{
		"user_id":         userID,
		"name":            name,
		"origins":         pq.Array(origins),
		"use_arweave":     useArweave,
		"reward_strategy": 0,
	})

	if err != nil {
		return 0, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	id, err := scanReturnID(rows)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *store) UpdateSite(ctx context.Context, id uint64, data map[string]interface{}) error {
	query, args, err := s.db.BindNamed(stmtUpdateSite, map[string]interface{}{
		"id":          id,
		"name":        data["name"],
		"origins":     pq.Array(data["origins"]),
		"use_arweave": data["use_arweave"],
	})

	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
