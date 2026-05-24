package user_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/user"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestUserStarService_List(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		userID int
		opts   []*core.APIParamOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
		wantLen                int
	}{
		"success-no-options": {
			userID: 1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1/stars", spath)
				return mock.NewResponse(`[{"id":10},{"id":20}]`), nil
			},
			wantLen: 2,
		},
		"success-with-count": {
			userID: 2,
			opts:   []*core.APIParamOption{o.WithCount(5)},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "5", query.Get("count"))
				return mock.NewResponse(`[{"id":1}]`), nil
			},
			wantLen: 1,
		},

		// --- validation errors ---
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
			opts:                   []*core.APIParamOption{o.WithCount(0)},
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			userID:                 0,
			opts:                   []*core.APIParamOption{o.WithCount(0)},
			wantValidationErrCount: 2,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			userID:                 1,
			opts:                   []*core.APIParamOption{o.WithCount(5), nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			userID:                 0,
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type": {
			userID:      1,
			opts:        []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
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
			wantErrType: errors.New(""),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}
			s := user.NewStarService(method)
			got, err := s.List(context.Background(), tc.userID, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, got)
				var target *core.InvalidOptionError
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
			assert.Len(t, got, tc.wantLen)
		})
	}
}

func TestUserStarService_Count(t *testing.T) {
	cases := map[string]struct {
		userID int

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantCount              int
	}{
		"success": {
			userID: 1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1/stars/count", spath)
				return mock.NewResponse(`{"count":42}`), nil
			},
			wantCount: 42,
		},

		// --- validation errors ---
		"error-validation-userID-zero": {
			userID:                 0,
			wantValidationErrCount: 1,
		},
		"error-validation-userID-negative": {
			userID:                 -1,
			wantValidationErrCount: 1,
		},

		// --- other errors ---
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
			wantErrType: errors.New(""),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}
			s := user.NewStarService(method)
			got, err := s.Count(context.Background(), tc.userID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Equal(t, 0, got)
				var ves core.ValidationErrors
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
