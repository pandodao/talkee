package reward_task_processer

import (
	"context"
	"talkee/core"
	"time"

	"github.com/fox-one/pkg/logger"
)

const ()

type Worker struct {
	rewardTaskz core.RewardTaskService
}

func New(
	rewardTaskz core.RewardTaskService,
) *Worker {
	return &Worker{
		rewardTaskz: rewardTaskz,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "reward_task_processer")
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
			dur = time.Second * 5
		}
	}
}

func (w *Worker) run(ctx context.Context) error {
	tasks, err := w.rewardTaskz.GetUnprocessedTasks(ctx, time.Now().Add(time.Hour*(-12)))
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if err := w.rewardTaskz.ProcessRewardTask(ctx, task); err != nil {
			return err
		}
	}

	return nil
}
