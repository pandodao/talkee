package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type (
	Asset struct {
		AssetID   string          `db:"asset_id" json:"asset_id"`
		Name      string          `db:"name" json:"name"`
		PriceUSD  decimal.Decimal `db:"price_usd" json:"price_usd"`
		Symbol    string          `db:"symbol" json:"symbol"`
		IconURL   string          `db:"icon_url" json:"icon_url"`
		Order     int64           `db:"order" json:"order"`
		CreatedAt *time.Time      `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time      `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time      `db:"deleted_at" json:"-"`
	}

	AssetStore interface {
		GetAssets(ctx context.Context) ([]*Asset, error)
		GetAsset(ctx context.Context, assetID string) (*Asset, error)
		SetAssets(ctx context.Context, assets []*Asset) error
		SetAsset(ctx context.Context, asset *Asset) error
	}

	AssetService interface {
		UpdateAssets(ctx context.Context) error
		GetCachedAssets(ctx context.Context, assetID string) (*Asset, error)
	}
)
