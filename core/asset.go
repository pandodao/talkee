package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type (
	Asset struct {
		AssetID   string          `json:"asset_id"`
		Name      string          `json:"name"`
		PriceUSD  decimal.Decimal `json:"price_usd"`
		Symbol    string          `json:"symbol"`
		IconURL   string          `json:"icon_url"`
		Order     int64           `json:"order"`
		CreatedAt *time.Time      `json:"created_at"`
		UpdatedAt *time.Time      `json:"updated_at"`
		DeletedAt *time.Time      `json:"-"`
	}

	AssetStore interface {
		// SELECT
		// * FROM @@table
		// WHERE deleted_at IS NULL;
		GetAssets(ctx context.Context) ([]*Asset, error)

		// SELECT
		// * FROM @@table
		// WEHRE
		//   asset_id = @assetID AND deleted_at IS NULL;
		GetAsset(ctx context.Context, assetID string) (*Asset, error)

		// INSERT INTO assets ("asset_id", "name", "symbol", "icon_url", "price_usd", "created_at", "updated_at")
		// VALUES (@asset.AssetID, @asset.Name, @asset.Symbol, @asset.IconURL, @asset.PriceUSD, NOW(), NOW())
		// ON CONFLICT (asset_id) DO
		// 	UPDATE SET
		// 		price_usd=EXCLUDED.price_usd,
		// 		name=EXCLUDED.name,
		// 		symbol=EXCLUDED.symbol,
		// 		icon_url=EXCLUDED.icon_url,
		// 		updated_at=NOW();
		SetAsset(ctx context.Context, asset *Asset) error
	}

	AssetService interface {
		UpdateAssets(ctx context.Context) error
		GetCachedAssets(ctx context.Context, assetID string) (*Asset, error)
	}
)
