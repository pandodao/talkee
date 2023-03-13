package reward_processer

import (
	"context"
	"talkee/core"
	"time"

	"github.com/fox-one/pkg/logger"
)

const ()

type Worker struct {
	rewardz core.RewardService
}

func New(
	rewardz core.RewardService,
) *Worker {
	return &Worker{
		rewardz: rewardz,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "reward_processer")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.run(ctx); err != nil {
				log.Error(err)
				dur = time.Second * 10
				break
			}
			dur = time.Second * 3
		}
	}
}

func (w *Worker) run(ctx context.Context) error {
	rewards, err := w.rewardz.GetCreatedRewards(ctx, 100)
	if err != nil {
		return err
	}

	for _, r := range rewards {
		if err := w.rewardz.TransferReward(ctx, r); err != nil {
			return err
		}
	}

	return nil
}
