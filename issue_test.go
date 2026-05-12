package backlog_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestIssueService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"All": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues", req.URL.Path)
				return mock.NewJSONResponse(fixture.Issue.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.All(ctx)
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 1, got[0].ID)
				assert.Equal(t, 2, got[1].ID)
			},
		},
		"All/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues", req.URL.Path)
				assert.Equal(t, "bug", req.URL.Query().Get("keyword"))
				assert.Equal(t, []string{"10", "20"}, req.URL.Query()["projectId[]"])
				return mock.NewJSONResponse(fixture.Issue.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.All(ctx,
					c.Issue.Option.WithKeyword("bug"),
					c.Issue.Option.WithProjectIDs([]int{10, 20}),
				)
				require.NoError(t, err)
				assert.Len(t, got, 2)
			},
		},
		"All/error": {
			doFunc: newInternalServerErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.All(ctx)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Count": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues/count", req.URL.Path)
				return mock.NewJSONResponse(`{"count":5}`), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Count(ctx)
				require.NoError(t, err)
				assert.Equal(t, 5, got)
			},
		},
		"Count/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues/count", req.URL.Path)
				assert.Equal(t, "bug", req.URL.Query().Get("keyword"))
				assert.Equal(t, []string{"10", "20"}, req.URL.Query()["projectId[]"])
				return mock.NewJSONResponse(`{"count":2}`), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Count(ctx,
					c.Issue.Option.WithKeyword("bug"),
					c.Issue.Option.WithProjectIDs([]int{10, 20}),
				)
				require.NoError(t, err)
				assert.Equal(t, 2, got)
			},
		},
		"Count/error": {
			doFunc: newInternalServerErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Count(ctx)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"One": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues/PRJ-1", req.URL.Path)
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.One(ctx, "PRJ-1")
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
				assert.Equal(t, "PRJ-1", got.IssueKey)
				assert.Equal(t, "First issue", got.Summary)
			},
		},
		"One/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.One(ctx, "PRJ-1")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Create": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/issues", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "10", req.PostForm.Get("projectId"))
				assert.Equal(t, "New issue", req.PostForm.Get("summary"))
				assert.Equal(t, "2", req.PostForm.Get("issueTypeId"))
				assert.Equal(t, "3", req.PostForm.Get("priorityId"))
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Create(ctx, 10, "New issue", 2, 3)
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
				assert.Equal(t, "First issue", got.Summary)
			},
		},
		"Create/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/issues", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "10", req.PostForm.Get("projectId"))
				assert.Equal(t, "New issue", req.PostForm.Get("summary"))
				assert.Equal(t, "2", req.PostForm.Get("issueTypeId"))
				assert.Equal(t, "3", req.PostForm.Get("priorityId"))
				assert.Equal(t, "details here", req.PostForm.Get("description"))
				assert.Equal(t, "5", req.PostForm.Get("assigneeId"))
				return mock.NewCreatedJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Create(ctx, 10, "New issue", 2, 3,
					c.Issue.Option.WithDescription("details here"),
					c.Issue.Option.WithAssigneeID(5),
				)
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
			},
		},
		"Create/error": {
			doFunc: newAuthErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Create(ctx, 10, "New issue", 2, 3)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Update": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/issues/PRJ-1", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "Updated summary", req.PostForm.Get("summary"))
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Update(ctx, "PRJ-1", c.Issue.Option.WithSummary("Updated summary"))
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
			},
		},
		"Update/with-options": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPatch, req.Method)
				assert.Equal(t, "/api/v2/issues/PRJ-1", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "Updated summary", req.PostForm.Get("summary"))
				assert.Equal(t, "2", req.PostForm.Get("statusId"))
				assert.Equal(t, "1", req.PostForm.Get("resolutionId"))
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Update(ctx, "PRJ-1",
					c.Issue.Option.WithSummary("Updated summary"),
					c.Issue.Option.WithStatusID(2),
					c.Issue.Option.WithResolutionID(1),
				)
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
			},
		},
		"Update/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Update(ctx, "PRJ-1", c.Issue.Option.WithSummary("Updated summary"))
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Delete": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/issues/PRJ-1", req.URL.Path)
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Delete(ctx, "PRJ-1")
				require.NoError(t, err)
				assert.Equal(t, 1, got.ID)
				assert.Equal(t, "PRJ-1", got.IssueKey)
			},
		},
		"Delete/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Delete(ctx, "PRJ-1")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Participants": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues/PRJ-1/participants", req.URL.Path)
				return mock.NewJSONResponse(fixture.User.ListJSON), nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Participants(ctx, "PRJ-1")
				require.NoError(t, err)
				assert.Len(t, got, 4)
				assert.Equal(t, 1, got[0].ID)
				assert.Equal(t, "admin", got[0].Name)
			},
		},
		"Participants/error": {
			doFunc: newNotFoundDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Participants(ctx, "PRJ-1")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)
			tc.call(t, c)
		})
	}
}

