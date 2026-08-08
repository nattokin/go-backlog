package issue_test

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
	"github.com/nattokin/go-backlog/internal/domain/issue"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestService_One(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantID                 int
	}{
		"success-by-key": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1", spath)
				return mock.NewResponse(fixture.Issue.SingleJSON), nil
			},
			wantID: 1,
		},
		"success-by-id": {
			issueIDOrKey: "1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1", spath)
				return mock.NewResponse(fixture.Issue.SingleJSON), nil
			},
			wantID: 1,
		},

		// --- validation errors ---
		"error-validation-issueIDOrKey-empty": {
			issueIDOrKey:           "",
			wantValidationErrCount: 1,
		},
		"error-validation-issueIDOrKey-zero": {
			issueIDOrKey:           "0",
			wantValidationErrCount: 1,
		},

		// --- other errors ---
		"error-client-network": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, &core.APIResponseError{}
			},
			wantErrType: &core.APIResponseError{},
		},
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
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
			s := issue.NewService(method)
			got, err := s.One(context.Background(), tc.issueIDOrKey)

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
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestService_Create(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectID   int
		summary     string
		issueTypeID int
		priorityID  int
		opts        []*core.APIParamOption

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
		wantID                 int
	}{
		"success-required-only": {
			projectID:   10,
			summary:     "New issue",
			issueTypeID: 2,
			priorityID:  3,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "10", form.Get("projectId"))
				assert.Equal(t, "New issue", form.Get("summary"))
				assert.Equal(t, "2", form.Get("issueTypeId"))
				assert.Equal(t, "3", form.Get("priorityId"))
				return mock.NewCreatedResponse(fixture.Issue.SingleJSON), nil
			},
			wantID: 1,
		},
		"success-with-options": {
			projectID:   10,
			summary:     "New issue",
			issueTypeID: 2,
			priorityID:  3,
			opts: []*core.APIParamOption{
				o.WithDescription("some description"),
				o.WithAssigneeID(5),
			},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "some description", form.Get("description"))
				assert.Equal(t, "5", form.Get("assigneeId"))
				return mock.NewCreatedResponse(fixture.Issue.SingleJSON), nil
			},
			wantID: 1,
		},

		// --- validation errors: argument only ---
		"error-validation-projectID-zero": {
			projectID:              0,
			summary:                "New issue",
			issueTypeID:            2,
			priorityID:             3,
			wantValidationErrCount: 1,
		},
		"error-validation-summary-empty": {
			projectID:              10,
			summary:                "",
			issueTypeID:            2,
			priorityID:             3,
			wantValidationErrCount: 1,
		},
		"error-validation-issueTypeID-zero": {
			projectID:              10,
			summary:                "New issue",
			issueTypeID:            0,
			priorityID:             3,
			wantValidationErrCount: 1,
		},
		"error-validation-priorityID-zero": {
			projectID:              10,
			summary:                "New issue",
			issueTypeID:            2,
			priorityID:             0,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectID:              0,
			summary:                "",
			issueTypeID:            0,
			priorityID:             0,
			wantValidationErrCount: 4,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			projectID:              10,
			summary:                "New issue",
			issueTypeID:            2,
			priorityID:             3,
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectID:              0,
			summary:                "",
			issueTypeID:            0,
			priorityID:             0,
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type-with-valid-values": {
			projectID:   10,
			summary:     "New issue",
			issueTypeID: 2,
			priorityID:  3,
			opts:        []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			projectID:   0,
			summary:     "",
			issueTypeID: 0,
			priorityID:  0,
			opts:        []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-client-network": {
			projectID:   10,
			summary:     "New issue",
			issueTypeID: 2,
			priorityID:  3,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			projectID:   10,
			summary:     "New issue",
			issueTypeID: 2,
			priorityID:  3,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, &core.APIResponseError{}
			},
			wantErrType: &core.APIResponseError{},
		},
		"error-response-invalid-json": {
			projectID:   10,
			summary:     "New issue",
			issueTypeID: 2,
			priorityID:  3,
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
			s := issue.NewService(method)
			got, err := s.Create(context.Background(), tc.projectID, tc.summary, tc.issueTypeID, tc.priorityID, tc.opts...)

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
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestService_Update(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		issueIDOrKey string
		option       *core.APIParamOption
		opts         []*core.APIParamOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
		wantID                 int
	}{
		"success-summary": {
			issueIDOrKey: "PRJ-1",
			option:       o.WithSummary("Updated summary"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1", spath)
				assert.Equal(t, "Updated summary", form.Get("summary"))
				return mock.NewResponse(fixture.Issue.SingleJSON), nil
			},
			wantID: 1,
		},
		"success-with-extra-options": {
			issueIDOrKey: "PRJ-1",
			option:       o.WithSummary("Updated summary"),
			opts: []*core.APIParamOption{
				o.WithStatusID(2),
				o.WithResolutionID(1),
				o.WithAssigneeID(5),
			},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1", spath)
				return mock.NewResponse(fixture.Issue.SingleJSON), nil
			},
			wantID: 1,
		},

		// --- validation errors: argument only ---
		"error-validation-issueIDOrKey-empty": {
			issueIDOrKey:           "",
			option:                 o.WithSummary("x"),
			wantValidationErrCount: 1,
		},
		"error-validation-issueIDOrKey-zero": {
			issueIDOrKey:           "0",
			option:                 o.WithSummary("x"),
			wantValidationErrCount: 1,
		},

		// --- validation errors: fixed option only ---
		"error-validation-fixed-option": {
			issueIDOrKey:           "PRJ-1",
			option:                 o.WithAssigneeID(0),
			wantValidationErrCount: 1,
		},

		// --- validation errors: all ---
		"error-validation-all": {
			issueIDOrKey:           "",
			option:                 o.WithAssigneeID(0),
			wantValidationErrCount: 2,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			issueIDOrKey:           "PRJ-1",
			option:                 o.WithSummary("x"),
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			issueIDOrKey:           "",
			option:                 o.WithSummary(""),
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type": {
			issueIDOrKey: "PRJ-1",
			option:       mock.NewInvalidTypeOption(),
			wantErrType:  &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-client-network": {
			issueIDOrKey: "PRJ-1",
			option:       o.WithSummary("x"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			issueIDOrKey: "PRJ-1",
			option:       o.WithSummary("x"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, &core.APIResponseError{}
			},
			wantErrType: &core.APIResponseError{},
		},
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			option:       o.WithSummary("x"),
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
			s := issue.NewService(method)
			got, err := s.Update(context.Background(), tc.issueIDOrKey, tc.option, tc.opts...)

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
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestService_Delete(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantID                 int
	}{
		"success-by-key": {
			issueIDOrKey: "PRJ-1",
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1", spath)
				return mock.NewResponse(fixture.Issue.SingleJSON), nil
			},
			wantID: 1,
		},

		// --- validation errors ---
		"error-validation-issueIDOrKey-empty": {
			issueIDOrKey:           "",
			wantValidationErrCount: 1,
		},
		"error-validation-issueIDOrKey-zero": {
			issueIDOrKey:           "0",
			wantValidationErrCount: 1,
		},

		// --- other errors ---
		"error-client-network": {
			issueIDOrKey: "PRJ-1",
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			issueIDOrKey: "PRJ-1",
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, &core.APIResponseError{}
			},
			wantErrType: &core.APIResponseError{},
		},
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
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
			s := issue.NewService(method)
			got, err := s.Delete(context.Background(), tc.issueIDOrKey)

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
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}
