package httpx

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/pkg/errors"
	"github.com/spy16/forge/pkg/log"
)

// WriteErr writes the error value to the ResponseWriter. HTTP status is
// inferred from the error value.
func WriteErr(w http.ResponseWriter, r *http.Request, err error) {
	e := errors.E(err)
	WriteJSON(w, r, e.Status, e)
}

// WriteJSON marshals 'v' as JSON value and writes to the ResponseWriter.
// No content is written if request-method is HEAD or status is 204.
func WriteJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if status == http.StatusNoContent || r.Method == http.MethodHead {
		return
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error(r.Context(), "failed to write JSON", err, core.M{
			"status":     status,
			"value_type": reflect.TypeOf(v).String(),
		})
	}
}
