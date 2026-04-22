package backlog_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

func TestSpaceService_One(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc  func(req *http.Request) (*http.Response, error)
		wantErr bool
	}{
		"success": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/space", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Space.SpaceJSON))),
				}, nil
			},
		},
		"error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"Authentication failure.","code":11,"moreInfo":""}]}`)),
				}, nil
			},
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)

			got, err := c.Space.One(ctx)

			if tc.wantErr {
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, "nulab", got.SpaceKey)
			assert.Equal(t, "Nulab Inc.", got.Name)
			assert.Equal(t, backlog.FormatMarkdown, got.TextFormattingRule)
		})
	}
}

func TestSpaceService_DiskUsage(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc  func(req *http.Request) (*http.Response, error)
		wantErr bool
	}{
		"success": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/space/diskUsage", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Space.DiskUsageJSON))),
				}, nil
			},
		},
		"error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"Authentication failure.","code":11,"moreInfo":""}]}`)),
				}, nil
			},
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)

			got, err := c.Space.DiskUsage(ctx)

			if tc.wantErr {
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, 1073741824, got.Capacity)
			assert.Equal(t, 119511, got.Issue)
			require.Len(t, got.Details, 1)
			assert.Equal(t, 1, got.Details[0].ProjectID)
			assert.Equal(t, 11931, got.Details[0].Issue)
		})
	}
}

func TestSpaceService_Notification(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc  func(req *http.Request) (*http.Response, error)
		wantErr bool
	}{
		"success": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/space/notification", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Space.NotificationJSON))),
				}, nil
			},
		},
		"error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"Authentication failure.","code":11,"moreInfo":""}]}`)),
				}, nil
			},
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)

			got, err := c.Space.Notification(ctx)

			if tc.wantErr {
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, "Backlog is a project management tool.", got.Content)
		})
	}
}

func TestSpaceService_UpdateNotification(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		content string
		doFunc  func(req *http.Request) (*http.Response, error)
		wantErr bool
	}{
		"success": {
			content: "Backlog is a project management tool.",
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPut, req.Method)
				assert.Equal(t, "/api/v2/space/notification", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "Backlog is a project management tool.", req.FormValue("content"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Space.NotificationJSON))),
				}, nil
			},
		},
		"error-validation-empty-content": {
			content: "",
			wantErr: true,
		},
		"error-api": {
			content: "content",
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"Authentication failure.","code":11,"moreInfo":""}]}`)),
				}, nil
			},
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)

			got, err := c.Space.UpdateNotification(ctx, tc.content)

			if tc.wantErr {
				require.Error(t, err)
				var apiErr *backlog.APIResponseError
				if tc.doFunc != nil {
					assert.True(t, errors.As(err, &apiErr))
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, "Backlog is a project management tool.", got.Content)
		})
	}
}

func TestSpaceActivityService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/space/activities", req.URL.Path)
				assert.Equal(t, "20", req.URL.Query().Get("count"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Activity.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Space.Activity.List(ctx, c.Space.Activity.Option.WithCount(20))
				require.NoError(t, err)
				assert.Len(t, got, 1)
			},
		},
		"List/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"Authentication failure.","code":11,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Space.Activity.List(ctx)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Get": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/activities/3153", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Activity.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Space.Activity.Get(ctx, 3153)
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.Equal(t, 3153, got.ID)
				assert.Equal(t, 2, got.Type)
			},
		},
		"Get/error-invalid-id": {
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Space.Activity.Get(ctx, 0)
				require.Error(t, err)
			},
		},
		"Get/error-api": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"Authentication failure.","code":11,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Space.Activity.Get(ctx, 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tc.doFunc == nil {
				c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("should not be called")
				}}))
				require.NoError(t, err)
				tc.call(t, c)
				return
			}

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)
			tc.call(t, c)
		})
	}
}

func TestSpaceAttachmentService(t *testing.T) {
	ctx := context.Background()

	t.Run("Upload", func(t *testing.T) {
		doFunc := func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Equal(t, "/api/v2/space/attachment", req.URL.Path)
			assert.True(t, strings.HasPrefix(req.Header.Get("Content-Type"), "multipart/form-data"))
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Attachment.UploadJSON))),
			}, nil
		}

		c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: doFunc}))
		require.NoError(t, err)

		f, err := os.Open("testdata/testfile")
		require.NoError(t, err)
		defer f.Close()

		got, err := c.Space.Attachment.Upload(ctx, "testfile", f)
		require.NoError(t, err)
		assert.Equal(t, 1, got.ID)
		assert.Equal(t, "test.txt", got.Name)
		assert.Equal(t, 8857, got.Size)
	})

	t.Run("Upload/error", func(t *testing.T) {
		doFunc := func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"Authentication failure.","code":11,"moreInfo":""}]}`)),
			}, nil
		}

		c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: doFunc}))
		require.NoError(t, err)

		_, err = c.Space.Attachment.Upload(ctx, "testfile", strings.NewReader("data"))
		require.Error(t, err)
		var target *backlog.APIResponseError
		assert.True(t, errors.As(err, &target))
	})
}
