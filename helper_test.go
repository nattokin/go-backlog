package backlog

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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

// --- WikiService ------------------------------------------------------------

// newWikiService returns a test instance of WikiService.
func newWikiService() *WikiService {
	return &WikiService{
		method:     newClientMethodMock(),
		Attachment: newWikiAttachmentService(),
		Option:     newWikiOptionService(),
	}
}

// newWikiAttachmentService returns a test instance of WikiAttachmentService.
func newWikiAttachmentService() *WikiAttachmentService {
	return &WikiAttachmentService{
		method: newClientMethodMock(),
	}
}

// --- ProjectService ------------------------------------------------------------

// newProjectService returns a test instance of ProjectService.
func newProjectService() *ProjectService {
	return &ProjectService{
		method:   newClientMethodMock(),
		Activity: newProjectActivityService(),
		User:     newProjectUserService(),
		Option:   newProjectOptionService(),
	}
}

// newProjectActivityService returns a test instance of ProjectActivityService.
func newProjectActivityService() *ProjectActivityService {
	return &ProjectActivityService{
		method: newClientMethodMock(),
		Option: newActivityOptionService(),
	}
}

// newProjectUserService returns a test instance of ProjectUserService.
func newProjectUserService() *ProjectUserService {
	return &ProjectUserService{
		method: newClientMethodMock(),
	}
}

// --- UserService ------------------------------------------------------------

// newUserService returns a test instance of UserService.
func newUserService() *UserService {
	return &UserService{
		method:   newClientMethodMock(),
		Activity: newUserActivityService(),
		Option:   newUserOptionService(),
	}
}

// newUserActivityService returns a test instance of UserActivityService.
func newUserActivityService() *UserActivityService {
	return &UserActivityService{
		method: newClientMethodMock(),
		Option: newActivityOptionService(),
	}
}

// --- SpaceService ------------------------------------------------------------

// newSpaceService returns a test instance of SpaceService.
func newSpaceService() *SpaceService {
	return &SpaceService{
		method:     newClientMethodMock(),
		Activity:   newSpaceActivityService(),
		Attachment: newSpaceAttachmentService(),
	}
}

// newSpaceActivityService returns a test instance of SpaceActivityService.
func newSpaceActivityService() *SpaceActivityService {
	return &SpaceActivityService{
		method: newClientMethodMock(),
		Option: newActivityOptionService(),
	}
}

// newSpaceAttachmentService returns a test instance of SpaceAttachmentService.
func newSpaceAttachmentService() *SpaceAttachmentService {
	return &SpaceAttachmentService{
		method: newClientMethodMock(),
	}
}

// --- IssueService ------------------------------------------------------------

// newIssueService returns a test instance of IssueService.
func newIssueService() *IssueService {
	return &IssueService{
		method:     newClientMethodMock(),
		Attachment: newIssueAttachmentService(),
	}
}

// newIssueAttachmentService returns a test instance of IssueAttachmentService.
func newIssueAttachmentService() *IssueAttachmentService {
	return &IssueAttachmentService{
		method: newClientMethodMock(),
	}
}

// --- PullRequestService ------------------------------------------------------------

// newPullRequestService returns a test instance of PullRequestService.
func newPullRequestService() *PullRequestService {
	return &PullRequestService{
		method:     newClientMethodMock(),
		Attachment: newPullRequestAttachmentService(),
	}
}

// newPullRequestAttachmentService returns a test instance of PullRequestAttachmentService.
func newPullRequestAttachmentService() *PullRequestAttachmentService {
	return &PullRequestAttachmentService{
		method: newClientMethodMock(),
	}
}

// newUnexpectedGetFn returns a mock function for http GET that fails if called.
func newUnexpectedGetFn(t *testing.T, reason string) func(spath string, query *QueryParams) (*http.Response, error) {
	t.Helper()
	return func(spath string, query *QueryParams) (*http.Response, error) {
		t.Helper()
		t.Errorf("Get must not be called when %s", reason)
		return nil, errors.New("unexpected call")
	}
}

