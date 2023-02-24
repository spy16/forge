package firebase

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokClaims struct {
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	Audience  string `json:"aud"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`

	Name          string `json:"name"`
	Email         string `json:"email"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

func (t *tokClaims) GetIssuer() (string, error)             { return t.Issuer, nil }
func (t *tokClaims) GetSubject() (string, error)            { return t.Subject, nil }
func (t *tokClaims) GetAudience() (jwt.ClaimStrings, error) { return []string{t.Audience}, nil }
func (t *tokClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(t.ExpiresAt, 0)), nil
}
func (t *tokClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(t.IssuedAt, 0)), nil
}
func (t *tokClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(t.IssuedAt, 0)), nil
}
