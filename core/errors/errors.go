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
	NotFound = Error{
		Code:    "not_found",
		Status:  http.StatusNotFound,
		Message: "Resource not found",
	}

	Conflict = Error{
		Code:    "conflict",
		Status:  http.StatusConflict,
		Message: "A conflicting state exists",
	}

	Forbidden = Error{
		Code:    "forbidden",
		Status:  http.StatusForbidden,
		Message: "You are not authorized",
	}
	Throttled = Error{
		Code:    "throttled",
		Status:  http.StatusTooManyRequests,
		Message: "You are doing way too much",
	}

	MissingAuth = Error{
		Code:    "missing_auth",
		Status:  http.StatusUnauthorized,
		Message: "You are not authenticated",
	}

	Unsupported = Error{
		Code:    "unsupported",
		Status:  http.StatusUnprocessableEntity,
		Message: "Requested action is not supported",
	}

	InvalidInput = Error{
		Code:    "invalid_input",
		Status:  http.StatusBadRequest,
		Message: "Your request is not valid",
	}

	InternalIssue = Error{
		Code:    "internal_issue",
		Status:  http.StatusInternalServerError,
		Message: "Oops, something went wrong",
	}
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
