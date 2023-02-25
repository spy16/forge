package forge

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/log"
)

const defRoutePrefix = "/forge"

func (app *forgedApp) setupRoutes() error {
	grp := app.fiber.Group(defRoutePrefix)

	grp.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Status(204).JSON(nil)
	})

	grp.Get("/me", func(ctx *fiber.Ctx) error {
		return nil
	})

	return nil
}

func reqLogger() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		t := time.Now()

		rc := core.FromCtx(ctx.UserContext())
		goCtx := ctx.UserContext()
		goCtx = log.Ctx(goCtx, core.M{
			"path":        ctx.Path(),
			"method":      ctx.Method(),
			"req_id":      rc.RequestID,
			"authn":       rc.Authenticated(),
			"remote_addr": rc.RemoteAddr,
		})

		err := ctx.Next()

		status := ctx.Response().StatusCode()
		if status >= 500 {
			log.Error(goCtx, "request finished with 5xx", err, core.M{
				"status":  status,
				"latency": time.Since(t),
			})
		} else if status >= 400 {
			log.Warn(goCtx, "request finished with 4xx", core.M{
				"error":   err,
				"status":  status,
				"latency": time.Since(t),
			})
		} else {
			log.Info(goCtx, "request finished", core.M{
				"status":  status,
				"latency": time.Since(t),
			})
		}

		return err
	}
}

func extractReqCtx() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		goCtx := ctx.UserContext()
		ctx.SetUserContext(core.NewCtx(goCtx, core.ReqCtx{
			Path:       ctx.Path(),
			Method:     ctx.Method(),
			Session:    nil,
			RequestID:  ctx.Locals("requestid").(string),
			RemoteAddr: ctx.IP(),
		}))
		return ctx.Next()
	}
}
