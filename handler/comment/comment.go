package comment

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

type AddCommentPayload struct {
	Content string `json:"content"`
}

func GetComments(favors core.FavouriteStore, comments core.CommentStore, commentz core.CommentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		siteID, slug, found := session.SiteInfoFrom(ctx)
		if !found {
			render.Error(w, http.StatusBadRequest, core.ErrInvalidSiteID)
			return
		}

		offset, _ := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
		limit, _ := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
		orderBy := strings.TrimSpace(r.URL.Query().Get("order_by"))

		if offset <= 0 {
			offset = 0
		}
		if limit <= 0 {
			limit = core.DefaultCommentsPerPage
		}

		order := "DESC"
		if orderBy != "created_at_asc" && orderBy != "created_at" && orderBy != "favor_count" && orderBy != "reply_count" {
			orderBy = "favor_count"
		}

		if orderBy == "created_at_asc" {
			orderBy = "created_at"
			order = "ASC"
		}

		count, err := comments.CountComments(ctx, siteID, slug)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}

		items, err := commentz.GetCommentsWithRewards(ctx, siteID, slug, offset, limit, orderBy, order)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}

		user, found := session.UserFrom(ctx)
		if found {
			if len(items) > 0 {
				cmtIDs := []uint64{}
				for _, item := range items {
					cmtIDs = append(cmtIDs, item.ID)
				}

				favs, err := favors.FindUserFavourites(r.Context(), user.ID, cmtIDs)
				if err != nil {
					render.Error(w, http.StatusInternalServerError, err)
					return
				}

				for _, item := range items {
					for _, f := range favs {
						if f.CommentID == item.ID {
							item.FavID = f.ID
						}
					}
				}
			}
		}

		render.JSON(w, map[string]interface{}{
			"total":    count,
			"comments": items,
		})
	}
}

func GetComment(favors core.FavouriteStore, commentz core.CommentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		commentIDStr := chi.URLParam(r, "commentID")
		commentID, _ := strconv.ParseUint(commentIDStr, 10, 64)
		if commentID <= 0 {
			render.Error(w, http.StatusBadRequest, core.ErrInvalidCommentID)
			return
		}

		comment, err := commentz.GetCommentWithReward(ctx, commentID)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}

		user, found := session.UserFrom(ctx)
		if found {
			if comment != nil {
				cmtIDs := []uint64{comment.ID}

				favs, err := favors.FindUserFavourites(r.Context(), user.ID, cmtIDs)
				if err != nil {
					render.Error(w, http.StatusInternalServerError, err)
					return
				}

				if len(favs) > 0 {
					comment.FavID = favs[0].ID
				}

			}
		}

		render.JSON(w, comment)
	}
}

func AddComment(commentz core.CommentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, found := session.UserFrom(ctx)
		if !found {
			render.Error(w, http.StatusUnauthorized, nil)
			return
		}

		siteID, slug, found := session.SiteInfoFrom(ctx)
		if !found {
			render.Error(w, http.StatusBadRequest, core.ErrInvalidSiteID)
			return
		}

		body := &AddCommentPayload{}
		if err := param.Binding(r, body); err != nil {
			render.Error(w, http.StatusBadRequest, err)
			return
		}

		body.Content = strings.TrimSpace(body.Content)
		if body.Content == "" {
			render.Error(w, http.StatusBadRequest, core.ErrInvalidCommentContent)
			return
		}

		comment, err := commentz.CreateComment(ctx, user.ID, siteID, slug, body.Content)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}
		render.JSON(w, comment)
	}
}

func FavComment(commentz core.CommentService, fav bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, found := session.UserFrom(ctx)
		if !found {
			render.Error(w, http.StatusUnauthorized, nil)
			return
		}

		commentIDStr := chi.URLParam(r, "commentID")
		commentID, _ := strconv.ParseUint(commentIDStr, 10, 64)
		if commentID <= 0 {
			render.Error(w, http.StatusBadRequest, core.ErrInvalidCommentID)
			return
		}

		comment, err := commentz.FavComment(ctx, commentID, user.ID, fav)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}

		render.JSON(w, comment)
	}
}
