package backlog_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestProjectIssueTypeService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/issueTypes", req.URL.Path)
				return mock.NewJSONResponse(fixture.IssueType.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.IssueType.List(ctx, "TEST")
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 1, got[0].ID)
				assert.Equal(t, "Bug", got[0].Name)
				assert.Equal(t, "#e30000", got[0].Color)
			},
		},
		"List/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.IssueType.List(ctx, "TEST")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Create": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/issueTypes", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "Bug", req.PostForm.Get("name"))
				assert.Equal(t, "#e30000", req.PostForm.Get("color"))
				return mock.NewJSONResponse(fixture.IssueType.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.IssueType.Create(ctx, "TEST", "Bug", "#e30000")
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
				assert.Equal(t, "Bug", got.Name)
				assert.Equal(t, "#e30000", got.Color)
			},
		},
		"Create/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "default summary", req.PostForm.Get("templateSummary"))
				assert.Equal(t, "default description", req.PostForm.Get("templateDescription"))
				return mock.NewJSONResponse(fixture.IssueType.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.IssueType.Create(ctx, "TEST", "Bug", "#e30000",
					c.Project.IssueType.Option.WithTemplateSummary("default summary"),
					c.Project.IssueType.Option.WithTemplateDescription("default description"),
				)
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
			},
		},
		"Create/error": {
			doFunc: newAuthErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.IssueType.Create(ctx, "TEST", "Bug", "#e30000")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Update": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/issueTypes/1", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "Bug Updated", req.PostForm.Get("name"))
				return mock.NewJSONResponse(fixture.IssueType.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.IssueType.Update(ctx, "TEST", 1,
					c.Project.IssueType.Option.WithName("Bug Updated"),
				)
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
			},
		},
		"Update/multiple-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "Bug Updated", req.PostForm.Get("name"))
				assert.Equal(t, "#990000", req.PostForm.Get("color"))
				assert.Equal(t, "new summary", req.PostForm.Get("templateSummary"))
				assert.Equal(t, "new description", req.PostForm.Get("templateDescription"))
				return mock.NewJSONResponse(fixture.IssueType.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.IssueType.Update(ctx, "TEST", 1,
					c.Project.IssueType.Option.WithName("Bug Updated"),
					c.Project.IssueType.Option.WithColor("#990000"),
					c.Project.IssueType.Option.WithTemplateSummary("new summary"),
					c.Project.IssueType.Option.WithTemplateDescription("new description"),
				)
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
			},
		},
		"Update/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.IssueType.Update(ctx, "TEST", 1,
					c.Project.IssueType.Option.WithName("Bug Updated"),
				)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Delete": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/issueTypes/1", req.URL.Path)
				return mock.NewJSONResponse(fixture.IssueType.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.IssueType.Delete(ctx, "TEST", 1, 2)
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
			},
		},
		"Delete/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.IssueType.Delete(ctx, "TEST", 1, 2)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
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
