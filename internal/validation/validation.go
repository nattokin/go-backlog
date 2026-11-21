// Package validation provides the Error/Errors vocabulary used
// across the domain layer.
package validation

import "strings"

// Error represents an argument validation error.
type Error struct {
	target  string
	message string
}

// NewError returns an *Error for target, the name of the invalid argument,
// with the given message.
func NewError(target, message string) *Error {
	return &Error{
		target:  target,
		message: message,
	}
}

// Target returns the name of the argument that failed validation.
func (e *Error) Target() string { return e.target }

// Message returns the validation error message.
func (e *Error) Message() string { return e.message }

// Error implements the error interface.
func (e *Error) Error() string { return e.message }

// Errors is a collection of Error values returned when
// multiple options fail validation simultaneously.
type Errors []*Error

// Error implements the error interface. It joins the messages of all errors
// with newlines.
func (es Errors) Error() string {
	msgs := make([]string, len(es))
	for i, e := range es {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "\n")
}
