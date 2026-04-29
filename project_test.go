package backlog_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestProjectService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"All": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects", req.URL.Path)
				assert.Equal(t, "true", req.URL.Query().Get("archived"))
				return mock.NewJSONResponse(fixture.Project.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.All(ctx, c.Project.Option.WithArchived(true))
				require.NoError(t, err)
				assert.Len(t, got, 3)
			},
		},
		"All/error": {
			doFunc: newAuthErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.All(ctx)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"One": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST", req.URL.Path)
				return mock.NewJSONResponse(fixture.Project.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.One(ctx, "TEST")
				require.NoError(t, err)
				assert.Equal(t, "TEST", got.ProjectKey)
			},
		},
		"One/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.One(ctx, "TEST")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Create": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/projects", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "TEST", req.PostForm.Get("key"))
				assert.Equal(t, "test", req.PostForm.Get("name"))
				assert.Equal(t, "true", req.PostForm.Get("chartEnabled"))
				return mock.NewJSONResponse(fixture.Project.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Create(ctx, "TEST", "test", c.Project.Option.WithChartEnabled(true))
				require.NoError(t, err)
				assert.Equal(t, "TEST", got.ProjectKey)
			},
		},
		"Create/error": {
			doFunc: newAuthErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Create(ctx, "TEST", "test")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Update": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "new-name", req.PostForm.Get("name"))
				return mock.NewJSONResponse(fixture.Project.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Update(ctx, "TEST", c.Project.Option.WithName("new-name"))
				require.NoError(t, err)
				assert.Equal(t, "TEST", got.ProjectKey)
			},
		},
		"Update/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Update(ctx, "TEST")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Delete": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST", req.URL.Path)
				return mock.NewJSONResponse(fixture.Project.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Delete(ctx, "TEST")
				require.NoError(t, err)
				assert.Equal(t, "TEST", got.ProjectKey)
			},
		},
		"Delete/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Delete(ctx, "TEST")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"DiskUsage": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/diskUsage", req.URL.Path)
				return mock.NewJSONResponse(fixture.Project.DiskUsageJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.DiskUsage(ctx, "TEST")
				require.NoError(t, err)
				assert.Equal(t, 1, got.ProjectID)
				assert.Equal(t, 11931, got.Issue)
			},
		},
		"DiskUsage/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.DiskUsage(ctx, "TEST")
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

func TestProjectCategoryService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"All": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/categories", req.URL.Path)
				return mock.NewJSONResponse(fixture.Category.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Category.All(ctx, "TEST")
				require.NoError(t, err)
				assert.Len(t, got, 2)
			},
		},
		"All/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Category.All(ctx, "TEST")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Create": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/categories", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "Bug", req.PostForm.Get("name"))
				return mock.NewJSONResponse(fixture.Category.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Category.Create(ctx, "TEST", "Bug")
				require.NoError(t, err)
				assert.Equal(t, "Bug", got.Name)
			},
		},
		"Create/error": {
			doFunc: newAuthErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Category.Create(ctx, "TEST", "Bug")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Update": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/categories/12", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "Bug Fixed", req.PostForm.Get("name"))
				return mock.NewJSONResponse(fixture.Category.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Category.Update(ctx, "TEST", 12, "Bug Fixed")
				require.NoError(t, err)
				assert.Equal(t, 12, got.ID)
			},
		},
		"Update/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Category.Update(ctx, "TEST", 12, "Bug Fixed")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Delete": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/categories/12", req.URL.Path)
				return mock.NewJSONResponse(fixture.Category.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Category.Delete(ctx, "TEST", 12)
				require.NoError(t, err)
				assert.Equal(t, 12, got.ID)
			},
		},
		"Delete/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Category.Delete(ctx, "TEST", 12)
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

func TestProjectSharedFileService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/files", req.URL.Path)
				return mock.NewJSONResponse(fixture.SharedFile.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.SharedFile.List(ctx, "TEST")
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 454403, got[0].ID)
				assert.Equal(t, "01_buz.png", got[0].Name)
				assert.Equal(t, 454404, got[1].ID)
				assert.Equal(t, "readme.md", got[1].Name)
			},
		},
		"List/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.SharedFile.List(ctx, "TEST")
				require.Error(t, err)
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

