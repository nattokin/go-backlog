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

func TestIssueService_One(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantID      int
	}{
		"success-by-key": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1", spath)
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			wantID: 1,
		},
		"success-by-id": {
			issueIDOrKey: "1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1", spath)
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			wantID: 1,
		},
		"error-empty-issueIDOrKey": {
			issueIDOrKey: "",
			wantErrType:  &core.ValidationError{},
		},
		"error-zero-issueIDOrKey": {
			issueIDOrKey: "0",
			wantErrType:  &core.ValidationError{},
		},
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

			s := issue.NewService(method)

			got, err := s.One(context.Background(), tc.issueIDOrKey)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestIssueService_Create(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectID   int
		summary     string
		issueTypeID int
		priorityID  int
		opts        []core.RequestOption

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantID      int
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
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
			wantID: 1,
		},
		"success-with-options": {
			projectID:   10,
			summary:     "New issue",
			issueTypeID: 2,
			priorityID:  3,
			opts: []core.RequestOption{
				o.WithDescription("some description"),
				o.WithAssigneeID(5),
				o.WithStartDate("2024-06-01"),
				o.WithDueDate("2024-06-30"),
			},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "10", form.Get("projectId"))
				assert.Equal(t, "New issue", form.Get("summary"))
				assert.Equal(t, "2", form.Get("issueTypeId"))
				assert.Equal(t, "3", form.Get("priorityId"))
				assert.Equal(t, "some description", form.Get("description"))
				assert.Equal(t, "5", form.Get("assigneeId"))
				assert.Equal(t, "2024-06-01", form.Get("startDate"))
				assert.Equal(t, "2024-06-30", form.Get("dueDate"))
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
			wantID: 1,
		},
		"error-invalid-projectID": {
			projectID:   0,
			summary:     "New issue",
			issueTypeID: 2,
			priorityID:  3,
			wantErrType: &core.ValidationError{},
		},
		"error-empty-summary": {
			projectID:   10,
			summary:     "",
			issueTypeID: 2,
			priorityID:  3,
			wantErrType: &core.ValidationError{},
		},
		"error-invalid-issueTypeID": {
			projectID:   10,
			summary:     "New issue",
			issueTypeID: 0,
			priorityID:  3,
			wantErrType: &core.ValidationError{},
		},
		"error-invalid-priorityID": {
			projectID:   10,
			summary:     "New issue",
			issueTypeID: 2,
			priorityID:  0,
			wantErrType: &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectID:   10,
			summary:     "New issue",
			issueTypeID: 2,
			priorityID:  3,
			opts:        []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},
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

			s := issue.NewService(method)

			got, err := s.Create(context.Background(), tc.projectID, tc.summary, tc.issueTypeID, tc.priorityID, tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestIssueService_Update(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		issueIDOrKey string
		option       core.RequestOption
		opts         []core.RequestOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantID      int
	}{
		"success-summary": {
			issueIDOrKey: "PRJ-1",
			option:       o.WithSummary("Updated summary"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1", spath)
				assert.Equal(t, "Updated summary", form.Get("summary"))
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			wantID: 1,
		},
		"success-with-extra-options": {
			issueIDOrKey: "PRJ-1",
			option:       o.WithSummary("Updated summary"),
			opts: []core.RequestOption{
				o.WithStatusID(2),
				o.WithResolutionID(1),
				o.WithAssigneeID(5),
			},
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1", spath)
				assert.Equal(t, "Updated summary", form.Get("summary"))
				assert.Equal(t, "2", form.Get("statusId"))
				assert.Equal(t, "1", form.Get("resolutionId"))
				assert.Equal(t, "5", form.Get("assigneeId"))
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			wantID: 1,
		},
		"success-by-numeric-id": {
			issueIDOrKey: "1",
			option:       o.WithDescription("updated desc"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1", spath)
				assert.Equal(t, "updated desc", form.Get("description"))
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			wantID: 1,
		},
		"error-empty-issueIDOrKey": {
			issueIDOrKey: "",
			option:       o.WithSummary("x"),
			wantErrType:  &core.ValidationError{},
		},
		"error-zero-issueIDOrKey": {
			issueIDOrKey: "0",
			option:       o.WithSummary("x"),
			wantErrType:  &core.ValidationError{},
		},
		"error-option-invalid-type": {
			issueIDOrKey: "PRJ-1",
			option:       mock.NewInvalidTypeOption(),
			wantErrType:  &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-assigneeID": {
			issueIDOrKey: "PRJ-1",
			option:       o.WithAssigneeID(0),
			wantErrType:  &core.ValidationError{},
		},
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

			s := issue.NewService(method)

			got, err := s.Update(context.Background(), tc.issueIDOrKey, tc.option, tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestIssueService_Delete(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantID      int
	}{
		"success-by-key": {
			issueIDOrKey: "PRJ-1",
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1", spath)
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			wantID: 1,
		},
		"success-by-id": {
			issueIDOrKey: "1",
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1", spath)
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			wantID: 1,
		},
		"error-empty-issueIDOrKey": {
			issueIDOrKey: "",
			wantErrType:  &core.ValidationError{},
		},
		"error-zero-issueIDOrKey": {
			issueIDOrKey: "0",
			wantErrType:  &core.ValidationError{},
		},
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

			s := issue.NewService(method)

			got, err := s.Delete(context.Background(), tc.issueIDOrKey)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}
