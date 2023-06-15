package forge

import (
	"github.com/go-chi/chi/v5"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/vipercfg"
)

// PreContext is the app state before fully initialised.
type PreContext interface {
	Configs() core.ConfLoader
	SetAuth(auth core.Auth)
	SetRouter(r chi.Router)
}

// PostContext is the app state after fully initialised.
type PostContext interface {
	Auth() core.Auth
	Router() chi.Router
	Configs() core.ConfLoader
	Authenticate() Middleware
}

// Option can be passed to Forge() to control the forging process.
type Option func(app *appForge) error

// WithConfLoader can be used to set a custom config loader.
func WithConfLoader(confL core.ConfLoader) Option {
	return func(app *appForge) error {
		if confL == nil {
			viperLoader, err := vipercfg.Init(vipercfg.WithName(app.name))
			if err != nil {
				return err
			}
			confL = viperLoader
		}

		app.confL = confL
		return nil
	}
}

// WithPreHook can be used to set a pre-hook for Forge(). This hook will be invoked
// when config-loader is initialized. Auth and other modules can be initialized
// here.
func WithPreHook(hook func(app PreContext) error) Option {
	return func(app *appForge) error {
		if hook == nil {
			hook = func(app PreContext) error { return nil }
		}
		app.pre = hook
		return nil
	}
}

// WithPostHook can be used to set a post-hook for Forge(). This hook will be invoked
// when all modules and base router is initialised. This can be used to set-up additional
// routes, etc.
func WithPostHook(hook func(app PostContext) error) Option {
	return func(app *appForge) error {
		if hook == nil {
			hook = func(app PostContext) error { return nil }
		}
		app.post = hook
		return nil
	}
}

func withDefaults(opts []Option) []Option {
	return append([]Option{
		WithConfLoader(nil),
		WithPreHook(nil),
		WithPostHook(nil),
	}, opts...)
}
