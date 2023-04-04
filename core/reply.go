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
		ID        uint64 `json:"id"`
		UserID    uint64 `json:"user_id"`
		CommentID string `json:"comment_id"`
		Content   string `json:"content"`

		CreatedAt *time.Time `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
		DeletedAt *time.Time `json:"-"`

		Creator *User `gorm:"-" json:"creator"`
	}

	ReplyStore interface {
		// SELECT
		// 	*
		// FROM "replies"
		// WHERE
		// 	"replies"."comment_id" = @commentID
		// 	AND "deleted_at" IS NULL
		// ORDER BY "replies"."created_at" DESC
		// OFFSET @offset
		// LIMIT @limit;
		GetReplies(ctx context.Context, commentID, offset, limit uint64) ([]*Reply, error)

		// SELECT
		// 	*
		// FROM "replies"
		// WHERE
		// 	"replies"."id" = @replyID
		// 	AND "replies"."deleted_at" IS NULL
		// LIMIT 1;
		GetReply(ctx context.Context, replyID uint64) (*Reply, error)

		// SELECT
		// 	COUNT("id")
		// FROM "replies"
		// WHERE
		//  "comment_id" = @commentID
		// 	AND "deleted_at" IS NULL;
		CountReplies(ctx context.Context, commentID uint64) (uint64, error)

		// INSERT INTO "replies"
		// 	(
		// 		"user_id",
		// 		"comment_id",
		// 		"content",
		// 		"created_at", "updated_at"
		// 	)
		// VALUES
		// 	(
		// 		@userID,
		// 		@commentID,
		// 		@content,
		// 		NOW(), NOW()
		// 	)
		// RETURNING id;
		CreateReply(ctx context.Context, userID, commentID uint64, content string) (uint64, error)
	}

	ReplyService interface {
		CreateReply(ctx context.Context, userID, commentID uint64, content string) (*Reply, error)
		GetReplies(ctx context.Context, commentID, offset, limit uint64) ([]*Reply, error)
	}
)
