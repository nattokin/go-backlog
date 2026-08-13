// Package core provides the Error/Errors vocabulary used
// across the domain layer, and offset-based pagination helpers.
package validation

import "strings"

// Error represents an argument validation error.
type Error struct {
	target  string
	message string
}

func NewError(target, message string) *Error {
	return &Error{
		target:  target,
		message: message,
	}
}

func (e *Error) Target() string  { return e.target }
func (e *Error) Message() string { return e.message }
func (e *Error) Error() string   { return e.message }

// Errors is a collection of Error values returned when
// multiple options fail validation simultaneously.
type Errors []*Error

func (es Errors) Error() string {
	msgs := make([]string, len(es))
	for i, e := range es {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "\n")
}
