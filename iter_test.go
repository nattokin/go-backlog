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
)

func TestPullRequestServiceAll(t *testing.T) {
	ctx := context.Background()

	t.Run("All", func(t *testing.T) {
		t.Parallel()

		// Single page response: verifies model->root conversion and error propagation via convertError.
		c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(newMockDoer(fixture.PullRequest.ListJSON)))
		require.NoError(t, err)

		var got []*backlog.PullRequest
		for pr, err := range c.PullRequest.All(ctx, "TEST", "repo", 100) {
			require.NoError(t, err)
			got = append(got, pr)
		}

		assert.Len(t, got, 2)
		assert.Equal(t, 2, got[0].ID)
		assert.Equal(t, 3, got[1].ID)
	})

	t.Run("All/error", func(t *testing.T) {
		t.Parallel()

		c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{
			do: newInternalServerErrorDoFunc(),
		}))
		require.NoError(t, err)

		for pr, err := range c.PullRequest.All(ctx, "TEST", "repo", 10) {
			assert.Nil(t, pr)
			require.Error(t, err)
			var target *backlog.APIResponseError
			assert.True(t, errors.As(err, &target))
			break
		}
	})

	t.Run("All/wrong-http-method", func(t *testing.T) {
		t.Parallel()

		c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{
			do: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				return newMockDoer(fixture.PullRequest.ListJSON).do(req)
			},
		}))
		require.NoError(t, err)

		for _, err := range c.PullRequest.All(ctx, "TEST", "repo", 100) {
			require.NoError(t, err)
			break
		}
	})
}
