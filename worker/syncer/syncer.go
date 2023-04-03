package syncer

import (
	"context"
	"fmt"
	"time"

	"talkee/core"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pkg/logger"
	"github.com/patrickmn/go-cache"
)

const (
	snapshotCheckpoint = "syncer:snapshot_checkpoint"

	snapshotKey = "snapshot:%s"
)

type (
	Config struct {
		ClientID string `valid:"required"`
	}

	Worker struct {
		cfg       Config
		propertys core.PropertyStore
		snapshots core.SnapshotStore
		users     core.UserStore
		snapshotz core.SnapshotService
		cache     *cache.Cache
	}
)

func New(
	cfg Config,
	propertys core.PropertyStore,
	snapshots core.SnapshotStore,
	users core.UserStore,

	snapshotz core.SnapshotService,
) *Worker {
	if _, err := govalidator.ValidateStruct(cfg); err != nil {
		panic(err)
	}

	return &Worker{
		cfg:       cfg,
		propertys: propertys,
		snapshots: snapshots,
		users:     users,
		snapshotz: snapshotz,
		cache:     cache.New(time.Hour, time.Hour),
	}
}

func (w *Worker) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "syncer")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.run(ctx); err == nil {
				dur = 100 * time.Millisecond
			} else {
				dur = time.Second
			}
		}
	}
}

func (w *Worker) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	v, err := w.propertys.Get(ctx, snapshotCheckpoint)
	if err != nil {
		log.WithError(err).Errorln("properties.Get")
		return err
	}
	var (
		offset    = v.Time()
		newOffset = offset
	)

	const LIMIT = 500

	var snapshots []*core.Snapshot

	{
		items, err := w.snapshotz.PollSnapshots(ctx, offset, LIMIT)
		if err != nil {
			log.WithError(err).Errorln("list snapshots")
			return err
		}
		// log.Println("Pull snapshot offset:", offset)

		for _, snapshot := range items {
			newOffset = snapshot.CreatedAt
			if _, ok := w.cache.Get(fmt.Sprintf(snapshotKey, snapshot.SnapshotID)); ok {
				continue
			}
			// change your condition here
			if snapshot.UserID == w.cfg.ClientID {
				snapshots = append(snapshots, snapshot)
			}
		}
	}

	err = w.ProcessSnapshots(ctx, snapshots)
	if err != nil {
		return err
	}

	for _, snapshot := range snapshots {
		w.cache.SetDefault(fmt.Sprintf(snapshotKey, snapshot.SnapshotID), true)
	}

	if newOffset.After(offset) {
		if _, err := w.propertys.Set(ctx, snapshotCheckpoint, newOffset.Format(time.RFC3339Nano)); err != nil {
			log.WithError(err).Errorln("properties.Set")
			return err
		}
	}
	return nil
}
