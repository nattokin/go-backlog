// Package model defines the types that Backlog API JSON responses are decoded
// into. Their fields and json tags follow the API response, and custom
// decoding of a response, such as CustomFieldBound.UnmarshalJSON, lives here.
//
// These types are not exposed to users. The root backlog package defines the
// public types and converts these into them with xxxFromModel functions, such
// as issueFromModel. Types and behavior meant for users, such as
// backlog.Timestamp, belong in the root package, not here.
package model
