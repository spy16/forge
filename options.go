package forge

import (
	"github.com/spy16/forge/core"
)

// Hook can be used with to further customise an app instance.
type Hook func(core.App, core.ConfLoader) error

// WithAuth sets a custom core.Auth implementation to be used.
func WithAuth(auth core.Auth) Option {
	return func(app *forgedApp) error {
		app.auth = auth
		return nil
	}
}

// WithConfLoader sets a custom core.ConfLoader implementation.
func WithConfLoader(cnfL core.ConfLoader) Option {
	return func(app *forgedApp) error {
		app.confL = cnfL
		return nil
	}
}

// WithUserRegistry sets a custom user registry to be used.
func WithUserRegistry(reg core.UserRegistry) Option {
	return func(app *forgedApp) error {
		app.users = reg
		return nil
	}
}

// WithPostHook will set a hook that will be invoked when the app
// is fully initialised. This hook can be used to set up custom
// routes, etc.
func WithPostHook(hook Hook) Option {
	return func(app *forgedApp) error {
		app.postHook = hook
		return nil
	}
}

func withDefaults(opts []Option) []Option {
	return append([]Option{
		// TODO: add default options.
	}, opts...)
}
