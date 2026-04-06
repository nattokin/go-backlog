package backlog

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
// ──────────────────────────────────────────────────────────────
//  TestOptionService
// ──────────────────────────────────────────────────────────────
//

func TestOptionService(t *testing.T) {
	o := newOptionService()

	// --- Boolean options -----------------------------------------------------------
	t.Run("boolean-options", func(t *testing.T) {
		cases := map[string]struct {
			option    RequestOption
			key       string
			wantValue bool
		}{
			"WithAll-true": {
				option:    o.WithAll(true),
				key:       paramAll.Value(),
				wantValue: true,
			},
			"WithAll-false": {
				option:    o.WithAll(false),
				key:       paramAll.Value(),
				wantValue: false,
			},
			"WithArchived-true": {
				option:    o.WithArchived(true),
				key:       paramArchived.Value(),
				wantValue: true,
			},
			"WithArchived-false": {
				option:    o.WithArchived(false),
				key:       paramArchived.Value(),
				wantValue: false,
			},
			"WithChartEnabled-true": {
				option:    o.WithChartEnabled(true),
				key:       paramChartEnabled.Value(),
				wantValue: true,
			},
			"WithChartEnabled-false": {
				option:    o.WithChartEnabled(false),
				key:       paramChartEnabled.Value(),
				wantValue: false,
			},
			"WithMailNotify-true": {
				option:    o.WithMailNotify(true),
				key:       paramMailNotify.Value(),
				wantValue: true,
			},
			"WithMailNotify-false": {
				option:    o.WithMailNotify(false),
				key:       paramMailNotify.Value(),
				wantValue: false,
			},
			"WithProjectLeaderCanEditProjectLeader-true": {
				option:    o.WithProjectLeaderCanEditProjectLeader(true),
				key:       paramProjectLeaderCanEditProjectLeader.Value(),
				wantValue: true,
			},
			"WithProjectLeaderCanEditProjectLeader-false": {
				option:    o.WithProjectLeaderCanEditProjectLeader(false),
				key:       paramProjectLeaderCanEditProjectLeader.Value(),
				wantValue: false,
			},
			"WithSendMail-true": {
				option:    o.WithSendMail(true),
				key:       paramSendMail.Value(),
				wantValue: true,
			},
			"WithSendMail-false": {
				option:    o.WithSendMail(false),
				key:       paramSendMail.Value(),
				wantValue: false,
			},
			"WithSubtaskingEnabled-true": {
				option:    o.WithSubtaskingEnabled(true),
				key:       paramSubtaskingEnabled.Value(),
				wantValue: true,
			},
			"WithSubtaskingEnabled-false": {
				option:    o.WithSubtaskingEnabled(false),
				key:       paramSubtaskingEnabled.Value(),
				wantValue: false,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := url.Values{}
				err := tc.option.Check()
				require.NoError(t, err)
				_ = tc.option.Set(form)
				assert.Equal(t, strconv.FormatBool(tc.wantValue), form.Get(tc.key))
			})
		}
	})

	// --- Integer options -----------------------------------------------------------
	t.Run("integer-options", func(t *testing.T) {
		cases := map[string]struct {
			option    RequestOption
			key       string
			wantValue int
			wantErr   bool
		}{
			"WithUserID-valid-1": {
				option:    o.WithUserID(1),
				key:       paramUserID.Value(),
				wantValue: 1,
			},
			"WithUserID-valid-2": {
				option:    o.WithUserID(2),
				key:       paramUserID.Value(),
				wantValue: 2,
			},
			"WithUserID-invalid-0": {
				option:  o.WithUserID(0),
				key:     paramUserID.Value(),
				wantErr: true,
			},
			"WithCount-valid-1": {
				option:    o.WithCount(1),
				key:       paramCount.Value(),
				wantValue: 1,
			},
			"WithCount-valid-100": {
				option:    o.WithCount(100),
				key:       paramCount.Value(),
				wantValue: 100,
			},
			"WithCount-invalid-0": {
				option:  o.WithCount(0),
				key:     paramCount.Value(),
				wantErr: true,
			},
			"WithCount-invalid-101": {
				option:  o.WithCount(101),
				key:     paramCount.Value(),
				wantErr: true,
			},
			"WithMinID-valid-1": {
				option:    o.WithMinID(1),
				key:       paramMinID.Value(),
				wantValue: 1,
			},
			"WithMinID-invalid-0": {
				option:  o.WithMinID(0),
				key:     paramMinID.Value(),
				wantErr: true,
			},
			"WithMaxID-valid-26": {
				option:    o.WithMaxID(26),
				key:       paramMaxID.Value(),
				wantValue: 26,
			},
			"WithMaxID-invalid-27": {
				option:  o.WithMaxID(27),
				key:     paramMaxID.Value(),
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
	})

	// --- String options ------------------------------------------------------------
	t.Run("string-options", func(t *testing.T) {
		cases := map[string]struct {
			option    RequestOption
			key       string
			wantValue string
			wantErr   bool
		}{
			"WithContent-valid": {
				option:    o.WithContent("Hello"),
				key:       paramContent.Value(),
				wantValue: "Hello",
			},
			"WithContent-empty": {
				option:  o.WithContent(""),
				key:     paramContent.Value(),
				wantErr: true,
			},
			"WithKey-valid": {
				option:    o.WithKey("ABC"),
				key:       paramKey.Value(),
				wantValue: "ABC",
			},
			"WithKey-empty": {
				option:  o.WithKey(""),
				key:     paramKey.Value(),
				wantErr: true,
			},
			"WithKeyword-non-empty": {
				option:    o.WithKeyword("backlog"),
				key:       paramKeyword.Value(),
				wantValue: "backlog",
			},
			"WithKeyword-empty": {
				option:    o.WithKeyword(""),
				key:       paramKeyword.Value(),
				wantValue: "",
			},
			"WithMailAddress-valid": {
				option:    o.WithMailAddress("test@example.com"),
				key:       paramMailAddress.Value(),
				wantValue: "test@example.com",
			},
			"WithMailAddress-empty": {
				option:  o.WithMailAddress(""),
				key:     paramMailAddress.Value(),
				wantErr: true,
			},
			"WithName-valid": {
				option:    o.WithName("testname"),
				key:       paramName.Value(),
				wantValue: "testname",
			},
			"WithName-empty": {
				option:  o.WithName(""),
				key:     paramName.Value(),
				wantErr: true,
			},
			"WithPassword-valid-8chars": {
				option:    o.WithPassword("abcdefgh"),
				key:       paramPassword.Value(),
				wantValue: "abcdefgh",
			},
			"WithPassword-valid-9chars": {
				option:    o.WithPassword("abcdefghi"),
				key:       paramPassword.Value(),
				wantValue: "abcdefghi",
			},
			"WithPassword-valid-7chars": {
				option:  o.WithPassword("abcdefg"),
				key:     paramPassword.Value(),
				wantErr: true,
			},
			"WithPassword-invalid-empty": {
				option:  o.WithPassword(""),
				key:     paramPassword.Value(),
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
				assert.Equal(t, tc.wantValue, form.Get(tc.key))
			})
		}
	})

	// --- Enum or special options ---------------------------------------------------
	t.Run("enum-or-special-options", func(t *testing.T) {
		cases := map[string]struct {
			option    RequestOption
			key       string
			wantValue string
			wantErr   bool
		}{
			"WithOrder-asc": {
				option:    o.WithOrder(OrderAsc),
				key:       paramOrder.Value(),
				wantValue: "asc",
			},
			"WithOrder-desc": {
				option:    o.WithOrder(OrderDesc),
				key:       paramOrder.Value(),
				wantValue: "desc",
			},
			"WithOrder-invalid": {
				option:  o.WithOrder("invalid"),
				key:     paramOrder.Value(),
				wantErr: true,
			},
			"WithOrder-empty": {
				option:  o.WithOrder(""),
				key:     paramOrder.Value(),
				wantErr: true,
			},
			"WithActivityTypeIDs-single-min": {
				option:    o.WithActivityTypeIDs([]int{1}),
				key:       paramActivityTypeIDs.Value(),
				wantValue: "1",
			},
			"WithActivityTypeIDs-single-max": {
				option:    o.WithActivityTypeIDs([]int{26}),
				key:       paramActivityTypeIDs.Value(),
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
				key: paramActivityTypeIDs.Value(),
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
				key:     paramActivityTypeIDs.Value(),
				wantErr: true,
			},
			"WithActivityTypeIDs-invalid-above": {
				option:  o.WithActivityTypeIDs([]int{27}),
				key:     paramActivityTypeIDs.Value(),
				wantErr: true,
			},
			"WithActivityTypeIDs-invalid-mixed-low": {
				option:  o.WithActivityTypeIDs([]int{0, 1}),
				key:     paramActivityTypeIDs.Value(),
				wantErr: true,
			},
			"WithActivityTypeIDs-invalid-mixed-high": {
				option:  o.WithActivityTypeIDs([]int{26, 27}),
				key:     paramActivityTypeIDs.Value(),
				wantErr: true,
			},
			"WithRoleType-valid-1": {
				option:    o.WithRoleType(1),
				key:       paramRoleType.Value(),
				wantValue: "1",
			},
			"WithRoleType-valid-6": {
				option:    o.WithRoleType(6),
				key:       paramRoleType.Value(),
				wantValue: "6",
			},
			"WithRoleType-invalid-0": {
				option:  o.WithRoleType(0),
				key:     paramRoleType.Value(),
				wantErr: true,
			},
			"WithRoleType-invalid-7": {
				option:  o.WithRoleType(7),
				key:     paramRoleType.Value(),
				wantErr: true,
			},
			"WithTextFormattingRule-valid-backlog": {
				option:    o.WithTextFormattingRule(FormatBacklog),
				key:       paramTextFormattingRule.Value(),
				wantValue: string(FormatBacklog),
			},
			"WithTextFormattingRule-valid-markdown": {
				option:    o.WithTextFormattingRule(FormatMarkdown),
				key:       paramTextFormattingRule.Value(),
				wantValue: string(FormatMarkdown),
			},
			"WithTextFormattingRule-invalid": {
				option:  o.WithTextFormattingRule("invalid"),
				key:     paramTextFormattingRule.Value(),
				wantErr: true,
			},
			"WithTextFormattingRule-invalid-empty": {
				option:  o.WithTextFormattingRule(""),
				key:     paramTextFormattingRule.Value(),
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
					assert.Empty(t, q.Get(tc.key))
					return
				}
				require.NoError(t, err)
				_ = tc.option.Set(q)

				if tc.key == paramActivityTypeIDs.Value() {
					values := (q)[tc.key]
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
