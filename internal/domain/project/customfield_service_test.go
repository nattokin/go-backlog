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

	"github.com/nattokin/go-backlog/internal/domain/project"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
	"github.com/nattokin/go-backlog/internal/validation"
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
				var ves validation.Errors
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
	o := &option.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		fieldType      int
		name           string
		opts           []*option.APIParamOption

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
			opts:           []*option.APIParamOption{o.WithDescription("sprint number"), o.WithRequired(true)},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "sprint number", form.Get("description"))
				assert.Equal(t, "true", form.Get("required"))
				return mock.NewResponse(fixture.CustomField.SingleJSON), nil
			},
		},
		"success-numeric-opts": {
			projectIDOrKey: "TEST",
			fieldType:      3,
			name:           "Points",
			opts: []*option.APIParamOption{
				o.WithApplicableIssueTypeIDs([]int{1, 2}),
				o.WithMin(1),
				o.WithMax(10),
				o.WithInitialValue(5),
				o.WithUnit("pt"),
			},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"1", "2"}, form["applicableIssueTypes[]"])
				assert.Equal(t, "1", form.Get("min"))
				assert.Equal(t, "10", form.Get("max"))
				assert.Equal(t, "5", form.Get("initialValue"))
				assert.Equal(t, "pt", form.Get("unit"))
				return mock.NewResponse(fixture.CustomField.SingleJSON), nil
			},
		},
		"success-date-opts": {
			projectIDOrKey: "TEST",
			fieldType:      4,
			name:           "Deadline",
			opts: []*option.APIParamOption{
				o.WithInitialValueType(3),
				o.WithInitialDate("2025-01-01"),
				o.WithInitialShift(7),
			},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "3", form.Get("initialValueType"))
				assert.Equal(t, "2025-01-01", form.Get("initialDate"))
				assert.Equal(t, "7", form.Get("initialShift"))
				return mock.NewResponse(fixture.CustomField.SingleJSON), nil
			},
		},
		"success-list-opts": {
			projectIDOrKey: "TEST",
			fieldType:      5,
			name:           "Priority",
			opts: []*option.APIParamOption{
				o.WithItems([]string{"High", "Low"}),
				o.WithAllowInput(true),
				o.WithAllowAddItem(false),
			},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"High", "Low"}, form["items[]"])
				assert.Equal(t, "true", form.Get("allowInput"))
				assert.Equal(t, "false", form.Get("allowAddItem"))
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
			opts:                   []*option.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			fieldType:              0,
			name:                   "",
			opts:                   []*option.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		"error-option-invalid-type-with-valid-values": {
			projectIDOrKey: "TEST",
			fieldType:      1,
			name:           "Sprint",
			opts:           []*option.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &option.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			projectIDOrKey: "",
			fieldType:      0,
			name:           "",
			opts:           []*option.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &option.InvalidOptionKeyError{},
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
				var target *option.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, field)
				var ves validation.Errors
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
	o := &option.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		customFieldID  int
		option         *option.APIParamOption
		opts           []*option.APIParamOption

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
			opts:           []*option.APIParamOption{o.WithRequired(true)},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "true", form.Get("required"))
				return mock.NewResponse(fixture.CustomField.SingleJSON), nil
			},
		},
		"success-applicable-issue-types": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			option:         o.WithApplicableIssueTypeIDs([]int{1, 2}),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"1", "2"}, form["applicableIssueTypes[]"])
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
			opts:                   []*option.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			customFieldID:          0,
			option:                 o.WithName(""),
			opts:                   []*option.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		"error-option-invalid-type-with-valid-values": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			option:         mock.NewInvalidTypeOption(),
			wantErrType:    &option.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			projectIDOrKey: "",
			customFieldID:  0,
			option:         mock.NewInvalidTypeOption(),
			wantErrType:    &option.InvalidOptionKeyError{},
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
				var target *option.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, field)
				var ves validation.Errors
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
				var ves validation.Errors
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
