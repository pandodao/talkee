package worker

import (
	"fmt"
	"net/http"
	"strconv"

	"talkee/config"
	"talkee/handler/hc"
	"talkee/internal/mixpay"
	assetServ "talkee/service/asset"
	commentServ "talkee/service/comment"
	rewardServ "talkee/service/reward"
	snapshotServ "talkee/service/snapshot"
	tipServ "talkee/service/tip"
	"talkee/session"
	"talkee/store"

	"talkee/store/asset"
	"talkee/store/comment"
	"talkee/store/favourite"
	"talkee/store/property"
	"talkee/store/reward"
	"talkee/store/snapshot"
	"talkee/store/tip"
	"talkee/store/user"
	"talkee/worker"
	"talkee/worker/arweavesyncer"
	"talkee/worker/reward_processer"
	"talkee/worker/syncer"
	"talkee/worker/timer"
	"talkee/worker/tip_deliver"

	"github.com/fox-one/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"

	"github.com/spf13/cobra"
)

func NewCmdWorker() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worker [health check port]",
		Short: "run workers",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			ctx := cmd.Context()
			cfg := config.C()

			h := store.MustInit(store.Config{
				Driver: cfg.DB.Driver,
				DSN:    cfg.DB.Datasource,
			})

			s := session.From(ctx)

			keystore, err := s.GetKeystore()
			if err != nil {
				return err
			}

			client, err := s.GetClient()
			if err != nil {
				return err
			}

			arWallet, err := s.GetArwallet()
			if err != nil {
				return err
			}

			mixpayClient := mixpay.New()

			propertys := property.New(h)

			assets := asset.New(h)
			snapshots := snapshot.New(h)
			users := user.New(h)
			comments := comment.New(h)
			rewards := reward.New(h)
			favourites := favourite.New(h)
			tips := tip.New(h)

			commentz := commentServ.New(arWallet, comments, rewards, users, commentServ.Config{
				AppName: cfg.AppName,
			})

			rewardServCfg := rewardServ.Config{}
			rewardServCfg.Pin, _ = s.GetPin()
			rewardz := rewardServ.New(rewards, favourites, comments, client, rewardServCfg)

			assetz := assetServ.New(client, assets)
			snapshotz := snapshotServ.New(client)
			tipz := tipServ.New(tipServ.Config{
				ClientID:          client.ClientID,
				MixpayPayeeID:     cfg.Mixpay.PayeeID,
				MixpayCallbackURL: cfg.Mixpay.CallbackURL,
			}, mixpayClient, tips, comments, rewards, commentz)

			workers := []worker.Worker{
				// reward task processer
				// reward_task_processer.New(rewardz),

				// tip delivery
				tip_deliver.New(tips, tipz),

				// reward process
				reward_processer.New(rewardz),

				// arweaver syncer
				arweavesyncer.New(commentz),

				// timer
				timer.New(timer.Config{}, propertys, comments, favourites, rewards, assets, assetz),

				// syncer
				syncer.New(syncer.Config{
					ClientID: keystore.ClientID,
				}, propertys, snapshots, users, snapshotz),
			}

			// run them all
			g, ctx := errgroup.WithContext(ctx)
			for idx := range workers {
				w := workers[idx]
				g.Go(func() error {
					return w.Run(ctx)
				})
			}

			// start the health check server
			g.Go(func() error {
				mux := chi.NewMux()
				mux.Use(middleware.Recoverer)
				mux.Use(middleware.StripSlashes)
				mux.Use(cors.AllowAll().Handler)
				mux.Use(logger.WithRequestID)
				mux.Use(middleware.Logger)

				{
					// hc for workers
					mux.Mount("/hc", hc.Handle(cmd.Version))
				}

				// launch server
				port := 8081
				if len(args) > 0 {
					port, err = strconv.Atoi(args[0])
					if err != nil {
						port = 8081
					}
				}

				addr := fmt.Sprintf(":%d", port)
				return http.ListenAndServe(addr, mux)
			})

			if err := g.Wait(); err != nil {
				cmd.PrintErrln("run worker", err)
			}

			return nil
		},
	}

	return cmd
}
