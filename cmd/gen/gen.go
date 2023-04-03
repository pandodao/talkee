package gen

import (
	"talkee/config"
	"talkee/store"
	_ "talkee/store/asset"
	_ "talkee/store/favourite"
	_ "talkee/store/property"
	_ "talkee/store/site"
	_ "talkee/store/snapshot"
	_ "talkee/store/user"

	"github.com/spf13/cobra"
)

func NewCmdGen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "generate database operation code",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.C()
			h := store.MustInit(store.Config{
				Driver: cfg.DB.Driver,
				DSN:    cfg.DB.Datasource,
			})
			h.Generate()
		},
	}

	return cmd
}
