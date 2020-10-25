package backlog_test

import (
	"strconv"
	"testing"

	backlog "github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestProjectOptionService_WithKey(t *testing.T) {
	s := backlog.ExportNewProjectService(nil).Option

	cases := map[string]struct {
		key       string
		wantError bool
	}{
		"valid": {
			key:       "TEST",
			wantError: false,
		},
		"empty": {
			key:       "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithKey(tc.key)
			params := backlog.ExportNewRequestParams()
			err := option(params)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.key, params.Get("key"))
			}
		})
	}
}

func TestProjectOptionService_WithName(t *testing.T) {
	s := backlog.ExportNewProjectService(nil).Option

	cases := map[string]struct {
		name      string
		wantError bool
	}{
		"valid": {
			name:      "test",
			wantError: false,
		},
		"empty": {
			name:      "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithName(tc.name)
			params := backlog.ExportNewRequestParams()
			err := option(params)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.name, params.Get("name"))
			}
		})
	}
}

func TestProjectOptionService_WithChartEnabled(t *testing.T) {
	s := backlog.ExportNewProjectService(nil).Option

	cases := map[string]struct {
		enabeld bool
	}{
		"true": {
			enabeld: true,
		},
		"false": {
			enabeld: false,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithChartEnabled(tc.enabeld)
			params := backlog.ExportNewRequestParams()
			err := option(params)
			assert.Nil(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), params.Get("chartEnabled"))
		})
	}
}

func TestProjectOptionService_WithSubtaskingEnabled(t *testing.T) {
	s := backlog.ExportNewProjectService(nil).Option

	cases := map[string]struct {
		enabeld bool
	}{
		"true": {
			enabeld: true,
		},
		"false": {
			enabeld: false,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithSubtaskingEnabled(tc.enabeld)
			params := backlog.ExportNewRequestParams()
			err := option(params)
			assert.Nil(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), params.Get("subtaskingEnabled"))
		})
	}
}

func TestProjectOptionService_WithProjectLeaderCanEditProjectLeader(t *testing.T) {
	s := backlog.ExportNewProjectService(nil).Option

	cases := map[string]struct {
		enabeld bool
	}{
		"true": {
			enabeld: true,
		},
		"false": {
			enabeld: false,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithProjectLeaderCanEditProjectLeader(tc.enabeld)
			params := backlog.ExportNewRequestParams()
			err := option(params)
			assert.Nil(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), params.Get("projectLeaderCanEditProjectLeader"))
		})
	}
}

func TestProjectOptionService_WithTextFormattingRule(t *testing.T) {
	s := backlog.ExportNewProjectService(nil).Option
	// TODO
	cases := map[string]struct {
		format    string
		wantError bool
	}{
		"backlog": {
			format:    backlog.FormatBacklog,
			wantError: false,
		},
		"markdown": {
			format:    backlog.FormatMarkdown,
			wantError: false,
		},
		"invalid": {
			format:    "test",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithTextFormattingRule(tc.format)
			params := backlog.ExportNewRequestParams()
			err := option(params)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.format, params.Get("textFormattingRule"))
			}
		})
	}
}

func TestProjectOptionService_WithArchived(t *testing.T) {
	s := backlog.ExportNewProjectService(nil).Option

	cases := map[string]struct {
		archived bool
	}{
		"true": {
			archived: true,
		},
		"false": {
			archived: false,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithArchived(tc.archived)
			params := backlog.ExportNewRequestParams()
			err := option(params)
			assert.Nil(t, err)
			assert.Equal(t, strconv.FormatBool(tc.archived), params.Get("archived"))
		})
	}
}

func TestWikiOptionService_WithName(t *testing.T) {
	s := backlog.ExportNewWikiService(nil).Option

	cases := map[string]struct {
		name      string
		wantError bool
	}{
		"valid": {
			name:      "test",
			wantError: false,
		},
		"empty": {
			name:      "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithName(tc.name)
			params := backlog.ExportNewRequestParams()
			err := option(params)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.name, params.Get("name"))
			}
		})
	}
}

func TestWikiOptionService_WithContent(t *testing.T) {
	s := backlog.ExportNewWikiService(nil).Option

	cases := map[string]struct {
		content   string
		wantError bool
	}{
		"valid": {
			content:   "content",
			wantError: false,
		},
		"empty": {
			content:   "",
			wantError: true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithContent(tc.content)
			params := backlog.ExportNewRequestParams()
			err := option(params)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.content, params.Get("content"))
			}
		})
	}
}

func TestWikiOptionService_WithMailNotify(t *testing.T) {
	s := backlog.ExportNewWikiService(nil).Option

	cases := map[string]struct {
		enabeld bool
	}{
		"true": {
			enabeld: true,
		},
		"false": {
			enabeld: false,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			option := s.WithMailNotify(tc.enabeld)
			params := backlog.ExportNewRequestParams()
			err := option(params)
			assert.Nil(t, err)
			assert.Equal(t, strconv.FormatBool(tc.enabeld), params.Get("mailNotify"))
		})
	}
}
