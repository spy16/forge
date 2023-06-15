package errors

import (
	"fmt"
)

// Error represents an error value with all the relevant context.
type Error struct {
	Code      string         `json:"code"`
	Cause     error          `json:"cause,omitempty"`
	Status    int            `json:"status"`
	Attribs   map[string]any `json:"attribs,omitempty"`
	Message   string         `json:"message"`
	DebugHint string         `json:"debug_hint,omitempty"`
}

// Coded returns a clone of the original error with the given code.
func (err Error) Coded(code string, attribs ...map[string]any) Error {
	cl := err.clone(attribs...)
	cl.Code = code
	return cl
}

// CausedBy returns a clone of the error with `e` set as the cause.
func (err Error) CausedBy(e error) Error {
	cloned := err.clone()
	cloned.Cause = e
	return cloned
}

// Msgf returns a clone of the error with a user-friendly message.
func (err Error) Msgf(format string, args ...any) Error {
	cloned := err.clone()
	cloned.Message = fmt.Sprintf(format, args...)
	return cloned
}

// Hintf returns a clone of the error with a debug hint.
func (err Error) Hintf(format string, args ...any) Error {
	cloned := err.clone()
	cloned.DebugHint = fmt.Sprintf(format, args...)
	return cloned
}

// Error represents technical description of the error.
func (err Error) Error() string {
	msg := err.Code
	if err.Cause != nil {
		msg += fmt.Sprintf(": %v", err.Cause)
	}
	if err.DebugHint != "" {
		msg += fmt.Sprintf(": %s", err.DebugHint)
	}
	return msg
}

// Is checks if 'other' is of type Error and has the same code.
// See https://blog.golang.org/go1.13-errors.
func (err Error) Is(other error) bool {
	if oe, ok := other.(Error); ok {
		return oe.Status == err.Status
	}

	// unknown error types are considered as internal errors.
	return err.Status == InternalIssue.Status
}

func (err Error) clone(withAttribs ...map[string]any) Error {
	cl := err
	cl.Attribs = cloneMerge(err.Attribs, withAttribs...)
	return cl
}

func cloneMerge(base map[string]any, maps ...map[string]any) map[string]any {
	res := map[string]any{}
	for k, v := range base {
		res[k] = v
	}
	for _, m := range maps {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}
