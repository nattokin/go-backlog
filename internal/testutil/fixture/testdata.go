package fixture

import "time"

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
