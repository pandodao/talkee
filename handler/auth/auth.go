package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"talkee/core"
	"talkee/handler/render"
	"talkee/session"

	"github.com/fox-one/pkg/httputil/param"
	"github.com/golang-jwt/jwt"
	"github.com/olahol/melody"
	"github.com/pandodao/passport-go/auth"
	"github.com/patrickmn/go-cache"
)

type LoginPayload struct {
	Method        string `json:"method"`
	MixinToken    string `json:"mixin_token"`
	Signature     string `json:"signature"`
	SignedMessage string `json:"signed_message"`
	Lang          string `json:"lang"`
}

func Login(s *session.Session, userz core.UserService, clientID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		body := &LoginPayload{}
		if err := param.Binding(r, body); err != nil {
			render.Error(w, http.StatusBadRequest, err)
			return
		}

		origin := r.Header.Get("Origin")
		domain := strings.Split(origin, "//")[1]

		authorizer := auth.New([]string{
			clientID,
		}, []string{
			domain,
		})

		switch body.Method {
		case "mixin_token":
			{
				authUser, err := authorizer.Authorize(ctx, &auth.AuthorizationParams{
					Method:     auth.AuthMethodMixinToken,
					MixinToken: body.MixinToken,
				})
				if err != nil {
					render.Error(w, http.StatusUnauthorized, err)
					return
				}
				user, token, err := s.LoginWithMixin(ctx, userz, authUser, body.Lang)
				if err != nil {
					render.Error(w, http.StatusUnauthorized, err)
					return
				}
				render.JSON(w, map[string]interface{}{
					"user":         user,
					"access_token": token,
				})
				return
			}
		case "mvm":
			{
				authUser, err := authorizer.Authorize(ctx, &auth.AuthorizationParams{
					Method:           auth.AuthMethodMvm,
					MvmSignature:     body.Signature,
					MvmSignedMessage: body.SignedMessage,
				})
				if err != nil {
					render.Error(w, http.StatusUnauthorized, err)
					return
				}
				user, token, err := s.LoginWithMixin(ctx, userz, authUser, body.Lang)
				if err != nil {
					render.Error(w, http.StatusUnauthorized, err)
					return
				}
				render.JSON(w, map[string]interface{}{
					"user":         user,
					"access_token": token,
				})
				return
			}
		default:
			render.JSON(w, nil)
			return
		}
	}
}

func HandleAuthentication(s *session.Session, users core.UserStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			accessToken := getBearerToken(r)
			if accessToken == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims := &session.JwtClaims{}

			tkn, err := jwt.ParseWithClaims(accessToken, claims,
				func(token *jwt.Token) (interface{}, error) {
					return s.JwtSecret, nil
				},
			)

			if err != nil {
				fmt.Println("jwt.ParseWithClaims", err)
				next.ServeHTTP(w, r)
				return
			}
			if !tkn.Valid {
				next.ServeHTTP(w, r)
				return
			}

			user, err := users.GetUser(ctx, claims.UserID)
			if err != nil {
				fmt.Println("users.GetUser", err)
				next.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(
				session.WithUser(ctx, user),
			))
		}

		return http.HandlerFunc(fn)
	}
}

func HandleWssAuthentication(ctx context.Context, m *melody.Melody, s *session.Session, users core.UserStore) func(*melody.Session, []byte) {
	return func(msession *melody.Session, data []byte) {

		var payload struct {
			Type  string
			Token string
		}

		if err := json.Unmarshal(data, &payload); err != nil {
			msession.Close()
			return
		}

		if payload.Type != "auth" {
			msession.Close()
			return
		}

		accessToken := payload.Token
		if accessToken == "" {
			msession.Close()
			return
		}

		claims := &session.JwtClaims{}

		tkn, err := jwt.ParseWithClaims(accessToken, claims,
			func(token *jwt.Token) (interface{}, error) {
				return s.JwtSecret, nil
			},
		)

		if err != nil {
			fmt.Println("jwt.ParseWithClaims", err)
			msession.Close()
			return
		}

		if !tkn.Valid {
			msession.Close()
			return
		}

		user, err := users.GetUser(ctx, claims.UserID)
		if err != nil {
			msession.Close()
			return
		}

		msession.Set("user_id", user.ID)
	}

}

func ApplySiteCors(originCache *cache.Cache, s *session.Session, sites core.SiteStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			origin := r.Header.Get("Origin")
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			pass := false
			origin = strings.ToLower(origin)
			val, found := originCache.Get(origin)
			if found {
				if val.(bool) {
					pass = true
				} else {
					pass = false
				}
			} else {
				site, err := sites.GetSiteByOrigin(ctx, origin)
				if err != nil || site == nil {
					pass = false
				} else {
					pass = true
				}
			}

			if !pass {
				originCache.Set(origin, false, cache.DefaultExpiration)
				next.ServeHTTP(w, r)
				return
			}

			originCache.Set(origin, true, cache.DefaultExpiration)

			w.Header().Add("Connection", "keep-alive")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin,Content-Length,Content-Type,Authorization,X-TALKEE-SITE-ID,X-TALKEE-SITE-SLUG")
			w.Header().Add("Access-Control-Max-Age", "86400")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				_, _ = w.Write([]byte(""))
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func getBearerToken(r *http.Request) string {
	s := r.Header.Get("Authorization")
	return strings.TrimPrefix(s, "Bearer ")
}
