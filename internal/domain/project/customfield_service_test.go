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

func TestCustomFieldService_List(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantLen                int
		wantErrType            error
		wantValidationErrCount int
	}{
		"success-key": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/customFields", spath)
				return mock.NewResponse(fixture.CustomField.ListJSON), nil
			},
			wantLen: 2,
		},
		"success-id": {
			projectIDOrKey: "6",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/6/customFields", spath)
				return mock.NewResponse(fixture.CustomField.ListJSON), nil
			},
			wantLen: 2,
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			wantValidationErrCount: 1,
		},

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
			s := project.NewCustomFieldService(method)
			fields, err := s.List(context.Background(), tc.projectIDOrKey)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, fields)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, fields)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, fields)
			assert.Len(t, fields, tc.wantLen)
		})
	}
}

func TestCustomFieldService_Create(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		fieldType      int
		name           string
		opts           []*core.APIParamOption

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success": {
			projectIDOrKey: "TEST",
			fieldType:      1,
			name:           "Sprint",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/customFields", spath)
				assert.Equal(t, "1", form.Get("typeId"))
				assert.Equal(t, "Sprint", form.Get("name"))
				return mock.NewResponse(fixture.CustomField.SingleJSON), nil
			},
		},
		"success-with-opts": {
			projectIDOrKey: "TEST",
			fieldType:      1,
			name:           "Sprint",
			opts:           []*core.APIParamOption{o.WithDescription("sprint number"), o.WithRequired(true)},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "sprint number", form.Get("description"))
				assert.Equal(t, "true", form.Get("required"))
				return mock.NewResponse(fixture.CustomField.SingleJSON), nil
			},
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			fieldType:              1,
			name:                   "Sprint",
			wantValidationErrCount: 1,
		},

		"error-validation-fieldType-zero": {
			projectIDOrKey:         "TEST",
			fieldType:              0,
			name:                   "Sprint",
			wantValidationErrCount: 1,
		},
		"error-validation-name-empty": {
			projectIDOrKey:         "TEST",
			fieldType:              1,
			name:                   "",
			wantValidationErrCount: 1,
		},
		"error-validation-fieldType-and-name": {
			projectIDOrKey:         "TEST",
			fieldType:              0,
			name:                   "",
			wantValidationErrCount: 2,
		},

		"error-validation-all": {
			projectIDOrKey:         "",
			fieldType:              0,
			name:                   "",
			wantValidationErrCount: 3,
		},

		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "TEST",
			fieldType:              1,
			name:                   "Sprint",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			fieldType:              0,
			name:                   "",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		"error-option-invalid-type-with-valid-values": {
			projectIDOrKey: "TEST",
			fieldType:      1,
			name:           "Sprint",
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			projectIDOrKey: "",
			fieldType:      0,
			name:           "",
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		"error-client-network": {
			projectIDOrKey: "TEST",
			fieldType:      1,
			name:           "Sprint",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			fieldType:      1,
			name:           "Sprint",
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
			s := project.NewCustomFieldService(method)
			field, err := s.Create(context.Background(), tc.projectIDOrKey, tc.fieldType, tc.name, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, field)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, field)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, field)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, field)
			assert.Equal(t, 1, field.ID)
			assert.Equal(t, "Sprint", field.Name)
		})
	}
}

func TestCustomFieldService_Update(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		customFieldID  int
		option         *core.APIParamOption
		opts           []*core.APIParamOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			option:         o.WithName("Sprint Updated"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/customFields/1", spath)
				assert.Equal(t, "Sprint Updated", form.Get("name"))
				return mock.NewResponse(fixture.CustomField.SingleJSON), nil
			},
		},
		"success-with-opts": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			option:         o.WithName("Sprint Updated"),
			opts:           []*core.APIParamOption{o.WithRequired(true)},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "true", form.Get("required"))
				return mock.NewResponse(fixture.CustomField.SingleJSON), nil
			},
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			customFieldID:          1,
			option:                 o.WithName("Sprint Updated"),
			wantValidationErrCount: 1,
		},
		"error-validation-customFieldID-zero": {
			projectIDOrKey:         "TEST",
			customFieldID:          0,
			option:                 o.WithName("Sprint Updated"),
			wantValidationErrCount: 1,
		},

		"error-validation-fixed-option": {
			projectIDOrKey:         "TEST",
			customFieldID:          1,
			option:                 o.WithName(""),
			wantValidationErrCount: 1,
		},

		"error-validation-all": {
			projectIDOrKey:         "",
			customFieldID:          0,
			option:                 o.WithName(""),
			wantValidationErrCount: 3,
		},

		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "TEST",
			customFieldID:          1,
			option:                 o.WithName("Sprint"),
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			customFieldID:          0,
			option:                 o.WithName(""),
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		"error-option-invalid-type-with-valid-values": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			option:         mock.NewInvalidTypeOption(),
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			projectIDOrKey: "",
			customFieldID:  0,
			option:         mock.NewInvalidTypeOption(),
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		"error-client-network": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			option:         o.WithName("Sprint Updated"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			option:         o.WithName("Sprint Updated"),
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
			s := project.NewCustomFieldService(method)
			field, err := s.Update(context.Background(), tc.projectIDOrKey, tc.customFieldID, tc.option, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, field)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, field)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, field)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, field)
			assert.Equal(t, 1, field.ID)
		})
	}
}

func TestCustomFieldService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		customFieldID  int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/customFields/1", spath)
				return mock.NewResponse(fixture.CustomField.SingleJSON), nil
			},
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			customFieldID:          1,
			wantValidationErrCount: 1,
		},
		"error-validation-customFieldID-zero": {
			projectIDOrKey:         "TEST",
			customFieldID:          0,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			customFieldID:          0,
			wantValidationErrCount: 2,
		},

		"error-client-network": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
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
			s := project.NewCustomFieldService(method)
			field, err := s.Delete(context.Background(), tc.projectIDOrKey, tc.customFieldID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, field)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, field)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, field)
			assert.Equal(t, 1, field.ID)
			assert.Equal(t, "Sprint", field.Name)
		})
	}
}
