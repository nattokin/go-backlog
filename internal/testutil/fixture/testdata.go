package fixture

import (
	"time"

	"github.com/nattokin/go-backlog"
)

// InvalidJSON is a malformed JSON string used to test JSON parse error handling.
const InvalidJSON = `
{invalid}
`

// mustTime parses an RFC3339 time string and panics on error.
// Used only for initializing package-level test fixture variables.
func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic("fixture: failed to parse time: " + err.Error())
	}
	return t
}

// mustTimestamp wraps mustTime as a backlog.Timestamp.
// Used only for initializing package-level test fixture variables.
func mustTimestamp(s string) backlog.Timestamp {
	return backlog.Timestamp{Time: mustTime(s)}
}

// ptrFloat64 returns a pointer to v.
// Used only for initializing package-level test fixture variables.
func ptrFloat64(v float64) *float64 {
	return &v
}
