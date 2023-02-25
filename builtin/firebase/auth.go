package firebase

import (
	"context"

	"github.com/spy16/forge/core"
)

// Auth implements auth module assuming user management is done
// externally by firebase and the token is issued by firebase.
type Auth struct {
}

// VerifyToken verifies the given firebase id-token.
func (fau *Auth) VerifyToken(ctx context.Context, token string) (*core.Session, error) {
	return nil, nil
}
