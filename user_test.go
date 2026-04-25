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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.All(ctx, "TEST", false)
				require.NoError(t, err)
				assert.Len(t, got, 4)
			},
		},
		"All/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such project.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.Add(ctx, "TEST", 1)
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
			},
		},
		"Add/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such project.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.Delete(ctx, "TEST", 1)
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
			},
		},
		"Delete/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such project.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.AddAdmin(ctx, "TEST", 1)
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
			},
		},
		"AddAdmin/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such project.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.AdminAll(ctx, "TEST")
				require.NoError(t, err)
				assert.Len(t, got, 4)
			},
		},
		"AdminAll/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such project.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Project.User.DeleteAdmin(ctx, "TEST", 1)
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
			},
		},
		"DeleteAdmin/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such project.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
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

func TestUserService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"All": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/users", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.User.All(ctx)
				require.NoError(t, err)
				assert.Len(t, got, 4)
			},
		},
		"All/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return newAuthErrorResponse(), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.User.All(ctx)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"One": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/users/1", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.User.One(ctx, 1)
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
			},
		},
		"One/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such user.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.User.One(ctx, 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Own": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/users/myself", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.User.Own(ctx)
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
			},
		},
		"Own/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return newAuthErrorResponse(), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.User.Own(ctx)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Add": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/users", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "newuser", req.PostForm.Get("userId"))
				assert.Equal(t, "password", req.PostForm.Get("password"))
				assert.Equal(t, "New User", req.PostForm.Get("name"))
				assert.Equal(t, "new@example.com", req.PostForm.Get("mailAddress"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.User.Add(ctx, "newuser", "password", "New User", "new@example.com", backlog.RoleAdministrator)
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
			},
		},
		"Add/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return newAuthErrorResponse(), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.User.Add(ctx, "newuser", "password", "New User", "new@example.com", backlog.RoleGuestReporter)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Update": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/users/1", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "updated-user", req.PostForm.Get("name"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.User.Update(ctx, 1, c.User.Option.WithName("updated-user"))
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
			},
		},
		"Update/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such user.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.User.Update(ctx, 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Delete": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/users/1", req.URL.Path)
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				form, err := url.ParseQuery(string(body))
				require.NoError(t, err)
				assert.Empty(t, form)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.User.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.User.Delete(ctx, 1)
				require.NoError(t, err)
				assert.Equal(t, "admin", got.UserID)
			},
		},
		"Delete/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such user.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.User.Delete(ctx, 1)
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

func TestUserActivityService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/users/1/activities", req.URL.Path)
				assert.Equal(t, "5", req.URL.Query().Get("minId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Activity.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.User.Activity.List(ctx, 1, c.User.Activity.Option.WithMinID(5))
				require.NoError(t, err)
				assert.Len(t, got, 1)
			},
		},
		"List/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such user.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.User.Activity.List(ctx, 1)
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

func TestUserOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	o := c.User.Option

	// --- Boolean options ------------------------------------------------------------
	t.Run("boolean-options", func(t *testing.T) {
		cases := map[string]struct {
			option    backlog.RequestOption
			key       string
			wantValue bool
		}{
			"with-form-send-mail": {
				option:    o.WithSendMail(true),
				key:       core.ParamSendMail.Value(),
				wantValue: true,
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
			"with-form-user-id": {
				option:    o.WithUserID(1),
				key:       core.ParamUserID.Value(),
				wantValue: 1,
			},
			"with-form-role-type": {
				option:    o.WithRoleType(2),
				key:       core.ParamRoleType.Value(),
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
			"with-form-name": {
				option:    o.WithName("example-user"),
				key:       core.ParamName.Value(),
				wantValue: "example-user",
			},
			"with-form-mail-address": {
				option:    o.WithMailAddress("user@example.com"),
				key:       core.ParamMailAddress.Value(),
				wantValue: "user@example.com",
			},
			"with-form-password": {
				option:    o.WithPassword("securepass"),
				key:       core.ParamPassword.Value(),
				wantValue: "securepass",
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
