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

		wantLen     int
		wantErrType error
	}{
		"success-projectIDOrKey-key": {
			projectIDOrKey: "TEST",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/customFields", spath)
				assert.Nil(t, query)
				return mock.NewJSONResponse(fixture.CustomField.ListJSON), nil
			},

			wantLen:     2,
			wantErrType: nil,
		},
		"success-projectIDOrKey-id": {
			projectIDOrKey: "6",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/6/customFields", spath)
				return mock.NewJSONResponse(fixture.CustomField.ListJSON), nil
			},

			wantLen:     2,
			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &core.ValidationError{},
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
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

func TestCustomFieldService_Create(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		fieldType      int
		name           string
		opts           []core.RequestOption

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			fieldType:      1,
			name:           "Sprint",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/customFields", spath)
				assert.Equal(t, "1", form.Get("typeId"))
				assert.Equal(t, "Sprint", form.Get("name"))
				assert.Empty(t, form.Get("description"))
				return mock.NewJSONResponse(fixture.CustomField.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"success-with-opts": {
			projectIDOrKey: "TEST",
			fieldType:      1,
			name:           "Sprint",
			opts:           []core.RequestOption{o.WithDescription("sprint number"), o.WithRequired(true)},

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/customFields", spath)
				assert.Equal(t, "1", form.Get("typeId"))
				assert.Equal(t, "Sprint", form.Get("name"))
				assert.Equal(t, "sprint number", form.Get("description"))
				assert.Equal(t, "true", form.Get("required"))
				return mock.NewJSONResponse(fixture.CustomField.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			fieldType:      1,
			name:           "Sprint",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-fieldType-zero": {
			projectIDOrKey: "TEST",
			fieldType:      0,
			name:           "Sprint",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-name-empty": {
			projectIDOrKey: "TEST",
			fieldType:      1,
			name:           "",
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: "TEST",
			fieldType:      1,
			name:           "Sprint",
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
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

func TestCustomFieldService_Update(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		customFieldID  int
		option         core.RequestOption
		opts           []core.RequestOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			option:         o.WithName("Sprint Updated"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/customFields/1", spath)
				assert.Equal(t, "Sprint Updated", form.Get("name"))
				return mock.NewJSONResponse(fixture.CustomField.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"success-with-opts": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			option:         o.WithName("Sprint Updated"),
			opts:           []core.RequestOption{o.WithRequired(true)},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/customFields/1", spath)
				assert.Equal(t, "Sprint Updated", form.Get("name"))
				assert.Equal(t, "true", form.Get("required"))
				return mock.NewJSONResponse(fixture.CustomField.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			customFieldID:  1,
			option:         o.WithName("Sprint Updated"),
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-customFieldID-zero": {
			projectIDOrKey: "TEST",
			customFieldID:  0,
			option:         o.WithName("Sprint Updated"),
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
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

func TestCustomFieldService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		customFieldID  int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/customFields/1", spath)
				assert.NotNil(t, form)
				return mock.NewJSONResponse(fixture.CustomField.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			customFieldID:  1,
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-customFieldID-zero": {
			projectIDOrKey: "TEST",
			customFieldID:  0,
			wantErrType:    &core.ValidationError{},
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
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

func TestCustomFieldService_AddListItem(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		customFieldID  int
		name           string

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			name:           "Item1",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/customFields/1/items", spath)
				assert.Equal(t, "Item1", form.Get("name"))
				return mock.NewJSONResponse(fixture.CustomField.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			customFieldID:  1,
			name:           "Item1",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-customFieldID-zero": {
			projectIDOrKey: "TEST",
			customFieldID:  0,
			name:           "Item1",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-name-empty": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			name:           "",
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			name:           "Item1",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			name:           "Item1",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
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

			field, err := s.AddListItem(context.Background(), tc.projectIDOrKey, tc.customFieldID, tc.name)

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

func TestCustomFieldService_UpdateListItem(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		customFieldID  int
		itemID         int
		name           string

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         10,
			name:           "Item1 Updated",

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/customFields/1/items/10", spath)
				assert.Equal(t, "Item1 Updated", form.Get("name"))
				return mock.NewJSONResponse(fixture.CustomField.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			customFieldID:  1,
			itemID:         10,
			name:           "Item1 Updated",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-customFieldID-zero": {
			projectIDOrKey: "TEST",
			customFieldID:  0,
			itemID:         10,
			name:           "Item1 Updated",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-itemID-zero": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         0,
			name:           "Item1 Updated",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-name-empty": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         10,
			name:           "",
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         10,
			name:           "Item1 Updated",

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         10,
			name:           "Item1 Updated",

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
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

			field, err := s.UpdateListItem(context.Background(), tc.projectIDOrKey, tc.customFieldID, tc.itemID, tc.name)

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

func TestCustomFieldService_DeleteListItem(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		customFieldID  int
		itemID         int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         10,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/customFields/1/items/10", spath)
				assert.NotNil(t, form)
				return mock.NewJSONResponse(fixture.CustomField.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			customFieldID:  1,
			itemID:         10,
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-customFieldID-zero": {
			projectIDOrKey: "TEST",
			customFieldID:  0,
			itemID:         10,
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-itemID-zero": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         0,
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         10,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         10,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
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

			field, err := s.DeleteListItem(context.Background(), tc.projectIDOrKey, tc.customFieldID, tc.itemID)

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
