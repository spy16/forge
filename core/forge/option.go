package forge

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/spy16/forge/core"
)

// Hook can be used with to further customise an app instance.
type Hook func(core.App, *gin.Engine) error

// WithAuth sets a custom core.Auth implementation to be used.
func WithAuth(auth core.Auth) Option {
	return func(app *forgeApp) error {
		app.auth = auth
		return nil
	}
}

// WithPostHook will set a hook that will be invoked when the app
// is fully initialised. This hook can be used to set up custom
// routes, etc.
func WithPostHook(hook Hook) Option {
	return func(app *forgeApp) error {
		app.postCb = hook
		return nil
	}
}

// WithStatic sets the static file system to be served on index
// route.
func WithStatic(static http.FileSystem) Option {
	return func(app *forgeApp) error {
		app.ginE.StaticFS("/", static)
		return nil
	}
}

func withDefaults(opts []Option) []Option {
	return append([]Option{
		WithAuth(nil),
		WithPostHook(func(app core.App, engine *gin.Engine) error { return nil }),
	}, opts...)
}
