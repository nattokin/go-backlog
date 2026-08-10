package project_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/project"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestService_One(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success-key": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST", spath)
				return mock.NewResponse(fixture.Project.SingleJSON), nil
			},
		},
		"success-id": {
			projectIDOrKey: strconv.Itoa(6),
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/6", spath)
				return mock.NewResponse(fixture.Project.SingleJSON), nil
			},
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			wantValidationErrCount: 1,
		},

		// --- other errors ---
		"error-client-network": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
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
			s := project.NewService(method)
			p, err := s.One(context.Background(), tc.projectIDOrKey)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, p)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, p)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, p)
			assert.Equal(t, 6, p.ID)
			assert.Equal(t, "TEST", p.ProjectKey)
			assert.Equal(t, "test", p.Name)
		})
	}
}

func TestService_Create(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		key  string
		name string
		opts []*core.APIParamOption

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success-minimum": {
			key:  "TEST",
			name: "test",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects", spath)
				assert.Equal(t, "TEST", form.Get("key"))
				assert.Equal(t, "test", form.Get("name"))
				return mock.NewResponse(fixture.Project.SingleJSON), nil
			},
		},
		"success-without-option": {
			key:  "TEST",
			name: "test",
			opts: []*core.APIParamOption{},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "", form.Get("chartEnabled"))
				assert.Equal(t, "", form.Get("subtaskingEnabled"))
				assert.Equal(t, "", form.Get("projectLeaderCanEditProjectLeader"))
				assert.Equal(t, "", form.Get("textFormattingRule"))
				return mock.NewResponse(fixture.Project.SingleJSON), nil
			},
		},
		"success-with-options": {
			key:  "TEST",
			name: "test",
			opts: []*core.APIParamOption{
				o.WithChartEnabled(true),
				o.WithSubtaskingEnabled(true),
				o.WithProjectLeaderCanEditProjectLeader(true),
				o.WithTextFormattingRule("backlog"),
			},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "true", form.Get("chartEnabled"))
				assert.Equal(t, "true", form.Get("subtaskingEnabled"))
				assert.Equal(t, "true", form.Get("projectLeaderCanEditProjectLeader"))
				assert.Equal(t, "backlog", form.Get("textFormattingRule"))
				return mock.NewResponse(fixture.Project.SingleJSON), nil
			},
		},

		// --- validation errors: argument only ---
		"error-validation-key-empty": {
			key:                    "",
			name:                   "test",
			wantValidationErrCount: 1,
		},
		"error-validation-name-empty": {
			key:                    "TEST",
			name:                   "",
			wantValidationErrCount: 1,
		},

		// --- validation errors: optional opts only ---
		"error-validation-opt-single": {
			key:                    "TEST",
			name:                   "test",
			opts:                   []*core.APIParamOption{o.WithTextFormattingRule("invalid")},
			wantValidationErrCount: 1,
		},

		// --- validation errors: all ---
		"error-validation-all": {
			key:                    "",
			name:                   "",
			opts:                   []*core.APIParamOption{o.WithTextFormattingRule("invalid")},
			wantValidationErrCount: 3,
		},

		// WithAll is not in Create's valid types → InvalidOptionKeyError (not InvalidOptionError)
		"error-nil-option-with-valid-values": {
			key:         "TEST",
			name:        "test",
			opts:        []*core.APIParamOption{o.WithAll(true), nil},
			wantErrType: &core.InvalidOptionKeyError{},
		},
		// --- fail-fast: nil option among invalid values ---
		"error-nil-option-with-invalid-values": {
			key:                    "",
			name:                   "",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type-with-valid-values": {
			key:         "TEST",
			name:        "test",
			opts:        []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			key:         "",
			name:        "",
			opts:        []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-client-network": {
			key:  "TEST",
			name: "test",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			key:  "TEST",
			name: "test",
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
			s := project.NewService(method)
			p, err := s.Create(context.Background(), tc.key, tc.name, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, p)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, p)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, p)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, p)
			assert.Equal(t, tc.key, p.ProjectKey)
			assert.Equal(t, tc.name, p.Name)
		})
	}
}

