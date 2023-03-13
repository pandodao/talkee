package reply

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

type AddReplyPayload struct {
	Content string `json:"content"`
}

func GetReplies(replys core.ReplyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		commentID, _ := strconv.ParseUint(chi.URLParam(r, "commentID"), 10, 64)
		offset, _ := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
		limit, _ := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)

		if commentID <= 0 {
			render.Error(w, http.StatusBadRequest, core.ErrInvalidCommentID)
			return
		}
		if offset <= 0 {
			offset = 0
		}
		if limit <= 0 {
			limit = core.DefaultRepliesPerPage
		}

		count, err := replys.CountReplies(ctx, commentID)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}
		replies, err := replys.GetReplies(ctx, commentID, offset, limit)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}
		render.JSON(w, map[string]interface{}{
			"total":   count,
			"replies": replies,
		})
	}
}

func AddReply(replyz core.ReplyService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, found := session.UserFrom(ctx)
		if !found {
			render.Error(w, http.StatusUnauthorized, nil)
			return
		}

		commentID, _ := strconv.ParseUint(chi.URLParam(r, "commentID"), 10, 64)
		if commentID <= 0 {
			render.Error(w, http.StatusBadRequest, core.ErrInvalidCommentID)
			return
		}

		body := &AddReplyPayload{}
		if err := param.Binding(r, body); err != nil {
			render.Error(w, http.StatusBadRequest, err)
			return
		}

		body.Content = strings.TrimSpace(body.Content)
		if body.Content == "" {
			render.Error(w, http.StatusBadRequest, core.ErrInvalidCommentContent)
			return
		}

		reply, err := replyz.CreateReply(ctx, user.ID, commentID, body.Content)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}
		render.JSON(w, reply)
	}
}
