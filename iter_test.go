package backlog_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

// pullRequestLastPageJSON is a single-element array used to simulate the last page of a paginated response.
const pullRequestLastPageJSON = `[{"id":4,"projectId":3,"repositoryId":5,"number":3,"summary":"last PR","description":"","base":"main","branch":"feature/baz","status":{"id":1,"name":"Open"},"assignee":null,"issue":null,"baseCommit":null,"branchCommit":null,"mergeCommit":null,"closeAt":null,"mergeAt":null,"createdUser":{"id":1,"userId":"admin","name":"admin","roleType":1,"lang":"ja","mailAddress":"admin@example.com"},"created":"2024-01-12T10:00:00Z","updatedUser":{"id":1,"userId":"admin","name":"admin","roleType":1,"lang":"ja","mailAddress":"admin@example.com"},"updated":"2024-01-12T10:00:00Z","attachments":[],"stars":[]}]`

func TestPullRequestServiceAll(t *testing.T) {
	ctx := context.Background()

	t.Run("All", func(t *testing.T) {
		t.Parallel()

		// page 1: full page (perPage=2), page 2: 1 item (signals last page)
		var callCount atomic.Int32
		c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{
			do: func(req *http.Request) (*http.Response, error) {
				n := callCount.Add(1)
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/projects/TEST/git/repositories/repo/pullRequests", req.URL.Path)
				assert.Equal(t, "2", req.URL.Query().Get("count"))
				switch n {
				case 1:
					assert.Equal(t, "0", req.URL.Query().Get("offset"))
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(fixture.PullRequest.ListJSON)),
					}, nil
				case 2:
					assert.Equal(t, "2", req.URL.Query().Get("offset"))
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(pullRequestLastPageJSON)),
					}, nil
				default:
					t.Errorf("unexpected request #%d", n)
					return nil, nil
				}
			},
		}))
		require.NoError(t, err)

		var got []*backlog.PullRequest
		for pr, err := range c.PullRequest.All(ctx, "TEST", "repo", 2) {
			require.NoError(t, err)
			got = append(got, pr)
		}

		assert.Equal(t, int32(2), callCount.Load())
		assert.Len(t, got, 3)
		assert.Equal(t, 2, got[0].ID)
		assert.Equal(t, 3, got[1].ID)
		assert.Equal(t, 4, got[2].ID)
	})

	t.Run("All/break", func(t *testing.T) {
		t.Parallel()

		var callCount atomic.Int32
		c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{
			do: func(req *http.Request) (*http.Response, error) {
				n := callCount.Add(1)
				if n > 1 {
					t.Errorf("unexpected request #%d after break", n)
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(fixture.PullRequest.ListJSON)),
				}, nil
			},
		}))
		require.NoError(t, err)

		var got []*backlog.PullRequest
		for pr, err := range c.PullRequest.All(ctx, "TEST", "repo", 2) {
			require.NoError(t, err)
			got = append(got, pr)
			break
		}

		assert.Equal(t, int32(1), callCount.Load())
		assert.Len(t, got, 1)
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
}
