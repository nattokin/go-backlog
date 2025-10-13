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

// --- Option Service Helpers ---

// newQueryOptionService returns a test instance of QueryOptionService.
func newQueryOptionService() *QueryOptionService {
	return &QueryOptionService{}
}

// newFormOptionService returns a test instance of FormOptionService.
func newFormOptionService() *FormOptionService {
	return &FormOptionService{}
}

// newActivityOptionService returns a test instance of ActivityOptionService.
func newActivityOptionService() *ActivityOptionService {
	return &ActivityOptionService{
		support: &optionSupport{
			query: newQueryOptionService(),
			form:  newFormOptionService(),
		},
	}
}

// newProjectOptionService returns a test instance of ProjectOptionService.
func newProjectOptionService() *ProjectOptionService {
	return &ProjectOptionService{
		support: &optionSupport{
			query: newQueryOptionService(),
			form:  newFormOptionService(),
		},
	}
}

// newUserOptionService returns a test instance of UserOptionService.
func newUserOptionService() *UserOptionService {
	return &UserOptionService{
		support: &optionSupport{
			query: newQueryOptionService(),
			form:  newFormOptionService(),
		},
	}
}

// newWikiOptionService returns a test instance of WikiOptionService.
func newWikiOptionService() *WikiOptionService {
	return &WikiOptionService{
		support: &optionSupport{
			query: newQueryOptionService(),
			form:  newFormOptionService(),
		},
	}
}

// newWikiService returns a test instance of WikiService.
func newWikiService() *WikiService {
	return &WikiService{
		method:     newClientMethodMock(),
		Attachment: ExportNewWikiAttachmentService(),
		Option:     newWikiOptionService(),
	}
}

// newProjectService returns a test instance of ProjectService.
func newProjectService() *ProjectService {
	return &ProjectService{
		method:   newClientMethodMock(),
		Activity: ExportNewProjectActivityService(),
		User:     ExportNewProjectUserService(),
		Option:   newProjectOptionService(),
	}
}
