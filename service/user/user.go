package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"talkee/core"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/passport-go/mvm"
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

func (s *UserService) LoginWithMixin(ctx context.Context, token, pubkey, lang string) (*core.User, error) {
	var cli *mixin.Client
	if lang == "" {
		lang = "en"
	}

	var user = &core.User{
		Lang: lang,
	}

	if token != "" {
		cli = mixin.NewFromAccessToken(token)
		profile, err := cli.UserMe(ctx)
		if err != nil {
			fmt.Printf("err cli.UserMe: %v\n", err)
			return nil, err
		}

		contractAddr, err := mvm.GetUserContract(ctx, profile.UserID)
		if err != nil {
			fmt.Printf("err mvm.GetUserContract: %v\n", err)
			return nil, err
		}

		// if contractAddr is not 0x000..00, it means the user has already registered a mvm account
		// we should not allow the user to login with mixin token
		emptyAddr := common.Address{}
		if contractAddr != emptyAddr {
			return nil, core.ErrBadMvmLoginMethod
		}

		user.MixinUserID = profile.UserID
		user.MixinIdentityNumber = profile.IdentityNumber
		user.FullName = profile.FullName
		user.AvatarURL = profile.AvatarURL

	} else if pubkey != "" {
		addr := common.HexToAddress(pubkey)
		mvmUser, err := mvm.GetBridgeUser(ctx, addr)
		if err != nil {
			fmt.Printf("err mvm.GetBridgeUser: %v\n", err)
			return nil, err
		}
		user.MixinUserID = mvmUser.UserID
		user.MixinIdentityNumber = "0"
		user.FullName = mvmUser.FullName
		user.AvatarURL = ""
		user.MvmPublicKey = pubkey

	} else {
		return nil, core.ErrInvalidAuthParams
	}

	existing, err := s.users.GetUserByMixinID(ctx, user.MixinUserID)
	if err != nil && err != sql.ErrNoRows {
		fmt.Printf("err users.GetUserByMixinID: %v\n", err)
		return nil, err
	}

	// create
	if err == sql.ErrNoRows {
		newUserId, err := s.users.Create(ctx, user)
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
	payload := map[string]interface{}{
		"full_name":  user.FullName,
		"avatar_url": user.AvatarURL,
		"lang":       lang,
	}

	if len(user.Lang) >= 2 {
		user.Lang = strings.ToLower(lang[:2])
	} else {
		user.Lang = "en"
	}

	if len(user.Lang) != 0 {
		payload["lang"] = user.Lang
	}

	if err := s.users.UpdateBasicInfo(ctx, existing.ID, payload); err != nil {
		fmt.Printf("err users.Updates: %v\n", err)
		return nil, err
	}

	return existing, nil
}
