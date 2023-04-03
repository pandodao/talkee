package snapshot

import (
	"talkee/core"
	"talkee/store"
	"talkee/store/snapshot/dao"

	"gorm.io/gen"
)

func init() {
	store.RegistGenerate(
		gen.Config{
			OutPath: "store/snapshot/dao",
		},
		func(g *gen.Generator) {
			g.ApplyInterface(func(core.SnapshotStore) {}, core.Snapshot{})
		},
	)
}

func New(h *store.Handler) core.SnapshotStore {
	var q *dao.Query
	if !dao.Q.Available() {
		dao.SetDefault(h.DB)
		q = dao.Q
	} else {
		q = dao.Use(h.DB)
	}

	v, ok := interface{}(q.Snapshot).(core.SnapshotStore)
	if !ok {
		panic("dao.Snapshot is not core.SnapshotStore")
	}

	return &storeImpl{
		SnapshotStore: v,
	}
}

type storeImpl struct {
	core.SnapshotStore
}
