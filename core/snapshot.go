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
		SnapshotID      string          `json:"snapshot_id"`
		TraceID         string          `json:"trace_id"`
		Source          string          `json:"source"`
		TransactionHash string          `json:"transaction_hash"`
		Receiver        string          `json:"receiver"`
		Sender          string          `json:"sender"`
		Type            string          `json:"type"`
		UserID          string          `json:"user_id"`
		OpponentID      string          `json:"opponent_id"`
		AssetID         string          `json:"asset_id"`
		Memo            string          `json:"memo"`
		Amount          decimal.Decimal `json:"amount"`
		CreatedAt       time.Time       `json:"created_at"`
	}

	SnapshotStore interface {
		// SELECT
		// 	*
		// FROM snapshots
		// WHERE
		// 	created_at >= @from
		// LIMIT @limit;
		GetSnapshots(ctx context.Context, from time.Time, limit int) ([]*Snapshot, error)

		// SELECT
		// 	*
		// FROM snapshots
		// WHERE
		// 	snapshot_id = @snapshotID
		// LIMIT 1;
		GetSnapshot(ctx context.Context, snapshotID string) (*Snapshot, error)

		// INSERT INTO snapshots
		// (
		// 	"snapshot_id",
		// 	"trace_id",
		// 	"source",
		// 	"transaction_hash",
		// 	"receiver",
		// 	"sender",
		// 	"type",
		// 	"user_id",
		// 	"opponent_id",
		// 	"asset_id",
		// 	"amount",
		// 	"memo",
		// 	"created_at"
		// ) VALUES (
		// 	@snapshot.SnapshotID,
		// 	@snapshot.TraceID,
		// 	@snapshot.Source,
		// 	@snapshot.TransactionHash,
		// 	@snapshot.Receiver,
		// 	@snapshot.Sender,
		// 	@snapshot.Type,
		// 	@snapshot.UserID,
		// 	@snapshot.OpponentID,
		// 	@snapshot.AssetID,
		// 	@snapshot.Amount,
		// 	@snapshot.Memo,
		// 	NOW()
		// );
		SetSnapshot(ctx context.Context, snapshot *Snapshot) error
	}
	SnapshotService interface {
		PollSnapshots(ctx context.Context, offset time.Time, limit int) ([]*Snapshot, error)
	}
)
