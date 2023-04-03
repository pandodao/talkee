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
	dao.SetDefault(h.DB)
	s := &storeImpl{}
	v, ok := interface{}(dao.Property).(core.PropertyStore)
	if !ok {
		panic("dao.Property is not core.PropertyStore")
	}
	s.PropertyStore = v
	return s
}

type storeImpl struct {
	core.PropertyStore
}
