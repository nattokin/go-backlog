package backlog_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestProjectCustomFieldService_List(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey         string
		doFunc                 func(*http.Request) (*http.Response, error)
		wantLen                int
		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			doFunc: func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, "/api/v2/projects/TEST/customFields", r.URL.Path)
				return mock.NewResponse(fixture.CustomField.ListJSON), nil
			},
			wantLen: 2,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			wantValidationErrCount: 1,
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			doFunc: func(r *http.Request) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			doFunc: func(r *http.Request) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "apiKey", backlog.WithDoer(&mock.Doer{DoFunc: tc.doFunc}))
			require.NoError(t, err)

			fields, err := c.Project.CustomField.List(context.Background(), tc.projectIDOrKey)

			if tc.wantValidationErrCount > 0 {
				require.Error(t, err)
				assert.Nil(t, fields)
				var target *backlog.ValidationError
				if assert.ErrorAs(t, err, &target) {
					assert.Equal(t, tc.wantValidationErrCount, 1)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, fields)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, fields)
			assert.Len(t, fields, tc.wantLen)
		})
	}
}

func TestProjectCustomFieldService_Create(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey         string
		fieldType              backlog.CustomFieldType
		name                   string
		opts                   []backlog.RequestOption
		doFunc                 func(*http.Request) (*http.Response, error)
		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			fieldType:      backlog.CustomFieldTypeText,
			name:           "Sprint",
			doFunc: func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, "/api/v2/projects/TEST/customFields", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)
				return mock.NewResponse(fixture.CustomField.SingleJSON), nil
			},
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			fieldType:              backlog.CustomFieldTypeText,
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
			fieldType:              backlog.CustomFieldTypeText,
			name:                   "",
			wantValidationErrCount: 1,
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			fieldType:      backlog.CustomFieldTypeText,
			name:           "Sprint",
			doFunc: func(r *http.Request) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			fieldType:      backlog.CustomFieldTypeText,
			name:           "Sprint",
			doFunc: func(r *http.Request) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "apiKey", backlog.WithDoer(&mock.Doer{DoFunc: tc.doFunc}))
			require.NoError(t, err)

			field, err := c.Project.CustomField.Create(context.Background(), tc.projectIDOrKey, tc.fieldType, tc.name, tc.opts...)

			if tc.wantValidationErrCount > 0 {
				require.Error(t, err)
				assert.Nil(t, field)
				var target *backlog.ValidationError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
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

func TestProjectCustomFieldService_Update(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey         string
		customFieldID          int
		opt                    func(c *backlog.Client) backlog.RequestOption
		doFunc                 func(*http.Request) (*http.Response, error)
		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			opt: func(c *backlog.Client) backlog.RequestOption {
				return c.Project.CustomField.Option.WithName("Sprint Updated")
			},
			doFunc: func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, "/api/v2/projects/TEST/customFields/1", r.URL.Path)
				return mock.NewResponse(fixture.CustomField.SingleJSON), nil
			},
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			customFieldID:          1,
			opt:                    func(c *backlog.Client) backlog.RequestOption { return c.Project.CustomField.Option.WithName("x") },
			wantValidationErrCount: 1,
		},
		"error-validation-customFieldID-zero": {
			projectIDOrKey:         "TEST",
			customFieldID:          0,
			opt:                    func(c *backlog.Client) backlog.RequestOption { return c.Project.CustomField.Option.WithName("x") },
			wantValidationErrCount: 1,
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			opt:            func(c *backlog.Client) backlog.RequestOption { return c.Project.CustomField.Option.WithName("x") },
			doFunc: func(r *http.Request) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			opt:            func(c *backlog.Client) backlog.RequestOption { return c.Project.CustomField.Option.WithName("x") },
			doFunc: func(r *http.Request) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "apiKey", backlog.WithDoer(&mock.Doer{DoFunc: tc.doFunc}))
			require.NoError(t, err)

			field, err := c.Project.CustomField.Update(context.Background(), tc.projectIDOrKey, tc.customFieldID, tc.opt(c))

			if tc.wantValidationErrCount > 0 {
				require.Error(t, err)
				assert.Nil(t, field)
				var target *backlog.ValidationError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, field)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, field)
			assert.Equal(t, 1, field.ID)
		})
	}
}

