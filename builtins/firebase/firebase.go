package firebase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/errors"
)

// Auth implements auth module for firebase-based user management.
type Auth struct {
	keys      KeySource
	Registry  core.UserRegistry
	ProjectID string
}

func (au *Auth) Authenticate(ctx context.Context, token string) (*core.Session, error) {
	tok, err := jwt.ParseWithClaims(token, &tokClaims{}, au.keyFunc)
	if err != nil {
		return nil, errors.MissingAuth.Hintf(err.Error())
	}
	claims := tok.Claims.(*tokClaims)

	if claims.Issuer != au.issuer() {
		return nil, errors.MissingAuth.Hintf("iss mismatch")
	}

	localUser, err := au.upsertLocalUser(ctx, *claims)
	if err != nil {
		return nil, err
	}

	return &core.Session{
		User:   *localUser,
		Token:  token,
		Expiry: time.Unix(claims.ExpiresAt, 0),
	}, nil
}

func (au *Auth) upsertLocalUser(ctx context.Context, claims tokClaims) (*core.User, error) {
	now := time.Now()

	u := core.User{
		ID: claims.Subject,
		Data: map[string]any{
			"name":    claims.Name,
			"picture": claims.Picture,
		},
		Email:     claims.Email,
		Username:  fmt.Sprintf("user%s", claims.Subject),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if claims.EmailVerified {
		u.VerifiedAt = &now
	}

	if au.Registry != nil {
		return au.Registry.Upsert(ctx, u)
	}
	return &u, nil
}

func (au *Auth) keyFunc(token *jwt.Token) (interface{}, error) {
	alg, ok := token.Header["alg"].(string)
	if !ok || alg != "RS256" {
		return nil, errors.MissingAuth.
			CausedBy(fmt.Errorf("invalid alg=%v", token.Header["alg"]))
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.MissingAuth.
			CausedBy(fmt.Errorf("invalid kid=%v", token.Header["kid"]))
	}

	return au.keys.Find(context.Background(), kid)
}

func (au *Auth) issuer() string {
	return fmt.Sprintf("https://securetoken.google.com/" + au.ProjectID)
}
