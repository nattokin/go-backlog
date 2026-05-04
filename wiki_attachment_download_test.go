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
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestWikiAttachmentService_Download(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"Download": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/wikis/1/attachments/20", req.URL.Path)
				return mock.NewBinaryResponse("doc.pdf", "application/pdf", []byte("PDF")), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Wiki.Attachment.Download(ctx, 1, 20)
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.Equal(t, "doc.pdf", got.Filename)
				assert.Equal(t, "application/pdf", got.ContentType)
				body, err := io.ReadAll(got.Body)
				require.NoError(t, err)
				assert.Equal(t, []byte("PDF"), body)
				got.Body.Close()
			},
		},
		"Download/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Wiki.Attachment.Download(ctx, 1, 20)
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
