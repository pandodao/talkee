package session

import (
	"context"
	"errors"
	"talkee/core"
	"time"

	"github.com/everFinance/goar"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pandodao/passport-go/auth"
	"github.com/patrickmn/go-cache"
)

var (
	ErrKeystoreNotProvided      = errors.New("keystore not provided, use --file or --stdin")
	ErrArweaveWalletNotProvided = errors.New("arweave keystore not provided, use --wallet or --stdin")
	ErrPinNotProvided           = errors.New("pin not provided, use --pin or include in keystore file")
)

type Session struct {
	Version string

	store             *mixin.Keystore
	arweaveWalletPath string
	token             string
	pin               string

	JwtSecret []byte
	cache     *cache.Cache
}

type JwtClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *Session) WithJWTSecret(secret []byte) *Session {
	s.JwtSecret = secret
	return s
}

func (s *Session) WithKeystore(store *mixin.Keystore) *Session {
	s.store = store
	return s
}

func (s *Session) WithAccessToken(token string) *Session {
	s.token = token
	return s
}

func (s *Session) WithPin(pin string) *Session {
	s.pin = pin
	return s
}

func (s *Session) GetKeystore() (*mixin.Keystore, error) {
	if s.store != nil {
		return s.store, nil
	}

	return nil, ErrKeystoreNotProvided
}

func (s *Session) GetPin() (string, error) {
	if s.pin != "" {
		return s.pin, nil
	}

	return "", ErrPinNotProvided
}

func (s *Session) GetClient() (*mixin.Client, error) {
	store, err := s.GetKeystore()
	if err != nil {
		return mixin.NewFromAccessToken(s.token), nil
	}

	return mixin.NewFromKeystore(store)
}

func (s *Session) GetArwallet() (*goar.Wallet, error) {
	if s.arweaveWalletPath == "" {
		s.arweaveWalletPath = "keystores/wallet.json"
	}

	wallet, err := goar.NewWalletFromPath(s.arweaveWalletPath, "https://arweave.net")
	if err != nil {
		return nil, ErrArweaveWalletNotProvided
	}

	return wallet, nil
}

func (s *Session) LoginWithMixin(ctx context.Context, userz core.UserService, authUser *auth.User, lang string) (*core.User, string, error) {
	user, err := userz.LoginWithMixin(ctx, authUser, lang)
	if err != nil {
		return nil, "", err
	}

	expirationTime := time.Now().Add(24 * 365 * time.Hour)
	claims := &JwtClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := newToken.SignedString(s.JwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, tokenString, nil
}
