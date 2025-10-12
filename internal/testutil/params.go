package testutil

import (
	"testing"

	"github.com/nattokin/go-backlog"
)

//
// ──────────────────────────────────────────────────────────────
//  Parameter initialization helpers
// ──────────────────────────────────────────────────────────────
//

// NewFormParams creates a new FormParams instance for testing.
//
// This helper is used to simplify form option tests, ensuring
// that each test case starts from a clean form state.
func NewFormParams(t *testing.T) *backlog.FormParams {
	t.Helper()
	return backlog.NewFormParams()
}

// NewQueryParams creates a new QueryParams instance for testing.
//
// This helper provides an initialized query parameter container
// to be used in QueryOptionService tests.
func NewQueryParams(t *testing.T) *backlog.QueryParams {
	t.Helper()
	return backlog.NewQueryParams()
}
