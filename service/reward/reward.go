package reward

import (
	"context"

	"talkee/core"

	"github.com/fox-one/mixin-sdk-go"
)

func New(
	rewards core.RewardStore,
	favors core.FavouriteStore,
	comments core.CommentStore,
	client *mixin.Client,
	cfg Config,
) *service {
	return &service{
		rewards:  rewards,
		favors:   favors,
		comments: comments,
		client:   client,
		cfg:      cfg,
	}
}

type Config struct {
	Pin string
}

type service struct {
	rewards  core.RewardStore
	comments core.CommentStore
	favors   core.FavouriteStore
	client   *mixin.Client
	cfg      Config
}

func (s *service) GetCreatedRewards(ctx context.Context, limit int) ([]*core.Reward, error) {
	return s.rewards.FindCreatedRewards(ctx, limit)
}

func (s *service) TransferReward(ctx context.Context, model *core.Reward) error {
	if model.RecipientID == s.client.ClientID {
		// ignore self transfer
		model.Status = core.RewardStatusRewarded
		return s.rewards.UpdateReward(ctx, model)
	}

	input := mixin.TransferInput{
		AssetID:    model.AssetID,
		OpponentID: model.RecipientID,
		Amount:     model.Amount,
		TraceID:    model.TraceID,
		Memo:       model.Memo,
	}

	snapshot, err := s.client.Transfer(ctx, &input, s.cfg.Pin)
	if err != nil {
		return err
	}

	model.SnapshotID = snapshot.SnapshotID
	model.Status = core.RewardStatusRewarded

	return s.rewards.UpdateReward(ctx, model)
}
