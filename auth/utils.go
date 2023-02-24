package auth

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/spy16/forge/pkg/errors"
)

const (
	bearerPrefix = "Bearer "
	headerAuthz  = "Authorization"
)

func extractToken(r *http.Request, cookieName string) string {
	var token string
	if authH := r.Header.Get(headerAuthz); strings.HasPrefix(authH, bearerPrefix) {
		token = strings.TrimPrefix(authH, bearerPrefix)
	} else {
		c, err := r.Cookie(cookieName)
		if err != nil || c == nil {
			return ""
		}
		token = c.Value
	}
	return strings.TrimSpace(token)
}

func translateErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
		return errors.NotFound.Coded("not_found")
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return errors.Conflict.Hintf(pgErr.Message)

		case pgerrcode.NoData:
			return errors.NotFound.Hintf(pgErr.Message)
		}
	}

	return errors.InternalIssue.CausedBy(err)
}
