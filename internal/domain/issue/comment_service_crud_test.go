package issue_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/issue"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestCommentService_Add(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		issueIDOrKey string
		content      string
		opts         []core.RequestOption

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantID      int
	}{
		"success-required-only": {
			issueIDOrKey: "PRJ-1",
			content:      "This is a comment.",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/comments", spath)
				assert.Equal(t, "This is a comment.", form.Get("content"))
				return mock.NewCreatedJSONResponse(fixture.Comment.SingleJSON), nil
			},
			wantID: 1,
		},
		"success-with-notifiedUserIDs": {
			issueIDOrKey: "PRJ-1",
			content:      "Notifying users.",
			opts:         []core.RequestOption{o.WithNotifiedUserIDs([]int{5, 6})},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/comments", spath)
				assert.Equal(t, "Notifying users.", form.Get("content"))
				assert.Equal(t, []string{"5", "6"}, form["notifiedUserId[]"])
				return mock.NewCreatedJSONResponse(fixture.Comment.SingleJSON), nil
			},
			wantID: 1,
		},
		"error-empty-issueIDOrKey": {
			issueIDOrKey: "",
			content:      "x",
			wantErrType:  &core.ValidationError{},
		},
		"error-zero-issueIDOrKey": {
			issueIDOrKey: "0",
			content:      "x",
			wantErrType:  &core.ValidationError{},
		},
		"error-empty-content": {
			issueIDOrKey: "PRJ-1",
			content:      "",
			wantErrType:  &core.ValidationError{},
		},
		"error-option-invalid-type": {
			issueIDOrKey: "PRJ-1",
			content:      "x",
			opts:         []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:  &core.InvalidOptionKeyError{},
		},
		"error-client-network": {
			issueIDOrKey: "PRJ-1",
			content:      "x",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			issueIDOrKey: "PRJ-1",
			content:      "x",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, &core.APIResponseError{}
			},
			wantErrType: &core.APIResponseError{},
		},
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			content:      "x",
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

			s := issue.NewCommentService(method)

			got, err := s.Add(context.Background(), tc.issueIDOrKey, tc.content, tc.opts...)

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

func TestCommentService_One(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string
		commentID    int

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantID      int
	}{
		"success": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/comments/42", spath)
				return mock.NewJSONResponse(fixture.Comment.SingleJSON), nil
			},
			wantID: 1,
		},
		"error-empty-issueIDOrKey": {
			issueIDOrKey: "",
			commentID:    1,
			wantErrType:  &core.ValidationError{},
		},
		"error-invalid-commentID": {
			issueIDOrKey: "PRJ-1",
			commentID:    0,
			wantErrType:  &core.ValidationError{},
		},
		"error-client-network": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, &core.APIResponseError{}
			},
			wantErrType: &core.APIResponseError{},
		},
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
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

			s := issue.NewCommentService(method)

			got, err := s.One(context.Background(), tc.issueIDOrKey, tc.commentID)

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

func TestCommentService_Delete(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string
		commentID    int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantID      int
	}{
		"success": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/comments/42", spath)
				return mock.NewJSONResponse(fixture.Comment.SingleJSON), nil
			},
			wantID: 1,
		},
		"error-empty-issueIDOrKey": {
			issueIDOrKey: "",
			commentID:    1,
			wantErrType:  &core.ValidationError{},
		},
		"error-invalid-commentID": {
			issueIDOrKey: "PRJ-1",
			commentID:    0,
			wantErrType:  &core.ValidationError{},
		},
		"error-client-network": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, &core.APIResponseError{}
			},
			wantErrType: &core.APIResponseError{},
		},
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockDeleteFn != nil {
				method.Delete = tc.mockDeleteFn
			}

			s := issue.NewCommentService(method)

			got, err := s.Delete(context.Background(), tc.issueIDOrKey, tc.commentID)

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

func TestCommentService_Update(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string
		commentID    int
		content      string

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantID      int
	}{
		"success": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			content:      "Updated content.",
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/comments/42", spath)
				assert.Equal(t, "Updated content.", form.Get("content"))
				return mock.NewJSONResponse(fixture.Comment.SingleJSON), nil
			},
			wantID: 1,
		},
		"error-empty-issueIDOrKey": {
			issueIDOrKey: "",
			commentID:    1,
			content:      "x",
			wantErrType:  &core.ValidationError{},
		},
		"error-invalid-commentID": {
			issueIDOrKey: "PRJ-1",
			commentID:    0,
			content:      "x",
			wantErrType:  &core.ValidationError{},
		},
		"error-empty-comment": {
			issueIDOrKey: "PRJ-1",
			commentID:    1,
			content:      "",
			wantErrType:  &core.ValidationError{},
		},
		"error-client-network": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			content:      "x",
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			content:      "x",
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, &core.APIResponseError{}
			},
			wantErrType: &core.APIResponseError{},
		},
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			content:      "x",
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

			s := issue.NewCommentService(method)

			got, err := s.Update(context.Background(), tc.issueIDOrKey, tc.commentID, tc.content)

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
