package forge

import (
	"regexp"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/spy16/forge/core/errors"
)

var namePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_]*[A-Za-z0-9]$`)

// Option values can be used to customise the forging process.
type Option func(app *forgeApp) error

// Forge initialises a new app instance from the given options. All modules
// will be initialised, all routes will be setup and a fiber app instance
// ready-for-use will be returned.
func Forge(name string, opts ...Option) (*gin.Engine, error) {
	if !namePattern.MatchString(name) {
		return nil, errors.InvalidInput.
			Hintf("name must match '%s'", namePattern)
	}

	ge := gin.New()
	ge.Use(
		gin.Recovery(),
		requestid.New(),
		extractReqCtx(),
		requestLogger(),
	)

	app := &forgeApp{
		name: name,
		ginE: ge,
	}

	for _, opt := range withDefaults(opts) {
		if err := opt(app); err != nil {
			return nil, err
		}
	}

	if err := app.postCb(app, app.ginE); err != nil {
		return nil, err
	}

	return app.ginE, nil
}
