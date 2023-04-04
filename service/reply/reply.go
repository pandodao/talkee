package reply

import (
	"context"

	"talkee/core"
	"talkee/store"
	"talkee/store/comment"
	"talkee/store/reply"
)

func New(
	replys core.ReplyStore,
	comments core.CommentStore,
	users core.UserStore,
	cfg Config,
) *service {
	return &service{
		replys:   replys,
		comments: comments,
		users:    users,
		cfg:      cfg,
	}
}

type Config struct {
}

type service struct {
	replys   core.ReplyStore
	comments core.CommentStore
	users    core.UserStore
	cfg      Config
}

func (s *service) CreateReply(ctx context.Context, userID, commentID uint64, content string) (*core.Reply, error) {
	com, err := s.comments.GetComment(ctx, commentID)
	if err != nil {
		return nil, err
	}

	var re *core.Reply

	if err := store.Transaction(func(tx *store.Handler) error {
		replys := reply.New(tx)
		comments := comment.New(tx)

		if err = comments.IncCommentReplyCount(ctx, com.ID); err != nil {
			return err
		}
		replyID, err := replys.CreateReply(ctx, userID, com.ID, content)
		if err != nil {
			return err
		}
		re, err = replys.GetReply(ctx, replyID)
		if err != nil {
			return err
		}
		return err
	}); err != nil {
		return nil, err
	}

	if err := s.WithUsers(ctx, re); err != nil {
		return nil, err
	}

	return re, nil
}

func (s *service) GetReplies(ctx context.Context, commentID, offset, limit uint64) ([]*core.Reply, error) {
	res, err := s.replys.GetReplies(ctx, commentID, offset, limit)
	if err != nil {
		return nil, err
	}

	if err := s.WithUsers(ctx, res...); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) WithUsers(ctx context.Context, replies ...*core.Reply) error {
	userIDs := []uint64{}
	for _, r := range replies {
		userIDs = append(userIDs, r.UserID)
	}

	users, err := s.users.GetUserByIDs(ctx, userIDs)
	if err != nil {
		return err
	}

	for _, r := range replies {
		for _, u := range users {
			if u.ID == r.UserID {
				r.Creator = u
			}
		}
	}

	return nil
}
