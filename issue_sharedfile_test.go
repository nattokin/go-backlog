package backlog_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestIssueSharedFileService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues/TEST-1/sharedFiles", req.URL.Path)
				return mock.NewJSONResponse(fixture.SharedFile.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.SharedFile.List(ctx, "TEST-1")
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 454403, got[0].ID)
				assert.Equal(t, 454404, got[1].ID)
			},
		},
		"List/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.SharedFile.List(ctx, "TEST-1")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Link": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/issues/TEST-1/sharedFiles", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, []string{"454403"}, req.PostForm["fileId[]"])
				return mock.NewJSONResponse(fixture.SharedFile.SingleListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.SharedFile.Link(ctx, "TEST-1", []int{454403})
				require.NoError(t, err)
				assert.Len(t, got, 1)
				assert.Equal(t, 454403, got[0].ID)
			},
		},
		"Link/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.SharedFile.Link(ctx, "TEST-1", []int{454403})
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Unlink": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/issues/TEST-1/sharedFiles/454403", req.URL.Path)
				return mock.NewJSONResponse(fixture.SharedFile.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.SharedFile.Unlink(ctx, "TEST-1", 454403)
				require.NoError(t, err)
				assert.Equal(t, 454403, got.ID)
				assert.Equal(t, "01_buz.png", got.Name)
			},
		},
		"Unlink/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.SharedFile.Unlink(ctx, "TEST-1", 454403)
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
