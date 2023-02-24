package forge

import (
	"context"
	_ "embed"

	"github.com/go-chi/chi/v5"

	"github.com/spy16/forge/core"
)

type forgedApp struct {
	ctx  context.Context
	name string

	auth      core.Auth
	confL     core.ConfLoader
	router    chi.Router
	postHook  Hook
	substrate core.Substrate
}

func (app *forgedApp) Auth() core.Auth { return app.auth }

func (app *forgedApp) Router() chi.Router { return app.router }

func (app *forgedApp) Configs() core.ConfLoader { return app.confL }

func (app *forgedApp) Substrate() core.Substrate { return app.substrate }
