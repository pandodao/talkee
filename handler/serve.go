package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"talkee/core"
	"talkee/handler/asset"
	"talkee/handler/auth"
	"talkee/handler/comment"
	"talkee/handler/echo"
	"talkee/handler/message"
	"talkee/handler/payment"
	"talkee/handler/render"
	"talkee/handler/reply"
	"talkee/handler/site"
	"talkee/handler/stat"
	"talkee/handler/tip"
	"talkee/handler/user"
	"talkee/handler/slug"
	"talkee/session"

	"github.com/go-chi/chi"
	"github.com/olahol/melody"
	"github.com/patrickmn/go-cache"
)

func New(cfg Config,
	session *session.Session,
	propertys core.PropertyStore,
	comments core.CommentStore,
	users core.UserStore,
	sites core.SiteStore,
	replys core.ReplyStore,
	assets core.AssetStore,
	favourites core.FavouriteStore,

	userz core.UserService,
	commentz core.CommentService,
	replyz core.ReplyService,
	assetz core.AssetService,
	tipz core.TipService,

	melody *melody.Melody,
) Server {
	// @TODO: make this longer, but need to update the cache on site update
	originCache := cache.New(5*time.Minute, 5*time.Minute)
	return Server{cfg: cfg, session: session,
		propertys:  propertys,
		comments:   comments,
		users:      users,
		sites:      sites,
		replys:     replys,
		assets:     assets,
		favourites: favourites,

		userz:    userz,
		commentz: commentz,
		replyz:   replyz,
		assetz:   assetz,
		tipz:     tipz,

		originCache: originCache,
		melody:      melody,
	}
}

type (
	Config struct {
		ClientID string
	}

	Server struct {
		cfg        Config
		session    *session.Session
		propertys  core.PropertyStore
		comments   core.CommentStore
		users      core.UserStore
		sites      core.SiteStore
		replys     core.ReplyStore
		assets     core.AssetStore
		favourites core.FavouriteStore

		userz    core.UserService
		commentz core.CommentService
		replyz   core.ReplyService
		assetz   core.AssetService
		tipz     core.TipService

		originCache *cache.Cache
		melody      *melody.Melody
	}
)

func (s Server) HandleRest() http.Handler {
	r := chi.NewRouter()
	r.Use(render.WrapResponse(true))
	r.Use(auth.HandleAuthentication(s.session, s.users))
	r.Use(auth.ApplySiteCors(s.originCache, s.session, s.sites))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Error(w, http.StatusNotFound, fmt.Errorf("not found"))
	})

	r.Route("/echo", func(r chi.Router) {
		r.Get("/{msg}", echo.Get())
		r.Post("/", echo.Post())
	})

	r.Route("/stat", func(r chi.Router) {
		r.Get("/global", stat.GetGlobal(s.propertys))
	})

	r.Route("/assets", func(r chi.Router) {
		r.Get("/{asset_id}", asset.GetAsset(s.assetz))
		r.Get("/", asset.GetAssets(s.assets))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", auth.Login(s.session, s.userz, s.cfg.ClientID))
	})

	r.Route("/sites", func(r chi.Router) {
		r.Get("/{siteID}", site.GetSite(s.sites))
		r.With(s.LoginRequired()).Post("/", site.CreateSite(s.sites))
		r.With(s.LoginRequired()).Get("/", site.GetMySites(s.sites))
		r.With(s.LoginRequired()).Put("/{siteID}", site.UpdateSite(s.sites))
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/{userID}", user.GetUser(s.users))
	})

	r.With(s.LoginRequired()).Route("/me", func(r chi.Router) {
		r.Get("/", user.GetMe(s.users))
	})

	r.With(s.SiteRequired()).Route("/airdrops", func(r chi.Router) {
		r.With(s.LoginRequired()).Post("/", tip.CreateTip(s.tipz))
	})

	r.With(s.SiteRequired()).Route("/slugs", func(r chi.Router) {
		r.Get("/", slug.GetSlug(s.tipz))
	})

	r.With(s.SiteRequired()).Route("/comments", func(r chi.Router) {
		r.Get("/", comment.GetComments(s.favourites, s.comments, s.commentz))
		r.With(s.LoginRequired()).Post("/", comment.AddComment(s.commentz))
		r.Get("/{commentID}", comment.GetComment(s.favourites, s.commentz))
		r.With(s.LoginRequired()).Put("/{commentID}/fav", comment.FavComment(s.commentz, true))
		r.With(s.LoginRequired()).Put("/{commentID}/unfav", comment.FavComment(s.commentz, false))
		r.Get("/{commentID}/replies", reply.GetReplies(s.replys, s.replyz))
		r.With(s.LoginRequired()).Post("/{commentID}/replies", reply.AddReply(s.replyz))
	})

	r.Route("/payments", func(r chi.Router) {
		r.Post("/mixpay", payment.HandleMixpayWebhook(s.tipz))
	})

	return r
}

func (s Server) WebsocketConnectHandler() http.Handler {
	r := chi.NewRouter()
	r.Use(auth.ApplySiteCors(s.originCache, s.session, s.sites))
	r.With(s.SiteRequired()).With(auth.HandleAuthentication(s.session, s.users)).Post("/messages", message.CreateMessage(s.melody))
	r.With(s.SiteRequired()).With(auth.HandleAuthentication(s.session, s.users)).Get("/messages", message.LatestMessages())
	r.With(s.SiteRequired()).Handle("/", s.HandleWebsocket())

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Error(w, http.StatusNotFound, fmt.Errorf("not found"))
	})

	return r
}

func (s Server) HandleWebsocket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		siteID, slug, found := session.SiteInfoFrom(ctx)
		if !found {
			render.Error(w, http.StatusBadRequest, core.ErrInvalidSiteID)
			return
		}

		err := s.melody.HandleRequestWithKeys(w, r, map[string]interface{}{
			"slug":    slug,
			"site_id": siteID,
		})

		if err != nil {
			render.Error(w, http.StatusBadRequest, err)
			return
		}
	}
}

func (s Server) LoginRequired() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if _, found := session.UserFrom(ctx); !found {
				render.Error(w, http.StatusUnauthorized, core.ErrUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func (s Server) SiteRequired() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// read site_id, slug from headers
			// and put them into ctx
			siteIDStr := strings.TrimSpace(r.Header.Get("X-TALKEE-SITE-ID"))
			slug := strings.TrimSpace(r.Header.Get("X-TALKEE-SITE-SLUG"))
			if siteIDStr == "" {
				siteIDStr = r.URL.Query().Get("site_id")
			}
			if slug == "" {
				slug = r.URL.Query().Get("slug")
			}

			siteID, err := strconv.ParseUint(siteIDStr, 10, 64)
			if err != nil || siteID <= 0 {
				render.Error(w, http.StatusBadRequest, core.ErrInvalidSiteID)
				return
			}
			if slug == "" {
				render.Error(w, http.StatusBadRequest, core.ErrInvalidSlug)
				return
			}

			ctx = session.WithSiteInfo(ctx, siteID, slug)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
