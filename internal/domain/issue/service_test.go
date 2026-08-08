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

func TestService_Count(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		opts []*core.APIParamOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
		wantCount              int
	}{
		"success-no-options": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/count", spath)
				return mock.NewResponse(`{"count":42}`), nil
			},
			wantCount: 42,
		},
		"success-with-projectIDs": {
			opts: []*core.APIParamOption{o.WithProjectIDs([]int{10, 20})},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"10", "20"}, query["projectId[]"])
				return mock.NewResponse(`{"count":5}`), nil
			},
			wantCount: 5,
		},

		// --- validation errors ---
		"error-validation-opt-single": {
			opts:                   []*core.APIParamOption{o.WithProjectIDs([]int{0})},
			wantValidationErrCount: 1,
		},
		// WithCount is not in the valid types for Count — results in InvalidOptionKeyError
		"error-validation-opt-multiple": {
			opts:        []*core.APIParamOption{o.WithProjectIDs([]int{0}), o.WithCount(0)},
			wantErrType: &core.InvalidOptionKeyError{},
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			opts:                   []*core.APIParamOption{o.WithProjectIDs([]int{1}), nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			opts:                   []*core.APIParamOption{o.WithProjectIDs([]int{0}), nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type": {
			opts:        []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-client-network": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, &core.APIResponseError{}
			},
			wantErrType: &core.APIResponseError{},
		},
		"error-response-invalid-json": {
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
			s := issue.NewService(method)
			count, err := s.Count(context.Background(), tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Zero(t, count)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

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

func TestService_Participants(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantIDs                []int
	}{
		"success-by-key": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/participants", spath)
				return mock.NewResponse(fixture.User.ListJSON), nil
			},
			wantIDs: []int{1, 2, 3, 4},
		},
		"success-by-id": {
			issueIDOrKey: "1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1/participants", spath)
				return mock.NewResponse(fixture.User.ListJSON), nil
			},
			wantIDs: []int{1, 2, 3, 4},
		},

		// --- validation errors ---
		"error-validation-issueIDOrKey-empty": {
			issueIDOrKey:           "",
			wantValidationErrCount: 1,
		},
		"error-validation-issueIDOrKey-zero": {
			issueIDOrKey:           "0",
			wantValidationErrCount: 1,
		},

		// --- other errors ---
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
				return nil, &core.APIResponseError{}
			},
			wantErrType: &core.APIResponseError{},
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
			s := issue.NewService(method)
			got, err := s.Participants(context.Background(), tc.issueIDOrKey)

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
