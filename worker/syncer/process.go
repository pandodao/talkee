package syncer

import (
	"context"
	"encoding/base64"
	"fmt"

	"talkee/core"

	"github.com/fox-one/pkg/logger"
)

func (w *Worker) ProcessSnapshots(ctx context.Context, snapshots []*core.Snapshot) error {
	for _, snapshot := range snapshots {
		err := w.ProcessSnapshot(ctx, snapshot)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Worker) ProcessSnapshot(ctx context.Context, snapshot *core.Snapshot) error {
	log := logger.FromContext(ctx)
	if err := w.PersistentSnapshot(ctx, snapshot); err != nil {
		log.Warn("failed to persistent snapshot", err)
		return err
	}

	return nil
}

func (w *Worker) PersistentSnapshot(ctx context.Context, snapshot *core.Snapshot) (err error) {
	log := logger.FromContext(ctx).WithField("worker", "cashier:processer")
	// handle transaction
	fmt.Printf("snapshot: %v\n", snapshot)
	action := core.SnapshotAction{}

	memo := snapshot.Memo
	var memoBytes []byte
	memoBytes, err = base64.StdEncoding.DecodeString(memo)
	if err != nil {
		log.Info("base64 decode memo error ignored: ", err, snapshot.Memo)
		memoBytes = []byte(memo)
	}

	if err := action.Unmarshal(memoBytes); err != nil {
		log.Warn("Unmarshal memo error ignored: ", err, snapshot.Memo)
		return nil
	}

	switch action.Action {
	case core.ActionDeposit:
		log.Info("deposit transaction accepted")
	}

	if err := w.snapshots.SetSnapshot(ctx, snapshot); err != nil {
		return err
	}
	return nil
}
