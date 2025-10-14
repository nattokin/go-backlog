package backlog

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
// ──────────────────────────────────────────────────────────────
//  TestQueryOption
// ──────────────────────────────────────────────────────────────
//

func TestQueryOption(t *testing.T) {
	cases := map[string]struct {
		option           *QueryOption
		expectCheckErr   bool
		expectSetErr     bool
		wantValue        string
		wantCheckErrType error
		wantSetErrType   error
	}{
		"Success": {
			option: &QueryOption{
				t:         queryKey,
				checkFunc: func() error { return nil },
				setFunc: func(query *QueryParams) error {
					query.Set(queryKey.Value(), "success")
					return nil
				},
			},
			expectCheckErr: false,
			expectSetErr:   false,
			wantValue:      "success",
		},
		"Check-error": {
			option:           newQueryOptionWithCheckError(queryKey),
			expectCheckErr:   true,
			expectSetErr:     false,
			wantCheckErrType: errors.New("check error"),
		},
		"set-error": {
			option:         newQueryOptionWithSetError(queryKey),
			expectCheckErr: false,
			expectSetErr:   true,
			wantSetErrType: errors.New("set error"),
		},
		"queryType-invalid": {
			option: &QueryOption{
				t: 0,
				setFunc: func(query *QueryParams) error {
					return nil
				},
			},
			expectCheckErr:   false,
			expectSetErr:     false,
			wantValue:        "",
			wantCheckErrType: nil,
		},
		"checkFunc-nil": {
			option: &QueryOption{
				t:         queryKey,
				checkFunc: nil,
				setFunc: func(query *QueryParams) error {
					query.Set(queryKey.Value(), "checkFunc nil")
					return nil
				},
			},
			expectCheckErr: false,
			expectSetErr:   false,
			wantValue:      "checkFunc nil",
		},
		"set-nil": {
			option: &QueryOption{
				t:         queryKey,
				checkFunc: func() error { return nil },
				setFunc:   nil,
			},
			expectCheckErr: false,
			expectSetErr:   true,
			wantSetErrType: newValidationError("set nil"),
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			query := NewQueryParams()
			err := tc.option.Check()
			if tc.expectCheckErr {
				require.Error(t, err)
				assert.IsType(t, tc.wantCheckErrType, err)
				return
			}
			require.NoError(t, err)

			if err := tc.option.set(query); tc.expectSetErr {
				require.Error(t, err)
				assert.IsType(t, tc.wantSetErrType, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, query.Get(tc.option.t.Value()))
			}
		})
	}

}

//
// ──────────────────────────────────────────────────────────────
//  TestQueryOptionService
// ──────────────────────────────────────────────────────────────
//

