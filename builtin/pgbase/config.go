package pgbase

import (
	"strings"
	"time"

	"github.com/spy16/forge/core/errors"
)

// Config represents the configuration options for postgres substrate.
type Config struct {
	PGSpec    string
	TokenTTL  time.Duration
	JWTSecret string
}

func (cfg *Config) sanitise() error {
	cfg.PGSpec = strings.TrimSpace(cfg.PGSpec)
	if cfg.PGSpec == "" {
		return errors.InvalidInput.Hintf("pg_spec must be set")
	}

	if len(cfg.JWTSecret) < 16 {
		return errors.InvalidInput.Hintf("jwt_secret must be at-least 16 char")
	}

	if cfg.TokenTTL <= 0 {
		cfg.TokenTTL = 1 * time.Hour
	}
	return nil
}
