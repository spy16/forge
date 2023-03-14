package servio

import (
	"encoding/json"
	"net/http"

	"github.com/spy16/forge/core/errors"
)

// BindJSON decodes the request body as JSON value into 'v'.
func BindJSON(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.InvalidInput.CausedBy(err).Hintf("invalid json body")
	}
	return nil
}
