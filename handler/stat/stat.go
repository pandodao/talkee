package stat

import (
	"net/http"
	"talkee/core"
	"talkee/handler/render"

	jsoniter "github.com/json-iterator/go"
)

const (
	statGlobalSummary = "stat:global_summary"
)

func GetGlobal(propertys core.PropertyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		val, err := propertys.Get(ctx, statGlobalSummary)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}

		stat := make(map[string]interface{})
		if err := jsoniter.Unmarshal(val.Bytes(), &stat); err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}

		render.JSON(w, stat)
	}
}
