package backlog_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestPullRequestCommentService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests/1/comments", req.URL.Path)
				return mock.NewJSONResponse(fixture.Comment.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Comment.List(ctx, "TEST", "repo", 1)
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 1, got[0].ID)
				assert.Equal(t, 2, got[1].ID)
			},
		},
		"List/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests/1/comments", req.URL.Path)
				assert.Equal(t, "20", req.URL.Query().Get("count"))
				assert.Equal(t, "asc", req.URL.Query().Get("order"))
				return mock.NewJSONResponse(fixture.Comment.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Comment.List(ctx, "TEST", "repo", 1,
					c.PullRequest.Comment.Option.WithCount(20),
					c.PullRequest.Comment.Option.WithOrder(backlog.OrderAsc),
				)
				require.NoError(t, err)
				assert.Len(t, got, 2)
			},
		},
		"All/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.PullRequest.Comment.List(ctx, "TEST", "repo", 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Add": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests/1/comments", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "This is a comment.", req.PostForm.Get("content"))
				return mock.NewCreatedJSONResponse(fixture.Comment.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Comment.Add(ctx, "TEST", "repo", 1, "This is a comment.")
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
				assert.Equal(t, "This is a comment.", got.Content)
			},
		},
		"Add/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests/1/comments", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "Notifying users.", req.PostForm.Get("content"))
				assert.Equal(t, []string{"5", "6"}, req.PostForm["notifiedUserId[]"])
				return mock.NewCreatedJSONResponse(fixture.Comment.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Comment.Add(ctx, "TEST", "repo", 1, "Notifying users.",
					c.PullRequest.Comment.Option.WithNotifiedUserIDs([]int{5, 6}),
				)
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
			},
		},
		"Add/error": {
			doFunc: newAuthErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.PullRequest.Comment.Add(ctx, "TEST", "repo", 1, "x")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Count": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests/1/comments/count", req.URL.Path)
				return mock.NewJSONResponse(`{"count":7}`), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Comment.Count(ctx, "TEST", "repo", 1)
				require.NoError(t, err)
				assert.Equal(t, 7, got)
			},
		},
		"Count/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.PullRequest.Comment.Count(ctx, "TEST", "repo", 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"One": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests/1/comments/42", req.URL.Path)
				return mock.NewJSONResponse(fixture.Comment.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Comment.One(ctx, "TEST", "repo", 1, 42)
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
				assert.Equal(t, "This is a comment.", got.Content)
			},
		},
		"One/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.PullRequest.Comment.One(ctx, "TEST", "repo", 1, 42)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Update": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests/1/comments/42", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "Updated content.", req.PostForm.Get("content"))
				return mock.NewJSONResponse(fixture.Comment.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Comment.Update(ctx, "TEST", "repo", 1, 42, "Updated content.")
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
			},
		},
		"Update/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.PullRequest.Comment.Update(ctx, "TEST", "repo", 1, 42, "Updated content.")
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

func TestPullRequestCommentOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	s := c.PullRequest.Comment.Option

	cases := map[string]struct {
		option  backlog.RequestOption
		wantKey string
	}{
		"WithAttachmentIDs":   {option: s.WithAttachmentIDs([]int{1}), wantKey: core.ParamAttachmentIDs.Value()},
		"WithCount":           {option: s.WithCount(20), wantKey: core.ParamCount.Value()},
		"WithMaxID":           {option: s.WithMaxID(100), wantKey: core.ParamMaxID.Value()},
		"WithMinID":           {option: s.WithMinID(1), wantKey: core.ParamMinID.Value()},
		"WithNotifiedUserIDs": {option: s.WithNotifiedUserIDs([]int{1}), wantKey: core.ParamNotifiedUserIDs.Value()},
		"WithOrder":           {option: s.WithOrder(backlog.OrderAsc), wantKey: core.ParamOrder.Value()},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}
