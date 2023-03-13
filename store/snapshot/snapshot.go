package snapshot

import (
	"context"
	_ "embed"
	"time"

	"talkee/core"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *snapshotStore {
	return &snapshotStore{db: db}
}

type snapshotStore struct {
	db *sqlx.DB
}

func (s *snapshotStore) GetSnapshots(ctx context.Context, from time.Time, limit int) ([]*core.Snapshot, error) {
	if !from.IsZero() {
		from, _ = time.Parse(time.RFC3339, time.RFC3339)
	}
	query, args, err := s.db.BindNamed(stmtGetMany, map[string]interface{}{
		"from":  from,
		"limit": limit,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	snaps, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	return snaps, nil
}

func (s *snapshotStore) GetSnapshot(ctx context.Context, snapshotID string) (*core.Snapshot, error) {
	query, args, err := s.db.BindNamed(stmtGetByID, map[string]interface{}{
		"id": snapshotID,
	})

	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	ss := &core.Snapshot{}

	if err := scanRow(rows, ss); err != nil {
		return nil, err
	}

	return ss, nil
}

func (s *snapshotStore) SetSnapshots(ctx context.Context, snapshots []*core.Snapshot) error {

	tx := s.db.MustBegin()
	for _, ss := range snapshots {
		query, args, err := s.db.BindNamed(stmtInsert, map[string]interface{}{
			"snapshot_id":      ss.SnapshotID,
			"trace_id":         ss.TraceID,
			"source":           ss.Source,
			"transaction_hash": ss.TransactionHash,
			"receiver":         ss.Receiver,
			"sender":           ss.Sender,
			"type":             ss.Type,
			"user_id":          ss.UserID,
			"opponent_id":      ss.OpponentID,
			"asset_id":         ss.AssetID,
			"amount":           ss.Amount,
			"memo":             ss.Memo,
		})

		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

//go:embed sql/get_many.sql
var stmtGetMany string

//go:embed sql/get_by_id.sql
var stmtGetByID string

//go:embed sql/insert.sql
var stmtInsert string
