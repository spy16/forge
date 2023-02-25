package pgbase

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/errors"
	"github.com/spy16/forge/pkg/strutils"
)

var (
	_ core.Auth     = (*PG)(nil)
	_ core.UserBase = (*PG)(nil)
)

// Connect loads postgres connection configs, creates a connection and
// executes any necessary migrations.
func Connect(ctx context.Context, cfg Config) (*PG, error) {
	if err := cfg.sanitise(); err != nil {
		return nil, err
	}

	conn, err := pgx.Connect(ctx, cfg.PGSpec)
	if err != nil {
		return nil, err
	}

	return &PG{cfg: cfg, conn: conn}, nil
}

// PG implements the core.Substrate using PostgresQL database.
type PG struct {
	cfg  Config
	conn *pgx.Conn
}

func (pg *PG) User(ctx context.Context, key string) (*core.User, error) {
	// TODO implement me
	panic("implement me")
}

func (pg *PG) Verify(ctx context.Context, uid, token string) (*core.User, error) {
	// TODO implement me
	panic("implement me")
}

func (pg *PG) SetPwd(ctx context.Context, uid string, pwd string) error {
	// TODO implement me
	panic("implement me")
}

func (pg *PG) SetData(ctx context.Context, uid string, data core.UserData) error {
	// TODO implement me
	panic("implement me")
}

func (pg *PG) Register(ctx context.Context, u core.User, keys []core.UserKey) (*core.User, error) {
	// TODO implement me
	panic("implement me")
}

func (pg *PG) IssueToken(ctx context.Context, u core.User) (*core.Session, string, error) {
	now := time.Now()
	expiresAt := now.Add(pg.cfg.TokenTTL)
	sessionID := strutils.RandStr(8)

	claims := tokenClaims{
		Subject:   u.ID,
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt.Unix(),
		SessionID: sessionID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &claims)

	tokenString, err := token.SignedString([]byte(pg.cfg.JWTSecret))
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

func (pg *PG) VerifyToken(ctx context.Context, token string) (*core.Session, error) {
	var errToken = errors.MissingAuth.Coded("invalid_token")

	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errToken.Hintf("empty token")
	}

	keyFn := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errToken.Hintf("invalid alg=%s", token.Method.Alg())
		}
		return []byte(pg.cfg.JWTSecret), nil
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

func (pg *PG) DB() *pgx.Conn { return pg.conn }
