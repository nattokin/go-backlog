package backlog

import (
	"fmt"
	"strings"

	"github.com/nattokin/go-backlog/internal/core"
)

// Error represents one of Backlog API response errors.
type Error struct {
	// Message is the detailed error message from the API.
	Message  string
	Code     int
	MoreInfo string
}

// Error returns the API error message.
func (e *Error) Error() string {
	msg := fmt.Sprintf("Message:%s, Code:%d", e.Message, e.Code)

	if e.MoreInfo == "" {
		return msg
	}

	return msg + ", MoreInfo:" + e.MoreInfo
}

// APIResponseError represents Error Response of Backlog API.
type APIResponseError struct {
	StatusCode int // HTTP status code (4xx or 5xx)
	Errors     []*Error
}

// Error returns all error messages in APIResponseError.
func (e *APIResponseError) Error() string {
	msgs := make([]string, len(e.Errors))

	for i, err := range e.Errors {
		msgs[i] = err.Error()
	}

	return fmt.Sprintf("Status Code:%d\n%s", e.StatusCode, strings.Join(msgs, "\n"))
}

// InvalidOptionKeyError represents an error for an invalid option value.
type InvalidOptionKeyError struct {
	Invalid   string
	ValidList []string
}

// Error returns the error message for an invalid option key.
func (e *InvalidOptionKeyError) Error() string {
	return fmt.Sprintf("invalid option key:%s, allowed option keys:%s", e.Invalid, strings.Join(e.ValidList, ","))
}

// ValidationError represents an argument validation error.
type ValidationError struct {
	message string
}

// Error returns the validation error message.
func (e *ValidationError) Error() string {
	return e.message
}

// InternalClientError represents client-side configuration or usage errors.
// It is distinct from API-level errors and indicates issues like missing Token
// or malformed base URL.
type InternalClientError struct {
	msg string
}

// Error returns the internal client error message.
func (e *InternalClientError) Error() string {
	return e.msg
}

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
		out := &APIResponseError{
			StatusCode: e.StatusCode,
			Errors:     make([]*Error, len(e.Errors)),
		}
		for i, ce := range e.Errors {
			out.Errors[i] = &Error{
				Message:  ce.Message,
				Code:     ce.Code,
				MoreInfo: ce.MoreInfo,
			}
		}
		return out
	case *core.InvalidOptionKeyError:
		return &InvalidOptionKeyError{
			Invalid:   e.Invalid,
			ValidList: e.ValidList,
		}
	case *core.ValidationError:
		return &ValidationError{message: e.Error()}
	case *core.InternalClientError:
		return &InternalClientError{msg: e.Error()}
	default:
		return err
	}
}
