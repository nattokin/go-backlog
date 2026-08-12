package backlog_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/option"
)

func TestIssueOptionService(t *testing.T) {
	c, err := backlog.NewClient("https://example.backlog.com", "token")
	require.NoError(t, err)
	s := c.Issue.Option

	date := "2024-01-01"

	cases := map[string]struct {
		option  backlog.RequestOption
		wantKey string
	}{
		"WithActualHours":       {option: s.WithActualHours(1.0), wantKey: option.ParamActualHours.Value()},
		"WithAssigneeID":        {option: s.WithAssigneeID(1), wantKey: option.ParamAssigneeID.Value()},
		"WithAssigneeIDs":       {option: s.WithAssigneeIDs([]int{1}), wantKey: option.ParamAssigneeIDs.Value()},
		"WithAttachment":        {option: s.WithAttachment(true), wantKey: option.ParamAttachment.Value()},
		"WithAttachmentIDs":     {option: s.WithAttachmentIDs([]int{1}), wantKey: option.ParamAttachmentIDs.Value()},
		"WithCategoryIDs":       {option: s.WithCategoryIDs([]int{1}), wantKey: option.ParamCategoryIDs.Value()},
		"WithComment":           {option: s.WithComment("comment"), wantKey: option.ParamComment.Value()},
		"WithCount":             {option: s.WithCount(20), wantKey: option.ParamCount.Value()},
		"WithCreatedSince":      {option: s.WithCreatedSince(date), wantKey: option.ParamCreatedSince.Value()},
		"WithCreatedUntil":      {option: s.WithCreatedUntil(date), wantKey: option.ParamCreatedUntil.Value()},
		"WithCreatedUserIDs":    {option: s.WithCreatedUserIDs([]int{1}), wantKey: option.ParamCreatedUserIDs.Value()},
		"WithCustomFieldItems":  {option: s.WithCustomFieldItems(1, []int{10}), wantKey: "customField"},
		"WithCustomFieldNum":    {option: s.WithCustomFieldNum(1, 3.0), wantKey: "customField"},
		"WithCustomFieldOther":  {option: s.WithCustomFieldOther(1, "other"), wantKey: "customField"},
		"WithCustomFieldString": {option: s.WithCustomFieldString(1, "string"), wantKey: "customField"},
		"WithCustomFieldTime":   {option: s.WithCustomFieldTime(1, time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), wantKey: "customField"},
		"WithDescription":       {option: s.WithDescription("desc"), wantKey: option.ParamDescription.Value()},
		"WithDueDate":           {option: s.WithDueDate(date), wantKey: option.ParamDueDate.Value()},
		"WithDueDateSince":      {option: s.WithDueDateSince(date), wantKey: option.ParamDueDateSince.Value()},
		"WithDueDateUntil":      {option: s.WithDueDateUntil(date), wantKey: option.ParamDueDateUntil.Value()},
		"WithEstimatedHours":    {option: s.WithEstimatedHours(1.0), wantKey: option.ParamEstimatedHours.Value()},
		"WithHasDueDate":        {option: s.WithHasDueDate(false), wantKey: option.ParamHasDueDate.Value()},
		"WithIDs":               {option: s.WithIDs([]int{1}), wantKey: option.ParamIDs.Value()},
		"WithIssueSort":         {option: s.WithIssueSort(backlog.IssueSortCreated), wantKey: option.ParamSort.Value()},
		"WithIssueTypeID":       {option: s.WithIssueTypeID(1), wantKey: option.ParamIssueTypeID.Value()},
		"WithIssueTypeIDs":      {option: s.WithIssueTypeIDs([]int{1}), wantKey: option.ParamIssueTypeIDs.Value()},
		"WithKeyword":           {option: s.WithKeyword("bug"), wantKey: option.ParamKeyword.Value()},
		"WithMilestoneIDs":      {option: s.WithMilestoneIDs([]int{1}), wantKey: option.ParamMilestoneIDs.Value()},
		"WithNotifiedUserIDs":   {option: s.WithNotifiedUserIDs([]int{1}), wantKey: option.ParamNotifiedUserIDs.Value()},
		"WithOffset":            {option: s.WithOffset(0), wantKey: option.ParamOffset.Value()},
		"WithOrder":             {option: s.WithOrder(backlog.OrderAsc), wantKey: option.ParamOrder.Value()},
		"WithParentChild":       {option: s.WithParentChild(0), wantKey: option.ParamParentChild.Value()},
		"WithParentIssueID":     {option: s.WithParentIssueID(1), wantKey: option.ParamParentIssueID.Value()},
		"WithParentIssueIDs":    {option: s.WithParentIssueIDs([]int{1}), wantKey: option.ParamParentIssueIDs.Value()},
		"WithPriorityID":        {option: s.WithPriorityID(1), wantKey: option.ParamPriorityID.Value()},
		"WithPriorityIDs":       {option: s.WithPriorityIDs([]int{1}), wantKey: option.ParamPriorityIDs.Value()},
		"WithProjectIDs":        {option: s.WithProjectIDs([]int{1}), wantKey: option.ParamProjectIDs.Value()},
		"WithResolutionID":      {option: s.WithResolutionID(1), wantKey: option.ParamResolutionID.Value()},
		"WithResolutionIDs":     {option: s.WithResolutionIDs([]int{1}), wantKey: option.ParamResolutionIDs.Value()},
		"WithSharedFile":        {option: s.WithSharedFile(true), wantKey: option.ParamSharedFile.Value()},
		"WithStartDate":         {option: s.WithStartDate(date), wantKey: option.ParamStartDate.Value()},
		"WithStartDateSince":    {option: s.WithStartDateSince(date), wantKey: option.ParamStartDateSince.Value()},
		"WithStartDateUntil":    {option: s.WithStartDateUntil(date), wantKey: option.ParamStartDateUntil.Value()},
		"WithStatusID":          {option: s.WithStatusID(1), wantKey: option.ParamStatusID.Value()},
		"WithStatusIDs":         {option: s.WithStatusIDs([]int{backlog.IssueStatusOpen}), wantKey: option.ParamStatusIDs.Value()},
		"WithSummary":           {option: s.WithSummary("summary"), wantKey: option.ParamSummary.Value()},
		"WithUpdatedSince":      {option: s.WithUpdatedSince(date), wantKey: option.ParamUpdatedSince.Value()},
		"WithUpdatedUntil":      {option: s.WithUpdatedUntil(date), wantKey: option.ParamUpdatedUntil.Value()},
		"WithVersionIDs":        {option: s.WithVersionIDs([]int{1}), wantKey: option.ParamVersionIDs.Value()},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}
