package core_test

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
)

func TestOptionService_int(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		option    core.RequestOption
		key       string
		wantValue int
		wantErr   bool
	}{
		"WithActualHours-valid-1": {
			option:    o.WithActualHours(1),
			key:       core.ParamActualHours.Value(),
			wantValue: 1,
		},
		"WithActualHours-invalid-0": {
			option:  o.WithActualHours(0),
			wantErr: true,
		},
		"WithAssigneeID-valid-1": {
			option:    o.WithAssigneeID(1),
			key:       core.ParamAssigneeID.Value(),
			wantValue: 1,
		},
		"WithAssigneeID-invalid-0": {
			option:  o.WithAssigneeID(0),
			wantErr: true,
		},
		"WithCount-valid-1": {
			option:    o.WithCount(1),
			key:       core.ParamCount.Value(),
			wantValue: 1,
		},
		"WithCount-valid-100": {
			option:    o.WithCount(100),
			key:       core.ParamCount.Value(),
			wantValue: 100,
		},
		"WithCount-invalid-0": {
			option:  o.WithCount(0),
			wantErr: true,
		},
		"WithCount-invalid-101": {
			option:  o.WithCount(101),
			wantErr: true,
		},
		"WithEstimatedHours-valid-1": {
			option:    o.WithEstimatedHours(1),
			key:       core.ParamEstimatedHours.Value(),
			wantValue: 1,
		},
		"WithEstimatedHours-invalid-0": {
			option:  o.WithEstimatedHours(0),
			wantErr: true,
		},
		"WithIssueID-valid-1": {
			option:    o.WithIssueID(1),
			key:       core.ParamIssueID.Value(),
			wantValue: 1,
		},
		"WithIssueID-invalid-0": {
			option:  o.WithIssueID(0),
			wantErr: true,
		},
		"WithIssueTypeID-valid-1": {
			option:    o.WithIssueTypeID(1),
			key:       core.ParamIssueTypeID.Value(),
			wantValue: 1,
		},
		"WithIssueTypeID-invalid-0": {
			option:  o.WithIssueTypeID(0),
			wantErr: true,
		},
		"WithMaxID-valid-26": {
			option:    o.WithMaxID(26),
			key:       core.ParamMaxID.Value(),
			wantValue: 26,
		},
		"WithMaxID-invalid-27": {
			option:  o.WithMaxID(27),
			wantErr: true,
		},
		"WithMinID-valid-1": {
			option:    o.WithMinID(1),
			key:       core.ParamMinID.Value(),
			wantValue: 1,
		},
		"WithMinID-invalid-0": {
			option:  o.WithMinID(0),
			wantErr: true,
		},
		"WithMinID-invalid-27": {
			option:  o.WithMinID(27),
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
		"WithParentChild-invalid-5": {
			option:  o.WithParentChild(5),
			wantErr: true,
		},
		"WithParentChild-invalid-negative": {
			option:  o.WithParentChild(-1),
			wantErr: true,
		},
		"WithParentIssueID-valid-1": {
			option:    o.WithParentIssueID(1),
			key:       core.ParamParentIssueID.Value(),
			wantValue: 1,
		},
		"WithParentIssueID-invalid-0": {
			option:  o.WithParentIssueID(0),
			wantErr: true,
		},
		"WithPriorityID-valid-1": {
			option:    o.WithPriorityID(1),
			key:       core.ParamPriorityID.Value(),
			wantValue: 1,
		},
		"WithPriorityID-invalid-0": {
			option:  o.WithPriorityID(0),
			wantErr: true,
		},
		"WithResolutionID-valid-1": {
			option:    o.WithResolutionID(1),
			key:       core.ParamResolutionID.Value(),
			wantValue: 1,
		},
		"WithResolutionID-invalid-0": {
			option:  o.WithResolutionID(0),
			wantErr: true,
		},
		"WithRoleType-valid-1": {
			option:    o.WithRoleType(1),
			key:       core.ParamRoleType.Value(),
			wantValue: 1,
		},
		"WithRoleType-valid-6": {
			option:    o.WithRoleType(6),
			key:       core.ParamRoleType.Value(),
			wantValue: 6,
		},
		"WithRoleType-invalid-0": {
			option:  o.WithRoleType(0),
			wantErr: true,
		},
		"WithRoleType-invalid-7": {
			option:  o.WithRoleType(7),
			wantErr: true,
		},
		"WithStatusID-valid-1": {
			option:    o.WithStatusID(1),
			key:       core.ParamStatusID.Value(),
			wantValue: 1,
		},
		"WithStatusID-invalid-0": {
			option:  o.WithStatusID(0),
			wantErr: true,
		},
		"WithUserID-valid-1": {
			option:    o.WithUserID(1),
			key:       core.ParamUserID.Value(),
			wantValue: 1,
		},
		"WithUserID-valid-2": {
			option:    o.WithUserID(2),
			key:       core.ParamUserID.Value(),
			wantValue: 2,
		},
		"WithUserID-invalid-0": {
			option:  o.WithUserID(0),
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
			assert.Equal(t, strconv.Itoa(tc.wantValue), form.Get(tc.key))
		})
	}

}
