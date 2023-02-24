package auth

import (
	"context"
)

type ctxKeyType string

var ctxKey = ctxKeyType("auth_session")

// NewCtx returns a new Go context with auth session injected.
func NewCtx(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, ctxKey, session)
}

// CurSession returns the current auth session from the go context.
// Returns guest session if no value found.
func CurSession(ctx context.Context) *Session {
	v, _ := ctx.Value(ctxKey).(*Session)
	return v
}
