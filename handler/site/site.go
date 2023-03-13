package site

import (
	"net/http"
	"strconv"
	"strings"
	"talkee/core"
	"talkee/handler/render"
	"talkee/session"

	"github.com/fox-one/pkg/httputil/param"
	"github.com/go-chi/chi"
)

type SitePayload struct {
	Name       string `json:"name"`
	Origin     string `json:"origin"`
	UseArweave bool   `json:"use_arweave"`
}

func GetSite(sites core.SiteStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		siteIDStr := chi.URLParam(r, "siteID")
		siteID, _ := strconv.ParseUint(siteIDStr, 10, 64)

		if siteID <= 0 {
			render.Error(w, http.StatusBadRequest, core.ErrInvalidSiteID)
			return
		}

		site, err := sites.GetSite(ctx, siteID)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}

		render.JSON(w, site)
	}
}

func UpdateSite(sites core.SiteStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, found := session.UserFrom(ctx)
		if !found {
			render.Error(w, http.StatusUnauthorized, core.ErrUnauthorized)
		}

		siteIDStr := chi.URLParam(r, "siteID")
		siteID, _ := strconv.ParseUint(siteIDStr, 10, 64)

		site, err := sites.GetSite(ctx, siteID)
		if err != nil || site.UserID != user.ID {
			render.Error(w, http.StatusNotFound, core.ErrNoRecord)
			return
		}

		body := &SitePayload{}
		if err := param.Binding(r, body); err != nil {
			render.Error(w, http.StatusBadRequest, err)
			return
		}

		err = sites.UpdateSite(ctx, site.ID, map[string]interface{}{
			"name":        body.Name,
			"origins":     []string{body.Origin},
			"use_arweave": body.UseArweave,
		})
		if err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}

		site.Name = body.Name
		site.Origins = []string{body.Origin}
		site.UseArweave = body.UseArweave

		render.JSON(w, site)
	}
}

func CreateSite(sites core.SiteStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, found := session.UserFrom(ctx)
		if !found {
			render.Error(w, http.StatusUnauthorized, core.ErrUnauthorized)
			return
		}

		body := &SitePayload{}
		if err := param.Binding(r, body); err != nil {
			render.Error(w, http.StatusBadRequest, err)
			return
		}

		siteArr, _ := sites.GetSitesByUser(ctx, user.ID)
		if len(siteArr) >= 3 {
			render.Error(w, http.StatusBadRequest, core.ErrSiteLimitReached)
			return
		}

		if !strings.HasPrefix(body.Origin, "https://") && !strings.HasPrefix(body.Origin, "http://") {
			render.Error(w, http.StatusBadRequest, core.ErrInvalidOrigin)
			return
		}

		siteID, err := sites.CreateSite(ctx, user.ID, body.Name, []string{body.Origin}, body.UseArweave)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}

		site, err := sites.GetSite(ctx, siteID)
		if err != nil {
			render.Error(w, http.StatusNotFound, core.ErrNoRecord)
			return
		}

		render.JSON(w, site)
	}
}

func GetMySites(sites core.SiteStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, found := session.UserFrom(ctx)
		if !found {
			render.Error(w, http.StatusUnauthorized, core.ErrUnauthorized)
		}

		siteArr, err := sites.GetSitesByUser(ctx, user.ID)
		if err != nil {
			render.JSON(w, []core.Site{})
			return
		}

		render.JSON(w, siteArr)
	}
}
