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

const (
	testProject = "PRJ"
	testRepo    = "repo1"
)

func TestPullRequestService_List(t *testing.T) {
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
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo1/pullRequests", spath)
				return mock.NewJSONResponse(fixture.PullRequest.ListJSON), nil
			},
			wantNumbers: []int{1, 2},
		},
		"success-with-statusIDs": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			opts:           []core.RequestOption{o.WithStatusIDs([]int{1, 2})},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"1", "2"}, query["statusId[]"])
				return mock.NewJSONResponse(fixture.PullRequest.ListJSON), nil
			},
			wantNumbers: []int{1, 2},
		},
		"success-with-count-and-offset": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
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
			repoIDOrName:   testRepo,
			wantErrType:    &core.ValidationError{},
		},
		"error-zero-projectIDOrKey": {
			projectIDOrKey: "0",
			repoIDOrName:   testRepo,
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-repoIDOrName": {
			projectIDOrKey: testProject,
			repoIDOrName:   "",
			wantErrType:    &core.ValidationError{},
		},
		"error-zero-repoIDOrName": {
			projectIDOrKey: testProject,
			repoIDOrName:   "0",
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-statusID": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			opts:           []core.RequestOption{o.WithStatusIDs([]int{0})},
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-count": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			opts:           []core.RequestOption{o.WithCount(0)},
			wantErrType:    &core.ValidationError{},
		},
		"error-option-set-failed": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			opts:           []core.RequestOption{mock.NewFailingSetOption(core.ParamOffset)},
			wantErrType:    errors.New(""),
		},
		"error-client-network": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
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

func TestPullRequestService_All(t *testing.T) {
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
		seq, err := s.All(ctx, 2, testProject, testRepo)
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
		seq, err := s.All(ctx, 2, testProject, testRepo)
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
		seq, err := s.All(ctx, 10, testProject, testRepo)
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
		_, err := s.All(ctx, 10, "", testRepo)
		require.Error(t, err)
		assert.IsType(t, &core.ValidationError{}, err)
	})

	t.Run("error-invalid-count", func(t *testing.T) {
		t.Parallel()

		s := pullrequest.NewService(mock.NewMethod(t))
		_, err := s.All(ctx, 0, testProject, testRepo)
		require.Error(t, err)
		assert.IsType(t, &core.ValidationError{}, err)
	})

	t.Run("error-invalid-option", func(t *testing.T) {
		t.Parallel()

		s := pullrequest.NewService(mock.NewMethod(t))
		_, err := s.All(ctx, 10, testProject, testRepo, mock.NewInvalidTypeOption())
		require.Error(t, err)
		assert.IsType(t, &core.InvalidOptionKeyError{}, err)
	})
}

func TestPullRequestService_Count(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		opts           []core.RequestOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantCount   int
	}{
		"success-no-options": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo1/pullRequests/count", spath)
				return mock.NewJSONResponse(`{"count":5}`), nil
			},
			wantCount: 5,
		},
		"success-with-assigneeIDs": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			opts:           []core.RequestOption{o.WithAssigneeIDs([]int{10, 20})},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"10", "20"}, query["assigneeId[]"])
				return mock.NewJSONResponse(`{"count":2}`), nil
			},
			wantCount: 2,
		},
		"error-empty-projectIDOrKey": {
			projectIDOrKey: "",
			repoIDOrName:   testRepo,
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-repoIDOrName": {
			projectIDOrKey: testProject,
			repoIDOrName:   "",
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-client-network": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
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
			count, err := s.Count(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Zero(t, count)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantCount, count)
		})
	}
}

