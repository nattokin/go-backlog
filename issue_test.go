package backlog_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestIssueService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues", req.URL.Path)
				return mock.NewJSONResponse(fixture.Issue.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.List(ctx)
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 1, got[0].ID)
				assert.Equal(t, 2, got[1].ID)
			},
		},
		"List/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues", req.URL.Path)
				assert.Equal(t, "bug", req.URL.Query().Get("keyword"))
				assert.Equal(t, []string{"10", "20"}, req.URL.Query()["projectId[]"])
				return mock.NewJSONResponse(fixture.Issue.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.List(ctx,
					c.Issue.Option.WithKeyword("bug"),
					c.Issue.Option.WithProjectIDs([]int{10, 20}),
				)
				require.NoError(t, err)
				assert.Len(t, got, 2)
			},
		},
		"All/error": {
			doFunc: newInternalServerErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.List(ctx)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Count": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues/count", req.URL.Path)
				return mock.NewJSONResponse(`{"count":5}`), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Count(ctx)
				require.NoError(t, err)
				assert.Equal(t, 5, got)
			},
		},
		"Count/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues/count", req.URL.Path)
				assert.Equal(t, "bug", req.URL.Query().Get("keyword"))
				assert.Equal(t, []string{"10", "20"}, req.URL.Query()["projectId[]"])
				return mock.NewJSONResponse(`{"count":2}`), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Count(ctx,
					c.Issue.Option.WithKeyword("bug"),
					c.Issue.Option.WithProjectIDs([]int{10, 20}),
				)
				require.NoError(t, err)
				assert.Equal(t, 2, got)
			},
		},
		"Count/error": {
			doFunc: newInternalServerErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Count(ctx)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"One": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues/PRJ-1", req.URL.Path)
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.One(ctx, "PRJ-1")
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
				assert.Equal(t, "PRJ-1", got.IssueKey)
				assert.Equal(t, "First issue", got.Summary)
			},
		},
		"One/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.One(ctx, "PRJ-1")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Create": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/issues", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "10", req.PostForm.Get("projectId"))
				assert.Equal(t, "New issue", req.PostForm.Get("summary"))
				assert.Equal(t, "2", req.PostForm.Get("issueTypeId"))
				assert.Equal(t, "3", req.PostForm.Get("priorityId"))
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Create(ctx, 10, "New issue", 2, 3)
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
				assert.Equal(t, "First issue", got.Summary)
			},
		},
		"Create/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/issues", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "10", req.PostForm.Get("projectId"))
				assert.Equal(t, "New issue", req.PostForm.Get("summary"))
				assert.Equal(t, "2", req.PostForm.Get("issueTypeId"))
				assert.Equal(t, "3", req.PostForm.Get("priorityId"))
				assert.Equal(t, "details here", req.PostForm.Get("description"))
				assert.Equal(t, "5", req.PostForm.Get("assigneeId"))
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Create(ctx, 10, "New issue", 2, 3,
					c.Issue.Option.WithDescription("details here"),
					c.Issue.Option.WithAssigneeID(5),
				)
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
			},
		},
		"Create/error": {
			doFunc: newAuthErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Create(ctx, 10, "New issue", 2, 3)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Update": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/issues/PRJ-1", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "Updated summary", req.PostForm.Get("summary"))
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Update(ctx, "PRJ-1", c.Issue.Option.WithSummary("Updated summary"))
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
			},
		},
		"Update/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/issues/PRJ-1", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "Updated summary", req.PostForm.Get("summary"))
				assert.Equal(t, "2", req.PostForm.Get("statusId"))
				assert.Equal(t, "1", req.PostForm.Get("resolutionId"))
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Update(ctx, "PRJ-1",
					c.Issue.Option.WithSummary("Updated summary"),
					c.Issue.Option.WithStatusID(2),
					c.Issue.Option.WithResolutionID(1),
				)
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
			},
		},
		"Update/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Update(ctx, "PRJ-1", c.Issue.Option.WithSummary("Updated summary"))
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Delete": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/issues/PRJ-1", req.URL.Path)
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Delete(ctx, "PRJ-1")
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
				assert.Equal(t, "PRJ-1", got.IssueKey)
			},
		},
		"Delete/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Delete(ctx, "PRJ-1")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Participants": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues/PRJ-1/participants", req.URL.Path)
				return mock.NewJSONResponse(fixture.User.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Participants(ctx, "PRJ-1")
				require.NoError(t, err)
				assert.Len(t, got, 4)
				assert.Equal(t, 1, got[0].ID)
				assert.Equal(t, "admin", got[0].Name)
			},
		},
		"Participants/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Participants(ctx, "PRJ-1")
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

func TestIssueStarService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"Add": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/stars", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "1", req.FormValue("issueId"))
				return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Issue.Star.Add(ctx, 1)
				require.NoError(t, err)
			},
		},
		"Add/error": {
			doFunc: newAuthErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Issue.Star.Add(ctx, 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Remove": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/stars", req.URL.Path)
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				form, err := url.ParseQuery(string(body))
				require.NoError(t, err)
				assert.Equal(t, "42", form.Get("id"))
				return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Issue.Star.Remove(ctx, 42)
				require.NoError(t, err)
			},
		},
		"Remove/error": {
			doFunc: newAuthErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Issue.Star.Remove(ctx, 42)
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
