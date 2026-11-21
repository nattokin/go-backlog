package backlog

import (
	"time"
)

// Timestamp represents a datetime value returned by the Backlog API.
// All exported methods of [time.Time] can be called on Timestamp.
// The zero Timestamp means that the field has no value.
type Timestamp struct {
	time.Time
}
