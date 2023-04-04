package property

import (
	"talkee/core"
	"talkee/store"
	"talkee/store/property/dao"

	"gorm.io/gen"
)

func init() {
	cfg := gen.Config{
		OutPath: "store/property/dao",
	}
	store.RegistGenerate(
		cfg,
		func(g *gen.Generator) {
			g.ApplyInterface(func(core.PropertyStore) {}, core.Property{})
		},
	)
}

func New(h *store.Handler) core.PropertyStore {
	var q *dao.Query
	if !dao.Q.Available() {
		dao.SetDefault(h.DB)
		q = dao.Q
	} else {
		q = dao.Use(h.DB)
	}

	v, ok := interface{}(q.Property).(core.PropertyStore)
	if !ok {
		panic("dao.Property is not core.PropertyStore")
	}

	return &storeImpl{
		PropertyStore: v,
	}
}

type storeImpl struct {
	core.PropertyStore
}
