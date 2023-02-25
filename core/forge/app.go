package forge

import (
	"github.com/gin-gonic/gin"

	"github.com/spy16/forge/core"
)

type forgeApp struct {
	name   string
	ginE   *gin.Engine
	auth   core.Auth
	postCb Hook
}

func (app *forgeApp) Auth() core.Auth { return app.auth }
