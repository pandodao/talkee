package slug

import (
	"net/http"
	"talkee/core"
	"talkee/handler/render"
	"talkee/session"
)

func GetSlug(tipz core.TipService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		siteID, slug, found := session.SiteInfoFrom(ctx)
		if !found {
			render.Error(w, http.StatusBadRequest, core.ErrInvalidSiteID)
			return
		}

		tips, err := tipz.GetTipsBySlug(ctx, siteID, slug, 64)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}

		render.JSON(w, map[string]interface{}{
			"site_id": siteID,
			"slug":    slug,
			"tips":    tips,
		})
	}
}
