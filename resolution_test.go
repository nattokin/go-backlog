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

func TestResolutionService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/resolutions", req.URL.Path)
				return mock.NewResponse(fixture.Resolution.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Resolution.List(ctx)
				require.NoError(t, err)
				assert.Equal(t, fixture.Resolution.List, got)
			},
		},
		"List/error": {
			doFunc: mock.NewInternalServerErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Resolution.List(ctx)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "apiKey", backlog.WithDoer(&mock.Doer{DoFunc: tc.doFunc}))
			require.NoError(t, err)
			tc.call(t, c)
		})
	}
}
