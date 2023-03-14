package forge

import (
	_ "embed"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/errors"
	"github.com/spy16/forge/core/servio"
)

type forgedApp struct {
	name string
	chi  chi.Router

	auth     core.Auth
	confL    core.ConfLoader
	users    core.UserRegistry
	postHook Hook
}

func (app *forgedApp) Chi() chi.Router          { return app.chi }
func (app *forgedApp) Auth() core.Auth          { return app.auth }
func (app *forgedApp) Users() core.UserRegistry { return app.users }

// Authenticate middleware can be included to restrict access to
// authenticated users only.
func (app *forgedApp) Authenticate() core.Middleware {
	errAuth := errors.MissingAuth
	cookieName := app.confL.String("auth.cookie_name", "_forge_auth")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// auth module is not enabled. all authenticated routes are inaccessible.
			if app.auth == nil {
				servio.JSONErr(w, r, errAuth.Hintf("auth module is disabled"))
				return
			}

			token := extractToken(r, cookieName)
			if token == "" {
				servio.JSONErr(w, r, errAuth.Hintf("invalid token"))
				return
			}

			s, err := app.auth.Authenticate(r.Context(), token)
			if err != nil {
				if errors.OneOf(err, []error{errors.NotFound, errors.InvalidInput, errors.MissingAuth}) {
					servio.JSONErr(w, r, errAuth.Hintf("invalid token"))
				} else {
					servio.JSONErr(w, r, errors.InternalIssue.CausedBy(err))
				}
				return
			}

			ctx := r.Context()
			rc := core.FromCtx(ctx)
			rc.Session = s
			ctx = core.NewCtx(ctx, rc)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request, cookieName string) string {
	var token string
	const bearerPrefix = "Bearer "
	if authH := r.Header.Get("Authorization"); strings.HasPrefix(authH, bearerPrefix) {
		return strings.TrimPrefix(authH, bearerPrefix)
	} else {
		authCookie, err := r.Cookie(cookieName)
		if err == nil && authCookie != nil {
			token = authCookie.Value
		}
	}
	return strings.TrimSpace(token)
}
