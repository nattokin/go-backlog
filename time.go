package backlog

import (
	"time"
)

// Timestamp represents a datetime value returned by the Backlog API.
// The Backlog API returns datetime fields as RFC3339 strings (e.g. "2022-09-01T06:35:39Z").
// All exported methods of [time.Time] can be called on Timestamp.
type Timestamp struct {
	time.Time
}
