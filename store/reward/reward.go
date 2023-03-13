package reward

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"talkee/core"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

func New(db *sqlx.DB) *store {
	return &store{db: db}
}

type Config struct {
}

type store struct {
	db *sqlx.DB
}

//go:embed sql/insert_reward.sql
var stmtCreateReward string

//go:embed sql/update_reward.sql
var stmtUpdateReward string

//go:embed sql/get_unprocessed_task.sql
var stmtGetUnprocessedTasks string

//go:embed sql/update_reward_task.sql
var stmtUpdateRewardTask string

//go:embed sql/get_default_strategy.sql
var stmtGetDefaultStrategy string

//go:embed sql/get_created_reward.sql
var stmtGetCreatedRewardList string

//go:embed sql/get_reward_task_existance.sql
var StmtGetRewardtaskExistence string

//go:embed sql/insert_reward_task.sql
var StmtInsertRewardTask string

//go:embed sql/get_rewards_by_comment_ids.sql
var StmtGetRewardsByCommentIDs string

//go:embed sql/sum_rewards_by_asset.sql
var StmtSumRewardsByAsset string

func (s *store) CreateReward(ctx context.Context, models ...*core.Reward) error {

	_, err := s.db.NamedExec(stmtCreateReward, models)
	return err
}

func (s *store) FinishRewardTask(ctx context.Context, model *core.RewardTask, rewards ...*core.Reward) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.NamedExecContext(ctx, stmtCreateReward, rewards); err != nil {
		return err
	}

	if _, err = tx.Exec(stmtUpdateRewardTask, true, model.ID); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *store) UpdateReward(ctx context.Context, model *core.Reward) error {
	query, args, err := s.db.BindNamed(stmtUpdateReward, map[string]interface{}{
		"id":          model.ID,
		"trace_id":    model.TraceID,
		"snapshot_id": model.SnapshotID,
		"status":      model.Status,
	})

	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *store) FindUnprocessedList(ctx context.Context, before time.Time, limit int) ([]*core.RewardTask, error) {
	var tasks []*core.RewardTask
	if err := s.db.Select(&tasks, stmtGetUnprocessedTasks, before, limit); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	return tasks, nil
}

func (s *store) UpdateRewardTask(ctx context.Context, model *core.RewardTask) error {
	_, err := s.db.Exec(stmtUpdateRewardTask, model.Processed, model.ID)
	return err
}

func (s *store) GetDefaultStrategy() (*core.DefaultStrategy, error) {
	var rewardStrategy core.RewardStrategy
	if err := s.db.Get(&rewardStrategy, stmtGetDefaultStrategy); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if rewardStrategy.Params != "" {
		var strategy core.DefaultStrategy
		if err := json.Unmarshal([]byte(rewardStrategy.Params), &strategy); err == nil {
			return &strategy, nil
		}
	}

	return nil, nil
}

func (s *store) FindCreatedRewards(ctx context.Context, limit int) ([]*core.Reward, error) {
	var rewards []*core.Reward

	if err := s.db.Select(&rewards, stmtGetCreatedRewardList, limit); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	return rewards, nil

}

func (s *store) FindRewardsByCommentID(ctx context.Context, commentID ...uint) ([]*core.Reward, error) {
	var rewards []*core.Reward

	if err := s.db.Select(&rewards, StmtGetRewardsByCommentIDs, commentID, core.RewardObjectTypeComment); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	return rewards, nil

}

func (s *store) SumRewardsByAsset(ctx context.Context) (map[string]decimal.Decimal, error) {
	var outputs []struct {
		AssetID string          `db:"asset_id"`
		Amount  decimal.Decimal `db:"amount"`
	}
	if err := s.db.Select(&outputs, StmtSumRewardsByAsset); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	result := make(map[string]decimal.Decimal)
	for _, item := range outputs {
		result[item.AssetID] = item.Amount
	}
	return result, nil
}
