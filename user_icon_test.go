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

func TestUserService_Icon(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"Icon": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/users/1/icon", req.URL.Path)
				return mock.NewBinaryResponse("avatar.png", "image/png", []byte("PNG")), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.User.Icon(ctx, 1)
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.Equal(t, "avatar.png", got.Filename)
				assert.Equal(t, "image/png", got.ContentType)
				body, err := io.ReadAll(got.Body)
				require.NoError(t, err)
				assert.Equal(t, []byte("PNG"), body)
				got.Body.Close()
			},
		},
		"Icon/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.User.Icon(ctx, 1)
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
