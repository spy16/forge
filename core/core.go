package core

import (
	"context"
	"net/http"
	"time"
)

// Auth implementation is responsible for validating access tokens
// and restoring user-session from it.
type Auth interface {
	Authenticate(ctx context.Context, token string) (*Session, error)
}

// UserRegistry implementation is responsible for maintaining user
// data.
type UserRegistry interface {
	Get(ctx context.Context, key string) (*User, error)
	Upsert(ctx context.Context, u User) (*User, error)
}

// Session represents a login-session for the contained user.
type Session struct {
	User   User      `json:"user"`
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
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

// Middleware is convenience type for http middlewares.
type Middleware func(http.Handler) http.Handler
