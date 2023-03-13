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
		ID                  uint64 `db:"id" json:"id"`
		MixinUserID         string `db:"mixin_user_id" json:"mixin_user_id"`
		MixinIdentityNumber string `db:"mixin_identity_number" json:"mixin_identity_number"`
		FullName            string `db:"full_name" json:"full_name"`
		AvatarURL           string `db:"avatar_url" json:"avatar_url"`
		MvmPublicKey        string `db:"mvm_public_key" json:"mvm_public_key"`
		Lang                string `db:"lang" json:"-"`

		CreatedAt *time.Time `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"-"`
		DeletedAt *time.Time `db:"deleted_at" json:"-"`
	}

	UserStore interface {
		GetUser(ctx context.Context, id uint64) (*User, error)
		GetUserByIDs(ctx context.Context, ids []uint64) ([]*User, error)
		GetUserByMixinID(ctx context.Context, mixinUserID string) (*User, error)
		Create(ctx context.Context, user *User) (uint64, error)
		UpdateBasicInfo(ctx context.Context, id uint64, data map[string]interface{}) error
	}
	UserService interface {
		LoginWithMixin(ctx context.Context, token, pubkey, lang string) (*User, error)
	}
)

func (u *User) IsMessengerUser() bool {
	return u.MixinIdentityNumber != "0" && len(u.MixinIdentityNumber) < 10
}
