package log

import (
	"context"

	"github.com/spy16/forge/core"
)

type ctxKey string

var fieldsKey = ctxKey("fields")

// Ctx returns a new context with fields injected.
func Ctx(ctx context.Context, fields core.M) context.Context {
	return context.WithValue(ctx, fieldsKey, fields)
}

func fromCtx(ctx context.Context) core.M {
	rc := reqCtxMap(core.FromCtx(ctx))

	f, _ := ctx.Value(fieldsKey).(core.M)
	return mergeFields(rc, []core.M{f})
}

func reqCtxMap(rc core.ReqCtx) core.M {
	return map[string]any{
		"path":        rc.Path,
		"route":       rc.Route,
		"method":      rc.Method,
		"authn":       rc.Authenticated(),
		"remote_addr": rc.RemoteAddr,
		"request_id":  rc.RequestID,
	}
}