// newMockGetFn returns a mock function for http GET that returns the given response.
func newMockGetFn(t *testing.T, wantPath string, res *http.Response) func(spath string, query *QueryParams) (*http.Response, error) {
	t.Helper()
	return func(spath string, query *QueryParams) (*http.Response, error) {
		t.Helper()
		assert.Equal(t, wantPath, spath)
		assert.Nil(t, query)
		return res, nil
	}
}

// newUnexpectedPostFn returns a mock function for http POST that fails if called.
func newUnexpectedPostFn(t *testing.T, reason string) func(spath string, form *FormParams) (*http.Response, error) {
	t.Helper()
	return func(spath string, form *FormParams) (*http.Response, error) {
		t.Helper()
		t.Errorf("Post must not be called when %s", reason)
		return nil, errors.New("unexpected call")
	}
}

// newMockPostFn returns a mock function for http POST that returns the given response.
func newMockPostFn(t *testing.T, wantPath string, res *http.Response) func(spath string, form *FormParams) (*http.Response, error) {
	t.Helper()
	return func(spath string, form *FormParams) (*http.Response, error) {
		t.Helper()
		assert.Equal(t, wantPath, spath)
		return res, nil
	}
}

// newUnexpectedPatchFn returns a mock function for http PATCH that fails if called.
func newUnexpectedPatchFn(t *testing.T, reason string) func(spath string, form *FormParams) (*http.Response, error) {
	t.Helper()
	return func(spath string, form *FormParams) (*http.Response, error) {
		t.Helper()
		t.Errorf("Patch must not be called when %s", reason)
		return nil, errors.New("unexpected call")
	}
}

// newMockPatchFn returns a mock function for http PATCH that returns the given response.
func newMockPatchFn(t *testing.T, wantPath string, res *http.Response) func(spath string, form *FormParams) (*http.Response, error) {
	t.Helper()
	return func(spath string, form *FormParams) (*http.Response, error) {
		t.Helper()
		assert.Equal(t, wantPath, spath)
		return res, nil
	}
}

// newUnexpectedDeleteFn returns a mock function for http DELETE that fails if called.
func newUnexpectedDeleteFn(t *testing.T, reason string) func(spath string, form *FormParams) (*http.Response, error) {
	t.Helper()
	return func(spath string, form *FormParams) (*http.Response, error) {
		t.Helper()
		t.Errorf("Delete must not be called when %s", reason)
		return nil, errors.New("unexpected call")
	}
}

// newMockDeleteFn returns a mock function for http DELETE that returns the given response.
func newMockDeleteFn(t *testing.T, wantPath string, res *http.Response) func(spath string, form *FormParams) (*http.Response, error) {
	t.Helper()
	return func(spath string, form *FormParams) (*http.Response, error) {
		t.Helper()
		assert.Equal(t, wantPath, spath)
		return res, nil
	}
}

// newUnexpectedUploadFn returns a mock function for http Upload that fails if called.
func newUnexpectedUploadFn(t *testing.T, reason string) func(spath, fileName string, r io.Reader) (*http.Response, error) {
	t.Helper()
	return func(spath, fileName string, r io.Reader) (*http.Response, error) {
		t.Helper()
		t.Errorf("Upload must not be called when %s", reason)
		return nil, errors.New("unexpected call")
	}
}

// newMockUploadFn returns a mock function for http Upload that returns the given response.
func newMockUploadFn(t *testing.T, wantPath, wantFileName string, body []byte) func(spath, fileName string, r io.Reader) (*http.Response, error) {
	t.Helper()
	return func(spath, fileName string, r io.Reader) (*http.Response, error) {
		t.Helper()
		assert.Equal(t, wantPath, spath)
		assert.Equal(t, wantFileName, fileName)
		data, _ := io.ReadAll(r)
		assert.NotNil(t, data)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}, nil
	}
}
