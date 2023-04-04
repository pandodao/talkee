package core

import (
	"context"
	"database/sql/driver"
	"time"

	"github.com/lib/pq"
)

type SiteOrigins []string

func (s *SiteOrigins) Scan(value interface{}) error {
	stringArray := pq.StringArray{}
	err := stringArray.Scan(value)
	if err != nil {
		return err
	}
	*s = SiteOrigins(stringArray)
	return nil
}

func (s SiteOrigins) Value() (driver.Value, error) {
	stringArray := pq.StringArray(s)
	if stringArray == nil {
		return nil, nil
	}
	return stringArray.Value()
}

type (
	Site struct {
		ID             uint64      `json:"id"`
		UserID         uint64      `json:"user_id"`
		Origins        SiteOrigins `gorm:"type:text[]" json:"origins"`
		Origin         string      `json:"origin"`
		Name           string      `json:"name"`
		UseArweave     bool        `json:"use_arweave"`
		RewardStrategy int         `json:"reward_strategy"`

		CreatedAt *time.Time `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
		DeletedAt *time.Time `json:"-"`
	}

	SiteStore interface {
		// SELECT
		// 	*
		// FROM @@table
		// WHERE
		// 	"id" = @siteID
		// AND deleted_at IS NULL
		// LIMIT 1;
		GetSite(ctx context.Context, siteID uint64) (*Site, error)

		// SELECT
		// 	*
		// FROM @@table
		// WHERE
		//  @origin=ANY("origins")
		// AND
		// 	deleted_at IS NULL
		// LIMIT 1;
		GetSiteByOrigin(ctx context.Context, origin string) (*Site, error)

		// SELECT
		// 	*
		// FROM @@table
		// WHERE
		// 	"user_id" = @userID
		// AND
		// 	deleted_at IS NULL;
		GetSitesByUser(ctx context.Context, userID uint64) ([]*Site, error)

		// INSERT INTO @@table
		// 	(
		// 		"user_id",
		// 		"name",
		// 		"origins",
		// 		"use_arweave",
		// 		"reward_strategy",
		// 		"created_at", "updated_at"
		// 	)
		// VALUES
		// 	(
		// 		@userID,
		// 		@name,
		// 		@origins,
		// 		@useArweave,
		// 		0,
		// 		NOW(), NOW()
		// 	)
		// RETURNING id
		// ;
		CreateSite(ctx context.Context, userID uint64, name string, origins SiteOrigins, useArweave bool) (uint64, error)

		// UPDATE @@table
		// 	{{set}}
		// 	  "name"=@site.Name,
		// 	  "origins"=@site.Origins,
		// 	  "use_arweave"=@site.UseArweave,
		// 		"updated_at"=NOW()
		// 	{{end}}
		// WHERE
		// 	 "id" = @id AND "deleted_at" IS NULL
		// ;
		UpdateSite(ctx context.Context, id uint64, site *Site) error
	}

	SiteService interface {
	}
)
