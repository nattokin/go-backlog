package backlog_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

func TestUserRecentlyViewedService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"ListIssues": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/users/myself/recentlyViewedIssues", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.RecentlyViewed.IssueListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.User.RecentlyViewed.ListIssues(ctx)
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, "TEST-1", got[0].IssueKey)
			},
		},
		"ListIssues/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				q := req.URL.Query()
				assert.Equal(t, "5", q.Get("count"))
				assert.Equal(t, "10", q.Get("offset"))
				assert.Equal(t, "asc", q.Get("order"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.RecentlyViewed.IssueListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				o := c.User.RecentlyViewed.Option
				got, err := c.User.RecentlyViewed.ListIssues(ctx,
					o.WithCount(5),
					o.WithOffset(10),
					o.WithOrder(backlog.OrderAsc),
				)
				require.NoError(t, err)
				assert.Len(t, got, 2)
			},
		},
		"ListIssues/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"Authentication failure.","code":11,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.User.RecentlyViewed.ListIssues(ctx)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"AddIssue": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/issues/1/recentlyViewedIssues", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.RecentlyViewed.IssueSingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.User.RecentlyViewed.AddIssue(ctx, 1)
				require.NoError(t, err)
				assert.Equal(t, "TEST-1", got.IssueKey)
			},
		},
		"AddIssue/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such issue.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.User.RecentlyViewed.AddIssue(ctx, 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"ListProjects": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/users/myself/recentlyViewedProjects", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.RecentlyViewed.ProjectListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.User.RecentlyViewed.ListProjects(ctx)
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, "TEST", got[0].ProjectKey)
			},
		},
		"ListProjects/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"Authentication failure.","code":11,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.User.RecentlyViewed.ListProjects(ctx)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"ListWikis": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/users/myself/recentlyViewedWikis", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.RecentlyViewed.WikiListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.User.RecentlyViewed.ListWikis(ctx)
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, "Home", got[0].Name)
			},
		},
		"ListWikis/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"Authentication failure.","code":11,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.User.RecentlyViewed.ListWikis(ctx)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"AddWiki": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/wikis/10/recentlyViewedWikis", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.RecentlyViewed.WikiSingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.User.RecentlyViewed.AddWiki(ctx, 10)
				require.NoError(t, err)
				assert.Equal(t, "Home", got.Name)
			},
		},
		"AddWiki/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such wiki.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.User.RecentlyViewed.AddWiki(ctx, 10)
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

func TestUserRecentlyViewedOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	o := c.User.RecentlyViewed.Option

	// --- Integer options ------------------------------------------------------------
	t.Run("integer-options", func(t *testing.T) {
		cases := map[string]struct {
			option    backlog.RequestOption
			key       string
			wantValue int
		}{
			"with-count": {
				option:    o.WithCount(20),
				key:       core.ParamCount.Value(),
				wantValue: 20,
			},
			"with-offset": {
				option:    o.WithOffset(10),
				key:       core.ParamOffset.Value(),
				wantValue: 10,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.Set(form)
				require.NoError(t, err)
				assert.Equal(t, strconv.Itoa(tc.wantValue), form.Get(tc.key))
			})
		}
	})

	// --- String options -------------------------------------------------------------
	t.Run("string-options", func(t *testing.T) {
		cases := map[string]struct {
			option    backlog.RequestOption
			key       string
			wantValue string
		}{
			"with-order-asc": {
				option:    o.WithOrder(backlog.OrderAsc),
				key:       core.ParamOrder.Value(),
				wantValue: "asc",
			},
			"with-order-desc": {
				option:    o.WithOrder(backlog.OrderDesc),
				key:       core.ParamOrder.Value(),
				wantValue: "desc",
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.Set(form)
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, form.Get(tc.key))
			})
		}
	})
}
