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
		ID            uint64 `json:"id"`
		UserID        uint64 `json:"user_id"`
		SiteID        uint64 `json:"site_id"`
		Slug          string `json:"slug"`
		FavorCount    uint64 `json:"favor_count"`
		ReplyCount    uint64 `json:"reply_count"`
		ArweaveTxHash string `json:"arweave_tx_hash"`
		Content       string `json:"content"`

		CreatedAt *time.Time `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
		DeletedAt *time.Time `json:"-"`

		Creator *User   `gorm:"-" json:"creator"`
		Reward  *Reward `gorm:"-" json:"reward"`
		FavID   uint64  `gorm:"-" json:"fav_id"`
	}

	CommentChian struct {
		Slug      string     `json:"slug"`
		Content   string     `json:"content"`
		CreatedAt *time.Time `json:"created_at"`

		MixinUserID         string `json:"mixin_user_id"`
		MixinIdentityNumber string `'json:"mixin_identity_number"`
		FullName            string `json:"full_name"`
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
