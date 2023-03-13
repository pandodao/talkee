package attachment

import (
	"context"
	_ "embed"
	"talkee/core"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) core.AttachmentStore {
	return &store{
		db: db,
	}
}

type store struct {
	db *sqlx.DB
}

func (s *store) GetAttachment(ctx context.Context, userID uint64) (*core.Attachment, error) {
	query, args, err := s.db.BindNamed(stmtGetByID, map[string]interface{}{
		"id": userID,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	att, err := scanRow(rows)
	if err != nil {
		return nil, err
	}

	return att, nil
}

func (s *store) GetAttachmentByProvider(ctx context.Context, provider int, itemID string) (*core.Attachment, error) {
	query, args, err := s.db.BindNamed(stmtGetByProvider, map[string]interface{}{
		"provider":         provider,
		"provider_item_id": itemID,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	att, err := scanRow(rows)
	if err != nil {
		return nil, err
	}

	return att, nil
}

func (s *store) CreateAttachment(ctx context.Context, att *core.Attachment) (*core.Attachment, error) {
	query, args, err := s.db.BindNamed(stmtInsert, map[string]interface{}{
		"mime_type": att.MimeType,
		"category":  att.Category,
		"provider":  att.Provider,
		"status":    att.Status,
		"lang":      att.Lang,
		"text":      att.Text,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	attID, err := scanReturnID(rows)
	if err != nil {
		return nil, err
	}

	att.ID = attID

	return att, nil
}

func (s *store) UpdateAttachment(ctx context.Context, id uint64, data map[string]interface{}) error {
	query, args, err := s.db.BindNamed(stmtUpdate, map[string]interface{}{
		"id":               id,
		"size":             data["size"],
		"persisted":        data["persisted"],
		"duration":         data["duration"],
		"path":             data["path"],
		"status":           data["status"],
		"provider_item_id": data["provider_item_id"],
	})

	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return err
}

func (s *store) GetAttachmentsByStatus(ctx context.Context, status int, limit uint64) ([]*core.Attachment, error) {
	query, args, err := s.db.BindNamed(stmtGetByStatus, map[string]interface{}{
		"status": status,
		"limit":  limit,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	atts, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	return atts, nil
}

//go:embed sql/get_by_id.sql
var stmtGetByID string

//go:embed sql/get_by_provider.sql
var stmtGetByProvider string

//go:embed sql/insert.sql
var stmtInsert string

//go:embed sql/update.sql
var stmtUpdate string

//go:embed sql/get_by_status.sql
var stmtGetByStatus string
