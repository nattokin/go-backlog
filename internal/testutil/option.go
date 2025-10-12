// Package testutil provides shared test helpers used across backlog package tests.
// It includes conversion utilities for RequestOption interfaces and
// other reusable test-level helpers.
//
// Note: This package is located under `internal/` to restrict its usage to tests only.
package testutil

import (
	"testing"

	"github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/require"
)

//
// ──────────────────────────────────────────────────────────────
//  Option conversion helpers
// ──────────────────────────────────────────────────────────────
//

// ToFormOptions converts a slice of RequestOption interfaces into
// a slice of *FormOption, asserting the correct type at runtime.
//
// It fails the test immediately if any element is not a *FormOption.
//
// Example:
//
//	formOpts := testutil.ToFormOptions(t, []backlog.RequestOption{
//		formService.WithName("test"),
//	})
//
//	form := backlog.NewFormParams()
//	err := formService.ApplyOptions(form, formOpts...)
//	require.NoError(t, err)
func ToFormOptions(t *testing.T, opts []backlog.RequestOption) []*backlog.FormOption {
	t.Helper()

	formOpts := make([]*backlog.FormOption, 0, len(opts))
	for _, opt := range opts {
		fopt, ok := opt.(*backlog.FormOption)
		require.Truef(t, ok, "expected *FormOption, got %T", opt)
		formOpts = append(formOpts, fopt)
	}
	return formOpts
}

// ToQueryOptions converts a slice of RequestOption interfaces into
// a slice of *QueryOption, asserting the correct type at runtime.
//
// It fails the test immediately if any element is not a *QueryOption.
//
// Example:
//
//	queryOpts := testutil.ToQueryOptions(t, []backlog.RequestOption{
//		queryService.WithAll(true),
//	})
//
//	query := backlog.NewQueryParams()
//	err := queryService.ApplyOptions(query, queryOpts...)
//	require.NoError(t, err)
func ToQueryOptions(t *testing.T, opts []backlog.RequestOption) []*backlog.QueryOption {
	t.Helper()

	queryOpts := make([]*backlog.QueryOption, 0, len(opts))
	for _, opt := range opts {
		qopt, ok := opt.(*backlog.QueryOption)
		require.Truef(t, ok, "expected *QueryOption, got %T", opt)
		queryOpts = append(queryOpts, qopt)
	}
	return queryOpts
}
