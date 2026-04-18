package core_test

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
)

func TestOptionService_slice(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		option    core.RequestOption
		key       string
		wantValue []string
		wantErr   bool
	}{
		"WithActivityTypeIDs-all-range": {
			option: o.WithActivityTypeIDs(func() []int {
				var all []int
				for i := 1; i <= 26; i++ {
					all = append(all, i)
				}
				return all
			}()),
			key: core.ParamActivityTypeIDs.Value(),
			wantValue: func() []string {
				v := []string{}
				for i := 1; i <= 26; i++ {
					v = append(v, strconv.Itoa(i))
				}
				return v
			}(),
		},
		"WithActivityTypeIDs-invalid-above": {
			option:  o.WithActivityTypeIDs([]int{27}),
			key:     core.ParamActivityTypeIDs.Value(),
			wantErr: true,
		},
		"WithActivityTypeIDs-invalid-below": {
			option:  o.WithActivityTypeIDs([]int{0}),
			key:     core.ParamActivityTypeIDs.Value(),
			wantErr: true,
		},
		"WithActivityTypeIDs-invalid-mixed-high": {
			option:  o.WithActivityTypeIDs([]int{26, 27}),
			key:     core.ParamActivityTypeIDs.Value(),
			wantErr: true,
		},
		"WithActivityTypeIDs-invalid-mixed-low": {
			option:  o.WithActivityTypeIDs([]int{0, 1}),
			key:     core.ParamActivityTypeIDs.Value(),
			wantErr: true,
		},
		"WithActivityTypeIDs-single-max": {
			option:    o.WithActivityTypeIDs([]int{26}),
			key:       core.ParamActivityTypeIDs.Value(),
			wantValue: []string{"26"},
		},
		"WithActivityTypeIDs-single-min": {
			option:    o.WithActivityTypeIDs([]int{1}),
			key:       core.ParamActivityTypeIDs.Value(),
			wantValue: []string{"1"},
		},

		"WithAssigneeIDs-invalid": {
			option:  o.WithAssigneeIDs([]int{0}),
			wantErr: true,
		},
		"WithAssigneeIDs-valid": {
			option:    o.WithAssigneeIDs([]int{100}),
			key:       core.ParamAssigneeIDs.Value(),
			wantValue: []string{"100"},
		},

		"WithCategoryIDs-invalid-zero": {
			option:  o.WithCategoryIDs([]int{0}),
			wantErr: true,
		},
		"WithCategoryIDs-valid": {
			option:    o.WithCategoryIDs([]int{5}),
			key:       core.ParamCategoryIDs.Value(),
			wantValue: []string{"5"},
		},

		"WithCreatedUserIDs-invalid": {
			option:  o.WithCreatedUserIDs([]int{0}),
			wantErr: true,
		},
		"WithCreatedUserIDs-valid": {
			option:    o.WithCreatedUserIDs([]int{1, 2}),
			key:       core.ParamCreatedUserIDs.Value(),
			wantValue: []string{"1", "2"},
		},

		"WithIDs-invalid": {
			option:  o.WithIDs([]int{0}),
			wantErr: true,
		},
		"WithIDs-valid": {
			option:    o.WithIDs([]int{10, 20, 30}),
			key:       core.ParamIDs.Value(),
			wantValue: []string{"10", "20", "30"},
		},

		"WithIssueTypeIDs-invalid-negative": {
			option:  o.WithIssueTypeIDs([]int{-1}),
			wantErr: true,
		},
		"WithIssueTypeIDs-valid": {
			option:    o.WithIssueTypeIDs([]int{10, 20}),
			key:       core.ParamIssueTypeIDs.Value(),
			wantValue: []string{"10", "20"},
		},

		"WithMilestoneIDs-invalid": {
			option:  o.WithMilestoneIDs([]int{0}),
			wantErr: true,
		},
		"WithMilestoneIDs-valid": {
			option:    o.WithMilestoneIDs([]int{7}),
			key:       core.ParamMilestoneIDs.Value(),
			wantValue: []string{"7"},
		},

		"WithParentIssueIDs-invalid": {
			option:  o.WithParentIssueIDs([]int{0}),
			wantErr: true,
		},
		"WithParentIssueIDs-valid": {
			option:    o.WithParentIssueIDs([]int{5, 6}),
			key:       core.ParamParentIssueIDs.Value(),
			wantValue: []string{"5", "6"},
		},

		"WithPriorityIDs-invalid": {
			option:  o.WithPriorityIDs([]int{0}),
			wantErr: true,
		},
		"WithPriorityIDs-valid": {
			option:    o.WithPriorityIDs([]int{2, 3}),
			key:       core.ParamPriorityIDs.Value(),
			wantValue: []string{"2", "3"},
		},

		"WithProjectIDs-invalid-zero": {
			option:  o.WithProjectIDs([]int{0}),
			wantErr: true,
		},
		"WithProjectIDs-valid": {
			option:    o.WithProjectIDs([]int{1, 2, 3}),
			key:       core.ParamProjectIDs.Value(),
			wantValue: []string{"1", "2", "3"},
		},

		"WithResolutionIDs-invalid": {
			option:  o.WithResolutionIDs([]int{0}),
			wantErr: true,
		},
		"WithResolutionIDs-valid": {
			option:    o.WithResolutionIDs([]int{1}),
			key:       core.ParamResolutionIDs.Value(),
			wantValue: []string{"1"},
		},

		"WithStatusIDs-invalid": {
			option:  o.WithStatusIDs([]int{0}),
			wantErr: true,
		},
		"WithStatusIDs-valid": {
			option:    o.WithStatusIDs([]int{1, 2}),
			key:       core.ParamStatusIDs.Value(),
			wantValue: []string{"1", "2"},
		},

		"WithVersionIDs-invalid": {
			option:  o.WithVersionIDs([]int{0}),
			wantErr: true,
		},
		"WithVersionIDs-valid": {
			option:    o.WithVersionIDs([]int{3, 4}),
			key:       core.ParamVersionIDs.Value(),
			wantValue: []string{"3", "4"},
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
			assert.Equal(t, tc.wantValue, q[tc.key])
		})
	}
}
