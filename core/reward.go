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
		ID         uint64     `db:"id" json:"id"`
		SiteID     uint64     `db:"site_id" json:"site_id"`
		Slug       string     `db:"slug" json:"slug"`
		Processed  bool       `db:"processed" json:"processed"`
		StrategyID uint64     `db:"strategy_id" json:"strategy_id"`
		CreatedAt  *time.Time `db:"created_at" json:"created_at"`
		UpdatedAt  *time.Time `db:"updated_at" json:"updated_at"`
	}

	Reward struct {
		ID          uint64          `db:"id" json:"id"`
		ObjectType  string          `db:"object_type" json:"object_type"`
		ObjectID    uint64          `db:"object_id" json:"object_id"`
		SiteID      uint64          `db:"site_id" json:"site_id"`
		RecipientID string          `db:"recipient_id" json:"recipient_id"`
		TraceID     string          `db:"trace_id" json:"trace_id"`
		SnapshotID  string          `db:"snapshot_id" json:"snapshot_id"`
		AssetID     string          `db:"asset_id" json:"asset_id"`
		Amount      decimal.Decimal `db:"amount" json:"amount"`
		Status      string          `db:"status" json:"status"`
		CreatedAt   *time.Time      `db:"created_at" json:"created_at"`
		UpdatedAt   *time.Time      `db:"updated_at" json:"updated_at"`
	}

	RewardStore interface {
		SumRewardsByAsset(ctx context.Context) (map[string]decimal.Decimal, error)
		CreateReward(ctx context.Context, models ...*Reward) error
		UpdateReward(ctx context.Context, model *Reward) error
		FindCreatedRewards(ctx context.Context, limit int) ([]*Reward, error)
		FindRewardsByCommentID(ctx context.Context, commentID ...uint) ([]*Reward, error)
	}

	RewardTaskStore interface {
		UpdateRewardTask(ctx context.Context, model *RewardTask) error
		FinishRewardTask(ctx context.Context, model *RewardTask, rewards ...*Reward) error
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
		ProcessRewardTask(ctx context.Context, model *RewardTask) error
		GetUnprocessedTasks(ctx context.Context, before time.Time) ([]*RewardTask, error)
	}
)
