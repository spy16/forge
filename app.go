package forge

import (
	"context"
	_ "embed"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/gofiber/fiber/v2"

	"github.com/spy16/forge/core"
)

type forgedApp struct {
	ctx   context.Context
	name  string
	fiber *fiber.App
	ginE  *gin.Engine

	auth     core.Auth
	confL    core.ConfLoader
	postHook Hook
}

func (app *forgedApp) Auth() core.Auth { return app.auth }

func (app *forgedApp) Router() chi.Router { return nil }

func (app *forgedApp) Fiber() *fiber.App { return app.fiber }
