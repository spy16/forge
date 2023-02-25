package pgbase

import "github.com/spy16/forge/core/errors"

type tokenClaims struct {
	Subject   string `json:"sub"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	SessionID string `json:"sid"`
}

func (tc tokenClaims) Valid() error {
	var errBadClaims = errors.InvalidInput.Coded("invalid_claims")

	if tc.IssuedAt >= tc.ExpiresAt {
		return errBadClaims.Hintf("iat > exp")
	} else if tc.Subject == "" {
		return errBadClaims.Hintf("empty sub")
	}
	return nil
}
