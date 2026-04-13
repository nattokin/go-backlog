package core_test

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

//
// ──────────────────────────────────────────────────────────────
//  TestOptionService
// ──────────────────────────────────────────────────────────────
//

func TestOptionService(t *testing.T) {
	o := &core.OptionService{}

	// --- Boolean options -----------------------------------------------------------
	t.Run("boolean-options", func(t *testing.T) {
		cases := map[string]struct {
			option    core.RequestOption
			key       string
			wantValue bool
		}{
			"WithAll-true": {
				option:    o.WithAll(true),
				key:       core.ParamAll.Value(),
				wantValue: true,
			},
			"WithAll-false": {
				option:    o.WithAll(false),
				key:       core.ParamAll.Value(),
				wantValue: false,
			},
			"WithArchived-true": {
				option:    o.WithArchived(true),
				key:       core.ParamArchived.Value(),
				wantValue: true,
			},
			"WithArchived-false": {
				option:    o.WithArchived(false),
				key:       core.ParamArchived.Value(),
				wantValue: false,
			},
			"WithChartEnabled-true": {
				option:    o.WithChartEnabled(true),
				key:       core.ParamChartEnabled.Value(),
				wantValue: true,
			},
			"WithChartEnabled-false": {
				option:    o.WithChartEnabled(false),
				key:       core.ParamChartEnabled.Value(),
				wantValue: false,
			},
			"WithMailNotify-true": {
				option:    o.WithMailNotify(true),
				key:       core.ParamMailNotify.Value(),
				wantValue: true,
			},
			"WithMailNotify-false": {
				option:    o.WithMailNotify(false),
				key:       core.ParamMailNotify.Value(),
				wantValue: false,
			},
			"WithProjectLeaderCanEditProjectLeader-true": {
				option:    o.WithProjectLeaderCanEditProjectLeader(true),
				key:       core.ParamProjectLeaderCanEditProjectLeader.Value(),
				wantValue: true,
			},
			"WithProjectLeaderCanEditProjectLeader-false": {
				option:    o.WithProjectLeaderCanEditProjectLeader(false),
				key:       core.ParamProjectLeaderCanEditProjectLeader.Value(),
				wantValue: false,
			},
			"WithSendMail-true": {
				option:    o.WithSendMail(true),
				key:       core.ParamSendMail.Value(),
				wantValue: true,
			},
			"WithSendMail-false": {
				option:    o.WithSendMail(false),
				key:       core.ParamSendMail.Value(),
				wantValue: false,
			},
			"WithSubtaskingEnabled-true": {
				option:    o.WithSubtaskingEnabled(true),
				key:       core.ParamSubtaskingEnabled.Value(),
				wantValue: true,
			},
			"WithSubtaskingEnabled-false": {
				option:    o.WithSubtaskingEnabled(false),
				key:       core.ParamSubtaskingEnabled.Value(),
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
			option    core.RequestOption
			key       string
			wantValue int
			wantErr   bool
		}{
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
				key:     core.ParamUserID.Value(),
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
				key:     core.ParamCount.Value(),
				wantErr: true,
			},
			"WithCount-invalid-101": {
				option:  o.WithCount(101),
				key:     core.ParamCount.Value(),
				wantErr: true,
			},
			"WithMinID-valid-1": {
				option:    o.WithMinID(1),
				key:       core.ParamMinID.Value(),
				wantValue: 1,
			},
			"WithMinID-invalid-0": {
				option:  o.WithMinID(0),
				key:     core.ParamMinID.Value(),
				wantErr: true,
			},
			"WithMaxID-valid-26": {
				option:    o.WithMaxID(26),
				key:       core.ParamMaxID.Value(),
				wantValue: 26,
			},
			"WithMaxID-invalid-27": {
				option:  o.WithMaxID(27),
				key:     core.ParamMaxID.Value(),
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
			option    core.RequestOption
			key       string
			wantValue string
			wantErr   bool
		}{
			"WithContent-valid": {
				option:    o.WithContent("Hello"),
				key:       core.ParamContent.Value(),
				wantValue: "Hello",
			},
			"WithContent-empty": {
				option:  o.WithContent(""),
				key:     core.ParamContent.Value(),
				wantErr: true,
			},
			"WithKey-valid": {
				option:    o.WithKey("ABC"),
				key:       core.ParamKey.Value(),
				wantValue: "ABC",
			},
			"WithKey-empty": {
				option:  o.WithKey(""),
				key:     core.ParamKey.Value(),
				wantErr: true,
			},
			"WithKeyword-non-empty": {
				option:    o.WithKeyword("backlog"),
				key:       core.ParamKeyword.Value(),
				wantValue: "backlog",
			},
			"WithKeyword-empty": {
				option:    o.WithKeyword(""),
				key:       core.ParamKeyword.Value(),
				wantValue: "",
			},
			"WithMailAddress-valid": {
				option:    o.WithMailAddress("test@example.com"),
				key:       core.ParamMailAddress.Value(),
				wantValue: "test@example.com",
			},
			"WithMailAddress-empty": {
				option:  o.WithMailAddress(""),
				key:     core.ParamMailAddress.Value(),
				wantErr: true,
			},
			"WithName-valid": {
				option:    o.WithName("testname"),
				key:       core.ParamName.Value(),
				wantValue: "testname",
			},
			"WithName-empty": {
				option:  o.WithName(""),
				key:     core.ParamName.Value(),
				wantErr: true,
			},
			"WithPassword-valid-8chars": {
				option:    o.WithPassword("abcdefgh"),
				key:       core.ParamPassword.Value(),
				wantValue: "abcdefgh",
			},
			"WithPassword-valid-9chars": {
				option:    o.WithPassword("abcdefghi"),
				key:       core.ParamPassword.Value(),
				wantValue: "abcdefghi",
			},
			"WithPassword-valid-7chars": {
				option:  o.WithPassword("abcdefg"),
				key:     core.ParamPassword.Value(),
				wantErr: true,
			},
			"WithPassword-invalid-empty": {
				option:  o.WithPassword(""),
				key:     core.ParamPassword.Value(),
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
			option    core.RequestOption
			key       string
			wantValue string
			wantErr   bool
		}{
			"WithOrder-asc": {
				option:    o.WithOrder(model.OrderAsc),
				key:       core.ParamOrder.Value(),
				wantValue: "asc",
			},
			"WithOrder-desc": {
				option:    o.WithOrder(model.OrderDesc),
				key:       core.ParamOrder.Value(),
				wantValue: "desc",
			},
			"WithOrder-invalid": {
				option:  o.WithOrder("invalid"),
				key:     core.ParamOrder.Value(),
				wantErr: true,
			},
			"WithOrder-empty": {
				option:  o.WithOrder(""),
				key:     core.ParamOrder.Value(),
				wantErr: true,
			},
			"WithActivityTypeIDs-single-min": {
				option:    o.WithActivityTypeIDs([]int{1}),
				key:       core.ParamActivityTypeIDs.Value(),
				wantValue: "1",
			},
			"WithActivityTypeIDs-single-max": {
				option:    o.WithActivityTypeIDs([]int{26}),
				key:       core.ParamActivityTypeIDs.Value(),
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
				key: core.ParamActivityTypeIDs.Value(),
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
				key:     core.ParamActivityTypeIDs.Value(),
				wantErr: true,
			},
			"WithActivityTypeIDs-invalid-above": {
				option:  o.WithActivityTypeIDs([]int{27}),
				key:     core.ParamActivityTypeIDs.Value(),
				wantErr: true,
			},
			"WithActivityTypeIDs-invalid-mixed-low": {
				option:  o.WithActivityTypeIDs([]int{0, 1}),
				key:     core.ParamActivityTypeIDs.Value(),
				wantErr: true,
			},
			"WithActivityTypeIDs-invalid-mixed-high": {
				option:  o.WithActivityTypeIDs([]int{26, 27}),
				key:     core.ParamActivityTypeIDs.Value(),
				wantErr: true,
			},
			"WithRoleType-valid-1": {
				option:    o.WithRoleType(1),
				key:       core.ParamRoleType.Value(),
				wantValue: "1",
			},
			"WithRoleType-valid-6": {
				option:    o.WithRoleType(6),
				key:       core.ParamRoleType.Value(),
				wantValue: "6",
			},
			"WithRoleType-invalid-0": {
				option:  o.WithRoleType(0),
				key:     core.ParamRoleType.Value(),
				wantErr: true,
			},
			"WithRoleType-invalid-7": {
				option:  o.WithRoleType(7),
				key:     core.ParamRoleType.Value(),
				wantErr: true,
			},
			"WithTextFormattingRule-valid-backlog": {
				option:    o.WithTextFormattingRule(model.FormatBacklog),
				key:       core.ParamTextFormattingRule.Value(),
				wantValue: string(model.FormatBacklog),
			},
			"WithTextFormattingRule-valid-markdown": {
				option:    o.WithTextFormattingRule(model.FormatMarkdown),
				key:       core.ParamTextFormattingRule.Value(),
				wantValue: string(model.FormatMarkdown),
			},
			"WithTextFormattingRule-invalid": {
				option:  o.WithTextFormattingRule("invalid"),
				key:     core.ParamTextFormattingRule.Value(),
				wantErr: true,
			},
			"WithTextFormattingRule-invalid-empty": {
				option:  o.WithTextFormattingRule(""),
				key:     core.ParamTextFormattingRule.Value(),
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

				if tc.key == core.ParamActivityTypeIDs.Value() {
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
