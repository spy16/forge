package forge

import (
	"context"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/spy16/forge/core/errors"
	"github.com/spy16/forge/core/servio"
	"github.com/spy16/forge/core/vipercfg"
)

var namePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_]*[A-Za-z0-9]$`)

// Option values can be used to customise the forging process.
type Option func(app *forgedApp) error

// Forge initialises a new app instance from the given options. All modules
// will be initialised, all routes will be setup and a fiber app instance
// ready-for-use will be returned.
func Forge(ctx context.Context, name string, opts ...Option) (chi.Router, error) {
	if !namePattern.MatchString(name) {
		return nil, errors.InvalidInput.
			Hintf("name must match '%s'", namePattern)
	}

	app := &forgedApp{
		name: name,
		chi:  newChi(),
	}

	for _, opt := range withDefaults(opts) {
		if err := opt(app); err != nil {
			return nil, err
		}
	}

	if app.confL == nil {
		v, err := vipercfg.Init(vipercfg.WithName(app.name))
		if err != nil {
			return nil, err
		}
		app.confL = v
	}

	if err := app.setupRoutes(); err != nil {
		return nil, err
	}

	if app.postHook != nil {
		if err := app.postHook(app, app.confL); err != nil {
			return nil, err
		}
	}

	return app.chi, nil
}

func newChi() chi.Router {
	ge := chi.NewRouter()

	ge.Use(
		middleware.Recoverer,
		middleware.RequestID,
		extractReqCtx(),
		requestLogger(),
	)

	ge.NotFound(func(w http.ResponseWriter, r *http.Request) {
		servio.JSONErr(w, r, errors.NotFound.Hintf("path not found"))
	})

	ge.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		err := errors.Error{Status: http.StatusMethodNotAllowed}.Hintf("method not allowed")
		servio.JSONErr(w, r, err)
	})

	return ge
}
