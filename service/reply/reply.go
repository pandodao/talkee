package reply

import (
	"context"

	"talkee/core"
)

func New(
	replys core.ReplyStore,
	comments core.CommentStore,
	cfg Config,
) *service {
	return &service{
		replys:   replys,
		comments: comments,
		cfg:      cfg,
	}
}

type Config struct {
}

type service struct {
	replys   core.ReplyStore
	comments core.CommentStore
	cfg      Config
}

func (s *service) CreateReply(ctx context.Context, userID, commentID uint64, content string) (*core.Reply, error) {
	comment, err := s.comments.GetComment(ctx, commentID)
	if err != nil {
		return nil, err
	}

	replyID, err := s.replys.CreateReply(ctx, userID, comment.ID, content)
	if err != nil {
		return nil, err
	}

	reply, err := s.replys.GetReply(ctx, replyID)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
