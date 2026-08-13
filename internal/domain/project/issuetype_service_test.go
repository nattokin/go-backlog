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

func TestIssueTypeService_List(t *testing.T) {
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
				assert.Equal(t, "projects/TEST/issueTypes", spath)
				return mock.NewResponse(fixture.IssueType.ListJSON), nil
			},
			wantLen: 2,
		},
		"success-id": {
			projectIDOrKey: "6",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/6/issueTypes", spath)
				return mock.NewResponse(fixture.IssueType.ListJSON), nil
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
			s := project.NewIssueTypeService(method)
			issueTypes, err := s.List(context.Background(), tc.projectIDOrKey)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, issueTypes)
				var ves validation.Errors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, issueTypes)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, issueTypes)
			assert.Len(t, issueTypes, tc.wantLen)
		})
	}
}

func TestIssueTypeService_Create(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		name           string
		color          string
		opts           []*option.APIParamOption

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success-minimum": {
			projectIDOrKey: "TEST",
			name:           "Bug",
			color:          "#e30000",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/issueTypes", spath)
				assert.Equal(t, "Bug", form.Get("name"))
				assert.Equal(t, "#e30000", form.Get("color"))
				return mock.NewResponse(fixture.IssueType.SingleJSON), nil
			},
		},
		"success-with-opts": {
			projectIDOrKey: "TEST",
			name:           "Bug",
			color:          "#e30000",
			opts:           []*option.APIParamOption{o.WithTemplateSummary("summary")},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "summary", form.Get("templateSummary"))
				return mock.NewResponse(fixture.IssueType.SingleJSON), nil
			},
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			name:                   "Bug",
			color:                  "#e30000",
			wantValidationErrCount: 1,
		},

		"error-validation-name-empty": {
			projectIDOrKey:         "TEST",
			name:                   "",
			color:                  "#e30000",
			wantValidationErrCount: 1,
		},
		"error-validation-color-empty": {
			projectIDOrKey:         "TEST",
			name:                   "Bug",
			color:                  "",
			wantValidationErrCount: 1,
		},
		"error-validation-name-and-color-empty": {
			projectIDOrKey:         "TEST",
			name:                   "",
			color:                  "",
			wantValidationErrCount: 2,
		},

		"error-validation-all": {
			projectIDOrKey:         "",
			name:                   "",
			color:                  "",
			wantValidationErrCount: 3,
		},

		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "TEST",
			name:                   "Bug",
			color:                  "#e30000",
			opts:                   []*option.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			name:                   "",
			color:                  "",
			opts:                   []*option.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		"error-option-invalid-type-with-valid-values": {
			projectIDOrKey: "TEST",
			name:           "Bug",
			color:          "#e30000",
			opts:           []*option.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &option.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			projectIDOrKey: "",
			name:           "",
			color:          "",
			opts:           []*option.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &option.InvalidOptionKeyError{},
		},

		"error-client-network": {
			projectIDOrKey: "TEST",
			name:           "Bug",
			color:          "#e30000",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			name:           "Bug",
			color:          "#e30000",
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
			s := project.NewIssueTypeService(method)
			issueType, err := s.Create(context.Background(), tc.projectIDOrKey, tc.name, tc.color, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, issueType)
				var target *option.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, issueType)
				var ves validation.Errors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, issueType)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, issueType)
			assert.Equal(t, 1, issueType.ID)
			assert.Equal(t, "Bug", issueType.Name)
			assert.Equal(t, "#e30000", issueType.Color)
		})
	}
}

