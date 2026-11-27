package watching_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/domain/watching"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
	"github.com/nattokin/go-backlog/internal/validation"
)

func TestService_List(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		userID int
		opts   []*option.APIParamOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
		wantIDs                []int
	}{
		"success-no-options": {
			userID: 1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1/watchings", spath)
				assert.Empty(t, query)
				return mock.NewResponse(fixture.Watching.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-options": {
			userID: 2,
			opts: []*option.APIParamOption{
				o.WithOrder("asc"),
				o.WithWatchingSort("created"),
				o.WithCount(10),
				o.WithOffset(20),
				o.WithResourceAlreadyRead(false),
				o.WithIssueIDs([]int{3, 4}),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/2/watchings", spath)
				assert.Equal(t, url.Values{
					"order":               {"asc"},
					"sort":                {"created"},
					"count":               {"10"},
					"offset":              {"20"},
					"resourceAlreadyRead": {"false"},
					"issueId[]":           {"3", "4"},
				}, query)
				return mock.NewResponse(fixture.Watching.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"success-empty": {
			userID: 1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewResponse(`[]`), nil
			},
			wantIDs: []int{},
		},

		"error-validation-userID-zero": {
			userID:                 0,
			wantValidationErrCount: 1,
		},
		"error-validation-userID-negative": {
			userID:                 -1,
			wantValidationErrCount: 1,
		},
		"error-validation-opt-single": {
			userID:                 1,
			opts:                   []*option.APIParamOption{o.WithWatchingSort("summary")},
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			userID:                 0,
			opts:                   []*option.APIParamOption{o.WithCount(0), o.WithWatchingSort("summary")},
			wantValidationErrCount: 3,
		},

		"error-nil-option-with-valid-values": {
			userID:                 1,
			opts:                   []*option.APIParamOption{o.WithCount(5), nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			userID:                 0,
			opts:                   []*option.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		"error-option-invalid-type": {
			userID:      1,
			opts:        []*option.APIParamOption{o.WithAlreadyRead(true)},
			wantErrType: &option.InvalidOptionKeyError{},
		},

		"error-client-network": {
			userID: 1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-json-decode": {
			userID: 1,
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
			s := watching.NewService(method)
			got, err := s.List(context.Background(), tc.userID, tc.opts...)

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
				var ves validation.Errors
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
			ids := make([]int, len(got))
			for i, w := range got {
				ids[i] = w.ID
			}
			assert.Equal(t, tc.wantIDs, ids)
		})
	}
}

func TestService_Count(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		userID int
		opts   []*option.APIParamOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
		wantCount              int
	}{
		"success": {
			userID: 1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1/watchings/count", spath)
				assert.Empty(t, query)
				return mock.NewResponse(fixture.Watching.CountJSON), nil
			},
			wantCount: 138,
		},
		"success-with-options": {
			userID: 1,
			opts: []*option.APIParamOption{
				o.WithResourceAlreadyRead(true),
				o.WithAlreadyRead(false),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1/watchings/count", spath)
				assert.Equal(t, url.Values{
					"resourceAlreadyRead": {"true"},
					"alreadyRead":         {"false"},
				}, query)
				return mock.NewResponse(fixture.Watching.CountJSON), nil
			},
			wantCount: 138,
		},

		"error-validation-userID-zero": {
			userID:                 0,
			wantValidationErrCount: 1,
		},
		"error-validation-userID-negative": {
			userID:                 -1,
			wantValidationErrCount: 1,
		},
		"error-validation-opt-single": {
			userID:                 1,
			opts:                   []*option.APIParamOption{mock.NewFailingCheckOption(option.ParamAlreadyRead)},
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			userID:                 0,
			opts:                   []*option.APIParamOption{mock.NewFailingCheckOption(option.ParamAlreadyRead)},
			wantValidationErrCount: 2,
		},

		"error-nil-option-with-valid-values": {
			userID:                 1,
			opts:                   []*option.APIParamOption{o.WithAlreadyRead(true), nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			userID:                 0,
			opts:                   []*option.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		"error-option-invalid-type": {
			userID:      1,
			opts:        []*option.APIParamOption{o.WithCount(10)},
			wantErrType: &option.InvalidOptionKeyError{},
		},

		"error-client-network": {
			userID: 1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-json-decode": {
			userID: 1,
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
			s := watching.NewService(method)
			got, err := s.Count(context.Background(), tc.userID, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Equal(t, 0, got)
				var target *option.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Equal(t, 0, got)
				var ves validation.Errors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Equal(t, 0, got)
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
		watchingID int

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			watchingID: 1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "watchings/1", spath)
				assert.Nil(t, query)
				return mock.NewResponse(fixture.Watching.SingleJSON), nil
			},
		},

		"error-validation-watchingID-zero": {
			watchingID:             0,
			wantValidationErrCount: 1,
		},
		"error-validation-watchingID-negative": {
			watchingID:             -1,
			wantValidationErrCount: 1,
		},

		"error-client-network": {
			watchingID: 1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-json-decode": {
			watchingID: 1,
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
			s := watching.NewService(method)
			got, err := s.One(context.Background(), tc.watchingID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves validation.Errors
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
			assert.Equal(t, 1, got.ID)
			assert.True(t, got.AlreadyRead)
			assert.Equal(t, "Check before release", got.Note)
			assert.Equal(t, "issue", got.Type)
			require.NotNil(t, got.Issue)
			assert.Equal(t, "PRJ-1", got.Issue.IssueKey)
		})
	}
}

func TestService_Create(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		issueIDOrKey string
		opts         []*option.APIParamOption

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success-issue-key": {
			issueIDOrKey: "PRJ-1",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "watchings", spath)
				assert.Equal(t, url.Values{"issueIdOrKey": {"PRJ-1"}}, form)
				return mock.NewCreatedResponse(fixture.Watching.SingleJSON), nil
			},
		},
		"success-issue-id-with-note": {
			issueIDOrKey: "1",
			opts:         []*option.APIParamOption{o.WithNote("Check before release")},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "watchings", spath)
				assert.Equal(t, url.Values{
					"issueIdOrKey": {"1"},
					"note":         {"Check before release"},
				}, form)
				return mock.NewCreatedResponse(fixture.Watching.SingleJSON), nil
			},
		},

		"error-validation-issueIDOrKey-empty": {
			issueIDOrKey:           "",
			wantValidationErrCount: 1,
		},
		"error-validation-opt-single": {
			issueIDOrKey:           "PRJ-1",
			opts:                   []*option.APIParamOption{mock.NewFailingCheckOption(option.ParamNote)},
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			issueIDOrKey:           "",
			opts:                   []*option.APIParamOption{mock.NewFailingCheckOption(option.ParamNote)},
			wantValidationErrCount: 2,
		},

		"error-nil-option-with-valid-values": {
			issueIDOrKey:           "PRJ-1",
			opts:                   []*option.APIParamOption{o.WithNote("note"), nil},
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
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-json-decode": {
			issueIDOrKey: "PRJ-1",
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
			s := watching.NewService(method)
			got, err := s.Create(context.Background(), tc.issueIDOrKey, tc.opts...)

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
				var ves validation.Errors
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
			assert.Equal(t, 1, got.ID)
		})
	}
}

