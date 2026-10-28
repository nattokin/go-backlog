package client

import "strings"

// PathJoin joins segments with "/" to build an endpoint path.
//
// Unlike path.Join, the result is never cleaned: dot segments are not
// resolved and no slash is added or removed. Caller-supplied segments must
// be validated beforehand (see validate.ValidateIDOrKey); NewRequest rejects
// any dot segment that still reaches it.
//
// Segments are not escaped here. NewRequest assigns the result to
// url.URL.Path, which escapes it on serialization.
func PathJoin(segments ...string) string {
	return strings.Join(segments, "/")
}
