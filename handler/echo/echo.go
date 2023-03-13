package echo

import (
	"net/http"

	"talkee/handler/render"

	"github.com/fox-one/pkg/httputil/param"
	"github.com/go-chi/chi"
)

func Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Message string `json:"msg,omitempty" valid:"minstringlength(1),required"`
		}

		if err := param.Binding(r, &body); err != nil {
			render.Error(w, http.StatusBadRequest, err)
			return
		}
		render.JSON(w, body)
	}
}

func Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := chi.URLParam(r, "msg")
		render.JSON(w, render.H{
			"msg": msg,
		})
	}
}
