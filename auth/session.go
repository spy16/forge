package auth

import (
	"time"

	"github.com/spy16/forge/core/errors"
)

type Session struct {
	ID        string
	Token     string
	UserID    string
	UserKind  string
	ExpiresAt time.Time
	RequestID string
}

type sessionClaims struct {
	ID        string `json:"tid"`
	Kind      string `json:"kind"`
	Subject   string `json:"sub"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

func (sc sessionClaims) Valid() error {
	var errInvalid = errors.InvalidInput.Coded("invalid_claims")

	if sc.Kind == "" {
		return errInvalid.Hintf("empty kind claim")
	} else if sc.IssuedAt >= sc.ExpiresAt {
		return errInvalid.Hintf("iat > exp")
	} else if sc.Subject == "" {
		return errInvalid.Hintf("empty sub claim")
	}
	return nil
}
