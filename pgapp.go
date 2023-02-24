package forge

import (
	"context"
	_ "embed"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/spy16/forge/core"
)

//go:embed schema.sql
var schema string

type pgApp struct {
	ctx      context.Context
	name     string
	conn     *pgx.Conn
	auth     core.Auth
	confL    core.ConfLoader
	router   chi.Router
	userRegs map[string]core.UserRegistry
}

func (app *pgApp) DB() *pgx.Conn { return app.conn }

func (app *pgApp) Auth() core.Auth { return app.auth }

func (app *pgApp) Configs() core.ConfLoader { return app.confL }

func (app *pgApp) Router() chi.Router { return app.router }

func (app *pgApp) UserRegistries() map[string]core.UserRegistry {
	clone := map[string]core.UserRegistry{}
	for kind, registry := range app.userRegs {
		clone[kind] = registry
	}
	return clone
}

func (app *pgApp) initDB(ctx context.Context) error {
	_, err := app.conn.Exec(ctx, schema)
	return err
}
