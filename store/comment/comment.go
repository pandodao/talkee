package comment

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"talkee/core"
	"talkee/store/reward"

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

//go:embed sql/get_comments.sql
var stmtGetComments string

//go:embed sql/get_comment.sql
var stmtGetComment string

//go:embed sql/count_all.sql
var stmtCountAllComments string

//go:embed sql/count_comments.sql
var stmtCountComments string

//go:embed sql/insert_comment.sql
var stmtInsertComment string

//go:embed sql/update_comment_meta.sql
var stmtUpdateCommentMeta string

//go:embed sql/get_arweave_syncer_list.sql
var stmtGetArweaveSyncerList string

//go:embed sql/update_comment_hash.sql
var stmtUpdateCommentTxhash string

//go:embed sql/fav_comment.sql
var stmtFavComment string

//go:embed sql/unfav_comment.sql
var stmtUnfavComment string

//go:embed sql/inc_comment_fav.sql
var stmtIncCommentFav string

//go:embed sql/dec_comment_fav.sql
var stmtDecCommentFav string

//go:embed sql/get_fav_comment_existence.sql
var stmtGetFavCommentExistence string

//go:embed sql/get_all_by_site_slug.sql
var stmtGetAllCommentsBySiteSlug string

//go:embed sql/get_site_strategy_id.sql
var stmtGetSiteStrategyID string

