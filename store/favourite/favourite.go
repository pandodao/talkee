package favourite

import (
	"context"
	"database/sql"
	_ "embed"
	"talkee/core"
	"talkee/store/user"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *store {
	return &store{db: db}
}

type Config struct {
}

type store struct {
	db *sqlx.DB
}

//go:embed sql/get_comment_favourites.sql
var stmtGetCommentFavourites string

//go:embed sql/get_user_comment_favourite_ids.sql
var StmtGetUserCommentFavouriteIDs string

//go:embed sql/count_all.sql
var StmtCountAll string

func (s *store) FindAllCommentFavourites(ctx context.Context, commentID uint64) ([]*core.CommentFavourite, error) {
	var favors []*core.CommentFavourite
	if err := s.db.Select(&favors, stmtGetCommentFavourites, commentID); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if len(favors) == 0 {
		return favors, nil
	}

	userIDs := []uint64{}
	for _, f := range favors {
		userIDs = append(userIDs, f.UserID)
	}

	var users []*core.User
	query, args, err := sqlx.In(user.StmtGetUserByIDs, userIDs)
	if err != nil {
		return nil, err
	}

	query = s.db.Rebind(query)
	if err := s.db.Select(&users, query, args...); err != nil {
		return nil, err
	}

	for _, f := range favors {
		for _, u := range users {
			if f.UserID == u.ID {
				f.Creator = u
			}
		}
	}

	return favors, nil
}

func (s *store) FindUserCommentFavourites(ctx context.Context, userID uint64, commentIDs ...uint64) ([]*core.CommentFavourite, error) {

	var favs []*core.CommentFavourite
	query, args, err := sqlx.In(StmtGetUserCommentFavouriteIDs, commentIDs, userID)
	if err != nil {
		return nil, err
	}

	query = s.db.Rebind(query)
	if err := s.db.Select(&favs, query, args...); err != nil {
		return nil, err
	}

	return favs, nil
}

func (s *store) CountAllFavourites(ctx context.Context) (uint64, error) {
	sum := []uint64{0}
	if err := s.db.Select(&sum, StmtCountAll); err != nil {
		if err != sql.ErrNoRows {
			return 0, err
		}
	}

	return sum[0], nil
}
