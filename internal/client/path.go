package client

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/nattokin/go-backlog/internal/validation"
)

// PathJoin joins segments with "/" to build an endpoint path.
//
// Every segment, fixed or caller-supplied, must be a single path segment:
// it must not be empty, "." or "..", and must not contain "/", "\" or
// control characters. The result is never cleaned or resolved, so each
// segment stays in its own position and a value cannot redirect the
// request to another endpoint.
//
// Segments are not escaped here. Assign the result to url.URL.Path, which
// escapes it on serialization.
func PathJoin(segments ...string) (string, error) {
	for _, s := range segments {
		if ve := validateSegment(s); ve != nil {
			return "", ve
		}
	}
	return strings.Join(segments, "/"), nil
}

func validateSegment(s string) *validation.Error {
	switch {
	case s == "":
		return validation.NewError("path", "invalid path segment: must not be empty")
	case s == "." || s == "..":
		return validation.NewError("path", fmt.Sprintf("invalid path segment %q: must not be a dot segment", s))
	case strings.ContainsAny(s, `/\`):
		return validation.NewError("path", fmt.Sprintf("invalid path segment %q: must not contain path separators", s))
	case strings.IndexFunc(s, unicode.IsControl) >= 0:
		return validation.NewError("path", fmt.Sprintf("invalid path segment %q: must not contain control characters", s))
	}
	return nil
}
