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

func TestProjectActivityService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/activities", req.URL.Path)
				assert.Equal(t, "10", req.URL.Query().Get("count"))
				return mock.NewJSONResponse(fixture.Activity.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Activity.List(ctx, "TEST", c.Project.Activity.Option.WithCount(10))
				require.NoError(t, err)
				assert.Len(t, got, 1)
			},
		},
		"List/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Activity.List(ctx, "TEST")
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

func TestProjectWebhookService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"All": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/webhooks", req.URL.Path)
				return mock.NewJSONResponse(fixture.Webhook.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Webhook.All(ctx, "TEST")
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 1, got[0].ID)
			},
		},
		"All/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Webhook.All(ctx, "TEST")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Create": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/webhooks", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "notify", req.PostForm.Get("name"))
				assert.Equal(t, "https://example.com/webhook", req.PostForm.Get("hookUrl"))
				assert.Equal(t, "true", req.PostForm.Get("allEvent"))
				return mock.NewJSONResponse(fixture.Webhook.AllEventJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Webhook.Create(
					ctx,
					"TEST",
					"notify",
					"https://example.com/webhook",
					c.Project.Webhook.Option.WithAllEvent(true),
				)
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
			},
		},
		"Create/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Webhook.Create(
					ctx,
					"TEST",
					"notify",
					"https://example.com/webhook",
					c.Project.Webhook.Option.WithAllEvent(true),
				)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"One": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/webhooks/1", req.URL.Path)
				return mock.NewJSONResponse(fixture.Webhook.AllEventJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Webhook.One(ctx, "TEST", 1)
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
			},
		},
		"One/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Webhook.One(ctx, "TEST", 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Update": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/webhooks/1", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "updated", req.PostForm.Get("name"))
				return mock.NewJSONResponse(fixture.Webhook.AllEventJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Webhook.Update(
					ctx,
					"TEST",
					1,
					c.Project.Webhook.Option.WithName("updated"),
				)
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
			},
		},
		"Update/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Webhook.Update(
					ctx,
					"TEST",
					1,
					c.Project.Webhook.Option.WithName("updated"),
				)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Delete": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/webhooks/1", req.URL.Path)
				return mock.NewJSONResponse(fixture.Webhook.AllEventJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Webhook.Delete(ctx, "TEST", 1)
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.Equal(t, 1, got.ID)
			},
		},
		"Delete/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Webhook.Delete(ctx, "TEST", 1)
				require.Error(t, err)
				assert.Nil(t, got)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c, err := backlog.NewClient(
				"https://example.backlog.com",
				"token",
				backlog.WithDoer(&mockDoer{do: tc.doFunc}),
			)
			require.NoError(t, err)
			tc.call(t, c)
		})
	}
}

func TestProjectWebhookOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)

	s := c.Project.Webhook.Option

	cases := map[string]struct {
		option  backlog.RequestOption
		wantKey string
	}{
		"WithAllEvent": {
			option:  s.WithAllEvent(true),
			wantKey: core.ParamAllEvent.Value(),
		},
		"WithActivityTypeIDs": {
			option:  s.WithActivityTypeIDs([]int{1, 2}),
			wantKey: core.ParamActivityTypeIDs.Value(),
		},
		"WithHookURL": {
			option:  s.WithHookURL("https://example.com/webhook"),
			wantKey: core.ParamHookURL.Value(),
		},
		"WithName": {
			option:  s.WithName("notify"),
			wantKey: core.ParamName.Value(),
		},
		"WithDescription": {
			option:  s.WithDescription("desc"),
			wantKey: core.ParamDescription.Value(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}
