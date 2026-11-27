package backlog_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestNotificationService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/notifications", req.URL.Path)
				return mock.NewResponse(fixture.Notification.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Notification.List(ctx)
				require.NoError(t, err)
				assert.Equal(t, fixture.Notification.List, got)
			},
		},
		"List/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				q := req.URL.Query()
				assert.Equal(t, "10", q.Get("minId"))
				assert.Equal(t, "20", q.Get("maxId"))
				assert.Equal(t, "5", q.Get("count"))
				assert.Equal(t, "asc", q.Get("order"))
				assert.Equal(t, "2", q.Get("senderId"))
				return mock.NewResponse(fixture.Notification.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				o := c.Notification.Option
				got, err := c.Notification.List(ctx,
					o.WithMinID(10),
					o.WithMaxID(20),
					o.WithCount(5),
					o.WithOrder(backlog.OrderAsc),
					o.WithSenderID(2),
				)
				require.NoError(t, err)
				assert.Len(t, got, 1)
			},
		},
		"List/error": {
			doFunc: mock.NewInternalServerErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Notification.List(ctx)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Count": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/notifications/count", req.URL.Path)
				q := req.URL.Query()
				assert.Equal(t, "false", q.Get("alreadyRead"))
				assert.Equal(t, "false", q.Get("resourceAlreadyRead"))
				return mock.NewResponse(fixture.Notification.CountJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				o := c.Notification.Option
				got, err := c.Notification.Count(ctx,
					o.WithAlreadyRead(false),
					o.WithResourceAlreadyRead(false),
				)
				require.NoError(t, err)
				assert.Equal(t, 138, got)
			},
		},
		"Count/error": {
			doFunc: mock.NewInternalServerErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Notification.Count(ctx)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"ResetUnreadCount": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/notifications/markAsRead", req.URL.Path)
				return mock.NewResponse(fixture.Notification.ResetUnreadCountJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Notification.ResetUnreadCount(ctx)
				require.NoError(t, err)
				assert.Equal(t, 4, got)
			},
		},
		"ResetUnreadCount/error": {
			doFunc: mock.NewInternalServerErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Notification.ResetUnreadCount(ctx)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"MarkAsRead": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/notifications/299/markAsRead", req.URL.Path)
				return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Notification.MarkAsRead(ctx, 299)
				require.NoError(t, err)
			},
		},
		"MarkAsRead/error": {
			doFunc: mock.NewInternalServerErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Notification.MarkAsRead(ctx, 299)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"MarkAsRead/validation-error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				t.Error("request must not be sent")
				return nil, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Notification.MarkAsRead(ctx, 0)
				var target *backlog.ValidationError
				assert.ErrorAs(t, err, &target)
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

func TestNotificationOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "apiKey")
	require.NoError(t, err)
	o := c.Notification.Option

	// --- Boolean options ------------------------------------------------------------
	t.Run("boolean-options", func(t *testing.T) {
		cases := map[string]struct {
			option    backlog.RequestOption
			key       string
			wantValue bool
		}{
			"with-already-read": {
				option:    o.WithAlreadyRead(true),
				key:       option.ParamAlreadyRead.Value(),
				wantValue: true,
			},
			"with-resource-already-read": {
				option:    o.WithResourceAlreadyRead(false),
				key:       option.ParamResourceAlreadyRead.Value(),
				wantValue: false,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.Set(form)
				require.NoError(t, err)
				assert.Equal(t, strconv.FormatBool(tc.wantValue), form.Get(tc.key))
			})
		}
	})

	// --- Integer options ------------------------------------------------------------
	t.Run("integer-options", func(t *testing.T) {
		cases := map[string]struct {
			option    backlog.RequestOption
			key       string
			wantValue int
		}{
			"with-count": {
				option:    o.WithCount(20),
				key:       option.ParamCount.Value(),
				wantValue: 20,
			},
			"with-max-id": {
				option:    o.WithMaxID(100),
				key:       option.ParamMaxID.Value(),
				wantValue: 100,
			},
			"with-min-id": {
				option:    o.WithMinID(10),
				key:       option.ParamMinID.Value(),
				wantValue: 10,
			},
			"with-sender-id": {
				option:    o.WithSenderID(2),
				key:       option.ParamSenderID.Value(),
				wantValue: 2,
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
				key:       option.ParamOrder.Value(),
				wantValue: "asc",
			},
			"with-order-desc": {
				option:    o.WithOrder(backlog.OrderDesc),
				key:       option.ParamOrder.Value(),
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
