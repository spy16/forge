package auth

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/spy16/forge/core/errors"
	"github.com/spy16/forge/pkg/httpx"
	"github.com/spy16/forge/pkg/strutils"
)

const keyIDSeparator = "/"

// AuthKey kinds.
const (
	KeyKindID       = "id"
	KeyKindEmail    = "email"
	KeyKindUsername = "username"
)

const defaultUserKind = "user"

var (
	idPattern       = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	usernamePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]+[A-Za-z0-9]$`)
	keyKindPattern  = regexp.MustCompile(`^[A-Za-z_]+$`)
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

// Key represents additional auth key for a user.
type Key struct {
	Key     string         `json:"key"`
	Attribs map[string]any `json:"attribs"`
}

type userCreds struct {
	Kind     string `json:"kind,omitempty"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

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
func (u *User) CheckPassword(pwd string) bool {
	if u.PwdHash == nil {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(*u.PwdHash), []byte(pwd))
	return err == nil
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

func strongPassword(pwd string) bool {
	// TODO: add more logic here.
	return len(pwd) >= 8
}

func (c *userCreds) readFrom(r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			return errors.InvalidInput.CausedBy(err)
		}
		*c = userCreds{
			Email:    r.FormValue("email"),
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}
		if c.Kind == "" {
			c.Kind = defaultUserKind
		}
	} else {
		if err := httpx.ReadJSON(r, c); err != nil {
			return err
		}
	}

	return c.sanitizeAndValidate()
}

func (c *userCreds) sanitizeAndValidate() error {
	var errInvalid = errors.InvalidInput.Coded("invalid_creds")

	if c.Email == "" && c.Username == "" {
		return errInvalid.Hintf("username or email must be specified")
	}

	if c.Email != "" && !strutils.IsValidEmail(c.Email) {
		return errInvalid.Hintf("invalid email")
	}

	if c.Username != "" && !usernamePattern.MatchString(c.Username) {
		return errInvalid.Hintf("invalid username")
	}

	return nil
}

// NewAuthKey returns a new formatted user login-key.
func NewAuthKey(kind, value string) string {
	return fmt.Sprintf("%s%s%s", kind, keyIDSeparator, value)
}

// SplitAuthKey splits the given key-id into its kind and actual value.
func SplitAuthKey(key string) (kind, value string) {
	parts := strings.SplitN(key, keyIDSeparator, 2)
	return parts[0], parts[1]
}

// ValidateAuthKey checks the validity of the login-key.
func ValidateAuthKey(key string) error {
	kind, val := SplitAuthKey(key)

	if !keyKindPattern.MatchString(kind) {
		return errors.InvalidInput.Coded("invalid_kind")
	}

	if len(strings.TrimSpace(val)) != len(val) {
		return errors.InvalidInput.Coded("invalid_value")
	}

	return nil
}
