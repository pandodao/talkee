package attachment_processer

import (
	"context"
	"talkee/core"
	"time"

	"github.com/fox-one/pkg/logger"
)

const ()

type Worker struct {
	attachmentz core.AttachmentService
}

func New(
	attachmentz core.AttachmentService,
) *Worker {
	return &Worker{
		attachmentz: attachmentz,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "attachment_processer")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			w.run(ctx)
			dur = time.Second * 10
		}
	}
}

func (w *Worker) run(ctx context.Context) {
	// TODO: process attachment
}
