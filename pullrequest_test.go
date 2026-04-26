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
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

func TestPullRequestService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"All": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests", req.URL.Path)
				return newJSONResponse(fixture.PullRequest.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.All(ctx, "TEST", "repo")
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 2, got[0].ID)
				assert.Equal(t, 3, got[1].ID)
			},
		},
		"All/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests", req.URL.Path)
				assert.Equal(t, []string{"1", "2"}, req.URL.Query()["statusId[]"])
				assert.Equal(t, "10", req.URL.Query().Get("count"))
				return newJSONResponse(fixture.PullRequest.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.All(ctx, "TEST", "repo",
					c.PullRequest.Option.WithStatusIDs([]int{1, 2}),
					c.PullRequest.Option.WithCount(10),
				)
				require.NoError(t, err)
				assert.Len(t, got, 2)
			},
		},
		"All/error": {
			doFunc: newInternalServerErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.PullRequest.All(ctx, "TEST", "repo")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Count": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests/count", req.URL.Path)
				return newJSONResponse(`{"count":5}`), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Count(ctx, "TEST", "repo")
				require.NoError(t, err)
				assert.Equal(t, 5, got)
			},
		},
		"Count/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests/count", req.URL.Path)
				assert.Equal(t, []string{"1"}, req.URL.Query()["statusId[]"])
				return newJSONResponse(`{"count":3}`), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Count(ctx, "TEST", "repo",
					c.PullRequest.Option.WithStatusIDs([]int{1}),
				)
				require.NoError(t, err)
				assert.Equal(t, 3, got)
			},
		},
		"Count/error": {
			doFunc: newInternalServerErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.PullRequest.Count(ctx, "TEST", "repo")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"One": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests/1", req.URL.Path)
				return newJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.One(ctx, "TEST", "repo", 1)
				require.NoError(t, err)
				assert.Equal(t, 2, got.ID)
				assert.Equal(t, "test PR", got.Summary)
			},
		},
		"One/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.PullRequest.One(ctx, "TEST", "repo", 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Create": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "new PR", req.PostForm.Get("summary"))
				assert.Equal(t, "details", req.PostForm.Get("description"))
				assert.Equal(t, "main", req.PostForm.Get("base"))
				assert.Equal(t, "feature/foo", req.PostForm.Get("branch"))
				return newCreatedJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Create(ctx, "TEST", "repo", "new PR", "details", "main", "feature/foo")
				require.NoError(t, err)
				assert.Equal(t, 2, got.ID)
				assert.Equal(t, "test PR", got.Summary)
			},
		},
		"Create/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "new PR", req.PostForm.Get("summary"))
				assert.Equal(t, "5", req.PostForm.Get("assigneeId"))
				assert.Equal(t, []string{"10", "20"}, req.PostForm["notifiedUserId[]"])
				return newCreatedJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Create(ctx, "TEST", "repo", "new PR", "", "main", "feature/foo",
					c.PullRequest.Option.WithAssigneeID(5),
					c.PullRequest.Option.WithNotifiedUserIDs([]int{10, 20}),
				)
				require.NoError(t, err)
				assert.Equal(t, 2, got.ID)
			},
		},
		"Create/error": {
			doFunc: newAuthErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.PullRequest.Create(ctx, "TEST", "repo", "new PR", "", "main", "feature/foo")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Update": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests/1", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "Updated summary", req.PostForm.Get("summary"))
				return newJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Update(ctx, "TEST", "repo", 1, c.PullRequest.Option.WithSummary("Updated summary"))
				require.NoError(t, err)
				assert.Equal(t, 2, got.ID)
			},
		},
		"Update/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests/1", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "Updated summary", req.PostForm.Get("summary"))
				assert.Equal(t, "a note", req.PostForm.Get("comment"))
				return newJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Update(ctx, "TEST", "repo", 1,
					c.PullRequest.Option.WithSummary("Updated summary"),
					c.PullRequest.Option.WithComment("a note"),
				)
				require.NoError(t, err)
				assert.Equal(t, 2, got.ID)
			},
		},
		"Update/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.PullRequest.Update(ctx, "TEST", "repo", 1, c.PullRequest.Option.WithSummary("Updated summary"))
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

func TestPullRequestAttachmentService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests/1/attachments", req.URL.Path)
				return newJSONResponse(fixture.Attachment.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Attachment.List(ctx, "TEST", "repo", 1)
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 2, got[0].ID)
				assert.Equal(t, 5, got[1].ID)
			},
		},
		"List/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.PullRequest.Attachment.List(ctx, "TEST", "repo", 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Remove": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests/1/attachments/8", req.URL.Path)
				return newJSONResponse(fixture.Attachment.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Attachment.Remove(ctx, "TEST", "repo", 1, 8)
				require.NoError(t, err)
				assert.Equal(t, 8, got.ID)
			},
		},
		"Remove/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.PullRequest.Attachment.Remove(ctx, "TEST", "repo", 1, 8)
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

func TestPullRequestStarService(t *testing.T) {
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
				assert.Equal(t, "2", req.FormValue("pullRequestId"))
				return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.PullRequest.Star.Add(ctx, 2)
				require.NoError(t, err)
			},
		},
		"Add/error": {
			doFunc: newAuthErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				err := c.PullRequest.Star.Add(ctx, 2)
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
				err := c.PullRequest.Star.Remove(ctx, 42)
				require.NoError(t, err)
			},
		},
		"Remove/error": {
			doFunc: newAuthErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				err := c.PullRequest.Star.Remove(ctx, 42)
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

func TestPullRequestOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	s := c.PullRequest.Option

	cases := map[string]struct {
		option  backlog.RequestOption
		wantKey string
	}{
		"WithAssigneeID":      {option: s.WithAssigneeID(1), wantKey: core.ParamAssigneeID.Value()},
		"WithAssigneeIDs":     {option: s.WithAssigneeIDs([]int{1}), wantKey: core.ParamAssigneeIDs.Value()},
		"WithAttachmentIDs":   {option: s.WithAttachmentIDs([]int{1}), wantKey: core.ParamAttachmentIDs.Value()},
		"WithComment":         {option: s.WithComment("note"), wantKey: core.ParamComment.Value()},
		"WithCount":           {option: s.WithCount(20), wantKey: core.ParamCount.Value()},
		"WithCreatedUserIDs":  {option: s.WithCreatedUserIDs([]int{1}), wantKey: core.ParamCreatedUserIDs.Value()},
		"WithDescription":     {option: s.WithDescription("desc"), wantKey: core.ParamDescription.Value()},
		"WithIssueID":         {option: s.WithIssueID(1), wantKey: core.ParamIssueID.Value()},
		"WithIssueIDs":        {option: s.WithIssueIDs([]int{1}), wantKey: core.ParamIssueIDs.Value()},
		"WithNotifiedUserIDs": {option: s.WithNotifiedUserIDs([]int{1}), wantKey: core.ParamNotifiedUserIDs.Value()},
		"WithOffset":          {option: s.WithOffset(0), wantKey: core.ParamOffset.Value()},
		"WithStatusIDs":       {option: s.WithStatusIDs([]int{1}), wantKey: core.ParamStatusIDs.Value()},
		"WithSummary":         {option: s.WithSummary("summary"), wantKey: core.ParamSummary.Value()},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}
