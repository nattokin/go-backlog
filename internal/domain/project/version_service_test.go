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

func TestVersionService_List(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		opts           []*core.APIParamOption

		mockGetFn func(context.Context, string, url.Values) (*http.Response, error)

		wantLen                int
		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success": {
			projectIDOrKey: "TEST",
			wantLen:        2,
			opts:           []*core.APIParamOption{o.WithArchived(true)},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/versions", spath)
				assert.Equal(t, "true", query.Get("archived"))
				return mock.NewResponse(fixture.Version.ListJSON), nil
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

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "TEST",
			opts:                   []*core.APIParamOption{o.WithArchived(true), nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type-with-valid-values": {
			projectIDOrKey: "TEST",
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			projectIDOrKey: "",
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-client-network": {
			projectIDOrKey: "TEST",
			wantErrType:    errors.New(""),
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			wantErrType:    &json.SyntaxError{},
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				m.Get = tc.mockGetFn
			}
			s := project.NewVersionService(m)
			got, err := s.List(context.Background(), tc.projectIDOrKey, tc.opts...)

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
			require.Len(t, got, tc.wantLen)
		})
	}
}

func TestVersionService_Add(t *testing.T) {
	o := &core.OptionService{}
	date := "2025-01-01"

	cases := map[string]struct {
		projectIDOrKey string
		name           string
		opts           []*core.APIParamOption

		mockPostFn func(context.Context, string, url.Values) (*http.Response, error)

		wantID                 int
		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success": {
			projectIDOrKey: "TEST",
			name:           "v1",
			opts:           []*core.APIParamOption{o.WithDescription("desc"), o.WithStartDate(date), o.WithReleaseDueDate(date)},
			wantID:         fixture.Version.Single.ID,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/versions", spath)
				assert.Equal(t, "v1", form.Get("name"))
				return mock.NewResponse(fixture.Version.SingleJSON), nil
			},
		},

		// --- validation errors: argument only ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			name:                   "v1",
			wantValidationErrCount: 1,
		},

		// --- validation errors: fixed option only ---
		"error-validation-name-empty": {
			projectIDOrKey:         "TEST",
			name:                   "",
			wantValidationErrCount: 1,
		},

		// --- validation errors: all ---
		"error-validation-all": {
			projectIDOrKey:         "",
			name:                   "",
			wantValidationErrCount: 2,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "TEST",
			name:                   "v1",
			opts:                   []*core.APIParamOption{o.WithDescription("desc"), nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			name:                   "",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type-with-valid-values": {
			projectIDOrKey: "TEST",
			name:           "v1",
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			projectIDOrKey: "",
			name:           "",
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-option-set-failed": {
			projectIDOrKey: "TEST",
			name:           "v1",
			opts:           []*core.APIParamOption{mock.NewFailingSetOption(core.ParamDescription)},
			wantErrType:    errors.New(""),
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			name:           "v1",
			wantErrType:    errors.New(""),
			mockPostFn:     func(context.Context, string, url.Values) (*http.Response, error) { return nil, errors.New("network") },
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			name:           "v1",
			wantErrType:    &json.SyntaxError{},
			mockPostFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := mock.NewMethod(t)
			if tc.mockPostFn != nil {
				m.Post = tc.mockPostFn
			}
			s := project.NewVersionService(m)
			got, err := s.Add(context.Background(), tc.projectIDOrKey, tc.name, tc.opts...)

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
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestVersionService_Update(t *testing.T) {
	o := &core.OptionService{}
	date := "2025-01-01"

	cases := map[string]struct {
		projectIDOrKey string
		versionID      int
		option         *core.APIParamOption
		opts           []*core.APIParamOption

		mockPatchFn func(context.Context, string, url.Values) (*http.Response, error)

		wantID                 int
		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success": {
			projectIDOrKey: "TEST",
			versionID:      1,
			option:         o.WithName("name"),
			opts:           []*core.APIParamOption{o.WithReleaseDueDate(date)},
			wantID:         fixture.Version.Single.ID,
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/versions/1", spath)
				return mock.NewResponse(fixture.Version.SingleJSON), nil
			},
		},

		// --- validation errors: argument only ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			versionID:              1,
			option:                 o.WithName("name"),
			wantValidationErrCount: 1,
		},
		"error-validation-versionID-zero": {
			projectIDOrKey:         "TEST",
			versionID:              0,
			option:                 o.WithName("name"),
			wantValidationErrCount: 1,
		},
		"error-validation-versionID-negative": {
			projectIDOrKey:         "TEST",
			versionID:              -1,
			option:                 o.WithName("name"),
			wantValidationErrCount: 1,
		},

		// --- validation errors: fixed option only ---
		"error-validation-fixed-option": {
			projectIDOrKey:         "TEST",
			versionID:              1,
			option:                 o.WithName(""),
			wantValidationErrCount: 1,
		},

		// --- validation errors: all ---
		"error-validation-all": {
			projectIDOrKey:         "",
			versionID:              0,
			option:                 o.WithName(""),
			wantValidationErrCount: 3,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "TEST",
			versionID:              1,
			option:                 o.WithName("name"),
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			versionID:              0,
			option:                 o.WithName(""),
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type-with-valid-values": {
			projectIDOrKey: "TEST",
			versionID:      1,
			option:         mock.NewInvalidTypeOption(),
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			projectIDOrKey: "",
			versionID:      0,
			option:         mock.NewInvalidTypeOption(),
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-option-set-failed": {
			projectIDOrKey: "TEST",
			versionID:      1,
			option:         mock.NewFailingSetOption(core.ParamArchived),
			wantErrType:    errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			versionID:      1,
			option:         o.WithName("name"),
			wantErrType:    &json.SyntaxError{},
			mockPatchFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := mock.NewMethod(t)
			if tc.mockPatchFn != nil {
				m.Patch = tc.mockPatchFn
			}
			s := project.NewVersionService(m)
			got, err := s.Update(context.Background(), tc.projectIDOrKey, tc.versionID, tc.option, tc.opts...)

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
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestVersionService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		versionID      int

		mockDeleteFn func(context.Context, string, url.Values) (*http.Response, error)

		wantID                 int
		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			versionID:      1,
			wantID:         fixture.Version.Single.ID,
			mockDeleteFn: func(ctx context.Context, spath string, _ url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/versions/1", spath)
				return mock.NewResponse(fixture.Version.SingleJSON), nil
			},
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			versionID:              1,
			wantValidationErrCount: 1,
		},
		"error-validation-versionID-zero": {
			projectIDOrKey:         "TEST",
			versionID:              0,
			wantValidationErrCount: 1,
		},
		"error-validation-versionID-negative": {
			projectIDOrKey:         "TEST",
			versionID:              -1,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			versionID:              0,
			wantValidationErrCount: 2,
		},

		// --- other errors ---
		"error-client-network": {
			projectIDOrKey: "TEST",
			versionID:      1,
			wantErrType:    errors.New(""),
			mockDeleteFn:   func(context.Context, string, url.Values) (*http.Response, error) { return nil, errors.New("network") },
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			versionID:      1,
			wantErrType:    &json.SyntaxError{},
			mockDeleteFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := mock.NewMethod(t)
			if tc.mockDeleteFn != nil {
				m.Delete = tc.mockDeleteFn
			}
			s := project.NewVersionService(m)
			got, err := s.Delete(context.Background(), tc.projectIDOrKey, tc.versionID)

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
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}
