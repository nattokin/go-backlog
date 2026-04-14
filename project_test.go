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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.All(ctx)
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Create(ctx, "TEST", "test")
				require.NoError(t, err)
				assert.Equal(t, "TEST", got.ProjectKey)
			},
		},
		"Update": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Update(ctx, "TEST")
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
