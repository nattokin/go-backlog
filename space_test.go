package backlog_test

import (
	"bytes"
	"context"
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
}
