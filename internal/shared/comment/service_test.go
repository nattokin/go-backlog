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

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/shared/comment"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
	"github.com/nattokin/go-backlog/internal/validation"
)

func TestService_ApplyListOptions(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		opts []*option.APIParamOption

		wantQuery              url.Values
		wantErrType            error
		wantValidationErrCount int
	}{
		"success-no-options": {
			wantQuery: url.Values{},
		},
		"success-all-options": {
			opts: []*option.APIParamOption{
				o.WithMinID(10),
				o.WithMaxID(20),
				o.WithCount(30),
				o.WithOrder("desc"),
			},
			wantQuery: url.Values{
				"minId": {"10"},
				"maxId": {"20"},
				"count": {"30"},
				"order": {"desc"},
			},
		},
		"error-validation-count": {
			opts:                   []*option.APIParamOption{o.WithCount(0)},
			wantValidationErrCount: 1,
		},
		"error-option-not-allowed": {
			opts:        []*option.APIParamOption{o.WithName("name")},
			wantErrType: &option.InvalidOptionKeyError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := comment.NewService(mock.NewMethod(t))
			query := url.Values{}
			err := s.ApplyListOptions(query, tc.opts...)

			if tc.wantValidationErrCount > 0 {
				var ves validation.Errors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

func TestService_FetchList(t *testing.T) {
	query := url.Values{"count": {"10"}}

	cases := map[string]struct {
		mockGetFn func(context.Context, string, url.Values) (*http.Response, error)

		wantFirstID int
		wantErrType error
	}{
		"success": {
			mockGetFn: func(ctx context.Context, spath string, q url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/comments", spath)
				assert.Equal(t, query, q)
				return mock.NewResponse(fixture.Comment.ListJSON), nil
			},
			wantFirstID: 1,
		},
		"error-client-network": {
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, errors.New("network")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, &client.APIResponseError{}
			},
			wantErrType: &client.APIResponseError{},
		},
		"error-response-invalid-json": {
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := mock.NewMethod(t)
			m.Get = tc.mockGetFn
			s := comment.NewService(m)
			got, err := s.FetchList(context.Background(), "issues/TEST-1/comments", query)

			if tc.wantErrType != nil {
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, got)
			assert.Equal(t, tc.wantFirstID, got[0].ID)
		})
	}
}

func TestService_ApplyAddOptions(t *testing.T) {
	o := &option.OptionService{}
	validTypes := []option.APIParamOptionType{option.ParamContent, option.ParamNotifiedUserIDs}

	cases := map[string]struct {
		content string
		opts    []*option.APIParamOption

		wantForm               url.Values
		wantErrType            error
		wantValidationErrCount int
	}{
		"success-content-only": {
			content:  "comment",
			wantForm: url.Values{"content": {"comment"}},
		},
		"success-with-options": {
			content: "comment",
			opts:    []*option.APIParamOption{o.WithNotifiedUserIDs([]int{1, 2})},
			wantForm: url.Values{
				"content":          {"comment"},
				"notifiedUserId[]": {"1", "2"},
			},
		},
		"error-validation-content-empty": {
			content:                "",
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			content:                "",
			opts:                   []*option.APIParamOption{o.WithNotifiedUserIDs([]int{0})},
			wantValidationErrCount: 2,
		},
		"error-option-not-allowed": {
			content:     "comment",
			opts:        []*option.APIParamOption{o.WithName("name")},
			wantErrType: &option.InvalidOptionKeyError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := comment.NewService(mock.NewMethod(t))
			form := url.Values{}
			err := s.ApplyAddOptions(form, validTypes, tc.content, tc.opts...)

			if tc.wantValidationErrCount > 0 {
				var ves validation.Errors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantForm, form)
		})
	}
}

