package core

import (
	"context"
	"time"
)

type (
	CommentFavourite struct {
		ID        uint64 `json:"id"`
		UserID    uint64 `json:"user_id"`
		CommentID uint64 `json:"comment_id"`

		CreatedAt *time.Time `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
		DeletedAt *time.Time `json:"-"`

		Creator *User   `gorm:"-" json:"creator"`
		Reward  *Reward `gorm:"-" json:"reward"`
	}

	CommentFavouriteStore interface {
		// SELECT
		// 	I
		// FROM "favourites"
		// WHERE
		// 	"comment_id" = @commentID
		// AND
		// 	"deleted_at" is NULL
		// ;
		FindAllCommentFavourites(ctx context.Context, commentID uint64) ([]*CommentFavourite, error)

		// SELECT
		// 	*
		// FROM "favourites"
		// WHERE
		//  "comment_id" IN (@commentIDs)
		// AND
		// 	"user_id"=@userID
		// AND
		// 	"deleted_at" is NULL;
		FindUserCommentFavourites(ctx context.Context, userID uint64, commentIDs ...uint64) ([]*CommentFavourite, error)

		// SELECT
		// 	count("id")
		// FROM @@table
		CountAllFavourites(ctx context.Context) (uint64, error)
	}
)
