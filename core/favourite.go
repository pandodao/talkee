package core

import (
	"context"
	"time"
)

type (
	Favourite struct {
		ID        uint64 `json:"id"`
		UserID    uint64 `json:"user_id"`
		CommentID uint64 `json:"comment_id"`

		CreatedAt *time.Time `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
		DeletedAt *time.Time `json:"-"`

		Creator *User   `gorm:"-" json:"creator"`
		Reward  *Reward `gorm:"-" json:"reward"`
	}

	FavouriteStore interface {
		// SELECT
		// 	*
		// FROM "favourites"
		// WHERE
		// 	"comment_id" = @commentID
		// AND
		// 	"deleted_at" is NULL
		// ;
		FindAllFavourites(ctx context.Context, commentID uint64) ([]*Favourite, error)

		// SELECT
		// 	*
		// FROM "favourites"
		// WHERE
		//  "comment_id" IN (@commentIDs)
		// AND
		// 	"user_id"=@userID
		// AND
		// 	"deleted_at" is NULL;
		FindUserFavourites(ctx context.Context, userID uint64, commentIDs []uint64) ([]*Favourite, error)

		// SELECT
		// 	count("id")
		// FROM @@table
		CountAllFavourites(ctx context.Context) (uint64, error)

		// INSERT INTO "favourites"
		// 	(
		// 		"comment_id",
		// 		"user_id",
		// 		"created_at",
		// 		"updated_at"
		// 	)
		// VALUES
		// 	(
		// 		@commentID,
		// 		@userID,
		// 		NOW(),
		// 		NOW()
		// 	)
		// ON CONFLICT ("comment_id", "user_id") DO
		// 	UPDATE
		// 	SET "deleted_at" = NULL, "updated_at" = NOW()
		// ;
		CreateFavourite(ctx context.Context, userID, commentID uint64) error

		// UPDATE "favourites"
		// {{set}}
		//   "deleted_at" = NOW()
		// {{end}}
		// WHERE
		// 	"user_id" = @userID
		// AND
		// 	"comment_id" = @commentID
		// AND
		// 	"deleted_at" is NULL
		// ;
		DeleteFavourite(ctx context.Context, userID, commentID uint64) error
	}
)
