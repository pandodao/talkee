package tip

import (
	"talkee/core"
	"talkee/store"
	"talkee/store/tip/dao"

	"gorm.io/gen"
)

func init() {
	store.RegistGenerate(
		gen.Config{
			OutPath: "store/tip/dao",
		},
		func(g *gen.Generator) {
			g.ApplyInterface(func(core.TipStore) {}, core.Tip{})
		},
	)
}

func New(h *store.Handler) core.TipStore {
	var q *dao.Query
	if !dao.Q.Available() {
		dao.SetDefault(h.DB)
		q = dao.Q
	} else {
		q = dao.Use(h.DB)
	}

	v, ok := interface{}(q.Tip).(core.TipStore)
	if !ok {
		panic("dao.Tip is not core.TipStore")
	}

	return &storeImpl{
		TipStore: v,
	}
}

type storeImpl struct {
	core.TipStore
}
