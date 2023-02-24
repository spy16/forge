package httpx

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/spy16/forge/pkg/errors"
)

// MarshalCookie marshals 'val' using JSON, encodes using Base64 and
// returns.
func MarshalCookie(val any) (string, error) {
	data, err := json.Marshal(val)
	if err != nil {
		return "", errors.InternalIssue.CausedBy(err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// UnmarshalCookie reads the cookie value as base64 encoded JSON value.
// Returns true if successful.
func UnmarshalCookie(r *http.Request, key string, into any) bool {
	c, err := r.Cookie(key)
	if err != nil || c == nil || c.Value == "" {
		return false
	}

	data, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		return true
	}

	return json.Unmarshal(data, into) == nil
}
