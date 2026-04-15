package backlog_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

func TestProjectService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"All": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects", req.URL.Path)
				assert.Equal(t, "true", req.URL.Query().Get("archived"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.All(ctx, c.Project.Option.WithArchived(true))
				require.NoError(t, err)
				assert.Len(t, got, 3)
			},
		},
		"One": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.One(ctx, "TEST")
				require.NoError(t, err)
				assert.Equal(t, "TEST", got.ProjectKey)
			},
		},
		"Create": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/projects", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "TEST", req.PostForm.Get("key"))
				assert.Equal(t, "test", req.PostForm.Get("name"))
				assert.Equal(t, "true", req.PostForm.Get("chartEnabled"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Create(ctx, "TEST", "test", c.Project.Option.WithChartEnabled(true))
				require.NoError(t, err)
				assert.Equal(t, "TEST", got.ProjectKey)
			},
		},
		"Update": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "new-name", req.PostForm.Get("name"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Update(ctx, "TEST", c.Project.Option.WithName("new-name"))
				require.NoError(t, err)
				assert.Equal(t, "TEST", got.ProjectKey)
			},
		},
		"Delete": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Delete(ctx, "TEST")
				require.NoError(t, err)
				assert.Equal(t, "TEST", got.ProjectKey)
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)
			tc.call(t, c)
		})
	}
}

func TestProjectActivityService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/activities", req.URL.Path)
				assert.Equal(t, "10", req.URL.Query().Get("count"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Activity.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Activity.List(ctx, "TEST", c.Project.Activity.Option.WithCount(10))
				require.NoError(t, err)
				assert.Len(t, got, 1)
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)
			tc.call(t, c)
		})
	}
}

func TestProjectOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	s := c.Project.Option

	cases := map[string]struct {
		option  backlog.RequestOption
		wantKey string
	}{
		"WithAll": {
			option:  s.WithAll(true),
			wantKey: core.ParamAll.Value(),
		},
		"WithArchived": {
			option:  s.WithArchived(true),
			wantKey: core.ParamArchived.Value(),
		},
		"WithChartEnabled": {
			option:  s.WithChartEnabled(true),
			wantKey: core.ParamChartEnabled.Value(),
		},
		"WithKey": {
			option:  s.WithKey("TEST"),
			wantKey: core.ParamKey.Value(),
		},
		"WithName": {
			option:  s.WithName("test"),
			wantKey: core.ParamName.Value(),
		},
		"WithProjectLeaderCanEditProjectLeader": {
			option:  s.WithProjectLeaderCanEditProjectLeader(true),
			wantKey: core.ParamProjectLeaderCanEditProjectLeader.Value(),
		},
		"WithSubtaskingEnabled": {
			option:  s.WithSubtaskingEnabled(true),
			wantKey: core.ParamSubtaskingEnabled.Value(),
		},
		"WithTextFormattingRule": {
			option:  s.WithTextFormattingRule(model.FormatBacklog),
			wantKey: core.ParamTextFormattingRule.Value(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}
