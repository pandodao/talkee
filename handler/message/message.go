package message

import (
	"container/list"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"talkee/core"
	"talkee/handler/render"
	"talkee/session"

	"github.com/fox-one/pkg/httputil/param"
	"github.com/olahol/melody"
)

type MessagePayload struct {
	MessageID string    `json:"id,omitempty"`
	Category  string    `json:"category,omitempty"`
	Data      string    `json:"data,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type messageList struct {
	*list.List
	mu sync.Mutex
}

var lastestMessageList = &messageList{List: list.New()}

type MessageResponse struct {
	MessagePayload
	User *core.User `json:"user,omitempty"`
}

func CreateMessage(m *melody.Melody) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, found := session.UserFrom(ctx)
		if !found {
			render.Error(w, http.StatusUnauthorized, core.ErrUnauthorized)
			return
		}

		siteID, slug, found := session.SiteInfoFrom(ctx)
		if !found {
			render.Error(w, http.StatusBadRequest, core.ErrInvalidSiteID)
			return
		}

		var payload MessagePayload
		if err := param.Binding(r, &payload); err != nil {
			render.Error(w, http.StatusBadRequest, err)
			return
		}

		var resp = &MessageResponse{}
		resp.MessagePayload = payload
		resp.CreatedAt = time.Now()
		resp.User = user

		data, err := json.Marshal(resp)
		if err != nil {
			render.Error(w, http.StatusBadRequest, err)
			return
		}

		if err := m.BroadcastFilter(data, func(s *melody.Session) bool {
			if val, ok := s.Keys["site_id"]; !ok || val != siteID {
				return false
			}

			if val, ok := s.Keys["slug"]; !ok || val != slug {
				return false
			}

			if val, ok := s.Keys["user_id"]; ok && val == user.ID {
				return false
			}

			return true
		}); err != nil {
			render.Error(w, http.StatusBadRequest, err)
			return
		}

		lastestMessageList.mu.Lock()
		defer lastestMessageList.mu.Unlock()

		lastestMessageList.PushBack(resp)
		if lastestMessageList.Len() > 100 {
			lastestMessageList.Remove(lastestMessageList.Front())
		}

		render.JSON(w, resp)
	}
}

func LatestMessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp = []*MessageResponse{}
		msg := lastestMessageList.Front()

		for {
			if msg == nil || msg.Value == nil {
				break
			}
			resp = append(resp, (msg.Value).(*MessageResponse))
			msg = msg.Next()
		}

		render.JSON(w, resp)
	}
}
