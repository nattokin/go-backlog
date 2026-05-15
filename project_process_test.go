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
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

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
				got, err := c.Project.Version.List(ctx, "TEST")
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, fixture.Version.Single.ID, got[0].ID)
			},
		},
		"All/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Version.List(ctx, "TEST")
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

func TestProjectStatusService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"All": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/statuses", req.URL.Path)
				return mock.NewJSONResponse(fixture.Status.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Status.All(ctx, "TEST")
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, fixture.Status.Single.ID, got[0].ID)
				assert.Equal(t, fixture.Status.Single.Name, got[0].Name)
				assert.Equal(t, fixture.Status.Single.Color, got[0].Color)
			},
		},
		"All/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Status.All(ctx, "TEST")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Create": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/statuses", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "Open", req.PostForm.Get("name"))
				assert.Equal(t, "#ed8077", req.PostForm.Get("color"))
				return mock.NewJSONResponse(fixture.Status.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Status.Create(ctx, "TEST", "Open", "#ed8077")
				require.NoError(t, err)
				assert.Equal(t, fixture.Status.Single.ID, got.ID)
				assert.Equal(t, fixture.Status.Single.Name, got.Name)
				assert.Equal(t, fixture.Status.Single.Color, got.Color)
			},
		},
		"Create/error": {
			doFunc: newAuthErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Status.Create(ctx, "TEST", "Open", "#ed8077")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Update": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/statuses/1", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "Closed", req.PostForm.Get("name"))
				assert.Equal(t, "#f5ab35", req.PostForm.Get("color"))
				return mock.NewJSONResponse(fixture.Status.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Status.Update(
					ctx,
					"TEST",
					1,
					c.Project.Status.Option.WithName("Closed"),
					c.Project.Status.Option.WithColor("#f5ab35"),
				)
				require.NoError(t, err)
				assert.Equal(t, fixture.Status.Single.ID, got.ID)
			},
		},
		"Update/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Status.Update(
					ctx,
					"TEST",
					1,
					c.Project.Status.Option.WithName("Closed"),
				)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Delete": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/statuses/1", req.URL.Path)
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				form, err := url.ParseQuery(string(body))
				require.NoError(t, err)
				assert.Equal(t, "2", form.Get("substituteStatusId"))
				return mock.NewJSONResponse(fixture.Status.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Status.Delete(ctx, "TEST", 1, 2)
				require.NoError(t, err)
				assert.Equal(t, fixture.Status.Single.ID, got.ID)
			},
		},
		"Delete/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Status.Delete(ctx, "TEST", 1, 2)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"UpdateOrder": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/statuses/updateDisplayOrder", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, []string{"2", "1"}, req.PostForm["statusId[]"])
				return mock.NewJSONResponse(fixture.Status.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.Status.UpdateOrder(ctx, "TEST", []int{2, 1})
				require.NoError(t, err)
				assert.Len(t, got, 2)
			},
		},
		"UpdateOrder/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Project.Status.UpdateOrder(ctx, "TEST", []int{2, 1})
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

func TestProjectStatusOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)

	s := c.Project.Status.Option

	cases := map[string]struct {
		option  backlog.RequestOption
		wantKey string
	}{
		"WithColor": {
			option:  s.WithColor("#ed8077"),
			wantKey: core.ParamColor.Value(),
		},
		"WithName": {
			option:  s.WithName("Open"),
			wantKey: core.ParamName.Value(),
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
			option:  s.WithReleaseDueDate("2024-03-15"),
			wantKey: core.ParamReleaseDueDate.Value(),
		},
		"WithStartDate": {
			option:  s.WithStartDate("2024-03-15"),
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
