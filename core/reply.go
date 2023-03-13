package core

import (
	"context"
	"time"
)

const (
	DefaultRepliesPerPage = 20
)

type (
	Reply struct {
		ID        uint64 `db:"id" json:"id"`
		UserID    uint64 `db:"user_id" json:"user_id"`
		CommentID string `db:"comment_id" json:"comment_id"`
		Content   string `db:"content" json:"content"`

		CreatedAt *time.Time `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"-"`

		Creator *User `db:"-" json:"creator"`
	}

	ReplyStore interface {
		GetReplies(ctx context.Context, commentID, offset, limit uint64) ([]*Reply, error)
		GetReply(ctx context.Context, replyID uint64) (*Reply, error)
		CountReplies(ctx context.Context, commentID uint64) (uint64, error)
		CreateReply(ctx context.Context, userID, commentID uint64, content string) (uint64, error)
	}

	ReplyService interface {
		CreateReply(ctx context.Context, userID, commentID uint64, content string) (*Reply, error)
	}
)
