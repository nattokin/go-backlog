package backlog

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/nattokin/go-backlog/internal/core"
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
func newOptionService() *core.OptionService {
	return &core.OptionService{}
}

// newActivityOptionService returns a test instance of ActivityOptionService.
func newActivityOptionService() *ActivityOptionService {
	return &ActivityOptionService{
		base: newOptionService(),
	}
}

// newProjectOptionService returns a test instance of ProjectOptionService.
func newProjectOptionService() *ProjectOptionService {
	return &ProjectOptionService{
		base: newOptionService(),
	}
}

// newUserOptionService returns a test instance of UserOptionService.
func newUserOptionService() *UserOptionService {
	return &UserOptionService{
		base: newOptionService(),
	}
}

// --- ProjectService ------------------------------------------------------------

// newProjectService returns a test instance of ProjectService.
func newProjectService() *ProjectService {
	return &ProjectService{
		method:   newClientMethod(),
		Activity: newProjectActivityService(),
		User:     newProjectUserService(),
		Option:   newProjectOptionService(),
	}
}

// newProjectActivityService returns a test instance of ProjectActivityService.
func newProjectActivityService() *ProjectActivityService {
	return &ProjectActivityService{
		method: newClientMethod(),
		Option: newActivityOptionService(),
	}
}

// newProjectUserService returns a test instance of ProjectUserService.
func newProjectUserService() *ProjectUserService {
	return &ProjectUserService{
		method: newClientMethod(),
	}
}

// --- UserService ------------------------------------------------------------

// newUserService returns a test instance of UserService.
func newUserService() *UserService {
	return &UserService{
		method:   newClientMethod(),
		Activity: newUserActivityService(),
		Option:   newUserOptionService(),
	}
}

// newUserActivityService returns a test instance of UserActivityService.
func newUserActivityService() *UserActivityService {
	return &UserActivityService{
		method: newClientMethod(),
		Option: newActivityOptionService(),
	}
}

// --- SpaceService ------------------------------------------------------------

// newSpaceService returns a test instance of SpaceService.
func newSpaceService() *SpaceService {
	return &SpaceService{
		method:     newClientMethod(),
		Activity:   newSpaceActivityService(),
		Attachment: newSpaceAttachmentService(),
	}
}

// newSpaceActivityService returns a test instance of SpaceActivityService.
func newSpaceActivityService() *SpaceActivityService {
	return &SpaceActivityService{
		method: newClientMethod(),
		Option: newActivityOptionService(),
	}
}

// newSpaceAttachmentService returns a test instance of SpaceAttachmentService.
func newSpaceAttachmentService() *SpaceAttachmentService {
	return &SpaceAttachmentService{
		method: newClientMethod(),
	}
}

// --- IssueService ------------------------------------------------------------

// newIssueService returns a test instance of IssueService.
func newIssueService() *IssueService {
	return &IssueService{
		method:     newClientMethod(),
		Attachment: newIssueAttachmentService(),
	}
}

// newIssueAttachmentService returns a test instance of IssueAttachmentService.
func newIssueAttachmentService() *IssueAttachmentService {
	return &IssueAttachmentService{
		method: newClientMethod(),
	}
}

// --- PullRequestService ------------------------------------------------------------

// newPullRequestService returns a test instance of PullRequestService.
func newPullRequestService() *PullRequestService {
	return &PullRequestService{
		method:     newClientMethod(),
		Attachment: newPullRequestAttachmentService(),
	}
}

// newPullRequestAttachmentService returns a test instance of PullRequestAttachmentService.
func newPullRequestAttachmentService() *PullRequestAttachmentService {
	return &PullRequestAttachmentService{
		method: newClientMethod(),
	}
}

// newClientMethod creates and returns a mock implementation of the `method` struct.
// Each API function (Get, Post, Patch, Delete) returns a default "not implemented" error.
func newClientMethod() *core.Method {
	return &core.Method{
		Get: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
		Post: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
		Patch: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
		Delete: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
		Upload: func(ctx context.Context, spath, fileName string, r io.Reader) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
	}
}
