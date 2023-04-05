package wss

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"talkee/config"
	"talkee/handler"
	"talkee/handler/auth"
	"talkee/handler/hc"
	commentServ "talkee/service/comment"
	replyServ "talkee/service/reply"
	userServ "talkee/service/user"
	"talkee/session"
	"talkee/store"
	"talkee/store/comment"
	"talkee/store/property"
	"talkee/store/reply"
	"talkee/store/reward"
	"talkee/store/site"
	"talkee/store/user"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/olahol/melody"

	"github.com/drone/signal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmdWss() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wss [port]",
		Short: "start the websocket daemon",
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

			propertys := property.New(h)
			users := user.New(h)
			comments := comment.New(h)
			sites := site.New(h)
			replys := reply.New(h)
			rewards := reward.New(h)

			userz := userServ.New(client, users, userServ.Config{
				MixinClientSecret: cfg.Auth.MixinClientSecret,
			})
			commentz := commentServ.New(nil, comments, rewards, users, commentServ.Config{
				AppName: cfg.AppName,
			})
			replyz := replyServ.New(replys, comments, users, replyServ.Config{})

			m := melody.New()
			m.HandleMessage(auth.HandleWssAuthentication(ctx, m, s, users))

			mux := chi.NewMux()
			mux.Use(middleware.Recoverer)
			mux.Use(middleware.StripSlashes)
			mux.Use(middleware.Logger)
			mux.Use(middleware.NewCompressor(5).Handler)
			// hc
			{
				mux.Mount("/hc", hc.Handle(cmd.Version))
			}

			// wss connect
			{
				cfg := handler.Config{}
				svr := handler.New(cfg, s,
					propertys,
					comments, users, sites, replys, nil, nil, userz, commentz, replyz, nil, nil, m)
				// wss
				mux.Mount("/ws", svr.WebsocketConnectHandler())
			}

			port := 8081
			if len(args) > 0 {
				port, err = strconv.Atoi(args[0])
				if err != nil {
					port = 8081
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
