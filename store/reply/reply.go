package reply

import (
	"talkee/core"
	"talkee/store"

	"talkee/store/reply/dao"

	"gorm.io/gen"
)

func init() {
	store.RegistGenerate(
		gen.Config{
			OutPath: "store/reply/dao",
		},
		func(g *gen.Generator) {
			g.ApplyInterface(func(core.ReplyStore) {}, core.Reply{})
		},
	)
}

func New(h *store.Handler) core.ReplyStore {
	var q *dao.Query
	if !dao.Q.Available() {
		dao.SetDefault(h.DB)
		q = dao.Q
	} else {
		q = dao.Use(h.DB)
	}

	v, ok := interface{}(q.Reply).(core.ReplyStore)
	if !ok {
		panic("dao.Reply is not core.ReplyStore")
	}

	return &storeImpl{
		ReplyStore: v,
	}
}

type storeImpl struct {
	core.ReplyStore
}
