package httpd

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"talkee/config"
	"talkee/handler"
	"talkee/handler/hc"
	"talkee/internal/mixpay"
	assetServ "talkee/service/asset"
	commentServ "talkee/service/comment"
	replyServ "talkee/service/reply"
	tipServ "talkee/service/tip"
	userServ "talkee/service/user"
	"talkee/session"
	"talkee/store"
	"talkee/store/asset"
	"talkee/store/comment"
	"talkee/store/favourite"
	"talkee/store/property"
	"talkee/store/reply"
	"talkee/store/reward"
	"talkee/store/site"
	"talkee/store/tip"
	"talkee/store/user"

	"github.com/fox-one/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/drone/signal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmdHttpd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "httpd [port]",
		Short: "start the httpd daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			ctx := cmd.Context()
			cfg := config.C()

			h := store.MustInit(store.Config{
				Driver: cfg.DB.Driver,
				DSN:    cfg.DB.Datasource,
			})

			s := session.From(ctx)
			s.WithJWTSecret([]byte(config.C().Auth.JwtSecret))

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
			users := user.New(h)
			comments := comment.New(h)
			sites := site.New(h)
			replys := reply.New(h)
			assets := asset.New(h)
			favourites := favourite.New(h)
			rewards := reward.New(h)
			tips := tip.New(h)

			userz := userServ.New(client, users, userServ.Config{
				MixinClientSecret: cfg.Auth.MixinClientSecret,
			})
			commentz := commentServ.New(arWallet, comments, rewards, users, commentServ.Config{
				AppName: cfg.AppName,
			})
			replyz := replyServ.New(replys, comments, users, replyServ.Config{})
			assetz := assetServ.New(client, assets)
			tipz := tipServ.New(tipServ.Config{
				ClientID:          client.ClientID,
				MixpayPayeeID:     cfg.Mixpay.PayeeID,
				MixpayCallbackURL: cfg.Mixpay.CallbackURL,
			}, mixpayClient, tips, comments, rewards, commentz)

			mux := chi.NewMux()
			mux.Use(middleware.Recoverer)
			mux.Use(middleware.StripSlashes)
			mux.Use(logger.WithRequestID)
			mux.Use(middleware.Logger)
			mux.Use(middleware.NewCompressor(5).Handler)

			{
				mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("hello world"))
				})
			}

			// hc
			{
				mux.Mount("/hc", hc.Handle(s.Version))
			}

			// rpc & api
			{
				cfg := handler.Config{
					ClientID: client.ClientID,
				}
				svr := handler.New(cfg, s,
					propertys,
					comments, users, sites, replys, assets, favourites, userz, commentz, replyz, assetz, tipz, nil)
				// api
				restHandler := svr.HandleRest()
				mux.Mount("/api", restHandler)
			}

			port := 8080
			if len(args) > 0 {
				port, err = strconv.Atoi(args[0])
				if err != nil {
					port = 8080
				}
			}

			// launch server
			if err != nil {
				panic(err)
			}
			addr := fmt.Sprintf(":%d", port)

			svr := &http.Server{
				Addr:    addr,
				Handler: mux,
			}

			done := make(chan struct{}, 1)
			ctx = signal.WithContextFunc(ctx, func() {
				logrus.Debug("shutdown server...")

				// create context with timeout
				ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
				defer cancel()

				if err := svr.Shutdown(ctx); err != nil {
					logrus.WithError(err).Error("graceful shutdown server failed")
				}

				close(done)
			})

			logrus.Infoln("serve at", addr)
			if err := svr.ListenAndServe(); err != http.ErrServerClosed {
				logrus.WithError(err).Fatal("server aborted")
			}

			<-done
			return nil
		},
	}

	return cmd
}
