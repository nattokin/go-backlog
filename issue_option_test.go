package backlog_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/core"
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
		"WithActualHours":       {option: s.WithActualHours(1.0), wantKey: core.ParamActualHours.Value()},
		"WithAssigneeID":        {option: s.WithAssigneeID(1), wantKey: core.ParamAssigneeID.Value()},
		"WithAssigneeIDs":       {option: s.WithAssigneeIDs([]int{1}), wantKey: core.ParamAssigneeIDs.Value()},
		"WithAttachment":        {option: s.WithAttachment(true), wantKey: core.ParamAttachment.Value()},
		"WithAttachmentIDs":     {option: s.WithAttachmentIDs([]int{1}), wantKey: core.ParamAttachmentIDs.Value()},
		"WithCategoryIDs":       {option: s.WithCategoryIDs([]int{1}), wantKey: core.ParamCategoryIDs.Value()},
		"WithComment":           {option: s.WithComment("comment"), wantKey: core.ParamComment.Value()},
		"WithCount":             {option: s.WithCount(20), wantKey: core.ParamCount.Value()},
		"WithCreatedSince":      {option: s.WithCreatedSince(date), wantKey: core.ParamCreatedSince.Value()},
		"WithCreatedUntil":      {option: s.WithCreatedUntil(date), wantKey: core.ParamCreatedUntil.Value()},
		"WithCreatedUserIDs":    {option: s.WithCreatedUserIDs([]int{1}), wantKey: core.ParamCreatedUserIDs.Value()},
		"WithCustomFieldItems":  {option: s.WithCustomFieldItems(1, []int{10}), wantKey: "customField"},
		"WithCustomFieldNum":    {option: s.WithCustomFieldNum(1, 3.0), wantKey: "customField"},
		"WithCustomFieldOther":  {option: s.WithCustomFieldOther(1, "other"), wantKey: "customField"},
		"WithCustomFieldString": {option: s.WithCustomFieldString(1, "string"), wantKey: "customField"},
		"WithCustomFieldTime":   {option: s.WithCustomFieldTime(1, time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)), wantKey: "customField"},
		"WithDescription":       {option: s.WithDescription("desc"), wantKey: core.ParamDescription.Value()},
		"WithDueDate":           {option: s.WithDueDate(date), wantKey: core.ParamDueDate.Value()},
		"WithDueDateSince":      {option: s.WithDueDateSince(date), wantKey: core.ParamDueDateSince.Value()},
		"WithDueDateUntil":      {option: s.WithDueDateUntil(date), wantKey: core.ParamDueDateUntil.Value()},
		"WithEstimatedHours":    {option: s.WithEstimatedHours(1.0), wantKey: core.ParamEstimatedHours.Value()},
		"WithHasDueDate":        {option: s.WithHasDueDate(false), wantKey: core.ParamHasDueDate.Value()},
		"WithIDs":               {option: s.WithIDs([]int{1}), wantKey: core.ParamIDs.Value()},
		"WithIssueSort":         {option: s.WithIssueSort(backlog.IssueSortCreated), wantKey: core.ParamSort.Value()},
		"WithIssueTypeID":       {option: s.WithIssueTypeID(1), wantKey: core.ParamIssueTypeID.Value()},
		"WithIssueTypeIDs":      {option: s.WithIssueTypeIDs([]int{1}), wantKey: core.ParamIssueTypeIDs.Value()},
		"WithKeyword":           {option: s.WithKeyword("bug"), wantKey: core.ParamKeyword.Value()},
		"WithMilestoneIDs":      {option: s.WithMilestoneIDs([]int{1}), wantKey: core.ParamMilestoneIDs.Value()},
		"WithNotifiedUserIDs":   {option: s.WithNotifiedUserIDs([]int{1}), wantKey: core.ParamNotifiedUserIDs.Value()},
		"WithOffset":            {option: s.WithOffset(0), wantKey: core.ParamOffset.Value()},
		"WithOrder":             {option: s.WithOrder(backlog.OrderAsc), wantKey: core.ParamOrder.Value()},
		"WithParentChild":       {option: s.WithParentChild(0), wantKey: core.ParamParentChild.Value()},
		"WithParentIssueID":     {option: s.WithParentIssueID(1), wantKey: core.ParamParentIssueID.Value()},
		"WithParentIssueIDs":    {option: s.WithParentIssueIDs([]int{1}), wantKey: core.ParamParentIssueIDs.Value()},
		"WithPriorityID":        {option: s.WithPriorityID(1), wantKey: core.ParamPriorityID.Value()},
		"WithPriorityIDs":       {option: s.WithPriorityIDs([]int{1}), wantKey: core.ParamPriorityIDs.Value()},
		"WithProjectIDs":        {option: s.WithProjectIDs([]int{1}), wantKey: core.ParamProjectIDs.Value()},
		"WithResolutionID":      {option: s.WithResolutionID(1), wantKey: core.ParamResolutionID.Value()},
		"WithResolutionIDs":     {option: s.WithResolutionIDs([]int{1}), wantKey: core.ParamResolutionIDs.Value()},
		"WithSharedFile":        {option: s.WithSharedFile(true), wantKey: core.ParamSharedFile.Value()},
		"WithStartDate":         {option: s.WithStartDate(date), wantKey: core.ParamStartDate.Value()},
		"WithStartDateSince":    {option: s.WithStartDateSince(date), wantKey: core.ParamStartDateSince.Value()},
		"WithStartDateUntil":    {option: s.WithStartDateUntil(date), wantKey: core.ParamStartDateUntil.Value()},
		"WithStatusID":          {option: s.WithStatusID(1), wantKey: core.ParamStatusID.Value()},
		"WithStatusIDs":         {option: s.WithStatusIDs([]int{1}), wantKey: core.ParamStatusIDs.Value()},
		"WithSummary":           {option: s.WithSummary("summary"), wantKey: core.ParamSummary.Value()},
		"WithUpdatedSince":      {option: s.WithUpdatedSince(date), wantKey: core.ParamUpdatedSince.Value()},
		"WithUpdatedUntil":      {option: s.WithUpdatedUntil(date), wantKey: core.ParamUpdatedUntil.Value()},
		"WithVersionIDs":        {option: s.WithVersionIDs([]int{1}), wantKey: core.ParamVersionIDs.Value()},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.wantKey, tc.option.Key())
		})
	}
}
