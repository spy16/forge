package core

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"
)

// Substrate is the base on which the app is forged. Persistence
// is one of the main responsibilities of substrate.
type Substrate interface {
	UserRegistry
}

// App is the contract guaranteed by Forge app implementations.
type App interface {
	Auth() Auth
	Router() chi.Router
	Substrate() Substrate
}

// Auth module provides facilities for issuing tokens, verifying
// tokens, etc.
type Auth interface {
	CreateSession(ctx context.Context, u User) (*Session, string, error)
	RestoreSession(ctx context.Context, token string) (*Session, error)
}

// UserRegistry provides facilities for managing users.
type UserRegistry interface {
	User(ctx context.Context, key string) (*User, error)
	Verify(ctx context.Context, uid, token string) (*User, error)
	SetPwd(ctx context.Context, uid string, pwd string) error
	SetData(ctx context.Context, uid string, data UserData) error
	Register(ctx context.Context, u User) (*User, error)
}

// ConfLoader is responsible for loading configurations during
// initial setup.
type ConfLoader interface {
	Int(key string, defVal int) int
	Bool(key string, defVal bool) bool
	String(key string, defVal string) string
	Strings(key string, defVal []string) []string
	Float64(key string, defVal float64) float64
	Duration(key string, defVal time.Duration) time.Duration
}

// M is an alias for the generic map provided for convenience.
type M = map[string]any
