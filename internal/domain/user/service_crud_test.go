package user_test

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
	"github.com/nattokin/go-backlog/internal/domain/user"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestUserService_Add(t *testing.T) {
	cases := map[string]struct {
		userID      string
		password    string
		name        string
		mailAddress string
		roleType    int

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			userID:      "admin",
			password:    "password",
			name:        "admin",
			mailAddress: "eguchi@nulab.example",
			roleType:    1,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users", spath)
				assert.Equal(t, "admin", form.Get("userId"))
				return mock.NewResponse(fixture.User.SingleJSON), nil
			},
		},

		// --- validation errors: argument only ---
		"error-validation-userID-empty": {
			userID:                 "",
			password:               "password",
			name:                   "admin",
			mailAddress:            "admin@example.com",
			roleType:               1,
			wantValidationErrCount: 1,
		},

		// --- validation errors: fixed options only ---
		"error-validation-password-too-short": {
			userID:                 "admin",
			password:               "short",
			name:                   "admin",
			mailAddress:            "admin@example.com",
			roleType:               1,
			wantValidationErrCount: 1,
		},
		"error-validation-name-empty": {
			userID:                 "admin",
			password:               "password",
			name:                   "",
			mailAddress:            "admin@example.com",
			roleType:               1,
			wantValidationErrCount: 1,
		},
		"error-validation-mailAddress-invalid": {
			userID:                 "admin",
			password:               "password",
			name:                   "admin",
			mailAddress:            "not-an-email",
			roleType:               1,
			wantValidationErrCount: 1,
		},
		"error-validation-roleType-0": {
			userID:                 "admin",
			password:               "password",
			name:                   "admin",
			mailAddress:            "admin@example.com",
			roleType:               0,
			wantValidationErrCount: 1,
		},
		"error-validation-roleType-7": {
			userID:                 "admin",
			password:               "password",
			name:                   "admin",
			mailAddress:            "admin@example.com",
			roleType:               7,
			wantValidationErrCount: 1,
		},
		"error-validation-multiple-options": {
			userID:                 "admin",
			password:               "short",
			name:                   "",
			mailAddress:            "not-an-email",
			roleType:               1,
			wantValidationErrCount: 3,
		},

		// --- validation errors: all ---
		"error-validation-all": {
			userID:                 "",
			password:               "short",
			name:                   "",
			mailAddress:            "not-an-email",
			roleType:               0,
			wantValidationErrCount: 4,
		},

		// --- other errors ---
		"error-client-network": {
			userID:      "admin",
			password:    "password",
			name:        "admin",
			mailAddress: "eguchi@nulab.example",
			roleType:    1,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network failure")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			userID:      "admin",
			password:    "password",
			name:        "admin",
			mailAddress: "eguchi@nulab.example",
			roleType:    1,
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
			s := user.NewService(method)
			got, err := s.Add(context.Background(), tc.userID, tc.password, tc.name, tc.mailAddress, tc.roleType)

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
			assert.Equal(t, "admin", got.UserID)
		})
	}
}

func TestUserService_One(t *testing.T) {
	cases := map[string]struct {
		id int

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success-id-1": {
			id: 1,
			mockGetFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				return mock.NewResponse(fixture.User.SingleJSON), nil
			},
		},

		// --- validation errors ---
		"error-validation-id-zero": {
			id:                     0,
			wantValidationErrCount: 1,
		},
		"error-validation-id-negative": {
			id:                     -1,
			wantValidationErrCount: 1,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}
			s := user.NewService(method)
			got, err := s.One(context.Background(), tc.id)

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
			assert.Equal(t, "admin", got.UserID)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		id     int
		option *core.APIParamOption
		opts   []*core.APIParamOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success": {
			id:     1,
			option: o.WithPassword("password"),
			opts: []*core.APIParamOption{
				o.WithName("admin"),
				o.WithMailAddress("eguchi@nulab.example"),
				o.WithRoleType(1),
			},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				return mock.NewResponse(fixture.User.SingleJSON), nil
			},
		},

		// --- validation errors: fixed option only ---
		"error-validation-fixed-option-name-empty": {
			id:                     1,
			option:                 o.WithName(""),
			wantValidationErrCount: 1,
		},

		// --- validation errors: optional opts only ---
		"error-validation-opt-single": {
			id:                     1,
			option:                 o.WithName("admin"),
			opts:                   []*core.APIParamOption{o.WithPassword("short")},
			wantValidationErrCount: 1,
		},
		"error-validation-opt-multiple": {
			id:                     1,
			option:                 o.WithName(""),
			opts:                   []*core.APIParamOption{o.WithPassword("short"), o.WithRoleType(0)},
			wantValidationErrCount: 3,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			id:                     1,
			option:                 o.WithName("admin"),
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			id:                     1,
			option:                 o.WithName(""),
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type": {
			id:          1,
			option:      mock.NewInvalidTypeOption(),
			wantErrType: &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-option-set-failed": {
			id:          1,
			option:      mock.NewFailingSetOption(core.ParamName),
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			id:     1,
			option: o.WithName("admin"),
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
			s := user.NewService(method)
			got, err := s.Update(context.Background(), tc.id, tc.option, tc.opts...)

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
			require.NotNil(t, got)
			assert.Equal(t, "admin", got.UserID)
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	cases := map[string]struct {
		id int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success-id-1": {
			id: 1,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				return mock.NewResponse(fixture.User.SingleJSON), nil
			},
		},

		// --- validation errors ---
		"error-validation-id-zero": {
			id:                     0,
			wantValidationErrCount: 1,
		},
		"error-validation-id-negative": {
			id:                     -1,
			wantValidationErrCount: 1,
		},

		// --- other errors ---
		"error-response-invalid-json": {
			id: 1,
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
			s := user.NewService(method)
			got, err := s.Delete(context.Background(), tc.id)

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

			assert.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}
