package backlog_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

func TestPullRequestAttachmentService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests/1/attachments", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Attachment.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Attachment.List(ctx, "TEST", "repo", 1)
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 2, got[0].ID)
				assert.Equal(t, 5, got[1].ID)
			},
		},
		"List/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such repository.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.PullRequest.Attachment.List(ctx, "TEST", "repo", 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Remove": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests/1/attachments/8", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Attachment.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.PullRequest.Attachment.Remove(ctx, "TEST", "repo", 1, 8)
				require.NoError(t, err)
				assert.Equal(t, 8, got.ID)
			},
		},
		"Remove/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such attachment.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.PullRequest.Attachment.Remove(ctx, "TEST", "repo", 1, 8)
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
