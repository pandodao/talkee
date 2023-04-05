package tip_deliver

import (
	"context"
	"talkee/core"
	"time"

	"github.com/fox-one/pkg/logger"
)

const ()

type Worker struct {
	tips core.TipStore
	tipz core.TipService
}

func New(
	tips core.TipStore,
	tipz core.TipService,
) *Worker {
	return &Worker{
		tips: tips,
		tipz: tipz,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "tip_deliver")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond
	circle := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.run(ctx, circle); err != nil {
				log.Error(err)
				dur = time.Second * 60
				break
			}
			dur = time.Second * 6
			circle++
		}
	}
}

func (w *Worker) run(ctx context.Context, circle int) error {
	w.runFilled(ctx)

	if circle%10 == 0 {
		go w.runPending(ctx)
	}

	return nil
}

func (w *Worker) runFilled(ctx context.Context) error {
	tps, err := w.tips.GetTipsByStatus(ctx, core.TipStatusFilled, 100)
	if err != nil {
		return err
	}

	for _, tp := range tps {
		if err := w.tipz.ProcessFilledTip(ctx, tp); err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) runPending(ctx context.Context) error {
	tps, err := w.tips.GetTipsByStatus(ctx, core.TipStatusPending, 100)
	if err != nil {
		return err
	}

	for _, tp := range tps {
		if err := w.tipz.ProcessPendingTip(ctx, tp); err != nil {
			return err
		}
	}

	return nil
}
