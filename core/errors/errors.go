package errors

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	Is     = errors.Is
	As     = errors.As
	New    = errors.New
	Errorf = fmt.Errorf
)

var (
	NotFound      = Error{Code: "not_found", Status: http.StatusNotFound}
	Conflict      = Error{Code: "conflict", Status: http.StatusConflict}
	Forbidden     = Error{Code: "forbidden", Status: http.StatusForbidden}
	Throttled     = Error{Code: "throttled", Status: http.StatusTooManyRequests}
	MissingAuth   = Error{Code: "missing_auth", Status: http.StatusUnauthorized}
	Unsupported   = Error{Code: "unsupported", Status: http.StatusUnprocessableEntity}
	InvalidInput  = Error{Code: "invalid_input", Status: http.StatusBadRequest}
	InternalIssue = Error{Code: "internal_issue", Status: http.StatusInternalServerError}
)

// E converts any given error to the Error type. Unknown are converted
// to ErrInternal.
func E(err error) Error {
	if e, ok := err.(Error); ok {
		return e
	}
	return InternalIssue.CausedBy(err)
}

// OneOf checks (in-order) if err is one of the given errors.
func OneOf(err error, others []error) bool {
	for _, other := range others {
		if Is(err, other) {
			return true
		}
	}
	return false
}
