// Package core provides the ValidationError/ValidationErrors vocabulary used
// across the domain layer, and offset-based pagination helpers.
package core

import "strings"

// ValidationError represents an argument validation error.
type ValidationError struct {
	target  string
	message string
}

func NewValidationError(target, message string) *ValidationError {
	return &ValidationError{
		target:  target,
		message: message,
	}
}

func (e *ValidationError) Target() string  { return e.target }
func (e *ValidationError) Message() string { return e.message }
func (e *ValidationError) Error() string   { return e.message }

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
