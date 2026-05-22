package core

import (
	"fmt"
	"strings"
)

// Error represents one of Backlog API response errors.
type Error struct {
	Message  string `json:"message,omitempty"`
	Code     int    `json:"code,omitempty"`
	MoreInfo string `json:"moreInfo,omitempty"`
}

func (e *Error) Error() string {
	msg := fmt.Sprintf("Message:%s, Code:%d", e.Message, e.Code)

	if e.MoreInfo == "" {
		return msg
	}

	return msg + ", MoreInfo:" + e.MoreInfo
}

// APIResponseError represents Error Response of Backlog API.
type APIResponseError struct {
	StatusCode int      `json:"-"` // HTTP status code (4xx or 5xx)
	Errors     []*Error `json:"errors,omitempty"`
}

func (e *APIResponseError) Error() string {
	msgs := make([]string, len(e.Errors))

	for i, err := range e.Errors {
		msgs[i] = err.Error()
	}

	return fmt.Sprintf("Status Code:%d\n%s", e.StatusCode, strings.Join(msgs, "\n"))
}

// InvalidOptionKeyError represents an error for an invalid option key.
type InvalidOptionKeyError struct {
	Invalid   string
	ValidList []string
}

func NewInvalidOptionKeyError(invalid string, validList []APIParamOptionType) *InvalidOptionKeyError {
	validKeys := []string{}
	for _, v := range validList {
		validKeys = append(validKeys, v.Value())
	}

	return &InvalidOptionKeyError{
		Invalid:   invalid,
		ValidList: validKeys,
	}
}

func (e *InvalidOptionKeyError) Error() string {
	return fmt.Sprintf("invalid option key:%s, allowed option keys:%s", e.Invalid, strings.Join(e.ValidList, ","))
}

// InvalidOptionError represents an error for an invalid option, such as a nil
// option or a Check() implementation that returned a nil ValidationResult.
type InvalidOptionError struct {
	message string
}

func NewInvalidOptionError(msg string) *InvalidOptionError {
	return &InvalidOptionError{message: msg}
}

func (e *InvalidOptionError) Error() string {
	return e.message
}

// ValidationResult is the return type of RequestOption.Check().
// Implementations must return a non-nil ValidationResult; returning nil is
// treated as a programming error and causes ApplyOptions to return an
// InvalidOptionError.
type ValidationResult interface {
	Valid() bool
	Target() string
	Message() string
}

// ValidationError represents an argument validation error.
// It implements ValidationResult.
type ValidationError struct {
	target  string
	message string
}

func NewValidationError(target, msg string) *ValidationError {
	return &ValidationError{
		target:  target,
		message: msg,
	}
}

func (e *ValidationError) Valid() bool      { return false }
func (e *ValidationError) Target() string   { return e.target }
func (e *ValidationError) Message() string  { return e.message }
func (e *ValidationError) Error() string    { return e.message }

// ValidationErrors is a collection of ValidationError values returned when
// multiple options fail validation simultaneously.
type ValidationErrors []*ValidationError

func (es ValidationErrors) Error() string {
	msgs := make([]string, len(es))
	for i, e := range es {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "\n")
}

// InternalClientError represents client-side configuration or usage errors.
// It is distinct from API-level errors and indicates issues like missing Token
// or malformed base URL.
type InternalClientError struct {
	msg string
}

func (e *InternalClientError) Error() string {
	return e.msg
}

func NewInternalClientError(msg string) *InternalClientError {
	return &InternalClientError{msg: msg}
}
