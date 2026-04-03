package backlog

import (
	"fmt"
	"strings"
)

// Error represents one of Backlog API response errors.
type Error struct {
	// Message is the detailed error message from the API.
	Message  string `json:"message,omitempty"`
	Code     int    `json:"code,omitempty"`
	MoreInfo string `json:"moreInfo,omitempty"`
}

// Error returns the API error message.
func (e *Error) Error() string {
	msg := fmt.Sprintf("Message:%s, Code:%d", e.Message, e.Code)

	if e.MoreInfo == "" {
		return msg
	}

	return msg + ", MoreInfo:" + e.MoreInfo
}

// InternalClientError represents client-side configuration or usage errors.
// It is distinct from API-level errors and indicates issues like missing token
// or malformed base URL.
type InternalClientError struct {
	msg string
}

func (e *InternalClientError) Error() string {
	return e.msg
}

func newInternalClientError(msg string) *InternalClientError {
	return &InternalClientError{msg: msg}
}

// APIResponseError represents Error Response of Backlog API.
type APIResponseError struct {
	StatusCode int      `json:"-"` // HTTP status code (4xx or 5xx)
	Errors     []*Error `json:"errors,omitempty"`
}

// Error returns all error messages in APIResponseError.
func (e *APIResponseError) Error() string {
	msgs := make([]string, len(e.Errors))

	for i, err := range e.Errors {
		msgs[i] = err.Error()
	}

	return fmt.Sprintf("Status Code:%d\n%s", e.StatusCode, strings.Join(msgs, "\n"))
}

// InvalidOptionError represents an error for an invalid option value.
type InvalidOptionError[T requestOptionType] struct {
	Invalid   T
	ValidList []T
}

func newInvalidOptionError[T requestOptionType](invalid T, validList []T) *InvalidOptionError[T] {
	return &InvalidOptionError[T]{
		Invalid:   invalid,
		ValidList: validList,
	}
}

func (e *InvalidOptionError[T]) Error() string {
	types := make([]string, len(e.ValidList))
	for k, v := range e.ValidList {
		types[k] = v.Value()
	}

	return fmt.Sprintf("invalid option:%s, allowed options:%s", e.Invalid.Value(), strings.Join(types, ","))
}

// ValidationError represents an argument validation error.
type ValidationError struct {
	Message string
}

func newValidationError(msg string) *ValidationError {
	return &ValidationError{
		Message: msg,
	}
}

func (e *ValidationError) Error() string {
	return e.Message
}
