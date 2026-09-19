package option_test

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/option"
)

func TestOptionService_int(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		option    *option.APIParamOption
		key       string
		wantValue int
		wantErr   bool
	}{
		"WithActualHours-valid-1": {
			option:    o.WithActualHours(1),
			key:       option.ParamActualHours.Value(),
			wantValue: 1,
		},
		"WithActualHours-invalid-0": {
			option:  o.WithActualHours(0),
			wantErr: true,
		},
		"WithAssigneeID-valid-1": {
			option:    o.WithAssigneeID(1),
			key:       option.ParamAssigneeID.Value(),
			wantValue: 1,
		},
		"WithAssigneeID-invalid-0": {
			option:  o.WithAssigneeID(0),
			wantErr: true,
		},
		"WithCommentID-valid-1": {
			option:    o.WithCommentID(1),
			key:       option.ParamCommentID.Value(),
			wantValue: 1,
		},
		"WithCommentID-invalid-0": {
			option:  o.WithCommentID(0),
			wantErr: true,
		},
		"WithCount-valid-1": {
			option:    o.WithCount(1),
			key:       option.ParamCount.Value(),
			wantValue: 1,
		},
		"WithCount-valid-100": {
			option:    o.WithCount(100),
			key:       option.ParamCount.Value(),
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
		"WithSharedFileCount-valid-1": {
			option:    o.WithSharedFileCount(1),
			key:       option.ParamCount.Value(),
			wantValue: 1,
		},
		"WithSharedFileCount-valid-1000": {
			option:    o.WithSharedFileCount(1000),
			key:       option.ParamCount.Value(),
			wantValue: 1000,
		},
		"WithSharedFileCount-invalid-0": {
			option:  o.WithSharedFileCount(0),
			wantErr: true,
		},
		"WithSharedFileCount-invalid-1001": {
			option:  o.WithSharedFileCount(1001),
			wantErr: true,
		},
		"WithEstimatedHours-valid-1": {
			option:    o.WithEstimatedHours(1),
			key:       option.ParamEstimatedHours.Value(),
			wantValue: 1,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tc.wantErr {
				require.Error(t, tc.option.Check())
				return
			}

			require.NoError(t, tc.option.Check())
			v := url.Values{}
			require.NoError(t, tc.option.Set(v))
			assert.Equal(t, strconv.Itoa(tc.wantValue), v.Get(tc.key))
		})
	}
}
