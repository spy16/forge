package core

import (
	"time"

	"github.com/spy16/forge/core/errors"
)

// Session represents a user-session.
type Session struct {
	ID        string
	UserID    string
	UserKind  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func (sess *Session) validate() error {
	if sess.ID == "" {
		return errors.InvalidInput.Hintf("id must be set")
	}

	if sess.UserID == "" {
		return errors.InvalidInput.Hintf("user_id must be set")
	}

	if sess.UserKind == "" {
		return errors.InvalidInput.Hintf("user_kind must be set")
	}

	return nil
}
