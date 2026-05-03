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
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestProjectCustomFieldService_All(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		doFunc         func(*http.Request) (*http.Response, error)
		wantLen        int
		wantErrType    error
	}{
		"success": {
			projectIDOrKey: "TEST",
			doFunc: func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, "/api/v2/projects/TEST/customFields", r.URL.Path)
				return mock.NewJSONResponse(fixture.CustomField.ListJSON), nil
			},
			wantLen:     2,
			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &backlog.ValidationError{},
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)

			fields, err := c.Project.CustomField.All(context.Background(), tc.projectIDOrKey)

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
		projectIDOrKey string
		fieldType      backlog.CustomFieldType
		name           string
		opts           []backlog.RequestOption
		doFunc         func(*http.Request) (*http.Response, error)
		wantErrType    error
	}{
		"success": {
			projectIDOrKey: "TEST",
			fieldType:      backlog.CustomFieldTypeText,
			name:           "Sprint",
			doFunc: func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, "/api/v2/projects/TEST/customFields", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)
				return mock.NewJSONResponse(fixture.CustomField.SingleJSON), nil
			},
			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			fieldType:      backlog.CustomFieldTypeText,
			name:           "Sprint",
			wantErrType:    &backlog.ValidationError{},
		},
		"error-validation-fieldType-zero": {
			projectIDOrKey: "TEST",
			fieldType:      0,
			name:           "Sprint",
			wantErrType:    &backlog.ValidationError{},
		},
		"error-validation-name-empty": {
			projectIDOrKey: "TEST",
			fieldType:      backlog.CustomFieldTypeText,
			name:           "",
			wantErrType:    &backlog.ValidationError{},
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)

			field, err := c.Project.CustomField.Create(context.Background(), tc.projectIDOrKey, tc.fieldType, tc.name, tc.opts...)

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
		projectIDOrKey string
		customFieldID  int
		opt            func(c *backlog.Client) backlog.RequestOption
		doFunc         func(*http.Request) (*http.Response, error)
		wantErrType    error
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			opt: func(c *backlog.Client) backlog.RequestOption {
				return c.Project.CustomField.Option.WithName("Sprint Updated")
			},
			doFunc: func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, "/api/v2/projects/TEST/customFields/1", r.URL.Path)
				return mock.NewJSONResponse(fixture.CustomField.SingleJSON), nil
			},
			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			customFieldID:  1,
			opt:            func(c *backlog.Client) backlog.RequestOption { return c.Project.CustomField.Option.WithName("x") },
			wantErrType:    &backlog.ValidationError{},
		},
		"error-validation-customFieldID-zero": {
			projectIDOrKey: "TEST",
			customFieldID:  0,
			opt:            func(c *backlog.Client) backlog.RequestOption { return c.Project.CustomField.Option.WithName("x") },
			wantErrType:    &backlog.ValidationError{},
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)

			field, err := c.Project.CustomField.Update(context.Background(), tc.projectIDOrKey, tc.customFieldID, tc.opt(c))

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

func TestProjectCustomFieldService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		customFieldID  int
		doFunc         func(*http.Request) (*http.Response, error)
		wantErrType    error
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			doFunc: func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, "/api/v2/projects/TEST/customFields/1", r.URL.Path)
				return mock.NewJSONResponse(fixture.CustomField.SingleJSON), nil
			},
			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			customFieldID:  1,
			wantErrType:    &backlog.ValidationError{},
		},
		"error-validation-customFieldID-zero": {
			projectIDOrKey: "TEST",
			customFieldID:  0,
			wantErrType:    &backlog.ValidationError{},
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)

			field, err := c.Project.CustomField.Delete(context.Background(), tc.projectIDOrKey, tc.customFieldID)

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
		projectIDOrKey string
		customFieldID  int
		name           string
		doFunc         func(*http.Request) (*http.Response, error)
		wantErrType    error
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			name:           "Item1",
			doFunc: func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, "/api/v2/projects/TEST/customFields/1/items", r.URL.Path)
				return mock.NewJSONResponse(fixture.CustomField.SingleJSON), nil
			},
			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			customFieldID:  1,
			name:           "Item1",
			wantErrType:    &backlog.ValidationError{},
		},
		"error-validation-customFieldID-zero": {
			projectIDOrKey: "TEST",
			customFieldID:  0,
			name:           "Item1",
			wantErrType:    &backlog.ValidationError{},
		},
		"error-validation-name-empty": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			name:           "",
			wantErrType:    &backlog.ValidationError{},
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)

			field, err := c.Project.CustomField.AddListItem(context.Background(), tc.projectIDOrKey, tc.customFieldID, tc.name)

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
		projectIDOrKey string
		customFieldID  int
		itemID         int
		name           string
		doFunc         func(*http.Request) (*http.Response, error)
		wantErrType    error
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         10,
			name:           "Item1 Updated",
			doFunc: func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, "/api/v2/projects/TEST/customFields/1/items/10", r.URL.Path)
				return mock.NewJSONResponse(fixture.CustomField.SingleJSON), nil
			},
			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			customFieldID:  1,
			itemID:         10,
			name:           "Item1 Updated",
			wantErrType:    &backlog.ValidationError{},
		},
		"error-validation-customFieldID-zero": {
			projectIDOrKey: "TEST",
			customFieldID:  0,
			itemID:         10,
			name:           "Item1 Updated",
			wantErrType:    &backlog.ValidationError{},
		},
		"error-validation-itemID-zero": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         0,
			name:           "Item1 Updated",
			wantErrType:    &backlog.ValidationError{},
		},
		"error-validation-name-empty": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         10,
			name:           "",
			wantErrType:    &backlog.ValidationError{},
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)

			field, err := c.Project.CustomField.UpdateListItem(context.Background(), tc.projectIDOrKey, tc.customFieldID, tc.itemID, tc.name)

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
		projectIDOrKey string
		customFieldID  int
		itemID         int
		doFunc         func(*http.Request) (*http.Response, error)
		wantErrType    error
	}{
		"success": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         10,
			doFunc: func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, "/api/v2/projects/TEST/customFields/1/items/10", r.URL.Path)
				return mock.NewJSONResponse(fixture.CustomField.SingleJSON), nil
			},
			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			customFieldID:  1,
			itemID:         10,
			wantErrType:    &backlog.ValidationError{},
		},
		"error-validation-customFieldID-zero": {
			projectIDOrKey: "TEST",
			customFieldID:  0,
			itemID:         10,
			wantErrType:    &backlog.ValidationError{},
		},
		"error-validation-itemID-zero": {
			projectIDOrKey: "TEST",
			customFieldID:  1,
			itemID:         0,
			wantErrType:    &backlog.ValidationError{},
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
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)

			field, err := c.Project.CustomField.DeleteListItem(context.Background(), tc.projectIDOrKey, tc.customFieldID, tc.itemID)

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
