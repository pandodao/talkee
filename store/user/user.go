package user

import (
	"context"
	_ "embed"

	"talkee/core"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *UserStore {
	return &UserStore{db: db}
}

type Config struct {
}

type UserStore struct {
	db *sqlx.DB
}

func (s *UserStore) GetUserByIDs(ctx context.Context, ids []uint64) ([]*core.User, error) {
	var user []*core.User
	err := s.db.Select(&user, StmtGetUserByIDs, ids)

	return user, err
}

func (s *UserStore) GetUser(ctx context.Context, id uint64) (*core.User, error) {
	query, args, err := s.db.BindNamed(stmtGetByID, map[string]interface{}{
		"id": id,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	user, err := scanRow(rows)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) GetUserByMixinID(ctx context.Context, mixinID string) (*core.User, error) {
	query, args, err := s.db.BindNamed(stmtGetByMixinID, map[string]interface{}{
		"mixin_user_id": mixinID,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	user, err := scanRow(rows)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Create(ctx context.Context, u *core.User) (uint64, error) {
	fullName := u.FullName
	if len(fullName) > 40 {
		fullName = fullName[:40]
	}

	query, args, err := s.db.BindNamed(stmtInsert, map[string]interface{}{
		"mixin_user_id":         u.MixinUserID,
		"full_name":             fullName,
		"mixin_identity_number": u.MixinIdentityNumber,
		"avatar_url":            u.AvatarURL,
		"lang":                  u.Lang,
		"mvm_public_key":        u.MvmPublicKey,
	})

	if err != nil {
		return 0, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	id, err := scanReturnID(rows)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *UserStore) UpdateBasicInfo(ctx context.Context, id uint64, data map[string]interface{}) error {
	fullName := data["full_name"].(string)
	if len(fullName) > 40 {
		fullName = fullName[:40]
	}
	query, args, err := s.db.BindNamed(stmtUpdate, map[string]interface{}{
		"id":         id,
		"full_name":  fullName,
		"avatar_url": data["avatar_url"],
		"lang":       data["lang"],
	})

	if err != nil {
		return err
	}

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

//go:embed sql/get_by_id.sql
var stmtGetByID string

//go:embed sql/get_by_ids.sql
var StmtGetUserByIDs string

//go:embed sql/get_by_mixin_id.sql
var stmtGetByMixinID string

//go:embed sql/insert.sql
var stmtInsert string

//go:embed sql/update.sql
var stmtUpdate string