func TestProjectVersionService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"All": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/versions", req.URL.Path)
				return mock.NewJSONResponse(fixture.Version.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Version.All(ctx, "TEST")
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, fixture.Version.Single.ID, got[0].ID)
			},
		},
		"All/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Version.All(ctx, "TEST")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Create": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/versions", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "v1.0.0", req.PostForm.Get("name"))
				assert.Equal(t, "first release", req.PostForm.Get("description"))
				return mock.NewJSONResponse(fixture.Version.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Version.Create(
					ctx,
					"TEST",
					"v1.0.0",
					c.Project.Version.Option.WithDescription("first release"),
				)
				require.NoError(t, err)
				assert.Equal(t, fixture.Version.Single.ID, got.ID)
			},
		},
		"Create/error": {
			doFunc: newAuthErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Version.Create(ctx, "TEST", "v1.0.0")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Update": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/versions/1", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "v1.0.1", req.PostForm.Get("name"))
				return mock.NewJSONResponse(fixture.Version.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Version.Update(
					ctx,
					"TEST",
					1,
					c.Project.Version.Option.WithName("v1.0.1"),
				)
				require.NoError(t, err)
				assert.Equal(t, fixture.Version.Single.ID, got.ID)
			},
		},
		"Update/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Version.Update(
					ctx,
					"TEST",
					1,
					c.Project.Version.Option.WithName("v1.0.1"),
				)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Delete": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/versions/1", req.URL.Path)
				return mock.NewJSONResponse(fixture.Version.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Version.Delete(ctx, "TEST", 1)
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.Equal(t, fixture.Version.Single.ID, got.ID)
			},
		},
		"Delete/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Version.Delete(ctx, "TEST", 1)
				require.Error(t, err)
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
				return mock.NewJSONResponse(fixture.User.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.All(ctx, "TEST", false)
				require.NoError(t, err)
				assert.Len(t, got, 4)
			},
		},
		"All/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.User.All(ctx, "TEST", false)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Add": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/users", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "1", req.PostForm.Get("userId"))
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.Add(ctx, "TEST", 1)
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
			},
		},
		"Add/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.User.Add(ctx, "TEST", 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
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
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.Delete(ctx, "TEST", 1)
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
			},
		},
		"Delete/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.User.Delete(ctx, "TEST", 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"AddAdmin": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/administrators", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "1", req.PostForm.Get("userId"))
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.AddAdmin(ctx, "TEST", 1)
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
			},
		},
		"AddAdmin/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.User.AddAdmin(ctx, "TEST", 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"AdminAll": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/administrators", req.URL.Path)
				return mock.NewJSONResponse(fixture.User.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.AdminAll(ctx, "TEST")
				require.NoError(t, err)
				assert.Len(t, got, 4)
			},
		},
		"AdminAll/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.User.AdminAll(ctx, "TEST")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
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
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.DeleteAdmin(ctx, "TEST", 1)
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
			},
		},
		"DeleteAdmin/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.User.DeleteAdmin(ctx, "TEST", 1)
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

func TestProjectOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	s := c.Project.Option

	cases := map[string]struct {
		option  backlog.RequestOption
		wantKey string
	}{
		"WithAll": {
			option:  s.WithAll(true),
			wantKey: core.ParamAll.Value(),
		},
		"WithArchived": {
			option:  s.WithArchived(true),
			wantKey: core.ParamArchived.Value(),
		},
		"WithChartEnabled": {
			option:  s.WithChartEnabled(true),
			wantKey: core.ParamChartEnabled.Value(),
		},
		"WithKey": {
			option:  s.WithKey("TEST"),
			wantKey: core.ParamKey.Value(),
		},
		"WithName": {
			option:  s.WithName("test"),
			wantKey: core.ParamName.Value(),
		},
		"WithProjectLeaderCanEditProjectLeader": {
			option:  s.WithProjectLeaderCanEditProjectLeader(true),
			wantKey: core.ParamProjectLeaderCanEditProjectLeader.Value(),
		},
		"WithSubtaskingEnabled": {
			option:  s.WithSubtaskingEnabled(true),
			wantKey: core.ParamSubtaskingEnabled.Value(),
		},
		"WithTextFormattingRule": {
			option:  s.WithTextFormattingRule(model.FormatBacklog),
			wantKey: core.ParamTextFormattingRule.Value(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantKey, tc.option.Key())
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

func TestProjectVersionOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)

	date := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	s := c.Project.Version.Option

	cases := map[string]struct {
		option  backlog.RequestOption
		wantKey string
	}{
		"WithArchived": {
			option:  s.WithArchived(true),
			wantKey: core.ParamArchived.Value(),
		},
		"WithDescription": {
			option:  s.WithDescription("desc"),
			wantKey: core.ParamDescription.Value(),
		},
		"WithName": {
			option:  s.WithName("v1.0.0"),
			wantKey: core.ParamName.Value(),
		},
		"WithReleaseDueDate": {
			option:  s.WithReleaseDueDate(date),
			wantKey: core.ParamReleaseDueDate.Value(),
		},
		"WithStartDate": {
			option:  s.WithStartDate(date),
			wantKey: core.ParamStartDate.Value(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}
