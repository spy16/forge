package forge

import (
	"context"
	"net/http"
	"regexp"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/spy16/forge/core/errors"
	"github.com/spy16/forge/core/vipercfg"
)

var namePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_]*[A-Za-z0-9]$`)

// Option values can be used to customise the forging process.
type Option func(app *forgedApp) error

// Forge initialises a new app instance from the given options. All modules
// will be initialised, all routes will be setup and a fiber app instance
// ready-for-use will be returned.
func Forge(ctx context.Context, name string, opts ...Option) (*gin.Engine, error) {
	if !namePattern.MatchString(name) {
		return nil, errors.InvalidInput.
			Hintf("name must match '%s'", namePattern)
	}

	app := &forgedApp{
		name: name,
		ginE: newGin(),
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

	return app.ginE, nil
}

func newGin() *gin.Engine {
	ge := gin.New()

	ge.Use(
		gin.Recovery(),
		requestid.New(),
		extractReqCtx(),
		requestLogger(),
	)

	ge.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound,
			errors.NotFound.Hintf("path not found"))
	})

	ge.HandleMethodNotAllowed = true
	ge.NoMethod(func(c *gin.Context) {
		err := errors.Error{Status: http.StatusMethodNotAllowed}
		c.JSON(http.StatusMethodNotAllowed,
			err.Hintf("method not allowed"))
	})

	return ge
}
