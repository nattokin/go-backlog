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

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/issue"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestCommentService_List(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		issueIDOrKey string
		opts         []*option.APIParamOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
		wantIDs                []int
	}{
		"success-no-options": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/comments", spath)
				return mock.NewResponse(fixture.Comment.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"success-by-numeric-id": {
			issueIDOrKey: "1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1/comments", spath)
				return mock.NewResponse(fixture.Comment.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-count-and-order": {
			issueIDOrKey: "PRJ-1",
			opts: []*option.APIParamOption{
				o.WithCount(20),
				o.WithOrder("asc"),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "20", query.Get("count"))
				assert.Equal(t, "asc", query.Get("order"))
				return mock.NewResponse(fixture.Comment.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-minID-maxID": {
			issueIDOrKey: "PRJ-1",
			opts: []*option.APIParamOption{
				o.WithMinID(10),
				o.WithMaxID(100),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "10", query.Get("minId"))
				assert.Equal(t, "100", query.Get("maxId"))
				return mock.NewResponse(fixture.Comment.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},

		"error-validation-issueIDOrKey-empty": {
			issueIDOrKey:           "",
			wantValidationErrCount: 1,
		},
		"error-validation-issueIDOrKey-zero": {
			issueIDOrKey:           "0",
			wantValidationErrCount: 1,
		},
		"error-validation-opt-count-zero": {
			issueIDOrKey:           "PRJ-1",
			opts:                   []*option.APIParamOption{o.WithCount(0)},
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			issueIDOrKey:           "",
			opts:                   []*option.APIParamOption{o.WithCount(0)},
			wantValidationErrCount: 2,
		},

		"error-nil-option-with-valid-values": {
			issueIDOrKey:           "PRJ-1",
			opts:                   []*option.APIParamOption{o.WithCount(10), nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			issueIDOrKey:           "",
			opts:                   []*option.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		"error-option-invalid-type": {
			issueIDOrKey: "PRJ-1",
			opts:         []*option.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:  &option.InvalidOptionKeyError{},
		},

		"error-client-network": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, &client.APIResponseError{}
			},
			wantErrType: &client.APIResponseError{},
		},
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
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
			s := issue.NewCommentService(method)
			got, err := s.List(context.Background(), tc.issueIDOrKey, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, got)
				var target *option.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
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

func TestCommentService_Count(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantCount              int
	}{
		"success": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/comments/count", spath)
				return mock.NewResponse(`{"count":7}`), nil
			},
			wantCount: 7,
		},

		"error-validation-issueIDOrKey-empty": {
			issueIDOrKey:           "",
			wantValidationErrCount: 1,
		},
		"error-validation-issueIDOrKey-zero": {
			issueIDOrKey:           "0",
			wantValidationErrCount: 1,
		},

		"error-client-network": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, &client.APIResponseError{}
			},
			wantErrType: &client.APIResponseError{},
		},
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
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
			s := issue.NewCommentService(method)
			count, err := s.Count(context.Background(), tc.issueIDOrKey)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Zero(t, count)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Zero(t, count)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantCount, count)
		})
	}
}

func TestCommentService_Notifications(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string
		commentID    int

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantLen                int
	}{
		"success": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/comments/42/notifications", spath)
				return mock.NewResponse(`[{"id":1},{"id":2}]`), nil
			},
			wantLen: 2,
		},

		"error-validation-issueIDOrKey-empty": {
			issueIDOrKey:           "",
			commentID:              1,
			wantValidationErrCount: 1,
		},
		"error-validation-commentID-zero": {
			issueIDOrKey:           "PRJ-1",
			commentID:              0,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			issueIDOrKey:           "",
			commentID:              0,
			wantValidationErrCount: 2,
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
				return nil, &client.APIResponseError{}
			},
			wantErrType: &client.APIResponseError{},
		},
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			mockGetFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
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
			s := issue.NewCommentService(method)
			got, err := s.Notifications(context.Background(), tc.issueIDOrKey, tc.commentID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			assert.Len(t, got, tc.wantLen)
		})
	}
}

func TestCommentService_Notify(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string
		commentID    int
		userIDs      []int

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantID                 int
	}{
		"success": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			userIDs:      []int{5, 6},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/comments/42/notifications", spath)
				assert.Equal(t, []string{"5", "6"}, form["notifiedUserId[]"])
				return mock.NewResponse(fixture.Comment.SingleJSON), nil
			},
			wantID: 1,
		},

		"error-validation-issueIDOrKey-empty": {
			issueIDOrKey:           "",
			commentID:              1,
			userIDs:                []int{5},
			wantValidationErrCount: 1,
		},
		"error-validation-commentID-zero": {
			issueIDOrKey:           "PRJ-1",
			commentID:              0,
			userIDs:                []int{5},
			wantValidationErrCount: 1,
		},
		"error-validation-userID-zero": {
			issueIDOrKey:           "PRJ-1",
			commentID:              42,
			userIDs:                []int{0},
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			issueIDOrKey:           "",
			commentID:              0,
			userIDs:                []int{0},
			wantValidationErrCount: 3,
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
		"error-client-api-error": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			userIDs:      []int{5},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, &client.APIResponseError{}
			},
			wantErrType: &client.APIResponseError{},
		},
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			commentID:    42,
			userIDs:      []int{5},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
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
			got, err := s.Notify(context.Background(), tc.issueIDOrKey, tc.commentID, tc.userIDs)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}
