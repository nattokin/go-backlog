package project_test

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
	"github.com/nattokin/go-backlog/internal/domain/project"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestProjectUserService_List(t *testing.T) {
	opt := &core.OptionService{}

	cases := map[string]struct {
		projectKey string
		opts       []*core.APIParamOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success-no-options": {
			projectKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/users", spath)
				return mock.NewResponse(fixture.User.ListJSON), nil
			},
		},
		"success-excludeGroupMembers-true": {
			projectKey: "TEST2",
			opts:       []*core.APIParamOption{opt.WithExcludeGroupMembers(true)},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "true", query.Get("excludeGroupMembers"))
				return mock.NewResponse(fixture.User.ListJSON), nil
			},
		},
		"success-excludeGroupMembers-false": {
			projectKey: "TEST3",
			opts:       []*core.APIParamOption{opt.WithExcludeGroupMembers(false)},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "false", query.Get("excludeGroupMembers"))
				return mock.NewResponse(fixture.User.ListJSON), nil
			},
		},

		// --- validation errors ---
		"error-validation-projectKey-empty": {
			projectKey:             "",
			wantValidationErrCount: 1,
		},
		"error-validation-projectKey-zero": {
			projectKey:             "0",
			wantValidationErrCount: 1,
		},
		"error-validation-opt": {
			projectKey:             "TEST",
			opts:                   []*core.APIParamOption{mock.NewFailingCheckOption(core.ParamExcludeGroupMembers)},
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectKey:             "",
			opts:                   []*core.APIParamOption{mock.NewFailingCheckOption(core.ParamExcludeGroupMembers)},
			wantValidationErrCount: 2,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			projectKey:             "TEST",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectKey:             "",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type": {
			projectKey:  "TEST",
			opts:        []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-response-invalid-json": {
			projectKey: "TEST",
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
			s := project.NewUserService(method)
			users, err := s.List(context.Background(), tc.projectKey, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, users)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, users)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, users)
				return
			}

			assert.NoError(t, err)
			require.Len(t, users, 4)
			require.NotNil(t, users[0])
			assert.Equal(t, "admin", users[0].UserID)
			assert.Equal(t, "admin", users[0].Name)
			assert.Equal(t, "eguchi@nulab.example", users[0].MailAddress)
			assert.Equal(t, 1, users[0].RoleType)
		})
	}
}

func TestProjectUserService_Add(t *testing.T) {
	cases := map[string]struct {
		projectKey string
		userID     int

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectKey: "TEST",
			userID:     1,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return mock.NewResponse(fixture.User.SingleJSON), nil
			},
		},

		// --- validation errors ---
		"error-validation-projectKey-empty": {
			projectKey:             "",
			userID:                 1,
			wantValidationErrCount: 1,
		},
		"error-validation-userID-zero": {
			projectKey:             "TEST",
			userID:                 0,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectKey:             "",
			userID:                 0,
			wantValidationErrCount: 2,
		},

		// --- other errors ---
		"error-response-invalid-json": {
			projectKey: "TEST",
			userID:     1,
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
			s := project.NewUserService(method)
			u, err := s.Add(context.Background(), tc.projectKey, tc.userID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, u)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, u)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, u)
			assert.Equal(t, "admin", u.UserID)
			assert.Equal(t, "admin", u.Name)
			assert.Equal(t, "eguchi@nulab.example", u.MailAddress)
			assert.Equal(t, 1, u.RoleType)
		})
	}
}

func TestProjectUserService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectKey string
		userID     int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectKey: "TEST",
			userID:     1,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return mock.NewResponse(fixture.User.SingleJSON), nil
			},
		},
		"success-projectIDOrKey-number": {
			projectKey: "1234",
			userID:     1,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/1234/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return mock.NewResponse(fixture.User.SingleJSON), nil
			},
		},

		// --- validation errors ---
		"error-validation-projectKey-empty": {
			projectKey:             "",
			userID:                 1,
			wantValidationErrCount: 1,
		},
		"error-validation-userID-zero": {
			projectKey:             "TEST",
			userID:                 0,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectKey:             "",
			userID:                 0,
			wantValidationErrCount: 2,
		},

		// --- other errors ---
		"error-response-invalid-json": {
			projectKey: "TEST",
			userID:     1,
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
			s := project.NewUserService(method)
			u, err := s.Delete(context.Background(), tc.projectKey, tc.userID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, u)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, u)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, u)
			assert.Equal(t, "admin", u.UserID)
			assert.Equal(t, "admin", u.Name)
			assert.Equal(t, "eguchi@nulab.example", u.MailAddress)
			assert.Equal(t, 1, u.RoleType)
		})
	}
}

func TestProjectUserService_AddAdmin(t *testing.T) {
	cases := map[string]struct {
		projectKey string
		userID     int

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectKey: "TEST",
			userID:     1,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/administrators", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return mock.NewResponse(fixture.User.SingleJSON), nil
			},
		},

		// --- validation errors ---
		"error-validation-projectKey-empty": {
			projectKey:             "",
			userID:                 1,
			wantValidationErrCount: 1,
		},
		"error-validation-userID-zero": {
			projectKey:             "TEST",
			userID:                 0,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectKey:             "",
			userID:                 0,
			wantValidationErrCount: 2,
		},

		// --- other errors ---
		"error-response-invalid-json": {
			projectKey: "TEST",
			userID:     1,
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
			s := project.NewUserService(method)
			u, err := s.AddAdmin(context.Background(), tc.projectKey, tc.userID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, u)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, u)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, u)
			assert.Equal(t, "admin", u.UserID)
			assert.Equal(t, "admin", u.Name)
			assert.Equal(t, "eguchi@nulab.example", u.MailAddress)
			assert.Equal(t, 1, u.RoleType)
		})
	}
}

func TestProjectUserService_AdminList(t *testing.T) {
	cases := map[string]struct {
		projectKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"error-client-network": {
			projectKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/administrators", spath)
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-validation-projectKey-empty": {
			projectKey:             "",
			wantValidationErrCount: 1,
		},
		"error-validation-projectKey-zero": {
			projectKey:             "0",
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
			s := project.NewUserService(method)
			users, err := s.AdminList(context.Background(), tc.projectKey)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, users)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			assert.Error(t, err)
			assert.ErrorAs(t, err, &tc.wantErrType)
			assert.Nil(t, users)
		})
	}
}

func TestProjectUserService_DeleteAdmin(t *testing.T) {
	cases := map[string]struct {
		projectKey string
		userID     int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"error-client-network": {
			projectKey: "TEST",
			userID:     1,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/administrators", spath)
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-validation-projectKey-empty": {
			projectKey:             "",
			userID:                 1,
			wantValidationErrCount: 1,
		},
		"error-validation-userID-zero": {
			projectKey:             "TEST",
			userID:                 0,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectKey:             "",
			userID:                 0,
			wantValidationErrCount: 2,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockDeleteFn != nil {
				method.Delete = tc.mockDeleteFn
			}
			s := project.NewUserService(method)
			u, err := s.DeleteAdmin(context.Background(), tc.projectKey, tc.userID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, u)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			assert.Error(t, err)
			assert.ErrorAs(t, err, &tc.wantErrType)
			assert.Nil(t, u)
		})
	}
}