func TestService_Update(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		option         *core.APIParamOption
		opts           []*core.APIParamOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success-key": {
			projectIDOrKey: "TEST",
			option:         o.WithName("test"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST", spath)
				return mock.NewResponse(fixture.Project.SingleJSON), nil
			},
		},
		"success-id": {
			projectIDOrKey: "1234",
			option:         o.WithName("test"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/1234", spath)
				return mock.NewResponse(fixture.Project.SingleJSON), nil
			},
		},
		"success-full-options": {
			projectIDOrKey: "TEST",
			option:         o.WithKey("TEST1"),
			opts: []*core.APIParamOption{
				o.WithName("test1"),
				o.WithChartEnabled(true),
				o.WithSubtaskingEnabled(true),
				o.WithProjectLeaderCanEditProjectLeader(true),
				o.WithTextFormattingRule("backlog"),
				o.WithArchived(true),
			},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "TEST1", form.Get("key"))
				assert.Equal(t, "test1", form.Get("name"))
				assert.Equal(t, "true", form.Get("chartEnabled"))
				assert.Equal(t, "true", form.Get("subtaskingEnabled"))
				assert.Equal(t, "true", form.Get("projectLeaderCanEditProjectLeader"))
				assert.Equal(t, "backlog", form.Get("textFormattingRule"))
				assert.Equal(t, "true", form.Get("archived"))
				return mock.NewResponse(fixture.Project.SingleJSON), nil
			},
		},

		// --- validation errors: argument only ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			option:                 o.WithName("test"),
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			option:                 o.WithName("test"),
			wantValidationErrCount: 1,
		},

		// --- validation errors: fixed option only ---
		"error-validation-fixed-option": {
			projectIDOrKey:         "TEST",
			option:                 o.WithTextFormattingRule("invalid"),
			wantValidationErrCount: 1,
		},

		// --- validation errors: optional opts only ---
		"error-validation-opt-single": {
			projectIDOrKey:         "TEST",
			option:                 o.WithName("test"),
			opts:                   []*core.APIParamOption{o.WithTextFormattingRule("invalid")},
			wantValidationErrCount: 1,
		},

		// --- validation errors: all ---
		"error-validation-all": {
			projectIDOrKey:         "",
			option:                 o.WithTextFormattingRule("invalid"),
			opts:                   []*core.APIParamOption{o.WithKey("")},
			wantValidationErrCount: 3,
		},

		// WithAll is not in Update's valid types → InvalidOptionKeyError (not InvalidOptionError)
		"error-nil-option-with-valid-values": {
			projectIDOrKey: "TEST",
			option:         o.WithName("test"),
			opts:           []*core.APIParamOption{o.WithAll(true), nil},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		// --- fail-fast: nil option among invalid values ---
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			option:                 o.WithTextFormattingRule("invalid"),
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type-with-valid-values": {
			projectIDOrKey: "TEST",
			option:         o.WithName("test"),
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			projectIDOrKey: "",
			option:         o.WithTextFormattingRule("invalid"),
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-client-network": {
			projectIDOrKey: "TEST",
			option:         o.WithName("test"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			option:         o.WithName("test"),
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
			s := project.NewService(method)
			p, err := s.Update(context.Background(), tc.projectIDOrKey, tc.option, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, p)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, p)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, p)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, p)
		})
	}
}

func TestService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success-key": {
			projectIDOrKey: "TEST",
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST", spath)
				return mock.NewResponse(fixture.Project.SingleJSON), nil
			},
		},
		"success-id": {
			projectIDOrKey: "1234",
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/1234", spath)
				return mock.NewResponse(fixture.Project.SingleJSON), nil
			},
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			wantValidationErrCount: 1,
		},

		// --- other errors ---
		"error-client-network": {
			projectIDOrKey: "TEST",
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
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
			s := project.NewService(method)
			p, err := s.Delete(context.Background(), tc.projectIDOrKey)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, p)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, p)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, p)
			assert.Equal(t, "TEST", p.ProjectKey)
		})
	}
}
