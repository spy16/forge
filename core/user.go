package core

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/spy16/forge/core/errors"
	"github.com/spy16/forge/pkg/strutils"
)

const defaultUserKind = "user"

var (
	idPattern       = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	usernamePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]+[A-Za-z0-9]$`)
)

// User represents a registered user in the system.
type User struct {
	ID          string         `json:"id"`
	Kind        string         `json:"kind"`
	Data        UserData       `json:"data"`
	Email       string         `json:"email"`
	PwdHash     *string        `json:"pwd_hash,omitempty"`
	Username    string         `json:"username"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	VerifiedAt  *time.Time     `json:"verified_at"`
	VerifyToken *string        `json:"verify_token,omitempty"`
	Attributes  map[string]any `json:"-"`
}

// UserData represents the standard user profile data.
// Refer https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims
type UserData map[string]any

// Validate validates the user object and returns error if invalid.
func (u *User) Validate() error {
	var errInvalid = errors.InvalidInput.Coded("invalid_user")

	if !idPattern.MatchString(u.ID) {
		return errInvalid.Hintf("invalid id")
	}

	if !usernamePattern.MatchString(u.Username) {
		return errInvalid.Hintf("invalid username")
	}

	if !strutils.IsValidEmail(u.Email) {
		return errInvalid.Hintf("invalid email")
	}
	return nil
}

// Clone returns a deep-clone of the user.
func (u *User) Clone(safe bool) User {
	cloned := User{
		ID:         u.ID,
		Data:       map[string]any{},
		Email:      u.Email,
		Username:   u.Username,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
		VerifiedAt: u.VerifiedAt,
	}

	for k, v := range u.Data {
		cloned.Data[k] = v
	}

	if !safe {
		cloned.PwdHash = u.PwdHash
		cloned.VerifyToken = u.VerifyToken
	}

	return cloned
}

// NewUser returns a new user value with sensible defaults set.
func NewUser(kind, username, email string) User {
	kind = strings.TrimSpace(kind)
	if kind == "" {
		kind = defaultUserKind
	}

	if username == "" {
		username = fmt.Sprintf("user%s", strutils.RandStr(8, strutils.CharsetNums))
	}

	now := time.Now()
	token := strutils.RandStr(10)

	return User{
		ID:          strutils.RandStr(16),
		Kind:        kind,
		Data:        map[string]any{},
		Email:       email,
		Username:    username,
		CreatedAt:   now,
		UpdatedAt:   now,
		VerifiedAt:  nil,
		VerifyToken: &token,
	}
}

// HashPassword hashes and returns the PwdHash value.
func HashPassword(pwd string) (string, error) {
	if !strongPassword(pwd) {
		return "", errors.InvalidInput.Hintf("password is not strong enough")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	if err != nil {
		return "", errors.InternalIssue.CausedBy(err)
	}

	return string(hash), nil
}

// CheckPassword returns true if the given password matches the hashed
// value of the password in the user object. Returns false if mismatch
// or no password is set for user.
func CheckPassword(hash *string, pwd string) bool {
	if hash == nil {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(*hash), []byte(pwd))
	return err == nil
}

func strongPassword(pwd string) bool {
	// TODO: add more logic here.
	return len(pwd) >= 8
}