func TestIssueTypeService_Update(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		issueTypeID    int
		option         *option.APIParamOption
		opts           []*option.APIParamOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success": {
			projectIDOrKey: "TEST",
			issueTypeID:    1,
			option:         o.WithName("Bug Updated"),
			opts:           []*option.APIParamOption{o.WithColor("#990000")},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/issueTypes/1", spath)
				assert.Equal(t, "Bug Updated", form.Get("name"))
				assert.Equal(t, "#990000", form.Get("color"))
				return mock.NewResponse(fixture.IssueType.SingleJSON), nil
			},
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			issueTypeID:            1,
			option:                 o.WithName("Bug Updated"),
			wantValidationErrCount: 1,
		},
		"error-validation-issueTypeID-zero": {
			projectIDOrKey:         "TEST",
			issueTypeID:            0,
			option:                 o.WithName("Bug Updated"),
			wantValidationErrCount: 1,
		},

		"error-validation-fixed-option": {
			projectIDOrKey:         "TEST",
			issueTypeID:            1,
			option:                 o.WithName(""),
			wantValidationErrCount: 1,
		},

		"error-validation-opt-single": {
			projectIDOrKey:         "TEST",
			issueTypeID:            1,
			option:                 o.WithName("Bug"),
			opts:                   []*option.APIParamOption{o.WithColor("")},
			wantValidationErrCount: 1,
		},

		"error-validation-all": {
			projectIDOrKey:         "",
			issueTypeID:            0,
			option:                 o.WithName(""),
			opts:                   []*option.APIParamOption{o.WithColor("")},
			wantValidationErrCount: 4,
		},

		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "TEST",
			issueTypeID:            1,
			option:                 o.WithName("Bug"),
			opts:                   []*option.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			issueTypeID:            0,
			option:                 o.WithName(""),
			opts:                   []*option.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		"error-option-invalid-type-with-valid-values": {
			projectIDOrKey: "TEST",
			issueTypeID:    1,
			option:         mock.NewInvalidTypeOption(),
			wantErrType:    &option.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			projectIDOrKey: "",
			issueTypeID:    0,
			option:         mock.NewInvalidTypeOption(),
			wantErrType:    &option.InvalidOptionKeyError{},
		},

		"error-client-network": {
			projectIDOrKey: "TEST",
			issueTypeID:    1,
			option:         o.WithName("Bug Updated"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			issueTypeID:    1,
			option:         o.WithName("Bug Updated"),
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
			s := project.NewIssueTypeService(method)
			issueType, err := s.Update(context.Background(), tc.projectIDOrKey, tc.issueTypeID, tc.option, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, issueType)
				var target *option.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, issueType)
				var ves validation.Errors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, issueType)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, issueType)
			assert.Equal(t, 1, issueType.ID)
		})
	}
}

func TestIssueTypeService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey        string
		issueTypeID           int
		substituteIssueTypeID int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			projectIDOrKey:        "TEST",
			issueTypeID:           1,
			substituteIssueTypeID: 2,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/issueTypes/1", spath)
				assert.Equal(t, "2", form.Get("substituteIssueTypeId"))
				return mock.NewResponse(fixture.IssueType.SingleJSON), nil
			},
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			issueTypeID:            1,
			substituteIssueTypeID:  2,
			wantValidationErrCount: 1,
		},
		"error-validation-issueTypeID-zero": {
			projectIDOrKey:         "TEST",
			issueTypeID:            0,
			substituteIssueTypeID:  2,
			wantValidationErrCount: 1,
		},
		"error-validation-substituteIssueTypeID-zero": {
			projectIDOrKey:         "TEST",
			issueTypeID:            1,
			substituteIssueTypeID:  0,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			issueTypeID:            0,
			substituteIssueTypeID:  0,
			wantValidationErrCount: 3,
		},

		"error-client-network": {
			projectIDOrKey:        "TEST",
			issueTypeID:           1,
			substituteIssueTypeID: 2,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey:        "TEST",
			issueTypeID:           1,
			substituteIssueTypeID: 2,
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
			s := project.NewIssueTypeService(method)
			issueType, err := s.Delete(context.Background(), tc.projectIDOrKey, tc.issueTypeID, tc.substituteIssueTypeID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, issueType)
				var ves validation.Errors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
				assert.Nil(t, issueType)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, issueType)
			assert.Equal(t, 1, issueType.ID)
		})
	}
}
