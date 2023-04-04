package favourite

import (
	"talkee/core"
	"talkee/store"
	"talkee/store/favourite/dao"

	"gorm.io/gen"
)

func init() {
	store.RegistGenerate(
		gen.Config{
			OutPath: "store/favourite/dao",
		},
		func(g *gen.Generator) {
			g.ApplyInterface(func(core.FavouriteStore) {}, core.Favourite{})
		},
	)
}

func New(h *store.Handler) core.FavouriteStore {
	var q *dao.Query
	if !dao.Q.Available() {
		dao.SetDefault(h.DB)
		q = dao.Q
	} else {
		q = dao.Use(h.DB)
	}

	v, ok := interface{}(q.Favourite).(core.FavouriteStore)
	if !ok {
		panic("dao.Favourite is not core.FavouriteStore")
	}

	return &storeImpl{
		FavouriteStore: v,
	}
}

type storeImpl struct {
	core.FavouriteStore
}
