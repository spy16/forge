package core

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spy16/forge/pkg/errors"
)

const keyIDSeparator = "/"

// AuthKey kinds.
const (
	KeyKindID       = "id"
	KeyKindEmail    = "email"
	KeyKindUsername = "username"
)

var keyKindPattern = regexp.MustCompile(`^[A-Za-z_]+$`)

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
