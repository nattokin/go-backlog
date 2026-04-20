package issue_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/issue"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestIssueService_All(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		opts []core.RequestOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantIDs     []int
	}{
		"success-no-options": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.ListJSON))),
				}, nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-projectIDs": {
			opts: []core.RequestOption{o.WithProjectIDs([]int{10, 20})},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, []string{"10", "20"}, query["projectId[]"])
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.ListJSON))),
				}, nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-keyword": {
			opts: []core.RequestOption{o.WithKeyword("bug")},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "bug", query.Get("keyword"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.ListJSON))),
				}, nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-sort-and-order": {
			opts: []core.RequestOption{
				o.WithIssueSort(model.IssueSortCreated),
				o.WithOrder(model.OrderAsc),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "created", query.Get("sort"))
				assert.Equal(t, "asc", query.Get("order"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.ListJSON))),
				}, nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-count-and-offset": {
			opts: []core.RequestOption{
				o.WithCount(50),
				o.WithOffset(100),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "50", query.Get("count"))
				assert.Equal(t, "100", query.Get("offset"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.ListJSON))),
				}, nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-date-filters": {
			opts: []core.RequestOption{
				o.WithCreatedSince(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				o.WithCreatedUntil(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "2024-01-01", query.Get("createdSince"))
				assert.Equal(t, "2024-12-31", query.Get("createdUntil"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.ListJSON))),
				}, nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-parentChild": {
			opts: []core.RequestOption{o.WithParentChild(1)},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "1", query.Get("parentChild"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.ListJSON))),
				}, nil
			},
			wantIDs: []int{1, 2},
		},
		"error-option-invalid-type": {
			opts:        []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-projectID": {
			opts:        []core.RequestOption{o.WithProjectIDs([]int{0})},
			wantErrType: &core.ValidationError{},
		},
		"error-option-invalid-sort": {
			opts:        []core.RequestOption{o.WithIssueSort("invalid")},
			wantErrType: &core.ValidationError{},
		},
		"error-option-invalid-parentChild": {
			opts:        []core.RequestOption{o.WithParentChild(5)},
			wantErrType: &core.ValidationError{},
		},
		"error-option-invalid-count": {
			opts:        []core.RequestOption{o.WithCount(0)},
			wantErrType: &core.ValidationError{},
		},
		"error-option-invalid-offset": {
			opts:        []core.RequestOption{o.WithOffset(-1)},
			wantErrType: &core.ValidationError{},
		},
		"error-option-set-failed": {
			opts:        []core.RequestOption{mock.NewFailingSetOption(core.ParamKeyword)},
			wantErrType: errors.New(""),
		},
		"error-client-network": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
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

			issues, err := s.All(context.Background(), tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, issues)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, issues)
			assert.Len(t, issues, len(tc.wantIDs))
			for i := range issues {
				assert.Equal(t, tc.wantIDs[i], issues[i].ID)
			}
		})
	}
}

func TestIssueService_Count(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		opts []core.RequestOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantCount   int
	}{
		"success-no-options": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/count", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"count":42}`))),
				}, nil
			},
			wantCount: 42,
		},
		"success-with-projectIDs": {
			opts: []core.RequestOption{o.WithProjectIDs([]int{10, 20})},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/count", spath)
				assert.Equal(t, []string{"10", "20"}, query["projectId[]"])
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"count":5}`))),
				}, nil
			},
			wantCount: 5,
		},
		"success-with-keyword": {
			opts: []core.RequestOption{o.WithKeyword("bug")},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/count", spath)
				assert.Equal(t, "bug", query.Get("keyword"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"count":3}`))),
				}, nil
			},
			wantCount: 3,
		},
		"error-option-invalid-type": {
			opts:        []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-projectID": {
			opts:        []core.RequestOption{o.WithProjectIDs([]int{0})},
			wantErrType: &core.ValidationError{},
		},
		"error-client-network": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
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

			count, err := s.Count(context.Background(), tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Zero(t, count)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantCount, count)
		})
	}
}

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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.SingleJSON))),
				}, nil
			},
			wantID: 1,
		},
		"success-by-id": {
			issueIDOrKey: "1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.SingleJSON))),
				}, nil
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
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
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
				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.SingleJSON))),
				}, nil
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
				o.WithStartDate(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)),
				o.WithDueDate(time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)),
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
				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.SingleJSON))),
				}, nil
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
		"error-response-invalid-json": {
			projectID:   10,
			summary:     "New issue",
			issueTypeID: 2,
			priorityID:  3,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.SingleJSON))),
				}, nil
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.SingleJSON))),
				}, nil
			},
			wantID: 1,
		},
		"success-by-numeric-id": {
			issueIDOrKey: "1",
			option:       o.WithDescription("updated desc"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1", spath)
				assert.Equal(t, "updated desc", form.Get("description"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.SingleJSON))),
				}, nil
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
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			option:       o.WithSummary("x"),
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.SingleJSON))),
				}, nil
			},
			wantID: 1,
		},
		"success-by-id": {
			issueIDOrKey: "1",
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.SingleJSON))),
				}, nil
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
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
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

func TestIssueService_contextPropagation(t *testing.T) {
	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	o := &core.OptionService{}

	cases := []struct {
		name string
		call func(t *testing.T, m *core.Method)
	}{
		{"All", func(t *testing.T, m *core.Method) {
			m.Get = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s := issue.NewService(m)
			s.All(ctx) //nolint:errcheck
		}},
		{"Count", func(t *testing.T, m *core.Method) {
			m.Get = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s := issue.NewService(m)
			s.Count(ctx) //nolint:errcheck
		}},
		{"One", func(t *testing.T, m *core.Method) {
			m.Get = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s := issue.NewService(m)
			s.One(ctx, "PRJ-1") //nolint:errcheck
		}},
		{"Create", func(t *testing.T, m *core.Method) {
			m.Post = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s := issue.NewService(m)
			s.Create(ctx, 10, "summary", 2, 3) //nolint:errcheck
		}},
		{"Update", func(t *testing.T, m *core.Method) {
			m.Patch = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s := issue.NewService(m)
			s.Update(ctx, "PRJ-1", o.WithSummary("x")) //nolint:errcheck
		}},
		{"Delete", func(t *testing.T, m *core.Method) {
			m.Delete = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s := issue.NewService(m)
			s.Delete(ctx, "PRJ-1") //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, mock.NewMethod(t))
		})
	}
}
