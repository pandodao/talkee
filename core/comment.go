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

	ArweaveSyncListItem struct {
		CommentID           uint64     `json:"comment_id"`
		UserID              uint64     `json:"user_id"`
		SiteID              uint64     `json:"site_id"`
		Slug                string     `json:"slug"`
		FavorCount          uint64     `json:"favor_count"`
		ReplyCount          uint64     `json:"reply_count"`
		ArweaveTxHash       string     `json:"arweave_tx_hash"`
		Content             string     `json:"content"`
		MixinUserID         string     `json:"mixin_user_id"`
		MixinIdentityNumber string     `json:"mixin_identity_number"`
		FullName            string     `json:"full_name"`
		AvatarURL           string     `json:"avatar_url"`
		MvmPublicKey        string     `json:"mvm_public_key"`
		CreatedAt           *time.Time `json:"created_at"`
		UpdatedAt           *time.Time `json:"updated_at"`
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
		// SELECT
		// 	*
		// FROM "comments"
		// WHERE
		// 	"comments"."site_id" = @siteID
		// 	AND "comments"."slug" = @slug
		// 	AND "comments"."deleted_at" IS NULL
		// {{if orderBy == "favor_count"}}
		//   ORDER BY "favor_count" {{if order == "DESC" }} DESC {{else}} ASC {{end}}
		// {{else}}
		//   ORDER BY "created_at" {{if order == "DESC" }} DESC {{else}} ASC {{end}}
		// {{end}}
		// OFFSET @offset
		// LIMIT @limit;
		GetComments(ctx context.Context, siteID uint64, slug string, offset, limit uint64, orderBy, order string) ([]*Comment, error)

		// SELECT
		//  *
		// FROM "comments"
		// WHERE
		// 	"comments"."id" = @id
		// 	AND "comments"."deleted_at" IS NULL
		// LIMIT 1;
		GetComment(ctx context.Context, id uint64) (*Comment, error)

		// INSERT INTO "comments"
		// 	(
		// 		"user_id",
		// 		"site_id",
		// 		"slug",
		// 		"content",
		// 		"created_at", "updated_at"
		// 	)
		// VALUES
		// 	(
		// 		@userID,
		// 		@siteID,
		// 		@slug,
		// 		@content,
		// 		NOW(), NOW()
		// 	)
		// RETURNING id;
		CreateComment(ctx context.Context, userID, siteID uint64, slug, content string) (uint64, error)

		// UPDATE
		// 	"comments"
		// SET
		//  "arweave_tx_hash" = @txHash,
		// 	"updated_at" = NOW()
		// WHERE
		// 	"id" = :id AND "deleted_at" IS NULL;
		// ;
		UpdateCommentTxHash(ctx context.Context, id uint64, txHash string) (err error)

		// SELECT
		// 	COUNT("id")
		// FROM "comments"
		// WHERE
		//  {{if siteID != 0}}
		//    "site_id"=@siteID AND "slug"=@slug AND
		//  {{end}}
		//  "deleted_at" IS NULL;
		CountComments(ctx context.Context, siteID uint64, slug string) (uint64, error)

		// SELECT
		// 	"id",
		// 	"user_id",
		// 	"site_id",
		// 	"slug",
		// 	"favor_count"
		// FROM "comments"
		// WHERE
		//   "site_id" = $1
		// AND
		//   "slug" = $2
		// AND "deleted_at" IS NULL;
		GetAllCommentsBySiteSlug(ctx context.Context, siteID uint64, slug string) ([]*Comment, error)

		// SELECT
		// 	"comments"."id",
		// 	"comments"."user_id",
		// 	"comments"."site_id",
		// 	"comments"."slug",
		// 	"comments"."favor_count", "reply_count",
		// 	"comments"."arweave_tx_hash",
		// 	"comments"."content",

		// 	"users"."mixin_user_id",
		// 	"users"."mixin_identity_number",
		// 	"users"."full_name",
		// 	"users"."avatar_url",
		// 	"users"."mvm_public_key",

		// 	"comments"."created_at",
		// 	"comments"."updated_at"
		// FROM "comments"
		// INNER JOIN "users" ON "comments"."user_id" = "users"."id"
		// WHERE
		// 	"comments"."site_id" in (
		// 		select "id" FROM "sites" WHERE  "use_arweave" = true
		// 	)
		// 	AND  ("comments"."arweave_tx_hash" is NULL OR "comments"."arweave_tx_hash" = '')
		// 	AND "comments"."deleted_at" IS NULL
		// ORDER BY "comments"."created_at" asc
		// LIMIT :limit
		// ;
		FindArweaveSyncList(ctx context.Context, limit uint64) ([]*ArweaveSyncListItem, error)

		// UPDATE "comments"
		// SET
		// 	  {{if fav}}"favor_count" = "favor_count" + 1 {{else}}"favor_count" = "favor_count" - 1{{end}},
		// 	  "updated_at" = NOW()
		// WHERE
		// 	"id" = @id AND "deleted_at" IS NULL;
		FavComment(ctx context.Context, id, userID uint64, fav bool) error

		// UPDATE
		// 	"comments"
		// SET
		// 	"reply_count" = "reply_count" + 1,
		// 	"updated_at" = NOW()
		// WHERE
		// 	"id" = @id AND "deleted_at" IS NULL;
		IncCommentReplyCount(ctx context.Context, id uint64) error
	}

	CommentService interface {
		FindArweaveSyncList(ctx context.Context, limit uint64) ([]*ArweaveSyncListItem, error)
		SyncToArweave(ctx context.Context, model *ArweaveSyncListItem) error

		CreateComment(ctx context.Context, userID, siteID uint64, slug, content string) (*Comment, error)
		FavComment(ctx context.Context, id, userID uint64, fav bool) (*Comment, error)
		GetCommentsWithRewards(ctx context.Context, siteID uint64, slug string, page, limit uint64, orderBy, order string) ([]*Comment, error)
		GetCommentWithReward(ctx context.Context, id uint64) (*Comment, error)
	}
)
