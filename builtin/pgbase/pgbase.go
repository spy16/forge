package pgbase

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/spy16/forge/core"
)

var _ core.Substrate = (*PG)(nil)

// Connect loads postgres connection configs, creates a connection and
// executes any necessary migrations.
func Connect(ctx context.Context, cfg Config) (*PG, error) {
	conn, err := pgx.Connect(ctx, cfg.PGSpec)
	if err != nil {
		return nil, err
	}

	return &PG{conn: conn}, nil
}

// Config represents the configuration options for postgres substrate.
type Config struct {
	PGSpec string
}

// PG implements the core.Substrate using PostgresQL database.
type PG struct {
	conn *pgx.Conn
}

func (pg *PG) DB() *pgx.Conn { return pg.conn }

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

func (pg *PG) Register(ctx context.Context, u core.User) (*core.User, error) {
	// TODO implement me
	panic("implement me")
}
