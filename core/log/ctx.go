package log

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/spy16/forge/core"
)

var fieldsKey = "forge:log_ctx"

// Ctx returns a new context with fields injected.
func Ctx(ctx context.Context, fieldArr ...core.M) context.Context {
	if len(fieldArr) == 0 {
		return ctx
	}
	return context.WithValue(ctx, fieldsKey, mergeFields(fieldArr[0], fieldArr[1:]))
}

func fromCtx(ctx context.Context) core.M {
	var fields core.M
	if gc, ok := ctx.(*gin.Context); ok {
		val, found := gc.Get(fieldsKey)
		if found {
			fields, _ = val.(core.M)
		}
	} else {
		fields, _ = ctx.Value(fieldsKey).(core.M)
	}

	rc := core.FromCtx(ctx)
	if rc.IsZero() {
		return fields
	}
	return mergeFields(reqCtxMap(rc), []core.M{fields})
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
