package forge

import (
	_ "embed"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/errors"
)

type forgedApp struct {
	name string
	ginE *gin.Engine

	auth     core.Auth
	confL    core.ConfLoader
	users    core.UserRegistry
	postHook Hook
}

func (app *forgedApp) Gin() *gin.Engine         { return app.ginE }
func (app *forgedApp) Auth() core.Auth          { return app.auth }
func (app *forgedApp) Users() core.UserRegistry { return app.users }

// Authenticate middleware can be included to restrict access to
// authenticated users only.
func (app *forgedApp) Authenticate() gin.HandlerFunc {
	errAuth := errors.MissingAuth
	cookieName := app.confL.String("auth.cookie_name", "_forge_auth")

	return func(ctx *gin.Context) {
		// auth module is not enabled. all authenticated routes are inaccessible.
		if app.auth == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ctx.Error(errAuth.Hintf("auth module is disabled")))
			return
		}

		token := extractToken(ctx, cookieName)
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ctx.Error(errAuth))
			return
		}

		s, err := app.auth.Authenticate(ctx, token)
		if err != nil {
			if errors.OneOf(err, []error{errors.NotFound, errors.InvalidInput, errors.MissingAuth}) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, ctx.Error(err))
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError,
					ctx.Error(errors.InternalIssue.CausedBy(err)))
			}
			return
		}

		rc := core.FromCtx(ctx)
		rc.Session = s
		ctx.Set(core.ReqCtxKey, rc)
		ctx.Next()
	}
}

func extractToken(c *gin.Context, cookieName string) string {
	var token string
	const bearerPrefix = "Bearer "
	if authH := c.GetHeader("Authorization"); strings.HasPrefix(authH, bearerPrefix) {
		return strings.TrimPrefix(authH, bearerPrefix)
	} else {
		authCookie, err := c.Cookie(cookieName)
		if err == nil {
			token = authCookie
		}
	}
	return strings.TrimSpace(token)
}
