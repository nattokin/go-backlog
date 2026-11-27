package backlog_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestWatchingService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/users/1/watchings", req.URL.Path)
				assert.Empty(t, req.URL.RawQuery)
				return mock.NewResponse(fixture.Watching.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Watching.List(ctx, 1)
				require.NoError(t, err)
				assert.Equal(t, fixture.Watching.List, got)
			},
		},
		"List/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, url.Values{
					"order":               {"asc"},
					"sort":                {"issueUpdated"},
					"count":               {"10"},
					"offset":              {"20"},
					"resourceAlreadyRead": {"false"},
					"issueId[]":           {"3", "4"},
				}, req.URL.Query())
				return mock.NewResponse(fixture.Watching.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				o := c.Watching.Option
				got, err := c.Watching.List(ctx, 1,
					o.WithOrder(backlog.OrderAsc),
					o.WithSort(backlog.WatchingSortIssueUpdated),
					o.WithCount(10),
					o.WithOffset(20),
					o.WithResourceAlreadyRead(false),
					o.WithIssueIDs([]int{3, 4}),
				)
				require.NoError(t, err)
				assert.Len(t, got, 2)
			},
		},
		"List/invalid": {
			doFunc: mock.NewUnexpectedDoFunc(t),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Watching.List(ctx, 0, c.Watching.Option.WithSort("summary"))
				require.Error(t, err)
				var target *backlog.ValidationError
				assert.ErrorAs(t, err, &target)
			},
		},
		"List/error": {
			doFunc: mock.NewInternalServerErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Watching.List(ctx, 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Count": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/users/1/watchings/count", req.URL.Path)
				assert.Equal(t, url.Values{
					"resourceAlreadyRead": {"true"},
					"alreadyRead":         {"false"},
				}, req.URL.Query())
				return mock.NewResponse(fixture.Watching.CountJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				o := c.Watching.Option
				got, err := c.Watching.Count(ctx, 1,
					o.WithResourceAlreadyRead(true),
					o.WithAlreadyRead(false),
				)
				require.NoError(t, err)
				assert.Equal(t, 138, got)
			},
		},
		"Count/error": {
			doFunc: mock.NewInternalServerErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Watching.Count(ctx, 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"One": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/watchings/1", req.URL.Path)
				return mock.NewResponse(fixture.Watching.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Watching.One(ctx, 1)
				require.NoError(t, err)
				assert.Equal(t, &fixture.Watching.Single, got)
			},
		},
		"One/error": {
			doFunc: mock.NewNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Watching.One(ctx, 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Create": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/watchings", req.URL.Path)
				assert.Empty(t, req.URL.RawQuery)
				assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
				require.NoError(t, req.ParseForm())
				assert.Equal(t, url.Values{
					"issueIdOrKey": {"PRJ-1"},
					"note":         {"Check before release"},
				}, req.PostForm)
				return mock.NewCreatedResponse(fixture.Watching.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Watching.Create(ctx, "PRJ-1", c.Watching.Option.WithNote("Check before release"))
				require.NoError(t, err)
				assert.Equal(t, &fixture.Watching.Single, got)
			},
		},
		"Create/invalid": {
			doFunc: mock.NewUnexpectedDoFunc(t),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Watching.Create(ctx, "")
				require.Error(t, err)
				var target *backlog.ValidationError
				assert.ErrorAs(t, err, &target)
			},
		},
		"Create/error": {
			doFunc: mock.NewNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Watching.Create(ctx, "PRJ-1")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Update": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/watchings/1", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, url.Values{"note": {"Check before release"}}, req.PostForm)
				return mock.NewResponse(fixture.Watching.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Watching.Update(ctx, 1, "Check before release")
				require.NoError(t, err)
				assert.Equal(t, &fixture.Watching.Single, got)
			},
		},
		"Update/error": {
			doFunc: mock.NewNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Watching.Update(ctx, 1, "note")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Delete": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/watchings/1", req.URL.Path)
				return mock.NewResponse(fixture.Watching.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Watching.Delete(ctx, 1)
				require.NoError(t, err)
				assert.Equal(t, &fixture.Watching.Single, got)
			},
		},
		"Delete/error": {
			doFunc: mock.NewNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Watching.Delete(ctx, 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"MarkAsRead": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/watchings/1/markAsRead", req.URL.Path)
				return mock.NewNoContentResponse(), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Watching.MarkAsRead(ctx, 1)
				require.NoError(t, err)
			},
		},
		"MarkAsRead/invalid": {
			doFunc: mock.NewUnexpectedDoFunc(t),
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Watching.MarkAsRead(ctx, 0)
				require.Error(t, err)
				var target *backlog.ValidationError
				assert.ErrorAs(t, err, &target)
			},
		},
		"MarkAsRead/error": {
			doFunc: mock.NewNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Watching.MarkAsRead(ctx, 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "apiKey", backlog.WithDoer(&mock.Doer{DoFunc: tc.doFunc}))
			require.NoError(t, err)
			tc.call(t, c)
		})
	}
}
