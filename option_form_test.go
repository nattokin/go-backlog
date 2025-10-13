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
//  TestFormOption
// ──────────────────────────────────────────────────────────────
//

func TestFormOption(t *testing.T) {
	cases := map[string]struct {
		option           *FormOption
		expectCheckErr   bool
		expectSetErr     bool
		wantValue        string
		wantCheckErrType error
		wantSetErrType   error
	}{
		"Success": {
			option: &FormOption{
				t:         formKey,
				checkFunc: func() error { return nil },
				setFunc: func(form *FormParams) error {
					form.Set(formKey.Value(), "success")
					return nil
				},
			},
			expectCheckErr: false,
			expectSetErr:   false,
			wantValue:      "success",
		},
		"Check-error": {
			option:           newFormOptionWithCheckError(formKey),
			expectCheckErr:   true,
			expectSetErr:     false,
			wantCheckErrType: errors.New("check error"),
		},
		"set-error": {
			option:         newFormOptionWithSetError(formKey),
			expectCheckErr: false,
			expectSetErr:   true,
			wantSetErrType: errors.New("set error"),
		},
		"queryType-invalid": {
			option: &FormOption{
				t: 0,
				setFunc: func(form *FormParams) error {
					return nil
				},
			},
			expectCheckErr:   false,
			expectSetErr:     false,
			wantValue:        "",
			wantCheckErrType: nil,
		},
		"checkFunc-nil": {
			option: &FormOption{
				t:         formKey,
				checkFunc: nil,
				setFunc: func(form *FormParams) error {
					form.Set(formKey.Value(), "checkFunc nil")
					return nil
				},
			},
			expectCheckErr: false,
			expectSetErr:   false,
			wantValue:      "checkFunc nil",
		},
		"set-nil": {
			option: &FormOption{
				t:         formKey,
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

			form := NewFormParams()
			err := tc.option.Check()
			if tc.expectCheckErr {
				require.Error(t, err)
				assert.IsType(t, tc.wantCheckErrType, err)
				return
			}
			require.NoError(t, err)

			if err := tc.option.set(form); tc.expectSetErr {
				require.Error(t, err)
				assert.IsType(t, tc.wantSetErrType, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantValue, form.Get(tc.option.t.Value()))
			}
		})
	}

}

// ──────────────────────────────────────────────────────────────
//	TestFormOptionService
// ──────────────────────────────────────────────────────────────

