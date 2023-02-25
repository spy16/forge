package core

import (
	"context"
)

var ReqCtxKey = "req_ctx"

// ReqCtx represents the context for the current request.
type ReqCtx struct {
	Path       string
	Route      string
	Method     string
	Session    *Session
	RequestID  string
	RemoteAddr string
}

// Authenticated returns true if rc contains authenticated user.
func (rc ReqCtx) Authenticated() bool { return rc.Session != nil }

// NewCtx returns a new Go context with given reqCtx injected.
func NewCtx(ctx context.Context, reqCtx ReqCtx) context.Context {
	return context.WithValue(ctx, ReqCtxKey, reqCtx)
}

// FromCtx returns the request context from Go context.
func FromCtx(ctx context.Context) ReqCtx {
	rc, _ := ctx.Value(ReqCtxKey).(ReqCtx)
	return rc
}
