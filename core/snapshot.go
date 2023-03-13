package core

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidTrace          = errors.New("invalid trace")
	ErrInvalidConversationID = errors.New("invalid conversation id")

	DefaultSnapshotFetchCount = 100
)

type (
	Snapshot struct {
		SnapshotID      string          `db:"snapshot_id" json:"snapshot_id"`
		TraceID         string          `db:"trace_id" json:"trace_id"`
		Source          string          `db:"source" json:"source"`
		TransactionHash string          `db:"transaction_hash" json:"transaction_hash"`
		Receiver        string          `db:"receiver" json:"receiver"`
		Sender          string          `db:"sender" json:"sender"`
		Type            string          `db:"type" json:"type"`
		CreatedAt       time.Time       `db:"created_at" json:"created_at"`
		UserID          string          `db:"user_id" json:"user_id"`
		OpponentID      string          `db:"opponent_id" json:"opponent_id"`
		AssetID         string          `db:"asset_id" json:"asset_id"`
		Amount          decimal.Decimal `db:"amount" json:"amount"`
		Memo            string          `db:"memo" json:"memo"`
	}

	SnapshotStore interface {
		GetSnapshots(ctx context.Context, from time.Time, limit int) ([]*Snapshot, error)
		GetSnapshot(ctx context.Context, snapshotID string) (*Snapshot, error)
		SetSnapshots(ctx context.Context, snapshots []*Snapshot) error
	}
	SnapshotService interface {
		PollSnapshots(ctx context.Context, offset time.Time, limit int) ([]*Snapshot, error)
	}
)
