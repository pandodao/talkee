package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

var (
	RewardStatusCreated     = "created"
	RewardStatusRewarded    = "rewarded"
	RewardObjectTypeComment = "comment"
	RewardObjectTypeFavour  = "favour"
)

type RewardStrategy struct {
	ID        uint64 `db:"id" json:"id"`
	Params    string `db:"params" json:"params"`
	IsDefault bool   `db:"is_default" json:"is_default"`
}

type DefaultStrategy struct {
	TimeLimit       string          `json:"time_limit,omitempty"`
	TopCommentCount int             `json:"top_comment_count,omitempty"`
	MinRewardAmount decimal.Decimal `json:"min_reward_amount,omitempty"`
	RewardAmount777 decimal.Decimal `json:"reward_amount_777,omitempty"`
	RewardAssetID   string          `json:"reward_asset_id,omitempty"`
	Ratio           decimal.Decimal `json:"ratio,omitempty"`
}

type (
	RewardTask struct {
		ID         uint64     `json:"id"`
		SiteID     uint64     `json:"site_id"`
		Slug       string     `json:"slug"`
		Processed  bool       `json:"processed"`
		StrategyID uint64     `json:"strategy_id"`
		CreatedAt  *time.Time `json:"created_at"`
		UpdatedAt  *time.Time `json:"updated_at"`
	}

	Reward struct {
		ID          uint64          `json:"id"`
		ObjectType  string          `json:"object_type"`
		ObjectID    uint64          `json:"object_id"`
		SiteID      uint64          `json:"site_id"`
		RecipientID string          `json:"recipient_id"`
		TraceID     string          `json:"trace_id"`
		SnapshotID  string          `json:"snapshot_id"`
		AssetID     string          `json:"asset_id"`
		Amount      decimal.Decimal `json:"amount"`
		Status      string          `json:"status"`
		CreatedAt   *time.Time      `json:"created_at"`
		UpdatedAt   *time.Time      `json:"updated_at"`
	}

	RewardStore interface {
		// SELECT
		// 	"asset_id",
		// 	SUM("amount") as "amount"
		// FROM
		// 	"rewards"
		// GROUP BY "asset_id";
		SumRewardsByAsset(ctx context.Context) (map[string]interface{}, error)

		// INSERT INTO "rewards"
		// 	(
		// 		"object_type",
		// 		"object_id",
		// 		"site_id",
		// 		"recipient_id",
		// 		"trace_id",
		// 		"snapshot_id",
		// 		"asset_id",
		// 		"amount",
		// 		"status",
		// 		"created_at",
		// 		"updated_at"
		// 	)
		// VALUES
		// 	(
		// 		@model.ObjectType,
		// 		@model.ObjectID,
		// 		@model.SiteID,
		// 		@model.RecipientID,
		// 		@model.TraceID,
		// 		@model.SnapshotID,
		// 		@model.AssetID,
		// 		@model.Amount,
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
		// 	"id" = :id;
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
		// "object_id" in (@commentIDs)
		// AND
		// "object_type" = 'comment'
		// ORDER BY "id" asc;
		FindRewardsByCommentIDs(ctx context.Context, commentIDs []uint64) ([]*Reward, error)
	}

	RewardTaskStore interface {
		// UPDATE
		// 	"reward_tasks"
		// SET
		// 	"processed" = @model.Processed,
		// 	"updated_at" = NOW()
		// WHERE
		// 	"id" = @model.ID;
		// ;
		UpdateRewardTask(ctx context.Context, model *RewardTask) error

		// SELECT
		//  *
		// FROM "reward_tasks"
		// WHERE
		// 	"created_at" < @before
		// 	AND "processed" = false
		// LIMIT @limit;
		FindUnprocessedList(ctx context.Context, before time.Time, limit int) ([]*RewardTask, error)
	}

	RewardStrategyStore interface {
		GetDefaultStrategy() (*DefaultStrategy, error)
	}

	RewardService interface {
		GetCreatedRewards(ctx context.Context, limit int) ([]*Reward, error)
		TransferReward(ctx context.Context, model *Reward) error
	}

	RewardTaskService interface {
		FinishRewardTask(ctx context.Context, model *RewardTask, rewards []*Reward) error
		ProcessRewardTask(ctx context.Context, model *RewardTask) error
		GetUnprocessedTasks(ctx context.Context, before time.Time) ([]*RewardTask, error)
	}
)
