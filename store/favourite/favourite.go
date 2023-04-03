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
			g.ApplyInterface(func(core.CommentFavouriteStore) {}, core.CommentFavourite{})
		},
	)
}

func New(h *store.Handler) core.CommentFavouriteStore {
	var q *dao.Query
	if !dao.Q.Available() {
		dao.SetDefault(h.DB)
		q = dao.Q
	} else {
		q = dao.Use(h.DB)
	}

	v, ok := interface{}(q.CommentFavourite).(core.CommentFavouriteStore)
	if !ok {
		panic("dao.CommentFavourite is not core.CommentFavouriteStore")
	}

	return &storeImpl{
		CommentFavouriteStore: v,
	}
}

type storeImpl struct {
	core.CommentFavouriteStore
}