// TestProjectCustomFieldService_Create_readBack verifies that the type-specific settings
// sent through options can be read back from the decoded response.
func TestProjectCustomFieldService_Create_readBack(t *testing.T) {
	cases := map[string]struct {
		fieldType backlog.CustomFieldType
		opts      func(o *backlog.ProjectCustomFieldOptionService) []backlog.RequestOption
		wantForm  map[string]string
		respJSON  string
		want      *backlog.CustomField
	}{
		"number": {
			fieldType: backlog.CustomFieldTypeNumber,
			opts: func(o *backlog.ProjectCustomFieldOptionService) []backlog.RequestOption {
				return []backlog.RequestOption{o.WithMin(0), o.WithMax(100), o.WithInitialValue(8), o.WithUnit("h")}
			},
			wantForm: map[string]string{"min": "0", "max": "100", "initialValue": "8", "unit": "h"},
			respJSON: fixture.CustomField.NumberJSON,
			want:     fixture.CustomField.Number,
		},
		"date": {
			fieldType: backlog.CustomFieldTypeDate,
			opts: func(o *backlog.ProjectCustomFieldOptionService) []backlog.RequestOption {
				return []backlog.RequestOption{
					o.WithMinDate("2024-01-15"), o.WithMaxDate("2024-12-31"), o.WithInitialValueType(2), o.WithInitialShift(7),
				}
			},
			wantForm: map[string]string{"min": "2024-01-15", "max": "2024-12-31", "initialValueType": "2", "initialShift": "7"},
			respJSON: fixture.CustomField.DateJSON,
			want:     fixture.CustomField.Date,
		},
		"checkbox": {
			fieldType: backlog.CustomFieldTypeCheckbox,
			opts: func(o *backlog.ProjectCustomFieldOptionService) []backlog.RequestOption {
				return []backlog.RequestOption{o.WithItems([]string{"Windows"}), o.WithAllowInput(true), o.WithAllowAddItem(true)}
			},
			wantForm: map[string]string{"items[]": "Windows", "allowInput": "true", "allowAddItem": "true"},
			respJSON: fixture.CustomField.CheckBoxJSON,
			want:     fixture.CustomField.CheckBox,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			doFunc := func(r *http.Request) (*http.Response, error) {
				require.NoError(t, r.ParseForm())
				for k, v := range tc.wantForm {
					assert.Equal(t, v, r.PostForm.Get(k), k)
				}
				return mock.NewResponse(tc.respJSON), nil
			}
			c, err := backlog.NewClient("https://example.backlog.com", "apiKey", backlog.WithDoer(&mock.Doer{DoFunc: doFunc}))
			require.NoError(t, err)

			got, err := c.Project.CustomField.Create(context.Background(), "TEST", tc.fieldType, tc.want.Name, tc.opts(c.Project.CustomField.Option)...)

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestProjectCustomFieldService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey         string
		customFieldID          int
		doFunc                 func(*http.Request) (*http.Response, error)
		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			doFunc: func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, "/api/v2/projects/TEST/customFields/1", r.URL.Path)
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
		"error-client-network": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			doFunc: func(r *http.Request) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			doFunc: func(r *http.Request) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "apiKey", backlog.WithDoer(&mock.Doer{DoFunc: tc.doFunc}))
			require.NoError(t, err)

			field, err := c.Project.CustomField.Delete(context.Background(), tc.projectIDOrKey, tc.customFieldID)

			if tc.wantValidationErrCount > 0 {
				require.Error(t, err)
				assert.Nil(t, field)
				var target *backlog.ValidationError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
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

func TestProjectCustomFieldService_AddListItem(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey         string
		customFieldID          int
		name                   string
		doFunc                 func(*http.Request) (*http.Response, error)
		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			name:           "Item1",
			doFunc: func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, "/api/v2/projects/TEST/customFields/1/items", r.URL.Path)
				return mock.NewResponse(fixture.CustomField.SingleJSON), nil
			},
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			customFieldID:          1,
			name:                   "Item1",
			wantValidationErrCount: 1,
		},
		"error-validation-customFieldID-zero": {
			projectIDOrKey:         "TEST",
			customFieldID:          0,
			name:                   "Item1",
			wantValidationErrCount: 1,
		},
		"error-validation-name-empty": {
			projectIDOrKey:         "TEST",
			customFieldID:          1,
			name:                   "",
			wantValidationErrCount: 1,
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			name:           "Item1",
			doFunc: func(r *http.Request) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			name:           "Item1",
			doFunc: func(r *http.Request) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "apiKey", backlog.WithDoer(&mock.Doer{DoFunc: tc.doFunc}))
			require.NoError(t, err)

			field, err := c.Project.CustomField.AddListItem(context.Background(), tc.projectIDOrKey, tc.customFieldID, tc.name)

			if tc.wantValidationErrCount > 0 {
				require.Error(t, err)
				assert.Nil(t, field)
				var target *backlog.ValidationError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, field)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, field)
			assert.Equal(t, 1, field.ID)
		})
	}
}

