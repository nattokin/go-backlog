package backlog

import (
	"fmt"
	"time"
)

// Timestamp represents a datetime value returned by the Backlog API.
// The Backlog API returns datetime fields as RFC3339 strings (e.g. "2022-09-01T06:35:39Z").
// All exported methods of [time.Time] can be called on Timestamp.
type Timestamp struct {
	time.Time
}

// Date represents a date-only value returned by the Backlog API.
// The Backlog API returns date-only fields as "YYYY-MM-DD" strings (e.g. "2023-04-15").
// The internal representation is kept unexported so it can change in the future
// (e.g. switching to time.Time) without breaking the public API.
type Date struct {
	value string
}

// InvalidDateStringError is returned when a string passed to [NewDate] is not
// a valid date in "YYYY-MM-DD" format.
type InvalidDateStringError struct {
	Value string
}

func (e *InvalidDateStringError) Error() string {
	return fmt.Sprintf("backlog: invalid date string %q: expected \"YYYY-MM-DD\" format", e.Value)
}

// NewDate returns a Date with the given "YYYY-MM-DD" string.
// Returns [*InvalidDateStringError] if s is not a valid date.
func NewDate(s string) (Date, error) {
	if _, err := time.Parse("2006-01-02", s); err != nil {
		return Date{}, &InvalidDateStringError{Value: s}
	}
	return Date{value: s}, nil
}

// String returns the date as a "YYYY-MM-DD" string.
// Returns an empty string when the date is unset.
func (d Date) String() string {
	return d.value
}

// IsZero reports whether d represents the zero Date (unset).
func (d Date) IsZero() bool {
	return d.value == ""
}
