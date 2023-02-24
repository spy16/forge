package auth

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/spy16/forge/pkg/errors"
	"github.com/spy16/forge/pkg/strutils"
)

// GetUser finds a user by given key.
func (auth *Auth) GetUser(ctx context.Context, authKey string) (*User, error) {
	var u User
	colNames := []string{
		"u.id", "u.kind", "u.user_data", "u.email", "u.pwd_hash", "u.username",
		"u.created_at", "u.updated_at", "u.verified_at", "u.verify_token",
		"u.attributes",
	}

	colPtrs := []any{
		&u.ID, &u.Kind, &u.Data, &u.Email, &u.PwdHash, &u.Username,
		&u.CreatedAt, &u.UpdatedAt, &u.VerifiedAt, &u.VerifyToken,
		&u.Attributes,
	}

	keyKind, val := SplitAuthKey(authKey)

	qb := sq.Select(colNames...).From("users AS u")
	if strutils.OneOf(keyKind, []string{"id", "email", "username"}) {
		qb = qb.Where(sq.Eq{keyKind: val})
	} else {
		qb = qb.
			InnerJoin("user_keys AS uk ON u.id=uk.user_id").
			Where(sq.Eq{"uk.key": authKey})
	}

	q, args, err := qb.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, errors.InternalIssue.CausedBy(err)
	}

	row := auth.conn.QueryRow(ctx, q, args...)
	if err := row.Scan(colPtrs...); err != nil {
		return nil, translateErr(err)
	}

	return &u, nil
}

func (auth *Auth) RegisterUser(ctx context.Context, u User, loginKeys []Key) (*User, error) {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	u.VerifiedAt = nil
	if err := u.Validate(); err != nil {
		return nil, err
	} else if !strutils.OneOf(u.Kind, auth.cfg.EnabledKinds) {
		return nil, errors.InvalidInput.Coded("invalid_kind").
			Hintf("user kind '%s' is not valid", u.Kind)
	}

	tx, err := auth.conn.Begin(ctx)
	if err != nil {
		return nil, translateErr(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Insert user data
	if err := insertUser(ctx, tx, u); err != nil {
		return nil, err
	}

	// Insert login keys
	if err := insertLoginKeys(ctx, tx, u.ID, loginKeys); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, translateErr(err)
	}

	return &u, nil
}

func (auth *Auth) VerifyUser(ctx context.Context, userID, token string) (*User, error) {
	now := time.Now()

	qb := sq.Update("users").
		Where(sq.Eq{
			"id":           userID,
			"verify_token": token,
		}).
		Set("verified_at", now).
		Set("updated_at", now).
		Set("verify_token", sq.Expr("null"))

	q, args, err := qb.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, errors.InternalIssue.CausedBy(err)
	}

	tag, err := auth.conn.Exec(ctx, q, args...)
	if err != nil {
		return nil, translateErr(err)
	} else if tag.RowsAffected() == 0 {
		return nil, errors.NotFound
	}

	return auth.GetUser(ctx, NewAuthKey(KeyKindID, userID))
}

func (auth *Auth) SetPassword(ctx context.Context, id, password string) error {
	now := time.Now()

	pwdHash, err := HashPassword(password)
	if err != nil {
		return err
	}

	q, args, err := sq.Update("users").
		Where(sq.Eq{"id": id}).
		Set("pwd_hash", pwdHash).
		Set("updated_at", now).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return errors.InternalIssue.CausedBy(err)
	}

	_, err = auth.conn.Exec(ctx, q, args...)
	return translateErr(err)
}

func (auth *Auth) SetUserData(ctx context.Context, id string, data UserData) error {
	now := time.Now()

	q, args, err := sq.Update("users").
		Where(sq.Eq{"id": id}).
		Set("user_data", data).
		Set("updated_at", now).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return errors.InternalIssue.CausedBy(err)
	}

	_, err = auth.conn.Exec(ctx, q, args...)
	return translateErr(err)
}

func insertUser(ctx context.Context, tx pgx.Tx, u User) error {
	// Insert into users table.
	colNames := []string{
		"id", "kind", "user_data", "email", "pwd_hash", "username",
		"created_at", "updated_at", "verify_token", "attributes",
	}

	colVals := []any{
		u.ID, u.Kind, u.Data, u.Email, u.PwdHash, u.Username,
		u.CreatedAt, u.UpdatedAt, u.VerifyToken, u.Attributes,
	}

	q, args, err := sq.Insert("users").Columns(colNames...).Values(colVals...).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return errors.InternalIssue.CausedBy(err)
	}

	_, err = tx.Exec(ctx, q, args...)
	return translateErr(err)
}

func insertLoginKeys(ctx context.Context, tx pgx.Tx, uid string, keys []Key) error {
	if len(keys) > 0 {
		insertQ := sq.Insert("user_keys").Columns("key", "user_id", "attribs")
		for _, key := range keys {
			insertQ = insertQ.Values(key.Key, uid, key.Attribs)
		}

		q, args, err := insertQ.PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return errors.InternalIssue.CausedBy(err)
		}

		if _, err := tx.Exec(ctx, q, args...); err != nil {
			return translateErr(err)
		}
	}

	return nil
}