// TestQueryOptionService verifies that each method of QueryOptionService
// correctly applies the expected query parameters.
func TestQueryOptionService(t *testing.T) {

	// --- Boolean options ------------------------------------------------------------
	t.Run("boolean-options", func(t *testing.T) {
		o := newQueryOptionService()

		cases := map[string]struct {
			option    *QueryOption
			key       string
			wantValue string
		}{
			"WithAll-true": {
				option:    o.WithAll(true),
				key:       queryAll.Value(),
				wantValue: "true",
			},
			"WithAll-false": {
				option:    o.WithAll(false),
				key:       queryAll.Value(),
				wantValue: "false",
			},
			"WithArchived-true": {
				option:    o.WithArchived(true),
				key:       queryArchived.Value(),
				wantValue: "true",
			},
			"WithArchived-false": {
				option:    o.WithArchived(false),
				key:       queryArchived.Value(),
				wantValue: "false",
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				q := NewQueryParams()
				err := tc.option.set(q)
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, q.Get(tc.key))
			})
		}
	})

	// --- Integer options ------------------------------------------------------------
	t.Run("integer-options", func(t *testing.T) {
		o := newQueryOptionService()

		cases := map[string]struct {
			option    *QueryOption
			key       string
			wantValue string
			wantErr   bool
		}{
			// WithCount
			"WithCount-valid-1": {
				option:    o.WithCount(1),
				key:       queryCount.Value(),
				wantValue: "1",
			},
			"WithCount-valid-100": {
				option:    o.WithCount(100),
				key:       queryCount.Value(),
				wantValue: "100",
			},
			"WithCount-invalid-0": {
				option:  o.WithCount(0),
				key:     queryCount.Value(),
				wantErr: true,
			},
			"WithCount-invalid-101": {
				option:  o.WithCount(101),
				key:     queryCount.Value(),
				wantErr: true,
			},

			// WithMinID / WithMaxID
			"WithMinID-valid-1": {
				option:    o.WithMinID(1),
				key:       queryMinID.Value(),
				wantValue: "1",
			},
			"WithMinID-invalid-0": {
				option:  o.WithMinID(0),
				key:     queryMinID.Value(),
				wantErr: true,
			},
			"WithMaxID-valid-26": {
				option:    o.WithMaxID(26),
				key:       queryMaxID.Value(),
				wantValue: "26",
			},
			"WithMaxID-invalid-27": {
				option:  o.WithMaxID(27),
				key:     queryMaxID.Value(),
				wantErr: true,
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				q := NewQueryParams()
				err := tc.option.Check()
				if tc.wantErr {
					assert.Error(t, err)
					assert.Empty(t, q.Get(tc.key))
					return
				}
				require.NoError(t, err)
				_ = tc.option.set(q)
				assert.Equal(t, tc.wantValue, q.Get(tc.key))
			})
		}
	})

	// --- String options ------------------------------------------------------------
	t.Run("string-options", func(t *testing.T) {
		o := newQueryOptionService()

		cases := map[string]struct {
			option    *QueryOption
			key       string
			wantValue string
		}{
			"WithKeyword-non-empty": {
				option:    o.WithKeyword("backlog"),
				key:       queryKeyword.Value(),
				wantValue: "backlog",
			},
			"WithKeyword-empty": {
				option:    o.WithKeyword(""),
				key:       queryKeyword.Value(),
				wantValue: "",
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				q := NewQueryParams()

				err := tc.option.Check()
				require.NoError(t, err)
				_ = tc.option.set(q)
				assert.Equal(t, tc.wantValue, q.Get(tc.key))
			})
		}
	})

	// --- Enum or special options ----------------------------------------------------
	t.Run("enum-or-special-options", func(t *testing.T) {
		o := newQueryOptionService()

		cases := map[string]struct {
			option    *QueryOption
			key       string
			wantValue string
			wantErr   bool
		}{
			"WithOrder-asc": {
				option:    o.WithOrder(OrderAsc),
				key:       queryOrder.Value(),
				wantValue: "asc",
			},
			"WithOrder-desc": {
				option:    o.WithOrder(OrderDesc),
				key:       queryOrder.Value(),
				wantValue: "desc",
			},
			"WithOrder-invalid": {
				option:  o.WithOrder(Order("invalid")),
				key:     queryOrder.Value(),
				wantErr: true,
			},
			"WithActivityTypeIDs-single-min": {
				option:    o.WithActivityTypeIDs([]int{1}),
				key:       queryActivityTypeIDs.Value(),
				wantValue: "1",
			},
			"WithActivityTypeIDs-single-max": {
				option:    o.WithActivityTypeIDs([]int{26}),
				key:       queryActivityTypeIDs.Value(),
				wantValue: "26",
			},
			"WithActivityTypeIDs-all-range": {
				option: o.WithActivityTypeIDs(func() []int {
					var all []int
					for i := 1; i <= 26; i++ {
						all = append(all, i)
					}
					return all
				}()),
				key: queryActivityTypeIDs.Value(),
				wantValue: func() string {
					s := ""
					for i := 1; i <= 26; i++ {
						if i > 1 {
							s += ","
						}
						s += strconv.Itoa(i)
					}
					return s
				}(),
			},
			"WithActivityTypeIDs-invalid-below": {
				option:  o.WithActivityTypeIDs([]int{0}),
				key:     queryActivityTypeIDs.Value(),
				wantErr: true,
			},
			"WithActivityTypeIDs-invalid-above": {
				option:  o.WithActivityTypeIDs([]int{27}),
				key:     queryActivityTypeIDs.Value(),
				wantErr: true,
			},
			"WithActivityTypeIDs-invalid-mixed-low": {
				option:  o.WithActivityTypeIDs([]int{0, 1}),
				key:     queryActivityTypeIDs.Value(),
				wantErr: true,
			},
			"WithActivityTypeIDs-invalid-mixed-high": {
				option:  o.WithActivityTypeIDs([]int{26, 27}),
				key:     queryActivityTypeIDs.Value(),
				wantErr: true,
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				q := NewQueryParams()
				err := tc.option.Check()
				if tc.wantErr {
					assert.Error(t, err)
					assert.Empty(t, q.Get(tc.key))
					return
				}
				require.NoError(t, err)
				_ = tc.option.set(q)

				if tc.key == queryActivityTypeIDs.Value() {
					// Compare joined values
					values := (*q.Values)[tc.key]
					got := ""
					for i, v := range values {
						if i > 0 {
							got += ","
						}
						got += v
					}
					assert.Equal(t, tc.wantValue, got)
				} else {
					assert.Equal(t, tc.wantValue, q.Get(tc.key))
				}
			})
		}
	})
}
