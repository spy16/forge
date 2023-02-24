package auth

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/spy16/forge/pkg/errors"
	"github.com/spy16/forge/pkg/httpx"
	"github.com/spy16/forge/pkg/strutils"
)

// CreateSession creates a new session for the given user and returns.
func (auth *Auth) CreateSession(_ context.Context, u User) (*Session, error) {
	now := time.Now()
	expiresAt := now.Add(auth.cfg.SessionTTL)
	sessionID := strutils.RandStr(8)

	claims := sessionClaims{
		ID:        sessionID,
		Kind:      u.Kind,
		Subject:   u.ID,
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &claims)

	tokenString, err := token.SignedString([]byte(auth.cfg.SigningSecret))
	if err != nil {
		return nil, errors.InternalIssue.CausedBy(err)
	}

	return &Session{
		ID:        sessionID,
		Token:     tokenString,
		UserID:    u.ID,
		UserKind:  u.Kind,
		ExpiresAt: expiresAt,
	}, nil
}

// RestoreSession verifies the given token, restores the session and returns.
// If token is not valid, errors.MissingAuth will be returned.
func (auth *Auth) RestoreSession(_ context.Context, token string) (*Session, error) {
	var errToken = errors.MissingAuth.Coded("invalid_token")

	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errToken.Hintf("empty token")
	}

	keyFn := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errToken.Hintf("invalid alg=%s", token.Method.Alg())
		}
		return []byte(auth.cfg.SigningSecret), nil
	}

	tok, err := jwt.ParseWithClaims(token, &sessionClaims{}, keyFn)
	if err != nil || !tok.Valid {
		return nil, errToken.CausedBy(err).Hintf("parse failed")
	}

	claims, ok := tok.Claims.(*sessionClaims)
	if !ok {
		return nil, errToken.Hintf("claims type='%s'", reflect.TypeOf(tok.Claims))
	}

	return &Session{
		ID:        claims.ID,
		Token:     token,
		UserID:    claims.Subject,
		UserKind:  claims.Kind,
		ExpiresAt: time.Unix(claims.ExpiresAt, 0),
	}, nil
}

// Authenticate returns a middleware that can authenticate incoming
// requests and inject the user into context.
func (auth *Auth) Authenticate() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return httpx.WrapErrH(func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()

			token := extractToken(r, auth.cfg.SessionCookie)
			if token == "" {
				ctx = NewCtx(ctx, nil)
			} else {
				sess, err := auth.RestoreSession(ctx, token)
				if err != nil {
					return err
				}
				ctx = NewCtx(ctx, sess)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
			return nil
		})
	}
}
