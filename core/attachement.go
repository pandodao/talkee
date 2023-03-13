package core

import (
	"context"
	"time"
)

const (
	AttachmentCategoryAudio = iota
)

const (
	AttachmentProviderPlayht = iota
)

const (
	// init .
	// converted -> pending -> download to s3, persist to db -> done
	AttachmentStatusInit = iota
	AttachmentStatusPending
	AttachmentStatusDone
	AttachmentStatusFailed
)

type (
	Attachment struct {
		ID             uint64  `db:"id" json:"id"`
		Size           uint64  `db:"size" json:"size"`
		MimeType       string  `db:"mime_type" json:"mime_type"`
		Persisted      bool    `db:"persisted" json:"persisted"`
		Duration       float64 `db:"duration" json:"duration"`
		Path           string  `db:"path" json:"path"`
		Category       int     `db:"category" json:"category"`
		Provider       int     `db:"provider" json:"provider"`
		ProviderItemID string  `db:"provider_item_id" json:"provider_item_id"`
		Status         int     `db:"status" json:"status"`
		Lang           string  `db:"lang" json:"lang"`
		Text           string  `db:"text" json:"text"`

		CreatedAt *time.Time `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"-"`
	}

	AttachmentStore interface {
		GetAttachment(ctx context.Context, id uint64) (*Attachment, error)
		GetAttachmentByProvider(ctx context.Context, provider int, providerItemID string) (*Attachment, error)
		GetAttachmentsByStatus(ctx context.Context, status int, limit uint64) ([]*Attachment, error)

		CreateAttachment(ctx context.Context, att *Attachment) (*Attachment, error)
		UpdateAttachment(ctx context.Context, id uint64, data map[string]interface{}) error
	}

	AttachmentService interface {
	}
)