func TestService_FetchAdd(t *testing.T) {
	form := url.Values{"content": {"comment"}}

	cases := map[string]struct {
		mockPostFn func(context.Context, string, url.Values) (*http.Response, error)

		wantID      int
		wantErrType error
	}{
		"success": {
			mockPostFn: func(ctx context.Context, spath string, f url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/comments", spath)
				assert.Equal(t, form, f)
				return mock.NewResponse(fixture.Comment.SingleJSON), nil
			},
			wantID: 1,
		},
		"error-client-network": {
			mockPostFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, errors.New("network")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			mockPostFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, &client.APIResponseError{}
			},
			wantErrType: &client.APIResponseError{},
		},
		"error-response-invalid-json": {
			mockPostFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := mock.NewMethod(t)
			m.Post = tc.mockPostFn
			s := comment.NewService(m)
			got, err := s.FetchAdd(context.Background(), "issues/TEST-1/comments", form)

			if tc.wantErrType != nil {
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestService_Count(t *testing.T) {
	cases := map[string]struct {
		mockGetFn func(context.Context, string, url.Values) (*http.Response, error)

		wantCount   int
		wantErrType error
	}{
		"success": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/comments/count", spath)
				assert.Nil(t, query)
				return mock.NewResponse(`{"count": 10}`), nil
			},
			wantCount: 10,
		},
		"error-client-network": {
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, errors.New("network")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, &client.APIResponseError{}
			},
			wantErrType: &client.APIResponseError{},
		},
		"error-response-invalid-json": {
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := mock.NewMethod(t)
			m.Get = tc.mockGetFn
			s := comment.NewService(m)
			got, err := s.Count(context.Background(), "issues/TEST-1/comments/count")

			if tc.wantErrType != nil {
				assert.Zero(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantCount, got)
		})
	}
}

func TestService_One(t *testing.T) {
	cases := map[string]struct {
		mockGetFn func(context.Context, string, url.Values) (*http.Response, error)

		wantID      int
		wantErrType error
	}{
		"success": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/comments/1", spath)
				assert.Nil(t, query)
				return mock.NewResponse(fixture.Comment.SingleJSON), nil
			},
			wantID: 1,
		},
		"error-client-network": {
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, errors.New("network")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, &client.APIResponseError{}
			},
			wantErrType: &client.APIResponseError{},
		},
		"error-response-invalid-json": {
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := mock.NewMethod(t)
			m.Get = tc.mockGetFn
			s := comment.NewService(m)
			got, err := s.One(context.Background(), "issues/TEST-1/comments/1")

			if tc.wantErrType != nil {
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestService_ApplyUpdateOptions(t *testing.T) {
	cases := map[string]struct {
		content string

		wantForm   url.Values
		wantErrMsg string
	}{
		"success": {
			content:  "updated",
			wantForm: url.Values{"content": {"updated"}},
		},
		"error-validation-content-empty": {
			content:    "",
			wantErrMsg: "invalid content: must not be empty",
		},
		"error-validation-content-whitespace": {
			content:    " \t\n ",
			wantErrMsg: "invalid content: must not be empty",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := comment.NewService(mock.NewMethod(t))
			form := url.Values{}
			ve := s.ApplyUpdateOptions(form, tc.content)

			if tc.wantErrMsg != "" {
				if assert.NotNil(t, ve) {
					assert.Equal(t, "content", ve.Target())
					assert.Equal(t, tc.wantErrMsg, ve.Message())
				}
				assert.Empty(t, form)
				return
			}

			assert.Nil(t, ve)
			assert.Equal(t, tc.wantForm, form)
		})
	}
}

func TestService_FetchUpdate(t *testing.T) {
	form := url.Values{"content": {"updated"}}

	cases := map[string]struct {
		mockPatchFn func(context.Context, string, url.Values) (*http.Response, error)

		wantID      int
		wantErrType error
	}{
		"success": {
			mockPatchFn: func(ctx context.Context, spath string, f url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/comments/1", spath)
				assert.Equal(t, form, f)
				return mock.NewResponse(fixture.Comment.SingleJSON), nil
			},
			wantID: 1,
		},
		"error-client-network": {
			mockPatchFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, errors.New("network")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			mockPatchFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, &client.APIResponseError{}
			},
			wantErrType: &client.APIResponseError{},
		},
		"error-response-invalid-json": {
			mockPatchFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := mock.NewMethod(t)
			m.Patch = tc.mockPatchFn
			s := comment.NewService(m)
			got, err := s.FetchUpdate(context.Background(), "issues/TEST-1/comments/1", form)

			if tc.wantErrType != nil {
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}