func TestPullRequestService_One(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		prNumber       int

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantNumber  int
	}{
		"success": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			prNumber:       1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo1/pullRequests/1", spath)
				return mock.NewJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},
		"error-empty-projectIDOrKey": {
			projectIDOrKey: "",
			repoIDOrName:   testRepo,
			prNumber:       1,
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-repoIDOrName": {
			projectIDOrKey: testProject,
			repoIDOrName:   "",
			prNumber:       1,
			wantErrType:    &core.ValidationError{},
		},
		"error-invalid-prNumber": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			prNumber:       0,
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			prNumber:       1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			prNumber:       1,
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
			got, err := s.One(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.prNumber)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantNumber, got.Number)
		})
	}
}

func TestPullRequestService_Create(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		summary        string
		description    string
		base           string
		branch         string
		opts           []core.RequestOption

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantNumber  int
	}{
		"success-required-only": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			summary:        "test PR",
			description:    "test description",
			base:           "main",
			branch:         "feature/foo",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo1/pullRequests", spath)
				assert.Equal(t, "test PR", form.Get("summary"))
				assert.Equal(t, "test description", form.Get("description"))
				assert.Equal(t, "main", form.Get("base"))
				assert.Equal(t, "feature/foo", form.Get("branch"))
				return mock.NewJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},
		"success-with-options": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			summary:        "test PR",
			description:    "test description",
			base:           "main",
			branch:         "feature/foo",
			opts: []core.RequestOption{
				o.WithAssigneeID(5),
				o.WithIssueID(10),
				o.WithNotifiedUserIDs([]int{1, 2}),
			},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "5", form.Get("assigneeId"))
				assert.Equal(t, "10", form.Get("issueId"))
				assert.Equal(t, []string{"1", "2"}, form["notifiedUserId[]"])
				return mock.NewJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},
		"error-empty-projectIDOrKey": {
			projectIDOrKey: "",
			repoIDOrName:   testRepo,
			summary:        "test PR",
			description:    "desc",
			base:           "main",
			branch:         "feature/foo",
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-repoIDOrName": {
			projectIDOrKey: testProject,
			repoIDOrName:   "",
			summary:        "test PR",
			description:    "desc",
			base:           "main",
			branch:         "feature/foo",
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-summary": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			summary:        "",
			description:    "desc",
			base:           "main",
			branch:         "feature/foo",
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-base": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			summary:        "test PR",
			description:    "desc",
			base:           "",
			branch:         "feature/foo",
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-branch": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			summary:        "test PR",
			description:    "desc",
			base:           "main",
			branch:         "",
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			summary:        "test PR",
			description:    "desc",
			base:           "main",
			branch:         "feature/foo",
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-assigneeID": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			summary:        "test PR",
			description:    "desc",
			base:           "main",
			branch:         "feature/foo",
			opts:           []core.RequestOption{o.WithAssigneeID(0)},
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			summary:        "test PR",
			description:    "desc",
			base:           "main",
			branch:         "feature/foo",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			summary:        "test PR",
			description:    "desc",
			base:           "main",
			branch:         "feature/foo",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockPostFn != nil {
				method.Post = tc.mockPostFn
			}

			s := pullrequest.NewService(method)
			got, err := s.Create(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.summary, tc.description, tc.base, tc.branch, tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantNumber, got.Number)
		})
	}
}

