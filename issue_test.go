package backlog_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

func TestIssueService_All(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"All": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.ListJSON))),
				}, nil
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.ListJSON))),
				}, nil
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
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"Internal Server Error","code":1,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.All(ctx)
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

func TestIssueAttachmentService(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		doFunc func(req *http.Request) (*http.Response, error)
		call   func(t *testing.T, c *backlog.Client)
	}{
		"List": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/api/v2/issues/TEST-1/attachments", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Attachment.ListJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Attachment.List(ctx, "TEST-1")
				require.NoError(t, err)
				assert.Len(t, got, 2)
				assert.Equal(t, 2, got[0].ID)
				assert.Equal(t, 5, got[1].ID)
			},
		},
		"List/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such issue.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Attachment.List(ctx, "TEST-1")
				require.Error(t, err)
				var target *backlog.APIResponseError
				assert.True(t, errors.As(err, &target))
			},
		},
		"Remove": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/api/v2/issues/TEST-1/attachments/8", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Attachment.SingleJSON))),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				got, err := c.Issue.Attachment.Remove(ctx, "TEST-1", 8)
				require.NoError(t, err)
				assert.Equal(t, 8, got.ID)
			},
		},
		"Remove/error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader(`{"errors":[{"message":"No such attachment.","code":6,"moreInfo":""}]}`)),
				}, nil
			},
			call: func(t *testing.T, c *backlog.Client) {
				_, err := c.Issue.Attachment.Remove(ctx, "TEST-1", 8)
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

	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	cases := map[string]struct {
		option  backlog.RequestOption
		wantKey string
	}{
		"WithProjectIDs":      {option: s.WithProjectIDs([]int{1}), wantKey: core.ParamProjectIDs.Value()},
		"WithIssueTypeIDs":    {option: s.WithIssueTypeIDs([]int{1}), wantKey: core.ParamIssueTypeIDs.Value()},
		"WithCategoryIDs":     {option: s.WithCategoryIDs([]int{1}), wantKey: core.ParamCategoryIDs.Value()},
		"WithVersionIDs":      {option: s.WithVersionIDs([]int{1}), wantKey: core.ParamVersionIDs.Value()},
		"WithMilestoneIDs":    {option: s.WithMilestoneIDs([]int{1}), wantKey: core.ParamMilestoneIDs.Value()},
		"WithStatusIDs":       {option: s.WithStatusIDs([]int{1}), wantKey: core.ParamStatusIDs.Value()},
		"WithPriorityIDs":     {option: s.WithPriorityIDs([]int{1}), wantKey: core.ParamPriorityIDs.Value()},
		"WithAssigneeIDs":     {option: s.WithAssigneeIDs([]int{1}), wantKey: core.ParamAssigneeIDs.Value()},
		"WithCreatedUserIDs":  {option: s.WithCreatedUserIDs([]int{1}), wantKey: core.ParamCreatedUserIDs.Value()},
		"WithResolutionIDs":   {option: s.WithResolutionIDs([]int{1}), wantKey: core.ParamResolutionIDs.Value()},
		"WithParentChild":     {option: s.WithParentChild(0), wantKey: core.ParamParentChild.Value()},
		"WithAttachment":      {option: s.WithAttachment(true), wantKey: core.ParamAttachment.Value()},
		"WithSharedFile":      {option: s.WithSharedFile(true), wantKey: core.ParamSharedFile.Value()},
		"WithIssueSort":       {option: s.WithIssueSort(backlog.IssueSortCreated), wantKey: core.ParamSort.Value()},
		"WithOrder":           {option: s.WithOrder(backlog.OrderAsc), wantKey: core.ParamOrder.Value()},
		"WithOffset":          {option: s.WithOffset(0), wantKey: core.ParamOffset.Value()},
		"WithCount":           {option: s.WithCount(20), wantKey: core.ParamCount.Value()},
		"WithCreatedSince":    {option: s.WithCreatedSince(date), wantKey: core.ParamCreatedSince.Value()},
		"WithCreatedUntil":    {option: s.WithCreatedUntil(date), wantKey: core.ParamCreatedUntil.Value()},
		"WithUpdatedSince":    {option: s.WithUpdatedSince(date), wantKey: core.ParamUpdatedSince.Value()},
		"WithUpdatedUntil":    {option: s.WithUpdatedUntil(date), wantKey: core.ParamUpdatedUntil.Value()},
		"WithStartDateSince":  {option: s.WithStartDateSince(date), wantKey: core.ParamStartDateSince.Value()},
		"WithStartDateUntil":  {option: s.WithStartDateUntil(date), wantKey: core.ParamStartDateUntil.Value()},
		"WithDueDateSince":    {option: s.WithDueDateSince(date), wantKey: core.ParamDueDateSince.Value()},
		"WithDueDateUntil":    {option: s.WithDueDateUntil(date), wantKey: core.ParamDueDateUntil.Value()},
		"WithHasDueDate":      {option: s.WithHasDueDate(false), wantKey: core.ParamHasDueDate.Value()},
		"WithIDs":             {option: s.WithIDs([]int{1}), wantKey: core.ParamIDs.Value()},
		"WithParentIssueIDs":  {option: s.WithParentIssueIDs([]int{1}), wantKey: core.ParamParentIssueIDs.Value()},
		"WithKeyword":         {option: s.WithKeyword("bug"), wantKey: core.ParamKeyword.Value()},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}
