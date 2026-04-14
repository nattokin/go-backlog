package backlog_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

func TestProjectUserService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"All": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/users", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.All(ctx, "TEST", false)
				require.NoError(t, err)
				assert.Len(t, got, 4)
			},
		},
		"Add": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/users", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "1", req.PostForm.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.Add(ctx, "TEST", 1)
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
			},
		},
		"Delete": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/users", req.URL.Path)
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				form, err := url.ParseQuery(string(body))
				require.NoError(t, err)
				assert.Equal(t, "1", form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.Delete(ctx, "TEST", 1)
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
			},
		},
		"AddAdmin": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/administrators", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "1", req.PostForm.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.AddAdmin(ctx, "TEST", 1)
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
			},
		},
		"AdminAll": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/administrators", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.AdminAll(ctx, "TEST")
				require.NoError(t, err)
				assert.Len(t, got, 4)
			},
		},
		"DeleteAdmin": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/administrators", req.URL.Path)
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				form, err := url.ParseQuery(string(body))
				require.NoError(t, err)
				assert.Equal(t, "1", form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.DeleteAdmin(ctx, "TEST", 1)
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
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
