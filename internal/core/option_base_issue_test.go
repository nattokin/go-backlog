package core_test

import (
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

func TestOptionService_Issue(t *testing.T) {
	o := &core.OptionService{}

	// --- ID slice options ---------------------------------------------------------
	t.Run("id-slice-options", func(t *testing.T) {
		cases := map[string]struct {
			option  core.RequestOption
			key     string
			wantIDs []string
			wantErr bool
		}{
			"WithProjectIDs-valid": {
				option:  o.WithProjectIDs([]int{1, 2, 3}),
				key:     core.ParamProjectIDs.Value(),
				wantIDs: []string{"1", "2", "3"},
			},
			"WithProjectIDs-invalid-zero": {
				option:  o.WithProjectIDs([]int{0}),
				wantErr: true,
			},
			"WithIssueTypeIDs-valid": {
				option:  o.WithIssueTypeIDs([]int{10, 20}),
				key:     core.ParamIssueTypeIDs.Value(),
				wantIDs: []string{"10", "20"},
			},
			"WithIssueTypeIDs-invalid-negative": {
				option:  o.WithIssueTypeIDs([]int{-1}),
				wantErr: true,
			},
			"WithCategoryIDs-valid": {
				option:  o.WithCategoryIDs([]int{5}),
				key:     core.ParamCategoryIDs.Value(),
				wantIDs: []string{"5"},
			},
			"WithCategoryIDs-invalid-zero": {
				option:  o.WithCategoryIDs([]int{0}),
				wantErr: true,
			},
			"WithVersionIDs-valid": {
				option:  o.WithVersionIDs([]int{3, 4}),
				key:     core.ParamVersionIDs.Value(),
				wantIDs: []string{"3", "4"},
			},
			"WithVersionIDs-invalid": {
				option:  o.WithVersionIDs([]int{0}),
				wantErr: true,
			},
			"WithMilestoneIDs-valid": {
				option:  o.WithMilestoneIDs([]int{7}),
				key:     core.ParamMilestoneIDs.Value(),
				wantIDs: []string{"7"},
			},
			"WithMilestoneIDs-invalid": {
				option:  o.WithMilestoneIDs([]int{0}),
				wantErr: true,
			},
			"WithStatusIDs-valid": {
				option:  o.WithStatusIDs([]int{1, 2}),
				key:     core.ParamStatusIDs.Value(),
				wantIDs: []string{"1", "2"},
			},
			"WithStatusIDs-invalid": {
				option:  o.WithStatusIDs([]int{0}),
				wantErr: true,
			},
			"WithPriorityIDs-valid": {
				option:  o.WithPriorityIDs([]int{2, 3}),
				key:     core.ParamPriorityIDs.Value(),
				wantIDs: []string{"2", "3"},
			},
			"WithPriorityIDs-invalid": {
				option:  o.WithPriorityIDs([]int{0}),
				wantErr: true,
			},
			"WithAssigneeIDs-valid": {
				option:  o.WithAssigneeIDs([]int{100}),
				key:     core.ParamAssigneeIDs.Value(),
				wantIDs: []string{"100"},
			},
			"WithAssigneeIDs-invalid": {
				option:  o.WithAssigneeIDs([]int{0}),
				wantErr: true,
			},
			"WithCreatedUserIDs-valid": {
				option:  o.WithCreatedUserIDs([]int{1, 2}),
				key:     core.ParamCreatedUserIDs.Value(),
				wantIDs: []string{"1", "2"},
			},
			"WithCreatedUserIDs-invalid": {
				option:  o.WithCreatedUserIDs([]int{0}),
				wantErr: true,
			},
			"WithResolutionIDs-valid": {
				option:  o.WithResolutionIDs([]int{1}),
				key:     core.ParamResolutionIDs.Value(),
				wantIDs: []string{"1"},
			},
			"WithResolutionIDs-invalid": {
				option:  o.WithResolutionIDs([]int{0}),
				wantErr: true,
			},
			"WithIDs-valid": {
				option:  o.WithIDs([]int{10, 20, 30}),
				key:     core.ParamIDs.Value(),
				wantIDs: []string{"10", "20", "30"},
			},
			"WithIDs-invalid": {
				option:  o.WithIDs([]int{0}),
				wantErr: true,
			},
			"WithParentIssueIDs-valid": {
				option:  o.WithParentIssueIDs([]int{5, 6}),
				key:     core.ParamParentIssueIDs.Value(),
				wantIDs: []string{"5", "6"},
			},
			"WithParentIssueIDs-invalid": {
				option:  o.WithParentIssueIDs([]int{0}),
				wantErr: true,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				q := url.Values{}
				err := tc.option.Check()
				if tc.wantErr {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				_ = tc.option.Set(q)
				assert.Equal(t, tc.wantIDs, q[tc.key])
			})
		}
	})

	// --- Integer options ----------------------------------------------------------
	t.Run("integer-options", func(t *testing.T) {
		cases := map[string]struct {
			option    core.RequestOption
			key       string
			wantValue int
			wantErr   bool
		}{
			"WithParentChild-valid-0": {
				option:    o.WithParentChild(0),
				key:       core.ParamParentChild.Value(),
				wantValue: 0,
			},
			"WithParentChild-valid-4": {
				option:    o.WithParentChild(4),
				key:       core.ParamParentChild.Value(),
				wantValue: 4,
			},
			"WithParentChild-invalid-negative": {
				option:  o.WithParentChild(-1),
				wantErr: true,
			},
			"WithParentChild-invalid-5": {
				option:  o.WithParentChild(5),
				wantErr: true,
			},
			"WithOffset-valid-0": {
				option:    o.WithOffset(0),
				key:       core.ParamOffset.Value(),
				wantValue: 0,
			},
			"WithOffset-valid-100": {
				option:    o.WithOffset(100),
				key:       core.ParamOffset.Value(),
				wantValue: 100,
			},
			"WithOffset-invalid-negative": {
				option:  o.WithOffset(-1),
				wantErr: true,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				q := url.Values{}
				err := tc.option.Check()
				if tc.wantErr {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				_ = tc.option.Set(q)
				assert.Equal(t, strconv.Itoa(tc.wantValue), q.Get(tc.key))
			})
		}
	})

	// --- Boolean options ----------------------------------------------------------
	t.Run("boolean-options", func(t *testing.T) {
		cases := map[string]struct {
			option    core.RequestOption
			key       string
			wantValue bool
		}{
			"WithAttachment-true": {
				option:    o.WithAttachment(true),
				key:       core.ParamAttachment.Value(),
				wantValue: true,
			},
			"WithAttachment-false": {
				option:    o.WithAttachment(false),
				key:       core.ParamAttachment.Value(),
				wantValue: false,
			},
			"WithSharedFile-true": {
				option:    o.WithSharedFile(true),
				key:       core.ParamSharedFile.Value(),
				wantValue: true,
			},
			"WithSharedFile-false": {
				option:    o.WithSharedFile(false),
				key:       core.ParamSharedFile.Value(),
				wantValue: false,
			},
			"WithHasDueDate-true": {
				option:    o.WithHasDueDate(true),
				key:       core.ParamHasDueDate.Value(),
				wantValue: true,
			},
			"WithHasDueDate-false": {
				option:    o.WithHasDueDate(false),
				key:       core.ParamHasDueDate.Value(),
				wantValue: false,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				q := url.Values{}
				err := tc.option.Check()
				require.NoError(t, err)
				_ = tc.option.Set(q)
				assert.Equal(t, strconv.FormatBool(tc.wantValue), q.Get(tc.key))
			})
		}
	})

	// --- IssueSort option ---------------------------------------------------------
	t.Run("WithIssueSort", func(t *testing.T) {
		cases := map[string]struct {
			sort    model.IssueSort
			wantErr bool
		}{
			"issueType":      {sort: model.IssueSortIssueType},
			"category":       {sort: model.IssueSortCategory},
			"version":        {sort: model.IssueSortVersion},
			"milestone":      {sort: model.IssueSortMilestone},
			"summary":        {sort: model.IssueSortSummary},
			"status":         {sort: model.IssueSortStatus},
			"priority":       {sort: model.IssueSortPriority},
			"attachment":     {sort: model.IssueSortAttachment},
			"sharedFile":     {sort: model.IssueSortSharedFile},
			"created":        {sort: model.IssueSortCreated},
			"createdUser":    {sort: model.IssueSortCreatedUser},
			"updated":        {sort: model.IssueSortUpdated},
			"updatedUser":    {sort: model.IssueSortUpdatedUser},
			"assignee":       {sort: model.IssueSortAssignee},
			"startDate":      {sort: model.IssueSortStartDate},
			"dueDate":        {sort: model.IssueSortDueDate},
			"estimatedHours": {sort: model.IssueSortEstimatedHours},
			"actualHours":    {sort: model.IssueSortActualHours},
			"childIssue":     {sort: model.IssueSortChildIssue},
			"invalid":        {sort: "invalid", wantErr: true},
			"empty":          {sort: "", wantErr: true},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				opt := o.WithIssueSort(tc.sort)
				q := url.Values{}
				err := opt.Check()
				if tc.wantErr {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				_ = opt.Set(q)
				assert.Equal(t, string(tc.sort), q.Get(core.ParamSort.Value()))
			})
		}
	})

	// --- Date options -------------------------------------------------------------
	t.Run("date-options", func(t *testing.T) {
		date := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
		want := "2024-03-15"

		cases := map[string]struct {
			option core.RequestOption
			key    string
		}{
			"WithCreatedSince":    {option: o.WithCreatedSince(date), key: core.ParamCreatedSince.Value()},
			"WithCreatedUntil":    {option: o.WithCreatedUntil(date), key: core.ParamCreatedUntil.Value()},
			"WithUpdatedSince":    {option: o.WithUpdatedSince(date), key: core.ParamUpdatedSince.Value()},
			"WithUpdatedUntil":    {option: o.WithUpdatedUntil(date), key: core.ParamUpdatedUntil.Value()},
			"WithStartDateSince":  {option: o.WithStartDateSince(date), key: core.ParamStartDateSince.Value()},
			"WithStartDateUntil":  {option: o.WithStartDateUntil(date), key: core.ParamStartDateUntil.Value()},
			"WithDueDateSince":    {option: o.WithDueDateSince(date), key: core.ParamDueDateSince.Value()},
			"WithDueDateUntil":    {option: o.WithDueDateUntil(date), key: core.ParamDueDateUntil.Value()},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				q := url.Values{}
				err := tc.option.Check()
				require.NoError(t, err)
				_ = tc.option.Set(q)
				assert.Equal(t, want, q.Get(tc.key))
			})
		}
	})
}
