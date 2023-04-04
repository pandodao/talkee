package reward

import (
	"talkee/core"
	"talkee/store"

	"gorm.io/gen"
)

func init() {
	store.RegistGenerate(
		gen.Config{
			OutPath: "store/reward_task/dao",
		},
		func(g *gen.Generator) {
			g.ApplyInterface(func(core.RewardTaskStore) {}, core.RewardTask{})
		},
	)
}

func New(h *store.Handler) core.RewardTaskStore {
	// var q *dao.Query
	// if !dao.Q.Available() {
	// 	dao.SetDefault(h.DB)
	// 	q = dao.Q
	// } else {
	// 	q = dao.Use(h.DB)
	// }

	// v, ok := interface{}(q.Reward).(core.RewardTaskStore)
	// if !ok {
	// 	panic("dao.Reward is not core.RewardTaskStore")
	// }

	// return &storeImpl{
	// 	RewardTaskStore: v,
	// }
	return &storeImpl{}
}

type storeImpl struct {
	core.RewardTaskStore
}
