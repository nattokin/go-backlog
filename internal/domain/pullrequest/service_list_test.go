package pullrequest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/pullrequest"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestService_List(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		opts           []core.RequestOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantNumbers []int
	}{
		"success-no-options": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo1/pullRequests", spath)
				return mock.NewJSONResponse(fixture.PullRequest.ListJSON), nil
			},
			wantNumbers: []int{1, 2},
		},
		"success-with-statusIDs": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			opts:           []core.RequestOption{o.WithStatusIDs([]int{1, 2})},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"1", "2"}, query["statusId[]"])
				return mock.NewJSONResponse(fixture.PullRequest.ListJSON), nil
			},
			wantNumbers: []int{1, 2},
		},
		"success-with-count-and-offset": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			opts: []core.RequestOption{
				o.WithCount(50),
				o.WithOffset(10),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "50", query.Get("count"))
				assert.Equal(t, "10", query.Get("offset"))
				return mock.NewJSONResponse(fixture.PullRequest.ListJSON), nil
			},
			wantNumbers: []int{1, 2},
		},
		"error-empty-projectIDOrKey": {
			projectIDOrKey: "",
			repoIDOrName:   "repo1",
			wantErrType:    &core.ValidationError{},
		},
		"error-zero-projectIDOrKey": {
			projectIDOrKey: "0",
			repoIDOrName:   "repo1",
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-repoIDOrName": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "",
			wantErrType:    &core.ValidationError{},
		},
		"error-zero-repoIDOrName": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "0",
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-statusID": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			opts:           []core.RequestOption{o.WithStatusIDs([]int{0})},
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-count": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			opts:           []core.RequestOption{o.WithCount(0)},
			wantErrType:    &core.ValidationError{},
		},
		"error-option-set-failed": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			opts:           []core.RequestOption{mock.NewFailingSetOption(core.ParamOffset)},
			wantErrType:    errors.New(""),
		},
		"error-client-network": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}

			s := pullrequest.NewService(method)
			prs, err := s.List(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, prs)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, prs)
			assert.Len(t, prs, len(tc.wantNumbers))
			for i := range prs {
				assert.Equal(t, tc.wantNumbers[i], prs[i].Number)
			}
		})
	}
}

func TestService_All(t *testing.T) {
	ctx := context.Background()

	t.Run("multi-page", func(t *testing.T) {
		t.Parallel()

		// page 1: full page (perPage=2), page 2: 1 item (signals last page)
		var callCount atomic.Int32
		method := mock.NewMethod(t)
		method.Get = func(_ context.Context, _ string, query url.Values) (*http.Response, error) {
			n := callCount.Add(1)
			assert.Equal(t, "2", query.Get("count"))
			switch n {
			case 1:
				assert.Equal(t, "0", query.Get("offset"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(fixture.PullRequest.ListJSON)),
				}, nil
			case 2:
				assert.Equal(t, "2", query.Get("offset"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(pullRequestLastPageJSON)),
				}, nil
			default:
				t.Errorf("unexpected request #%d", n)
				return nil, nil
			}
		}

		s := pullrequest.NewService(method)
		seq, err := s.All(ctx, 2, "PRJ", "repo1")
		require.NoError(t, err)
		var got []int
		for pr, err := range seq {
			require.NoError(t, err)
			got = append(got, pr.Number)
		}

		assert.Equal(t, int32(2), callCount.Load())
		assert.Equal(t, []int{1, 2, 3}, got)
	})

	t.Run("break", func(t *testing.T) {
		t.Parallel()

		var callCount atomic.Int32
		method := mock.NewMethod(t)
		method.Get = func(_ context.Context, _ string, _ url.Values) (*http.Response, error) {
			n := callCount.Add(1)
			if n > 1 {
				t.Errorf("unexpected request #%d after break", n)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(fixture.PullRequest.ListJSON)),
			}, nil
		}

		s := pullrequest.NewService(method)
		seq, err := s.All(ctx, 2, "PRJ", "repo1")
		require.NoError(t, err)
		var got []int
		for pr, err := range seq {
			require.NoError(t, err)
			got = append(got, pr.Number)
			break
		}

		assert.Equal(t, int32(1), callCount.Load())
		assert.Len(t, got, 1)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		method := mock.NewMethod(t)
		method.Get = func(_ context.Context, _ string, _ url.Values) (*http.Response, error) {
			return nil, errors.New("network error")
		}

		s := pullrequest.NewService(method)
		seq, err := s.All(ctx, 10, "PRJ", "repo1")
		require.NoError(t, err)
		for pr, err := range seq {
			assert.Nil(t, pr)
			require.Error(t, err)
			break
		}
	})

	t.Run("error-invalid-project", func(t *testing.T) {
		t.Parallel()

		s := pullrequest.NewService(mock.NewMethod(t))
		_, err := s.All(ctx, 10, "", "repo1")
		require.Error(t, err)
		assert.IsType(t, &core.ValidationError{}, err)
	})

	t.Run("error-invalid-count", func(t *testing.T) {
		t.Parallel()

		s := pullrequest.NewService(mock.NewMethod(t))
		_, err := s.All(ctx, 0, "PRJ", "repo1")
		require.Error(t, err)
		assert.IsType(t, &core.ValidationError{}, err)
	})

	t.Run("error-invalid-option", func(t *testing.T) {
		t.Parallel()

		s := pullrequest.NewService(mock.NewMethod(t))
		_, err := s.All(ctx, 10, "PRJ", "repo1", mock.NewInvalidTypeOption())
		require.Error(t, err)
		assert.IsType(t, &core.InvalidOptionKeyError{}, err)
	})

	t.Run("error-offset-passed-to-all", func(t *testing.T) {
		t.Parallel()

		o := &core.OptionService{}
		s := pullrequest.NewService(mock.NewMethod(t))
		_, err := s.All(ctx, 10, "PRJ", "repo1", o.WithOffset(5))
		require.Error(t, err)
		assert.IsType(t, &core.InvalidOptionKeyError{}, err)
	})

	t.Run("error-count-passed-to-all", func(t *testing.T) {
		t.Parallel()

		o := &core.OptionService{}
		s := pullrequest.NewService(mock.NewMethod(t))
		_, err := s.All(ctx, 10, "PRJ", "repo1", o.WithCount(50))
		require.Error(t, err)
		assert.IsType(t, &core.InvalidOptionKeyError{}, err)
	})
}
