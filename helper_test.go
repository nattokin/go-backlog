package backlog

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/user"
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

// newProjectOptionService returns a test instance of ProjectOptionService.
func newProjectOptionService() *ProjectOptionService {
	return &ProjectOptionService{
		base: newOptionService(),
	}
}

// --- ProjectService ------------------------------------------------------------

// newProjectService returns a test instance of ProjectService.
func newProjectService() *ProjectService {
	return &ProjectService{
		method:   newClientMethod(),
		Activity: newProjectActivityService(),
		User:     user.NewProjectUserService(newClientMethod(), &core.OptionService{}),
		Option:   newProjectOptionService(),
	}
}

// newProjectActivityService returns a test instance of ProjectActivityService.
func newProjectActivityService() *ProjectActivityService {
	return &ProjectActivityService{
		method: newClientMethod(),
		Option: activity.NewActivityOptionService(&core.OptionService{}),
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
