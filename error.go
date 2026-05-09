package backlog

import (
	"fmt"

	"github.com/nattokin/go-backlog/internal/core"
)

// Error represents one of the individual error entries in a Backlog API response.
// It is used as an element of [APIResponseError.Errors] and is not an error itself.
type Error struct {
	// Message is the detailed error message from the API.
	Message  string
	Code     int
	MoreInfo string
}

// APIResponseError represents an error response from the Backlog API.
// Use [errors.As] to check whether a returned error is an *APIResponseError.
type APIResponseError struct {
	core *core.APIResponseError
}

// Error implements the error interface.
func (e *APIResponseError) Error() string { return e.core.Error() }

// StatusCode returns the HTTP status code of the error response.
func (e *APIResponseError) StatusCode() int { return e.core.StatusCode }

// Errors returns the individual error entries in the response.
func (e *APIResponseError) Errors() []*Error {
	out := make([]*Error, len(e.core.Errors))
	for i, ce := range e.core.Errors {
		out[i] = &Error{
			Message:  ce.Message,
			Code:     ce.Code,
			MoreInfo: ce.MoreInfo,
		}
	}
	return out
}

// InvalidOptionKeyError is returned when an option method is called with a key
// that is not valid for the target service method.
// Use [errors.As] to check whether a returned error is an *InvalidOptionKeyError.
type InvalidOptionKeyError struct {
	core *core.InvalidOptionKeyError
}

// Error implements the error interface.
func (e *InvalidOptionKeyError) Error() string { return e.core.Error() }

// InvalidKey returns the invalid option key that was provided.
func (e *InvalidOptionKeyError) InvalidKey() string { return e.core.Invalid }

// AllowKeys returns the list of allowed option keys.
func (e *InvalidOptionKeyError) AllowKeys() []string { return e.core.ValidList }

// ValidationError is returned when a required argument fails validation
// (e.g. an empty string where a non-empty value is required).
// Use [errors.As] to check whether a returned error is a *ValidationError.
type ValidationError struct {
	core *core.ValidationError
}

// Error implements the error interface.
func (e *ValidationError) Error() string { return e.core.Error() }

// InternalClientError represents client-side configuration or usage errors.
// It is distinct from API-level errors and indicates issues like a missing token
// or a malformed base URL.
// Use [errors.As] to check whether a returned error is an *InternalClientError.
type InternalClientError struct {
	core *core.InternalClientError
}

// Error implements the error interface.
func (e *InternalClientError) Error() string { return e.core.Error() }

// convertError converts an error returned from internal packages into the
// corresponding root-package error type. This prevents internal types from
// leaking into the public API surface.
//
// Only error types that core returns directly (without wrapping) are converted
// here; *Error is excluded because it is never returned standalone.
func convertError(err error) error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *core.APIResponseError:
		return &APIResponseError{core: e}
	case *core.InvalidOptionKeyError:
		return &InvalidOptionKeyError{core: e}
	case *core.ValidationError:
		return &ValidationError{core: e}
	case *core.InternalClientError:
		return &InternalClientError{core: e}
	default:
		return err
	}
}

// InvalidDateStringError is returned when a string passed to [NewDate] is not
// a valid date in "YYYY-MM-DD" format.
type InvalidDateStringError struct {
	Value string
}

func (e *InvalidDateStringError) Error() string {
	return fmt.Sprintf("backlog: invalid date string %q: expected \"YYYY-MM-DD\" format", e.Value)
}