// TestFormOptionService verifies that each method of FormOptionService
// correctly applies its expected value to FormParams.
//
// The test is organized into logical groups: Boolean, Integer, String,
// and Enum/Special options.
func TestFormOptionService(t *testing.T) {
	o := newFormOptionService()

	// --- Boolean options ------------------------------------------------------------
	t.Run("Boolean-options", func(t *testing.T) {
		cases := map[string]struct {
			option    *FormOption
			key       string
			wantValue bool
		}{
			"WithArchived-true": {
				option:    o.WithArchived(true),
				key:       formArchived.Value(),
				wantValue: true,
			},
			"WithArchived-false": {
				option:    o.WithArchived(false),
				key:       formArchived.Value(),
				wantValue: false,
			},
			"WithChartEnabled-true": {
				option:    o.WithChartEnabled(true),
				key:       formChartEnabled.Value(),
				wantValue: true,
			},
			"WithChartEnabled-false": {
				option:    o.WithChartEnabled(false),
				key:       formChartEnabled.Value(),
				wantValue: false,
			},
			"WithMailNotify-true": {
				option:    o.WithMailNotify(true),
				key:       formMailNotify.Value(),
				wantValue: true,
			},
			"WithMailNotify-false": {
				option:    o.WithMailNotify(false),
				key:       formMailNotify.Value(),
				wantValue: false,
			},
			"WithProjectLeaderCanEditProjectLeader-true": {
				option:    o.WithProjectLeaderCanEditProjectLeader(true),
				key:       formProjectLeaderCanEditProjectLeader.Value(),
				wantValue: true,
			},
			"WithProjectLeaderCanEditProjectLeader-false": {
				option:    o.WithProjectLeaderCanEditProjectLeader(false),
				key:       formProjectLeaderCanEditProjectLeader.Value(),
				wantValue: false,
			},
			"WithSendMail-true": {
				option:    o.WithSendMail(true),
				key:       formSendMail.Value(),
				wantValue: true,
			},
			"WithSendMail-false": {
				option:    o.WithSendMail(false),
				key:       formSendMail.Value(),
				wantValue: false,
			},
			"WithSubtaskingEnabled-true": {
				option:    o.WithSubtaskingEnabled(true),
				key:       formSubtaskingEnabled.Value(),
				wantValue: true,
			},
			"WithSubtaskingEnabled-false": {
				option:    o.WithSubtaskingEnabled(false),
				key:       formSubtaskingEnabled.Value(),
				wantValue: false,
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := NewFormParams()
				err := tc.option.Check()
				require.NoError(t, err)
				_ = tc.option.set(form)
				assert.Equal(t, strconv.FormatBool(tc.wantValue), form.Get(tc.key))
			})
		}
	})

	// --- Integer options ------------------------------------------------------------
	t.Run("Integer-options", func(t *testing.T) {
		cases := map[string]struct {
			option    *FormOption
			key       string
			wantValue int
			wantErr   bool
		}{
			"WithUserID-valid": {
				option:    o.WithUserID(42),
				key:       formUserID.Value(),
				wantValue: 42,
				wantErr:   false,
			},
			"WithUserID-invalid": {
				option:  o.WithUserID(0),
				key:     formUserID.Value(),
				wantErr: true,
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := NewFormParams()
				err := tc.option.Check()
				if tc.wantErr {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				_ = tc.option.set(form)
				assert.Equal(t, strconv.Itoa(tc.wantValue), form.Get(tc.key))
			})
		}
	})

	// --- String options ------------------------------------------------------------
	t.Run("String-options", func(t *testing.T) {
		cases := map[string]struct {
			option    *FormOption
			key       string
			wantValue string
			wantErr   bool
		}{
			"WithContent-valid": {
				option:    o.WithContent("Hello"),
				key:       formContent.Value(),
				wantValue: "Hello",
			},
			"WithContent-empty": {
				option:  o.WithContent(""),
				key:     formContent.Value(),
				wantErr: true,
			},
			"WithKey-valid": {
				option:    o.WithKey("ABC"),
				key:       formKey.Value(),
				wantValue: "ABC",
			},
			"WithKey-empty": {
				option:  o.WithKey(""),
				key:     formKey.Value(),
				wantErr: true,
			},
			"WithMailAddress-valid": {
				option:    o.WithMailAddress("test@example.com"),
				key:       formMailAddress.Value(),
				wantValue: "test@example.com",
			},
			"WithMailAddress-empty": {
				option:  o.WithMailAddress(""),
				key:     formMailAddress.Value(),
				wantErr: true,
			},
			"WithName-valid": {
				option:    o.WithName("testname"),
				key:       formName.Value(),
				wantValue: "testname",
			},
			"WithName-empty": {
				option:  o.WithName(""),
				key:     formName.Value(),
				wantErr: true,
			},
			"WithPassword-valid": {
				option:    o.WithPassword("abcdefgh"),
				key:       formPassword.Value(),
				wantValue: "abcdefgh",
			},
			"WithPassword-short": {
				option:  o.WithPassword("short"),
				key:     formPassword.Value(),
				wantErr: true,
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := NewFormParams()
				err := tc.option.Check()
				if tc.wantErr {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				_ = tc.option.set(form)
				assert.Equal(t, tc.wantValue, form.Get(tc.key))
			})
		}
	})

	// --- Enum or special options ------------------------------------------------------------
	t.Run("Enum-or-special-options", func(t *testing.T) {
		cases := map[string]struct {
			option    *FormOption
			key       string
			wantValue string
			wantErr   bool
		}{
			"WithRoleType-valid": {
				option:    o.WithRoleType(RoleAdministrator),
				key:       formRoleType.Value(),
				wantValue: strconv.Itoa(int(RoleAdministrator)),
			},
			"WithRoleType-invalid": {
				option:  o.WithRoleType(99),
				key:     formRoleType.Value(),
				wantErr: true,
			},
			"WithTextFormattingRule-valid-backlog": {
				option:    o.WithTextFormattingRule(FormatBacklog),
				key:       formTextFormattingRule.Value(),
				wantValue: string(FormatBacklog),
			},
			"WithTextFormattingRule-valid-markdown": {
				option:    o.WithTextFormattingRule(FormatMarkdown),
				key:       formTextFormattingRule.Value(),
				wantValue: string(FormatMarkdown),
			},
			"WithTextFormattingRule-invalid": {
				option:  o.WithTextFormattingRule("invalid"),
				key:     formTextFormattingRule.Value(),
				wantErr: true,
			},
		}

		for name, tc := range cases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				form := NewFormParams()
				err := tc.option.Check()
				if tc.wantErr {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				_ = tc.option.set(form)
				assert.Equal(t, tc.wantValue, form.Get(tc.key))
			})
		}
	})
}