func TestPullRequestService_Update(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		prNumber       int
		option         core.RequestOption
		opts           []core.RequestOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantNumber  int
	}{
		"success-summary": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			prNumber:       1,
			option:         o.WithSummary("Updated PR"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo1/pullRequests/1", spath)
				assert.Equal(t, "Updated PR", form.Get("summary"))
				return mock.NewJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},
		"success-with-comment": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			prNumber:       1,
			option:         o.WithSummary("Updated PR"),
			opts:           []core.RequestOption{o.WithComment("looks good")},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "Updated PR", form.Get("summary"))
				assert.Equal(t, "looks good", form.Get("comment"))
				return mock.NewJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},
		"success-with-issueID": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			prNumber:       1,
			option:         o.WithIssueID(42),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "42", form.Get("issueId"))
				return mock.NewJSONResponse(fixture.PullRequest.SingleJSON), nil
			},
			wantNumber: 1,
		},
		"error-empty-projectIDOrKey": {
			projectIDOrKey: "",
			repoIDOrName:   testRepo,
			prNumber:       1,
			option:         o.WithSummary("x"),
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-repoIDOrName": {
			projectIDOrKey: testProject,
			repoIDOrName:   "",
			prNumber:       1,
			option:         o.WithSummary("x"),
			wantErrType:    &core.ValidationError{},
		},
		"error-invalid-prNumber": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			prNumber:       0,
			option:         o.WithSummary("x"),
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			prNumber:       1,
			option:         mock.NewInvalidTypeOption(),
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-assigneeID": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			prNumber:       1,
			option:         o.WithAssigneeID(0),
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			prNumber:       1,
			option:         o.WithSummary("x"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: testProject,
			repoIDOrName:   testRepo,
			prNumber:       1,
			option:         o.WithSummary("x"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockPatchFn != nil {
				method.Patch = tc.mockPatchFn
			}

			s := pullrequest.NewService(method)
			got, err := s.Update(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.prNumber, tc.option, tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantNumber, got.Number)
		})
	}
}

func Test_contextPropagation(t *testing.T) {
	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	makeMockFn := func(t *testing.T) func(context.Context, string, url.Values) (*http.Response, error) {
		return func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
			assert.Same(t, sentinel, got.Value(ctxKey{}))
			return nil, errors.New("stop")
		}
	}

	o := &core.OptionService{}

	cases := []struct {
		name string
		call func(t *testing.T, m *core.Method)
	}{
		{"Service.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := pullrequest.NewService(m)
			s.List(ctx, testProject, testRepo) //nolint:errcheck
		}},
		{"Service.All", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := pullrequest.NewService(m)
			seq, err := s.All(ctx, 10, testProject, testRepo)
			require.NoError(t, err)
			for range seq {
				break
			}
		}},
		{"Service.Count", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := pullrequest.NewService(m)
			s.Count(ctx, testProject, testRepo) //nolint:errcheck
		}},
		{"Service.One", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := pullrequest.NewService(m)
			s.One(ctx, testProject, testRepo, 1) //nolint:errcheck
		}},
		{"Service.Create", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := pullrequest.NewService(m)
			s.Create(ctx, testProject, testRepo, "summary", "desc", "main", "feature/foo") //nolint:errcheck
		}},
		{"Service.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := pullrequest.NewService(m)
			s.Update(ctx, testProject, testRepo, 1, o.WithSummary("x")) //nolint:errcheck
		}},
		{"AttachmentService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := pullrequest.NewAttachmentService(m)
			s.List(ctx, "TEST", "repo", 1) //nolint:errcheck
		}},
		{"AttachmentService.Remove", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := pullrequest.NewAttachmentService(m)
			s.Remove(ctx, "TEST", "repo", 1, 1) //nolint:errcheck
		}},
		{"AttachmentService.Download", func(t *testing.T, m *core.Method) {
			m.Download = makeMockFn(t)
			s := pullrequest.NewAttachmentService(m)
			s.Download(ctx, "TEST", "repo", 1, 1) //nolint:errcheck
		}},
		{"CommentService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := pullrequest.NewCommentService(m)
			s.List(ctx, "PRJ-1", "REPO-1", 1) //nolint:errcheck
		}},
		{"CommentService.Add", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := pullrequest.NewCommentService(m)
			s.Add(ctx, "PRJ-1", "REPO-1", 1, "comment") //nolint:errcheck
		}},
		{"CommentService.Count", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := pullrequest.NewCommentService(m)
			s.Count(ctx, "PRJ-1", "REPO-1", 1) //nolint:errcheck
		}},
		{"CommentService.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := pullrequest.NewCommentService(m)
			s.Update(ctx, "PRJ-1", "REPO-1", 1, 1, "content") //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
