package site

import (
	"talkee/core"
	"talkee/store"
	"talkee/store/site/dao"

	"gorm.io/gen"
)

func init() {
	store.RegistGenerate(
		gen.Config{
			OutPath: "store/site/dao",
		},
		func(g *gen.Generator) {
			g.ApplyInterface(func(core.SiteStore) {}, core.Site{})
		},
	)
}

func New(h *store.Handler) core.SiteStore {
	var q *dao.Query
	if !dao.Q.Available() {
		dao.SetDefault(h.DB)
		q = dao.Q
	} else {
		q = dao.Use(h.DB)
	}

	v, ok := interface{}(q.Site).(core.SiteStore)
	if !ok {
		panic("dao.Site is not core.SiteStore")
	}

	return &storeImpl{
		SiteStore: v,
	}
}

type storeImpl struct {
	core.SiteStore
}
