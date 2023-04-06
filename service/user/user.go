package user

import (
	"context"
	"fmt"
	"strings"

	"talkee/core"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/pandodao/passport-go/auth"
	"gorm.io/gorm"
)

func New(client *mixin.Client,
	users core.UserStore,
	cfg Config,
) *UserService {
	return &UserService{
		client: client,
		users:  users,
		cfg:    cfg,
	}
}

type Config struct {
	MixinClientSecret string
}

type UserService struct {
	client *mixin.Client
	users  core.UserStore
	cfg    Config
}

func (s *UserService) LoginWithMixin(ctx context.Context, authUser *auth.User, lang string) (*core.User, error) {
	if lang == "" {
		lang = "en"
	}

	var user = &core.User{
		Lang:                lang,
		MixinUserID:         authUser.UserID,
		MixinIdentityNumber: authUser.IdentityNumber,
		FullName:            authUser.FullName,
		AvatarURL:           authUser.AvatarURL,
		MvmPublicKey:        authUser.MvmAddress.Hex(),
	}

	existing, err := s.users.GetUserByMixinID(ctx, user.MixinUserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Printf("err users.GetUserByMixinID: %v\n", err)
		return nil, err
	}

	// create
	if err == gorm.ErrRecordNotFound {
		newUserId, err := s.users.CreateUser(ctx, user)
		if err != nil {
			fmt.Printf("err users.Create: %v\n", err)
			return nil, err
		}

		user.ID = newUserId

		// create conversation for messenger user.
		if user.IsMessengerUser() {
			if _, err := s.client.CreateContactConversation(ctx, user.MixinUserID); err != nil {
				return nil, err
			}
		}
		return user, nil
	}

	// update
	if len(user.Lang) >= 2 {
		user.Lang = strings.ToLower(lang[:2])
	} else {
		user.Lang = "en"
	}

	if err := s.users.UpdateBasicInfo(ctx, existing.ID, user); err != nil {
		fmt.Printf("err users.Updates: %v\n", err)
		return nil, err
	}

	return existing, nil
}
