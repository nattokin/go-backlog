package backlog

import (
	"github.com/nattokin/go-backlog/internal/core"
)

// Error represents one of the individual error entries in a Backlog API response.
// It is a data structure used for decoding API responses and is not an error itself.
type Error struct {
	// Message is the detailed error message from the API.
	Message  string
	Code     int
	MoreInfo string
}

// APIResponseError represents an error response from the Backlog API.
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

// InvalidOptionKeyError represents an error for an invalid option key.
type InvalidOptionKeyError struct {
	core *core.InvalidOptionKeyError
}

// Error implements the error interface.
func (e *InvalidOptionKeyError) Error() string { return e.core.Error() }

// InvalidKey returns the invalid option key that was provided.
func (e *InvalidOptionKeyError) InvalidKey() string { return e.core.Invalid }

// AllowKeys returns the list of allowed option keys.
func (e *InvalidOptionKeyError) AllowKeys() []string { return e.core.ValidList }

// ValidationError represents an argument validation error.
type ValidationError struct {
	core *core.ValidationError
}

// Error implements the error interface.
func (e *ValidationError) Error() string { return e.core.Error() }

// InternalClientError represents client-side configuration or usage errors.
// It is distinct from API-level errors and indicates issues like missing Token
// or malformed base URL.
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
