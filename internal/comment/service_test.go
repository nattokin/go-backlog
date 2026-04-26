package comment_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
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

func TestIssueCommentService_All(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		issueIDOrKey string
		opts         []core.RequestOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantIDs     []int
	}{
		"success-no-options": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/comments", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Comment.ListJSON))),
				}, nil
			},
			wantIDs: []int{1, 2},
		},
		"success-by-numeric-id": {
			issueIDOrKey: "1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1/comments", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Comment.ListJSON))),
				}, nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-count-and-order": {
			issueIDOrKey: "PRJ-1",
			opts: []core.RequestOption{
				o.WithCount(20),
				o.WithOrder("asc"),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/comments", spath)
				assert.Equal(t, "20", query.Get("count"))
				assert.Equal(t, "asc", query.Get("order"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Comment.ListJSON))),
				}, nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-minID-maxID": {
			issueIDOrKey: "PRJ-1",
			opts: []core.RequestOption{
				o.WithMinID(10),
				o.WithMaxID(100),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/comments", spath)
				assert.Equal(t, "10", query.Get("minId"))
				assert.Equal(t, "100", query.Get("maxId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Comment.ListJSON))),
				}, nil
			},
			wantIDs: []int{1, 2},
		},
		"error-empty-issueIDOrKey": {
			issueIDOrKey: "",
			wantErrType:  &core.ValidationError{},
		},
		"error-zero-issueIDOrKey": {
			issueIDOrKey: "0",
			wantErrType:  &core.ValidationError{},
		},
		"error-option-invalid-type": {
			issueIDOrKey: "PRJ-1",
			opts:         []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:  &core.InvalidOptionKeyError{},
		},
		"error-client-network": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
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

			s := comment.NewIssueService(method)

			got, err := s.All(context.Background(), tc.issueIDOrKey, tc.opts...)

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

func TestIssueCommentService_Add(t *testing.T) {
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
				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Comment.SingleJSON))),
				}, nil
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
				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Comment.SingleJSON))),
				}, nil
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
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			content:      "x",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
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

			s := comment.NewIssueService(method)

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

func TestIssueCommentService_Count(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantCount   int
	}{
		"success": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/comments/count", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"count":7}`))),
				}, nil
			},
			wantCount: 7,
		},
		"error-empty-issueIDOrKey": {
			issueIDOrKey: "",
			wantErrType:  &core.ValidationError{},
		},
		"error-zero-issueIDOrKey": {
			issueIDOrKey: "0",
			wantErrType:  &core.ValidationError{},
		},
		"error-client-network": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
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

			s := comment.NewIssueService(method)

			count, err := s.Count(context.Background(), tc.issueIDOrKey)

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

func TestIssueCommentService_One(t *testing.T) {
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Comment.SingleJSON))),
				}, nil
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
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
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

			s := comment.NewIssueService(method)

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

func TestIssueCommentService_Delete(t *testing.T) {
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Comment.SingleJSON))),
				}, nil
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
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
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

			s := comment.NewIssueService(method)

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

func TestIssueCommentService_Update(t *testing.T) {
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Comment.SingleJSON))),
				}, nil
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
		"error-client-network": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			content:      "x",
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			content:      "x",
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
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

			s := comment.NewIssueService(method)

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

func TestIssueCommentService_Notifications(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string
		commentID    int

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantLen     int
	}{
		"success": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/comments/42/notifications", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`[{"id":1},{"id":2}]`))),
				}, nil
			},
			wantLen: 2,
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
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			mockGetFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
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

			s := comment.NewIssueService(method)

			got, err := s.Notifications(context.Background(), tc.issueIDOrKey, tc.commentID)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, got, tc.wantLen)
		})
	}
}

func TestIssueCommentService_Notify(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string
		commentID    int
		userIDs      []int

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantID      int
	}{
		"success": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			userIDs:      []int{5, 6},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/comments/42/notifications", spath)
				assert.Equal(t, []string{"5", "6"}, form["notifiedUserId[]"])
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Comment.SingleJSON))),
				}, nil
			},
			wantID: 1,
		},
		"error-empty-issueIDOrKey": {
			issueIDOrKey: "",
			commentID:    1,
			userIDs:      []int{5},
			wantErrType:  &core.ValidationError{},
		},
		"error-invalid-commentID": {
			issueIDOrKey: "PRJ-1",
			commentID:    0,
			userIDs:      []int{5},
			wantErrType:  &core.ValidationError{},
		},
		"error-invalid-userID": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			userIDs:      []int{0},
			wantErrType:  &core.ValidationError{},
		},
		"error-client-network": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			userIDs:      []int{5},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			userIDs:      []int{5},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
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

			s := comment.NewIssueService(method)

			got, err := s.Notify(context.Background(), tc.issueIDOrKey, tc.commentID, tc.userIDs)

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
