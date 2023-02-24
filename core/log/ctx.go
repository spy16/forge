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
	f, _ := ctx.Value(fieldsKey).(core.M)
	return f
}
