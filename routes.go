package forge

import (
	"github.com/gin-gonic/gin"

	"github.com/spy16/forge/core"
)

const defRoutePrefix = "/forge"

func (app *forgedApp) setupRoutes() error {
	grp := app.ginE.Group(defRoutePrefix)

	grp.GET("/ping", func(ctx *gin.Context) {
		ctx.Status(204)
	})

	grp.GET("/me", app.Authenticate(), func(ctx *gin.Context) {
		rc := core.FromCtx(ctx)
		ctx.JSON(200, rc.Session.User)
	})

	return nil
}
