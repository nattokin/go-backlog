package core_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
)

func TestOptionService_slice(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		option   core.RequestOption
		key      string
		wantVals []string
		wantErr  bool
	}{
		"WithActivityTypeIDs-valid": {
			option:   o.WithActivityTypeIDs([]int{1, 2, 26}),
			key:      core.ParamActivityTypeIDs.Value(),
			wantVals: []string{"1", "2", "26"},
		},
		"WithActivityTypeIDs-invalid-0": {
			option:  o.WithActivityTypeIDs([]int{0}),
			wantErr: true,
		},
		"WithActivityTypeIDs-invalid-27": {
			option:  o.WithActivityTypeIDs([]int{27}),
			wantErr: true,
		},
		"WithApplicableIssueTypeIDs-valid": {
			option:   o.WithApplicableIssueTypeIDs([]int{1, 2, 3}),
			key:      core.ParamApplicableIssueTypeIDs.Value(),
			wantVals: []string{"1", "2", "3"},
		},
		"WithApplicableIssueTypeIDs-invalid-0": {
			option:  o.WithApplicableIssueTypeIDs([]int{0}),
			wantErr: true,
		},
		"WithAttachmentIDs-valid": {
			option:   o.WithAttachmentIDs([]int{10, 20}),
			key:      core.ParamAttachmentIDs.Value(),
			wantVals: []string{"10", "20"},
		},
		"WithAttachmentIDs-invalid-0": {
			option:  o.WithAttachmentIDs([]int{0}),
			wantErr: true,
		},
		"WithItems-valid": {
			option:   o.WithItems([]string{"High", "Medium", "Low"}),
			key:      core.ParamItems.Value(),
			wantVals: []string{"High", "Medium", "Low"},
		},
		"WithItems-invalid-empty-element": {
			option:  o.WithItems([]string{"High", "", "Low"}),
			wantErr: true,
		},
		"WithItems-invalid-all-empty": {
			option:  o.WithItems([]string{""}),
			wantErr: true,
		},
		"WithProjectIDs-valid": {
			option:   o.WithProjectIDs([]int{1, 2}),
			key:      core.ParamProjectIDs.Value(),
			wantVals: []string{"1", "2"},
		},
		"WithProjectIDs-invalid-0": {
			option:  o.WithProjectIDs([]int{0}),
			wantErr: true,
		},
		"WithIssueTypeIDs-valid": {
			option:   o.WithIssueTypeIDs([]int{1}),
			key:      core.ParamIssueTypeIDs.Value(),
			wantVals: []string{"1"},
		},
		"WithIssueTypeIDs-invalid-0": {
			option:  o.WithIssueTypeIDs([]int{0}),
			wantErr: true,
		},
		"WithCategoryIDs-valid": {
			option:   o.WithCategoryIDs([]int{1, 2}),
			key:      core.ParamCategoryIDs.Value(),
			wantVals: []string{"1", "2"},
		},
		"WithCategoryIDs-invalid-0": {
			option:  o.WithCategoryIDs([]int{0}),
			wantErr: true,
		},
		"WithVersionIDs-valid": {
			option:   o.WithVersionIDs([]int{1}),
			key:      core.ParamVersionIDs.Value(),
			wantVals: []string{"1"},
		},
		"WithVersionIDs-invalid-0": {
			option:  o.WithVersionIDs([]int{0}),
			wantErr: true,
		},
		"WithMilestoneIDs-valid": {
			option:   o.WithMilestoneIDs([]int{1}),
			key:      core.ParamMilestoneIDs.Value(),
			wantVals: []string{"1"},
		},
		"WithMilestoneIDs-invalid-0": {
			option:  o.WithMilestoneIDs([]int{0}),
			wantErr: true,
		},
		"WithIssueIDs-valid": {
			option:   o.WithIssueIDs([]int{1, 2}),
			key:      core.ParamIssueIDs.Value(),
			wantVals: []string{"1", "2"},
		},
		"WithIssueIDs-invalid-0": {
			option:  o.WithIssueIDs([]int{0}),
			wantErr: true,
		},
		"WithNotifiedUserIDs-valid": {
			option:   o.WithNotifiedUserIDs([]int{1}),
			key:      core.ParamNotifiedUserIDs.Value(),
			wantVals: []string{"1"},
		},
		"WithNotifiedUserIDs-invalid-0": {
			option:  o.WithNotifiedUserIDs([]int{0}),
			wantErr: true,
		},
		"WithStatusIDs-valid": {
			option:   o.WithStatusIDs([]int{1, 2}),
			key:      core.ParamStatusIDs.Value(),
			wantVals: []string{"1", "2"},
		},
		"WithStatusIDs-invalid-0": {
			option:  o.WithStatusIDs([]int{0}),
			wantErr: true,
		},
		"WithPriorityIDs-valid": {
			option:   o.WithPriorityIDs([]int{2}),
			key:      core.ParamPriorityIDs.Value(),
			wantVals: []string{"2"},
		},
		"WithPriorityIDs-invalid-0": {
			option:  o.WithPriorityIDs([]int{0}),
			wantErr: true,
		},
		"WithAssigneeIDs-valid": {
			option:   o.WithAssigneeIDs([]int{1}),
			key:      core.ParamAssigneeIDs.Value(),
			wantVals: []string{"1"},
		},
		"WithAssigneeIDs-invalid-0": {
			option:  o.WithAssigneeIDs([]int{0}),
			wantErr: true,
		},
		"WithCreatedUserIDs-valid": {
			option:   o.WithCreatedUserIDs([]int{1}),
			key:      core.ParamCreatedUserIDs.Value(),
			wantVals: []string{"1"},
		},
		"WithCreatedUserIDs-invalid-0": {
			option:  o.WithCreatedUserIDs([]int{0}),
			wantErr: true,
		},
		"WithResolutionIDs-valid": {
			option:   o.WithResolutionIDs([]int{1}),
			key:      core.ParamResolutionIDs.Value(),
			wantVals: []string{"1"},
		},
		"WithResolutionIDs-invalid-0": {
			option:  o.WithResolutionIDs([]int{0}),
			wantErr: true,
		},
		"WithIDs-valid": {
			option:   o.WithIDs([]int{1, 2, 3}),
			key:      core.ParamIDs.Value(),
			wantVals: []string{"1", "2", "3"},
		},
		"WithIDs-invalid-0": {
			option:  o.WithIDs([]int{0}),
			wantErr: true,
		},
		"WithParentIssueIDs-valid": {
			option:   o.WithParentIssueIDs([]int{1}),
			key:      core.ParamParentIssueIDs.Value(),
			wantVals: []string{"1"},
		},
		"WithParentIssueIDs-invalid-0": {
			option:  o.WithParentIssueIDs([]int{0}),
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			form := url.Values{}
			err := tc.option.Check()
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			_ = tc.option.Set(form)
			assert.Equal(t, tc.wantVals, form[tc.key])
		})
	}
}
