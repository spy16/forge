package forge

import "github.com/spy16/forge/core"

// Hook can be used with CLI to further customise app instance.
type Hook func(core.App, core.ConfLoader) error

// WithSubstrate sets a custom core.Substrate implementation to be used.
func WithSubstrate(subs core.Substrate) Option {
	return func(app *forgedApp) error {
		app.substrate = subs
		return nil
	}
}

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

// WithPostHook will set a hook that will be invoked when the app
// is initialised.
func WithPostHook(hook Hook) Option {
	return func(app *forgedApp) error {
		app.postHook = hook
		return nil
	}
}

func withDefaults(opts []Option) []Option {
	return append([]Option{
		// TODO: add any default options here.
	}, opts...)
}
