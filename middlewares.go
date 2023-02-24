package forge

import (
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/log"
)

func extractReqCtx() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(core.ReqCtxKey, core.ReqCtx{
			Path:       ctx.Request.URL.Path,
			Route:      ctx.FullPath(),
			Method:     ctx.Request.Method,
			Session:    nil,
			RequestID:  requestid.Get(ctx),
			RemoteAddr: ctx.ClientIP(),
		})

		ctx.Next()
	}
}

func requestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()
		ctx.Next() // call next handler
		status := ctx.Writer.Status()
		fields := core.M{
			"status":  status,
			"latency": time.Since(t),
		}

		if status >= 500 {
			fields["errors"] = ctx.Errors
			log.Error(ctx, "request finished with 5xx", nil, fields)
		} else if status >= 400 {
			fields["errors"] = ctx.Errors
			log.Warn(ctx, "request finished with 4xx", fields)
		} else {
			log.Info(ctx, "request finished", fields)
		}
	}
}
