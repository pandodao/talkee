package core

import (
	"context"
	"time"
)

const (
	TrialPeriod = 7 * 24 * time.Hour

	DefaultPeriod = 30 * 24 * time.Hour
)

type (
	User struct {
		ID                  uint64 `json:"id"`
		MixinUserID         string `json:"mixin_user_id"`
		MixinIdentityNumber string `json:"mixin_identity_number"`
		FullName            string `json:"full_name"`
		AvatarURL           string `json:"avatar_url"`
		MvmPublicKey        string `json:"mvm_public_key"`
		Lang                string `json:"-"`

		CreatedAt *time.Time `json:"created_at"`
		UpdatedAt *time.Time `json:"-"`
		DeletedAt *time.Time `json:"-"`
	}

	UserStore interface {
		// SELECT
		// 	*
		// FROM @@table
		// WHERE id = @id AND deleted_at IS NULL;
		GetUser(ctx context.Context, id uint64) (*User, error)

		// SELECT
		// 	*
		// FROM @@table
		// WHERE id in (@ids) AND deleted_at IS NULL;
		GetUserByIDs(ctx context.Context, ids []uint64) ([]*User, error)

		// SELECT
		// 	*
		// FROM @@table
		// WHERE mixin_user_id = @mixinUserID;
		GetUserByMixinID(ctx context.Context, mixinUserID string) (*User, error)

		// INSERT INTO users
		// 	(
		// 		"full_name", "avatar_url",
		// 		"mixin_user_id", "mixin_identity_number",
		// 		"lang", "mvm_public_key",
		// 		"created_at", "updated_at"
		// 	)
		// VALUES
		// 	(
		// 		@user.FullName, @user.AvatarURL,
		// 		@user.MixinUserID, @user.MixinIdentityNumber,
		// 		@user.Lang, @user.MvmPublicKey,
		// 		NOW(), NOW()
		// 	)
		// RETURNING id;
		CreateUser(ctx context.Context, user *User) (uint64, error)

		// UPDATE @@table
		// 	{{set}}
		// 	  "name"=@user.FullName,
		// 	  "avatar_url"=@user.AvatarURL,
		// 	  "lang"=@user.Lang,
		// 		"updated_at"=NOW()
		// 	{{end}}
		// WHERE
		// 	"id" = @id AND "deleted_at" IS NULL;
		UpdateBasicInfo(ctx context.Context, id uint64, user *User) error
	}
	UserService interface {
		LoginWithMixin(ctx context.Context, token, pubkey, lang string) (*User, error)
	}
)

func (u *User) IsMessengerUser() bool {
	return u.MixinIdentityNumber != "0" && len(u.MixinIdentityNumber) < 10
}
