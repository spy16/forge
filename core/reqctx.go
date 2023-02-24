package core

import (
	"context"
)

var ctxKey = ctxKeyType("req_ctx")

type ctxKeyType string

// ReqCtx represents the context for the current request.
type ReqCtx struct {
	Session    *Session
	Path       string
	Method     string
	RequestID  string
	RemoteAddr string
}

// Authenticated returns true if rc contains authenticated user.
func (rc ReqCtx) Authenticated() bool { return rc.Session != nil }

// NewCtx returns a new Go context with given reqCtx injected.
func NewCtx(ctx context.Context, reqCtx ReqCtx) context.Context {
	return context.WithValue(ctx, ctxKey, reqCtx)
}

// FromCtx returns the request context from Go context.
func FromCtx(ctx context.Context) ReqCtx {
	v, _ := ctx.Value(ctxKey).(ReqCtx)
	return v
}
