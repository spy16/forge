package jwtauth

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/errors"
	"github.com/spy16/forge/pkg/strutils"
)

func New(ttl time.Duration, secret string) (*JWTAuth, error) {
	if ttl <= 0 {
		ttl = 1 * time.Hour
	}
	if secret == "" {
		return nil, errors.InvalidInput.Hintf("jwt secret must be set")
	}
	return &JWTAuth{
		jwtTTL:    ttl,
		jwtSecret: secret,
	}, nil
}

type JWTAuth struct {
	jwtTTL    time.Duration
	jwtSecret string
}

func (auth *JWTAuth) CreateSession(ctx context.Context, u core.User) (*core.Session, string, error) {
	now := time.Now()
	expiresAt := now.Add(auth.jwtTTL)
	sessionID := strutils.RandStr(8)

	claims := tokenClaims{
		Subject:   u.ID,
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt.Unix(),
		SessionID: sessionID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &claims)

	tokenString, err := token.SignedString([]byte(auth.jwtSecret))
	if err != nil {
		return nil, "", errors.InternalIssue.CausedBy(err)
	}

	return &core.Session{
		ID:        sessionID,
		UserID:    u.ID,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}, tokenString, nil
}

func (auth *JWTAuth) RestoreSession(ctx context.Context, token string) (*core.Session, error) {
	var errToken = errors.MissingAuth.Coded("invalid_token")

	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errToken.Hintf("empty token")
	}

	keyFn := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errToken.Hintf("invalid alg=%s", token.Method.Alg())
		}
		return []byte(auth.jwtSecret), nil
	}

	tok, err := jwt.ParseWithClaims(token, &tokenClaims{}, keyFn)
	if err != nil || !tok.Valid {
		return nil, errToken.CausedBy(err).Hintf("parse failed")
	}

	claims, ok := tok.Claims.(*tokenClaims)
	if !ok {
		return nil, errToken.Hintf("wrong claims type='%s'", reflect.TypeOf(tok.Claims))
	}

	return &core.Session{
		ID:        claims.SessionID,
		UserID:    claims.Subject,
		CreatedAt: time.Unix(claims.IssuedAt, 0),
		ExpiresAt: time.Unix(claims.ExpiresAt, 0),
	}, nil
}
