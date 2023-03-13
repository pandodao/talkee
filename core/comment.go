package core

import (
	"context"
	"time"
)

const (
	DefaultCommentsPerPage = 20
)

type (
	Comment struct {
		ID            uint64 `db:"id" json:"id"`
		UserID        uint64 `db:"user_id" json:"user_id"`
		SiteID        uint64 `db:"site_id" json:"site_id"`
		Slug          string `db:"slug" json:"slug"`
		FavorCount    uint64 `db:"favor_count" json:"favor_count"`
		ReplyCount    uint64 `db:"reply_count" json:"reply_count"`
		ArweaveTxHash string `db:"arweave_tx_hash" json:"arweave_tx_hash"`
		Content       string `db:"content" json:"content"`

		CreatedAt *time.Time `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"-"`

		Creator *User   `db:"-" json:"creator"`
		Reward  *Reward `db:"-" json:"reward"`
		FavID   uint64  `db:"-" json:"fav_id"`
	}

	CommentChian struct {
		Slug      string     `json:"slug"`
		Content   string     `json:"content"`
		CreatedAt *time.Time `db:"created_at" json:"created_at"`

		MixinUserID         string `json:"mixin_user_id"`
		MixinIdentityNumber string `'json:"mixin_identity_number"`
		FullName            string `json:"full_name"`
	}

	CommentFavourite struct {
		ID        uint64 `db:"id" json:"id"`
		UserID    uint64 `db:"user_id" json:"user_id"`
		CommentID uint64 `db:"comment_id" json:"comment_id"`

		CreatedAt *time.Time `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"-"`

		Creator *User   `db:"-" json:"creator"`
		Reward  *Reward `db:"-" json:"reward"`
	}

	CommentFavouriteStore interface {
		FindAllCommentFavourites(ctx context.Context, commentID uint64) ([]*CommentFavourite, error)
		FindUserCommentFavourites(ctx context.Context, userID uint64, commentIDs ...uint64) ([]*CommentFavourite, error)
		CountAllFavourites(ctx context.Context) (uint64, error)
	}

	CommentStore interface {
		GetComments(ctx context.Context, siteID uint64, slug string, page, limit uint64, orderBy, order string) ([]*Comment, error)
		GetCommentsWithRewards(ctx context.Context, siteID uint64, slug string, page, limit uint64, orderBy, order string) ([]*Comment, error)
		GetComment(ctx context.Context, id uint64) (*Comment, error)
		GetCommentWithReward(ctx context.Context, id uint64) (*Comment, error)
		CreateComment(ctx context.Context, userID, siteID uint64, slug, content string) (uint64, error)
		UpdateCommentMeta(ctx context.Context, id uint64, data map[string]interface{}) error
		UpdateCommentTxHash(ctx context.Context, id uint64, txHash string) (err error)
		CountComments(ctx context.Context, siteID uint64, slug string) (uint64, error)
		GetAllCommentsBySiteSlug(ctx context.Context, siteID uint64, slug string) ([]*Comment, error)

		FindArweaveSyncList(ctx context.Context, limit uint64) ([]*Comment, error)
		FavComment(ctx context.Context, id, userID uint64, fav bool) error
	}

	CommentService interface {
		FindArweaveSyncList(ctx context.Context, limit uint64) ([]*Comment, error)
		SyncToArweave(ctx context.Context, model *Comment) error

		CreateComment(ctx context.Context, userID, siteID uint64, slug, content string) (*Comment, error)
		FavComment(ctx context.Context, id, userID uint64, fav bool) (*Comment, error)
	}
)
