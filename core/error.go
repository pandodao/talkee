package core

import "errors"

var (
	ErrCorruptedSnapshotAction = errors.New("corrupted snapshot action")
	ErrUnauthorized            = errors.New("unauthorized")
	ErrInvalidSiteID           = errors.New("invalid site id")
	ErrInvalidSlug             = errors.New("invalid slug")
	ErrInvalidCommentContent   = errors.New("invalid comment content")
	ErrInvalidCommentID        = errors.New("invalid comment id")
	ErrInvalidUserID           = errors.New("invalid user id")
	ErrInvalidAssetID          = errors.New("invalid asset id")
	ErrNoRecord                = errors.New("no record")
	ErrSiteLimitReached        = errors.New("site limit reached")
	ErrInvalidAuthParams       = errors.New("invalid auth params")
	ErrInvalidOrigin           = errors.New("invalid origin")
	ErrBadMvmLoginMethod       = errors.New("bad mvm login method")
	ErrBadMvmLoginSignature    = errors.New("bad mvm login signature")
	ErrBadMvmLoginMessage      = errors.New("bad mvm login message")
)
