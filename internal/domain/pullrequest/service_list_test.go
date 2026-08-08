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
		opts           []*core.APIParamOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
		wantNumbers            []int
	}{
		"success-no-options": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo1/pullRequests", spath)
				return mock.NewResponse(fixture.PullRequest.ListJSON), nil
			},
			wantNumbers: []int{1, 2},
		},
		"success-with-all-options": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			opts: []*core.APIParamOption{
				o.WithStatusIDs([]int{1, 2}),
				o.WithAssigneeIDs([]int{10}),
				o.WithIssueIDs([]int{100}),
				o.WithCreatedUserIDs([]int{11}),
				o.WithOffset(5),
				o.WithCount(50),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"1", "2"}, query["statusId[]"])
				assert.Equal(t, []string{"10"}, query["assigneeId[]"])
				assert.Equal(t, []string{"100"}, query["issueId[]"])
				assert.Equal(t, []string{"11"}, query["createdUserId[]"])
				assert.Equal(t, "5", query.Get("offset"))
				assert.Equal(t, "50", query.Get("count"))
				return mock.NewResponse(fixture.PullRequest.ListJSON), nil
			},
			wantNumbers: []int{1, 2},
		},
		"success-with-statusIDs": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			opts:           []*core.APIParamOption{o.WithStatusIDs([]int{1, 2})},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"1", "2"}, query["statusId[]"])
				return mock.NewResponse(fixture.PullRequest.ListJSON), nil
			},
			wantNumbers: []int{1, 2},
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			repoIDOrName:           "repo1",
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			repoIDOrName:           "repo1",
			wantValidationErrCount: 1,
		},
		"error-validation-repoIDOrName-empty": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "",
			wantValidationErrCount: 1,
		},
		"error-validation-repoIDOrName-zero": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "0",
			wantValidationErrCount: 1,
		},
		"error-validation-opt-invalid-statusID": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo1",
			opts:                   []*core.APIParamOption{o.WithStatusIDs([]int{0})},
			wantValidationErrCount: 1,
		},
		"error-validation-opt-invalid-count": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo1",
			opts:                   []*core.APIParamOption{o.WithCount(0)},
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			repoIDOrName:           "",
			opts:                   []*core.APIParamOption{o.WithCount(0)},
			wantValidationErrCount: 3,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo1",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			repoIDOrName:           "",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-option-set-failed": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			opts:           []*core.APIParamOption{mock.NewFailingSetOption(core.ParamOffset)},
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
				return mock.NewResponse(fixture.InvalidJSON), nil
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

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, prs)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, prs)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, prs)
				assert.ErrorAs(t, err, &tc.wantErrType)
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
		var ves core.ValidationErrors
		assert.ErrorAs(t, err, &ves)
	})

	t.Run("error-invalid-repo", func(t *testing.T) {
		t.Parallel()

		s := pullrequest.NewService(mock.NewMethod(t))
		_, err := s.All(ctx, 10, "PRJ", "")
		require.Error(t, err)
		var ves core.ValidationErrors
		assert.ErrorAs(t, err, &ves)
	})

	t.Run("error-invalid-count", func(t *testing.T) {
		t.Parallel()

		s := pullrequest.NewService(mock.NewMethod(t))
		_, err := s.All(ctx, 0, "PRJ", "repo1")
		require.Error(t, err)
		var ves core.ValidationErrors
		assert.ErrorAs(t, err, &ves)
	})

	t.Run("error-invalid-option", func(t *testing.T) {
		t.Parallel()

		s := pullrequest.NewService(mock.NewMethod(t))
		_, err := s.All(ctx, 10, "PRJ", "repo1", mock.NewInvalidTypeOption())
		require.Error(t, err)
		var target *core.InvalidOptionKeyError
		assert.ErrorAs(t, err, &target)
	})

	t.Run("error-offset-passed-to-all", func(t *testing.T) {
		t.Parallel()

		o := &core.OptionService{}
		s := pullrequest.NewService(mock.NewMethod(t))
		_, err := s.All(ctx, 10, "PRJ", "repo1", o.WithOffset(5))
		require.Error(t, err)
		var target *core.InvalidOptionKeyError
		assert.ErrorAs(t, err, &target)
	})

	t.Run("error-count-passed-to-all", func(t *testing.T) {
		t.Parallel()

		o := &core.OptionService{}
		s := pullrequest.NewService(mock.NewMethod(t))
		_, err := s.All(ctx, 10, "PRJ", "repo1", o.WithCount(50))
		require.Error(t, err)
		var target *core.InvalidOptionKeyError
		assert.ErrorAs(t, err, &target)
	})
}
