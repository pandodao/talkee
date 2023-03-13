package session

import (
	"context"
	"talkee/core"
)

type contextKey struct{}

type key int

const (
	userKey key = iota
	siteIDKey
	siteSlugKey
)

func With(ctx context.Context, s *Session) context.Context {
	return context.WithValue(ctx, contextKey{}, s)
}

func From(ctx context.Context) *Session {
	return ctx.Value(contextKey{}).(*Session)
}

func WithUser(ctx context.Context, user *core.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func WithSiteInfo(ctx context.Context, siteID uint64, siteSlug string) context.Context {
	ctx = context.WithValue(ctx, siteIDKey, siteID)
	return context.WithValue(ctx, siteSlugKey, siteSlug)
}

func UserFrom(ctx context.Context) (*core.User, bool) {
	u, ok := ctx.Value(userKey).(*core.User)
	return u, ok
}

func SiteInfoFrom(ctx context.Context) (uint64, string, bool) {
	siteID, ok1 := ctx.Value(siteIDKey).(uint64)
	slug, ok2 := ctx.Value(siteSlugKey).(string)
	return siteID, slug, ok1 && ok2
}
