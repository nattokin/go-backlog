package backlog_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

func TestUserStarService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/users/1/stars", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Star.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.User.Star.List(ctx, 1)
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 10, got[0].ID)
				assert.Equal(t, "admin", got[0].Presenter.UserID)
			},
		},
		"List/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, "/api/v2/users/1/stars", req.URL.Path)
				q := req.URL.Query()
				assert.Equal(t, "5", q.Get("count"))
				assert.Equal(t, "100", q.Get("minId"))
				assert.Equal(t, "200", q.Get("maxId"))
				assert.Equal(t, "asc", q.Get("order"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Star.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				o := c.User.Star.Option
				got, err := c.User.Star.List(ctx, 1,
					o.WithCount(5),
					o.WithMinID(100),
					o.WithMaxID(200),
					o.WithOrder(backlog.OrderAsc),
				)
				require.NoError(t, err)
				assert.Len(t, got, 2)
			},
		},
		"List/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.User.Star.List(ctx, 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Count": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/users/1/stars/count", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Star.CountJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.User.Star.Count(ctx, 1)
				require.NoError(t, err)
				assert.Equal(t, 42, got)
			},
		},
		"Count/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.User.Star.Count(ctx, 1)
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

func TestUserStarOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	o := c.User.Star.Option

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
			"with-min-id": {
				option:    o.WithMinID(10),
				key:       core.ParamMinID.Value(),
				wantValue: 10,
			},
			"with-max-id": {
				option:    o.WithMaxID(99),
				key:       core.ParamMaxID.Value(),
				wantValue: 99,
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
