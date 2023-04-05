package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

var (
	RewardStatusCreated  = "created"
	RewardStatusRewarded = "rewarded"

	RewardObjectTypeSlug    = "slug"
	RewardObjectTypeComment = "comment"
	RewardObjectTypeFavour  = "favour"
)

type (
	Reward struct {
		ID          uint64          `json:"id"`
		TipID       uint64          `json:"tip_id"`
		ObjectType  string          `json:"object_type"`
		ObjectID    uint64          `json:"object_id"`
		SiteID      uint64          `json:"site_id"`
		RecipientID string          `json:"recipient_id"`
		TraceID     string          `json:"trace_id"`
		SnapshotID  string          `json:"snapshot_id"`
		AssetID     string          `json:"asset_id"`
		Amount      decimal.Decimal `json:"amount"`
		Memo        string          `json:"memo"`
		Status      string          `json:"status"`
		CreatedAt   *time.Time      `json:"created_at"`
		UpdatedAt   *time.Time      `json:"updated_at"`
	}

	SumRewardItem struct {
		AssetID string          `json:"asset_id"`
		Amount  decimal.Decimal `json:"amount"`
	}

	RewardStore interface {
		// SELECT
		// 	"asset_id",
		// 	SUM("amount") as "amount"
		// FROM
		// 	"rewards"
		// GROUP BY "asset_id";
		SumRewardsByAsset(ctx context.Context) ([]*SumRewardItem, error)

		// INSERT INTO "rewards"
		// 	(
		//    "tip_id",
		// 		"object_type",
		// 		"object_id",
		// 		"site_id",
		// 		"recipient_id",
		// 		"trace_id",
		// 		"snapshot_id",
		// 		"asset_id",
		// 		"amount",
		//    "memo",
		// 		"status",
		// 		"created_at",
		// 		"updated_at"
		// 	)
		// VALUES
		// 	(
		// 		@model.TipID,
		// 		@model.ObjectType,
		// 		@model.ObjectID,
		// 		@model.SiteID,
		// 		@model.RecipientID,
		// 		@model.TraceID,
		// 		@model.SnapshotID,
		// 		@model.AssetID,
		// 		@model.Amount,
		//    @model.Memo,
		// 		@model.Status,
		// 		NOW(), NOW()
		// 	);
		CreateReward(ctx context.Context, model *Reward) error

		// UPDATE
		// 	"rewards"
		// SET
		// 	"status" = @model.Status,
		// 	"snapshot_id"= @model.SnapshotID,
		// 	"trace_id"= @model.TraceID,
		// 	"updated_at" = NOW()
		// WHERE
		// 	"id" = @model.ID;
		UpdateReward(ctx context.Context, model *Reward) error

		// SELECT
		//  *
		// FROM
		//  "rewards"
		// WHERE
		//  "status" = 'created'
		// ORDER BY "id" asc
		// LIMIT @limit;
		FindCreatedRewards(ctx context.Context, limit int) ([]*Reward, error)

		// SELECT
		//  *
		// FROM
		//  "rewards"
		// WHERE
		//  "tip_id" = @tipID AND "status" = @status
		// ORDER BY "id" asc;
		// ;
		GetRewardsByTipIDAndStatus(ctx context.Context, tipID uint64, status string) ([]*Reward, error)

		// SELECT
		//  *
		// FROM
		//  "rewards"
		// WHERE
		// "object_id" in (@commentIDs)
		// AND
		// "object_type" = 'comment'
		// ORDER BY "id" asc;
		FindRewardsByCommentIDs(ctx context.Context, commentIDs []uint64) ([]*Reward, error)
	}

	RewardService interface {
		GetCreatedRewards(ctx context.Context, limit int) ([]*Reward, error)
		TransferReward(ctx context.Context, model *Reward) error
	}
)
