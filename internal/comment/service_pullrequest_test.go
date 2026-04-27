package comment_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/comment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestPullRequestCommentService_All(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		prNumber       int
		opts           []core.RequestOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantIDs     []int
	}{
		"success-no-options": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo/pullRequests/1/comments", spath)
				return mock.NewJSONResponse(fixture.Comment.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-count-and-order": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			opts: []core.RequestOption{
				o.WithCount(20),
				o.WithOrder("asc"),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo/pullRequests/1/comments", spath)
				assert.Equal(t, "20", query.Get("count"))
				assert.Equal(t, "asc", query.Get("order"))
				return mock.NewJSONResponse(fixture.Comment.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-minID-maxID": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			opts: []core.RequestOption{
				o.WithMinID(10),
				o.WithMaxID(100),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo/pullRequests/1/comments", spath)
				assert.Equal(t, "10", query.Get("minId"))
				assert.Equal(t, "100", query.Get("maxId"))
				return mock.NewJSONResponse(fixture.Comment.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"error-empty-projectIDOrKey": {
			projectIDOrKey: "",
			repoIDOrName:   "repo",
			prNumber:       1,
			wantErrType:    &core.ValidationError{},
		},
		"error-zero-projectIDOrKey": {
			projectIDOrKey: "0",
			repoIDOrName:   "repo",
			prNumber:       1,
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-repoIDOrName": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "",
			prNumber:       1,
			wantErrType:    &core.ValidationError{},
		},
		"error-zero-repoIDOrName": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "0",
			prNumber:       1,
			wantErrType:    &core.ValidationError{},
		},
		"error-invalid-prNumber": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       0,
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-client-network": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
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

			s := comment.NewPullRequestService(method)

			got, err := s.All(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.prNumber, tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Len(t, got, len(tc.wantIDs))
			for i := range got {
				assert.Equal(t, tc.wantIDs[i], got[i].ID)
			}
		})
	}
}

func TestPullRequestCommentService_Add(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		prNumber       int
		content        string
		opts           []core.RequestOption

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantID      int
	}{
		"success-required-only": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			content:        "This is a comment.",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo/pullRequests/1/comments", spath)
				assert.Equal(t, "This is a comment.", form.Get("content"))
				return mock.NewCreatedJSONResponse(fixture.Comment.SingleJSON), nil
			},
			wantID: 1,
		},
		"success-with-notifiedUserIDs": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			content:        "Notifying users.",
			opts:           []core.RequestOption{o.WithNotifiedUserIDs([]int{5, 6})},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo/pullRequests/1/comments", spath)
				assert.Equal(t, "Notifying users.", form.Get("content"))
				assert.Equal(t, []string{"5", "6"}, form["notifiedUserId[]"])
				return mock.NewCreatedJSONResponse(fixture.Comment.SingleJSON), nil
			},
			wantID: 1,
		},
		"error-empty-projectIDOrKey": {
			projectIDOrKey: "",
			repoIDOrName:   "repo",
			prNumber:       1,
			content:        "x",
			wantErrType:    &core.ValidationError{},
		},
		"error-zero-projectIDOrKey": {
			projectIDOrKey: "0",
			repoIDOrName:   "repo",
			prNumber:       1,
			content:        "x",
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-repoIDOrName": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "",
			prNumber:       1,
			content:        "x",
			wantErrType:    &core.ValidationError{},
		},
		"error-zero-repoIDOrName": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "0",
			prNumber:       1,
			content:        "x",
			wantErrType:    &core.ValidationError{},
		},
		"error-invalid-prNumber": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       0,
			content:        "x",
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-content": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			content:        "",
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			content:        "x",
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-client-network": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			content:        "x",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			content:        "x",
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

			s := comment.NewPullRequestService(method)

			got, err := s.Add(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.prNumber, tc.content, tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestPullRequestCommentService_Count(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		prNumber       int

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantCount   int
	}{
		"success": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo/pullRequests/1/comments/count", spath)
				return mock.NewJSONResponse(`{"count":7}`), nil
			},
			wantCount: 7,
		},
		"error-empty-projectIDOrKey": {
			projectIDOrKey: "",
			repoIDOrName:   "repo",
			prNumber:       1,
			wantErrType:    &core.ValidationError{},
		},
		"error-zero-projectIDOrKey": {
			projectIDOrKey: "0",
			repoIDOrName:   "repo",
			prNumber:       1,
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-repoIDOrName": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "",
			prNumber:       1,
			wantErrType:    &core.ValidationError{},
		},
		"error-zero-repoIDOrName": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "0",
			prNumber:       1,
			wantErrType:    &core.ValidationError{},
		},
		"error-invalid-prNumber": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       0,
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
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

			s := comment.NewPullRequestService(method)

			count, err := s.Count(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.prNumber)

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

func TestPullRequestCommentService_Update(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		prNumber       int
		commentID      int
		content        string

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantID      int
	}{
		"success": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			commentID:      42,
			content:        "Updated content.",
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo/pullRequests/1/comments/42", spath)
				assert.Equal(t, "Updated content.", form.Get("content"))
				return mock.NewJSONResponse(fixture.Comment.SingleJSON), nil
			},
			wantID: 1,
		},
		"error-empty-projectIDOrKey": {
			projectIDOrKey: "",
			repoIDOrName:   "repo",
			prNumber:       1,
			commentID:      1,
			content:        "x",
			wantErrType:    &core.ValidationError{},
		},
		"error-zero-projectIDOrKey": {
			projectIDOrKey: "0",
			repoIDOrName:   "repo",
			prNumber:       1,
			commentID:      1,
			content:        "x",
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-repoIDOrName": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "",
			prNumber:       1,
			commentID:      1,
			content:        "x",
			wantErrType:    &core.ValidationError{},
		},
		"error-zero-repoIDOrName": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "0",
			prNumber:       1,
			commentID:      1,
			content:        "x",
			wantErrType:    &core.ValidationError{},
		},
		"error-invalid-prNumber": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       0,
			commentID:      1,
			content:        "x",
			wantErrType:    &core.ValidationError{},
		},
		"error-invalid-commentID": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			commentID:      0,
			content:        "x",
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-content": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			commentID:      1,
			content:        "",
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			commentID:      42,
			content:        "x",
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			commentID:      42,
			content:        "x",
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

			s := comment.NewPullRequestService(method)

			got, err := s.Update(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.prNumber, tc.commentID, tc.content)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}
