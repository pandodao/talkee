package asset

import (
	"talkee/core"
	"talkee/store"
	"talkee/store/asset/dao"

	"gorm.io/gen"
)

func init() {
	store.RegistGenerate(
		gen.Config{
			OutPath: "store/asset/dao",
		},
		func(g *gen.Generator) {
			g.ApplyInterface(func(core.AssetStore) {}, core.Asset{})
		},
	)
}

func New(h *store.Handler) core.AssetStore {
	var q *dao.Query
	if !dao.Q.Available() {
		dao.SetDefault(h.DB)
		q = dao.Q
	} else {
		q = dao.Use(h.DB)
	}

	v, ok := interface{}(q.Asset).(core.AssetStore)
	if !ok {
		panic("dao.Asset is not core.AssetStore")
	}

	return &storeImpl{
		AssetStore: v,
	}
}

type storeImpl struct {
	core.AssetStore
}
