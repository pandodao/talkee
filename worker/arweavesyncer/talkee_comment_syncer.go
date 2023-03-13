package arweavesyncer

import (
	"context"
	"talkee/core"
	"time"

	"github.com/fox-one/pkg/logger"
)

const ()

type Worker struct {
	commnetz core.CommentService
}

func New(
	commnetz core.CommentService,
) *Worker {
	return &Worker{
		commnetz: commnetz,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "talkeeCommentSyncerWorker")
	ctx = logger.WithContext(ctx, log)

	log.Info("talkee comment syncer worker start")
	dur := time.Second * 10
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			err := w.syncCommentToArweave(ctx)
			if err != nil {
				log.Error(err)
				dur = time.Second * 10
				break
			}
			dur = time.Second * 5
		}

	}
}

func (w *Worker) syncCommentToArweave(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "syncCommentToArweave")
	comments, err := w.commnetz.FindArweaveSyncList(ctx, 100)
	if err != nil {
		return err
	}

	for _, comment := range comments {
		if err := w.commnetz.SyncToArweave(ctx, comment); err != nil {
			log.Error(err)
		}
	}

	return nil
}
