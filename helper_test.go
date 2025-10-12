package backlog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

//
// ──────────────────────────────────────────────────────────────
//  Internal test helpers
// ──────────────────────────────────────────────────────────────
//
// This file provides shared helper functions for unit tests within
// the backlog package. These helpers are intended for internal tests
// that need access to unexported structs or methods.
//
// Note:
// Do not import external `testutil` here — this file is for
// package-local (internal) unit tests only.
//

// toFormOptions converts a slice of RequestOption interfaces into
// a slice of *FormOption safely.
//
// This helper is useful when testing form-related applyOptions
// or when verifying mixed RequestOption slices.
//
// Example:
//
//	opts := []RequestOption{
//		formService.WithName("example"),
//		formService.WithArchived(true),
//	}
//	formOpts := toFormOptions(t, opts)
//	err := formService.applyOptions(form, formOpts...)
//	require.NoError(t, err)
func toFormOptions(t *testing.T, opts []RequestOption) []*FormOption {
	t.Helper()

	formOpts := make([]*FormOption, 0, len(opts))
	for _, opt := range opts {
		fopt, ok := opt.(*FormOption)
		require.Truef(t, ok, "expected *FormOption, got %T", opt)
		formOpts = append(formOpts, fopt)
	}
	return formOpts
}

// toQueryOptions converts a slice of RequestOption interfaces into
// a slice of *QueryOption safely.
//
// It is mainly used for verifying query parameter behaviors within
// internal option service tests.
//
// Example:
//
//	opts := []RequestOption{
//		queryService.WithAll(true),
//	}
//	queryOpts := toQueryOptions(t, opts)
//	err := queryService.applyOptions(query, queryOpts...)
//	require.NoError(t, err)
func toQueryOptions(t *testing.T, opts []RequestOption) []*QueryOption {
	t.Helper()

	queryOpts := make([]*QueryOption, 0, len(opts))
	for _, opt := range opts {
		qopt, ok := opt.(*QueryOption)
		require.Truef(t, ok, "expected *QueryOption, got %T", opt)
		queryOpts = append(queryOpts, qopt)
	}
	return queryOpts
}
