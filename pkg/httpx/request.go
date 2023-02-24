package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/spy16/forge/core/errors"
)

// ReadJSON decodes the body of the request as JSON into 'ptr'.
func ReadJSON(r *http.Request, ptr any) error {
	if err := json.NewDecoder(r.Body).Decode(ptr); err != nil {
		return errors.InvalidInput.Coded("bad_json").CausedBy(err)
	}
	return nil
}
