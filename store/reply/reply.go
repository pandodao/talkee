package reply

import (
	"context"
	_ "embed"
	"talkee/core"

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

//go:embed sql/get_replies.sql
var stmtGetReplies string

//go:embed sql/get_reply.sql
var stmtGetReply string

//go:embed sql/count_replies.sql
var stmtCountReplies string

//go:embed sql/insert_reply.sql
var stmtInsertReply string

//go:embed sql/inc_comment_reply.sql
var stmtIncCommentReply string

func (s *store) GetReplies(ctx context.Context, commentID, offset, limit uint64) ([]*core.Reply, error) {
	query, args, err := s.db.BindNamed(stmtGetReplies, map[string]interface{}{
		"comment_id": commentID,
		"limit":      limit,
		"offset":     offset,
	})
	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	replies, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	return replies, nil
}

func (s *store) GetReply(ctx context.Context, replyID uint64) (*core.Reply, error) {
	query, args, err := s.db.BindNamed(stmtGetReply, map[string]interface{}{
		"reply_id": replyID,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	reply, err := scanRow(rows)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (s *store) CountReplies(ctx context.Context, commentID uint64) (uint64, error) {
	query, args, err := s.db.BindNamed(stmtCountReplies, map[string]interface{}{
		"comment_id": commentID,
	})

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

func (s *store) CreateReply(ctx context.Context, userID, commentID uint64, content string) (uint64, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	query1, args1, err := tx.BindNamed(stmtIncCommentReply, map[string]interface{}{
		"comment_id": commentID,
	})
	if err != nil {
		return 0, err
	}

	if _, err := tx.ExecContext(ctx, query1, args1...); err != nil {
		return 0, err
	}

	query, args, err := tx.BindNamed(stmtInsertReply, map[string]interface{}{
		"user_id":    userID,
		"comment_id": commentID,
		"content":    content,
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

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}
