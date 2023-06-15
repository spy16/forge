package servio

import (
	"encoding/json"
	"net/http"

	"github.com/spy16/forge/core/errors"
	"github.com/spy16/forge/core/log"
)

// JSON encodes 'v' as JSON into the writer.
func JSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.WriteHeader(status)
	if status != http.StatusNoContent {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			log.Error(r.Context(), "failed to write json", err)
		}
	}
}

// JSONErr writes the given error as JSON output. Status code is
// inferred from the error value.
func JSONErr(w http.ResponseWriter, r *http.Request, err error) {
	e := errors.E(err)
	JSON(w, r, e.Status, e)
}