func TestIssueStarService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"Add": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/api/v2/stars", req.URL.Path)
				require.NoError(t, req.ParseForm())
				assert.Equal(t, "1", req.FormValue("issueId"))
				return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Issue.Star.Add(ctx, 1)
				require.NoError(t, err)
			},
		},
		"Add/error": {
			doFunc: newAuthErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Issue.Star.Add(ctx, 1)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Remove": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/stars", req.URL.Path)
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				form, err := url.ParseQuery(string(body))
				require.NoError(t, err)
				assert.Equal(t, "42", form.Get("id"))
				return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Issue.Star.Remove(ctx, 42)
				require.NoError(t, err)
			},
		},
		"Remove/error": {
			doFunc: newAuthErrorDoFunc(),
			call: func(t *testing.T, c *backlog.Client) {
				err := c.Issue.Star.Remove(ctx, 42)
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := backlog.NewClient("https://example.backlog.com", "token", backlog.WithDoer(&mockDoer{do: tc.doFunc}))
			require.NoError(t, err)
			tc.call(t, c)
		})
	}
}

func TestIssueOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	s := c.Issue.Option

	date := "2024-01-01"

	cases := map[string]struct {
		option  backlog.RequestOption
		wantKey string
	}{
		"WithActualHours":      {option: s.WithActualHours(1.0), wantKey: core.ParamActualHours.Value()},
		"WithAssigneeID":       {option: s.WithAssigneeID(1), wantKey: core.ParamAssigneeID.Value()},
		"WithAssigneeIDs":      {option: s.WithAssigneeIDs([]int{1}), wantKey: core.ParamAssigneeIDs.Value()},
		"WithAttachment":       {option: s.WithAttachment(true), wantKey: core.ParamAttachment.Value()},
		"WithAttachmentIDs":    {option: s.WithAttachmentIDs([]int{1}), wantKey: core.ParamAttachmentIDs.Value()},
		"WithCategoryIDs":      {option: s.WithCategoryIDs([]int{1}), wantKey: core.ParamCategoryIDs.Value()},
		"WithComment":          {option: s.WithComment("comment"), wantKey: core.ParamComment.Value()},
		"WithCount":            {option: s.WithCount(20), wantKey: core.ParamCount.Value()},
		"WithCreatedSince":     {option: s.WithCreatedSince(date), wantKey: core.ParamCreatedSince.Value()},
		"WithCreatedUntil":     {option: s.WithCreatedUntil(date), wantKey: core.ParamCreatedUntil.Value()},
		"WithCreatedUserIDs":   {option: s.WithCreatedUserIDs([]int{1}), wantKey: core.ParamCreatedUserIDs.Value()},
		"WithCustomFieldItem":  {option: s.WithCustomFieldItem(1, 10), wantKey: "customField"},
		"WithCustomFieldOther": {option: s.WithCustomFieldOther(1, "other"), wantKey: "customField"},
		"WithDescription":      {option: s.WithDescription("desc"), wantKey: core.ParamDescription.Value()},
		"WithDueDate":          {option: s.WithDueDate(date), wantKey: core.ParamDueDate.Value()},
		"WithDueDateSince":     {option: s.WithDueDateSince(date), wantKey: core.ParamDueDateSince.Value()},
		"WithDueDateUntil":     {option: s.WithDueDateUntil(date), wantKey: core.ParamDueDateUntil.Value()},
		"WithEstimatedHours":   {option: s.WithEstimatedHours(1.0), wantKey: core.ParamEstimatedHours.Value()},
		"WithHasDueDate":       {option: s.WithHasDueDate(false), wantKey: core.ParamHasDueDate.Value()},
		"WithIDs":              {option: s.WithIDs([]int{1}), wantKey: core.ParamIDs.Value()},
		"WithIssueSort":        {option: s.WithIssueSort(backlog.IssueSortCreated), wantKey: core.ParamSort.Value()},
		"WithIssueTypeID":      {option: s.WithIssueTypeID(1), wantKey: core.ParamIssueTypeID.Value()},
		"WithIssueTypeIDs":     {option: s.WithIssueTypeIDs([]int{1}), wantKey: core.ParamIssueTypeIDs.Value()},
		"WithKeyword":          {option: s.WithKeyword("bug"), wantKey: core.ParamKeyword.Value()},
		"WithMilestoneIDs":     {option: s.WithMilestoneIDs([]int{1}), wantKey: core.ParamMilestoneIDs.Value()},
		"WithNotifiedUserIDs":  {option: s.WithNotifiedUserIDs([]int{1}), wantKey: core.ParamNotifiedUserIDs.Value()},
		"WithOffset":           {option: s.WithOffset(0), wantKey: core.ParamOffset.Value()},
		"WithOrder":            {option: s.WithOrder(backlog.OrderAsc), wantKey: core.ParamOrder.Value()},
		"WithParentChild":      {option: s.WithParentChild(0), wantKey: core.ParamParentChild.Value()},
		"WithParentIssueID":    {option: s.WithParentIssueID(1), wantKey: core.ParamParentIssueID.Value()},
		"WithParentIssueIDs":   {option: s.WithParentIssueIDs([]int{1}), wantKey: core.ParamParentIssueIDs.Value()},
		"WithPriorityID":       {option: s.WithPriorityID(1), wantKey: core.ParamPriorityID.Value()},
		"WithPriorityIDs":      {option: s.WithPriorityIDs([]int{1}), wantKey: core.ParamPriorityIDs.Value()},
		"WithProjectIDs":       {option: s.WithProjectIDs([]int{1}), wantKey: core.ParamProjectIDs.Value()},
		"WithResolutionID":     {option: s.WithResolutionID(1), wantKey: core.ParamResolutionID.Value()},
		"WithResolutionIDs":    {option: s.WithResolutionIDs([]int{1}), wantKey: core.ParamResolutionIDs.Value()},
		"WithSharedFile":       {option: s.WithSharedFile(true), wantKey: core.ParamSharedFile.Value()},
		"WithStartDate":        {option: s.WithStartDate(date), wantKey: core.ParamStartDate.Value()},
		"WithStartDateSince":   {option: s.WithStartDateSince(date), wantKey: core.ParamStartDateSince.Value()},
		"WithStartDateUntil":   {option: s.WithStartDateUntil(date), wantKey: core.ParamStartDateUntil.Value()},
		"WithStatusID":         {option: s.WithStatusID(1), wantKey: core.ParamStatusID.Value()},
		"WithStatusIDs":        {option: s.WithStatusIDs([]int{1}), wantKey: core.ParamStatusIDs.Value()},
		"WithSummary":          {option: s.WithSummary("summary"), wantKey: core.ParamSummary.Value()},
		"WithUpdatedSince":     {option: s.WithUpdatedSince(date), wantKey: core.ParamUpdatedSince.Value()},
		"WithUpdatedUntil":     {option: s.WithUpdatedUntil(date), wantKey: core.ParamUpdatedUntil.Value()},
		"WithVersionIDs":       {option: s.WithVersionIDs([]int{1}), wantKey: core.ParamVersionIDs.Value()},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}
