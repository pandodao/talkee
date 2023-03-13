package core

import (
	"context"
	"time"
)

type (
	Site struct {
		ID             uint64   `db:"id" json:"id"`
		UserID         uint64   `db:"user_id" json:"user_id"`
		Origins        []string `db:"origins" json:"origins"`
		Origin         string   `db:"origin" json:"origin"`
		Name           string   `db:"name" json:"name"`
		UseArweave     bool     `db:"use_arweave" json:"use_arweave"`
		RewardStrategy int      `db:"reward_strategy" json:"reward_strategy"`

		CreatedAt *time.Time `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"-"`
	}

	SiteStore interface {
		GetSite(ctx context.Context, siteID uint64) (*Site, error)
		GetSiteByOrigin(ctx context.Context, origin string) (*Site, error)
		GetSitesByUser(ctx context.Context, userID uint64) ([]*Site, error)
		CreateSite(ctx context.Context, userID uint64, name string, origin []string, useArweave bool) (uint64, error)
		UpdateSite(ctx context.Context, id uint64, data map[string]interface{}) error
	}

	SiteService interface {
	}
)
