package forge

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/errors"
	"github.com/spy16/forge/core/servio"
)

const defRoutePrefix = "/forge"

var (
	namePattern    = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_]*[A-Za-z0-9]$`)
	errInvalidName = errors.InvalidInput.Hintf("name must match '%s'", namePattern)
)

// Forge forges a new application using given options. If 'conf' is nil, viper-based config
// loader will be initialised. Config file discovery will be done based on the 'name'.
func Forge(name string, opts ...Option) (chi.Router, error) {
	if !namePattern.MatchString(name) {
		return nil, errInvalidName
	}

	forger := &appForge{name: name}
	for _, opt := range withDefaults(opts) {
		if err := opt(forger); err != nil {
			return nil, err
		}
	}

	if err := forger.pre(forger); err != nil {
		return nil, err
	}

	if forger.chi == nil {
		forger.SetRouter(nil)
	}

	if err := forger.setupRoutes(); err != nil {
		return nil, err
	}

	if err := forger.post(forger); err != nil {
		return nil, err
	}

	return forger.chi, nil
}

type Middleware func(http.Handler) http.Handler

type appForge struct {
	name string
	pre  func(preCtx PreContext) error
	post func(postCtx PostContext) error

	// dependencies. set during pre-event. used during post.
	chi   chi.Router
	auth  core.Auth
	confL core.ConfLoader
}

func (app *appForge) Auth() core.Auth          { return app.auth }
func (app *appForge) Router() chi.Router       { return app.chi }
func (app *appForge) Configs() core.ConfLoader { return app.confL }

func (app *appForge) SetAuth(auth core.Auth) { app.auth = auth }
func (app *appForge) SetRouter(r chi.Router) {
	if r == nil {
		r = newChi()
	}

	r.Use(
		middleware.Recoverer,
		middleware.RequestID,
		extractReqCtx(),
		requestLogger(),
	)

	app.chi = r
}

// Authenticate middleware can be included to restrict access to
// authenticated users only.
func (app *appForge) Authenticate() Middleware {
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

			session, err := app.auth.Authenticate(r.Context(), token)
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
			rc.Session = session
			ctx = core.NewCtx(ctx, rc)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (app *appForge) setupRoutes() error {
	app.chi.Route(defRoutePrefix, func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			servio.JSON(w, r, http.StatusNoContent, nil)
		})

		// authenticated routes.
		r.Group(func(r chi.Router) {
			r.Use(app.Authenticate())

			r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
				rc := core.FromCtx(r.Context())
				servio.JSON(w, r, 200, rc.Session.User)
			})
		})
	})
	return nil
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

func newChi() chi.Router {
	ge := chi.NewRouter()

	ge.NotFound(func(w http.ResponseWriter, r *http.Request) {
		servio.JSONErr(w, r, errors.NotFound.Hintf("path not found"))
	})

	ge.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		err := errors.Error{Status: http.StatusMethodNotAllowed}.Hintf("method not allowed")
		servio.JSONErr(w, r, err)
	})

	return ge
}
