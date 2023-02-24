package forge

import (
	"context"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/errors"
	"github.com/spy16/forge/core/router"
)

const defDBSpec = "postgres://postgres@localhost:5432/forge?sslmode=disable"

var namePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_]*[A-Za-z0-9]$`)

// Option values can be used to customise the forging process.
// Note: Currently no values exist. Provided for future usage.
type Option func(app *pgApp) error

// Forge initialises a new app instance from the given configs. This
// includes connecting to postgres & executing the necessary migrations.
func Forge(ctx context.Context, name string, confL core.ConfLoader, opts ...Option) (core.App, error) {
	if !namePattern.MatchString(name) {
		return nil, errors.InvalidInput.
			Hintf("name must match '%s'", namePattern)
	}

	app := &pgApp{
		ctx:   ctx,
		name:  name,
		confL: confL,
	}
	for _, opt := range withDefaults(opts) {
		if err := opt(app); err != nil {
			return nil, err
		}
	}

	pgSpec := app.confL.String("forge.db.spec", defDBSpec)
	if !strings.HasPrefix(pgSpec, "postgres://") {
		return nil, errors.New("db.spec must be valid postgres address")
	}

	conn, err := pgx.Connect(ctx, pgSpec)
	if err != nil {
		return nil, err
	}
	app.conn = conn

	// TODO: initialise all service modules.

	appRouter, err := router.New(app)
	if err != nil {
		_ = conn.Close(ctx)
		return nil, err
	}
	app.router = appRouter

	return app, nil
}

func withDefaults(opts []Option) []Option {
	return append([]Option{
		// TODO: add any default options here.
	}, opts...)
}
