package core

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

const (
	TipStrategyAvg  = "avg"
	TipStrategyTopN = "topn"

	TipTypeSlug     = "slug"
	TipTypeComments = "comments"
	TipTypeComment  = "comment"
	TipTypeUser     = "user"
)

const (
	TipStatusInit = iota
	TipStatusFilled
	TipStatusPending
	TipStatusDelivered
	TipStatusFailed
)

type (
	StrategyParams map[string]any
	Tip            struct {
		ID             uint64          `json:"id"`
		UUID           string          `json:"uuid"`
		UserID         uint64          `json:"user_id"`
		SiteID         uint64          `json:"site_id"`
		Slug           string          `json:"slug"`
		TipType        string          `json:"tip_type"`
		StrategyName   string          `json:"strategy_name"`
		StrategyParams StrategyParams  `gorm:"type:jsonb" json:"strategy_params"`
		AssetID        string          `json:"asset_id"`
		Amount         decimal.Decimal `json:"amount"`
		Memo           string          `json:"memo"`
		Status         int             `json:"status"`

		CreatedAt *time.Time `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at"`

		MixpayCode string `gorm:"-" json:"mixpay_code"`
	}

	TipStore interface {
		// INSERT INTO "tips"
		// (
		//  "uuid", "user_id",
		//   "site_id", "slug", "tip_type",
		//   "strategy_name", "strategy_params",
		//   "asset_id", "amount", "memo",
		//   "status", "created_at", "updated_at"
		// )
		// VALUES
		// (
		//   @tip.UUID, @tip.UserID,
		//   @tip.SiteID, @tip.Slug, @tip.TipType,
		//   @tip.StrategyName, @tip.StrategyParams,
		//   @tip.AssetID, @tip.Amount, @tip.Memo,
		//   0, NOW(), NOW()
		// )
		// RETURNING "id";
		CreateTip(ctx context.Context, tip *Tip) (uint64, error)

		// SELECT
		//  *
		// FROM "tips" WHERE
		//  "id" = @id AND "deleted_at" IS NULL
		// LIMIT 1;
		GetTip(ctx context.Context, id uint64) (*Tip, error)

		// SELECT
		//  *
		// FROM "tips" WHERE
		//  "uuid" = @uuid AND "deleted_at" IS NULL
		// LIMIT 1;
		GetTipByUUID(ctx context.Context, uuid string) (*Tip, error)

		// SELECT
		//  *
		// FROM "tips" WHERE
		//  "status" = @status AND "deleted_at" IS NULL
		// ORDER BY "id" DESC
		// LIMIT @limit
		GetTipsByStatus(ctx context.Context, status, limit int) ([]*Tip, error)

		// UPDATE "tips" SET
		// "status" = @status,
		// "updated_at" = NOW()
		// WHERE "id" = @id AND "deleted_at" IS NULL;
		UpdateTipStatus(ctx context.Context, id uint64, status int) error
	}

	TipService interface {
		CreateTip(ctx context.Context, tip *Tip, redirectURL string) (*Tip, error)
		FillTipByMixpay(ctx context.Context, tipUUID string) (*Tip, error)

		ProcessInitTip(ctx context.Context, tip *Tip) error
		ProcessCommentsType(ctx context.Context, tip *Tip) error
		ProcessPendingTip(ctx context.Context, tip *Tip) error
	}
)

func (t Tip) Validate() error {
	if t.SiteID == 0 {
		return ErrInvalidTipParams
	}
	if t.Slug == "" {
		return ErrInvalidTipParams
	}
	if t.AssetID == "" {
		return ErrInvalidTipParams
	}
	if t.Amount.IsZero() {
		return ErrInvalidTipParams
	}
	if t.TipType != TipTypeSlug && t.TipType != TipTypeComments && t.TipType != TipTypeComment && t.TipType != TipTypeUser {
		return ErrInvalidTipParams
	}
	if t.StrategyName != TipStrategyAvg && t.StrategyName != TipStrategyTopN {
		return ErrInvalidTipParams
	}
	return nil
}

func (a StrategyParams) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *StrategyParams) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, a)
}

func (a StrategyParams) ToMap() map[string]any {
	return a
}
