package forge

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/errors"
	"github.com/spy16/forge/core/log"
)

// Authenticate middleware can be included to restrict access to
// authenticated users only.
func (app *forgeApp) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errAuth := ctx.Error(errors.MissingAuth)
		if app.auth == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errAuth)
			return
		}

		token := extractToken(ctx)
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errAuth)
			return
		}

		se, err := app.auth.Authenticate(ctx, token)
		if err != nil {
			if errors.OneOf(err, []error{errors.NotFound, errors.InvalidInput, errors.MissingAuth}) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errAuth)
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.InternalIssue.CausedBy(err))
			}
			return
		}

		rc := core.FromCtx(ctx)
		rc.Session = se
		ctx.Set(core.ReqCtxKey, rc)
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

func extractToken(c *gin.Context) string {
	var token string
	const bearerPrefix = "Bearer "
	if authH := c.GetHeader("Authorization"); strings.HasPrefix(authH, bearerPrefix) {
		return strings.TrimPrefix(authH, bearerPrefix)
	} else {
		authCookie, err := c.Cookie("_forge_auth")
		if err == nil {
			token = authCookie
		}
	}
	return strings.TrimSpace(token)
}