func TestService_Update(t *testing.T) {
	cases := map[string]struct {
		watchingID int
		note       string

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			watchingID: 1,
			note:       "Check before release",
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "watchings/1", spath)
				assert.Equal(t, url.Values{"note": {"Check before release"}}, form)
				return mock.NewResponse(fixture.Watching.SingleJSON), nil
			},
		},
		"success-empty-note": {
			watchingID: 1,
			note:       "",
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "watchings/1", spath)
				assert.Equal(t, url.Values{"note": {""}}, form)
				return mock.NewResponse(fixture.Watching.SingleJSON), nil
			},
		},

		"error-validation-watchingID-zero": {
			watchingID:             0,
			wantValidationErrCount: 1,
		},
		"error-validation-watchingID-negative": {
			watchingID:             -1,
			wantValidationErrCount: 1,
		},

		"error-client-network": {
			watchingID: 1,
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-json-decode": {
			watchingID: 1,
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
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
			s := watching.NewService(method)
			got, err := s.Update(context.Background(), tc.watchingID, tc.note)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves validation.Errors
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
			assert.Equal(t, 1, got.ID)
		})
	}
}

func TestService_Delete(t *testing.T) {
	cases := map[string]struct {
		watchingID int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			watchingID: 1,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "watchings/1", spath)
				assert.Nil(t, form)
				return mock.NewResponse(fixture.Watching.SingleJSON), nil
			},
		},

		"error-validation-watchingID-zero": {
			watchingID:             0,
			wantValidationErrCount: 1,
		},
		"error-validation-watchingID-negative": {
			watchingID:             -1,
			wantValidationErrCount: 1,
		},

		"error-client-network": {
			watchingID: 1,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-json-decode": {
			watchingID: 1,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
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
			s := watching.NewService(method)
			got, err := s.Delete(context.Background(), tc.watchingID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves validation.Errors
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
			assert.Equal(t, 1, got.ID)
		})
	}
}

func TestService_MarkAsRead(t *testing.T) {
	cases := map[string]struct {
		watchingID int

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			watchingID: 1,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "watchings/1/markAsRead", spath)
				assert.Nil(t, form)
				return mock.NewNoContentResponse(), nil
			},
		},

		"error-validation-watchingID-zero": {
			watchingID:             0,
			wantValidationErrCount: 1,
		},
		"error-validation-watchingID-negative": {
			watchingID:             -1,
			wantValidationErrCount: 1,
		},

		"error-client-network": {
			watchingID: 1,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockPostFn != nil {
				method.Post = tc.mockPostFn
			}
			s := watching.NewService(method)
			err := s.MarkAsRead(context.Background(), tc.watchingID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				var ves validation.Errors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
		})
	}
}
