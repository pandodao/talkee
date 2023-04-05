package reward

import (
	"context"

	"talkee/core"

	"github.com/fox-one/mixin-sdk-go"
)

func New(
	rewards core.RewardStore,
	rewardTasks core.RewardTaskStore,
	rewardStrategys core.RewardStrategyStore,
	favors core.FavouriteStore,
	comments core.CommentStore,
	client *mixin.Client,
	cfg Config,
) *service {
	return &service{
		rewards:         rewards,
		rewardTasks:     rewardTasks,
		rewardStrategys: rewardStrategys,
		favors:          favors,
		comments:        comments,
		client:          client,
		cfg:             cfg,
	}
}

type Config struct {
	Pin string
}

type service struct {
	rewards         core.RewardStore
	rewardTasks     core.RewardTaskStore
	comments        core.CommentStore
	rewardStrategys core.RewardStrategyStore
	favors          core.FavouriteStore
	client          *mixin.Client
	cfg             Config
}

// func (s *service) GetUnprocessedTasks(ctx context.Context, before time.Time) ([]*core.RewardTask, error) {
// 	return s.rewardTasks.FindUnprocessedList(ctx, before, 100)
// }

// func (s *service) ProcessRewardTask(ctx context.Context, model *core.RewardTask) error {
// 	defaultStrategy, err := s.rewardStrategys.GetDefaultStrategy()
// 	if err != nil {
// 		return err
// 	}

// 	if defaultStrategy == nil {
// 		return fmt.Errorf("default reward strategy not exist")
// 	}

// 	comments, err := s.comments.GetAllCommentsBySiteSlug(ctx, model.SiteID, model.Slug)
// 	if err != nil {
// 		return err
// 	}

// 	sort.SliceStable(comments, func(i, j int) bool {
// 		return comments[i].FavorCount > comments[j].FavorCount
// 	})

// 	rewardedUsers := make(map[string]bool)
// 	var rewards []*core.Reward
// 	length := int(math.Min(float64(defaultStrategy.TopCommentCount), float64(len(comments))))
// 	for i := 0; i < length; i++ {
// 		cmt := comments[i]
// 		r := &core.Reward{
// 			ObjectType:  core.RewardObjectTypeComment,
// 			ObjectID:    cmt.ID,
// 			SiteID:      cmt.SiteID,
// 			RecipientID: cmt.Creator.MixinUserID,
// 			Status:      core.RewardStatusCreated,
// 			Amount:      defaultStrategy.MinRewardAmount.Mul(defaultStrategy.Ratio),
// 			AssetID:     defaultStrategy.RewardAssetID,
// 		}

// 		if decimal.NewFromInt(int64(cmt.FavorCount)).GreaterThan(defaultStrategy.MinRewardAmount) {
// 			r.Amount = decimal.NewFromInt(int64(cmt.FavorCount)).Mul(defaultStrategy.Ratio)
// 		}
// 		rewards = append(rewards, r)
// 		rewardedUsers[cmt.Creator.MixinUserID] = true
// 	}

// 	if length == defaultStrategy.TopCommentCount {
// 		cmt := comments[defaultStrategy.TopCommentCount-1]
// 		if cmt.FavorCount > 0 {
// 			favors, err := s.favors.FindAllFavourites(ctx, cmt.ID)
// 			if err != nil {
// 				return err
// 			}
// 			if len(favors) > 0 {
// 				s := rand.NewSource(time.Now().Unix())
// 				r := rand.New(s)
// 				luckyFavor := favors[r.Intn(len(favors))]
// 				rewards = append(rewards, &core.Reward{
// 					ObjectType:  core.RewardObjectTypeFavour,
// 					ObjectID:    luckyFavor.ID,
// 					SiteID:      cmt.SiteID,
// 					RecipientID: luckyFavor.Creator.MixinUserID,
// 					Amount:      decimal.NewFromInt(int64(cmt.FavorCount)).Mul(defaultStrategy.Ratio),
// 					Status:      core.RewardStatusCreated,
// 					AssetID:     defaultStrategy.RewardAssetID,
// 				})
// 			}
// 		}
// 	}

// 	if len(comments) > defaultStrategy.TopCommentCount {
// 		restComments := comments[defaultStrategy.TopCommentCount:]
// 		s := rand.NewSource(time.Now().Unix())
// 		r := rand.New(s)
// 		cmt := restComments[r.Intn(len(restComments))]
// 		rewards = append(rewards, &core.Reward{
// 			ObjectType:  core.RewardObjectTypeComment,
// 			ObjectID:    cmt.ID,
// 			SiteID:      cmt.SiteID,
// 			RecipientID: cmt.Creator.MixinUserID,
// 			//set amount to 777 satoshi
// 			Amount:  defaultStrategy.RewardAmount777.Mul(defaultStrategy.Ratio),
// 			Status:  core.RewardStatusCreated,
// 			AssetID: defaultStrategy.RewardAssetID,
// 		})
// 		rewardedUsers[cmt.Creator.MixinUserID] = true
// 	}

// 	if len(comments) > defaultStrategy.TopCommentCount {
// 		for _, cmt := range comments[defaultStrategy.TopCommentCount:] {
// 			_, alreadyRewarded := rewardedUsers[cmt.Creator.MixinUserID]
// 			if !alreadyRewarded {
// 				rewards = append(rewards, &core.Reward{
// 					ObjectType:  core.RewardObjectTypeComment,
// 					ObjectID:    cmt.ID,
// 					SiteID:      cmt.SiteID,
// 					RecipientID: cmt.Creator.MixinUserID,
// 					Amount:      defaultStrategy.MinRewardAmount.Mul(defaultStrategy.Ratio),
// 					Status:      core.RewardStatusCreated,
// 					AssetID:     defaultStrategy.RewardAssetID,
// 				})
// 				rewardedUsers[cmt.Creator.MixinUserID] = true
// 			}
// 		}
// 	}

// 	if len(rewards) > 0 {
// 		if err := s.FinishRewardTask(ctx, model, rewards); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (s *service) GetCreatedRewards(ctx context.Context, limit int) ([]*core.Reward, error) {
	return s.rewards.FindCreatedRewards(ctx, limit)
}

func (s *service) TransferReward(ctx context.Context, model *core.Reward) error {
	// @TODO
	// traceID := uuid.Modify(model.RecipientID, model.ObjectType+"-reward-"+strconv.FormatUint(model.ObjectID, 10)+model.AssetID)

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

// func (s *service) FinishRewardTask(ctx context.Context, model *core.RewardTask, rewards []*core.Reward) error {
// 	if err := store.Transaction(func(tx *store.Handler) error {
// 		rewardStore := reward.New(tx)
// 		rewardTaskStore := rewardTask.New(tx)

// 		for _, re := range rewards {
// 			if err := rewardStore.CreateReward(ctx, re); err != nil {
// 				return err
// 			}

// 			if err := rewardTaskStore.UpdateRewardTask(ctx, model); err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	}); err != nil {
// 		return err
// 	}
// 	return nil
// }
