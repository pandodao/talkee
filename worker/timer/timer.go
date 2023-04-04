package timer

import (
	"context"
	"time"

	"talkee/core"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pkg/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
)

const (
	statGlobalSummary = "stat:global_summary"
)

type (
	Config struct {
	}

	Worker struct {
		cfg        Config
		propertys  core.PropertyStore
		comments   core.CommentStore
		favourites core.FavouriteStore
		rewards    core.RewardStore
		assets     core.AssetStore
		assetz     core.AssetService
	}
)

func New(
	cfg Config,
	propertys core.PropertyStore,
	comments core.CommentStore,
	favourites core.FavouriteStore,
	rewards core.RewardStore,
	assets core.AssetStore,
	assetz core.AssetService,
) *Worker {
	if _, err := govalidator.ValidateStruct(cfg); err != nil {
		panic(err)
	}

	return &Worker{
		cfg:        cfg,
		propertys:  propertys,
		comments:   comments,
		favourites: favourites,
		rewards:    rewards,
		assets:     assets,
		assetz:     assetz,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "timer")
	ctx = logger.WithContext(ctx, log)

	log.Println("timer worker started")
	dur := time.Millisecond
	var circle int64
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.run(ctx, circle); err == nil {
				dur = time.Second
				circle += 1
			} else {
				dur = 10 * time.Second
				circle += 10
			}
		}
	}
}

func (w *Worker) run(ctx context.Context, circle int64) error {

	if circle%60 == 0 {
		if err := w.updateAssets(ctx); err != nil {
			return err
		}
	}

	if circle%3600 == 0 {
		if err := w.updateGlobalSummary(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) updateAssets(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "timer")

	// update assets and fiats
	if err := w.assetz.UpdateAssets(ctx); err != nil {
		log.Warn("update assets error", err)
	}

	return nil
}

func (w *Worker) updateGlobalSummary(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "timer")

	commentCount, err := w.comments.CountComments(ctx, 0, "")
	if err != nil {
		log.WithError(err).Warn("comments.CountComments")
		return err
	}

	totalUsd := decimal.Zero
	rewardSumMap, err := w.rewards.SumRewardsByAsset(ctx)
	if err != nil {
		log.WithError(err).Warn("rewards.SumRewardsByAsset")
		return err
	}
	for assetID, _amount := range rewardSumMap {
		amount := _amount.(decimal.Decimal)
		asset, err := w.assets.GetAsset(ctx, assetID)
		if err != nil {
			log.WithError(err).Warn("assets.GetAsset")
			continue
		}
		totalUsd = totalUsd.Add(amount.Mul(asset.PriceUSD))
	}

	favCounts, err := w.favourites.CountAllFavourites(ctx)
	if err != nil {
		log.WithError(err).Warn("favourites.CountAllFavourites")
		return err
	}

	// add the numbers from links_talkee
	commentCount += 735325
	favCounts += 1660076
	totalUsd = totalUsd.Add(decimal.NewFromFloat(2121.87))

	summaryData := map[string]interface{}{
		"total_comments":    commentCount,
		"total_favors":      favCounts,
		"total_rewards_usd": totalUsd.StringFixed(2),
	}

	buf, err := jsoniter.Marshal(summaryData)
	if err != nil {
		log.WithError(err).Errorln("jsoniter.Marshal")
		return err
	}

	if _, err := w.propertys.Set(ctx, statGlobalSummary, string(buf)); err != nil {
		log.WithError(err).Errorln("properties.Set")
		return err
	}
	return nil
}
