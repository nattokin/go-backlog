package backlog_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestIssueAttachmentService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues/TEST-1/attachments", req.URL.Path)
				return mock.NewJSONResponse(fixture.Attachment.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Attachment.List(ctx, "TEST-1")
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 2, got[0].ID)
				assert.Equal(t, 5, got[1].ID)
			},
		},
		"List/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Attachment.List(ctx, "TEST-1")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Remove": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/issues/TEST-1/attachments/8", req.URL.Path)
				return mock.NewJSONResponse(fixture.Attachment.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Attachment.Remove(ctx, "TEST-1", 8)
				require.NoError(t, err)
				assert.Equal(t, 8, got.ID)
			},
		},
		"Remove/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Attachment.Remove(ctx, "TEST-1", 8)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Download": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues/TEST-1/attachments/10", req.URL.Path)
				return mock.NewBinaryResponse("file.png", "image/png", []byte("PNG")), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Attachment.Download(ctx, "TEST-1", 10)
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.Equal(t, "file.png", got.Filename)
				assert.Equal(t, "image/png", got.ContentType)
				body, err := io.ReadAll(got.Body)
				require.NoError(t, err)
				assert.Equal(t, []byte("PNG"), body)
				got.Body.Close()
			},
		},
		"Download/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Attachment.Download(ctx, "TEST-1", 10)
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
