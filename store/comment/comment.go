package comment

import (
	"talkee/core"
	"talkee/store"
	"talkee/store/comment/dao"

	"gorm.io/gen"
)

func init() {
	store.RegistGenerate(
		gen.Config{
			OutPath: "store/comment/dao",
		},
		func(g *gen.Generator) {
			g.ApplyInterface(func(core.CommentStore) {}, core.Comment{}, core.ArweaveSyncListItem{})
		},
	)
}

func New(h *store.Handler) core.CommentStore {
	var q *dao.Query
	if !dao.Q.Available() {
		dao.SetDefault(h.DB)
		q = dao.Q
	} else {
		q = dao.Use(h.DB)
	}

	v, ok := interface{}(q.Comment).(core.CommentStore)
	if !ok {
		panic("dao.Comment is not core.CommentStore")
	}

	return &storeImpl{
		CommentStore: v,
	}
}

type storeImpl struct {
	core.CommentStore
}
