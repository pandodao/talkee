package reward

import (
	"talkee/core"
	"talkee/store"
	"talkee/store/reward/dao"

	"gorm.io/gen"
)

func init() {
	store.RegistGenerate(
		gen.Config{
			OutPath: "store/reward/dao",
		},
		func(g *gen.Generator) {
			g.ApplyInterface(func(core.RewardStore) {}, core.Reward{})
		},
	)
}

func New(h *store.Handler) core.RewardStore {
	var q *dao.Query
	if !dao.Q.Available() {
		dao.SetDefault(h.DB)
		q = dao.Q
	} else {
		q = dao.Use(h.DB)
	}

	v, ok := interface{}(q.Reward).(core.RewardStore)
	if !ok {
		panic("dao.Reward is not core.RewardStore")
	}

	return &storeImpl{
		RewardStore: v,
	}
}

type storeImpl struct {
	core.RewardStore
}
