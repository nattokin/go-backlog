package backlog

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
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

// --- Option Service Helpers ---

// newOptionService returns a test instance of OptionService.
func newOptionService() *OptionService {
	return &OptionService{}
}

// newActivityOptionService returns a test instance of ActivityOptionService.
func newActivityOptionService() *ActivityOptionService {
	return &ActivityOptionService{
		registry: newOptionService(),
	}
}

// newProjectOptionService returns a test instance of ProjectOptionService.
func newProjectOptionService() *ProjectOptionService {
	return &ProjectOptionService{
		registry: newOptionService(),
	}
}

// newUserOptionService returns a test instance of UserOptionService.
func newUserOptionService() *UserOptionService {
	return &UserOptionService{
		registry: newOptionService(),
	}
}

// newWikiOptionService returns a test instance of WikiOptionService.
func newWikiOptionService() *WikiOptionService {
	return &WikiOptionService{
		registry: newOptionService(),
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
func newUnexpectedGetFn(t *testing.T) func(spath string, query url.Values) (*http.Response, error) {
	t.Helper()
	return func(spath string, query url.Values) (*http.Response, error) {
		t.Helper()
		t.Error("Get must not be called")
		return nil, errors.New("unexpected call")
	}
}

// newUnexpectedPostFn returns a mock function for http POST that fails if called.
func newUnexpectedPostFn(t *testing.T) func(spath string, form url.Values) (*http.Response, error) {
	t.Helper()
	return func(spath string, form url.Values) (*http.Response, error) {
		t.Helper()
		t.Error("Post must not be called")
		return nil, errors.New("unexpected call")
	}
}

// newUnexpectedPatchFn returns a mock function for http PATCH that fails if called.
func newUnexpectedPatchFn(t *testing.T) func(spath string, form url.Values) (*http.Response, error) {
	t.Helper()
	return func(spath string, form url.Values) (*http.Response, error) {
		t.Helper()
		t.Error("Patch must not be called")
		return nil, errors.New("unexpected call")
	}
}

// newUnexpectedDeleteFn returns a mock function for http DELETE that fails if called.
func newUnexpectedDeleteFn(t *testing.T) func(spath string, form url.Values) (*http.Response, error) {
	t.Helper()
	return func(spath string, form url.Values) (*http.Response, error) {
		t.Helper()
		t.Error("Delete must not be called")
		return nil, errors.New("unexpected call")
	}
}

// newUnexpectedUploadFn returns a mock function for http Upload that fails if called.
func newUnexpectedUploadFn(t *testing.T) func(spath, fileName string, r io.Reader) (*http.Response, error) {
	t.Helper()
	return func(spath, fileName string, r io.Reader) (*http.Response, error) {
		t.Helper()
		t.Error("Upload must not be called")
		return nil, errors.New("unexpected call")
	}
}
