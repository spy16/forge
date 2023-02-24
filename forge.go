package forge

import (
	"context"
	"regexp"
	"time"

	"github.com/spy16/forge/builtin/jwtauth"
	"github.com/spy16/forge/builtin/pgbase"
	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/errors"
	"github.com/spy16/forge/core/router"
	"github.com/spy16/forge/pkg/vipercfg"
)

var namePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_]*[A-Za-z0-9]$`)

// Option values can be used to customise the forging process.
type Option func(app *forgedApp) error

// Forge initialises a new app instance from the given configs. This
// includes connecting to postgres & executing the necessary migrations.
func Forge(ctx context.Context, name string, opts ...Option) (core.App, error) {
	if !namePattern.MatchString(name) {
		return nil, errors.InvalidInput.
			Hintf("name must match '%s'", namePattern)
	}

	app := &forgedApp{
		ctx:  ctx,
		name: name,
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

	if app.substrate == nil {
		cfg := pgbase.Config{
			PGSpec: app.confL.String("forge.pg_spec", "postgres://postgres@localhost:5432/forge?sslmode=disable"),
		}
		pgSub, err := pgbase.Connect(ctx, cfg)
		if err != nil {
			return nil, err
		}
		app.substrate = pgSub
	}

	if app.auth == nil {
		au, err := jwtauth.New(
			app.confL.Duration("forge.auth.session_ttl", 1*time.Hour),
			app.confL.String("forge.auth.jwt_secret", ""),
		)
		if err != nil {
			return nil, err
		}
		app.auth = au
	}

	appRouter, err := router.New(app, app.confL)
	if err != nil {
		return nil, err
	}
	app.router = appRouter

	if app.postHook != nil {
		if err := app.postHook(app, app.confL); err != nil {
			return nil, err
		}
	}
	return app, nil
}