func (s *store) GetComments(ctx context.Context, siteID uint64, slug string, offset, limit uint64, orderBy, order string) ([]*core.Comment, error) {
	query, args, err := s.db.BindNamed(fmt.Sprintf(stmtGetComments, orderBy, order), map[string]interface{}{
		"site_id": siteID,
		"slug":    slug,
		"limit":   limit,
		"offset":  offset,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	comments, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *store) GetCommentsWithRewards(ctx context.Context, siteID uint64, slug string, offset, limit uint64, orderBy, order string) ([]*core.Comment, error) {
	comments, err := s.GetComments(ctx, siteID, slug, offset, limit, orderBy, order)
	if len(comments) == 0 {
		return comments, err
	}

	if err := s.WithRewards(comments...); err != nil {
		return nil, err
	}

	return comments, err
}

func (s *store) GetComment(ctx context.Context, id uint64) (*core.Comment, error) {
	query, args, err := s.db.BindNamed(stmtGetComment, map[string]interface{}{
		"id": id,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	comment, err := scanRow(rows)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *store) GetCommentWithReward(ctx context.Context, id uint64) (*core.Comment, error) {
	query, args, err := s.db.BindNamed(stmtGetComment, map[string]interface{}{
		"id": id,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	comment, err := scanRow(rows)
	if err != nil {
		return nil, err
	}

	var rewards []*core.Reward
	query, args, err = sqlx.In(reward.StmtGetRewardsByCommentIDs, []uint64{comment.ID}, core.RewardObjectTypeComment)
	if err != nil {
		return nil, err
	}
	query = s.db.Rebind(query)
	if err := s.db.Select(&rewards, query, args...); err != nil {
		return nil, err
	}

	if len(rewards) > 0 {
		comment.Reward = rewards[0]
	}

	return comment, nil
}

func (s *store) CountComments(ctx context.Context, siteID uint64, slug string) (uint64, error) {
	stmt := stmtCountComments
	argMap := map[string]interface{}{
		"site_id": siteID,
		"slug":    slug,
	}
	if siteID == 0 {
		stmt = stmtCountAllComments
		argMap = map[string]interface{}{}
	}

	query, args, err := s.db.BindNamed(stmt, argMap)

	if err != nil {
		return 0, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	count, err := scanReturnID(rows)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *store) GetAllCommentsBySiteSlug(ctx context.Context, siteID uint64, slug string) ([]*core.Comment, error) {
	var comments []*core.Comment

	if err := s.db.Select(&comments, stmtGetAllCommentsBySiteSlug, siteID, slug); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if len(comments) == 0 {
		return comments, nil
	}

	userIDs := []uint64{}
	for _, c := range comments {
		userIDs = append(userIDs, c.UserID)
	}

	var users []*core.User
	// query, args, err := sqlx.In(user.StmtGetUserByIDs, userIDs)
	// if err != nil {
	// 	return nil, err
	// }

	// query = s.db.Rebind(query)
	// if err := s.db.Select(&users, query, args...); err != nil {
	// 	return nil, err
	// }

	for _, c := range comments {
		for _, u := range users {
			if c.UserID == u.ID {
				c.Creator = u
			}
		}
	}

	return comments, nil
}

func (s *store) CreateComment(ctx context.Context, userID, siteID uint64, slug, content string) (uint64, error) {

	tx, err := s.db.Beginx()
	if err != nil {
		return 0, nil
	}
	defer tx.Rollback()

	query, args, err := tx.BindNamed(stmtInsertComment, map[string]interface{}{
		"user_id": userID,
		"site_id": siteID,
		"slug":    slug,
		"content": content,
	})

	if err != nil {
		return 0, err
	}

	rows, err := tx.QueryxContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	id, err := scanReturnID(rows)
	if err != nil {
		return 0, err
	}

	// insert reward task when first comment / site reward strategy set 1
	var strategyID uint
	if err := tx.Get(&strategyID, stmtGetSiteStrategyID, siteID); err != nil {
		return 0, err
	}

	if strategyID > 0 {
		var task core.RewardTask
		if err := tx.Get(&task, reward.StmtGetRewardtaskExistence, siteID, slug); err != nil {
			if err != sql.ErrNoRows {
				return 0, err
			}
		}

		if task.ID == 0 {
			_, err := tx.Exec(reward.StmtInsertRewardTask, siteID, slug, 0)
			if err != nil {
				return 0, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *store) FavComment(ctx context.Context, id, userID uint64, fav bool) (err error) {
	query, args, err := s.db.BindNamed(stmtGetFavCommentExistence, map[string]interface{}{
		"comment_id": id,
		"user_id":    userID,
	})
	if err != nil {
		return err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return err
	}

	existed, _ := scanRowFav(rows)

	if fav && existed != nil {
		return nil
	}

	if !fav && existed == nil {
		return nil
	}

	stmt1 := stmtFavComment
	stmt2 := stmtIncCommentFav
	if !fav {
		stmt1 = stmtUnfavComment
		stmt2 = stmtDecCommentFav
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()

	query1, args1, err := tx.BindNamed(stmt1, map[string]interface{}{
		"comment_id": id,
		"user_id":    userID,
	})
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, query1, args1...); err != nil {
		return err
	}

	query2, args2, err := tx.BindNamed(stmt2, map[string]interface{}{
		"comment_id": id,
	})
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, query2, args2...); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *store) UpdateCommentMeta(ctx context.Context, id uint64, data map[string]interface{}) (err error) {
	query, args, err := s.db.BindNamed(stmtUpdateCommentMeta, data)
	if err != nil {
		return err
	}

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (s *store) UpdateCommentTxHash(ctx context.Context, id uint64, txHash string) (err error) {
	query, args, err := s.db.BindNamed(stmtUpdateCommentTxhash, map[string]interface{}{"arweave_tx_hash": txHash, "id": id})
	if err != nil {
		return err
	}

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (s *store) FindArweaveSyncList(ctx context.Context, limit uint64) ([]*core.Comment, error) {
	query, args, err := s.db.BindNamed(stmtGetArweaveSyncerList, map[string]interface{}{
		"limit": limit,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	comments, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *store) WithRewards(comments ...*core.Comment) error {
	cmtIDs := []uint64{}
	for _, c := range comments {
		cmtIDs = append(cmtIDs, c.ID)
	}

	var rewards []*core.Reward
	query, args, err := sqlx.In(reward.StmtGetRewardsByCommentIDs, cmtIDs, core.RewardObjectTypeComment)
	if err != nil {
		return err
	}

	query = s.db.Rebind(query)
	if err := s.db.Select(&rewards, query, args...); err != nil {
		return err
	}

	for _, r := range rewards {
		for _, c := range comments {
			if r.ObjectID == c.ID {
				c.Reward = r
			}
		}
	}

	return nil
}

func (s *store) WithFavourites(userID uint, comments ...core.Comment) error {

	cmtIDs := []uint64{}
	for _, c := range comments {
		cmtIDs = append(cmtIDs, c.ID)
	}

	var favs []*core.CommentFavourite
	query, args, err := sqlx.In("", cmtIDs, userID)
	if err != nil {
		return err
	}

	query = s.db.Rebind(query)
	if err := s.db.Select(&favs, query, args...); err != nil {
		return err
	}

	for _, f := range favs {
		for _, c := range comments {
			if f.CommentID == c.ID {
				c.FavID = f.ID
			}
		}
	}

	return nil
}
