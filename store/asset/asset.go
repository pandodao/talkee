package asset

import (
	"context"
	_ "embed"
	"talkee/core"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) core.AssetStore {
	return &store{
		db: db,
	}
}

type store struct {
	db *sqlx.DB
}

func (s *store) GetAssets(ctx context.Context) ([]*core.Asset, error) {
	rows, err := s.db.QueryxContext(ctx, stmtGetAll, nil...)
	if err != nil {
		return nil, err
	}

	as, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	return as, nil
}

func (s *store) GetAsset(ctx context.Context, assetID string) (*core.Asset, error) {
	query, args, err := s.db.BindNamed(stmtGetByID, map[string]interface{}{
		"asset_id": assetID,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	asset := &core.Asset{}

	if err := scanRow(rows, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

func (s *store) SetAssets(ctx context.Context, assets []*core.Asset) error {

	tx := s.db.MustBegin()
	for _, asset := range assets {
		query, args, err := s.db.BindNamed(stmtUpdate, map[string]interface{}{
			"name":      asset.Name,
			"symbol":    asset.Symbol,
			"icon_url":  asset.IconURL,
			"price_usd": asset.PriceUSD,
			"asset_id":  asset.AssetID,
		})

		if err != nil {
			return err
		}

		tx.MustExecContext(ctx, query, args...)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *store) SetAsset(ctx context.Context, asset *core.Asset) error {

	query, args, err := s.db.BindNamed(stmtUpdate, map[string]interface{}{
		"name":      asset.Name,
		"symbol":    asset.Symbol,
		"icon_url":  asset.IconURL,
		"price_usd": asset.PriceUSD,
		"asset_id":  asset.AssetID,
	})

	if err != nil {
		return err
	}

	s.db.MustExecContext(ctx, query, args...)

	return nil
}

//go:embed sql/get_all.sql
var stmtGetAll string

//go:embed sql/get_by_id.sql
var stmtGetByID string

//go:embed sql/update.sql
var stmtUpdate string
