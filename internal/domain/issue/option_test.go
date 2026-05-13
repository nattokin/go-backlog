package issue_test

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/issue"
)

func TestWithCustomField(t *testing.T) {
	cases := map[string]struct {
		option    core.RequestOption
		key       string
		wantValue string
		wantErr   bool
	}{
		"string": {
			option:    issue.WithCustomField(1, "test"),
			key:       "customField_1",
			wantValue: "test",
		},
		"string-invalid-empty": {
			option:  issue.WithCustomField(2, ""),
			wantErr: true,
		},
		"string-invalid-id-zero": {
			option:  issue.WithCustomField(0, "test"),
			wantErr: true,
		},
		"string-invalid-id-negative": {
			option:  issue.WithCustomField(-1, "test"),
			wantErr: true,
		},
		"float-1": {
			option:    issue.WithCustomField(1, 1.0),
			key:       "customField_1",
			wantValue: "1",
		},
		"float-1.5": {
			option:    issue.WithCustomField(2, 1.5),
			key:       "customField_2",
			wantValue: "1.5",
		},
		"float-negative-1": {
			option:    issue.WithCustomField(3, -1.0),
			key:       "customField_3",
			wantValue: "-1",
		},
		"float-negative-1.5": {
			option:    issue.WithCustomField(4, -1.5),
			key:       "customField_4",
			wantValue: "-1.5",
		},
		"float-invalid-id-zero": {
			option:  issue.WithCustomField(0, 1.0),
			wantErr: true,
		},
		"float-invalid-id-negative": {
			option:  issue.WithCustomField(-1, 1.0),
			wantErr: true,
		},
		"time": {
			option:    issue.WithCustomField(1, time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC)),
			key:       "customField_1",
			wantValue: "2024-01-10",
		},
		"time-invalid-zero": {
			option:  issue.WithCustomField(2, time.Time{}),
			wantErr: true,
		},
		"time-invalid-id-zero": {
			option:  issue.WithCustomField(0, time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC)),
			wantErr: true,
		},
		"time-invalid-id-negative": {
			option:  issue.WithCustomField(-1, time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC)),
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			form := url.Values{}
			err := tc.option.Check()
			if tc.wantErr {
				require.Error(t, err)
				errType := &core.ValidationError{}
				assert.ErrorAs(t, err, &errType)
				return
			}
			require.NoError(t, err)

			_ = tc.option.Set(form)
			_, exists := form[tc.key]
			assert.True(t, exists)
			assert.Equal(t, tc.wantValue, form.Get(tc.key))
		})
	}
}

func TestWithCustomFieldItems(t *testing.T) {
	cases := map[string]struct {
		id        int
		itemIDs   []int
		wantValue []string
		wantErr   bool
	}{
		"single": {
			id:        5,
			itemIDs:   []int{101},
			wantValue: []string{"101"},
		},
		"multiple": {
			id:        5,
			itemIDs:   []int{101, 202},
			wantValue: []string{"101", "202"},
		},
		"invalid-id-zero": {
			id:      0,
			itemIDs: []int{101},
			wantErr: true,
		},
		"invalid-id-negative": {
			id:      -1,
			itemIDs: []int{101},
			wantErr: true,
		},
		"invalid-itemID-zero": {
			id:      5,
			itemIDs: []int{0},
			wantErr: true,
		},
		"invalid-itemID-negative": {
			id:      5,
			itemIDs: []int{-1},
			wantErr: true,
		},
		"invalid-itemID-in-middle": {
			id:      5,
			itemIDs: []int{101, -1, 202},
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			opt := issue.WithCustomFieldItems(tc.id, tc.itemIDs)
			assert.Equal(t, "customField", opt.Key())

			err := opt.Check()
			if tc.wantErr {
				require.Error(t, err)
				errType := &core.ValidationError{}
				assert.ErrorAs(t, err, &errType)
				return
			}
			require.NoError(t, err)

			v := url.Values{}
			require.NoError(t, opt.Set(v))
			assert.Equal(t, tc.wantValue, v[fmt.Sprintf("customField_%d", tc.id)])
		})
	}
}

func TestWithCustomFieldOther(t *testing.T) {
	opt := issue.WithCustomFieldOther(5, "other text")
	assert.Equal(t, "customField", opt.Key())
	require.NoError(t, opt.Check())

	v := url.Values{}
	require.NoError(t, opt.Set(v))
	assert.Equal(t, "other text", v.Get("customField_5_otherValue"))
}
