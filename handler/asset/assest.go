package asset

import (
	"net/http"
	"talkee/core"
	"talkee/handler/render"

	"github.com/go-chi/chi"
)

func GetAssets(assets core.AssetStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assets, err := assets.GetAssets(r.Context())
		if err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}

		render.JSON(w, assets)
	}
}

func GetAsset(assetz core.AssetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		assetID := chi.URLParam(r, "asset_id")
		if assetID == "" {
			render.Error(w, http.StatusBadRequest, core.ErrInvalidAssetID)
			return
		}

		asset, err := assetz.GetCachedAssets(r.Context(), assetID)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}
		render.JSON(w, asset)
	}
}
