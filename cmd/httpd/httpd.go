package httpd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"talkee/config"
	"talkee/handler"
	"talkee/handler/hc"
	assetServ "talkee/service/asset"
	commentServ "talkee/service/comment"
	replyServ "talkee/service/reply"
	userServ "talkee/service/user"
	"talkee/session"
	"talkee/store/asset"
	"talkee/store/comment"
	"talkee/store/favourite"
	"talkee/store/property"
	"talkee/store/reply"
	"talkee/store/site"
	"talkee/store/user"

	"github.com/fox-one/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"

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

			conn, err := sqlx.Connect(config.C().DB.Driver, config.C().DB.Datasource)
			if err != nil {
				log.Fatalln("connect to database failed", err)
			}
			conn.SetMaxIdleConns(2)

			defer conn.Close()

			cfg := config.C()

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

			propertys := property.New(conn)
			users := user.New(conn)
			comments := comment.New(conn)
			sites := site.New(conn)
			replys := reply.New(conn)
			assets := asset.New(conn)
			favourites := favourite.New(conn)

			userz := userServ.New(client, users, userServ.Config{
				MixinClientSecret: cfg.Auth.MixinClientSecret,
			})
			commentz := commentServ.New(arWallet, comments, commentServ.Config{
				AppName: cfg.AppName,
			})
			replyz := replyServ.New(replys, comments, replyServ.Config{})
			assetz := assetServ.New(client, assets)

			s = s.Setup(userz, []string{client.ClientID})

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
				mux.Mount("/hc", hc.Handle(cmd.Version))
			}

			// rpc & api
			{
				cfg := handler.Config{}
				svr := handler.New(cfg, s,
					propertys,
					comments, users, sites, replys, assets, favourites, userz, commentz, replyz, assetz, nil)
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