func TestProjectCustomFieldService_UpdateListItem(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey         string
		customFieldID          int
		itemID                 int
		name                   string
		doFunc                 func(*http.Request) (*http.Response, error)
		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         10,
			name:           "Item1 Updated",
			doFunc: func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, "/api/v2/projects/TEST/customFields/1/items/10", r.URL.Path)
				return mock.NewResponse(fixture.CustomField.SingleJSON), nil
			},
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			customFieldID:          1,
			itemID:                 10,
			name:                   "Item1 Updated",
			wantValidationErrCount: 1,
		},
		"error-validation-customFieldID-zero": {
			projectIDOrKey:         "TEST",
			customFieldID:          0,
			itemID:                 10,
			name:                   "Item1 Updated",
			wantValidationErrCount: 1,
		},
		"error-validation-itemID-zero": {
			projectIDOrKey:         "TEST",
			customFieldID:          1,
			itemID:                 0,
			name:                   "Item1 Updated",
			wantValidationErrCount: 1,
		},
		"error-validation-name-empty": {
			projectIDOrKey:         "TEST",
			customFieldID:          1,
			itemID:                 10,
			name:                   "",
			wantValidationErrCount: 1,
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         10,
			name:           "Item1 Updated",
			doFunc: func(r *http.Request) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         10,
			name:           "Item1 Updated",
			doFunc: func(r *http.Request) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "apiKey", backlog.WithDoer(&mock.Doer{DoFunc: tc.doFunc}))
			require.NoError(t, err)

			field, err := c.Project.CustomField.UpdateListItem(context.Background(), tc.projectIDOrKey, tc.customFieldID, tc.itemID, tc.name)

			if tc.wantValidationErrCount > 0 {
				require.Error(t, err)
				assert.Nil(t, field)
				var target *backlog.ValidationError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, field)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, field)
			assert.Equal(t, 1, field.ID)
		})
	}
}

func TestProjectCustomFieldService_DeleteListItem(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey         string
		customFieldID          int
		itemID                 int
		doFunc                 func(*http.Request) (*http.Response, error)
		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         10,
			doFunc: func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, "/api/v2/projects/TEST/customFields/1/items/10", r.URL.Path)
				return mock.NewResponse(fixture.CustomField.SingleJSON), nil
			},
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			customFieldID:          1,
			itemID:                 10,
			wantValidationErrCount: 1,
		},
		"error-validation-customFieldID-zero": {
			projectIDOrKey:         "TEST",
			customFieldID:          0,
			itemID:                 10,
			wantValidationErrCount: 1,
		},
		"error-validation-itemID-zero": {
			projectIDOrKey:         "TEST",
			customFieldID:          1,
			itemID:                 0,
			wantValidationErrCount: 1,
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         10,
			doFunc: func(r *http.Request) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         10,
			doFunc: func(r *http.Request) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "apiKey", backlog.WithDoer(&mock.Doer{DoFunc: tc.doFunc}))
			require.NoError(t, err)

			field, err := c.Project.CustomField.DeleteListItem(context.Background(), tc.projectIDOrKey, tc.customFieldID, tc.itemID)

			if tc.wantValidationErrCount > 0 {
				require.Error(t, err)
				assert.Nil(t, field)
				var target *backlog.ValidationError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, field)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, field)
			assert.Equal(t, 1, field.ID)
		})
	}
}

func TestProjectCustomFieldOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "apiKey")
	require.NoError(t, err)

	s := c.Project.CustomField.Option

	cases := map[string]struct {
		option  backlog.RequestOption
		wantKey string
	}{
		"WithName": {
			option:  s.WithName("Sprint"),
			wantKey: option.ParamName.Value(),
		},
		"WithDescription": {
			option:  s.WithDescription("sprint number"),
			wantKey: option.ParamDescription.Value(),
		},
		"WithRequired": {
			option:  s.WithRequired(true),
			wantKey: option.ParamRequired.Value(),
		},
		"WithApplicableIssueTypeIDs": {
			option:  s.WithApplicableIssueTypeIDs([]int{1, 2}),
			wantKey: option.ParamApplicableIssueTypeIDs.Value(),
		},
		"WithMin": {
			option:  s.WithMin(0.5),
			wantKey: option.ParamMin.Value(),
		},
		"WithMax": {
			option:  s.WithMax(100.0),
			wantKey: option.ParamMax.Value(),
		},
		"WithInitialValue": {
			option:  s.WithInitialValue(1.0),
			wantKey: option.ParamInitialValue.Value(),
		},
		"WithUnit": {
			option:  s.WithUnit("pt"),
			wantKey: option.ParamUnit.Value(),
		},
		"WithInitialValueType": {
			option:  s.WithInitialValueType(1),
			wantKey: option.ParamInitialValueType.Value(),
		},
		"WithInitialDate": {
			option:  s.WithInitialDate("2024-03-15"),
			wantKey: option.ParamInitialDate.Value(),
		},
		"WithInitialShift": {
			option:  s.WithInitialShift(7),
			wantKey: option.ParamInitialShift.Value(),
		},
		"WithMinDate": {
			option:  s.WithMinDate("2024-01-01"),
			wantKey: option.ParamMin.Value(),
		},
		"WithMaxDate": {
			option:  s.WithMaxDate("2024-12-31"),
			wantKey: option.ParamMax.Value(),
		},
		"WithItems": {
			option:  s.WithItems([]string{"High", "Low"}),
			wantKey: option.ParamItems.Value(),
		},
		"WithAllowInput": {
			option:  s.WithAllowInput(true),
			wantKey: option.ParamAllowInput.Value(),
		},
		"WithAllowAddItem": {
			option:  s.WithAllowAddItem(true),
			wantKey: option.ParamAllowAddItem.Value(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}
